package internal

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/stellhub/stellar"
	stellargrpc "github.com/stellhub/stellar/transport/grpc"
	stellarhttp "github.com/stellhub/stellar/transport/http"
)

type registryStarter struct {
	cfg stellar.Config
}

func NewRegistryStarter() *registryStarter {
	return &registryStarter{}
}

func (s *registryStarter) Name() string {
	return "registry-register-example"
}

func (s *registryStarter) Condition(ctx stellar.StarterContext) bool {
	s.cfg = ctx.Config()
	return true
}

func (s *registryStarter) Init(_ context.Context, app *stellar.App) error {
	api := app.HTTP().Group("/api/v1")
	api.GET("/ping", handleHTTPPing)
	api.GET("/registry/instance", s.handleConfiguredInstance)
	api.GET("/registry/discover", s.handleDiscover(app))

	rpc := app.RPC()
	if rpc == nil {
		return errors.New("registry example requires grpc.server.enabled=true")
	}
	return rpc.Register(stellargrpc.Service{
		Description:    registryServiceDesc(),
		Implementation: &registryService{},
	})
}

func (s *registryStarter) Start(context.Context) error {
	return nil
}

func (s *registryStarter) Stop(context.Context) error {
	return nil
}

func handleHTTPPing(context.Context, *stellarhttp.Request) (*stellarhttp.Response, error) {
	return stellarhttp.JSON(http.StatusOK, map[string]any{
		"message":   "pong",
		"transport": "http",
	}), nil
}

func (s *registryStarter) handleConfiguredInstance(context.Context, *stellarhttp.Request) (*stellarhttp.Response, error) {
	if s.cfg.Registry == nil {
		return stellarhttp.JSON(http.StatusNotFound, map[string]string{
			"message": "registry is not configured",
		}), nil
	}
	return stellarhttp.JSON(http.StatusOK, map[string]any{
		"adapter":    s.cfg.Registry.Adapter,
		"namespace":  s.cfg.Registry.Namespace,
		"service":    s.cfg.Registry.Service,
		"instanceId": s.cfg.Registry.InstanceID,
		"zone":       s.cfg.Registry.Zone,
		"endpoints":  s.cfg.Registry.ServiceEndpoints,
	}), nil
}

func (s *registryStarter) handleDiscover(app *stellar.App) stellarhttp.Endpoint {
	return func(ctx context.Context, req *stellarhttp.Request) (*stellarhttp.Response, error) {
		registry, ok := app.ServiceRegistry()
		if !ok {
			return stellarhttp.JSON(http.StatusNotFound, map[string]string{
				"message": "registry is not configured",
			}), nil
		}
		namespace := "default"
		service := ""
		if s.cfg.Registry != nil {
			namespace = s.cfg.Registry.Namespace
			service = s.cfg.Registry.Service
		}
		if value := strings.TrimSpace(req.Query.Get("namespace")); value != "" {
			namespace = value
		}
		if value := strings.TrimSpace(req.Query.Get("service")); value != "" {
			service = value
		}
		instances, err := registry.Discover(ctx, stellar.ServiceQuery{
			Namespace:   namespace,
			Service:     service,
			PassingOnly: true,
		})
		if err != nil {
			return nil, err
		}
		return stellarhttp.JSON(http.StatusOK, map[string]any{
			"namespace": namespace,
			"service":   service,
			"instances": instances,
		}), nil
	}
}
