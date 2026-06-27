package routing

import (
	"testing"
	"time"
)

func TestPreviewRouteConfigMatchesOpenAI(t *testing.T) {
	cfg := BuildPreviewConfig(time.Unix(1, 0), ConfigOptions{})

	result := PreviewRouteConfig(cfg, PreviewRequest{Domain: "api.openai.com"})

	if !result.Matched {
		t.Fatalf("expected openai domain to match")
	}
	if result.Rule == nil || result.Rule.ID != "rule_openai_suffix" {
		t.Fatalf("unexpected matched rule: %#v", result.Rule)
	}
	if result.Action.Type != "outbound" {
		t.Fatalf("unexpected action type: %s", result.Action.Type)
	}
	if result.OutboundTag != "unlock:openai:us" {
		t.Fatalf("unexpected outbound tag: %s", result.OutboundTag)
	}
	if result.DNSResolverTag != "dns:cloudflare-doh" {
		t.Fatalf("unexpected resolver tag: %s", result.DNSResolverTag)
	}
	if result.FallbackPolicy != "fallback_default" {
		t.Fatalf("unexpected fallback policy: %s", result.FallbackPolicy)
	}
	if result.ExecutionEnabled {
		t.Fatalf("P0 preview must not enable execution")
	}
}

func TestPreviewRouteConfigFallsBackToDefaultProxy(t *testing.T) {
	cfg := BuildPreviewConfig(time.Unix(1, 0), ConfigOptions{})

	result := PreviewRouteConfig(cfg, PreviewRequest{Domain: "example.com"})

	if result.Matched {
		t.Fatalf("did not expect example.com to match")
	}
	if result.Action.Type != "proxy" {
		t.Fatalf("unexpected default action: %s", result.Action.Type)
	}
	if result.OutboundTag != "proxy:default" {
		t.Fatalf("unexpected default outbound: %s", result.OutboundTag)
	}
}

func TestMissingRequiredFeaturesTreatsEmptySupportedAsMissingAll(t *testing.T) {
	missing := MissingRequiredFeatures([]string{"routing_profile_v1", "dot"}, nil)

	if len(missing) != 2 || missing[0] != "routing_profile_v1" || missing[1] != "dot" {
		t.Fatalf("missing = %#v, want all required features", missing)
	}
}

func TestStableHashIgnoresGeneratedAtAndHealth(t *testing.T) {
	first := BuildPreviewConfig(time.Unix(1, 0), ConfigOptions{})
	second := BuildPreviewConfig(time.Unix(100, 0), ConfigOptions{})

	if first.RoutingHash != second.RoutingHash {
		t.Fatalf("hash should be stable across generated_at/health changes: %s != %s", first.RoutingHash, second.RoutingHash)
	}
}
