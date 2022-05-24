package factory

import (
	"github.com/kiling91/telegram-email-assistant/internal/bot"
	"github.com/kiling91/telegram-email-assistant/internal/config"
	"github.com/kiling91/telegram-email-assistant/internal/email"
	printmsg "github.com/kiling91/telegram-email-assistant/internal/print_msg"
)

type Factory interface {
	Config() *config.Config
	Bot() bot.Bot
	ImapEmail() email.ReadEmail
	PrintMsg() printmsg.PrintMsg
}
