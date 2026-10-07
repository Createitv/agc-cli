package domain

import "testing"

func TestProvisioningLifecycleRequiredContracts(t *testing.T) {
	cases := []struct{ id, name, location, kind string }{
		{"provision-api-add-fingerprints", "appId", "header", ""},
		{"provision-api-add-fingerprints", "fingerprintList", "body", "array"},
		{"provision-api-delete-fingerprints", "fingerprintList", "body", "array"},
		{"provision-api-apply-cent", "csr", "body", "string"},
		{"provision-api-apply-cent", "certName", "body", "string"},
		{"provision-api-apply-cent", "certType", "body", "integer"},
		{"provision-api-delete-cent", "certIds", "body", "array"},
		{"provision-api-apply-provision", "provisionName", "body", "string"},
		{"provision-api-apply-provision", "provisionType", "body", "integer"},
		{"provision-api-apply-provision", "certId", "body", "string"},
		{"provision-api-apply-provision", "appId", "body", "string"},
		{"provision-api-delete-provision", "id", "query", "string"},
		{"provision-api-update-provision", "provisionId", "body", "string"},
		{"provision-api-update-provision", "deviceIdList", "body", "array"},
		{"provision-api-query-provision", "appId", "header", ""},
	}
	for _, c := range cases {
		t.Run(c.id+"/"+c.name, func(t *testing.T) {
			ep, _ := EndpointByID("provisioning", c.id)
			for _, p := range ep.Parameters {
				if p.Name == c.name && p.In == c.location && p.Required && p.Type == c.kind {
					return
				}
			}
			t.Fatalf("缺少官方字段 %s@%s，类型 %s", c.name, c.location, c.kind)
		})
	}
}
