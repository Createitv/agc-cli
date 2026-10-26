package command

import (
	"path/filepath"
	"strings"
	"testing"
)

func TestHvigorLocalCommandDoesNotRequireHuaweiAuthorization(t *testing.T) {
	t.Setenv("AGC_CONFIG_DIR", t.TempDir())
	_, err := execute("cicd", "build-hvigor", "--invoke", "--dry-run=false", "--credentials-path", filepath.Join(t.TempDir(), "absent.json"), "--field", "project="+t.TempDir())
	if err == nil || !strings.Contains(err.Error(), "build-profile.json5") {
		t.Fatalf("expected local project error without credential loading: %v", err)
	}
}

func TestHvigorDryRunDoesNotExecuteProject(t *testing.T) {
	t.Setenv("AGC_CONFIG_DIR", t.TempDir())
	output, err := execute("cicd", "build-hvigor", "--invoke", "--field", "project=/does/not/exist")
	if err != nil || !strings.Contains(output, `"dryRun":true`) {
		t.Fatalf("dry run should describe local build without executing: %s %v", output, err)
	}
}
