package handler

import "time"

type Delay interface {
	GetDelay() time.Duration
	GetPlan() time.Time
}

type TemplateMessage interface {
	GetProvider() string
	GetTemplate() string
	GetTime() Delay
	GetTo() []string
	GetParams() map[string]string
}

type CustomMessage interface {
	GetProvider() string
	GetTime() Delay
	GetMessage() Message
}

type Message interface {
	GetFrom() string
	GetTo() []string
	GetSubject() string
	GetBody() string
}

type Provider interface {
	Send(message Message) error
}

type ProviderRegistry interface {
	Get(code string) (Provider, error)
}

type TemplateEngine interface {
	GetTemplate(code string, params map[string]interface{}) (Message, error)
}

type Handler struct {
	providerRegistry ProviderRegistry
	inputTmplCh      chan TemplateMessage
	inputCustomCh    chan CustomMessage
}

func (h *Handler) AddProviderRegistry(registry ProviderRegistry) {
	h.providerRegistry = registry
}

func (h *Handler) AddInputs(tmplCh chan TemplateMessage, customCh chan CustomMessage) {
	h.inputTmplCh = tmplCh
	h.inputCustomCh = customCh
}
