package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type CommitType struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type TicketKeyWord struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type Config struct {
	Types            []CommitType    `json:"types"`
	Scopes           []string        `json:"scopes"`
	RequireScope     bool            `json:"requireScope"`
	SubjectCharLimit int             `json:"subjectCharLimit"`
	TicketKeyWords   []TicketKeyWord `json:"ticketKeyWords"`
	Theme            string          `json:"theme"`
}

var defaultConfig = Config{
	Types: []CommitType{
		{"feat", "A new feature"},
		{"fix", "A bug fix"},
		{"chore", "Build process or auxiliary tool changes"},
		{"docs", "Documentation only changes"},
		{"style", "Markup, white-space, formatting, missing semi-colons..."},
		{"refactor", "A code change that neither fixes a bug nor adds a feature"},
		{"perf", "A code change that improves performance"},
		{"test", "Adding missing tests"},
		{"build", "Changes that affect the build system or external dependencies"},
		{"ci", "CI related changes"},
	},
	Scopes:           []string{},
	RequireScope:     false,
	SubjectCharLimit: 100,
	TicketKeyWords: []TicketKeyWord{
		{"Closes", "Closes the issue when merged"},
		{"Fixes", "Fixes a bug and closes the issue"},
		{"Refs", "References without closing"},
	},
	Theme: "base16",
}

func globalConfigPath() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".config", "snap", "config.json")
}

func localConfigPath() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	for {
		gitDir := filepath.Join(dir, ".git")
		if _, err := os.Stat(gitDir); err == nil {
			return filepath.Join(dir, ".snap.json"), nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return "", fmt.Errorf("not a git repository")
		}
		dir = parent
	}
}

func loadConfig() Config {
	if localPath, err := localConfigPath(); err == nil {
		if data, err := os.ReadFile(localPath); err == nil {
			var cfg Config
			if err := json.Unmarshal(data, &cfg); err == nil {
				return cfg
			}
		}
	}
	if data, err := os.ReadFile(globalConfigPath()); err == nil {
		var cfg Config
		if err := json.Unmarshal(data, &cfg); err == nil {
			return cfg
		}
	}
	return defaultConfig
}

func saveConfig(cfg Config, path string) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}

	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0o644)
}
