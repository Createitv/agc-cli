package domain

import "testing"

func TestEndpointRegistryCoversEveryCapability(t *testing.T) {
	if got := EndpointCount(); got < 150 {
		t.Fatalf("endpoint count = %d, want at least 150", got)
	}
	for _, capability := range Capabilities {
		endpoints := EndpointsByFamily(capability.ID)
		if len(endpoints) == 0 {
			t.Fatalf("%s has no endpoints", capability.ID)
		}
		for _, endpoint := range endpoints {
			if endpoint.FamilyID != capability.ID {
				t.Fatalf("endpoint family = %q, want %q", endpoint.FamilyID, capability.ID)
			}
			if endpoint.Method == "" || endpoint.Path == "" || endpoint.Command == "" {
				t.Fatalf("endpoint has incomplete interface data: %#v", endpoint)
			}
			if endpoint.SourceURL == "" || endpoint.OfficialSlug == "" {
				t.Fatalf("endpoint has no official source metadata: %#v", endpoint)
			}
			if endpoint.Links["self"].Href == "" {
				t.Fatalf("endpoint %s has no self link", endpoint.ID)
			}
		}
	}
}

func TestEndpointByID(t *testing.T) {
	endpoint, ok := EndpointByID("publishing", "app-submit")
	if !ok {
		t.Fatal("app-submit endpoint not found")
	}
	if endpoint.Method != "POST" {
		t.Fatalf("method = %q, want POST", endpoint.Method)
	}
	if endpoint.Path != "/api/publish/v2/app-submit" {
		t.Fatalf("path = %q", endpoint.Path)
	}
	if _, ok := EndpointByID("publishing", "missing"); ok {
		t.Fatal("missing endpoint unexpectedly found")
	}
}

func TestEndpointParametersAreClassified(t *testing.T) {
	review, ok := EndpointByID("comments", "com-getreviewinfo")
	if !ok {
		t.Fatal("com-getreviewinfo endpoint not found")
	}
	locations := map[string]string{}
	for _, parameter := range review.Parameters {
		locations[parameter.Name] = parameter.In
	}
	if locations["reviewId"] != "path" {
		t.Fatalf("reviewId location = %q, want path", locations["reviewId"])
	}
	callback, ok := EndpointByID("game-items", "propapi-order")
	if !ok {
		t.Fatal("propapi-order endpoint not found")
	}
	if callback.Parameters[0].Name != "callbackUrl" || callback.Parameters[0].In != "path" {
		t.Fatalf("callback parameters = %#v", callback.Parameters)
	}
	update, ok := EndpointByID("pms", "updatepromotion-harmonyosnext")
	if !ok {
		t.Fatal("updatepromotion-harmonyosnext endpoint not found")
	}
	if update.Body == "" {
		t.Fatal("updatepromotion-harmonyosnext should declare a body")
	}
}

func TestVerifiedEndpointParameterContracts(t *testing.T) {
	cases := []struct {
		family, id, name, location string
		required                   bool
	}{
		{"provisioning", "provision-api-get-fingerprints", "appId", "header", true},
		{"provisioning", "provision-api-getacl", "appId", "header", true},
		{"domains", "domain-api-get-domain", "appId", "header", true},
		{"domains", "domain-api-get-domain", "category", "query", true},
		{"testing", "test-api-get-test-grouplist", "appId", "header", true},
		{"testing", "test-api-query-test-user", "appId", "header", true},
		{"testing", "test-api-query-test-user", "groupId", "header", true},
		{"upload", "upload-url-new", "contentLength", "query", true},
		{"upload", "upload-url-new", "appId", "query", false},
		{"upload", "upload-url-new", "suffix", "query", false},
		{"upload", "upload-url-new", "fileName", "query", false},
		{"publishing", "app-info-query", "appId", "query", true},
		{"publishing", "app-info-query", "lang", "query", false},
		{"projects", "queryservice", "projectId", "query", true},
		{"projects", "queryservice", "appID", "query", true},
		{"reports", "appdownloadexport", "language", "query", false},
		{"reports", "appdownloadexport", "startTime", "query", false},
		{"reports", "appdownloadexport", "endTime", "query", false},
		{"reports", "iapexport", "currency", "query", true},
		{"reports", "orderanalysisexport", "teamId", "header", true},
	}
	for _, tc := range cases {
		t.Run(tc.id+"/"+tc.name, func(t *testing.T) {
			ep, ok := EndpointByID(tc.family, tc.id)
			if !ok {
				t.Fatal("missing endpoint")
			}
			for _, p := range ep.Parameters {
				if p.Name == tc.name {
					if p.In != tc.location || p.Required != tc.required {
						t.Fatalf("parameter = %#v", p)
					}
					return
				}
			}
			t.Fatalf("missing parameter %s", tc.name)
		})
	}
}

func TestHarmonyOSCommentsRequireMillisecondQueryWindow(t *testing.T) {
	for _, id := range []string{"comapi-getreviews-harmonyos", "com-rating-harmonyos"} {
		ep, _ := EndpointByID("comments", id)
		found := map[string]Parameter{}
		for _, parameter := range ep.Parameters {
			found[parameter.Name] = parameter
		}
		for _, name := range []string{"appId", "beginTime", "endTime", "countries"} {
			parameter, ok := found[name]
			if !ok || parameter.In != "query" || !parameter.Required {
				t.Fatalf("%s %s = %#v", id, name, parameter)
			}
		}
		for _, name := range []string{"beginTime", "endTime"} {
			if found[name].Type != "integer" || found[name].Format != "int64" {
				t.Fatalf("%s timestamp = %#v", id, found[name])
			}
		}
	}
}
