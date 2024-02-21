package providers

import (
	"strconv"

	"notifier/internal/handler"
	"notifier/pkg/parse"
)

const (
	SmtpProviderCode = "smtp"
)

type smtpProviderConfig struct {
	host     string
	port     int
	username string
	password string
}

func newSmtpProviderConfig(fields map[string]string) (smtpProviderConfig, error) {
	cfg := smtpProviderConfig{}
	if h, exists := fields["host"]; exists {
		cfg.host = parse.GetValue(h)
	} else {
		return cfg, ErrEmptyProviderConfigRequiredValue
	}
	if p, exists := fields["port"]; exists {
		cfg.port, _ = strconv.Atoi(p)
	} else {
		return cfg, ErrEmptyProviderConfigRequiredValue
	}

	if u, exists := fields["username"]; exists {
		cfg.username = parse.GetValue(u)
	} else {
		return cfg, ErrEmptyProviderConfigRequiredValue
	}
	if pwd, exists := fields["password"]; exists {
		cfg.password = parse.GetValue(pwd)
	} else {
		return cfg, ErrEmptyProviderConfigRequiredValue
	}

	return cfg, nil
}

type SmtpProvider struct {
	cfg smtpProviderConfig
}

func NewSmtpProvider(cfg *Config) (*SmtpProvider, error) {
	providerCfg, cfgErr := newSmtpProviderConfig(cfg.Params)
	if cfgErr != nil {
		return nil, cfgErr
	}

	return &SmtpProvider{cfg: providerCfg}, nil
}

func (s *SmtpProvider) Send(message handler.Message) error {
	return nil
}