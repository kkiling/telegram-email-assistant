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
	fact  factory.Factory
	users map[int64][]string
}

func NewApp(configFile string) *App {
	fact := factory_impl.NewFactory(configFile)

	tUsers := fact.Config().Telegram.Users
	users := make(map[int64][]string)
	for _, u := range tUsers {
		users[u.UserId] = u.ImapLogin
	}

	return &App{
		fact:  fact,
		users: users,
	}
}

func (a *App) allowedUser(ctx bot.Context) bool {
	allowed := false
	userId := ctx.UserId()

	if _, ok := a.users[userId]; ok {
		allowed = true
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

	readers := make([]*Reader, 0)

	for _, u := range cfg.Imap {

		users := make([]int64, 0)
		for _, t := range cfg.Telegram.Users {
			for _, i := range t.ImapLogin {
				if i == u.Login {
					users = append(users, t.UserId)
					break
				}
			}
		}

		imapUser := &email.ImapUser{
			ImapServer: u.ImapServer,
			Login:      u.Login,
			Password:   u.Password,
		}

		reader := NewReader(a.fact, users, imapUser)
		readers = append(readers, reader)
	}

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
			for _, reader := range readers {
				// TODO: Convert to groutins
				reader.Start(ctx)
			}
		}
	}
}

func (a *App) startCommand() {
	b := a.fact.Bot()

	b.Handle("/start", func(ctx bot.Context) error {
		if !a.allowedUser(ctx) {
			return nil
		}

		msg := "✌ Hey! I am your personal email assistant.\n"
		msg += "📧 I will send notifications of new email in your mailbox:\n"
		for _, login := range a.users[ctx.UserId()] {
			msg += fmt.Sprintf("\t - %s\n", login)
		}
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
