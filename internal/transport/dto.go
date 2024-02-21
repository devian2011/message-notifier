package transport

import (
	"time"
)

type TemplateCommand struct {
	Provider string            `json:"provider"`
	Template string            `json:"template"`
	Time     Delay             `json:"time"`
	To       []string          `json:"to"`
	Params   map[string]string `json:"params"`
}

func (c *TemplateCommand) GetProvider() string {
	return c.Provider
}

func (c *TemplateCommand) GetTemplate() string {
	return c.Template
}

func (c *TemplateCommand) GetTime() Delay {
	return c.Time
}

func (c *TemplateCommand) GetTo() []string {
	return c.To
}

func (c *TemplateCommand) GetParams() map[string]string {
	return c.Params
}

type CustomCommand struct {
	Provider string
	Time     Delay   `json:"time"`
	Message  Message `json:"message"`
}

func (c *CustomCommand) GetProvider() string {
	return c.Provider
}

func (c *CustomCommand) GetTime() Delay {
	return c.Time
}

func (c *CustomCommand) GetMessage() Message {
	return c.Message
}

type Delay struct {
	Delay time.Duration `json:"delay"`
	Plan  time.Time     `json:"plan"`
}

func (d *Delay) GetDelay() time.Duration {
	return d.Delay
}

func (d *Delay) GetPlan() time.Time {
	return d.Plan
}

type Message struct {
	From    string   `json:"from"`
	To      []string `json:"to"`
	Subject string   `json:"subject"`
	Body    string   `json:"body"`
}

func (m *Message) GetFrom() string {
	return m.From
}

func (m *Message) GetTo() []string {
	return m.To
}

func (m *Message) GetSubject() string {
	return m.Subject
}

func (m *Message) GetBody() string {
	return m.Body
}
