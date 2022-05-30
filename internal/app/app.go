package app

import (
	"context"
	"fmt"
	"sort"

	"github.com/kiling91/telegram-email-assistant/internal/email"
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

func (a *App) readEmails(userIds []int64, imapUser *email.ImapUser) {
	logrus.Info("Start read unseen emails")
	imap := a.fact.ImapEmail()
	b := a.fact.Bot()
	pnt := a.fact.PrintMsg()
	storage := a.fact.Storage()

	emails, err := imap.ReadUnseenEmails(context.Background(), imapUser)
	if err != nil {
		logrus.Fatalln(err)
	}

	sort.Slice(emails, func(i, j int) bool {
		return emails[i].Date.Before(emails[j].Date)
	})

	for _, e := range emails {
		if contains, err := storage.MsgIdContains(imapUser.Login, e.Uid); err != nil {
			logrus.Warnf("error get msg contains from storage: %v", err)
		} else if contains {
			continue
		}

		msg := pnt.PrintMsgEnvelope(e)
		for _, id := range userIds {
			if _, err := b.Send(id, msg); err != nil {
				logrus.Warnf("error send msg: %v", err)
			} else {
				if err := storage.SaveMsgId(imapUser.Login, e.Uid); err != nil {
					logrus.Warnf("error save msg id to storage: %v", err)
				}
			}
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

func (a *App) Start() {
	a.startCommand()

	// ***
	cfg := a.fact.Config()
	userIds := cfg.Telegram.AllowedUserIds
	imapUser := &email.ImapUser{
		ImapServer: cfg.Imap.ImapServer,
		Login:      cfg.Imap.Login,
		Password:   cfg.Imap.Password,
	}
	// ***

	go a.readEmails(userIds, imapUser)
	a.fact.Bot().Start()
}

func (a *App) Shutdown() {
	a.fact.Bot().Stop()
	a.fact.Storage().ShutDown()
}
