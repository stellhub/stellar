package internal

import (
	"context"
	"log/slog"
	"time"

	"github.com/stellhub/stellar"
	grpcgoadapter "github.com/stellhub/stellar/transport/grpc/adapters/grpcgo"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/structpb"
)

const mockGRPCServiceAddr = "127.0.0.1:19097"

type clientStarter struct {
	config stellar.Config
	logger *slog.Logger
	conn   *grpc.ClientConn
	mock   *MockGRPCService
}

func NewClientStarter() *clientStarter {
	return &clientStarter{}
}

func (s *clientStarter) Name() string {
	return "grpc-client-interceptor-example"
}

func (s *clientStarter) Condition(ctx stellar.StarterContext) bool {
	s.config = ctx.Config()
	return true
}

func (s *clientStarter) Init(ctx context.Context, app *stellar.App) error {
	s.logger = app.Logger()
	conn, _, err := grpcgoadapter.NewNamedClientConnFromConfig(
		ctx,
		s.config.GRPC.Client,
		"user-service",
		app.Observability(),
		grpcgoadapter.WithInterceptors(app.Interceptors()),
	)
	if err != nil {
		return err
	}

	s.conn = conn
	s.mock = NewMockGRPCService(mockGRPCServiceAddr, s.logger)
	return nil
}

func (s *clientStarter) Start(ctx context.Context) error {
	if err := s.mock.Start(ctx); err != nil {
		return err
	}
	defer func() {
		stopCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := s.mock.Stop(stopCtx); err != nil {
			s.logger.Error("mock grpc service shutdown failed", "error", err)
		}
	}()

	requestCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	reply := new(structpb.Struct)
	if err := s.conn.Invoke(requestCtx, "/"+interceptorServiceName+"/Ping", &emptypb.Empty{}, reply); err != nil {
		return err
	}
	s.logger.Info("grpc response loaded",
		"message", reply.Fields["message"].GetStringValue(),
		"interceptor_header", reply.Fields["interceptor_header"].GetStringValue(),
	)
	return nil
}

func (s *clientStarter) Stop(context.Context) error {
	if s.conn == nil {
		return nil
	}
	return s.conn.Close()
}
