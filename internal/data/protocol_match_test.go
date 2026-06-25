package data

import (
	"encoding/json"
	"testing"

	servermodel "github.com/npanel-dev/NPanel-backend/internal/model/server"
)

func TestMatchNodeProtocolConfigPrefersSameTypeAndPort(t *testing.T) {
	protocols := []*servermodel.Protocol{
		{Type: "mx", Port: 443, Enable: true, Transport: "mc1"},
		{Type: "mx", Port: 3389, Enable: true, Transport: "mundordp"},
		{Type: "mx", Port: 3306, Enable: true, Transport: "mundosql"},
	}

	for _, tc := range []struct {
		name      string
		port      uint16
		transport string
	}{
		{name: "mc1", port: 443, transport: "mc1"},
		{name: "rdp", port: 3389, transport: "mundordp"},
		{name: "sql", port: 3306, transport: "mundosql"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			matched, _, _ := matchNodeProtocolConfig(protocols, "mx", tc.port)
			if matched == nil {
				t.Fatal("matched protocol is nil")
			}
			if matched.Transport != tc.transport {
				t.Fatalf("matched transport = %q, want %q", matched.Transport, tc.transport)
			}
		})
	}
}

func TestMatchNodeProtocolConfigHandlesDifferentProtocolsAndPorts(t *testing.T) {
	protocols := []*servermodel.Protocol{
		{Type: "shadowsocks", Port: 443, Enable: true, Cipher: "chacha20-ietf-poly1305"},
		{Type: "vless", Port: 8443, Enable: true, Transport: "ws"},
	}

	ss, _, _ := matchNodeProtocolConfig(protocols, "shadowsocks", 443)
	if ss == nil {
		t.Fatal("shadowsocks:443 matched protocol is nil")
	}
	if ss.Cipher != "chacha20-ietf-poly1305" {
		t.Fatalf("shadowsocks cipher = %q, want chacha20-ietf-poly1305", ss.Cipher)
	}

	vless, _, _ := matchNodeProtocolConfig(protocols, "vless", 8443)
	if vless == nil {
		t.Fatal("vless:8443 matched protocol is nil")
	}
	if vless.Transport != "ws" {
		t.Fatalf("vless transport = %q, want ws", vless.Transport)
	}
}

func TestMatchNodeProtocolConfigRequiresSameTypeAndPortWhenNodeHasPort(t *testing.T) {
	protocols := []*servermodel.Protocol{
		{Type: "mx", Port: 443, Enable: true, Transport: "mc1"},
		{Type: "vless", Port: 8443, Enable: true, Transport: "tcp"},
	}

	matched, firstEnabled, firstAvailable := matchNodeProtocolConfig(protocols, "mx", 9443)
	if matched != nil {
		t.Fatalf("matched = %+v, want nil for missing mx:9443 instance", matched)
	}
	if firstEnabled != protocols[0] || firstAvailable != protocols[0] {
		t.Fatalf("unexpected fallbacks: enabled=%+v available=%+v", firstEnabled, firstAvailable)
	}
}

func TestMatchNodeProtocolConfigFallsBackToSameTypeWhenNodePortIsUnset(t *testing.T) {
	protocols := []*servermodel.Protocol{
		{Type: "mx", Port: 443, Enable: true, Transport: "mc1"},
		{Type: "vless", Port: 8443, Enable: true, Transport: "tcp"},
	}

	matched, _, _ := matchNodeProtocolConfig(protocols, "mx", 0)
	if matched == nil || matched.Transport != "mc1" {
		t.Fatalf("matched = %+v, want mx/mc1 fallback", matched)
	}
}

func TestMatchNodeProtocolConfigReturnsFallbacksWhenTypeMissing(t *testing.T) {
	protocols := []*servermodel.Protocol{
		{Type: "vless", Port: 443, Enable: false},
		{Type: "trojan", Port: 8443, Enable: true},
	}

	matched, firstEnabled, firstAvailable := matchNodeProtocolConfig(protocols, "mx", 443)
	if matched != nil {
		t.Fatalf("matched = %+v, want nil for missing type", matched)
	}
	if firstEnabled != protocols[1] {
		t.Fatalf("firstEnabled = %+v, want trojan", firstEnabled)
	}
	if firstAvailable != protocols[0] {
		t.Fatalf("firstAvailable = %+v, want vless", firstAvailable)
	}
}

func TestLegacyMatchedServerProtocolUsesPortInstance(t *testing.T) {
	protocols := []*servermodel.Protocol{
		{Type: "mx", Port: 443, Enable: true, Transport: "mc1"},
		{Type: "mx", Port: 3389, Enable: true, Transport: "mundordp"},
	}
	raw, err := json.Marshal(protocols)
	if err != nil {
		t.Fatal(err)
	}

	matched := legacyMatchedServerProtocol(string(raw), "mx", 3389)
	if matched == nil {
		t.Fatal("matched protocol is nil")
	}
	if matched.Transport != "mundordp" {
		t.Fatalf("matched transport = %q, want mundordp", matched.Transport)
	}
}
