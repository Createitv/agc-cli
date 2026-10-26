package server

import (
	"context"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/Createitv/agc-cli/pkg/agcapi"
	"github.com/Createitv/agc-cli/pkg/cloudfunctions"
	"github.com/Createitv/agc-cli/pkg/domain"
)

func Handler() http.Handler {
	return HandlerWithToken("")
}

// HandlerWithToken protects the local command API when the operator exposes it
// beyond loopback. This token is separate from Huawei authorization credentials.
func HandlerWithToken(token string) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/v1", root)
	mux.HandleFunc("/api/v1/capabilities", capabilities)
	mux.HandleFunc("/api/v1/endpoints", endpoints)
	mux.HandleFunc("/api/v1/openapi.json", openapi)
	mux.HandleFunc("/api/v1/", resource)
	mux.HandleFunc("/api/run", run)
	if token == "" {
		return mux
	}
	expected := sha256.Sum256([]byte(token))
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		provided := sha256.Sum256([]byte(r.Header.Get("X-AGC-Server-Token")))
		if subtle.ConstantTimeCompare(expected[:], provided[:]) != 1 {
			writeJSON(w, http.StatusUnauthorized, domain.ErrorEnvelope{Error: domain.ErrorDetail{Code: "server_auth_required", Message: "valid X-AGC-Server-Token required"}})
			return
		}
		mux.ServeHTTP(w, r)
	})
}

func root(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/api/v1" {
		resource(w, r)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"name":         "agc-cli",
		"capabilities": domain.DecoratedCapabilities(),
		"_links": map[string]domain.Link{
			"capabilities": {Href: "/api/v1/capabilities", Method: "GET"},
			"openapi":      {Href: "/api/v1/openapi.json", Method: "GET"},
			"run":          {Href: "/api/run", Method: "POST"},
		},
	})
}

func capabilities(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, domain.Envelope[[]domain.Capability]{
		Data: domain.DecoratedCapabilities(),
		Affordances: map[string]string{
			"cli": "agc capabilities",
		},
	})
}

func endpoints(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, domain.Envelope[[]domain.Endpoint]{
		Data: domain.AllEndpoints(),
		Affordances: map[string]string{
			"cli": "agc endpoints",
		},
	})
}

func openapi(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, domain.OpenAPISpec())
}

func resource(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/api/v1/")
	id = strings.Trim(id, "/")
	if id == "" {
		root(w, r)
		return
	}
	segments := strings.Split(id, "/")
	if len(segments) >= 2 && segments[1] == "endpoints" {
		resourceEndpoints(w, r, segments)
		return
	}
	capability, ok := domain.CapabilityByID(id)
	if !ok {
		writeJSON(w, http.StatusNotFound, domain.ErrorEnvelope{Error: domain.ErrorDetail{Code: "not_found", Message: "unknown API family"}})
		return
	}
	writeJSON(w, http.StatusOK, domain.Envelope[domain.Capability]{
		Data: capability,
		Warnings: []string{
			"This module is scaffolded for the public AppGallery Connect API. Wire concrete endpoint adapters before using it against production data.",
		},
	})
}

func resourceEndpoints(w http.ResponseWriter, r *http.Request, segments []string) {
	familyID := segments[0]
	if _, ok := domain.CapabilityByID(familyID); !ok {
		writeJSON(w, http.StatusNotFound, domain.ErrorEnvelope{Error: domain.ErrorDetail{Code: "not_found", Message: "unknown API family"}})
		return
	}
	if len(segments) == 2 {
		writeJSON(w, http.StatusOK, domain.Envelope[[]domain.Endpoint]{
			Data: domain.EndpointsByFamily(familyID),
			Affordances: map[string]string{
				"cli": "agc " + familyID + " endpoints",
			},
		})
		return
	}
	if len(segments) == 4 && segments[3] == "invoke" {
		invokeEndpoint(w, r, familyID, segments[2])
		return
	}
	if len(segments) == 3 {
		endpoint, ok := domain.EndpointByID(familyID, segments[2])
		if !ok {
			writeJSON(w, http.StatusNotFound, domain.ErrorEnvelope{Error: domain.ErrorDetail{Code: "not_found", Message: "unknown endpoint"}})
			return
		}
		writeJSON(w, http.StatusOK, domain.Envelope[domain.Endpoint]{
			Data: endpoint,
		})
		return
	}
	writeJSON(w, http.StatusNotFound, domain.ErrorEnvelope{Error: domain.ErrorDetail{Code: "not_found", Message: "unknown endpoint route"}})
}

