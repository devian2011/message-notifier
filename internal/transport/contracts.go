package transport

import (
	"github.com/google/uuid"

	"notifier/internal/entity"
)

type Handler interface {
	HandleTemplateMessage(command entity.TemplateCommand) (uuid.UUID, error)
	HandleCustomMessage(command entity.CustomCommand) (uuid.UUID, error)
}

type TemplateStore interface {
	SaveTemplate(template entity.MessageTemplate) error
	GetTemplates() ([]entity.MessageTemplate, error)
	GetTemplate(code string) (entity.MessageTemplate, error)
}

type ProvidersStore interface {
	GetProviders() map[string]string
}

type TaskStore interface {
	Get(id uuid.UUID) (entity.MessageTask, error)
	Search(params entity.SearchParams) ([]entity.MessageTask, error)
}
