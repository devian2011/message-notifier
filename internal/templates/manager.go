package templates

import (
	"errors"
	"fmt"

	"notifier/internal/entity"
)

type TemplateStorage interface {
	Get(code string) (entity.MessageTemplate, error)
}

type TemplateManager struct {
	store TemplateStorage
}

func NewTemplateManager(store TemplateStorage) *TemplateManager {
	return &TemplateManager{store: store}
}

func (t *TemplateManager) BuildMessage(template string, params map[string]interface{}) (entity.Message, error) {
	tmpl, cfgErr := t.store.Get(template)
	if cfgErr != nil {
		return entity.Message{}, errors.Join(
			fmt.Errorf("template with code: %s not exists", template),
			cfgErr)
	}

	msg := entity.Message{
		From: tmpl.From,
		To:   tmpl.To,
	}

	var parseSubjectErr error
	msg.Subject, parseSubjectErr = Parse(tmpl.Subject, params)
	if parseSubjectErr != nil {
		return msg, parseSubjectErr
	}

	var parseBodyErr error
	msg.Body, parseBodyErr = Parse(tmpl.Body, params)
	if parseBodyErr != nil {
		return entity.Message{}, parseBodyErr
	}

	return msg, nil
}
