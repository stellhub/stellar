package http_test

import (
	"context"
	"encoding/json"
	stdhttp "net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/stellhub/stellar/interceptor"
	"github.com/stellhub/stellar/middleware"
	stellarhttp "github.com/stellhub/stellar/transport/http"
	ginadapter "github.com/stellhub/stellar/transport/http/adapters/gin"
)

type pingRequest struct{}

type pingResponse struct {
	Message string `json:"message"`
}

func TestRouterUsesInterceptorsBeforeRouteMiddleware(t *testing.T) {
	registry := interceptor.NewRegistry()
	order := make([]string, 0, 3)
	registry.Register(interceptor.Business(interceptor.KindHTTPServer, "business", 10, interceptor.New("business", func(ctx context.Context, inv *interceptor.Invocation, req any, next interceptor.Handler) (any, error) {
		order = append(order, "interceptor")
		return next(ctx, inv, req)
	})))
	router := stellarhttp.NewRouter(stellarhttp.WithInterceptors(registry))
	router.Use(func(next middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req any) (any, error) {
			order = append(order, "middleware")
			return next(ctx, req)
		}
	})
	router.GET("/ping", func(context.Context, *stellarhttp.Request) (*stellarhttp.Response, error) {
		order = append(order, "handler")
		return stellarhttp.JSON(stdhttp.StatusOK, pingResponse{Message: "pong"}), nil
	})

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(stdhttp.MethodGet, "/ping", nil)

	adapter := ginadapter.New("")
	adapter.UseRouter(router)
	adapter.Handler().ServeHTTP(recorder, request)

	if recorder.Code != stdhttp.StatusOK {
		t.Fatalf("expected status %d, got %d", stdhttp.StatusOK, recorder.Code)
	}
	if !reflect.DeepEqual(order, []string{"interceptor", "middleware", "handler"}) {
		t.Fatalf("unexpected order: %#v", order)
	}
}

func TestTypedRouteRunsBusinessInterceptorsAfterBinding(t *testing.T) {
	registry := interceptor.NewRegistry()
	order := make([]string, 0, 4)
	registry.Register(
		interceptor.Framework(interceptor.KindHTTPServer, interceptor.StageAdmission, "framework", interceptor.New("framework", func(ctx context.Context, inv *interceptor.Invocation, req any, next interceptor.Handler) (any, error) {
			order = append(order, "framework")
			return next(ctx, inv, req)
		})),
		interceptor.Framework(interceptor.KindHTTPServer, interceptor.StageDecodeValidate, "decode", interceptor.New("decode", func(ctx context.Context, inv *interceptor.Invocation, req any, next interceptor.Handler) (any, error) {
			order = append(order, "decode")
			return next(ctx, inv, req)
		})),
		interceptor.Business(interceptor.KindHTTPServer, "business", 10, interceptor.New("business", func(ctx context.Context, inv *interceptor.Invocation, req any, next interceptor.Handler) (any, error) {
			order = append(order, "business")
			return next(ctx, inv, req)
		})),
	)
	router := stellarhttp.NewRouter(stellarhttp.WithInterceptors(registry))
	stellarhttp.Handle(router, stdhttp.MethodGet, "/typed", func(*stellarhttp.Request) (*pingRequest, error) {
		order = append(order, "binder")
		return &pingRequest{}, nil
	}, func(context.Context, *pingRequest) (*pingResponse, error) {
		order = append(order, "handler")
		return &pingResponse{Message: "pong"}, nil
	}, stellarhttp.JSONEncoder[pingResponse])

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(stdhttp.MethodGet, "/typed", nil)

	adapter := ginadapter.New("")
	adapter.UseRouter(router)
	adapter.Handler().ServeHTTP(recorder, request)

	if recorder.Code != stdhttp.StatusOK {
		t.Fatalf("expected status %d, got %d", stdhttp.StatusOK, recorder.Code)
	}
	if !reflect.DeepEqual(order, []string{"framework", "binder", "decode", "business", "handler"}) {
		t.Fatalf("unexpected order: %#v", order)
	}
}

func TestTypedRouteAndMiddleware(t *testing.T) {
	router := stellarhttp.NewRouter()
	order := make([]string, 0, 2)
	first := func(next middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req any) (any, error) {
			order = append(order, "first")
			return next(ctx, req)
		}
	}
	second := func(next middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req any) (any, error) {
			order = append(order, "second")
			return next(ctx, req)
		}
	}
	router.Use(first)

	stellarhttp.Handle(router, stdhttp.MethodGet, "/ping", stellarhttp.EmptyBinder[pingRequest](), func(context.Context, *pingRequest) (*pingResponse, error) {
		return &pingResponse{Message: "pong"}, nil
	}, stellarhttp.JSONEncoder[pingResponse], second)

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(stdhttp.MethodGet, "/ping", nil)

	adapter := ginadapter.New("")
	adapter.UseRouter(router)
	adapter.Handler().ServeHTTP(recorder, request)

	if recorder.Code != stdhttp.StatusOK {
		t.Fatalf("expected status %d, got %d", stdhttp.StatusOK, recorder.Code)
	}
	if !reflect.DeepEqual(order, []string{"first", "second"}) {
		t.Fatalf("unexpected middleware order: %#v", order)
	}

	var response pingResponse
	if err := json.NewDecoder(recorder.Body).Decode(&response); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if response.Message != "pong" {
		t.Fatalf("expected pong, got %q", response.Message)
	}
}
