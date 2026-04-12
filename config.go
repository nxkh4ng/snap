package main

import (
	"encoding/json"
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

func loadConfig() Config {
	data, err := os.ReadFile(globalConfigPath())
	if err != nil {
		return defaultConfig
	}

	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return defaultConfig
	}

	return cfg
}

func saveConfig(cfg Config) error {
	path := globalConfigPath()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}

	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0o644)
}
