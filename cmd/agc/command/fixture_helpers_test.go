package command

import (
	"encoding/json"
	"fmt"
	"github.com/Createitv/agc-cli/pkg/domain"
	"strings"
	"time"
)

// Registry-wide smoke tests need syntactically valid values; the independent
// official-reference test owns the method, location, and serialization oracle.
func parameterFixture(p domain.Parameter) string {
	if p.Type == "string" && p.Format == "date" {
		return time.Now().UTC().AddDate(0, 0, 31).Format("2006-01-02")
	}
	if p.EncodedJSON != nil {
		return parameterFixture(*p.EncodedJSON)
	}
	if p.Name == "currency" {
		return "CNY"
	}
	if p.In == "query" && (p.Name == "startTime" || p.Name == "endTime") && strings.Contains(p.Description, "YYYYMMDD") {
		return "20261001"
	}
	if len(p.Enum) > 0 {
		return p.Enum[0]
	}
	if p.Type == "integer" || p.Type == "number" {
		if p.Minimum != nil {
			return fmt.Sprint(*p.Minimum)
		}
		return "1"
	}
	if p.Type == "array" {
		switch p.ItemsType {
		case "object":
			if len(p.ItemsProperties) > 0 {
				fields := map[string]any{}
				for _, child := range p.ItemsProperties {
					raw := parameterFixture(child)
					var value any
					if child.Type == "string" {
						value = raw
					} else {
						_ = json.Unmarshal([]byte(raw), &value)
					}
					fields[child.Name] = value
				}
				encoded, _ := json.Marshal([]any{fields})
				return string(encoded)
			}
			return `[{"key":"fixture"}]`
		case "integer":
			return `[1]`
		default:
			return `["fixture"]`
		}
	}
	if p.Type == "object" {
		if len(p.Properties) > 0 {
			fields := map[string]any{}
			for _, child := range p.Properties {
				raw := parameterFixture(child)
				var value any
				if child.Type == "string" {
					value = raw
				} else {
					_ = json.Unmarshal([]byte(raw), &value)
				}
				fields[child.Name] = value
			}
			encoded, _ := json.Marshal(fields)
			return string(encoded)
		}
		return `{"key":"fixture"}`
	}
	return "fixture"
}
