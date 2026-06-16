package main

import (
	"log"

	"github.com/stellhub/stellar"
	"github.com/stellhub/stellar/examples/message/internal"
)

func main() {
	if err := stellar.Run(stellar.WithStarter(internal.NewMessageStarter("message-stellflow-quickstart"))); err != nil {
		log.Fatal(err)
	}
}
