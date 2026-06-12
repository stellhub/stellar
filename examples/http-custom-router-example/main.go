package main

import (
	"context"
	"log"
	"net/http"
	"strings"

	"github.com/stellhub/stellar"
	stellarhttp "github.com/stellhub/stellar/transport/http"
)

type customRouterStarter struct{}

type pingResponse struct {
	Message string `json:"message"`
}

type helloResponse struct {
	Message string `json:"message"`
}

type createItemRequest struct {
	Name string `json:"name"`
}

type createItemResponse struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func main() {
	if err := stellar.Start(stellar.WithStarter(newCustomRouterStarter())); err != nil {
		log.Fatal(err)
	}
}

func newCustomRouterStarter() *customRouterStarter {
	return &customRouterStarter{}
}

func (s *customRouterStarter) Name() string {
	return "http-custom-router-example"
}

func (s *customRouterStarter) Condition(stellar.StarterContext) bool {
	return true
}

func (s *customRouterStarter) Init(_ context.Context, app *stellar.App) error {
	api := app.HTTP().Group("/api/v1")

	api.GET("/ping", handlePing)
	api.GET("/hello", handleHello)
	stellarhttp.Handle(
		api,
		http.MethodPost,
		"/items",
		stellarhttp.JSONBinder[createItemRequest](),
		createItem,
		stellarhttp.JSONEncoder[createItemResponse],
	)

	return nil
}

func (s *customRouterStarter) Start(context.Context) error {
	return nil
}

func (s *customRouterStarter) Stop(context.Context) error {
	return nil
}

func handlePing(context.Context, *stellarhttp.Request) (*stellarhttp.Response, error) {
	return stellarhttp.JSON(http.StatusOK, pingResponse{Message: "pong"}), nil
}

func handleHello(_ context.Context, req *stellarhttp.Request) (*stellarhttp.Response, error) {
	name := strings.TrimSpace(req.Query.Get("name"))
	if name == "" {
		name = "stellar"
	}
	return stellarhttp.JSON(http.StatusOK, helloResponse{
		Message: "hello, " + name,
	}), nil
}

func createItem(_ context.Context, req *createItemRequest) (*createItemResponse, error) {
	return &createItemResponse{
		ID:   "item-001",
		Name: req.Name,
	}, nil
}
