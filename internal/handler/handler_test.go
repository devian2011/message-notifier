package handler

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	"notifier/internal/entity"
)

type testStore struct {
	task entity.MessageTask
}

func (t *testStore) Save(task entity.MessageTask) error {
	t.task = task
	return nil
}

type testTManager struct{}

var tMessage = entity.Message{
	From:    "test@test.test",
	To:      []string{"to_test"},
	Subject: "Test Subject",
	Body:    "Test Body",
}

func (t *testTManager) BuildMessage(template string, params map[string]interface{}) (entity.Message, error) {
	if template == "err" {
		return entity.Message{}, errors.New("error message")
	} else {
		return tMessage, nil
	}
}

func TestHandler_HandleCustomMessage(t *testing.T) {
	store := &testStore{}
	h := NewHandler(&testTManager{}, store)

	cmd := entity.CustomCommand{
		Provider: "tg",
		Code:     "code",
		Retries:  5,
		Time:     entity.Delay{},
		Message:  tMessage,
	}

	id, handleErr := h.HandleCustomMessage(cmd)

	assert.Nil(t, handleErr)
	assert.Equal(t, id, store.task.Id)
	assert.Equal(t, cmd.Provider, store.task.Provider)
	assert.Equal(t, cmd.Message, store.task.Message)
	assert.Equal(t, entity.TaskStatusToProcess, store.task.Status)
	assert.Equal(t, cmd.Retries, store.task.MaxRetryCnt)
	assert.Equal(t, uint8(0), store.task.RetryCnt)
	assert.True(t, store.task.LastExecutionTime.IsZero())

}

func TestHandler_HandleTemplateMessage(t *testing.T) {
	store := &testStore{}
	h := NewHandler(&testTManager{}, store)

	errCommand := entity.TemplateCommand{
		Provider: "mail",
		Template: "err",
		Retries:  1,
		Time:     entity.Delay{},
		To:       []string{"to"},
		Params:   map[string]interface{}{"1": 1, "m": "m"},
	}

	_, errCommandErr := h.HandleTemplateMessage(errCommand)

	assert.Equal(t, "error message", errCommandErr.Error())

	cmd := entity.TemplateCommand{
		Provider: "tg",
		Template: "command",
		Retries:  3,
		Time:     entity.Delay{},
		To:       []string{"to"},
		Params:   map[string]interface{}{},
	}

	id, errCommandErr := h.HandleTemplateMessage(cmd)

	assert.Nil(t, errCommandErr)
	assert.Equal(t, id, store.task.Id)
	assert.Equal(t, cmd.Provider, store.task.Provider)

	assert.Equal(t, cmd.To, store.task.Message.To)
	assert.Equal(t, tMessage.Body, store.task.Message.Body)
	assert.Equal(t, tMessage.From, store.task.Message.From)
	assert.Equal(t, tMessage.Subject, store.task.Message.Subject)

	assert.Equal(t, entity.TaskStatusToProcess, store.task.Status)
	assert.Equal(t, cmd.Retries, store.task.MaxRetryCnt)
	assert.Equal(t, uint8(0), store.task.RetryCnt)
	assert.True(t, store.task.LastExecutionTime.IsZero())
}
