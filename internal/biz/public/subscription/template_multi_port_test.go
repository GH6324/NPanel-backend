package subscription

import (
	"strings"
	"testing"
)

func TestRenderTemplateKeepsSameServerSameProtocolPortsAsDistinctNodes(t *testing.T) {
	nodes := []*NodeInfo{
		{
			Name:      "local-8080",
			Server:    "127.0.0.1",
			Port:      8080,
			Type:      "shadowsocks",
			Method:    "chacha20-ietf-poly1305",
			ServerKey: "server-key-8080",
		},
		{
			Name:      "local-9090",
			Server:    "127.0.0.1",
			Port:      9090,
			Type:      "shadowsocks",
			Method:    "aes-256-gcm",
			ServerKey: "server-key-9090",
		},
	}

	out, err := RenderTemplate(
		`{{- range .Proxies }}{{ .Server }} {{ .Port }} {{ $.UserInfo.Password }} {{ .Type }} {{ .Method }} {{ .ServerKey }}{{ "\n" }}{{- end }}`,
		"text",
		"site",
		"subscribe",
		nodes,
		&UserSubscribe{UUID: "12346"},
		UserInfo{Password: "12346"},
		nil,
	)
	if err != nil {
		t.Fatal(err)
	}

	lines := strings.Split(strings.TrimSpace(string(out)), "\n")
	if len(lines) != 2 {
		t.Fatalf("rendered line count = %d, want 2: %q", len(lines), string(out))
	}
	if lines[0] != "127.0.0.1 8080 12346 shadowsocks chacha20-ietf-poly1305 server-key-8080" {
		t.Fatalf("first line = %q", lines[0])
	}
	if lines[1] != "127.0.0.1 9090 12346 shadowsocks aes-256-gcm server-key-9090" {
		t.Fatalf("second line = %q", lines[1])
	}
}
