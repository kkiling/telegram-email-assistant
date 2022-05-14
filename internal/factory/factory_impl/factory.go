package factory_impl

import (
	"github.com/kiling91/telegram-email-assistant/internal/config"
	"github.com/kiling91/telegram-email-assistant/internal/email"
	"github.com/kiling91/telegram-email-assistant/internal/email/imap_msg"
	"github.com/kiling91/telegram-email-assistant/internal/factory"
	printmsg "github.com/kiling91/telegram-email-assistant/internal/print_msg"
	telegrammsg "github.com/kiling91/telegram-email-assistant/internal/print_msg/telegram_msg"
)

type fact struct {
	config    *config.Config
	imapEmail email.ReadEmail
	printMsg  printmsg.PrintMsg
}

func NewFactory() factory.Factory {
	return &fact{}
}

func (f *fact) Config() *config.Config {
	if f.config == nil {
		f.config = config.NewConfig()
	}
	return f.config
}

func (f *fact) ImapEmail() email.ReadEmail {
	if f.imapEmail == nil {
		f.imapEmail = imapmsg.NewReadEmail(f)
	}
	return f.imapEmail
}

func (f *fact) PrintMsg() printmsg.PrintMsg {
	if f.printMsg == nil {
		f.printMsg = telegrammsg.NewPrintEmail(f)
	}
	return f.printMsg
}