type invokeEndpointRequest struct {
	SessionFile string              `json:"sessionFile,omitempty"`
	ProjectID   string              `json:"projectId,omitempty"`
	FilePath    string              `json:"filePath,omitempty"`
	QueryValues map[string][]string `json:"queryValues,omitempty"`
	BaseURL     string              `json:"baseUrl,omitempty"`
	Params      map[string]string   `json:"params,omitempty"`
	Query       map[string]string   `json:"query,omitempty"`
	Headers     map[string]string   `json:"headers,omitempty"`
	Fields      map[string]string   `json:"fields,omitempty"`
	Body        json.RawMessage     `json:"body,omitempty"`
	Token       string              `json:"token,omitempty"`
	DryRun      *bool               `json:"dryRun,omitempty"`
}

func invokeEndpoint(w http.ResponseWriter, r *http.Request, familyID, endpointID string) {
	if r.Method != http.MethodPost {
		writeJSON(w, http.StatusMethodNotAllowed, domain.ErrorEnvelope{Error: domain.ErrorDetail{Code: "method_not_allowed", Message: "POST required"}})
		return
	}
	endpoint, ok := domain.EndpointByID(familyID, endpointID)
	if !ok {
		writeJSON(w, http.StatusNotFound, domain.ErrorEnvelope{Error: domain.ErrorDetail{Code: "not_found", Message: "unknown endpoint"}})
		return
	}
	if endpoint.Authentication == "agcli-user-oauth" {
		payload := domain.FunctionRequest{DryRun: true}
		dec := json.NewDecoder(r.Body)
		dec.DisallowUnknownFields()
		if dec.Decode(&payload) != nil {
			writeJSON(w, 400, domain.ErrorEnvelope{Error: domain.ErrorDetail{Code: "invalid_json", Message: "Invalid native function request"}})
			return
		}
		result, err := (&cloudfunctions.Client{}).Execute(r.Context(), endpointID, payload)
		if err != nil {
			detail := cloudfunctions.Feedback(endpointID, err)
			status := http.StatusBadGateway
			if detail.Code == "invalid_request" {
				status = http.StatusBadRequest
			}
			if detail.Code == "authentication_required" {
				status = http.StatusUnauthorized
			}
			writeJSON(w, status, domain.ErrorEnvelope{Error: detail})
			return
		}
		writeJSON(w, 200, domain.Envelope[cloudfunctions.Response]{Data: result})
		return
	}
	var payload invokeEndpointRequest
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		writeJSON(w, http.StatusBadRequest, domain.ErrorEnvelope{Error: domain.ErrorDetail{Code: "invalid_json", Message: "request body must be JSON", Details: err.Error()}})
		return
	}
	payload.Query = agcapi.ScalarQuery(payload.Query, payload.QueryValues)
	body, err := requestBody(payload)
	if len(payload.Body) == 0 && len(payload.Fields) > 0 {
		body, err = agcapi.MarshalEndpointFields(endpoint, payload.Fields)
	}
	if err != nil {
		writeJSON(w, http.StatusBadRequest, domain.ErrorEnvelope{Error: domain.ErrorDetail{Code: "invalid_body", Message: err.Error()}})
		return
	}
	if err := validateInvokeParameters(endpoint, payload.Params, payload.Query, payload.Headers, payload.Fields, body); err != nil {
		writeJSON(w, http.StatusBadRequest, domain.ErrorEnvelope{Error: domain.ErrorDetail{Code: "missing_parameter", Message: err.Error()}})
		return
	}
	validationErr := agcapi.ValidateEndpointRequest(endpoint, payload.Query, payload.Headers, body, payload.QueryValues)
	if validationErr == nil {
		validationErr = agcapi.ValidateEndpointPathInput(endpoint, payload.Params)
	}
	if validationErr == nil {
		validationErr = agcapi.ValidateEndpointFileInput(endpoint, payload.FilePath, payload.Headers)
	}
	if err := validationErr; err != nil {
		detail := domain.ErrorDetail{Code: "invalid_request", Message: err.Error(), Endpoint: endpoint.FamilyID + "/" + endpoint.ID, SourceURL: endpoint.SourceURL, Suggestion: "Correct the fields using the documented requirements before retrying.", DocumentationCommand: endpoint.Command + " --docs --pretty"}
		var invalid *agcapi.RequestValidationError
		if errors.As(err, &invalid) {
			detail.Violations = invalid.Violations
		}
		writeJSON(w, http.StatusBadRequest, domain.ErrorEnvelope{Error: detail})
		return
	}
	if err := agcapi.ValidateEndpointBody(endpoint, body); err != nil {
		writeJSON(w, http.StatusBadRequest, domain.ErrorEnvelope{Error: domain.ErrorDetail{Code: "invalid_body", Message: err.Error()}})
		return
	}
	dryRun := true
	if payload.DryRun != nil {
		dryRun = *payload.DryRun
	}
	ctx, cancel := context.WithTimeout(r.Context(), 60*time.Second)
	defer cancel()
	resp, err := agcapi.InvokeEndpoint(ctx, nil, agcapi.InvokeRequest{
		SessionFile: payload.SessionFile,
		ProjectID:   payload.ProjectID,
		FilePath:    payload.FilePath,
		Endpoint:    endpoint,
		BaseURL:     payload.BaseURL,
		Params:      payload.Params,
		Query:       payload.Query,
		QueryValues: payload.QueryValues,
		Headers:     payload.Headers,
		Body:        body,
		AccessToken: payload.Token,
		DryRun:      dryRun,
	})
	if err != nil {
		var invalid *agcapi.RequestValidationError
		if errors.As(err, &invalid) {
			writeJSON(w, 400, domain.ErrorEnvelope{Error: domain.ErrorDetail{Code: "invalid_request", Message: err.Error(), Endpoint: endpoint.FamilyID + "/" + endpoint.ID, SourceURL: endpoint.SourceURL, Violations: invalid.Violations, DocumentationCommand: endpoint.Command + " --docs --pretty", Suggestion: "Correct the affected field before invoking."}})
			return
		}
		writeInvokeError(w, resp, err)
		return
	}
	writeJSON(w, http.StatusOK, domain.Envelope[agcapi.InvokeResponse]{
		Data: resp,
		Affordances: map[string]string{
			"cli": endpoint.Command + " --invoke",
		},
	})
}

