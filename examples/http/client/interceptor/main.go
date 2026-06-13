package main

import (
	"log"

	"github.com/stellhub/stellar"
	"github.com/stellhub/stellar/examples/http/client/interceptor/internal"
)

func main() {
	if _, err := stellar.Start(
		stellar.WithInterceptor(internal.ClientInterceptors()...),
		stellar.WithStarter(internal.NewClientStarter()),
	); err != nil {
		log.Fatal(err)
	}
}
