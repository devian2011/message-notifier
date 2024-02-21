package internal

import (
	"context"
	"time"
)

type App struct {
	ctx       context.Context
	cfg       *Config
	timeStart time.Time
}

func NewApp(ctx context.Context, filePath string) (*App, error) {
	cfg, cfgErr := loadConfig(filePath)
	if cfgErr != nil {
		return nil, cfgErr
	}

	return &App{
		ctx:       ctx,
		cfg:       cfg,
		timeStart: time.Now(),
	}, nil
}

func (a *App) Run() error {
	errCh := make(chan error)

	select {
	case err := <-errCh:
		return err
	case <-a.ctx.Done():
		return nil
	}
}
