package providers

import "errors"

var (
	ErrEmptyProviderConfigRequiredValue = errors.New("empty provider config required value")
	ErrUnknownProvider                  = errors.New("unknown provider with code")
)

type Config struct {
	Provider string            `json:"provider" yaml:"provider"`
	Params   map[string]string `json:"params" yaml:"params"`
}
