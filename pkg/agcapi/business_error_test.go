package agcapi

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Createitv/agc-cli/pkg/domain"
)

func TestInvokeEndpointBusinessErrors(t *testing.T) {
	for _, tc := range []struct{ name, body, code, message string }{
		{"ret numeric", `{"ret":{"code":204144703,"msg":"permission denied"},"appList":null}`, "204144703", "permission denied"},
		{"ret string", `{"ret":{"code":"204144703","msg":"permission denied"}}`, "204144703", "permission denied"},
		{"testing", `{"rtnCode":123,"rtnDesc":"testing denied"}`, "123", "testing denied"},
		{"PMS", `{"error":{"errorCode":"123","errorMsg":"PMS denied"}}`, "123", "PMS denied"},
		{"flat", `{"code":123,"message":"denied"}`, "123", "denied"},
		{"mixed envelopes", `{"ret":{"code":0},"error":{"errorCode":123,"errorMsg":"denied"}}`, "123", "denied"},
		{"business subcode", `{"rtnCode":0,"businessCode":100,"rtnDesc":"invitation already active"}`, "100", "invitation already active"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { w.Write([]byte(tc.body)) }))
			defer server.Close()
			resp, err := InvokeEndpoint(context.Background(), server.Client(), InvokeRequest{Endpoint: domain.Endpoint{Method: "GET", Path: "/test"}, BaseURL: server.URL})
			var agcErr AGCError
			if !errors.As(err, &agcErr) {
				t.Fatalf("expected AGCError, got %v", err)
			}
			if agcErr.StatusCode != 200 || agcErr.Code != tc.code || agcErr.Message != tc.message {
				t.Fatalf("unexpected error: %#v", agcErr)
			}
			if string(resp.RawBody) != tc.body || string(agcErr.RawBody) != tc.body {
				t.Fatal("response body lost")
			}
		})
	}
}

func TestInvokeEndpointBusinessSuccess(t *testing.T) {
	for _, body := range []string{`{"ret":{"code":0}}`, `{"ret":{"code":"0"}}`, `{"rtnCode":"0"}`, `{"error":{"errorCode":0}}`, `{"code":0}`, `{"data":[]}`, `{"code":null}`, "binary\x00data"} {
		t.Run(body, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { w.Write([]byte(body)) }))
			defer server.Close()
			_, err := InvokeEndpoint(context.Background(), server.Client(), InvokeRequest{Endpoint: domain.Endpoint{Method: "GET", Path: "/test"}, BaseURL: server.URL})
			if err != nil {
				t.Fatalf("unexpected failure: %v", err)
			}
		})
	}
}
