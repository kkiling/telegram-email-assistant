package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"syscall"

	"github.com/kiling91/telegram-email-assistant/internal/app"
	"github.com/sirupsen/logrus"
)

func main() {
	logrus.SetLevel(logrus.DebugLevel)

	configFile := flag.String("config", "configs/config.yml", "Path to config file.")
	flag.Parse()

	ctx, cancel := context.WithCancel(context.Background())
	a := app.NewApp(*configFile)
	// Gracefully shutdown
	go func() {
		sig := make(chan os.Signal, 1)
		signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
		<-sig
		logrus.Println("Shutdown bot")
		cancel()
		a.Shutdown()
	}()

	logrus.Println("Start bot")
	a.Start(ctx)
}
