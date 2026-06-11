package main

import (
	"context"
	"log"

	stellflux "github.com/stellhub/stellflux-go"
)

func main() {
	app := stellflux.New(stellflux.Config{
		AppName:     "basic-example",
		Environment: stellflux.EnvDev,
		Zone:        "local",
	})
	app.Use(stellflux.StandardModules()...)

	if err := app.Start(context.Background()); err != nil {
		log.Fatal(err)
	}
	defer app.Stop(context.Background())
}
