package main

import (
	"fmt"
	"log"
	"os"

	"github.com/ifreddyrondon/bastion"

	"github.com/ifreddyrondon/capture/config"
	"github.com/ifreddyrondon/capture/pkg/http/rest"
)

func main() {
	cfg, err := config.New()
	if err != nil {
		log.Panicln("Configuration error", err)
	}

	app := bastion.New()
	app.Mount("/", rest.Router(cfg.Resources))
	app.RegisterOnShutdown(cfg.OnShutdown)
	fmt.Fprintln(os.Stderr, app.Serve(cfg.ADDR))
}
