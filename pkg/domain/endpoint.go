package domain

import "strings"

type Endpoint struct {
	DefaultBaseURL        string            `json:"defaultBaseUrl,omitempty"`
	Authentication        string            `json:"authentication,omitempty"`
	DynamicBody           *Parameter        `json:"dynamicBody,omitempty"`
	DynamicBodyKeyPattern string            `json:"dynamicBodyKeyPattern,omitempty"`
	RequestContentType    string            `json:"requestContentType,omitempty"`
	Documentation         *ReferenceSummary `json:"documentation,omitempty"`
	Rules                 []ReferenceRule   `json:"rules,omitempty"`
	ID                    string            `json:"id"`
	FamilyID              string            `json:"familyId"`
	Name                  string            `json:"name"`
	Description           string            `json:"description"`
	Method                string            `json:"method"`
	Path                  string            `json:"path"`
	Command               string            `json:"command"`
	RequiredParams        []string          `json:"requiredParams,omitempty"`
	Parameters            []Parameter       `json:"parameters,omitempty"`
	Body                  string            `json:"body,omitempty"`
	Status                string            `json:"status"`
	Source                string            `json:"source,omitempty"`
	SourceURL             string            `json:"sourceUrl,omitempty"`
	OfficialSlug          string            `json:"officialSlug,omitempty"`
	Direction             string            `json:"direction,omitempty"`
	Links                 map[string]Link   `json:"_links,omitempty"`
	Affordances           map[string]string `json:"affordances,omitempty"`
}

type Parameter struct {
	AdditionalProperties *Parameter      `json:"additionalProperties,omitempty"`
	PropertyNamesPattern string          `json:"propertyNamesPattern,omitempty"`
	MinProperties        int             `json:"minProperties,omitempty"`
	MaxProperties        int             `json:"maxProperties,omitempty"`
	Pattern              string          `json:"pattern,omitempty"`
	ItemsEnum            []string        `json:"itemsEnum,omitempty"`
	EncodedJSON          *Parameter      `json:"encodedJson,omitempty"`
	ItemsPattern         string          `json:"itemsPattern,omitempty"`
	ItemsMaxLength       int             `json:"itemsMaxLength,omitempty"`
	Properties           []Parameter     `json:"properties,omitempty"`
	ItemsRules           []ReferenceRule `json:"itemsRules,omitempty"`
	ItemsProperties      []Parameter     `json:"itemsProperties,omitempty"`
	SourceURL            string          `json:"sourceUrl,omitempty"`
	MaxLength            int             `json:"maxLength,omitempty"`
	MinLength            int             `json:"minLength,omitempty"`
	Minimum              *int64          `json:"minimum,omitempty"`
	Maximum              *int64          `json:"maximum,omitempty"`
	Enum                 []string        `json:"enum,omitempty"`
	AllowEmpty           bool            `json:"allowEmpty,omitempty"`
	MinItems             int             `json:"minItems,omitempty"`
	MaxItems             int             `json:"maxItems,omitempty"`
	ItemsType            string          `json:"itemsType,omitempty"`
	Type                 string          `json:"type,omitempty"`
	Format               string          `json:"format,omitempty"`
	Name                 string          `json:"name"`
	In                   string          `json:"in"`
	Required             bool            `json:"required"`
	Description          string          `json:"description,omitempty"`
}

func endpoint(familyID, id, name, description, method, path string, requiredParams []string) Endpoint {
	parameters := verifiedParameters(familyID, id, inferParameters(method, path, requiredParams))
	requiredParams = nil
	for _, parameter := range parameters {
		if parameter.Required {
			requiredParams = append(requiredParams, parameter.Name)
		}
	}
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
	if familyID == "gameplay" && id == "resources" {
		endpoint.DefaultBaseURL = "https://jos-open-api.cloud.huawei.com"
		endpoint.RequestContentType = "application/x-www-form-urlencoded"
		endpoint.Authentication = "body-signature"
	}
	if familyID == "game-items" && (id == "propapi-getroleinfo" || id == "propapi-order") {
		endpoint.RequestContentType = "application/x-www-form-urlencoded; charset=UTF-8"
	}
	if direction == "inbound-callback" {
		endpoint.Source = "Huawei AppGallery Connect inbound callback reference"
	}
	if direction == "local" {
		endpoint.Source = "agc-cli local integration"
	}
	return applyReferenceParameters(endpoint)
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
	fields := []string{}
	for _, parameter := range parameters {
		if parameter.In == "body" || parameter.In == "file" {
			fields = append(fields, parameter.Name)
		}
	}
	if len(fields) == 0 {
		if !methodNeedsBody(method) {
			return ""
		}
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
	endpoint = attachReference(endpoint)
	href := "/api/v1/" + endpoint.FamilyID + "/endpoints/" + endpoint.ID
	endpoint.Links = map[string]Link{
		"self":   {Href: href, Method: "GET"},
		"family": {Href: "/api/v1/" + endpoint.FamilyID + "/endpoints", Method: "GET"},
	}
	endpoint.Affordances = map[string]string{
		"show": endpoint.Command + " --pretty",
	}
	if strings.HasPrefix(endpoint.RequestContentType, "application/x-www-form-urlencoded") {
		endpoint.Affordances["formInput"] = "Supply raw unescaped values using --field or a JSON object using --body; the CLI performs UTF-8 form encoding once."
	}
	if endpoint.Authentication == "agcli-user-oauth" {
		endpoint.Affordances["invoke"] = endpoint.Command + " --invoke --name FUNCTION --session-file SESSION.json"
		endpoint.Affordances["docs"] = endpoint.Command + " --docs --pretty"
		return endpoint
	}
	if strings.HasPrefix(endpoint.Authentication, "agcli-management-") {
		endpoint.Affordances["invoke"] = endpoint.Command + " --invoke --session-file SESSION.json"
		if endpoint.Body != "" {
			endpoint.Affordances["invoke"] += " --body body.json"
		}
		endpoint.Affordances["docs"] = endpoint.Command + " --docs --pretty"
		if endpoint.FamilyID == "clouddb" && endpoint.ID == "import-object-types" {
			endpoint.Affordances["formInput"] = "Supply raw schema JSON with --body. The reviewed vendor uses raw JSON bytes with a form Content-Type; no form encoding is performed."
		}
		return endpoint
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
