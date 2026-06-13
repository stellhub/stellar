package internal

import (
	"context"
	"errors"
	"log/slog"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/structpb"
)

const interceptorServiceName = "stellar.examples.grpc.Interceptor"

type MockGRPCService struct {
	addr     string
	logger   *slog.Logger
	server   *grpc.Server
	listener net.Listener
}

func NewMockGRPCService(addr string, logger *slog.Logger) *MockGRPCService {
	if logger == nil {
		logger = slog.Default()
	}
	return &MockGRPCService{
		addr:   addr,
		logger: logger,
	}
}

func (s *MockGRPCService) Start(context.Context) error {
	listener, err := net.Listen("tcp", s.addr)
	if err != nil {
		return err
	}
	s.listener = listener
	s.server = grpc.NewServer()
	s.server.RegisterService(interceptorServiceDesc(), &mockInterceptorService{})

	go func() {
		if err := s.server.Serve(listener); err != nil && !errors.Is(err, grpc.ErrServerStopped) {
			s.logger.Error("mock grpc service stopped", "error", err)
		}
	}()
	s.logger.Info("mock grpc service started", "addr", s.addr)
	return nil
}

func (s *MockGRPCService) Stop(context.Context) error {
	if s.server == nil {
		return nil
	}
	s.server.GracefulStop()
	return nil
}

type mockInterceptorService struct{}

func interceptorServiceDesc() *grpc.ServiceDesc {
	return &grpc.ServiceDesc{
		ServiceName: interceptorServiceName,
		HandlerType: (*interceptorServer)(nil),
		Methods: []grpc.MethodDesc{
			{
				MethodName: "Ping",
				Handler:    pingHandler,
			},
		},
		Streams:  []grpc.StreamDesc{},
		Metadata: "examples/grpc/client/interceptor",
	}
}

type interceptorServer interface {
	Ping(context.Context, *emptypb.Empty) (*structpb.Struct, error)
}

func (s *mockInterceptorService) Ping(ctx context.Context, _ *emptypb.Empty) (*structpb.Struct, error) {
	header := ""
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		values := md.Get("x-example-interceptor")
		if len(values) > 0 {
			header = values[0]
		}
	}
	return structpb.NewStruct(map[string]any{
		"message":            "pong",
		"interceptor_header": header,
	})
}

func pingHandler(
	srv any,
	ctx context.Context,
	dec func(any) error,
	interceptor grpc.UnaryServerInterceptor,
) (any, error) {
	req := new(emptypb.Empty)
	if err := dec(req); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(interceptorServer).Ping(ctx, req)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/" + interceptorServiceName + "/Ping",
	}
	handler := func(ctx context.Context, req any) (any, error) {
		return srv.(interceptorServer).Ping(ctx, req.(*emptypb.Empty))
	}
	return interceptor(ctx, req, info, handler)
}
