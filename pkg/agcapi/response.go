package agcapi

import (
	"encoding/json"
	"fmt"
)

type AGCError struct {
	StatusCode int             `json:"statusCode"`
	Code       string          `json:"code,omitempty"`
	Message    string          `json:"message"`
	RawBody    json.RawMessage `json:"rawBody,omitempty"`
}

func (e AGCError) Error() string {
	if e.Code != "" {
		return fmt.Sprintf("agc endpoint returned HTTP %d (%s): %s", e.StatusCode, e.Code, e.Message)
	}
	return fmt.Sprintf("agc endpoint returned HTTP %d: %s", e.StatusCode, e.Message)
}

type PaginatedResponse[T any] struct {
	Data       []T    `json:"data"`
	NextCursor string `json:"nextCursor,omitempty"`
	Total      int    `json:"total,omitempty"`
}

// responseResults recognizes the result envelopes used by AGC services. A missing
// or null code is not a business failure (downloads may not be JSON at all).
func responseResults(data []byte) []responseResult {
	var root map[string]json.RawMessage
	if json.Unmarshal(data, &root) != nil {
		return nil
	}
	var results []responseResult
	add := func(obj map[string]json.RawMessage, codeKey string, messageKeys ...string) {
		var code string
		raw := obj[codeKey]
		if json.Unmarshal(raw, &code) != nil {
			var number json.Number
			if json.Unmarshal(raw, &number) == nil {
				code = number.String()
			}
		}
		var message string
		for _, key := range messageKeys {
			if json.Unmarshal(obj[key], &message) == nil && message != "" {
				break
			}
		}
		results = append(results, responseResult{code: code, message: message})
	}
	nested := func(key string) map[string]json.RawMessage {
		var obj map[string]json.RawMessage
		_ = json.Unmarshal(root[key], &obj)
		return obj
	}
	add(nested("ret"), "code", "msg", "message")
	add(root, "rtnCode", "rtnDesc", "rtnMsg", "message", "msg")
	add(nested("error"), "errorCode", "errorMsg", "errorMessage", "message", "msg")
	add(root, "code", "message", "msg", "error")
	add(nested("errorDetail"), "code", "message", "msg")
	return results
}

type responseResult struct{ code, message string }

func businessFailure(data []byte) bool {
	for _, result := range responseResults(data) {
		if result.code != "" && result.code != "0" {
			return true
		}
	}
	return false
}

func ParseAGCError(statusCode int, data []byte) AGCError {
	errResp := AGCError{StatusCode: statusCode, Message: string(data)}
	if len(data) == 0 {
		errResp.Message = "empty error response"
		return errResp
	}
	if json.Valid(data) {
		errResp.RawBody = append(json.RawMessage(nil), data...)
	}
	results := responseResults(data)
	// Prefer the failing result when multiple envelopes occur in one response.
	for _, result := range results {
		if result.code != "" && result.code != "0" {
			errResp.Code = result.code
			if result.message != "" {
				errResp.Message = result.message
			}
			return errResp
		}
	}
	for _, result := range results {
		if result.code != "" || result.message != "" {
			errResp.Code = result.code
			if result.message != "" {
				errResp.Message = result.message
			}
			return errResp
		}
	}
	return errResp
}
