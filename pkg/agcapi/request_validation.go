package agcapi

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/Createitv/agc-cli/pkg/domain"
)

type RequestValidationError struct{ Violations []domain.InputViolation }

// ValidateEndpointPathInput shares the source parameter constraints with the
// JSON/query validator, preserving path field names in Agent feedback.
func ValidateEndpointPathInput(e domain.Endpoint, params map[string]string) error {
	var fields []domain.Parameter
	for _, p := range e.Parameters {
		if p.In == "path" {
			p.In = "query"
			fields = append(fields, p)
		}
	}
	if len(fields) == 0 {
		return nil
	}
	err := ValidateEndpointRequest(domain.Endpoint{Parameters: fields, SourceURL: e.SourceURL}, params, nil, nil)
	if invalid, ok := err.(*RequestValidationError); ok {
		for i := range invalid.Violations {
			invalid.Violations[i].Field = strings.Replace(invalid.Violations[i].Field, "query.", "params.", 1)
			invalid.Violations[i].Message = strings.ReplaceAll(invalid.Violations[i].Message, "query.", "params.")
		}
	}
	return err
}

func (e *RequestValidationError) Error() string {
	if len(e.Violations) == 0 {
		return "invalid request"
	}
	return e.Violations[0].Message
}

func ValidateEndpointRequest(e domain.Endpoint, query, headers map[string]string, body []byte, repeated ...map[string][]string) error {
	var object map[string]any
	if len(body) > 0 {
		decoder := json.NewDecoder(strings.NewReader(string(body)))
		decoder.UseNumber()
		if err := decoder.Decode(&object); err != nil {
			return err
		}
	}
	queryLists := map[string][]string{}
	if len(repeated) > 0 {
		queryLists = repeated[0]
	}
	arrayQuery := map[string]bool{}
	for _, p := range e.Parameters {
		if p.In == "query" && p.Type == "array" {
			arrayQuery[p.Name] = true
		}
	}
	value := func(path string) (any, bool) {
		pieces := strings.SplitN(path, ".", 2)
		if len(pieces) != 2 {
			return nil, false
		}
		switch pieces[0] {
		case "body":
			var current any = object
			for _, name := range strings.Split(pieces[1], ".") {
				parent, ok := current.(map[string]any)
				if !ok {
					return nil, false
				}
				current, ok = parent[name]
				if !ok {
					return nil, false
				}
			}
			return current, true
		case "query":
			name := pieces[1]
			if arrayQuery[name] {
				entries, ok := queryLists[name]
				if !ok {
					if text, present := query[name]; present {
						entries = []string{text}
						ok = true
					}
				}
				if !ok {
					return nil, false
				}
				values := make([]any, len(entries))
				for i, text := range entries {
					values[i] = text
				}
				return values, true
			}
			v, ok := query[name]
			return v, ok
		case "header":
			for name, v := range headers {
				if strings.EqualFold(name, pieces[1]) {
					return v, true
				}
			}
		}
		return nil, false
	}
	present := func(path string) bool {
		v, ok := value(path)
		if !ok || v == nil {
			return false
		}
		switch x := v.(type) {
		case string:
			return strings.TrimSpace(x) != ""
		case []any:
			return len(x) > 0
		}
		return true
	}
	var issues []domain.InputViolation
	if domain.IsServerlessFamily(e.FamilyID) {
		issues = append(issues, serverlessViolations(e, object)...)
		issues = append(issues, resumableViolations(e, headers)...)
	}
	add := func(code, field, message, source, suggestion string) {
		if source == "" {
			source = e.SourceURL
		}
		issues = append(issues, domain.InputViolation{Code: code, Field: field, Message: message, SourceURL: source, Suggestion: suggestion})
	}
	if e.DynamicBody != nil {
		if e.DynamicBody.Required && len(object) == 0 {
			add("required", "body", "A numbered part object is required", e.DynamicBody.SourceURL, "Supply additionalProp1 containing the documented part fields.")
		}
		for name, item := range object {
			field := "body." + name
			matched, err := regexp.MatchString(e.DynamicBodyKeyPattern, name)
			if err != nil || !matched {
				add("part_key", field, "Invalid numbered part key", e.DynamicBody.SourceURL, "Use additionalProp1, additionalProp2 and subsequent numbered part keys.")
				continue
			}
			child, ok := item.(map[string]any)
			if !ok {
				add("object_required", field, field+" must be an object", e.DynamicBody.SourceURL, "Supply a JSON object containing the documented part fields.")
				continue
			}
			encoded, _ := json.Marshal(child)
			if err := ValidateEndpointRequest(domain.Endpoint{Parameters: e.DynamicBody.Properties, SourceURL: e.DynamicBody.SourceURL}, nil, nil, encoded); err != nil {
				if nested, ok := err.(*RequestValidationError); ok {
					for _, issue := range nested.Violations {
						issue.Field = field + "." + strings.TrimPrefix(issue.Field, "body.")
						issues = append(issues, issue)
					}
				}
			}
		}
	}
	for _, rule := range e.Rules {
		if rule.UnlessPresent != "" && present(rule.UnlessPresent) {
			continue
		}
		active := true
		if rule.When != nil {
			v, ok := value(rule.When.Field)
			if !ok && rule.When.Default != "" {
				v, ok = rule.When.Default, true
			}
			active = ok && fmt.Sprint(v) == rule.When.Equals
		}
		if !active {
			continue
		}
		switch rule.Kind {
		case "csv_count":
			if !present(rule.Field) {
				continue
			}
			v, _ := value(rule.Field)
			entries := strings.Split(fmt.Sprint(v), ",")
			invalid := rule.MinItems > 0 && len(entries) < rule.MinItems || rule.MaxItems > 0 && len(entries) > rule.MaxItems
			for _, entry := range entries {
				invalid = invalid || strings.TrimSpace(entry) == ""
			}
			if invalid {
				add(rule.ID, rule.Field, rule.Message, rule.SourceURL, rule.Suggestion)
			}
		case "minimum_remaining_days":
			if !present(rule.Field) {
				continue
			}
			v, _ := value(rule.Field)
			date, err := time.Parse("2006-01-02", fmt.Sprint(v))
			if err != nil {
				continue
			} // The field's date-format check reports this.
			today := time.Now().UTC().Truncate(24 * time.Hour)
			if !date.After(today.AddDate(0, 0, rule.MinDays)) {
				add(rule.ID, rule.Field, rule.Message, rule.SourceURL, rule.Suggestion)
			}
		case "forbidden_if":
			if present(rule.Field) {
				add(rule.ID, rule.Field, rule.Message, rule.SourceURL, rule.Suggestion)
			}
		case "date_window":
			if len(rule.Fields) != 2 {
				continue
			}
			dates := make([]time.Time, 2)
			valid := true
			for i, field := range rule.Fields {
				raw, exists := value(field)
				if !exists {
					valid = false
					continue
				}
				date, err := time.Parse("20060102", fmt.Sprint(raw))
				if err != nil {
					add(rule.ID, field, rule.Message, rule.SourceURL, rule.Suggestion)
					valid = false
				} else {
					dates[i] = date
				}
			}
			if valid && (dates[1].Before(dates[0]) || dates[1].Sub(dates[0]) > time.Duration(rule.MaxDays)*24*time.Hour) {
				add(rule.ID, rule.Fields[1], rule.Message, rule.SourceURL, rule.Suggestion)
			}
		case "equal_length":
			if len(rule.Fields) != 2 {
				continue
			}
			length := func(field string) int {
				v, ok := value(field)
				if !ok {
					return 0
				}
				if items, ok := v.([]any); ok {
					return len(items)
				}
				return 1
			}
			if length(rule.Fields[0]) != length(rule.Fields[1]) {
				add(rule.ID, rule.Fields[1], rule.Message, rule.SourceURL, rule.Suggestion)
			}
		case "required_with":
			active := false
			for _, field := range rule.Fields {
				active = active || present(field)
			}
			if active && !present(rule.Field) {
				add(rule.ID, rule.Field, rule.Message, rule.SourceURL, rule.Suggestion)
			}
		case "required_if":
			if !present(rule.Field) {
				add(rule.ID, rule.Field, rule.Message, rule.SourceURL, rule.Suggestion)
			}
		case "at_most_one":
			count := 0
			for _, field := range rule.Fields {
				if present(field) {
					count++
				}
			}
			if count > 1 {
				add(rule.ID, strings.Join(rule.Fields, "|"), rule.Message, rule.SourceURL, rule.Suggestion)
			}
		case "at_least_one":
			found := false
			for _, field := range rule.Fields {
				found = found || present(field)
			}
			if !found {
				add(rule.ID, strings.Join(rule.Fields, "|"), rule.Message, rule.SourceURL, rule.Suggestion)
			}
		}
	}
	for _, p := range e.Parameters {
		if p.In != "body" && p.In != "query" && p.In != "header" {
			continue
		}
		path := p.In + "." + p.Name
		if p.In == "query" && p.Type != "array" && len(queryLists[p.Name]) > 1 {
			add("scalar_query", path, path+" accepts one query value", p.SourceURL, "Provide this scalar query key once; repeat only documented array parameters.")
		}
		v, ok := value(path)
		if !ok || v == nil {
			if p.Required {
				add("required", path, path+" is required", p.SourceURL, "Supply this required field in the documented location.")
			}
			continue
		}
		if p.In == "body" {
			valid := true
			switch p.Type {
			case "string":
				text, ok := v.(string)
				valid = ok && (!p.Required || p.AllowEmpty || strings.TrimSpace(text) != "")
			case "array":
				_, valid = v.([]any)
			case "object":
				_, valid = v.(map[string]any)
			case "integer", "number":
				_, valid = v.(json.Number)
			case "boolean":
				_, valid = v.(bool)
			}
			if !valid {
				add("type", path, path+" must be "+p.Type, p.SourceURL, "Use the documented JSON type for this field.")
				continue
			}
		}
		if p.EncodedJSON != nil {
			text, ok := v.(string)
			var decoded any
			decoder := json.NewDecoder(strings.NewReader(text))
			decoder.UseNumber()
			if !ok || decoder.Decode(&decoded) != nil || !json.Valid([]byte(text)) {
				add("encoded_json", path, path+" must contain valid JSON", p.SourceURL, "Supply the documented JSON array as a string.")
			} else {
				model := *p.EncodedJSON
				model.Name = "value"
				model.In = "body"
				encoded, _ := json.Marshal(map[string]any{"value": decoded})
				if err := ValidateEndpointRequest(domain.Endpoint{Parameters: []domain.Parameter{model}, SourceURL: p.SourceURL}, nil, nil, encoded); err != nil {
					if nested, ok := err.(*RequestValidationError); ok {
						for _, issue := range nested.Violations {
							issue.Field = path + strings.TrimPrefix(issue.Field, "body.value")
							issues = append(issues, issue)
						}
					}
				}
			}
		}
		if child, ok := v.(map[string]any); ok {
			if (p.MinProperties > 0 && len(child) < p.MinProperties) || (p.MaxProperties > 0 && len(child) > p.MaxProperties) {
				add("object_length", path, path+" violates documented map entry bounds", p.SourceURL, "Supply the documented number of dynamic map entries.")
			}
			for key, value := range child {
				fieldPath := path + "[" + strconv.Quote(key) + "]"
				if p.PropertyNamesPattern != "" {
					matched, err := regexp.MatchString(p.PropertyNamesPattern, key)
					if err != nil || !matched {
						add("property_name", fieldPath, "Dynamic map key violates the documented format", p.SourceURL, "Use a documented version or identifier as the map key.")
					}
				}
				if p.AdditionalProperties != nil {
					model := *p.AdditionalProperties
					model.In = "body"
					model.Name = "value"
					data, _ := json.Marshal(map[string]any{"value": value})
					if err := ValidateEndpointRequest(domain.Endpoint{Parameters: []domain.Parameter{model}}, nil, nil, data); err != nil {
						if invalid, ok := err.(*RequestValidationError); ok {
							for _, issue := range invalid.Violations {
								issue.Field = fieldPath + strings.TrimPrefix(issue.Field, "body.value")
								issues = append(issues, issue)
							}
						}
					}
				}
			}
		}
		if child, ok := v.(map[string]any); ok && len(p.Properties) > 0 {
			encoded, _ := json.Marshal(child)
			if err := ValidateEndpointRequest(domain.Endpoint{Parameters: p.Properties, SourceURL: p.SourceURL}, nil, nil, encoded); err != nil {
				if nested, ok := err.(*RequestValidationError); ok {
					for _, issue := range nested.Violations {
						issue.Field = path + "." + strings.TrimPrefix(issue.Field, "body.")
						issues = append(issues, issue)
					}
				}
			}
		}
		if text, ok := v.(string); ok {
			if p.Format == "date" {
				if _, err := time.Parse("2006-01-02", text); err != nil {
					add("date_format", path, path+" must be a valid yyyy-MM-dd date", p.SourceURL, "Use the documented yyyy-MM-dd format with a real calendar date.")
				}
			}
			if p.Pattern != "" {
				matched, err := regexp.MatchString(p.Pattern, text)
				if err != nil || !matched {
					add("pattern", path, path+" violates the documented field format", p.SourceURL, "Use a value matching the documented pattern: "+p.Pattern)
				}
			}
			length := utf8.RuneCountInString(text)
			if p.MaxLength > 0 && length > p.MaxLength {
				add("max_length", path, fmt.Sprintf("%s must contain at most %d characters", path, p.MaxLength), p.SourceURL, "Shorten this field to the documented maximum.")
			}
			if p.MinLength > 0 && length < p.MinLength && !(p.AllowEmpty && length == 0) {
				add("min_length", path, fmt.Sprintf("%s must contain at least %d characters when supplied", path, p.MinLength), p.SourceURL, "Provide the documented minimum length or omit the optional field.")
			}
		}
		if items, ok := v.([]any); ok {
			for index, item := range items {
				valid := true
				switch p.ItemsType {
				case "string":
					_, valid = item.(string)
				case "object":
					_, valid = item.(map[string]any)
				case "integer":
					number, ok := item.(json.Number)
					valid = ok
					if ok {
						_, err := number.Int64()
						valid = err == nil
					}
				case "number":
					_, valid = item.(json.Number)
				case "boolean":
					_, valid = item.(bool)
				}
				itemPath := fmt.Sprintf("%s[%d]", path, index)
				if valid && len(p.ItemsEnum) > 0 {
					allowed := false
					for _, candidate := range p.ItemsEnum {
						allowed = allowed || fmt.Sprint(item) == candidate
					}
					if !allowed {
						add("item_enum", itemPath, itemPath+" must be one of "+strings.Join(p.ItemsEnum, ", "), p.SourceURL, "Use a documented array item value: "+strings.Join(p.ItemsEnum, ", "))
					}
				}
				if !valid {
					add("item_type", itemPath, itemPath+" must be "+p.ItemsType, p.SourceURL, "Use the documented JSON type for every array item.")
				}
				if text, ok := item.(string); ok && p.ItemsPattern != "" {
					matched, err := regexp.MatchString(p.ItemsPattern, text)
					if err != nil || !matched {
						add("item_pattern", itemPath, itemPath+" violates the documented identifier format", p.SourceURL, "Use an identifier matching the documented pattern: "+p.ItemsPattern)
					}
				}
				if text, ok := item.(string); ok && p.ItemsMaxLength > 0 && utf8.RuneCountInString(text) > p.ItemsMaxLength {
					add("item_max_length", itemPath, itemPath+" exceeds the documented character limit", p.SourceURL, fmt.Sprintf("Use at most %d characters for this array item.", p.ItemsMaxLength))
				}
			}
			if len(p.ItemsProperties) > 0 {
				for index, item := range items {
					child, ok := item.(map[string]any)
					if !ok {
						add("object_required", fmt.Sprintf("%s[%d]", path, index), path+" entries must be objects", p.SourceURL, "Supply a JSON object for each array entry.")
						continue
					}
					encoded, _ := json.Marshal(child)
					childEndpoint := domain.Endpoint{Parameters: p.ItemsProperties, Rules: p.ItemsRules, SourceURL: p.SourceURL}
					if err := ValidateEndpointRequest(childEndpoint, nil, nil, encoded); err != nil {
						if nested, ok := err.(*RequestValidationError); ok {
							for _, issue := range nested.Violations {
								issue.Field = fmt.Sprintf("%s[%d].%s", path, index, strings.TrimPrefix(issue.Field, "body."))
								issues = append(issues, issue)
							}
						}
					}
				}
			}
			if (p.MaxItems > 0 && len(items) > p.MaxItems) || (p.MinItems > 0 && len(items) < p.MinItems) {
				add("array_length", path, fmt.Sprintf("%s violates documented array length (%d..%d)", path, p.MinItems, p.MaxItems), p.SourceURL, "Adjust the number of array entries to the documented bounds.")
			}
		}
		if p.Type == "integer" {
			number, err := strconv.ParseInt(fmt.Sprint(v), 10, 64)
			if err != nil {
				add("integer_required", path, path+" must be an integer", p.SourceURL, "Use a base-10 integer.")
				continue
			}
			if (p.Minimum != nil && number < *p.Minimum) || (p.Maximum != nil && number > *p.Maximum) {
				add("numeric_range", path, path+" is outside the documented range", p.SourceURL, "Use a value within the minimum and maximum listed by --docs.")
			}
		}
		if len(p.Enum) > 0 {
			valid := false
			for _, allowed := range p.Enum {
				valid = valid || fmt.Sprint(v) == allowed
			}
			if !valid {
				add("enum", path, path+" must be one of "+strings.Join(p.Enum, ", "), p.SourceURL, "Use a documented value: "+strings.Join(p.Enum, ", "))
			}
		}
	}
	issues = append(issues, validateManagementContent(e, query, body, object)...)
	if len(issues) > 0 {
		return &RequestValidationError{Violations: issues}
	}
	return nil
}
