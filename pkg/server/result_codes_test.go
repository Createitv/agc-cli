package server

import (
	"encoding/json"
	"github.com/Createitv/agc-cli/pkg/agcapi"
	"github.com/Createitv/agc-cli/pkg/domain"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestRESTPreservesCodeOriginWithoutRawResponse(t *testing.T) {
	recorder := httptest.NewRecorder()
	err := agcapi.ParseAGCError(200, []byte(`{"rtnCode":2,"businessCode":0,"invitationCode":"raw-secret"}`))
	writeInvokeError(recorder, agcapi.InvokeResponse{URL: "https://example.com"}, err)
	var envelope domain.ErrorEnvelope
	if json.Unmarshal(recorder.Body.Bytes(), &envelope) != nil || envelope.Error.CodeField != "rtnCode" || envelope.Error.ResultCodes["businessCode"] != "0" {
		t.Fatal("REST discarded original result codes")
	}
	if strings.Contains(recorder.Body.String(), "raw-secret") {
		t.Fatal("REST leaked raw response")
	}
}
