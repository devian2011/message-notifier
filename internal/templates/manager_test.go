package templates

import (
	"errors"
	"notifier/internal/entity"
	"testing"

	"github.com/stretchr/testify/assert"
)

type testStore struct {
	cfg map[string]entity.MessageTemplate
}

func (t *testStore) Get(code string) (entity.MessageTemplate, error) {
	cfg, exists := t.cfg[code]
	if !exists {
		return entity.MessageTemplate{}, errors.New("not found")
	}
	return cfg, nil
}

func TestTemplateManager_BuildMessage(t *testing.T) {
	store := &testStore{
		cfg: map[string]entity.MessageTemplate{
			"test": {
				From:    "from@from",
				Subject: "subject {{.sub}}",
				Body:    "Body: {{.body}}",
			},
		},
	}

	m := NewTemplateManager(store)

	_, err := m.BuildMessage("no", map[string]interface{}{})

	assert.NotNil(t, err)

	msg, noErr := m.BuildMessage("test", map[string]interface{}{"sub": "Subject", "body": "Hello World"})

	assert.Nil(t, noErr)

	assert.Equal(t, store.cfg["test"].From, msg.From)
	assert.Equal(t, store.cfg["test"].To, msg.To)
	assert.Equal(t, "subject Subject", msg.Subject)
	assert.Equal(t, "Body: Hello World", msg.Body)
}
