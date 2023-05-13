package config

import (
	"io"
	"os"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

type Store struct {
	URI string `yaml:"uri"`
}

type Runtime struct {
	LogLevel    string `yaml:"log_level"`
	ExternalURL string `yaml:"external_url"`
	SendMetrics bool   `yaml:"send_metrics"`
}

type App struct {
	ImageTimelineDays int `yaml:"image_timeline_days"`
	ImageVersionLimit int `yaml:"image_version_limit"`
}

// Config represents app config object.
type Config struct {
	Name      string  `yaml:"name"`
	Version   string  `yaml:"version"`
	ProjectID string  `yaml:"project_id"`
	Runtime   Runtime `yaml:"runtime"`
	Store     Store   `yaml:"store"`
	App       App     `yaml:"app"`
}

// ReadFromEnvVarFile reads app config from file path in env var CONFIG.
func ReadFromEnvVarFile() (*Config, error) {
	return ReadFromFile(GetEnv("CONFIG", ""))
}

// ReadFromFile reads app config from file.
func ReadFromFile(path string) (*Config, error) {
	if path == "" {
		return nil, errors.New("config value must be existing file path")
	}

	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		return nil, errors.Wrapf(err, "provided config file does not exists: %s", path)
	}

	j, err := os.Open(path)
	if err != nil {
		return nil, errors.Wrapf(err, "error opening config file: %s", path)
	}
	defer j.Close()

	b, err := io.ReadAll(j)
	if err != nil {
		return nil, errors.Wrapf(err, "error reading config file %v", j)
	}

	var c Config
	if err := yaml.Unmarshal(b, &c); err != nil {
		return nil, errors.Wrapf(err, "error unmarshalling config file %v", j)
	}
	return &c, nil
}
