package agcapi

import (
	"context"
	"github.com/Createitv/agc-cli/pkg/domain"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestGameplayResourceUsesDocumentedHostAndSignedForm(t *testing.T) {
	e, _ := domain.EndpointByID("gameplay", "resources")
	body := []byte(`{"method":"external.hms.gs.syncGameResources","appId":"fixture","resources":"[{\"index\":\"offline\",\"name\":\"offline\",\"type\":1,\"url\":\"https://example.test/image.png\",\"updateTime\":1,\"operType\":0}]","ts":1,"sign":"offline-sign"}`)
	dry, err := InvokeEndpoint(context.Background(), nil, InvokeRequest{Endpoint: e, Body: body, DryRun: true})
	if err != nil {
		t.Fatal(err)
	}
	if !strings.HasPrefix(dry.URL, "https://jos-open-api.cloud.huawei.com/") {
		t.Fatalf("wrong documented service host: %s", dry.URL)
	}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "" {
			t.Error("AGC token sent to body-signature service")
		}
		if !strings.HasPrefix(r.Header.Get("Content-Type"), "application/x-www-form-urlencoded") {
			t.Error("resource body was not a form")
		}
		if err := r.ParseForm(); err != nil || r.PostForm.Get("method") != "external.hms.gs.syncGameResources" || !strings.Contains(r.PostForm.Get("resources"), `"index":"offline"`) {
			t.Error("resource JSON string changed during form encoding")
		}
		_, _ = w.Write([]byte(`{"code":0}`))
	}))
	defer server.Close()
	if _, err := InvokeEndpoint(context.Background(), server.Client(), InvokeRequest{Endpoint: e, BaseURL: server.URL, Body: body, AccessToken: "offline-agc-token"}); err != nil {
		t.Fatal(err)
	}
}
