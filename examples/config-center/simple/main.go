package main

import (
	"log"

	"github.com/stellhub/stellar"
	"github.com/stellhub/stellar/examples/config-center/simple/internal"
)

func main() {
	if err := stellar.Run(stellar.WithStarter(internal.NewConfigCenterStarter())); err != nil {
		log.Fatal(err)
	}
}
