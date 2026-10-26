package agcapi

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/Createitv/agc-cli/pkg/domain"
)

func serverlessViolations(e domain.Endpoint, object map[string]any) []domain.InputViolation {
	var issues []domain.InputViolation
	add := func(field, msg string) {
		issues = append(issues, domain.InputViolation{Code: "serverless_contract", Field: field, Message: msg, SourceURL: e.SourceURL, Suggestion: "Correct this field using " + e.Command + " --docs --pretty before authorization or network calls."})
	}
	if e.FamilyID == "cloud-auth" && e.ID == "import-users" {
		users, _ := object["users"].([]any)
		for i, u := range users {
			user, _ := u.(map[string]any)
			providers, _ := user["providers"].([]any)
			for _, v := range providers {
				provider, _ := v.(map[string]any)
				p, _ := provider["provider"].(string)
				prefix := fmt.Sprintf("body.users[%d].", i)
				if p == "12" {
					if user["email"] == nil || user["email"] == "" {
						add(prefix+"email", "Email provider 12 requires email.")
					}
					if _, ok := user["emailVerified"].(bool); !ok {
						add(prefix+"emailVerified", "Email provider 12 requires boolean emailVerified.")
					}
				}
				if p == "11" && (user["phone"] == nil || user["phone"] == "") {
					add(prefix+"phone", "Phone provider 11 requires phone.")
				}
			}
		}
	}
	if e.FamilyID == "clouddb" && e.ID == "query-data" {
		query, _ := object["queryReqMsg"].(map[string]any)
		kind := fmt.Sprint(query["queryType"])
		if allowed, ok := map[string][]string{"1": {"AVG", "Average"}, "2": {"SUM"}, "3": {"MAX"}, "4": {"MIN"}, "5": {"COUNT"}}[kind]; ok {
			var cond struct {
				Conditions []struct {
					Type string `json:"conditionType"`
				} `json:"queryConditions"`
			}
			encoded, _ := query["queryCond"].(string)
			if json.Unmarshal([]byte(encoded), &cond) == nil {
				found := false
				for _, c := range cond.Conditions {
					for _, a := range allowed {
						found = found || strings.EqualFold(c.Type, a)
					}
				}
				if !found {
					add("body.queryReqMsg.queryCond", "Aggregate queryType requires a matching conditionType.")
				}
			}
		}
	}
	if e.FamilyID == "clouddb" && (e.ID == "update-data" || e.ID == "delete-data") {
		schemas, _ := object["schemas"].([]any)
		ops, _ := object["opData"].([]any)
		if len(schemas) != len(ops) {
			add("body.opData", "schemas and opData must align by object type.")
		}
		for i, s := range schemas {
			schema, _ := s.(map[string]any)
			fields, _ := schema["fs"].([]any)
			version, deleted := -1, -1
			for j, f := range fields {
				field, _ := f.(map[string]any)
				switch field["n"] {
				case "naturalbase_version":
					version = j
				case "naturalbase_deleted":
					deleted = j
				}
			}
			if version < 0 {
				add(fmt.Sprintf("body.schemas[%d].fs", i), "Write schemas require naturalbase_version (TYPE_LONG).")
			}
			if deleted < 0 {
				add(fmt.Sprintf("body.schemas[%d].fs", i), "Write schemas require naturalbase_deleted (TYPE_BOOLEAN).")
			}
			if i >= len(ops) {
				continue
			}
			op, _ := ops[i].(map[string]any)
			rows, _ := op["os"].([]any)
			for j, row := range rows {
				r, _ := row.(map[string]any)
				values, _ := r["fs"].([]any)
				prefix := fmt.Sprintf("body.opData[%d].os[%d].fs", i, j)
				if len(values) != len(fields) {
					add(prefix, "Field values must match schema field order and count.")
					continue
				}
				if deleted >= 0 {
					v, _ := values[deleted].(map[string]any)
					want := e.ID == "delete-data"
					if flag, ok := v["bl"].(bool); !ok || flag != want {
						add(fmt.Sprintf("%s[%d].bl", prefix, deleted), fmt.Sprintf("naturalbase_deleted must be %t for this operation.", want))
					}
				}
			}
		}
	}
	return issues
}
