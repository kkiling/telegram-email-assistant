package factory

import (
	"github.com/kiling91/telegram-email-assistant/internal/config"
	"github.com/kiling91/telegram-email-assistant/internal/email"
	printmsg "github.com/kiling91/telegram-email-assistant/internal/print_msg"
)

type Factory interface {
	Config() *config.Config
	ImapEmail() email.ReadEmail
	PrintMsg() printmsg.PrintMsg
}
