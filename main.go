package main

import (
	"context"
	"github.com/kiling91/telegram-email-assistant/internal/email"
	"github.com/kiling91/telegram-email-assistant/internal/factory/factory_impl"
	"log"
)

func main() {
	fact := factory_impl.NewFactory()

	user := &email.ImapUser{
		ImapServer: "imap.yandex.ru:993",
		Login:      "kirillkiling@yandex.ru",
		Password:   "zxishjxtaufdfnvk",
	}

	/*storage := fact.GetStorage()
	userUID, err := storage.AddUser(&types.EmailUser{
		ImapServer: "",
		Login:      "",
		Password:   "",
	})
	if err != nil {
		log.Fatalln(err)
	}*/

	imap := fact.ImapEmail()
	/*emails, err := imap.ReadUnseenEmails(user)
	if err != nil {
		log.Fatalln(err)
	}*/

	msg, err := imap.ReadEmailBody(context.Background(), user, 37)
	if err != nil {
		return
	}
	log.Println(msg)
}
