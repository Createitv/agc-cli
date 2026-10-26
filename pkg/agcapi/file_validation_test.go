package agcapi

import (
	"errors"
	"path/filepath"
	"strings"
	"testing"

	"github.com/Createitv/agc-cli/pkg/domain"
)

func TestResourceFileErrorsDoNotExposeLocalPaths(t *testing.T) {
	e, _ := domain.EndpointByID("pms", "resourcemanage-harmonyosnext")
	dir := t.TempDir()
	for _, path := range []string{dir, filepath.Join(dir, "private-missing-file.png")} {
		err := ValidateEndpointFileInput(e, path)
		var invalid *RequestValidationError
		if !errors.As(err, &invalid) || len(invalid.Violations) != 1 {
			t.Fatalf("unclassified local file failure: %v", err)
		}
		v := invalid.Violations[0]
		if v.Code != "invalid_file" || v.Field != "filePath" || v.SourceURL == "" || strings.Contains(v.Message, path) {
			t.Fatalf("wrong file feedback: %+v", v)
		}
	}
}
