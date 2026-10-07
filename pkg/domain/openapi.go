package domain

import "strings"

func OpenAPISpec() map[string]any {
	paths := map[string]any{
		"/api/v1": map[string]any{
			"get": operation("getRoot", "Command Center root", "Discover local AGC CLI REST affordances.", nil),
		},
		"/api/v1/capabilities": map[string]any{
			"get": operation("listCapabilities", "List API families", "List every AppGallery Connect API family registered in the CLI.", nil),
		},
		"/api/v1/endpoints": map[string]any{
			"get": operation("listEndpoints", "List registered endpoints", "List every registered AppGallery Connect endpoint.", nil),
		},
		"/api/v1/openapi.json": map[string]any{
			"get": operation("getOpenAPI", "Export OpenAPI", "Export this local REST interface contract.", nil),
		},
	}
	for _, capability := range DecoratedCapabilities() {
		paths[capability.RESTPath] = map[string]any{
			"get": operation("get"+exportedID(capability.ID)+"Capability", capability.Name, capability.Description, []string{capability.Name}),
		}
		paths[capability.RESTPath+"/endpoints"] = map[string]any{
			"get": operation("list"+exportedID(capability.ID)+"Endpoints", "List "+capability.Name+" endpoints", capability.Description, []string{capability.Name}),
		}
	}
	for _, endpoint := range AllEndpoints() {
		tags := []string{endpoint.FamilyID}
		showPath := "/api/v1/" + endpoint.FamilyID + "/endpoints/" + endpoint.ID
		invokePath := showPath + "/invoke"
		paths[showPath] = map[string]any{
			"get": endpointOperation("get", endpoint, "Show endpoint", tags),
		}
		paths[invokePath] = map[string]any{
			"post": endpointOperation("invoke", endpoint, "Invoke endpoint", tags),
		}
	}
	return map[string]any{
		"openapi": "3.1.0",
		"info": map[string]any{
			"title":       "agc-cli local REST API",
			"version":     "0.1.0",
			"description": "Local REST and OpenAPI contract for the Huawei AppGallery Connect CLI command center. All routes require X-AGC-Server-Token when AGC_SERVER_TOKEN is configured; otherwise local server authentication is optional. Huawei upstream credentials are supplied in the invoke JSON token or headers fields.",
		},
		"servers": []map[string]string{
			{"url": "http://localhost:8421"},
		},
		"paths":                    paths,
		"security":                 []map[string][]string{{}, {"localServerToken": {}}},
		"x-agc-auth-required-when": "AGC_SERVER_TOKEN is configured",
		"components": map[string]any{
			"securitySchemes": map[string]any{
				"localServerToken": map[string]string{
					"type":        "apiKey",
					"in":          "header",
					"name":        "X-AGC-Server-Token",
					"description": "Required on every local route when AGC_SERVER_TOKEN is configured. This token authenticates the local server, not Huawei.",
				},
			},
		},
	}
}

func endpointOperation(prefix string, endpoint Endpoint, summaryPrefix string, tags []string) map[string]any {
	op := operation(
		prefix+exportedID(endpoint.FamilyID)+exportedID(endpoint.ID),
		summaryPrefix+": "+endpoint.Name,
		endpoint.Description,
		tags,
	)
	op["x-agc-family"] = endpoint.FamilyID
	op["x-agc-command"] = endpoint.Command
	op["x-huawei-method"] = endpoint.Method
	op["x-huawei-path"] = endpoint.Path
	op["x-agc-parameters"] = endpoint.Parameters
	if prefix == "invoke" {
		op["requestBody"] = map[string]any{
			"required": true,
			"content": map[string]any{
				"application/json": map[string]any{
					"schema": map[string]any{
						"type":       "object",
						"properties": invokeRequestSchemaProperties(endpoint),
						"required":   invokeRequiredContainers(endpoint),
					},
				},
			},
		}
	}
	return op
}

func operation(operationID, summary, description string, tags []string) map[string]any {
	op := map[string]any{
		"operationId": operationID,
		"summary":     summary,
		"description": description,
		"responses": map[string]any{
			"200": map[string]string{"description": "OK"},
		},
	}
	if len(tags) > 0 {
		op["tags"] = tags
	}
	return op
}

