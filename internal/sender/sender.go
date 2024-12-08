package sender

import (
	"context"
	"sync"

	"notifier/internal/entity"
)

type Storage interface {
	Save(task entity.MessageTask) error
	GetTask() <-chan entity.MessageTask
}

type Provider interface {
	Send(message entity.Message) error
}

type ProviderRegistry interface {
	Get(code string) (Provider, error)
}

type Config struct {
	ThreadCnt uint8 `json:"thread_cnt"`
}

type Sender struct {
	ctx              context.Context
	cfg              *Config
	wg               *sync.WaitGroup
	store            Storage
	providerRegistry ProviderRegistry
	errCh            chan<- error
}

func NewSender(
	ctx context.Context,
	cfg *Config,
	store Storage,
	providerRegistry ProviderRegistry,
	errCh chan<- error,
) *Sender {
	return &Sender{
		ctx:              ctx,
		cfg:              cfg,
		wg:               &sync.WaitGroup{},
		store:            store,
		providerRegistry: providerRegistry,
		errCh:            errCh,
	}
}

func (s *Sender) Run() {
	s.wg.Add(int(s.cfg.ThreadCnt))
	for c := uint8(0); c < s.cfg.ThreadCnt; c++ {
		go s.send()
	}
}

func (s *Sender) Wait() {
	s.wg.Wait()
}

func (s *Sender) send() {
	for {
		select {
		case <-s.ctx.Done():
			return

		}
	}
}
