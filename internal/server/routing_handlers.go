package server

import (
	"encoding/json"
	"fmt"
	nethttp "net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/go-kratos/kratos/v2/transport/http"
	publicrouting "github.com/npanel-dev/NPanel-backend/internal/biz/public/routing"
	"github.com/npanel-dev/NPanel-backend/internal/pkg/middleware"
	adminroutingservice "github.com/npanel-dev/NPanel-backend/internal/service/admin/routing"
)

type publicRoutingCache struct {
	mu    sync.Mutex
	items map[string]publicRoutingCacheItem
}

type publicRoutingCacheItem struct {
	envelope  publicrouting.Envelope
	expiresAt time.Time
}

func registerRoutingPreviewRoutes(srv *http.Server, routing *adminroutingservice.RoutingService) {
	cache := &publicRoutingCache{}
	srv.HandleFunc("/v1/public/routing/config", handleRoutingConfig(routing, cache))
	srv.HandleFunc("/v1/public/routing/preview", handleRoutingPreview(routing, cache))
}

func handleRoutingConfig(routing *adminroutingservice.RoutingService, cache *publicRoutingCache) nethttp.HandlerFunc {
	return func(w nethttp.ResponseWriter, r *nethttp.Request) {
		if r.Method != nethttp.MethodGet {
			writeRoutingError(w, nethttp.StatusMethodNotAllowed, 405, "method not allowed")
			return
		}

		cfg, fallback := loadPublicRoutingConfig(routing, cache, r, publicrouting.PreviewRequest{})
		if fallback != "" {
			w.Header().Set("X-Routing-Fallback", fallback)
		}

		if r.Header.Get("If-None-Match") == cfg.RoutingHash {
			w.WriteHeader(nethttp.StatusNotModified)
			return
		}

		writeRoutingOK(w, cfg)
	}
}

func handleRoutingPreview(routing *adminroutingservice.RoutingService, cache *publicRoutingCache) nethttp.HandlerFunc {
	return func(w nethttp.ResponseWriter, r *nethttp.Request) {
		if r.Method != nethttp.MethodPost {
			writeRoutingError(w, nethttp.StatusMethodNotAllowed, 405, "method not allowed")
			return
		}

		var req publicrouting.PreviewRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeRoutingError(w, nethttp.StatusBadRequest, 400, "invalid preview request")
			return
		}
		if len(req.SupportedFeatures) == 0 {
			req.SupportedFeatures = publicrouting.ParseFeatureList(r.Header.Get("X-Routing-Features"))
		}
		req.Domain = strings.TrimSpace(req.Domain)
		if req.Domain == "" && req.IP == "" {
			writeRoutingError(w, nethttp.StatusBadRequest, 422, "domain or ip is required")
			return
		}

		fillPreviewScopeFromRequest(r, &req)
		cfg, fallback := loadPublicRoutingConfig(routing, cache, r, req)
		if fallback != "" {
			w.Header().Set("X-Routing-Fallback", fallback)
		}
		result := publicrouting.PreviewRouteConfig(cfg, req)
		writeRoutingOK(w, result)
	}
}

func loadPublicRoutingConfig(routing *adminroutingservice.RoutingService, cache *publicRoutingCache, r *nethttp.Request, req publicrouting.PreviewRequest) (publicrouting.Envelope, string) {
	now := time.Now()
	scopeOpts := routingConfigOptionsFromRequest(r, req)
	cacheKey := publicRoutingCacheKey(scopeOpts)
	if cfg, ok := cache.get(cacheKey, now); ok {
		return cfg, ""
	}

	features := publicrouting.ParseFeatureList(r.Header.Get("X-Routing-Features"))
	if len(scopeOpts.SupportedFeatures) == 0 {
		scopeOpts.SupportedFeatures = features
	}
	cfg, err := routing.BuildPublicConfig(r.Context(), now, publicrouting.ConfigOptions{
		UserID:            scopeOpts.UserID,
		SubscribeID:       scopeOpts.SubscribeID,
		UserSubscribeID:   scopeOpts.UserSubscribeID,
		SubscribeToken:    scopeOpts.SubscribeToken,
		NodeID:            scopeOpts.NodeID,
		UserAgent:         r.UserAgent(),
		SupportedFeatures: scopeOpts.SupportedFeatures,
	})
	if err != nil {
		return publicrouting.BuildPreviewConfig(now, publicrouting.ConfigOptions{
			UserID:            scopeOpts.UserID,
			SubscribeID:       scopeOpts.SubscribeID,
			UserSubscribeID:   scopeOpts.UserSubscribeID,
			SubscribeToken:    scopeOpts.SubscribeToken,
			NodeID:            scopeOpts.NodeID,
			UserAgent:         r.UserAgent(),
			SupportedFeatures: scopeOpts.SupportedFeatures,
		}), "fixture"
	}
	cache.set(cacheKey, cfg, now.Add(15*time.Second))
	return cfg, ""
}

