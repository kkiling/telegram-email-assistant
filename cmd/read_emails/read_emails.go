package main

import (
	"context"
	"sort"
	"time"

	"github.com/kiling91/telegram-email-assistant/internal/email"
	"github.com/kiling91/telegram-email-assistant/internal/factory/factory_impl"
	"github.com/sirupsen/logrus"
)

func main() {
	fact := factory_impl.NewFactory("config/config.yml")

	cfg := fact.Config()
	user := &email.ImapUser{
		ImapServer: cfg.Imap.ImapServer,
		Login:      cfg.Imap.Login,
		Password:   cfg.Imap.Password,
	}

	imap := fact.ImapEmail()
	emails, err := imap.ReadUnseenEmails(context.Background(), user)
	if err != nil {
		logrus.Fatalln(err)

	}

	sort.Slice(emails, func(i, j int) bool {
		return emails[i].Date.Before(emails[j].Date)
	})

	for _, e := range emails {
		start := time.Now()

		msg, err := imap.ReadEmail(context.Background(), user, e.Uid)
		if err != nil {
			logrus.Fatalf("Error read #%d", e.Uid)
		}

		pnt := fact.PrintMsg()
		_, err = pnt.PrintMsgWithBody(msg, user.Login)
		if err != nil {
			logrus.Fatalf("Error read #%d", e.Uid)
		}

		elapsed := time.Since(start)
		logrus.Printf("#%d - %s %s (%fs)", e.Uid, e.FromAddress, e.Subject, elapsed.Seconds())
	}
}
