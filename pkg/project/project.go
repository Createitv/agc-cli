package project

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"
)

type Config struct {
	AppID       string    `json:"appId"`
	ProjectID   string    `json:"projectId,omitempty"`
	PackageName string    `json:"packageName,omitempty"`
	Profile     string    `json:"profile,omitempty"`
	CreatedAt   time.Time `json:"createdAt"`
}

func Path(cwd string) string {
	return filepath.Join(cwd, ".agc", "project.json")
}

func Save(cwd string, config Config) error {
	if config.CreatedAt.IsZero() {
		config.CreatedAt = time.Now().UTC()
	}
	path := Path(cwd)
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

func Load(cwd string) (Config, error) {
	data, err := os.ReadFile(Path(cwd))
	if err != nil {
		return Config{}, err
	}
	var config Config
	return config, json.Unmarshal(data, &config)
}
