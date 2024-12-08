package handler

import (
	"time"

	"github.com/google/uuid"

	"notifier/internal/entity"
)

type Storage interface {
	Save(task entity.MessageTask) error
}

type TemplateEngine interface {
	BuildMessage(template string, params map[string]interface{}) (entity.Message, error)
}

type Handler struct {
	templateEngine TemplateEngine
	store          Storage
}

func NewHandler(
	templateEngine TemplateEngine,
	store Storage,
) *Handler {
	return &Handler{
		templateEngine: templateEngine,
		store:          store,
	}
}

func (h *Handler) HandleTemplateMessage(command entity.TemplateCommand) (uuid.UUID, error) {
	id, err := uuid.NewV7()
	if err != nil {
		return uuid.UUID{}, err
	}

	msg, msgBuildErr := h.templateEngine.BuildMessage(command.Template, command.Params)

	if msgBuildErr != nil {
		return id, msgBuildErr
	}

	msg.To = append(msg.To, command.To...)

	return id, h.store.Save(entity.MessageTask{
		Id:                 id,
		Provider:           command.Provider,
		Message:            msg,
		Status:             entity.TaskStatusToProcess,
		MaxRetryCnt:        command.Retries,
		StartExecutionTime: getStartExecutionTime(command.Time),
		CreatedAt:          time.Now(),
		RetryCnt:           0,
	})
}

func (h *Handler) HandleCustomMessage(command entity.CustomCommand) (uuid.UUID, error) {
	id, err := uuid.NewV7()
	if err != nil {
		return uuid.UUID{}, err
	}

	return id, h.store.Save(entity.MessageTask{
		Id:                 id,
		Provider:           command.Provider,
		Message:            command.Message,
		Status:             entity.TaskStatusToProcess,
		MaxRetryCnt:        command.Retries,
		StartExecutionTime: getStartExecutionTime(command.Time),
		CreatedAt:          time.Now(),
		RetryCnt:           0,
	})
}

func getStartExecutionTime(d entity.Delay) time.Time {
	runOn := time.Now()
	if d.Delay.Nanoseconds() != 0 {
		runOn = runOn.Add(d.Delay)
	}
	if !d.Plan.IsZero() && d.Plan.After(time.Now()) {
		runOn = d.Plan
	}

	return runOn
}
