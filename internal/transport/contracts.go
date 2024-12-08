package transport

import (
	"github.com/google/uuid"

	"notifier/internal/entity"
)

type Handler interface {
	HandleTemplateMessage(command entity.TemplateCommand) (uuid.UUID, error)
	HandleCustomMessage(command entity.CustomCommand) (uuid.UUID, error)
}
