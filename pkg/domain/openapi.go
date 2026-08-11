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
			"description": "Local REST and OpenAPI contract for the Huawei AppGallery Connect CLI command center.",
		},
		"servers": []map[string]string{
			{"url": "http://localhost:8421"},
		},
		"paths": paths,
		"components": map[string]any{
			"securitySchemes": map[string]any{
				"bearerAuth": map[string]string{
					"type":         "http",
					"scheme":       "bearer",
					"bearerFormat": "Huawei AppGallery Connect access token",
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
					},
				},
			},
		}
		op["security"] = []map[string][]string{{"bearerAuth": []string{}}}
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
		"token":   map[string]string{"type": "string"},
		"dryRun":  map[string]string{"type": "boolean"},
	}
	if methodNeedsBody(endpoint.Method) {
		properties["body"] = map[string]any{
			"description": "Raw JSON request body for " + endpoint.Method + " " + endpoint.Path + ".",
		}
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
