package main

import (
	"log"

	"github.com/stellhub/stellar"
	"github.com/stellhub/stellar/examples/opentelemetry/custom-metrics/internal"
)

func main() {
	if err := stellar.Run(stellar.WithStarter(internal.NewMetricsStarter())); err != nil {
		log.Fatal(err)
	}
}
