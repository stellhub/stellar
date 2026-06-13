package internal

import (
	"context"
	"log/slog"

	"github.com/stellhub/stellar"
)

type serverInterceptorKey struct{}

func ServerInterceptors() []stellar.InterceptorDefinition {
	return []stellar.InterceptorDefinition{
		stellar.GRPCServerInterceptor("example.grpc.server.audit", 10, func(ctx context.Context, inv *stellar.InterceptorInvocation, req any, next stellar.InterceptorHandler) (any, error) {
			if inv != nil {
				slog.InfoContext(ctx, "grpc server interceptor before",
					"name", "example.grpc.server.audit",
					"method", inv.Path,
				)
			}
			resp, err := next(ctx, inv, req)
			slog.InfoContext(ctx, "grpc server interceptor after",
				"name", "example.grpc.server.audit",
				"error", err,
			)
			return resp, err
		}),
		stellar.GRPCServerInterceptor("example.grpc.server.context", 20, func(ctx context.Context, inv *stellar.InterceptorInvocation, req any, next stellar.InterceptorHandler) (any, error) {
			ctx = context.WithValue(ctx, serverInterceptorKey{}, "grpc-server")
			return next(ctx, inv, req)
		}),
	}
}

func serverInterceptorValue(ctx context.Context) string {
	value, _ := ctx.Value(serverInterceptorKey{}).(string)
	return value
}
