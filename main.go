package main

import (
	"github.com/kiling91/telegram-email-assistant/pkg/factory/factory_impl"
	"github.com/kiling91/telegram-email-assistant/pkg/types"
	"log"
)

func main() {
	fact := factory_impl.NewFactory()

	user := &types.EmailUser{
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

	imap := fact.ImapServer()
	err := imap.ReadUnseenEmails(user)
	if err != nil {
		log.Fatalln(err)
	}
}
