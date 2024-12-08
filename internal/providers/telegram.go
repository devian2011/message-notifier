package providers

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"notifier/internal/entity"

	"notifier/pkg/config"
)

const (
	TgProviderCode = "TG"
)

type tgProviderConfig struct {
	token    string
	username string
	chatId   int
}

func newTgProviderConfig(fields map[string]string) (tgProviderConfig, error) {
	cfg := tgProviderConfig{}
	if h, exists := fields["host"]; exists {
		cfg.token = config.GetValue(h)
	} else {
		return cfg, ErrEmptyProviderConfigRequiredValue
	}

	return cfg, nil
}

type TgProvider struct {
	cfg tgProviderConfig
	bot *tgbotapi.BotAPI
}

func NewTgProvider(cfg *Config) (*TgProvider, error) {
	providerCfg, cfgErr := newTgProviderConfig(cfg.Params)
	if cfgErr != nil {
		return nil, cfgErr
	}

	bot, botInitErr := tgbotapi.NewBotAPI(providerCfg.token)
	if botInitErr != nil {
		return nil, botInitErr
	}

	return &TgProvider{cfg: providerCfg, bot: bot}, nil
}

func (tg *TgProvider) Send(message entity.Message) error {
	return nil
}
