package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/motru4/legendary-currency-wallet/services/exchanger/internal/app"
	"github.com/motru4/legendary-currency-wallet/services/exchanger/internal/config"
)

func main() {
	// Flag to specify the path to .env (local) - default config.env
	configPath := flag.String("c", "config.env", "path to config env file")
	flag.Parse()

	cfg := config.New(*configPath)

	application, err := app.New(cfg)
	if err != nil {
		log.Fatalf("failed to init app: %v", err)
	}

	go func() {
		if err := application.Run(); err != nil {
			log.Fatalf("application run error: %v", err)
		}
	}()

	fmt.Println("exchanger service started")

	// We catch signals for graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	<-stop
	fmt.Println("shutting down exchanger...")

	if err := application.Stop(); err != nil {
		log.Printf("error during shutdown: %v", err)
	}
}
