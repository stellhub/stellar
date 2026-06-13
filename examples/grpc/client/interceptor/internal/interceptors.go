package internal

import (
	"context"
	"log/slog"

	"github.com/stellhub/stellar"
	"google.golang.org/grpc/metadata"
)

func ClientInterceptors() []stellar.InterceptorDefinition {
	return []stellar.InterceptorDefinition{
		stellar.GRPCClientInterceptor("example.grpc.client.metadata", 10, func(ctx context.Context, inv *stellar.InterceptorInvocation, req any, next stellar.InterceptorHandler) (any, error) {
			ctx = metadata.AppendToOutgoingContext(ctx, "x-example-interceptor", "grpc-client")
			if inv != nil {
				slog.InfoContext(ctx, "grpc client interceptor before",
					"name", "example.grpc.client.metadata",
					"method", inv.Path,
					"target", inv.Target,
				)
			}
			return next(ctx, inv, req)
		}),
		stellar.GRPCClientInterceptor("example.grpc.client.audit", 20, func(ctx context.Context, inv *stellar.InterceptorInvocation, req any, next stellar.InterceptorHandler) (any, error) {
			resp, err := next(ctx, inv, req)
			slog.InfoContext(ctx, "grpc client interceptor after",
				"name", "example.grpc.client.audit",
				"error", err,
			)
			return resp, err
		}),
	}
}
