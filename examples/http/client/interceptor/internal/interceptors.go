package internal

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/stellhub/stellar"
)

func ClientInterceptors() []stellar.InterceptorDefinition {
	return []stellar.InterceptorDefinition{
		stellar.HTTPClientInterceptor("example.http.client.header", 10, func(ctx context.Context, inv *stellar.InterceptorInvocation, req any, next stellar.InterceptorHandler) (any, error) {
			if request, ok := req.(*http.Request); ok {
				request.Header.Set("X-Example-Interceptor", "http-client")
			}
			if inv != nil {
				slog.InfoContext(ctx, "http client interceptor before",
					"name", "example.http.client.header",
					"service", inv.Service,
					"target", inv.Target,
				)
			}
			return next(ctx, inv, req)
		}),
		stellar.HTTPClientInterceptor("example.http.client.audit", 20, func(ctx context.Context, inv *stellar.InterceptorInvocation, req any, next stellar.InterceptorHandler) (any, error) {
			resp, err := next(ctx, inv, req)
			slog.InfoContext(ctx, "http client interceptor after",
				"name", "example.http.client.audit",
				"error", err,
			)
			return resp, err
		}),
	}
}
