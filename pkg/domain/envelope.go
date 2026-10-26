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
	CodeField            string            `json:"codeField,omitempty"`
	ResultCodes          map[string]string `json:"resultCodes,omitempty"`
	Endpoint             string            `json:"endpoint,omitempty"`
	SourceURL            string            `json:"sourceUrl,omitempty"`
	Suggestion           string            `json:"suggestion,omitempty"`
	DocumentationCommand string            `json:"documentationCommand,omitempty"`
	Violations           []InputViolation  `json:"violations,omitempty"`
	Code                 string            `json:"code"`
	Message              string            `json:"message"`
	Details              string            `json:"details,omitempty"`
}
