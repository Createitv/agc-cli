package domain

import "encoding/json"

// FunctionRequest is independent of HTTP, session persistence and CLI presentation.
// Nil configuration pointers mean preserve existing values, including empty maps.
type FunctionRequest struct {
	Name                string             `json:"name,omitempty"`
	SessionFile         string             `json:"sessionFile,omitempty"`
	BaseURL             string             `json:"baseUrl,omitempty"`
	ProjectID           string             `json:"projectId,omitempty"`
	CodePath            string             `json:"codePath,omitempty"`
	Runtime             string             `json:"runtime,omitempty"`
	Handler             *string            `json:"handler,omitempty"`
	Description         string             `json:"description,omitempty"`
	VersionDescription  string             `json:"versionDescription,omitempty"`
	Version             string             `json:"version,omitempty"`
	AliasName           string             `json:"aliasName,omitempty"`
	Memory              *int               `json:"memory,omitempty"`
	Timeout             *int               `json:"timeout,omitempty"`
	InstanceConcurrency *int               `json:"instanceConcurrency,omitempty"`
	Environment         *map[string]string `json:"environment,omitempty"`
	Layers              *[]string          `json:"layers,omitempty"`
	Payload             json.RawMessage    `json:"payload,omitempty"`
	FuzzyName           string             `json:"fuzzyName,omitempty"`
	PageSize            int                `json:"pageSize,omitempty"`
	PageIndex           int                `json:"pageIndex,omitempty"`
	DryRun              bool               `json:"dryRun"`
}

var FunctionOperations = []string{"create-function", "update-function-code", "update-function-config", "publish-function", "update-and-publish-function", "invoke", "list-functions", "query-function-version", "list-function-versions"}

func init() {
	Capabilities = append(Capabilities, Capability{ID: "functions", Name: "Cloud Functions", Description: "Native Go workflows based on Huawei AGCLI 1.12.0 source; explicit user OAuth session required.", Command: "agc functions", RESTPath: "/api/v1/functions", Status: "registered"})
	for _, id := range FunctionOperations {
		e := Endpoint{ID: id, FamilyID: "functions", Name: id, Command: "agc functions " + id, Method: "POST", Path: "native:functions/" + id, Direction: "native-workflow", Authentication: "agcli-user-oauth", Status: "registered", OfficialSlug: "agcli-1.12.0-" + id, Source: "Huawei official AGCLI 1.12.0 source", SourceURL: "https://www.npmjs.com/package/@hw-agconnect/agc-cli/v/1.12.0"}
		e.Documentation = &ReferenceSummary{Title: id, SourceURL: e.SourceURL, CheckedAt: "2026-10-26", Status: "vendor-source-reviewed", ReadCommand: e.Command + " --docs --pretty"}
		EndpointFamilies["functions"] = append(EndpointFamilies["functions"], e)
	}
}
func FunctionSchema() map[string]any {
	props := map[string]any{}
	for _, n := range []string{"name", "sessionFile", "baseUrl", "projectId", "codePath", "runtime", "handler", "description", "versionDescription", "version", "aliasName", "fuzzyName"} {
		props[n] = map[string]any{"type": "string"}
	}
	for _, n := range []string{"memory", "timeout", "instanceConcurrency", "pageSize", "pageIndex"} {
		props[n] = map[string]any{"type": "integer"}
	}
	props["name"].(map[string]any)["pattern"] = "^[a-z][a-z0-9-]{0,30}[a-z0-9]$"
	props["runtime"].(map[string]any)["enum"] = []string{"nodejs14.x", "java1.8", "python3", "custom-runtime"}
	props["memory"].(map[string]any)["enum"] = []int{500, 1000, 2000, 4000}
	for n, maximum := range map[string]int{"timeout": 1800, "instanceConcurrency": 1000, "pageSize": 1000} {
		props[n].(map[string]any)["minimum"] = 1
		props[n].(map[string]any)["maximum"] = maximum
	}
	props["pageIndex"].(map[string]any)["minimum"] = 1
	props["version"].(map[string]any)["pattern"] = `^(\$latest|[1-9][0-9]{0,2}|1000)$`
	for _, n := range []string{"description", "versionDescription"} {
		props[n].(map[string]any)["maxLength"] = 64
		props[n].(map[string]any)["x-agcli-display-width-max"] = 64
	}
	props["environment"] = map[string]any{"type": "object", "additionalProperties": map[string]string{"type": "string"}, "writeOnly": true}
	props["layers"] = map[string]any{"type": "array", "maxItems": 5, "items": map[string]string{"type": "string"}}
	props["payload"] = map[string]any{"description": "Arbitrary valid JSON; numbers retain precision."}
	props["dryRun"] = map[string]any{"type": "boolean", "default": true}
	return map[string]any{"type": "object", "properties": props}
}
