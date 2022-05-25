package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/kiling91/telegram-email-assistant/internal/app"
)

func main() {
	configFile := flag.String("config", "config/config.yml", "Path to config file.")
	a := app.NewApp(*configFile)
	// Gracefully shutdown
	go func() {
		sig := make(chan os.Signal, 1)
		signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
		<-sig
		log.Println("Shutdown bot")
		a.Shutdown()
	}()

	log.Println("Start bot")
	a.Start()
}
