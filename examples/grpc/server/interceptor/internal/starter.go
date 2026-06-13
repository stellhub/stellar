package internal

import (
	"context"
	"errors"

	"github.com/stellhub/stellar"
	stellargrpc "github.com/stellhub/stellar/transport/grpc"
)

type serverStarter struct{}

func NewServerStarter() *serverStarter {
	return &serverStarter{}
}

func (s *serverStarter) Name() string {
	return "grpc-server-interceptor-example"
}

func (s *serverStarter) Condition(stellar.StarterContext) bool {
	return true
}

func (s *serverStarter) Init(_ context.Context, app *stellar.App) error {
	rpc := app.RPC()
	if rpc == nil {
		return errors.New("grpc interceptor example requires grpc.server.enabled=true")
	}

	return rpc.Register(stellargrpc.Service{
		Description:    interceptorServiceDesc(),
		Implementation: &interceptorService{},
	})
}

func (s *serverStarter) Start(context.Context) error {
	return nil
}

func (s *serverStarter) Stop(context.Context) error {
	return nil
}
