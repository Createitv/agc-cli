package command

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

type publishingReference struct {
	Endpoints []struct {
		ID         string `json:"id"`
		Method     string `json:"method"`
		Path       string `json:"path"`
		Parameters []struct {
			Name     string `json:"name"`
			In       string `json:"in"`
			Type     string `json:"type"`
			Required bool   `json:"required"`
		} `json:"parameters"`
	} `json:"endpoints"`
}

// This fixture was read from official Huawei parameter tables, independently
// of the production registry. Wrong paths, absent commands and misplaced or
// stringified body values must fail at the actual CLI -> HTTP boundary.
func TestCurrentPublishingRequestsMatchOfficialReference(t *testing.T) {
	raw, err := os.ReadFile("testdata/publishing-current.json")
	if err != nil {
		t.Fatal(err)
	}
	var reference publishingReference
	if err = json.Unmarshal(raw, &reference); err != nil {
		t.Fatal(err)
	}
	if len(reference.Endpoints) != 31 {
		t.Fatalf("expected 31 official operations, got %d", len(reference.Endpoints))
	}
	for _, e := range reference.Endpoints {
		t.Run(e.ID, func(t *testing.T) {
			args := []string{"--project", t.TempDir(), "publishing", e.ID, "--invoke", "--dry-run=false", "--token", "offline-publishing-token"}
			values := map[string]any{}
			for _, p := range e.Parameters {
				var value any = "fixture"
				encoded := "fixture"
				if p.In == "body" {
					switch {
					case strings.HasPrefix(p.Type, "Integer"):
						value = float64(1)
						encoded = "1"
					case strings.HasPrefix(p.Type, "List<String"):
						value = []any{"fixture"}
						encoded = `["fixture"]`
					case strings.HasPrefix(p.Type, "List<Integer"):
						value = []any{float64(1)}
						encoded = `[1]`
					case strings.HasPrefix(p.Type, "List<"):
						value = []any{map[string]any{"fixture": "value"}}
						encoded = `[{"fixture":"value"}]`
					case !strings.HasPrefix(p.Type, "String"):
						value = map[string]any{"fixture": "value"}
						encoded = `{"fixture":"value"}`
					}
				}
				values[p.In+"/"+p.Name] = value
				flag := map[string]string{"query": "--query", "header": "--header", "body": "--field"}[p.In]
				args = append(args, flag, p.Name+"="+encoded)
			}
			called := false
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				called = true
				if r.Method != e.Method || r.URL.Path != e.Path {
					t.Errorf("got %s %s, want %s %s", r.Method, r.URL.Path, e.Method, e.Path)
				}
				if r.Header.Get("Authorization") != "Bearer offline-publishing-token" {
					t.Error("authorization missing")
				}
				var body map[string]any
				raw, _ := io.ReadAll(r.Body)
				if len(raw) > 0 && json.Unmarshal(raw, &body) != nil {
					t.Error("invalid JSON body")
				}
				for _, p := range e.Parameters {
					want := values[p.In+"/"+p.Name]
					switch p.In {
					case "query":
						if r.URL.Query().Get(p.Name) != want {
							t.Errorf("query %s missing", p.Name)
						}
					case "header":
						if r.Header.Get(p.Name) != want {
							t.Errorf("header %s missing", p.Name)
						}
					case "body":
						got, _ := json.Marshal(body[p.Name])
						expected, _ := json.Marshal(want)
						if string(got) != string(expected) {
							t.Errorf("body %s=%s, want %s", p.Name, got, expected)
						}
					}
				}
				w.Header().Set("Content-Type", "application/json")
				io.WriteString(w, `{"ret":{"code":0,"msg":"success"}}`)
			}))
			defer server.Close()
			args = append(args, "--base-url", server.URL)
			_, err := execute(args...)
			if err != nil {
				t.Fatal(err)
			}
			if !called {
				t.Fatal("CLI never sent request")
			}
		})
	}
}

func TestPublishingVersionListSendsEmptyJSONObject(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		raw, _ := io.ReadAll(r.Body)
		if strings.TrimSpace(string(raw)) != "{}" {
			t.Errorf("got body %q, want {}", raw)
		}
		io.WriteString(w, `{"ret":{"code":0}}`)
	}))
	defer server.Close()
	_, err := execute("--project", t.TempDir(), "publishing", "app-version-list", "--invoke", "--dry-run=false", "--token", "offline", "--base-url", server.URL, "--header", "appId=fixture")
	if err != nil {
		t.Fatal(err)
	}
}

func TestPublishingAgreementListsAcceptExplicitEmptyName(t *testing.T) {
	for _, id := range []string{"user-agreement-list", "privacy-agreement-list"} {
		t.Run(id, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				var body map[string]any
				json.NewDecoder(r.Body).Decode(&body)
				if name, ok := body["name"]; !ok || name != "" {
					t.Errorf("expected explicit empty name, got %v", body)
				}
				io.WriteString(w, `{"ret":{"code":0}}`)
			}))
			defer server.Close()
			_, err := execute("--project", t.TempDir(), "publishing", id, "--invoke", "--dry-run=false", "--token", "offline", "--base-url", server.URL, "--header", "appId=fixture", "--field", "name=")
			if err != nil {
				t.Fatal(err)
			}
		})
	}
}

// Huawei's creation table says AgreementInfo, but its description and request
// example explicitly use an array of localized agreements with nested baseInfo.
func TestPublishingUserAgreementAcceptsOfficialArrayExample(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var body map[string]any
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatal(err)
		}
		locales, ok := body["agreementInfos"].([]any)
		if !ok || len(locales) != 1 {
			t.Errorf("expected localized agreement array, got %v", body)
		}
		if body["channel"] != float64(1) {
			t.Errorf("channel must remain a JSON integer, got %v", body["channel"])
		}
		io.WriteString(w, `{"ret":{"code":0}}`)
	}))
	defer server.Close()
	_, err := execute("--project", t.TempDir(), "publishing", "user-agreement-create", "--invoke", "--dry-run=false", "--token", "offline", "--base-url", server.URL, "--header", "appId=fixture", "--field", "type=2", "--field", "channel=1", "--field", `agreementInfos=[{"lang":"zh","baseInfo":{"name":"Fixture","agreementTitle":"Fixture terms","userAgreementUrl":"https://example.com/terms"}}]`)
	if err != nil {
		t.Fatal(err)
	}
}

func TestPublishingHelpExplainsRequiredParameterLocations(t *testing.T) {
	out, err := execute("publishing", "app-package-info-update", "--help")
	if err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{"--query appId=", "--field fileName=", "--field objectId="} {
		if !strings.Contains(out, want) {
			t.Errorf("help missing %q", want)
		}
	}
}
