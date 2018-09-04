package main

import (
	"fmt"
	"log"
	"os"

	"github.com/ifreddyrondon/capture/app"
	"github.com/ifreddyrondon/capture/internal/config"
)

func main() {
	cfg, err := config.New()
	if err != nil {
		log.Panicln("Configuration error", err)
	}
	bastion := app.New(cfg)
	bastion.RegisterOnShutdown(cfg.OnShutdown)
	if err := bastion.Serve(); err != nil {
		fmt.Fprintf(os.Stderr, "%v", err)
		os.Exit(1)
	}
}
