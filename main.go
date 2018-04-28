package main

import (
	"log"

	"github.com/ifreddyrondon/gocapture/app"
	"github.com/ifreddyrondon/gocapture/database"
)

func main() {
	ds := database.Open("postgres://localhost/captures_app?sslmode=disable")
	app := app.New(ds.DB)
	app.RegisterOnShutdown(ds.OnShutdown)
	if err := app.Serve(); err != nil {
		log.Printf("%v", err)
	}
}
