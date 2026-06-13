package internal

import (
	"context"
	"net/http"

	"github.com/stellhub/stellar"
	stellarhttp "github.com/stellhub/stellar/transport/http"
)

type serverStarter struct{}

func NewServerStarter() *serverStarter {
	return &serverStarter{}
}

func (s *serverStarter) Name() string {
	return "http-server-interceptor-example"
}

func (s *serverStarter) Condition(stellar.StarterContext) bool {
	return true
}

func (s *serverStarter) Init(_ context.Context, app *stellar.App) error {
	api := app.HTTP().Group("/api/v1")
	api.GET("/ping", handlePing)
	stellarhttp.Handle(
		api,
		http.MethodPost,
		"/echo",
		stellarhttp.JSONBinder[echoRequest](),
		echo,
		stellarhttp.JSONEncoder[echoResponse],
	)
	return nil
}

func (s *serverStarter) Start(context.Context) error {
	return nil
}

func (s *serverStarter) Stop(context.Context) error {
	return nil
}
