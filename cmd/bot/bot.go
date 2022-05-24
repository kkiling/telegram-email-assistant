package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/kiling91/telegram-email-assistant/internal/bot"
	"github.com/kiling91/telegram-email-assistant/internal/factory/factory_impl"
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
