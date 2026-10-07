package hvigor

import (
	"archive/zip"
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestBuildReturnsOnlySelectedProductOutputArtifacts(t *testing.T) {
	project := t.TempDir()
	if err := os.WriteFile(filepath.Join(project, "build-profile.json5"), []byte("{}"), 0600); err != nil {
		t.Fatal(err)
	}
	executable := filepath.Join(project, "hvigorw")
	if err := os.WriteFile(executable, []byte("#!/bin/sh\nexit 0\n"), 0700); err != nil {
		t.Fatal(err)
	}
	selected := filepath.Join(project, "entry", "build", "default", "outputs", "default", "entry.hap")
	for _, path := range []string{selected, filepath.Join(project, "entry", "build", "other", "outputs", "other", "entry.hap"), filepath.Join(project, "unrelated.hap")} {
		if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
			t.Fatal(err)
		}
		file, err := os.Create(path)
		if err != nil {
			t.Fatal(err)
		}
		archive := zip.NewWriter(file)
		entry, err := archive.Create("module.json")
		if err != nil {
			t.Fatal(err)
		}
		if _, err := entry.Write([]byte("{}")); err != nil {
			t.Fatal(err)
		}
		if err := archive.Close(); err != nil {
			t.Fatal(err)
		}
		if err := file.Close(); err != nil {
			t.Fatal(err)
		}
	}
	result, err := Build(context.Background(), Request{Project: project, Executable: executable})
	if err != nil {
		t.Fatal(err)
	}
	if len(result.Artifacts) != 1 || result.Artifacts[0].Path != selected {
		t.Fatalf("unexpected product artifacts: %+v", result.Artifacts)
	}
}

func TestAppBuildSelectsProjectAPPInsteadOfModuleHAP(t *testing.T) {
	project := t.TempDir()
	if err := os.WriteFile(filepath.Join(project, "build-profile.json5"), []byte("{}"), 0600); err != nil {
		t.Fatal(err)
	}
	executable := filepath.Join(project, "hvigorw")
	if err := os.WriteFile(executable, []byte("#!/bin/sh\nexit 0\n"), 0700); err != nil {
		t.Fatal(err)
	}
	for _, relative := range []string{"build/outputs/default/fixture.app", "entry/build/default/outputs/default/entry.hap", "build/outputs/other/fixture.app"} {
		path := filepath.Join(project, relative)
		if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
			t.Fatal(err)
		}
		file, err := os.Create(path)
		if err != nil {
			t.Fatal(err)
		}
		archive := zip.NewWriter(file)
		entry, err := archive.Create("pack.info")
		if err != nil {
			t.Fatal(err)
		}
		if _, err := entry.Write([]byte("{}")); err != nil {
			t.Fatal(err)
		}
		if err := archive.Close(); err != nil {
			t.Fatal(err)
		}
		if err := file.Close(); err != nil {
			t.Fatal(err)
		}
	}
	// JSON 输入与 CLI 一样，避免测试直接依赖新增 Go 字段而无法先复现旧行为。
	var request Request
	if err := json.Unmarshal([]byte(`{"task":"assembleApp"}`), &request); err != nil {
		t.Fatal(err)
	}
	request.Project, request.Executable = project, executable
	result, err := Build(context.Background(), request)
	if err != nil {
		t.Fatal(err)
	}
	if len(result.Artifacts) != 1 || !strings.HasSuffix(result.Artifacts[0].Path, "fixture.app") {
		t.Fatalf("APP build returned wrong artifacts: %+v", result.Artifacts)
	}
}

func TestBuildDoesNotAcceptSuccessfulProcessWithoutArtifact(t *testing.T) {
	project := t.TempDir()
	if err := os.WriteFile(filepath.Join(project, "build-profile.json5"), []byte("{}"), 0600); err != nil {
		t.Fatal(err)
	}
	executable := filepath.Join(project, "hvigorw")
	if err := os.WriteFile(executable, []byte("#!/bin/sh\nexit 0\n"), 0700); err != nil {
		t.Fatal(err)
	}
	_, err := Build(context.Background(), Request{Project: project, Executable: executable})
	if err == nil || !strings.Contains(err.Error(), "no HAP") {
		t.Fatalf("expected missing artifact error, got %v", err)
	}
}

func TestBuildFailureDoesNotExposeBuildOutput(t *testing.T) {
	project := t.TempDir()
	os.WriteFile(filepath.Join(project, "build-profile.json5"), []byte("{}"), 0600)
	executable := filepath.Join(project, "hvigorw")
	os.WriteFile(executable, []byte("#!/bin/sh\necho SECRET_SIGNING_CONFIGURATION >&2\nexit 1\n"), 0700)
	_, err := Build(context.Background(), Request{Project: project, Executable: executable})
	if err == nil || strings.Contains(err.Error(), "SECRET") {
		t.Fatalf("expected sanitized build failure: %v", err)
	}
}
