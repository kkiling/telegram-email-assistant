package main

import (
	"context"
	"time"

	"github.com/kiling91/telegram-email-assistant/internal/email"
	"github.com/kiling91/telegram-email-assistant/internal/factory/factory_impl"
	log "github.com/sirupsen/logrus"
)

func main() {
	fact := factory_impl.NewFactory()

	user := &email.ImapUser{
		ImapServer: "imap.yandex.ru:993",
		Login:      "kirillkiling@yandex.ru",
		Password:   "hvitldgmynqhsvol",
	}

	imap := fact.ImapEmail()
	emails, err := imap.ReadUnseenEmails(context.Background(), user)
	if err != nil {
		log.Fatalln(err)

	}

	for _, email := range emails {
		if email.Uid == 651 || email.Uid == 677 {

		} else {
			continue
		}

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