func routingConfigOptionsFromRequest(r *nethttp.Request, req publicrouting.PreviewRequest) publicrouting.ConfigOptions {
	query := r.URL.Query()
	token := firstNonEmptyString(req.SubscribeToken, r.Header.Get("token"), r.Header.Get("X-Subscribe-Token"), query.Get("subscribe_token"), query.Get("token"))
	userID := firstNonZeroInt64(parseInt64String(req.UserID), parseInt64String(query.Get("user_id")), middleware.GetUserID(r.Context()))
	return publicrouting.ConfigOptions{
		UserID:            userID,
		SubscribeID:       firstNonZeroInt64(parseInt64String(req.SubscribeID), parseInt64String(query.Get("subscribe_id"))),
		UserSubscribeID:   firstNonZeroInt64(parseInt64String(req.UserSubscribeID), parseInt64String(query.Get("user_subscribe_id"))),
		SubscribeToken:    token,
		NodeID:            firstNonZeroInt64(parseInt64String(req.NodeID), parseInt64String(query.Get("node_id"))),
		SupportedFeatures: req.SupportedFeatures,
	}
}

func fillPreviewScopeFromRequest(r *nethttp.Request, req *publicrouting.PreviewRequest) {
	query := r.URL.Query()
	req.SubscribeToken = firstNonEmptyString(req.SubscribeToken, r.Header.Get("token"), r.Header.Get("X-Subscribe-Token"), query.Get("subscribe_token"), query.Get("token"))
	req.UserID = firstNonEmptyString(req.UserID, query.Get("user_id"))
	req.SubscribeID = firstNonEmptyString(req.SubscribeID, query.Get("subscribe_id"))
	req.UserSubscribeID = firstNonEmptyString(req.UserSubscribeID, query.Get("user_subscribe_id"))
	req.NodeID = firstNonEmptyString(req.NodeID, query.Get("node_id"))
}

func publicRoutingCacheKey(opts publicrouting.ConfigOptions) string {
	return fmt.Sprintf("u:%d|s:%d|us:%d|n:%d|t:%s|f:%s", opts.UserID, opts.SubscribeID, opts.UserSubscribeID, opts.NodeID, opts.SubscribeToken, strings.Join(opts.SupportedFeatures, ","))
}

func parseInt64String(value string) int64 {
	parsed, _ := strconv.ParseInt(strings.TrimSpace(value), 10, 64)
	return parsed
}

func firstNonZeroInt64(values ...int64) int64 {
	for _, value := range values {
		if value != 0 {
			return value
		}
	}
	return 0
}

func firstNonEmptyString(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return strings.TrimSpace(value)
		}
	}
	return ""
}

func (c *publicRoutingCache) get(key string, now time.Time) (publicrouting.Envelope, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.items == nil {
		return publicrouting.Envelope{}, false
	}
	item, ok := c.items[key]
	if !ok || now.After(item.expiresAt) {
		return publicrouting.Envelope{}, false
	}
	return item.envelope, true
}

func (c *publicRoutingCache) set(key string, envelope publicrouting.Envelope, expiresAt time.Time) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.items == nil {
		c.items = map[string]publicRoutingCacheItem{}
	}
	c.items[key] = publicRoutingCacheItem{envelope: envelope, expiresAt: expiresAt}
}

func writeRoutingOK(w nethttp.ResponseWriter, data any) {
	writeRoutingJSON(w, nethttp.StatusOK, map[string]any{
		"code":    200,
		"message": "success",
		"data":    data,
	})
}

func writeRoutingError(w nethttp.ResponseWriter, httpStatus, code int, message string) {
	writeRoutingJSON(w, httpStatus, map[string]any{
		"code":    code,
		"message": message,
		"data":    map[string]any{},
	})
}

func writeRoutingJSON(w nethttp.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(body)
}
