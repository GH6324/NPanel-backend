package server

import (
	"encoding/json"
	nethttp "net/http"
	"strings"
	"time"

	"github.com/go-kratos/kratos/v2/transport/http"
	publicrouting "github.com/npanel-dev/NPanel-backend/internal/biz/public/routing"
	"github.com/npanel-dev/NPanel-backend/internal/pkg/middleware"
)

func registerRoutingPreviewRoutes(srv *http.Server) {
	srv.HandleFunc("/v1/public/routing/config", handleRoutingConfig)
	srv.HandleFunc("/v1/public/routing/preview", handleRoutingPreview)
}

func handleRoutingConfig(w nethttp.ResponseWriter, r *nethttp.Request) {
	if r.Method != nethttp.MethodGet {
		writeRoutingError(w, nethttp.StatusMethodNotAllowed, 405, "method not allowed")
		return
	}

	features := publicrouting.ParseFeatureList(r.Header.Get("X-Routing-Features"))
	cfg := publicrouting.BuildPreviewConfig(time.Now(), publicrouting.ConfigOptions{
		UserID:            middleware.GetUserID(r.Context()),
		UserAgent:         r.UserAgent(),
		SupportedFeatures: features,
	})

	if r.Header.Get("If-None-Match") == cfg.RoutingHash {
		w.WriteHeader(nethttp.StatusNotModified)
		return
	}

	writeRoutingOK(w, cfg)
}

func handleRoutingPreview(w nethttp.ResponseWriter, r *nethttp.Request) {
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

	cfg := publicrouting.BuildPreviewConfig(time.Now(), publicrouting.ConfigOptions{
		UserID:            middleware.GetUserID(r.Context()),
		UserAgent:         r.UserAgent(),
		SupportedFeatures: req.SupportedFeatures,
	})
	result := publicrouting.PreviewRouteConfig(cfg, req)
	writeRoutingOK(w, result)
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
