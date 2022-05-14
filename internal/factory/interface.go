package factory

import (
	"github.com/kiling91/telegram-email-assistant/internal/config"
	"github.com/kiling91/telegram-email-assistant/internal/email"
)

type Factory interface {
	Config() *config.Config
	ImapEmail() email.ReadEmail
}
