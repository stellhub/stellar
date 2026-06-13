package internal

import (
	"context"
	"log/slog"

	"github.com/stellhub/stellar"
)

func ServerInterceptors() []stellar.InterceptorDefinition {
	return []stellar.InterceptorDefinition{
		stellar.HTTPServerInterceptor("example.http.server.audit", 10, func(ctx context.Context, inv *stellar.InterceptorInvocation, req any, next stellar.InterceptorHandler) (any, error) {
			if inv != nil {
				slog.InfoContext(ctx, "http server interceptor before",
					"name", "example.http.server.audit",
					"method", inv.Method,
					"path", inv.Path,
				)
			}
			resp, err := next(ctx, inv, req)
			slog.InfoContext(ctx, "http server interceptor after",
				"name", "example.http.server.audit",
				"error", err,
			)
			return resp, err
		}),
		stellar.HTTPServerInterceptor("example.http.server.marker", 20, func(ctx context.Context, inv *stellar.InterceptorInvocation, req any, next stellar.InterceptorHandler) (any, error) {
			if inv != nil {
				inv.SetAttribute("example.interceptor", "http-server")
			}
			return next(ctx, inv, req)
		}),
	}
}
