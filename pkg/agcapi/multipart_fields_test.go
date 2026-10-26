package agcapi

import (
	"encoding/json"
	"github.com/Createitv/agc-cli/pkg/domain"
	"testing"
)

func TestMultipartFieldsEncodeDocumentedNumberedObjects(t *testing.T) {
	for _, tc := range []struct{ id, value string }{{"obbfile-uploadurl", `{"length":1,"sha256":"fixture"}`}, {"obbfile-compose", `{"partObjectId":"fixture","etag":"fixture"}`}} {
		e, _ := domain.EndpointByID("upload", tc.id)
		body, err := MarshalEndpointFields(e, map[string]string{"additionalProp1": tc.value})
		if err != nil {
			t.Fatal(err)
		}
		var fields map[string]any
		if json.Unmarshal(body, &fields) != nil {
			t.Fatal("invalid JSON")
		}
		if _, ok := fields["additionalProp1"].(map[string]any); !ok {
			t.Fatalf("numbered part was encoded as a string: %s", body)
		}
	}
}

func TestMultipartFieldsReturnSpecificMissingPartFields(t *testing.T) {
	for _, tc := range []struct{ id, value, field string }{{"obbfile-uploadurl", `{"sha256":"fixture"}`, "body.additionalProp1.length"}, {"obbfile-compose", `{"partObjectId":"fixture"}`, "body.additionalProp1.etag"}} {
		e, _ := domain.EndpointByID("upload", tc.id)
		_, err := MarshalEndpointFields(e, map[string]string{"additionalProp1": tc.value})
		problem, ok := err.(*RequestValidationError)
		if !ok || len(problem.Violations) == 0 || problem.Violations[0].Field != tc.field || problem.Violations[0].SourceURL == "" {
			t.Fatalf("missing actionable part feedback: %v", err)
		}
	}
}
