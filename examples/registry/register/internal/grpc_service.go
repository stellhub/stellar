package internal

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/structpb"
)

const registryServiceName = "stellar.examples.registry.RegistryExample"

type registryService struct{}

type registryServer interface {
	Ping(context.Context, *emptypb.Empty) (*structpb.Struct, error)
}

func registryServiceDesc() *grpc.ServiceDesc {
	return &grpc.ServiceDesc{
		ServiceName: registryServiceName,
		HandlerType: (*registryServer)(nil),
		Methods: []grpc.MethodDesc{
			{
				MethodName: "Ping",
				Handler:    registryPingHandler,
			},
		},
		Streams:  []grpc.StreamDesc{},
		Metadata: "examples/registry/register",
	}
}

func (s *registryService) Ping(context.Context, *emptypb.Empty) (*structpb.Struct, error) {
	return structpb.NewStruct(map[string]any{
		"message":   "pong",
		"transport": "grpc",
		"service":   registryServiceName,
	})
}

func registryPingHandler(
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
		return srv.(registryServer).Ping(ctx, req)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/" + registryServiceName + "/Ping",
	}
	handler := func(ctx context.Context, req any) (any, error) {
		return srv.(registryServer).Ping(ctx, req.(*emptypb.Empty))
	}
	return interceptor(ctx, req, info, handler)
}
