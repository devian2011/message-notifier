package internal

import (
	"notifier/internal/providers"
	"notifier/internal/templates"
	"notifier/internal/transport"
)

type Config struct {
	Version   string `json:"version" yaml:"version"`
	Transport struct {
		Http transport.HttpConfig `json:"http" yaml:"http"`
	} `json:"transport" yaml:"transport"`
	Providers []providers.Config      `json:"providers" yaml:"providers"`
	Templates []templates.TemplateCfg `json:"templates" yaml:"templates"`
}

func loadConfig(filePath string) (*Config, error) {
	return &Config{}, nil
}
