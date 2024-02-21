package providers

import "notifier/internal/handler"

type Registry struct {
	providers map[string]handler.Provider
}

func NewRegistry(cfg []Config) (*Registry, error) {
	registry := &Registry{providers: make(map[string]handler.Provider, len(cfg))}
	for _, c := range cfg {
		switch c.Provider {
		case SmtpProviderCode:
			var smtpInitErr error
			registry.providers[c.Code], smtpInitErr = NewSmtpProvider(&c)
			if smtpInitErr != nil {
				return nil, smtpInitErr
			}
		case TgProviderCode:
			var tgInitErr error
			registry.providers[c.Code], tgInitErr = NewTgProvider(&c)
			if tgInitErr != nil {
				return nil, tgInitErr
			}
		}
	}

	return registry, nil
}

func (r *Registry) Get(code string) (handler.Provider, error) {
	if val, exists := r.providers[code]; exists {
		return val, nil
	}

	return nil, ErrUnknownProviderWithCode
}
