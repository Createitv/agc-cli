// Package managementtest contains independent source-based fixtures for offline
// CLI and REST transport assertions. No fixture is derived from the registry.
package managementtest

import (
	"archive/zip"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

type Fixture struct {
	Method, Path  string
	Params, Query map[string]string
	Body          map[string]any
}

func Fixtures() map[string]Fixture {
	out := map[string]Fixture{}
	put := func(key, method, path string, params, query map[string]string, body map[string]any) {
		out[key] = Fixture{method, path, params, query, body}
	}
	put("cache/execute-command", "POST", "/agc/edge/scs/scs-service/v1/eval", nil, nil, map[string]any{"command": "GET fixture-key"})
	put("cache/connection-info", "GET", "/agc/edge/scs/scs-service/v1/cache-invoke-address", nil, nil, nil)
	put("cache/usage", "POST", "/agc/edge/spms/v1/query/metric/scs-metric", nil, map[string]string{"views": "query_cache_metric_trend", "startTime": "1792962000000", "endTime": "1792965600000"}, nil)
	put("cache/enable-service", "POST", "/agc/edge/cpms/project-service/v1/services/order", nil, map[string]string{"projectId": "offline-project"}, map[string]any{"serviceName": "cloudCache"})
	for _, v := range []struct{ id, method string }{{"set-password", "POST"}, {"reset-password", "PUT"}} {
		put("cache/"+v.id, v.method, "/agc/edge/scs/scs-service/v1/cache-password", nil, nil, map[string]any{"password": "Offline9!fixture"})
	}
	db := "/agc/edge/clouddb/clouddbservice/v1/"
	object := map[string]any{"objectTypeName": "Book", "fields": []any{map[string]any{"fieldName": "id", "fieldType": "Integer", "belongPrimaryKey": true, "notNull": true}}, "indexes": []any{}, "permissions": []any{map[string]any{"role": "Administrator", "rights": []any{"Read", "Upsert", "Delete"}}}}
	put("clouddb/schema-versions", "GET", db+"schemaVersion", nil, nil, nil)
	put("clouddb/get-object-type", "GET", db+"objectType", nil, map[string]string{"schemaVersion": "1", "objectTypeName": "Book"}, nil)
	put("clouddb/create-object-type", "POST", db+"objectType", nil, nil, map[string]any{"objectType": object})
	put("clouddb/update-object-type", "PUT", db+"objectType", nil, nil, map[string]any{"objectType": object, "originObjectType": object})
	put("clouddb/delete-object-type", "DELETE", db+"objectType", nil, nil, map[string]any{"objectTypeName": "Book"})
	put("clouddb/import-object-types", "POST", db+"importObjectType", nil, nil, map[string]any{"schemaVersion": 1, "objectTypes": []any{map[string]any{"objectTypeName": "Book", "fields": object["fields"], "indexes": []any{}}}})
	put("clouddb/export-object-types", "POST", db+"exportObjectType", nil, nil, map[string]any{"schemaVersion": 1, "type": "json"})
	put("clouddb/list-zones", "GET", db+"naturalStore", nil, nil, nil)
	put("clouddb/create-zone", "POST", db+"naturalStore", nil, nil, map[string]any{"cloudDBZoneName": "Offline"})
	put("clouddb/delete-zone", "DELETE", db+"naturalStore", nil, nil, map[string]any{"cloudDBZoneName": "Offline"})
	put("clouddb/process-status", "POST", db+"processPercentage", nil, nil, map[string]any{"processId": "offline-process"})
	fn := "/api/wisefunction/serverless/v1/"
	fp := map[string]string{"functionName": "test-func"}
	ap := map[string]string{"functionName": "test-func", "aliasName": "test-alias"}
	put("functions/create-alias", "POST", fn+"functions/test-func/aliases", fp, nil, map[string]any{"name": "test-alias", "functionVersion": "1", "routingConfig": map[string]any{"1": 100}})
	put("functions/list-aliases", "GET", fn+"functions/test-func/aliases", fp, nil, nil)
	put("functions/get-alias", "GET", fn+"functions/test-func/aliases/test-alias", ap, nil, nil)
	put("functions/delete-alias", "DELETE", fn+"functions/test-func/aliases/test-alias", ap, nil, nil)
	put("functions/update-alias", "PUT", fn+"functions/test-func/aliases/test-alias", ap, nil, map[string]any{"functionVersion": "1", "revisionId": "revision-1", "routingConfig": map[string]any{"1": 80, "2": 20}})
	put("functions/create-http-trigger", "POST", fn+"trigger", nil, nil, map[string]any{"funcId": "urn:offline:test-func:1", "mode": "NO_PATH", "authFlag": true, "authAlgor": "HDA-SYSTEM", "authType": "cloudgw-client", "enableUrlDecode": true})
	cron := map[string]any{"funcId": "urn:offline:test-func:1", "triggerType": "Cron", "triggerMode": "ASync", "triggerSwitch": false, "spec": map[string]any{"schedule": "@every 6m", "payload": map[string]any{"hello": "world"}}}
	put("functions/create-cron-trigger", "POST", fn+"triggers", nil, nil, cron)
	update := map[string]any{"funcId": cron["funcId"], "triggerMode": "ASync", "triggerSwitch": true, "spec": cron["spec"], "revisionId": "revision-1"}
	put("functions/update-cron-trigger", "PUT", fn+"triggers/offline-trigger", map[string]string{"triggerId": "offline-trigger"}, nil, update)
	for _, v := range []struct{ kind, route, list string }{{"http-trigger", "trigger", "http-triggers"}, {"trigger", "triggers", "triggers"}} {
		for _, op := range []struct{ id, method string }{{"get", "GET"}, {"delete", "DELETE"}} {
			put("functions/"+op.id+"-"+v.kind, op.method, fn+v.route+"/offline-trigger", map[string]string{"triggerId": "offline-trigger"}, nil, nil)
		}
		put("functions/list-"+v.list, "GET", fn+v.route+"/list/urn:offline:test-func:1", map[string]string{"funcId": "urn:offline:test-func:1"}, nil, nil)
	}
	put("functions/get-function-policy", "GET", fn+"functions/test-func/policy", fp, map[string]string{"qualifier": "1"}, nil)
	put("functions/create-function-policy", "POST", fn+"functions/test-func/policy", fp, map[string]string{"qualifier": "1"}, map[string]any{"effect": "Allow", "action": "InvokeFunction", "sourceProvider": "Cron", "sourceId": "offline-trigger"})
	put("functions/list-layers", "GET", fn+"layers", nil, nil, nil)
	put("functions/list-layer-versions", "GET", fn+"layers/test-layer/versions", map[string]string{"layerName": "test-layer"}, nil, nil)
	put("functions/get-layer-version", "GET", fn+"layers/test-layer/versions/1", map[string]string{"layerName": "test-layer", "layerVersion": "1"}, nil, nil)
	put("functions/create-layer-version", "POST", fn+"layers/test-layer/versions", map[string]string{"layerName": "test-layer"}, nil, map[string]any{"sharedScope": 0, "compatibleRuntimes": []any{"nodejs14.x"}, "zipFile": "UEsFBgAAAAAAAAAAAAAAAAAAAAAAAA=="})
	sam := map[string]any{"name": "offline-app", "components": []any{map[string]any{"name": "test-func", "type": "Huawei::AGC::CloudFunctions", "properties": map[string]any{"codeUri": "offline-bucket/test.zip", "runtime": "nodejs14.x", "memorySize": 500, "timeout": 55, "handler": "handler.myHandler", "instanceConcurrency": 100}}}}
	samProperties := sam["components"].([]any)[0].(map[string]any)["properties"].(map[string]any)
	samProperties["description"] = "offline SAM fixture"
	samProperties["env"] = []any{map[string]any{"name": "ENV_ONE", "value": "offline-value"}}
	samProperties["events"] = []any{map[string]any{"name": "http-event", "type": "http", "properties": map[string]any{"authType": "cloudgw-client", "enableUrlDecode": false}}, map[string]any{"name": "cron-event", "type": "cron", "properties": map[string]any{"schedule": "@every 6m", "input": `{"test":true}`}}}
	put("templates/search-applications", "GET", "/api/adm/v1/applications/search", nil, map[string]string{"appName": "offline-app"}, nil)
	put("templates/create-application", "POST", "/api/adm/v1/applications", nil, nil, sam)
	put("templates/create-plan", "POST", "/api/adm/v1/applications/offline-app/plans", map[string]string{"applicationId": "offline-app"}, nil, sam)
	put("templates/get-plan", "GET", "/api/adm/v1/applications/plans/offline-plan", map[string]string{"planId": "offline-plan"}, nil, nil)
	put("templates/apply-plan", "PUT", "/api/adm/v1/applications/offline-app/plans/offline-plan", map[string]string{"applicationId": "offline-app", "planId": "offline-plan"}, nil, nil)
	return out
}

func SessionFile(t *testing.T, base string) string {
	t.Helper()
	p := filepath.Join(t.TempDir(), "session.json")
	b, _ := json.Marshal(map[string]any{"uid": "offline-user", "teamId": "offline-team", "productId": "offline-project", "accessToken": "offline-oauth", "authInfo": "offline-auth", "csrfToken": "offline-csrf", "expiresAt": time.Now().Add(time.Hour), "projectBaseUrl": base, "projectEdgeBaseUrl": base, "teamBaseUrl": base})
	var session map[string]any
	json.Unmarshal(b, &session)
	session["projectStorageBaseUrl"] = base
	b, _ = json.Marshal(session)
	if err := os.WriteFile(p, b, 0600); err != nil {
		t.Fatal(err)
	}
	return p
}

func Deployment(t *testing.T) (*httptest.Server, map[string]any) {
	t.Helper()
	file := filepath.Join(t.TempDir(), "function.zip")
	f, err := os.Create(file)
	if err != nil {
		t.Fatal(err)
	}
	z := zip.NewWriter(f)
	entry, _ := z.Create("handler.js")
	entry.Write([]byte("exports.myHandler=()=>1"))
	z.Close()
	f.Close()
	body := Fixtures()["templates/create-application"].Body
	body["storageInstance"] = "offline-bucket"
	body["waitForReady"] = true
	body["components"].([]any)[0].(map[string]any)["properties"].(map[string]any)["codeUri"] = file
	created, uploaded := false, false
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if r.Header.Get("oauth2Token") != "offline-oauth" {
			t.Error("deployment session lost")
		}
		switch {
		case r.Method == "GET" && r.URL.Path == "/api/adm/v1/applications/search":
			if created {
				io.WriteString(w, `[{"id":"application:offline","name":"offline-app","status":{"ready":true}}]`)
			} else {
				io.WriteString(w, `[]`)
			}
		case r.Method == "PUT" && strings.HasPrefix(r.URL.Path, "/v0/offline-bucket/"):
			data, _ := io.ReadAll(r.Body)
			expected, _ := os.ReadFile(file)
			if string(data) != string(expected) {
				t.Error("deployment ZIP changed")
			}
			uploaded = true
			io.WriteString(w, `{}`)
		case r.Method == "POST" && r.URL.Path == "/api/adm/v1/applications":
			if !uploaded {
				t.Error("application before upload")
			}
			created = true
			io.WriteString(w, `{"id":"application:offline"}`)
		default:
			t.Errorf("unexpected deployment request %s %s", r.Method, r.URL)
			http.Error(w, "unexpected", 500)
		}
	}))
	return server, body
}

