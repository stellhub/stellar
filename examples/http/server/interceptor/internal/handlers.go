package internal

import (
	"context"
	"fmt"
	"net/http"

	stellarhttp "github.com/stellhub/stellar/transport/http"
)

type pingResponse struct {
	Message     string `json:"message"`
	Interceptor string `json:"interceptor"`
}

type echoRequest struct {
	Message string `json:"message"`
}

type echoResponse struct {
	Message string `json:"message"`
}

func handlePing(context.Context, *stellarhttp.Request) (*stellarhttp.Response, error) {
	return stellarhttp.JSON(http.StatusOK, pingResponse{
		Message:     "pong",
		Interceptor: "http-server",
	}), nil
}

func echo(_ context.Context, req *echoRequest) (*echoResponse, error) {
	return &echoResponse{
		Message: fmt.Sprintf("echo: %s", req.Message),
	}, nil
}
