package interceptor_test

import (
	"context"
	"reflect"
	"testing"

	"github.com/stellhub/stellar/interceptor"
)

func TestRegistryKeepsFrameworkStageOrderAndBusinessOrder(t *testing.T) {
	registry := interceptor.NewRegistry()
	order := make([]string, 0)

	registry.Register(
		interceptor.Framework(interceptor.KindHTTPServer, interceptor.StageSecurity, "security", record("security", &order)),
		interceptor.Business(interceptor.KindHTTPServer, "biz-20", 20, record("biz-20", &order)),
		interceptor.Framework(interceptor.KindHTTPServer, interceptor.StageAdmission, "admission", record("admission", &order)),
		interceptor.Business(interceptor.KindHTTPServer, "biz-10", 10, record("biz-10", &order)),
		interceptor.Framework(interceptor.KindHTTPServer, interceptor.StageDecodeValidate, "decode", record("decode", &order)),
	)

	handler := registry.Chain(interceptor.KindHTTPServer, func(context.Context, *interceptor.Invocation, any) (any, error) {
		order = append(order, "handler")
		return "ok", nil
	})
	resp, err := handler(context.Background(), &interceptor.Invocation{Kind: interceptor.KindHTTPServer}, nil)
	if err != nil {
		t.Fatalf("chain failed: %v", err)
	}
	if resp != "ok" {
		t.Fatalf("unexpected response %v", resp)
	}

	expected := []string{"admission", "security", "decode", "biz-10", "biz-20", "handler"}
	if !reflect.DeepEqual(order, expected) {
		t.Fatalf("unexpected order: %#v", order)
	}
}

func TestDefaultFrameworkAddsRequestID(t *testing.T) {
	registry := interceptor.NewRegistry(interceptor.WithDefaultFramework())
	handler := registry.Chain(interceptor.KindGRPCClient, func(ctx context.Context, inv *interceptor.Invocation, _ any) (any, error) {
		if inv.RequestID == "" {
			t.Fatalf("expected invocation request id")
		}
		if interceptor.RequestIDFromContext(ctx) != inv.RequestID {
			t.Fatalf("expected request id in context")
		}
		return nil, nil
	})

	if _, err := handler(context.Background(), &interceptor.Invocation{Kind: interceptor.KindGRPCClient}, nil); err != nil {
		t.Fatalf("chain failed: %v", err)
	}
}

func record(name string, order *[]string) interceptor.Interceptor {
	return interceptor.New(name, func(ctx context.Context, inv *interceptor.Invocation, req any, next interceptor.Handler) (any, error) {
		*order = append(*order, name)
		return next(ctx, inv, req)
	})
}
