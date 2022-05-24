package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/kiling91/telegram-email-assistant/internal/factory/factory_impl"
	"github.com/kiling91/telegram-email-assistant/pkg/bot"
	"github.com/sirupsen/logrus"
)

func main() {
	fact := factory_impl.NewFactory()

	b := fact.Bot()
	b.Handle("/start", func(ctx bot.Context) error {
		_, err := b.Send(ctx.UserId(), "hello")
		return err
	})

	b.Handle("/test", func(ctx bot.Context) error {

		inline := bot.NewInline(2, func(ctx bot.BtnContext) error {
			b.Send(ctx.UserId(), ctx.Data())
			return nil
		})

		inline.Add("⚙ Settings", "btn_settings", "{settings}")
		inline.Add("? Help", "btn_help", "{help}")

		_, err := b.Send(ctx.UserId(), "test btn", inline)
		return err
	})

	b.Handle("/time", func(ctx bot.Context) error {
		edit, err := b.Send(ctx.UserId(), "test timer ...")
		go func() {
			index := 0
			for {
				time.Sleep(time.Second)
				index++
				_, err = b.Edit(edit, fmt.Sprintf("test timer %d", index))
				logrus.Warnf("error edit: %v", err)
			}
		}()
		return err
	})

	b.Handle("/photo", func(ctx bot.Context) error {
		_, err := b.SendPhoto(ctx.UserId(), &bot.Photo{
			Filename: "test.jpg",
			Caption:  "Какой то текст....",
		})
		return err
	})

	b.Handle("/doc", func(ctx bot.Context) error {
		return b.SendDocument(ctx.UserId(), "test.jpg")
	})

	// Gracefully shutdown
	go func() {
		sig := make(chan os.Signal, 1)
		signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
		<-sig
		log.Println("shutdown bot")
		b.Stop()
	}()

	log.Println("start bot")
	b.Start()
}
