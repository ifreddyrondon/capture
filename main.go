package main

import (
	"fmt"
	"os"

	"github.com/ifreddyrondon/capture/app"
	"github.com/ifreddyrondon/capture/database"
)

func main() {
	ds := database.Open("postgres://localhost/captures_app?sslmode=disable")
	bastion := app.New(ds.DB)
	bastion.RegisterOnShutdown(ds.OnShutdown)
	if err := bastion.Serve(); err != nil {
		fmt.Fprintf(os.Stderr, "%v", err)
		os.Exit(1)
	}
}
