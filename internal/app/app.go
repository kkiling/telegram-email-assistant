package app

import (
	"context"
	"fmt"
	"time"

	"github.com/kiling91/telegram-email-assistant/internal/email"
	"github.com/kiling91/telegram-email-assistant/internal/factory"
	"github.com/kiling91/telegram-email-assistant/internal/factory/factory_impl"
	"github.com/kiling91/telegram-email-assistant/pkg/bot"
	"github.com/sirupsen/logrus"
)

type BtnType = string

const (
	BtnMark BtnType = "btn_mark"
	BtnRead BtnType = "btn_read"
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
	allowedUserIds := a.fact.Config().Telegram.AllowedUserIds
	allowed := false
	userId := ctx.UserId()

	for _, id := range allowedUserIds {
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

func (a *App) readEmailsLoop(ctx context.Context) {
	cfg := a.fact.Config()
	userIds := cfg.Telegram.AllowedUserIds
	imapUser := &email.ImapUser{
		ImapServer: cfg.Imap.ImapServer,
		Login:      cfg.Imap.Login,
		Password:   cfg.Imap.Password,
	}
	reader := NewReader(a.fact, userIds, imapUser)
	isFirst := true
	for alive := true; alive; {
		var timer *time.Timer
		if isFirst {
			timer = time.NewTimer(time.Second)
			isFirst = false
		} else {
			timer = time.NewTimer(time.Duration(cfg.App.MailCheckTimeout) * time.Second)
		}
		select {
		case <-ctx.Done():
			alive = false
		case <-timer.C:
			reader.Start(ctx)
		}
	}
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

func (a *App) Start(ctx context.Context) {
	a.startCommand()
	go a.readEmailsLoop(ctx)
	a.fact.Bot().Start()
}

func (a *App) Shutdown() {
	a.fact.Bot().Stop()
	a.fact.Storage().ShutDown()
}
