package config

import (
	"fmt"
	"os"
	"gopkg.in/yaml.v3"
)

type Limits struct {
	WallTimeS    int `yaml:"wall_time_s"`
	MemoryKB     int `yaml:"memory_kb"`
	MaxProcesses int `yaml:"max_processes"`
}

type Command struct {
	Cmd           string   `yaml:"cmd"`
	Args          []string `yaml:"args"`
	Limits        Limits   `yaml:"limits"`
	FlagAllowlist []string `yaml:"flag_allowlist,omitempty"`
}

type Language struct {
	ID                       string  `yaml:"id"`
	Name                     string  `yaml:"name"`
	SourceFilename           string  `yaml:"source_filename"`
	SourceFilenameStrategy   string  `yaml:"source_filename_strategy,omitempty"`
	Artifact                 string  `yaml:"artifact,omitempty"`
	ArtifactFilenameStrategy string  `yaml:"artifact_filename_strategy,omitempty"`
	Build                    *Command `yaml:"build,omitempty"`
	Run                      Command  `yaml:"run"`
}

type Config struct {
	Languages []Language `yaml:"languages"`
}

// Load reads the YAML and returns a map keyed by language ID for O(1) lookups.
func Load(path string) (map[string]Language, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse yaml: %w", err)
	}

	registry := make(map[string]Language)
	for _, lang := range cfg.Languages {
		registry[lang.ID] = lang
	}
	return registry, nil
}