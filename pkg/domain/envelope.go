package domain

type Envelope[T any] struct {
	Data        T                 `json:"data"`
	NextCursor  string            `json:"nextCursor,omitempty"`
	Affordances map[string]string `json:"affordances,omitempty"`
	Warnings    []string          `json:"warnings,omitempty"`
}

type ErrorEnvelope struct {
	Error ErrorDetail `json:"error"`
}

type ErrorDetail struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}
