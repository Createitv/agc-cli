package project

import "testing"

func TestSaveAndLoadProjectConfig(t *testing.T) {
	dir := t.TempDir()
	want := Config{AppID: "123", ProjectID: "p1", PackageName: "com.example.app", Profile: "prod"}
	if err := Save(dir, want); err != nil {
		t.Fatal(err)
	}
	got, err := Load(dir)
	if err != nil {
		t.Fatal(err)
	}
	if got.AppID != want.AppID || got.ProjectID != want.ProjectID || got.PackageName != want.PackageName || got.Profile != want.Profile {
		t.Fatalf("config = %#v, want %#v", got, want)
	}
}