func HDCFile(t *testing.T) string {
	t.Helper()
	t.Setenv("AGC_OFFLINE_HDC_STUB", "1")
	p, err := os.Executable()
	if err != nil {
		t.Fatal(err)
	}
	return p
}

func AssertRequest(t *testing.T, r *http.Request, f Fixture) {
	t.Helper()
	if f.Path == "/agc/edge/clouddb/clouddbservice/v1/importObjectType" && r.Header.Get("Content-Type") != "application/x-www-form-urlencoded" {
		t.Error("vendor raw JSON-string import content type changed")
	}
	if r.Method != f.Method || r.URL.Path != f.Path {
		t.Errorf("wire=%s %s want %s %s", r.Method, r.URL.Path, f.Method, f.Path)
	}
	for k, v := range f.Query {
		if r.URL.Query().Get(k) != v {
			t.Errorf("query %s lost/changed", k)
		}
	}
	if len(f.Body) > 0 {
		var got any
		if json.NewDecoder(r.Body).Decode(&got) != nil {
			t.Error("missing JSON body")
		}
		want, _ := json.Marshal(f.Body)
		raw, _ := json.Marshal(got)
		if string(raw) != string(want) {
			t.Errorf("JSON wire differs: %s", raw)
		}
	}
	if r.Header.Get("oauth2Token") != "offline-oauth" || r.Header.Get("Authorization") != "" {
		t.Error("wrong authentication branch")
	}
}

// ExitIfHDCStub handles only the explicit offline device command in a child test process.
func ExitIfHDCStub() {
	if os.Getenv("AGC_OFFLINE_HDC_STUB") == "1" && len(os.Args) == 8 && os.Args[1] == "-t" && os.Args[3] == "shell" && os.Args[4] == "cf_prefetch" && os.Args[5] == "getcache" && os.Args[6] == "-m" {
		fmt.Println(strings.Join(os.Args[1:], "\n"))
		os.Exit(0)
	}
}
