package agcapi

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/Createitv/agc-cli/pkg/domain"
	"strings"
	"testing"
)

func TestSchemaUpdatesPreserveImmutableFieldFlagsAndIndexDefinitions(t *testing.T) {
	e, _ := domain.EndpointByID("clouddb", "update-object-type")
	original := map[string]any{"objectTypeName": "Book", "fields": []any{map[string]any{"fieldName": "id", "fieldType": "Integer", "belongPrimaryKey": true}, map[string]any{"fieldName": "secret", "fieldType": "String", "isNeedEncrypt": false, "isSensitive": false}}, "indexes": []any{map[string]any{"indexName": "byId", "indexList": []any{map[string]any{"fieldName": "id", "sortType": "ASC"}}}}}
	for _, test := range []struct {
		name, field string
		change      func(map[string]any)
	}{
		{"encryption", "body.objectType.fields[1].isNeedEncrypt", func(m map[string]any) { m["fields"].([]any)[1].(map[string]any)["isNeedEncrypt"] = true }},
		{"sensitivity", "body.objectType.fields[1].isSensitive", func(m map[string]any) { m["fields"].([]any)[1].(map[string]any)["isSensitive"] = true }},
		{"index", "body.objectType.indexes[0].indexList", func(m map[string]any) {
			m["indexes"].([]any)[0].(map[string]any)["indexList"].([]any)[0].(map[string]any)["sortType"] = "DESC"
		}},
	} {
		t.Run(test.name, func(t *testing.T) {
			raw, _ := json.Marshal(original)
			var changed map[string]any
			json.Unmarshal(raw, &changed)
			test.change(changed)
			b, _ := json.Marshal(map[string]any{"originObjectType": original, "objectType": changed})
			err := ValidateEndpointRequest(e, nil, nil, b)
			var invalid *RequestValidationError
			if !errors.As(err, &invalid) {
				t.Fatalf("immutable mutation accepted: %v", err)
			}
			for _, issue := range invalid.Violations {
				if issue.Field == test.field {
					return
				}
			}
			t.Fatalf("missing precise child path %s", test.field)
		})
	}
}

func TestSAMNestedEnvironmentAndHandlerRulesRejectBeforeAuthorization(t *testing.T) {
	e, _ := domain.EndpointByID("templates", "deploy")
	for _, test := range []struct {
		name string
		prop map[string]any
	}{
		{"handler", map[string]any{"handler": "missing-dot"}},
		{"environment-key", map[string]any{"env": []any{map[string]any{"name": "bad name", "value": "valid"}}}},
		{"environment-value", map[string]any{"env": []any{map[string]any{"name": "GOOD", "value": strings.Repeat("x", 2049)}}}},
		{"http-event-prerequisites", map[string]any{"events": []any{map[string]any{"name": "test-event", "type": "http", "properties": map[string]any{}}}}},
		{"cron-event-json", map[string]any{"events": []any{map[string]any{"name": "test-event", "type": "cron", "properties": map[string]any{"schedule": "@every 6m", "input": "not-json"}}}}},
	} {
		t.Run(test.name, func(t *testing.T) {
			properties := map[string]any{"codeUri": "/offline.zip", "runtime": "nodejs14.x", "memorySize": 500, "timeout": 55, "handler": "handler.myHandler", "instanceConcurrency": 100}
			for k, v := range test.prop {
				properties[k] = v
			}
			body, _ := json.Marshal(map[string]any{"name": "offline", "storageInstance": "offline-bucket", "components": []any{map[string]any{"name": "test-func", "type": "Huawei::AGC::CloudFunctions", "properties": properties}}})
			if ValidateEndpointRequest(e, nil, nil, body) == nil {
				t.Fatal("invalid SAM child field accepted")
			}
		})
	}
}

func TestManagementPathViolationPrecedesSessionLookup(t *testing.T) {
	e, _ := domain.EndpointByID("functions", "get-alias")
	_, err := InvokeEndpoint(context.Background(), nil, InvokeRequest{Endpoint: e, Params: map[string]string{"functionName": "INVALID NAME", "aliasName": "test-alias"}, SessionFile: "/must-not-read"})
	var invalid *RequestValidationError
	if !errors.As(err, &invalid) || len(invalid.Violations) == 0 || invalid.Violations[0].Field != "params.functionName" {
		t.Fatalf("expected indexed path validation before session lookup: %v", err)
	}
}

func TestReviewedManagementRejectsInvalidContent(t *testing.T) {
	cases := []struct{ family, id, body string }{
		{"cache", "execute-command", `{"command":"FLUSHALL"}`},
		{"cache", "execute-command", `{"command":"GET 中文"}`},
		{"cache", "reset-password", `{"password":"alllowercase"}`},
		{"functions", "create-alias", `{"name":"my-alias","functionVersion":"1","routingConfig":{"1":80,"2":80}}`},
		{"functions", "create-layer-version", `{"compatibleRuntimes":["nodejs14.x"],"sharedScope":0,"zipFile":"not-a-zip"}`},
		{"functions", "create-cron-trigger", `{"funcId":"id","triggerType":"Cron","triggerMode":"ASync","triggerSwitch":false,"spec":{"schedule":"@every 0m"}}`},
		{"clouddb", "create-object-type", `{"objectType":{"objectTypeName":"Thing","fields":[{"fieldName":"id","fieldType":"Boolean","belongPrimaryKey":true}]}}`},
	}
	for _, c := range cases {
		t.Run(c.family+"/"+c.id, func(t *testing.T) {
			e, ok := domain.EndpointByID(c.family, c.id)
			if !ok {
				t.Fatal("operation absent")
			}
			if err := ValidateEndpointRequest(e, nil, nil, []byte(c.body)); err == nil {
				t.Fatal("invalid documented content accepted")
			}
		})
	}
}
