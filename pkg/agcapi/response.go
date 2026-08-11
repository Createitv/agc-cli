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

func ParseAGCError(statusCode int, data []byte) AGCError {
	errResp := AGCError{StatusCode: statusCode, Message: string(data)}
	if len(data) == 0 {
		errResp.Message = "empty error response"
		return errResp
	}
	if json.Valid(data) {
		errResp.RawBody = data
		var candidates []struct {
			Code    string `json:"code"`
			Message string `json:"message"`
			Msg     string `json:"msg"`
			Error   string `json:"error"`
		}
		var flat struct {
			Code    string `json:"code"`
			Message string `json:"message"`
			Msg     string `json:"msg"`
			Error   string `json:"error"`
			Err     *struct {
				Code    string `json:"code"`
				Message string `json:"message"`
			} `json:"errorDetail"`
		}
		if json.Unmarshal(data, &flat) == nil {
			candidates = append(candidates, struct {
				Code    string `json:"code"`
				Message string `json:"message"`
				Msg     string `json:"msg"`
				Error   string `json:"error"`
			}{Code: flat.Code, Message: flat.Message, Msg: flat.Msg, Error: flat.Error})
			if flat.Err != nil {
				candidates = append(candidates, struct {
					Code    string `json:"code"`
					Message string `json:"message"`
					Msg     string `json:"msg"`
					Error   string `json:"error"`
				}{Code: flat.Err.Code, Message: flat.Err.Message})
			}
		}
		for _, candidate := range candidates {
			if candidate.Code != "" {
				errResp.Code = candidate.Code
			}
			switch {
			case candidate.Message != "":
				errResp.Message = candidate.Message
			case candidate.Msg != "":
				errResp.Message = candidate.Msg
			case candidate.Error != "":
				errResp.Message = candidate.Error
			}
			if errResp.Code != "" || errResp.Message != string(data) {
				return errResp
			}
		}
	}
	return errResp
}
