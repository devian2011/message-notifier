package templates

import (
	"errors"
	"fmt"

	"notifier/internal/entity"
)

var (
	ErrTemplateNotExists     = errors.New("template not exists")
	ErrTemplateFileNotExists = errors.New("template file not exists")
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

type TemplateManager struct {
	cfg map[string]TemplateCfg
}

func NewTemplateManager(cfg map[string]TemplateCfg) *TemplateManager {
	return &TemplateManager{cfg: cfg}
}

func (t *TemplateManager) BuildMessage(template string, params map[string]interface{}) (entity.Message, error) {
	cfg, exists := t.cfg[template]
	if !exists {
		return entity.Message{}, errors.Join(
			fmt.Errorf("template with code: %s not exists", template), ErrTemplateNotExists)
	}

	msg := entity.Message{
		From: cfg.From,
		To:   cfg.To,
	}

	var parseSubjectErr error
	msg.Subject, parseSubjectErr = Parse(cfg.Subject, params)
	if parseSubjectErr != nil {
		return msg, parseSubjectErr
	}

	var parseBodyErr error

	if cfg.Body.Text != "" {
		msg.Body, parseBodyErr = Parse(cfg.Body.Text, params)
		if parseBodyErr != nil {
			return entity.Message{}, parseBodyErr
		}
	} else {
		if cfg.Body.Path != "" {
			msg.Body, parseBodyErr = ParseFile(cfg.Body.Text, params)
			if parseBodyErr != nil {
				return entity.Message{}, parseBodyErr
			}
		}
	}

	return msg, nil
}
