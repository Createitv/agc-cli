package domain

import (
	"strconv"
	"strings"
)

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
	op["x-huawei-default-base-url"] = endpoint.DefaultBaseURL
	op["x-huawei-authentication"] = endpoint.Authentication
	op["x-huawei-content-type"] = endpoint.RequestContentType
	op["x-agc-parameters"] = endpoint.Parameters
	op["x-agc-rules"] = endpoint.Rules
	op["x-agc-documentation"] = endpoint.Documentation
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
	if prefix == "invoke" {
		schema := op["requestBody"].(map[string]any)["content"].(map[string]any)["application/json"].(map[string]any)["schema"].(map[string]any)
		if (endpoint.FamilyID == "pms" && endpoint.ID == "resourcemanage-harmonyosnext") || (IsServerlessFamily(endpoint.FamilyID) && parameterFileRequired(endpoint)) {
			if IsServerlessFamily(endpoint.FamilyID) {
				schema["required"] = invokeRequiredContainers(endpoint)
			} else {
				schema["required"] = []string{"filePath"}
			}
		}
		alternatives := []any{}
		for _, parameter := range endpoint.Parameters {
			if parameter.In != "query" || !parameter.Required {
				continue
			}
			branches := []any{}
			for _, container := range []string{"query", "queryValues"} {
				branches = append(branches, map[string]any{"required": []string{container}, "properties": map[string]any{container: map[string]any{"required": []string{parameter.Name}}}})
			}
			alternatives = append(alternatives, map[string]any{"anyOf": branches})
		}
		if endpoint.Authentication == "agcli-user-oauth" {
			required := []string{}
			if endpoint.ID != "list-functions" {
				required = append(required, "name")
			}
			if endpoint.ID == "create-function" || endpoint.ID == "update-function-code" {
				required = append(required, "codePath")
			}
			if strings.Contains(endpoint.ID, "publish") {
				required = append(required, "versionDescription")
			}
			schema["required"] = required
			alternatives = append(alternatives, map[string]any{"if": map[string]any{"properties": map[string]any{"dryRun": map[string]any{"const": false}}, "required": []string{"dryRun"}}, "then": map[string]any{"required": []string{"sessionFile"}}})
		}
		if strings.HasPrefix(endpoint.Authentication, "agcli-management-") {
			alternatives = append(alternatives, map[string]any{"if": map[string]any{"properties": map[string]any{"dryRun": map[string]any{"const": false}}, "required": []string{"dryRun"}}, "then": map[string]any{"required": []string{"sessionFile"}}})
		}
		if len(alternatives) > 0 {
			schema["allOf"] = alternatives
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
	if endpoint.Authentication == "agcli-user-oauth" {
		return FunctionSchema()["properties"].(map[string]any)
	}
	properties := map[string]any{
		"sessionFile": map[string]any{"type": "string", "writeOnly": true, "description": "Explicit user management session file; not a publishing profile."},
		"projectId":   map[string]any{"type": "string", "description": "Project must match the selected management session."},
		"filePath":    map[string]string{"type": "string", "description": "Explicit local file path on the CLI server host for supported file-upload endpoints. Multipart metadata is supplied through fields or body."},
		"queryValues": map[string]any{"type": "object", "description": "Ordered repeated upstream query keys. Each array replaces the scalar query entry with the same name.", "additionalProperties": map[string]any{"type": "array", "items": map[string]string{"type": "string"}}},
		"baseUrl":     map[string]string{"type": "string"},
		"params":      map[string]any{"type": "object", "additionalProperties": map[string]string{"type": "string"}},
		"query":       map[string]any{"type": "object", "additionalProperties": map[string]string{"type": "string"}},
		"headers":     map[string]any{"type": "object", "additionalProperties": map[string]string{"type": "string"}},
		"fields":      map[string]any{"type": "object", "additionalProperties": map[string]string{"type": "string"}},
		"body":        map[string]any{"description": "Raw JSON request body for the Huawei endpoint."},
		"token":       map[string]string{"type": "string", "description": "Huawei upstream access token; separate from the local X-AGC-Server-Token header."},
		"dryRun":      map[string]string{"type": "boolean"},
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
			if (endpoint.FamilyID == "pms" && endpoint.ID == "resourcemanage-harmonyosnext") || IsServerlessFamily(endpoint.FamilyID) {
				continue
			}
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
		if parameter.Pattern != "" {
			known[parameter.Name].(map[string]any)["pattern"] = parameter.Pattern
		}
		if parameter.Required && container != "query" {
			required, _ := object["required"].([]string)
			object["required"] = append(required, parameter.Name)
		}
	}
	if methodNeedsBody(endpoint.Method) || endpoint.Body != "" {
		bodySchema := map[string]any{
			"description": "Raw JSON request body for " + endpoint.Method + " " + endpoint.Path + ".",
		}
		if endpoint.DynamicBody != nil {
			child := invokeRequestSchemaProperties(Endpoint{Method: "POST", Parameters: endpoint.DynamicBody.Properties})["body"]
			bodySchema["type"] = "object"
			bodySchema["patternProperties"] = map[string]any{endpoint.DynamicBodyKeyPattern: child}
			bodySchema["additionalProperties"] = false
			if endpoint.DynamicBody.Required {
				bodySchema["minProperties"] = 1
			}
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
			if parameter.EncodedJSON != nil {
				model := *parameter.EncodedJSON
				model.Name = "value"
				model.In = "body"
				nested := invokeRequestSchemaProperties(Endpoint{Method: "POST", Parameters: []Parameter{model}})["body"].(map[string]any)["properties"].(map[string]any)["value"]
				field["contentMediaType"] = "application/json"
				field["contentSchema"] = nested
			}
			if parameter.Format != "" {
				field["format"] = parameter.Format
			}
			if parameter.Pattern != "" {
				field["pattern"] = parameter.Pattern
			}
			if parameter.MaxLength > 0 {
				field["maxLength"] = parameter.MaxLength
			}
			if parameter.MinLength > 0 {
				field["minLength"] = parameter.MinLength
			}
			if parameter.Minimum != nil {
				field["minimum"] = *parameter.Minimum
			}
			if parameter.Maximum != nil {
				field["maximum"] = *parameter.Maximum
			}
			if len(parameter.Enum) > 0 {
				if parameter.Type == "integer" {
					values := []int64{}
					for _, text := range parameter.Enum {
						if value, err := strconv.ParseInt(text, 10, 64); err == nil {
							values = append(values, value)
						}
					}
					field["enum"] = values
				} else {
					field["enum"] = parameter.Enum
				}
			}
			if parameter.Type == "string" && parameter.Required && !parameter.AllowEmpty {
				field["minLength"] = 1
			}
			if parameter.Type == "object" && len(parameter.Properties) > 0 {
				nested := invokeRequestSchemaProperties(Endpoint{Method: "POST", Parameters: parameter.Properties})["body"].(map[string]any)
				for name, value := range nested {
					field[name] = value
				}
			}
			if parameter.Type == "object" {
				if parameter.MinProperties > 0 {
					field["minProperties"] = parameter.MinProperties
				}
				if parameter.MaxProperties > 0 {
					field["maxProperties"] = parameter.MaxProperties
				}
				if parameter.PropertyNamesPattern != "" {
					field["propertyNames"] = map[string]any{"pattern": parameter.PropertyNamesPattern}
				}
				if parameter.AdditionalProperties != nil {
					value := *parameter.AdditionalProperties
					value.In = "body"
					value.Name = "value"
					value.Required = false
					nested := invokeRequestSchemaProperties(Endpoint{Method: "POST", Parameters: []Parameter{value}})["body"].(map[string]any)["properties"].(map[string]any)["value"]
					field["additionalProperties"] = nested
				}
			}
			if parameter.Type == "array" {
				if parameter.ItemsType != "" {
					itemSchema := map[string]any{"type": parameter.ItemsType}
					if len(parameter.ItemsEnum) > 0 {
						if parameter.ItemsType == "integer" {
							values := []int64{}
							for _, text := range parameter.ItemsEnum {
								if value, err := strconv.ParseInt(text, 10, 64); err == nil {
									values = append(values, value)
								}
							}
							itemSchema["enum"] = values
						} else {
							itemSchema["enum"] = parameter.ItemsEnum
						}
					}
					if parameter.ItemsPattern != "" {
						itemSchema["pattern"] = parameter.ItemsPattern
					}
					if parameter.ItemsMaxLength > 0 {
						itemSchema["maxLength"] = parameter.ItemsMaxLength
					}
					field["items"] = itemSchema
					if len(parameter.ItemsProperties) > 0 {
						nested := invokeRequestSchemaProperties(Endpoint{Method: "POST", Parameters: parameter.ItemsProperties})
						field["items"] = nested["body"]
						if len(parameter.ItemsRules) > 0 {
							field["items"].(map[string]any)["x-agc-rules"] = parameter.ItemsRules
						}
					}
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
	if endpoint.DynamicBody != nil && endpoint.DynamicBody.Required {
		required = append(required, "body")
		seen["body"] = true
	}
	for _, parameter := range endpoint.Parameters {
		if !parameter.Required {
			continue
		}
		container := ""
		switch parameter.In {
		case "header":
			container = "headers"
		case "query":
			continue
		case "path":
			container = "params"
		case "file":
			if IsServerlessFamily(endpoint.FamilyID) {
				container = "filePath"
				break
			}
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

func parameterFileRequired(e Endpoint) bool {
	for _, p := range e.Parameters {
		if p.In == "file" && p.Required {
			return true
		}
	}
	return false
}
