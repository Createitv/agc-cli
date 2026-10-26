package command

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestSkillsList(t *testing.T) {
	out, err := execute("skills", "list", "--pretty")
	if err != nil {
		t.Fatal(err)
	}
	var result struct {
		Data []struct {
			Name    string `json:"name"`
			Command string `json:"command"`
		} `json:"data"`
	}
	if err := json.Unmarshal([]byte(out), &result); err != nil {
		t.Fatal(err)
	}
	if len(result.Data) != 1 || result.Data[0].Name != "agc-cli" || result.Data[0].Command != "agc skills show agc-cli" {
		t.Fatalf("unexpected catalog: %s", out)
	}
}

func TestSkillsShow(t *testing.T) {
	out, err := execute("skills", "show", "agc-cli")
	if err != nil {
		t.Fatal(err)
	}
	var result struct {
		Data struct {
			Content string `json:"content"`
		} `json:"data"`
	}
	if err := json.Unmarshal([]byte(out), &result); err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{"name: agc-cli", "agc endpoints", "--dry-run=false"} {
		if !strings.Contains(result.Data.Content, want) {
			t.Fatalf("missing %q: %s", want, out)
		}
	}
	out, err = execute("skills", "show", "agc-cli", "--output", "markdown")
	if err != nil || !strings.HasPrefix(out, "---\nname: agc-cli") {
		t.Fatalf("markdown: %s, %v", out, err)
	}
}

func TestSkillsInvalidArguments(t *testing.T) {
	for _, args := range [][]string{{"skills", "list", "extra"}, {"skills", "show"}, {"skills", "show", "missing"}, {"skills", "show", "agc-cli", "extra"}, {"skills", "list", "--output", "invalid"}, {"skills", "show", "agc-cli", "--output", "invalid"}} {
		if _, err := execute(args...); err == nil {
			t.Fatalf("expected error for %v", args)
		}
	}
}

func TestSkillsDiscovery(t *testing.T) {
	for _, args := range [][]string{{"--help"}, {"skills", "--help"}, {"skills", "list", "--output", "table"}, {"skills", "list", "--output", "markdown"}} {
		out, err := execute(args...)
		if err != nil {
			t.Fatal(err)
		}
		if !strings.Contains(out, "skills") {
			t.Fatalf("missing skills discovery: %s", out)
		}
	}
}
