package main

import (
	"context"
	"sort"
	"time"

	"github.com/kiling91/telegram-email-assistant/internal/email"
	"github.com/kiling91/telegram-email-assistant/internal/factory/factory_impl"
	log "github.com/sirupsen/logrus"
)

func main() {
	fact := factory_impl.NewFactory()

	cfg := fact.Config()
	user := &email.ImapUser{
		ImapServer: cfg.Imap.ImapServer,
		Login:      cfg.Imap.Login,
		Password:   cfg.Imap.Password,
	}

	imap := fact.ImapEmail()
	emails, err := imap.ReadUnseenEmails(context.Background(), user)
	if err != nil {
		log.Fatalln(err)

	}

	sort.Slice(emails, func(i, j int) bool {
		return emails[i].Date.Before(emails[j].Date)
	})

	for _, email := range emails {
		start := time.Now()

		msg, err := imap.ReadEmail(context.Background(), user, email.Uid)
		if err != nil {
			log.Fatalf("Error read #%d", email.Uid)
		}

		pnt := fact.PrintMsg()
		_, err = pnt.PrintMsgWithBody(msg, user.Login)
		if err != nil {
			log.Fatalf("Error read #%d", email.Uid)
		}

		elapsed := time.Since(start)
		log.Printf("#%d - %s %s (%fs)", email.Uid, email.FromAddress, email.Subject, elapsed.Seconds())
	}
}
