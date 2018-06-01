package main

import (
	"log"

	"github.com/ifreddyrondon/capture/app"
	"github.com/ifreddyrondon/capture/database"
)

func main() {
	ds := database.Open("postgres://localhost/captures_app?sslmode=disable")
	app := app.New(ds.DB)
	app.RegisterOnShutdown(ds.OnShutdown)
	if err := app.Serve(); err != nil {
		log.Printf("%v", err)
	}
}
