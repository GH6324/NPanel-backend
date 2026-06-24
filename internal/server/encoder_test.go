package server

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	adminsystemv1 "github.com/npanel-dev/NPanel-backend/api/admin/system/v1"
)

func TestShouldEmitUnpopulatedFields(t *testing.T) {
	tests := []struct {
		name string
		path string
		want bool
	}{
		{name: "admin subscribe config", path: "/v1/admin/system/subscribe_config", want: true},
		{name: "admin subscribe config with api prefix", path: "/api/v1/admin/system/subscribe_config", want: true},
		{name: "public site config", path: "/v1/common/site/config", want: true},
		{name: "other route", path: "/v1/admin/system/site_config", want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := shouldEmitUnpopulatedFields(&http.Request{URL: &url.URL{Path: tt.path}}, nil)
			if got != tt.want {
				t.Fatalf("shouldEmitUnpopulatedFields(%q) = %v, want %v", tt.path, got, tt.want)
			}
		})
	}
}

func TestCustomResponseEncoderKeepsSubscribeConfigFalseFields(t *testing.T) {
	recorder := httptest.NewRecorder()
	req := &http.Request{URL: &url.URL{Path: "/unexpected/proxy/path"}}
	reply := &adminsystemv1.GetSubscribeConfigReply{
		Code:    200,
		Message: "Success",
		Data: &adminsystemv1.SubscribeConfig{
			SubscribePath: "/api/subscribe",
			PanDomain:     false,
		},
	}

	if err := CustomResponseEncoder(recorder, req, reply); err != nil {
		t.Fatalf("CustomResponseEncoder returned error: %v", err)
	}

	body := recorder.Body.String()
	for _, want := range []string{
		`"pan_domain":false`,
		`"subscribe_domain":""`,
		`"single_model":false`,
		`"user_agent_limit":false`,
		`"user_agent_list":""`,
	} {
		if !strings.Contains(body, want) {
			t.Fatalf("response body missing %s: %s", want, body)
		}
	}
}