func invokeRequestSchemaProperties(endpoint Endpoint) map[string]any {
	properties := map[string]any{
		"baseUrl": map[string]string{"type": "string"},
		"params":  map[string]any{"type": "object", "additionalProperties": map[string]string{"type": "string"}},
		"query":   map[string]any{"type": "object", "additionalProperties": map[string]string{"type": "string"}},
		"headers": map[string]any{"type": "object", "additionalProperties": map[string]string{"type": "string"}},
		"fields":  map[string]any{"type": "object", "additionalProperties": map[string]string{"type": "string"}},
		"body":    map[string]any{"description": "Raw JSON request body for the Huawei endpoint."},
		"token":   map[string]string{"type": "string", "description": "Huawei upstream access token; separate from the local X-AGC-Server-Token header."},
		"dryRun":  map[string]string{"type": "boolean"},
	}
	for _, parameter := range endpoint.Parameters {
		container := ""
		switch parameter.In {
		case "header":
			container = "headers"
		case "query":
			container = "query"
		case "path":
			container = "params"
		case "file":
			container = "fields"
		}
		if container == "" {
			continue
		}
		object := properties[container].(map[string]any)
		known, ok := object["properties"].(map[string]any)
		if !ok {
			known = map[string]any{}
			object["properties"] = known
		}
		known[parameter.Name] = map[string]any{"type": "string", "description": parameter.Description, "x-huawei-required": parameter.Required}
		// Local query/header maps hold strings; preserve upstream types separately.
		if parameter.Type != "" {
			known[parameter.Name].(map[string]any)["x-huawei-type"] = parameter.Type
		}
		if parameter.Format != "" {
			known[parameter.Name].(map[string]any)["x-huawei-format"] = parameter.Format
		}
		if parameter.Required {
			required, _ := object["required"].([]string)
			object["required"] = append(required, parameter.Name)
		}
	}
	if methodNeedsBody(endpoint.Method) {
		bodySchema := map[string]any{
			"description": "Raw JSON request body for " + endpoint.Method + " " + endpoint.Path + ".",
		}
		fields := map[string]any{}
		required := []string{}
		for _, parameter := range endpoint.Parameters {
			if parameter.In != "body" {
				continue
			}
			field := map[string]any{"description": parameter.Description}
			if parameter.Type != "" {
				field["type"] = parameter.Type
			}
			if parameter.Format != "" {
				field["format"] = parameter.Format
			}
			if parameter.Type == "array" {
				if parameter.ItemsType != "" {
					field["items"] = map[string]string{"type": parameter.ItemsType}
				}
				if parameter.MinItems > 0 {
					field["minItems"] = parameter.MinItems
				}
				if parameter.MaxItems > 0 {
					field["maxItems"] = parameter.MaxItems
				}
			}
			fields[parameter.Name] = field
			if parameter.Required {
				required = append(required, parameter.Name)
			}
		}
		if len(fields) > 0 {
			bodySchema["type"] = "object"
			bodySchema["properties"] = fields
			if len(required) > 0 {
				bodySchema["required"] = required
			}
		}
		properties["body"] = bodySchema
	}
	return properties
}

func exportedID(id string) string {
	parts := strings.FieldsFunc(id, func(r rune) bool {
		return r == '-' || r == '_' || r == ' '
	})
	for i, part := range parts {
		if part == "" {
			continue
		}
		parts[i] = strings.ToUpper(part[:1]) + part[1:]
	}
	return strings.Join(parts, "")
}

func invokeRequiredContainers(endpoint Endpoint) []string {
	required := []string{}
	seen := map[string]bool{}
	for _, parameter := range endpoint.Parameters {
		if !parameter.Required {
			continue
		}
		container := ""
		switch parameter.In {
		case "header":
			container = "headers"
		case "query":
			container = "query"
		case "path":
			container = "params"
		case "file":
			container = "fields"
		case "body":
			container = "body"
		}
		if container != "" && !seen[container] {
			required = append(required, container)
			seen[container] = true
		}
	}
	return required
}
