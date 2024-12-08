package entity

import (
	"github.com/google/uuid"
	"time"
)

type TaskStatus int

const (
	TaskStatusToProcess TaskStatus = iota
	TaskStatusInProgress
	TaskStatusDone
	TaskStatusError
)

type TemplateCommand struct {
	Provider string                 `json:"provider"`
	Template string                 `json:"template"`
	Retries  uint8                  `json:"retries"`
	Time     Delay                  `json:"time"`
	To       []string               `json:"to"`
	Params   map[string]interface{} `json:"params"`
}

type CustomCommand struct {
	Provider string  `json:"provider"`
	Code     string  `json:"code"`
	Retries  uint8   `json:"retries"`
	Time     Delay   `json:"time"`
	Message  Message `json:"message"`
}

type Delay struct {
	Delay time.Duration `json:"delay"`
	Plan  time.Time     `json:"plan"`
}

type Message struct {
	From    string   `json:"from"`
	To      []string `json:"to"`
	Subject string   `json:"subject"`
	Body    string   `json:"body"`
}

type MessageTemplate struct {
	From    string   `json:"from" yaml:"from"`
	To      []string `json:"to" yaml:"to"`
	Subject string   `json:"subject" yaml:"subject"`
	Body    string   `json:"body" yaml:"body"`
}

type MessageTask struct {
	Id                 uuid.UUID  `json:"id"`
	Provider           string     `json:"provider"`
	Message            Message    `json:"message"`
	Status             TaskStatus `json:"status"`
	MaxRetryCnt        uint8      `json:"max_retry_cnt"`
	StartExecutionTime time.Time  `json:"start_execution_time"`
	LastExecutionTime  time.Time  `json:"last_execution_time"`

	CreatedAt time.Time `json:"created_at"`
	RetryCnt  uint8     `json:"retry_cnt"`
}

type SearchParams struct {
	Provider []string
	Status   []TaskStatus
	OrderBy  []string
}
