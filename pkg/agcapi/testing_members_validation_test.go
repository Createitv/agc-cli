package agcapi

import (
	"encoding/json"
	"github.com/Createitv/agc-cli/pkg/domain"
	"strings"
	"testing"
)

func TestTestingMemberDeletionRejectsWrongOrOversizedIDArrays(t *testing.T) {
	e, _ := domain.EndpointByID("testing", "test-api-delete-test-group-user")
	ids := make([]string, 101)
	for i := range ids {
		ids[i] = "owned-fixture"
	}
	oversized, _ := json.Marshal(map[string]any{"groupId": "owned-group", "deleteTesterIds": ids})
	for _, body := range [][]byte{[]byte(`{"groupId":"owned-group","deleteTesterIds":"member"}`), []byte(`{"groupId":"owned-group","deleteTesterIds":[1]}`), oversized, []byte(`{"groupId":"owned-group","deleteTesterIds":["` + strings.Repeat("x", 101) + `"]}`), []byte(`{"deleteTesterIds":["member"]}`)} {
		if err := ValidateEndpointRequest(e, nil, map[string]string{"appId": "fixture"}, body); err == nil {
			t.Errorf("invalid official request accepted: %s", body)
		}
	}
	if err := ValidateEndpointRequest(e, nil, map[string]string{"appId": "fixture"}, []byte(`{"groupId":"owned-group","deleteTesterIds":["member"]}`)); err != nil {
		t.Fatal(err)
	}
}

func TestMemberFieldsReturnStructuredArrayIndex(t *testing.T) {
	e, _ := domain.EndpointByID("testing", "test-api-delete-test-group-user")
	_, err := MarshalEndpointFields(e, map[string]string{"groupId": "fixture", "deleteTesterIds": "[1]"})
	problem, ok := err.(*RequestValidationError)
	if !ok || len(problem.Violations) == 0 || problem.Violations[0].Field != "body.deleteTesterIds[0]" {
		t.Fatalf("missing structured array index: %v", err)
	}
}
