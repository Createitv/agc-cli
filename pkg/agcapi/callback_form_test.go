package agcapi

import (
	"context"
	"github.com/Createitv/agc-cli/pkg/domain"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestGameCallbackUsesOfficialFormEncodingWithoutAGCToken(t *testing.T) {
	for _, id := range []string{"propapi-getroleinfo", "propapi-order"} {
		t.Run(id, func(t *testing.T) {
			endpoint, _ := domain.EndpointByID("game-items", id)
			called := false
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				called = true
				if !strings.HasPrefix(r.Header.Get("Content-Type"), "application/x-www-form-urlencoded") {
					t.Errorf("wrong callback content type: %s", r.Header.Get("Content-Type"))
				}
				if r.Header.Get("Authorization") != "" {
					t.Error("AGC token leaked to callback server")
				}
				if err := r.ParseForm(); err != nil {
					t.Fatal(err)
				}
				if r.PostForm.Get("sign") != "a+b/=" || r.PostForm.Get("userId") != "测试+player" || r.PostForm.Get("userIdType") != "1" {
					t.Errorf("callback form lost original values: %v", r.PostForm)
				}
				_, _ = w.Write([]byte(`{"code":0}`))
			}))
			defer server.Close()
			payload := `{"userId":"测试+player","userIdType":1,"ts":123,"sign":"a+b/="}`
			if id == "propapi-order" {
				payload = `{"requestId":"offline-request","userId":"测试+player","userIdType":1,"zoneId":"offline-zone","roleId":"offline-role","productId":"offline-product","price":199,"quantity":1,"currency":"CNY","ts":123,"sign":"a+b/="}`
			}
			_, err := InvokeEndpoint(context.Background(), server.Client(), InvokeRequest{Endpoint: endpoint, Params: map[string]string{"callbackUrl": server.URL}, AccessToken: "offline-agc-token", Body: []byte(payload)})
			if err != nil {
				t.Fatal(err)
			}
			if !called {
				t.Fatal("callback transport was not exercised")
			}
		})
	}
}
