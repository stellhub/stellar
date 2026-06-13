package main

import (
	"log"

	"github.com/stellhub/stellar"
	"github.com/stellhub/stellar/examples/grpc/server/interceptor/internal"
)

func main() {
	if err := stellar.Run(
		stellar.WithInterceptor(internal.ServerInterceptors()...),
		stellar.WithStarter(internal.NewServerStarter()),
	); err != nil {
		log.Fatal(err)
	}
}
