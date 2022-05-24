package factory_impl

import (
	"log"

	"github.com/kiling91/telegram-email-assistant/internal/bot"
	"github.com/kiling91/telegram-email-assistant/internal/bot/tgbot"
	"github.com/kiling91/telegram-email-assistant/internal/config"
	"github.com/kiling91/telegram-email-assistant/internal/email"
	imapmsg "github.com/kiling91/telegram-email-assistant/internal/email/imap_msg"
	"github.com/kiling91/telegram-email-assistant/internal/factory"
	printmsg "github.com/kiling91/telegram-email-assistant/internal/print_msg"
	telegrammsg "github.com/kiling91/telegram-email-assistant/internal/print_msg/telegram_msg"
)

type fact struct {
	config    *config.Config
	imapEmail email.ReadEmail
	printMsg  printmsg.PrintMsg
	bot       bot.Bot
}

func NewFactory() factory.Factory {
	return &fact{}
}

func (f *fact) Config() *config.Config {
	if f.config == nil {
		cfg, err := config.NewConfig()
		if err != nil {
			log.Fatal(err)
		}
		f.config = cfg
	}
	return f.config
}

func (f *fact) Bot() bot.Bot {
	if f.bot == nil {
		cfg := f.Config()
		bot, err := tgbot.NewTbBot(&cfg.Telegram)
		if err != nil {
			log.Fatalf("error init tgbot: %v", bot)
		}
		f.bot = bot
	}
	return f.bot
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
