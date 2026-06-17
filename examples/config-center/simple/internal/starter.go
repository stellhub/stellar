package internal

import (
	"context"
	"net/http"

	"github.com/stellhub/stellar"
	stellarhttp "github.com/stellhub/stellar/transport/http"
)

type configCenterStarter struct {
	cfg stellar.Config
}

func NewConfigCenterStarter() *configCenterStarter {
	return &configCenterStarter{}
}

func (s *configCenterStarter) Name() string {
	return "config-center-simple-example"
}

func (s *configCenterStarter) Condition(ctx stellar.StarterContext) bool {
	s.cfg = ctx.Config()
	return true
}

func (s *configCenterStarter) Init(_ context.Context, app *stellar.App) error {
	api := app.HTTP().Group("/api/v1")
	api.GET("/config-center/status", s.handleStatus(app))
	api.GET("/config-center/config", s.handleConfig)
	api.GET("/config-center/sources", s.handleSources(app))
	return nil
}

func (s *configCenterStarter) Start(context.Context) error {
	return nil
}

func (s *configCenterStarter) Stop(context.Context) error {
	return nil
}

func (s *configCenterStarter) handleStatus(app *stellar.App) stellarhttp.Endpoint {
	return func(_ context.Context, _ *stellarhttp.Request) (*stellarhttp.Response, error) {
		configCenter, ok := app.ConfigCenter()
		if !ok {
			return stellarhttp.JSON(http.StatusNotFound, map[string]any{
				"configured": false,
				"message":    "config center is not configured",
			}), nil
		}
		return stellarhttp.JSON(http.StatusOK, map[string]any{
			"configured": true,
			"adapter":    configCenter.AdapterName(),
			"bootstrap":  bootstrapSummary(s.cfg.ConfigCenter),
		}), nil
	}
}

func (s *configCenterStarter) handleConfig(_ context.Context, _ *stellarhttp.Request) (*stellarhttp.Response, error) {
	return stellarhttp.JSON(http.StatusOK, map[string]any{
		"app": map[string]any{
			"name":    s.cfg.AppName,
			"env":     s.cfg.Environment,
			"zone":    s.cfg.Zone,
			"version": s.cfg.Version,
		},
		"http": map[string]any{
			"server": httpServerSummary(s.cfg.HTTP.Server),
		},
		"mq":       mqSummary(s.cfg.MQ),
		"registry": registrySummary(s.cfg.Registry),
	}), nil
}

func (s *configCenterStarter) handleSources(app *stellar.App) stellarhttp.Endpoint {
	return func(ctx context.Context, _ *stellarhttp.Request) (*stellarhttp.Response, error) {
		configCenter, ok := app.ConfigCenter()
		if !ok {
			return stellarhttp.JSON(http.StatusNotFound, map[string]any{
				"configured": false,
				"message":    "config center is not configured",
			}), nil
		}
		sources, err := configCenter.Load(ctx)
		if err != nil {
			return nil, err
		}
		return stellarhttp.JSON(http.StatusOK, map[string]any{
			"adapter": configCenter.AdapterName(),
			"count":   len(sources),
			"sources": sourceSummaries(sources),
		}), nil
	}
}
