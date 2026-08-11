package output

import (
	"bytes"
	"strings"
	"testing"

	"github.com/your-org/agc-cli/pkg/domain"
)

func TestWriteJSON(t *testing.T) {
	var buf bytes.Buffer
	err := Write(&buf, domain.Envelope[map[string]string]{Data: map[string]string{"id": "publishing"}}, JSON, false)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(buf.String(), `"id":"publishing"`) {
		t.Fatalf("json output = %s", buf.String())
	}
}

func TestWriteMarkdown(t *testing.T) {
	var buf bytes.Buffer
	err := Write(&buf, domain.Envelope[[]map[string]string]{Data: []map[string]string{{"id": "publishing", "status": "planned"}}}, Markdown, false)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(buf.String(), "| id | status |") {
		t.Fatalf("markdown output = %s", buf.String())
	}
}

func TestWriteTable(t *testing.T) {
	var buf bytes.Buffer
	err := Write(&buf, domain.Envelope[[]map[string]string]{Data: []map[string]string{{"id": "publishing", "status": "planned"}}}, Table, false)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(buf.String(), "id\tstatus") {
		t.Fatalf("table output = %s", buf.String())
	}
}

func TestWriteNoData(t *testing.T) {
	var buf bytes.Buffer
	if err := Write(&buf, []string{}, Table, false); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(buf.String(), "No data") {
		t.Fatalf("output = %s", buf.String())
	}
}

func TestWriteRejectsUnknownFormat(t *testing.T) {
	err := Write(&bytes.Buffer{}, map[string]string{}, Format("xml"), false)
	if err == nil {
		t.Fatal("expected unsupported format error")
	}
}
