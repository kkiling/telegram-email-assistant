package main

import (
	"fmt"
	"strconv"
	"time"

	"github.com/kiling91/telegram-email-assistant/internal/factory/factory_impl"
	log "github.com/sirupsen/logrus"
	tg "gopkg.in/telebot.v3"
)

type UserProfile struct {
	TelegramID int64 `json:"telegram_id"`
}

func (p *UserProfile) Recipient() string {
	return strconv.FormatInt(p.TelegramID, 10)
}

func mainProcess(allowedUserId int64, b *tg.Bot) {

	time.Sleep(1 * time.Second)
	for {
		menu := &tg.ReplyMarkup{}

		btnHelp := menu.Data("ℹ Help", "btn_help")
		btnSettings := menu.Data("⚙ Settings", "btn_settings")
		menu.Inline(
			menu.Row(btnHelp),
			menu.Row(btnSettings),
		)

		// On inline button pressed (callback)
		b.Handle(&btnHelp, func(c tg.Context) error {
			log.Println(c.Callback().Unique)
			return c.Respond()
		})

		b.Handle(&btnSettings, func(c tg.Context) error {
			log.Println(c.Callback().Unique)
			return c.Respond()
		})

		if _, err := b.Send(&UserProfile{
			TelegramID: 594785598,
		}, "time", menu); err != nil {
			log.Errorf("Unable to send message %v", err)
		}
		return
		// time.Sleep(10 * time.Second)
	}
}

func main() {
	fact := factory_impl.NewFactory()
	cfg := fact.Config()

	pref := tg.Settings{
		Token:  cfg.Telegram.BotToken,
		Poller: &tg.LongPoller{Timeout: 10 * time.Second},
	}

	b, err := tg.NewBot(pref)
	if err != nil {
		log.Fatal(err)
		return
	}

	b.Handle("/start", func(c tg.Context) error {
		userId := c.Sender().ID
		if userId != cfg.Telegram.AllowedUserId {
			return c.Send(fmt.Sprintf("❗ Access is denied: your id #%d", userId))
		}
		return c.Send("Hello!")
	})

	b.Handle(tg.OnText, func(c tg.Context) error {
		return c.Send("🚫 I am not trained to respond to messages or commands")
	})

	go mainProcess(cfg.Telegram.AllowedUserId, b)

	b.Start()
}
