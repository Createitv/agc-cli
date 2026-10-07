package agcapi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/Createitv/agc-cli/pkg/domain"
)

// MarshalEndpointFields 将 CLI/REST 的字符串输入按已经核对的业务类型编码。
// 未声明类型的旧接口维持字符串行为；嵌套 JSON 推荐使用完整 body。
func MarshalEndpointFields(endpoint domain.Endpoint, fields map[string]string) ([]byte, error) {
	values := map[string]json.RawMessage{}
	for name, value := range fields {
		encoded, _ := json.Marshal(value)
		for _, parameter := range endpoint.Parameters {
			if parameter.Name == name && parameter.In == "body" && parameter.Type != "" && parameter.Type != "string" {
				encoded = []byte(value)
			}
		}
		values[name] = encoded
	}
	encoded, err := json.Marshal(values)
	if err != nil {
		return nil, fmt.Errorf("body fields must contain valid JSON values")
	}
	if err := ValidateEndpointBody(endpoint, encoded); err != nil {
		return nil, err
	}
	return encoded, nil
}

// ValidateEndpointBody 校验已知字段类型，避免 JSON 文件绕过整数/字符串约束。
// 不把尚未核对的嵌套模型臆测成完整 schema。
func ValidateEndpointBody(endpoint domain.Endpoint, body []byte) error {
	if len(body) == 0 {
		for _, parameter := range endpoint.Parameters {
			if parameter.In == "body" && parameter.Type != "" && parameter.Required {
				return fmt.Errorf("missing required body field %s", parameter.Name)
			}
		}
		return nil
	}
	var fields map[string]json.RawMessage
	if err := json.Unmarshal(body, &fields); err != nil || fields == nil {
		return fmt.Errorf("body must be a JSON object")
	}
	for _, parameter := range endpoint.Parameters {
		if parameter.In != "body" || parameter.Type == "" {
			continue
		}
		raw, ok := fields[parameter.Name]
		if !ok || bytes.Equal(bytes.TrimSpace(raw), []byte("null")) {
			if parameter.Required {
				return fmt.Errorf("missing required body field %s", parameter.Name)
			}
			continue
		}
		decoder := json.NewDecoder(bytes.NewReader(raw))
		decoder.UseNumber()
		var value any
		if err := decoder.Decode(&value); err != nil {
			return fmt.Errorf("body field %s must be %s", parameter.Name, parameter.Type)
		}
		valid := false
		switch parameter.Type {
		case "string":
			text, ok := value.(string)
			valid = ok && (!parameter.Required || strings.TrimSpace(text) != "")
		case "integer":
			if number, ok := value.(json.Number); ok {
				_, err := number.Int64()
				valid = err == nil
			}
		case "number":
			_, valid = value.(json.Number)
		case "boolean":
			_, valid = value.(bool)
		case "array":
			items, ok := value.([]any)
			valid = ok && len(items) >= parameter.MinItems && (parameter.MaxItems == 0 || len(items) <= parameter.MaxItems)
			if valid && parameter.ItemsType == "string" {
				for _, item := range items {
					if _, ok := item.(string); !ok {
						valid = false
					}
				}
			}
		case "object":
			_, valid = value.(map[string]any)
		default:
			valid = true
		}
		if !valid {
			return fmt.Errorf("body field %s must be %s", parameter.Name, parameter.Type)
		}
	}
	return nil
}
