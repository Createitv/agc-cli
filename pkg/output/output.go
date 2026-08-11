package output

import (
	"encoding/json"
	"fmt"
	"io"
	"reflect"
	"strings"
)

type Format string

const (
	JSON     Format = "json"
	Table    Format = "table"
	Markdown Format = "markdown"
)

func Write(w io.Writer, value any, format Format, pretty bool) error {
	switch format {
	case "", JSON:
		encoder := json.NewEncoder(w)
		if pretty {
			encoder.SetIndent("", "  ")
		}
		return encoder.Encode(value)
	case Table:
		return writeTable(w, value, false)
	case Markdown:
		return writeTable(w, value, true)
	default:
		return fmt.Errorf("unsupported output format %q", format)
	}
}

func writeTable(w io.Writer, value any, markdown bool) error {
	rows := rowsFromValue(value)
	if len(rows) == 0 {
		_, err := fmt.Fprintln(w, "No data")
		return err
	}
	headers := orderedKeys(rows[0])
	if markdown {
		fmt.Fprintf(w, "| %s |\n", strings.Join(headers, " | "))
		separators := make([]string, len(headers))
		for i := range separators {
			separators[i] = "---"
		}
		fmt.Fprintf(w, "| %s |\n", strings.Join(separators, " | "))
		for _, row := range rows {
			values := make([]string, len(headers))
			for i, header := range headers {
				values[i] = fmt.Sprint(row[header])
			}
			fmt.Fprintf(w, "| %s |\n", strings.Join(values, " | "))
		}
		return nil
	}
	fmt.Fprintln(w, strings.Join(headers, "\t"))
	for _, row := range rows {
		values := make([]string, len(headers))
		for i, header := range headers {
			values[i] = fmt.Sprint(row[header])
		}
		fmt.Fprintln(w, strings.Join(values, "\t"))
	}
	return nil
}

func rowsFromValue(value any) []map[string]any {
	encoded, _ := json.Marshal(value)
	var envelope struct {
		Data json.RawMessage `json:"data"`
	}
	if json.Unmarshal(encoded, &envelope) == nil && len(envelope.Data) > 0 {
		encoded = envelope.Data
	}
	var rows []map[string]any
	if json.Unmarshal(encoded, &rows) == nil {
		return rows
	}
	var row map[string]any
	if json.Unmarshal(encoded, &row) == nil {
		return []map[string]any{row}
	}
	if reflect.ValueOf(value).Kind() == reflect.Slice {
		return rows
	}
	return nil
}

func orderedKeys(row map[string]any) []string {
	preferred := []string{"id", "name", "description", "command", "restPath", "status"}
	keys := make([]string, 0, len(row))
	seen := map[string]bool{}
	for _, key := range preferred {
		if _, ok := row[key]; ok {
			keys = append(keys, key)
			seen[key] = true
		}
	}
	for key := range row {
		if !seen[key] && !strings.HasPrefix(key, "_") && key != "affordances" {
			keys = append(keys, key)
		}
	}
	return keys
}
