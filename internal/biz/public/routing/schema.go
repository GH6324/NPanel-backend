package routing

type Envelope struct {
	Schema                 string                 `json:"schema"`
	RoutingVersion         int                    `json:"routing_version"`
	RoutingHash            string                 `json:"routing_hash"`
	GeneratedAt            string                 `json:"generated_at"`
	ExpiresAt              string                 `json:"expires_at,omitempty"`
	Mode                   string                 `json:"mode"`
	ManagedBy              string                 `json:"managed_by"`
	AllowUserOverride      bool                   `json:"allow_user_override"`
	Profile                Profile                `json:"profile"`
	CapabilityRequirements CapabilityRequirements `json:"capability_requirements"`
	DNSResolvers           []DNSResolver          `json:"dns_resolvers"`
	Outbounds              []RouteOutbound        `json:"outbounds"`
	UnlockServices         []UnlockService        `json:"unlock_services"`
	Rules                  []Rule                 `json:"rules"`
	HealthSnapshot         HealthSnapshot         `json:"health_snapshot,omitempty"`
	Compat                 Compat                 `json:"compat,omitempty"`
}

type Profile struct {
	ID                    string       `json:"id"`
	Code                  string       `json:"code"`
	Name                  string       `json:"name"`
	Description           string       `json:"description"`
	Scope                 ProfileScope `json:"scope"`
	Priority              int          `json:"priority"`
	DefaultAction         RouteAction  `json:"default_action"`
	DefaultDNSResolverTag string       `json:"default_dns_resolver_tag"`
	DefaultFallbackPolicy string       `json:"default_fallback_policy"`
	Enabled               bool         `json:"enabled"`
}

type ProfileScope struct {
	Type string `json:"type"`
	ID   string `json:"id"`
}

type CapabilityRequirements struct {
	MinSchema            string   `json:"min_schema"`
	MinOmnxtSDK          string   `json:"min_omnxt_sdk"`
	MinFlutterPlugin     string   `json:"min_flutter_plugin"`
	RequiredFeatures     []string `json:"required_features"`
	OptionalFeatures     []string `json:"optional_features"`
	NodeRequiredFeatures []string `json:"node_required_features"`
}

type DNSResolver struct {
	Tag             string          `json:"tag"`
	Name            string          `json:"name"`
	Proto           string          `json:"proto"`
	Address         string          `json:"address"`
	Port            int             `json:"port"`
	Path            string          `json:"path,omitempty"`
	ServerName      string          `json:"server_name,omitempty"`
	Bootstrap       []string        `json:"bootstrap"`
	Detour          ResolverDetour  `json:"detour"`
	TimeoutMS       int             `json:"timeout_ms"`
	Retry           int             `json:"retry"`
	CacheTTLSeconds int             `json:"cache_ttl_seconds"`
	HealthCheck     HealthCheckSpec `json:"health_check"`
	Enabled         bool            `json:"enabled"`
}

type ResolverDetour struct {
	Type string `json:"type"`
	Tag  string `json:"tag,omitempty"`
}

type HealthCheckSpec struct {
	Enabled         bool   `json:"enabled"`
	Domain          string `json:"domain,omitempty"`
	URL             string `json:"url,omitempty"`
	IntervalSeconds int    `json:"interval_seconds"`
	TimeoutMS       int    `json:"timeout_ms,omitempty"`
}

type RouteOutbound struct {
	Tag              string          `json:"tag"`
	Name             string          `json:"name"`
	Type             string          `json:"type"`
	NodeID           string          `json:"node_id,omitempty"`
	NodeGroupID      string          `json:"node_group_id,omitempty"`
	Region           string          `json:"region,omitempty"`
	ServiceTags      []string        `json:"service_tags"`
	SelectionPolicy  string          `json:"selection_policy"`
	FailPolicy       string          `json:"fail_policy"`
	FallbackPoolTags []string        `json:"fallback_pool_tags"`
	HealthCheck      HealthCheckSpec `json:"health_check"`
	Enabled          bool            `json:"enabled"`
}

