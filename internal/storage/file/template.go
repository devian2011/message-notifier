package file

import (
	"errors"
	"os"

	"notifier/internal/entity"
)

var (
	ErrTemplateNotFound  = errors.New("template not found")
	ErrTemplateWrongPath = errors.New("template wrong path")
)

type TemplateCfg struct {
	From    string   `json:"from"`
	To      []string `json:"to" yaml:"to"`
	Subject string   `json:"subject" yaml:"subject"`
	Body    struct {
		Text string `json:"text" yaml:"text"`
		Path string `json:"path" yaml:"path"`
	} `json:"body" yaml:"body"`
}

type TemplateStore struct {
	templates map[string]TemplateCfg
}

func (s *TemplateStore) Get(code string) (entity.MessageTemplate, error) {
	cfg, exists := s.templates[code]
	if !exists {
		return entity.MessageTemplate{}, ErrTemplateNotFound
	}

	tmpl := entity.MessageTemplate{
		From:    cfg.From,
		To:      cfg.To,
		Subject: cfg.Subject,
		Body:    cfg.Body.Text,
	}

	if cfg.Body.Text == "" && cfg.Body.Path != "" {
		data, readErr := os.ReadFile(cfg.Body.Path)
		if readErr != nil {
			return entity.MessageTemplate{}, errors.Join(ErrTemplateWrongPath, readErr)
		}
		tmpl.Body = string(data)
	}

	return tmpl, nil
}