func requestBody(payload invokeEndpointRequest) ([]byte, error) {
	if len(payload.Body) > 0 {
		return payload.Body, nil
	}
	if len(payload.Fields) == 0 {
		return nil, nil
	}
	return json.Marshal(payload.Fields)
}

func validateInvokeParameters(endpoint domain.Endpoint, params, query, headers, fields map[string]string, body []byte) error {
	for _, parameter := range endpoint.Parameters {
		if !parameter.Required {
			continue
		}
		switch parameter.In {
		case "path":
			if params[parameter.Name] == "" {
				return endpointParamError("params", parameter.Name)
			}
		case "query":
			if query[parameter.Name] == "" {
				return endpointParamError("query", parameter.Name)
			}
		case "header":
			found := false
			for name, value := range headers {
				if strings.EqualFold(name, parameter.Name) && value != "" {
					found = true
				}
			}
			if !found {
				return endpointParamError("headers", parameter.Name)
			}
		case "body":
			var object map[string]json.RawMessage
			if json.Unmarshal(body, &object) != nil || len(object[parameter.Name]) == 0 || string(object[parameter.Name]) == "null" || (!parameter.AllowEmpty && string(object[parameter.Name]) == `""`) {
				return endpointParamError("fields", parameter.Name)
			}
		case "file":
			if domain.IsServerlessFamily(endpoint.FamilyID) {
				continue
			} // validated through explicit filePath
			if params[parameter.Name] == "" && fields[parameter.Name] == "" && len(body) == 0 {
				return endpointParamError("fields", parameter.Name)
			}
		}
	}
	return nil
}

func endpointParamError(location, name string) error {
	return fmt.Errorf("missing %s.%s", location, name)
}

func writeInvokeError(w http.ResponseWriter, resp agcapi.InvokeResponse, err error) {
	var authErr *agcapi.AuthenticationError
	if errors.As(err, &authErr) {
		writeJSON(w, 401, domain.ErrorEnvelope{Error: domain.ErrorDetail{Code: "authentication_failed", Message: authErr.Message}})
		return
	}
	if agcErr, ok := err.(agcapi.AGCError); ok {
		writeJSON(w, http.StatusBadGateway, domain.ErrorEnvelope{Error: domain.ErrorDetail{Code: agcErr.Code, CodeField: agcErr.CodeField, ResultCodes: agcErr.ResultCodes, Message: agcErr.Message, Details: resp.URL}})
		return
	}
	status := http.StatusBadGateway
	if strings.Contains(err.Error(), "missing required") || strings.Contains(err.Error(), "missing --") {
		status = http.StatusBadRequest
	}
	writeJSON(w, status, domain.ErrorEnvelope{Error: domain.ErrorDetail{Code: "invoke_error", Message: err.Error(), Details: resp.URL}})
}

func run(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSON(w, http.StatusMethodNotAllowed, domain.ErrorEnvelope{Error: domain.ErrorDetail{Code: "method_not_allowed", Message: "POST required"}})
		return
	}
	writeJSON(w, http.StatusAccepted, domain.Envelope[map[string]string]{
		Data: map[string]string{
			"status": "accepted",
		},
		Warnings: []string{"Command execution endpoint is intentionally inert until command sandboxing is implemented."},
	})
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}
