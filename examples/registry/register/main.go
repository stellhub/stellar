package main

import (
	"log"

	"github.com/stellhub/stellar"
	"github.com/stellhub/stellar/examples/registry/register/internal"
)

func main() {
	if err := stellar.Run(stellar.WithStarter(internal.NewRegistryStarter())); err != nil {
		log.Fatal(err)
	}
}
