package providers

import "notifier/internal/sender"

type Registry struct {
	providers map[string]sender.Provider
}

func NewRegistry(cfg map[string]Config) (*Registry, error) {
	registry := &Registry{providers: make(map[string]sender.Provider, len(cfg))}
	for code, c := range cfg {
		switch c.Provider {
		case SmtpProviderCode:
			var smtpInitErr error
			registry.providers[code], smtpInitErr = NewSmtpProvider(&c)
			if smtpInitErr != nil {
				return nil, smtpInitErr
			}
		case TgProviderCode:
			var tgInitErr error
			registry.providers[code], tgInitErr = NewTgProvider(&c)
			if tgInitErr != nil {
				return nil, tgInitErr
			}
		}
	}

	return registry, nil
}

func (r *Registry) Get(code string) (sender.Provider, error) {
	if val, exists := r.providers[code]; exists {
		return val, nil
	}

	return nil, ErrUnknownProvider
}
