package main

import (
	"log"

	"github.com/ifreddyrondon/gocapture/app"
	"github.com/ifreddyrondon/gocapture/database"
)

func main() {
	ds := database.Open("postgres://localhost/captures_app?sslmode=disable")
	err := app.New(ds.DB).Serve()
	if err != nil {
		log.Printf("%v", err)
	}
}
