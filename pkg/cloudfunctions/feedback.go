package cloudfunctions

import (
	"errors"

	"github.com/Createitv/agc-cli/pkg/agcapi"
	"github.com/Createitv/agc-cli/pkg/domain"
)

// Feedback never includes OAuth sessions, function configuration or raw upstream
// credential responses. Wrapped publish failures retain the upstream code.
func Feedback(operation string, err error) domain.ErrorDetail {
	detail := domain.ErrorDetail{Code: "function_transport_failed", Message: err.Error(), Endpoint: "functions/" + operation, SourceURL: "https://www.npmjs.com/package/@hw-agconnect/agc-cli/v/1.12.0", DocumentationCommand: "agc functions " + operation + " --docs --pretty", Suggestion: "Check the selected project service URL and transport; do not blindly retry a write."}
	var invalid *agcapi.RequestValidationError
	if errors.As(err, &invalid) {
		detail.Code = "invalid_request"
		detail.Violations = invalid.Violations
		detail.Suggestion = "Correct the affected field using the documentation command before retrying."
		for _, v := range invalid.Violations {
			if v.Field == "sessionFile" {
				detail.Code = "authentication_required"
				detail.Suggestion = "Explicitly select an unexpired user OAuth session for the dedicated project."
			}
		}
	}
	var upstream agcapi.AGCError
	if errors.As(err, &upstream) {
		detail.Code = upstream.Code
		if detail.Code == "" {
			detail.Code = "upstream_error"
		}
		detail.CodeField = upstream.CodeField
		detail.ResultCodes = upstream.ResultCodes
		detail.Suggestion = "Review the upstream role, function state and documented prerequisites; retain the revision after a successful update."
	}
	return detail
}
