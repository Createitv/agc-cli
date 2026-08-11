package domain

import "strings"

type Endpoint struct {
	ID             string            `json:"id"`
	FamilyID       string            `json:"familyId"`
	Name           string            `json:"name"`
	Description    string            `json:"description"`
	Method         string            `json:"method"`
	Path           string            `json:"path"`
	Command        string            `json:"command"`
	RequiredParams []string          `json:"requiredParams,omitempty"`
	Parameters     []Parameter       `json:"parameters,omitempty"`
	Body           string            `json:"body,omitempty"`
	Status         string            `json:"status"`
	Source         string            `json:"source,omitempty"`
	SourceURL      string            `json:"sourceUrl,omitempty"`
	OfficialSlug   string            `json:"officialSlug,omitempty"`
	Direction      string            `json:"direction,omitempty"`
	Links          map[string]Link   `json:"_links,omitempty"`
	Affordances    map[string]string `json:"affordances,omitempty"`
}

type Parameter struct {
	Name        string `json:"name"`
	In          string `json:"in"`
	Required    bool   `json:"required"`
	Description string `json:"description,omitempty"`
}

func endpoint(familyID, id, name, description, method, path string, requiredParams []string) Endpoint {
	parameters := inferParameters(method, path, requiredParams)
	return Endpoint{
		ID:             id,
		FamilyID:       familyID,
		Name:           name,
		Description:    description,
		Method:         method,
		Path:           path,
		Command:        "agc " + familyID + " " + id,
		RequiredParams: requiredParams,
		Parameters:     parameters,
		Body:           inferBody(method, parameters),
		Status:         "registered",
		Source:         "Huawei AppGallery Connect API reference",
		Direction:      "developer-to-huawei",
	}
}

func officialEndpoint(familyID, id, name, description, method, path string, requiredParams []string, slug, sourceURL, direction string) Endpoint {
	endpoint := endpoint(familyID, id, name, description, method, path, requiredParams)
	endpoint.OfficialSlug = slug
	endpoint.SourceURL = sourceURL
	endpoint.Direction = direction
	if direction == "inbound-callback" {
		endpoint.Source = "Huawei AppGallery Connect inbound callback reference"
	}
	if direction == "local" {
		endpoint.Source = "agc-cli local integration"
	}
	return endpoint
}

func inferParameters(method, path string, requiredParams []string) []Parameter {
	parameters := make([]Parameter, 0, len(requiredParams))
	for _, name := range requiredParams {
		location := "query"
		if strings.Contains(path, "{"+name+"}") || name == "uploadUrl" || name == "callbackUrl" {
			location = "path"
		} else if name == "file" {
			location = "file"
		} else if methodNeedsBody(method) {
			location = "body"
		}
		parameters = append(parameters, Parameter{
			Name:        name,
			In:          location,
			Required:    true,
			Description: name + " parameter",
		})
	}
	return parameters
}

func inferBody(method string, parameters []Parameter) string {
	if !methodNeedsBody(method) {
		return ""
	}
	fields := []string{}
	for _, parameter := range parameters {
		if parameter.In == "body" || parameter.In == "file" {
			fields = append(fields, parameter.Name)
		}
	}
	if len(fields) == 0 {
		return "json"
	}
	return strings.Join(fields, ",")
}

func EndpointsByFamily(familyID string) []Endpoint {
	endpoints := EndpointFamilies[familyID]
	out := make([]Endpoint, len(endpoints))
	for i, endpoint := range endpoints {
		out[i] = decorateEndpoint(endpoint)
	}
	return out
}

func AllEndpoints() []Endpoint {
	out := []Endpoint{}
	for _, capability := range Capabilities {
		out = append(out, EndpointsByFamily(capability.ID)...)
	}
	return out
}

func EndpointByID(familyID, endpointID string) (Endpoint, bool) {
	for _, endpoint := range EndpointFamilies[familyID] {
		if endpoint.ID == endpointID {
			return decorateEndpoint(endpoint), true
		}
	}
	return Endpoint{}, false
}

func EndpointCount() int {
	return len(AllEndpoints())
}

func decorateEndpoint(endpoint Endpoint) Endpoint {
	href := "/api/v1/" + endpoint.FamilyID + "/endpoints/" + endpoint.ID
	endpoint.Links = map[string]Link{
		"self":   {Href: href, Method: "GET"},
		"family": {Href: "/api/v1/" + endpoint.FamilyID + "/endpoints", Method: "GET"},
	}
	endpoint.Affordances = map[string]string{
		"show": endpoint.Command + " --pretty",
	}
	if methodNeedsBody(endpoint.Method) {
		endpoint.Affordances["invoke"] = endpoint.Command + " --invoke --body body.json"
	} else {
		endpoint.Affordances["invoke"] = endpoint.Command + " --invoke"
	}
	return endpoint
}

func methodNeedsBody(method string) bool {
	switch strings.ToUpper(method) {
	case "POST", "PUT", "PATCH":
		return true
	default:
		return false
	}
}
