package data

import "testing"

func TestSelectServerProtocolConfigDifferentProtocolsDifferentPorts(t *testing.T) {
	protocols := []map[string]interface{}{
		{"type": "shadowsocks", "port": float64(443), "cipher": "chacha20-ietf-poly1305"},
		{"type": "vless", "port": float64(8443), "transport": "ws"},
	}

	ss := selectServerProtocolConfig(protocols, "shadowsocks", 443)
	if ss == nil {
		t.Fatal("shadowsocks:443 config is nil")
	}
	if got := ss["cipher"]; got != "chacha20-ietf-poly1305" {
		t.Fatalf("shadowsocks cipher = %v, want chacha20-ietf-poly1305", got)
	}

	vless := selectServerProtocolConfig(protocols, "vless", 8443)
	if vless == nil {
		t.Fatal("vless:8443 config is nil")
	}
	if got := vless["transport"]; got != "ws" {
		t.Fatalf("vless transport = %v, want ws", got)
	}
}

func TestSelectServerProtocolConfigSameProtocolDifferentPorts(t *testing.T) {
	protocols := []map[string]interface{}{
		{"type": "shadowsocks", "port": float64(443), "cipher": "chacha20-ietf-poly1305"},
		{"type": "shadowsocks", "port": float64(8443), "cipher": "aes-256-gcm"},
	}

	matched := selectServerProtocolConfig(protocols, "shadowsocks", 8443)
	if matched == nil {
		t.Fatal("shadowsocks:8443 config is nil")
	}
	if got := matched["cipher"]; got != "aes-256-gcm" {
		t.Fatalf("shadowsocks:8443 cipher = %v, want aes-256-gcm", got)
	}
}

func TestSelectServerProtocolConfigReturnsNilForMissingPortInstance(t *testing.T) {
	protocols := []map[string]interface{}{
		{"type": "shadowsocks", "port": float64(443), "cipher": "chacha20-ietf-poly1305"},
		{"type": "shadowsocks", "port": float64(8443), "cipher": "aes-256-gcm"},
	}

	matched := selectServerProtocolConfig(protocols, "shadowsocks", 9443)
	if matched != nil {
		t.Fatalf("matched = %+v, want nil for missing shadowsocks:9443", matched)
	}
}
