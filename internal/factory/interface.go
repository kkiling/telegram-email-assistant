package factory

import (
	"github.com/kiling91/telegram-email-assistant/internal/config"
	"github.com/kiling91/telegram-email-assistant/internal/email"
	"github.com/kiling91/telegram-email-assistant/internal/printmsg"
	"github.com/kiling91/telegram-email-assistant/internal/storage"
	"github.com/kiling91/telegram-email-assistant/pkg/bot"
)

type Factory interface {
	Config() *config.Config
	Bot() bot.Bot
	ImapEmail() email.ReadEmail
	PrintMsg() printmsg.PrintMsg
	Storage() storage.Storage
}
