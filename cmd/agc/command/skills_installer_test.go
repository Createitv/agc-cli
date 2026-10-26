package command

import (
	"context"
	"io"
	"strings"
	"testing"
)

func TestSkillsListDelegatesToSkillsInstaller(t *testing.T) {
	original := runSkillsCLI
	t.Cleanup(func() { runSkillsCLI = original })
	var got []string
	runSkillsCLI = func(_ context.Context, stdout, _ io.Writer, args ...string) error {
		got = append([]string(nil), args...)
		_, _ = io.WriteString(stdout, "agc-cli-usage\n")
		return nil
	}
	out, err := execute("skills", "list", "--remote")
	if err != nil {
		t.Fatal(err)
	}
	if strings.Join(got, " ") != "--yes skills add Createitv/agc-cli --list" {
		t.Fatalf("installer args = %q", got)
	}
	if out != "agc-cli-usage\n" {
		t.Fatalf("output = %q", out)
	}
}

func TestSkillsInstallSelectsSkillAgentAndScope(t *testing.T) {
	original := runSkillsCLI
	t.Cleanup(func() { runSkillsCLI = original })
	var got []string
	runSkillsCLI = func(_ context.Context, _, _ io.Writer, args ...string) error {
		got = append([]string(nil), args...)
		return nil
	}
	_, err := execute("skills", "install", "--name", "agc-auth", "--agent", "codex", "--global")
	if err != nil {
		t.Fatal(err)
	}
	want := "--yes skills add Createitv/agc-cli --skill agc-auth --agent codex --global"
	if strings.Join(got, " ") != want {
		t.Fatalf("installer args = %q, want %q", got, want)
	}
}
