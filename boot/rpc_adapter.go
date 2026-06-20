package boot

import (
	"context"

	"github.com/stellhub/stellar/loadbalancer"
	bootgrpc "github.com/stellhub/stellar/transport/grpc"
	grpcgoadapter "github.com/stellhub/stellar/transport/grpc/adapters/grpcgo"
	"google.golang.org/grpc"
)

func WithRPCAdapter(adapter bootgrpc.Adapter) Option {
	return func(app *App) {
		app.setRPCAdapter(adapter, true)
	}
}

func WithRPCServer(addr string) Option {
	return func(app *App) {
		adapter := app.rpcAdapter
		if adapter == nil {
			adapter = grpcgoadapter.New(addr)
		}
		if setter, ok := adapter.(bootgrpc.AddrSetter); ok {
			setter.SetAddr(addr)
		}
		app.setRPCAdapter(adapter, true)
	}
}

func (a *App) NewGRPCClient(ctx context.Context, name string) (*grpc.ClientConn, string, error) {
	cfg := a.Config()
	var router loadbalancer.Router
	if a.governanceRouteEnabled {
		governanceRouter := loadbalancer.NewGovernanceRouter(a.governance, nil)
		governanceRouter.Metrics = a.governanceMetrics
		router = governanceRouter
	}
	return grpcgoadapter.NewNamedClientConnFromConfig(
		ctx,
		grpcClientConfigWithDiscovery(cfg),
		name,
		a.observability,
		grpcgoadapter.WithInterceptors(a.interceptors),
		grpcgoadapter.WithRouter(router),
	)
}

func (a *App) setRPCAdapter(adapter bootgrpc.Adapter, registerTransport bool) {
	if adapter == nil {
		return
	}
	if consumer, ok := adapter.(observabilityConsumer); ok {
		consumer.UseObservability(a.observability)
	}
	if consumer, ok := adapter.(interceptorConsumer); ok {
		consumer.UseInterceptors(a.interceptors)
	}
	if consumer, ok := adapter.(serviceNameConsumer); ok {
		consumer.UseServiceName(a.config.AppName)
	}
	a.rpcAdapter = adapter
	if registerTransport {
		a.addTransport(adapter)
	}
}
