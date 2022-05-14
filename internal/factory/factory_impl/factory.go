package factory_impl

import (
	"github.com/kiling91/telegram-email-assistant/internal/config"
	"github.com/kiling91/telegram-email-assistant/internal/email"
	"github.com/kiling91/telegram-email-assistant/internal/email/imapmsg"
	"github.com/kiling91/telegram-email-assistant/internal/factory"
)

type fact struct {
	config    *config.Config
	imapEmail email.ReadEmail
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
