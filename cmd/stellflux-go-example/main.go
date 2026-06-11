package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	stellflux "github.com/stellhub/stellflux-go"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	app := stellflux.New(stellflux.Config{
		AppName:     "stellflux-go-example",
		Environment: stellflux.EnvDev,
		Zone:        "local",
	}, stellflux.WithLogger(logger))
	app.Use(stellflux.StandardModules()...)

	if err := app.Start(context.Background()); err != nil {
		logger.Error("failed to start StellFlux Go", "error", err)
		os.Exit(1)
	}

	server := &http.Server{
		Addr:              getenv("STELLFLUX_HTTP_ADDR", ":8080"),
		Handler:           app.Handler(),
		ReadHeaderTimeout: 5 * time.Second,
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	go func() {
		logger.Info("starting StellFlux Go example", "address", server.Addr)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Error("example server stopped unexpectedly", "error", err)
			os.Exit(1)
		}
	}()

	<-ctx.Done()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		logger.Error("failed to stop example server", "error", err)
		os.Exit(1)
	}
	if err := app.Stop(shutdownCtx); err != nil {
		logger.Error("failed to stop StellFlux Go", "error", err)
		os.Exit(1)
	}
}

func getenv(key string, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}