type UnlockService struct {
	Code                  string   `json:"code"`
	Name                  string   `json:"name"`
	Category              string   `json:"category"`
	Regions               []string `json:"regions"`
	DefaultRegion         string   `json:"default_region"`
	DefaultOutboundTag    string   `json:"default_outbound_tag"`
	DefaultDNSResolverTag string   `json:"default_dns_resolver_tag"`
	DefaultFailPolicy     string   `json:"default_fail_policy"`
	HealthCheckURL        string   `json:"health_check_url"`
	Enabled               bool     `json:"enabled"`
}

type Rule struct {
	ID          string      `json:"id"`
	Name        string      `json:"name"`
	Priority    int         `json:"priority"`
	Enabled     bool        `json:"enabled"`
	ServiceCode string      `json:"service_code,omitempty"`
	Matcher     Matcher     `json:"matcher"`
	Action      RouteAction `json:"action"`
}

type Matcher struct {
	Type  string    `json:"type"`
	Value any       `json:"value,omitempty"`
	Op    string    `json:"op,omitempty"`
	Items []Matcher `json:"items,omitempty"`
}

type RouteAction struct {
	Type           string `json:"type"`
	OutboundTag    string `json:"outbound_tag,omitempty"`
	DNSResolverTag string `json:"dns_resolver_tag,omitempty"`
	FailPolicy     string `json:"fail_policy,omitempty"`
	Reason         string `json:"reason,omitempty"`
}

type HealthSnapshot struct {
	GeneratedAt  string         `json:"generated_at"`
	Outbounds    []HealthStatus `json:"outbounds"`
	DNSResolvers []HealthStatus `json:"dns_resolvers"`
	Services     []HealthStatus `json:"services"`
}

type HealthStatus struct {
	Tag                 string `json:"tag,omitempty"`
	Code                string `json:"code,omitempty"`
	Region              string `json:"region,omitempty"`
	Status              string `json:"status"`
	Source              string `json:"source"`
	RTTMS               int    `json:"rtt_ms,omitempty"`
	CheckedAt           string `json:"checked_at"`
	ConsecutiveFailures int    `json:"consecutive_failures"`
	LastError           string `json:"last_error,omitempty"`
	OutboundTag         string `json:"outbound_tag,omitempty"`
	DNSResolverTag      string `json:"dns_resolver_tag,omitempty"`
}

type Compat struct {
	LegacyDNS      []string `json:"legacy_dns"`
	LegacyOutbound []string `json:"legacy_outbound"`
}

type PreviewRequest struct {
	Domain            string   `json:"domain"`
	IP                string   `json:"ip,omitempty"`
	Port              int      `json:"port,omitempty"`
	UserID            string   `json:"user_id,omitempty"`
	SubscribeID       string   `json:"subscribe_id,omitempty"`
	UserSubscribeID   string   `json:"user_subscribe_id,omitempty"`
	SubscribeToken    string   `json:"subscribe_token,omitempty"`
	NodeID            string   `json:"node_id,omitempty"`
	SupportedFeatures []string `json:"supported_features,omitempty"`
}

type PreviewResult struct {
	RoutingHash      string      `json:"routing_hash"`
	Profile          Profile     `json:"profile"`
	ScopeType        string      `json:"scope_type,omitempty"`
	ScopeID          string      `json:"scope_id,omitempty"`
	Matched          bool        `json:"matched"`
	Rule             *Rule       `json:"rule,omitempty"`
	Action           RouteAction `json:"action"`
	DNSResolverTag   string      `json:"dns_resolver_tag"`
	OutboundTag      string      `json:"outbound_tag"`
	FallbackPolicy   string      `json:"fallback_policy"`
	Unsupported      []string    `json:"unsupported_features"`
	EffectiveMode    string      `json:"effective_mode"`
	ExecutionEnabled bool        `json:"execution_enabled"`
}

type ConfigOptions struct {
	UserID            int64
	SubscribeID       int64
	UserSubscribeID   int64
	SubscribeToken    string
	NodeID            int64
	SupportedFeatures []string
	UserAgent         string
}
