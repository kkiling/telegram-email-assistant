package app

import (
	"fmt"

	"github.com/kiling91/telegram-email-assistant/internal/factory"
	"github.com/kiling91/telegram-email-assistant/internal/factory/factory_impl"
	"github.com/kiling91/telegram-email-assistant/pkg/bot"
	"github.com/sirupsen/logrus"
)

type App struct {
	fact factory.Factory
}

func NewApp(configFile string) *App {
	return &App{
		fact: factory_impl.NewFactory(configFile),
	}
}

func (a *App) allowedUser(ctx bot.Context) bool {
	allowedUserId := a.fact.Config().Telegram.AllowedUserId
	allowed := false
	userId := ctx.UserId()

	for _, id := range allowedUserId {
		if id == userId {
			allowed = true
			break
		}
	}

	if !allowed {
		_, err := a.fact.Bot().Send(userId, fmt.Sprintf("❗ Access is denied: your id #%d", userId))
		logrus.Warn(err)
		return false
	}

	return true
}

func (a *App) startCommand() {
	login := a.fact.Config().Imap.Login
	b := a.fact.Bot()
	b.Handle("/start", func(ctx bot.Context) error {
		if !a.allowedUser(ctx) {
			return nil
		}

		msg := "✌ Hey! I am your personal email assistant.\n"
		msg += fmt.Sprintf("📧 I will send notifications of new email in your mailbox: %s", login)

		_, err := b.Send(ctx.UserId(), msg)
		return err
	})
}

func (a *App) Start() {
	a.startCommand()
	a.fact.Bot().Start()
}

func (a *App) Shutdown() {
	a.fact.Bot().Stop()
}
