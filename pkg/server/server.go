package server

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/Createitv/agc-cli/pkg/agcapi"
	"github.com/Createitv/agc-cli/pkg/domain"
)

func Handler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/v1", root)
	mux.HandleFunc("/api/v1/capabilities", capabilities)
	mux.HandleFunc("/api/v1/endpoints", endpoints)
	mux.HandleFunc("/api/v1/openapi.json", openapi)
	mux.HandleFunc("/api/v1/", resource)
	mux.HandleFunc("/api/run", run)
	return mux
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
	BaseURL string            `json:"baseUrl,omitempty"`
	Params  map[string]string `json:"params,omitempty"`
	Query   map[string]string `json:"query,omitempty"`
	Headers map[string]string `json:"headers,omitempty"`
	Fields  map[string]string `json:"fields,omitempty"`
	Body    json.RawMessage   `json:"body,omitempty"`
	Token   string            `json:"token,omitempty"`
	DryRun  *bool             `json:"dryRun,omitempty"`
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
	var payload invokeEndpointRequest
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		writeJSON(w, http.StatusBadRequest, domain.ErrorEnvelope{Error: domain.ErrorDetail{Code: "invalid_json", Message: "request body must be JSON", Details: err.Error()}})
		return
	}
	body, err := requestBody(payload)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, domain.ErrorEnvelope{Error: domain.ErrorDetail{Code: "invalid_body", Message: err.Error()}})
		return
	}
	if err := validateInvokeParameters(endpoint, payload.Params, payload.Query, payload.Headers, payload.Fields, body); err != nil {
		writeJSON(w, http.StatusBadRequest, domain.ErrorEnvelope{Error: domain.ErrorDetail{Code: "missing_parameter", Message: err.Error()}})
		return
	}
	dryRun := true
	if payload.DryRun != nil {
		dryRun = *payload.DryRun
	}
	ctx, cancel := context.WithTimeout(r.Context(), 60*time.Second)
	defer cancel()
	resp, err := agcapi.InvokeEndpoint(ctx, nil, agcapi.InvokeRequest{
		Endpoint:    endpoint,
		BaseURL:     payload.BaseURL,
		Params:      payload.Params,
		Query:       payload.Query,
		Headers:     payload.Headers,
		Body:        body,
		AccessToken: payload.Token,
		DryRun:      dryRun,
	})
	if err != nil {
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
			if headers[parameter.Name] == "" {
				return endpointParamError("headers", parameter.Name)
			}
		case "body":
			if len(body) == 0 && fields[parameter.Name] == "" {
				return endpointParamError("fields", parameter.Name)
			}
		case "file":
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
	if agcErr, ok := err.(agcapi.AGCError); ok {
		writeJSON(w, http.StatusBadGateway, domain.ErrorEnvelope{Error: domain.ErrorDetail{Code: agcErr.Code, Message: agcErr.Message, Details: resp.URL}})
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
