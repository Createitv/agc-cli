package agcapi

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/Createitv/agc-cli/pkg/domain"
)

func TestPrefetchDeviceUsesExplicitDeviceWithoutShellInterpolation(t *testing.T) {
	e, ok := domain.EndpointByID("prefetch", "device-getcache")
	if !ok {
		t.Fatal("documented device command missing")
	}
	p := filepath.Join(t.TempDir(), "hdc")
	// A test-owned executable proves the exact argv; no handset is contacted.
	if err := os.WriteFile(p, []byte("#!/bin/sh\nprintf '%s\\n' \"$@\"\n"), 0700); err != nil {
		t.Fatal(err)
	}
	b, _ := json.Marshal(map[string]string{"bundleName": "com.example.offline", "deviceId": "offline-device", "hdcPath": p})
	got, err := InvokeEndpoint(context.Background(), nil, InvokeRequest{Endpoint: e, Body: b})
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(got.Body), `-t\noffline-device\nshell\ncf_prefetch\ngetcache\n-m\ncom.example.offline\n`) {
		t.Fatalf("wrong hdc argv %s", got.Body)
	}
	if !strings.Contains(string(got.Body), `"dataFreshnessVerified":false`) {
		t.Fatal("command success was misreported as fresh device data")
	}
}

func TestPrefetchDeviceDryRunNeverExecutesHDC(t *testing.T) {
	e, ok := domain.EndpointByID("prefetch", "device-getcache")
	if !ok {
		t.Fatal("device operation missing")
	}
	_, err := InvokeEndpoint(context.Background(), nil, InvokeRequest{Endpoint: e, Body: []byte(`{"bundleName":"com.example.offline","deviceId":"offline-device","hdcPath":"/must-not-execute"}`), DryRun: true})
	if err != nil {
		t.Fatal(err)
	}
}
