package templates

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTemplateManager_BuildMessage(t *testing.T) {
	cfg := map[string]TemplateCfg{
		"test": {
			From:    "from@from",
			Subject: "subject {{.sub}}",
			Body: struct {
				Text string `json:"text" yaml:"text"`
				Path string `json:"path" yaml:"path"`
			}{Text: "Body: {{.body}}", Path: ""},
		},
	}

	m := NewTemplateManager(cfg)

	_, err := m.BuildMessage("no", map[string]interface{}{})

	assert.True(t, errors.Is(err, ErrTemplateNotExists))

	msg, noErr := m.BuildMessage("test", map[string]interface{}{"sub": "Subject", "body": "Hello World"})

	assert.Nil(t, noErr)

	assert.Equal(t, cfg["test"].From, msg.From)
	assert.Equal(t, cfg["test"].To, msg.To)
	assert.Equal(t, "subject Subject", msg.Subject)
	assert.Equal(t, "Body: Hello World", msg.Body)
}
