package providers

import "errors"

var (
	ErrEmptyProviderConfigRequiredValue = errors.New("empty provider config required value")
	ErrUnknownProviderWithCode          = errors.New("unknown provider with code")
)

type Config struct {
	Code     string            `json:"code" yaml:"code"`
	Provider string            `json:"provider" yaml:"provider"`
	Params   map[string]string `json:"params" yaml:"params"`
}
