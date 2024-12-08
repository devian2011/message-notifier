package internal

import (
	"os"

	"gopkg.in/yaml.v3"

	"notifier/internal/providers"
	"notifier/internal/transport"
)

type Config struct {
	Version   string `json:"version" yaml:"version"`
	Transport struct {
		Http transport.HttpConfig `json:"http" yaml:"http"`
	} `json:"transport" yaml:"transport"`
	Providers map[string]providers.Config `json:"providers" yaml:"providers"`
}

func loadConfig(filePath string) (*Config, error) {
	file, fileOpenErr := os.Open(filePath)
	if fileOpenErr != nil {
		return nil, fileOpenErr
	}

	cfg := &Config{}

	decodeErr := yaml.NewDecoder(file).Decode(cfg)

	return cfg, decodeErr
}
