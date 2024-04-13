package handler

import (
	"context"
	"errors"
	"os"

	"log/slog"

	"github.com/aws/aws-lambda-go/events"
	"github.com/jsfelipearaujo/lambda-authorizer/src/policy"
	"github.com/jsfelipearaujo/lambda-authorizer/src/token"
)

func init() {
	opts := &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}

	handler := slog.NewTextHandler(os.Stdout, opts)
	log := slog.New(handler)
	slog.SetDefault(log)
}

var (
	ErrUnauthorized = errors.New("Unauthorized")
)

func getTokenFromRequestHeaders(request events.APIGatewayCustomAuthorizerRequestTypeRequest) string {
	headers := []string{"authorization", "Authorization"}

	for _, header := range headers {
		tokenString := request.Headers[header]
		if tokenString != "" {
			return tokenString
		}
	}

	return ""
}

func HandleRequest(ctx context.Context, request events.APIGatewayCustomAuthorizerRequestTypeRequest) (events.APIGatewayCustomAuthorizerResponse, error) {
	tokenString := getTokenFromRequestHeaders(request)

	if tokenString == "" {
		slog.Error("no token found in request")
		return events.APIGatewayCustomAuthorizerResponse{}, ErrUnauthorized
	}

	isValid, err := token.Validator(tokenString)
	if err != nil {
		slog.Error("error validating token: %v", err)
		return events.APIGatewayCustomAuthorizerResponse{}, ErrUnauthorized
	}

	if isValid {
		return policy.GenerateAllowPolicy(), nil
	}

	slog.Error("token is not valid")
	return events.APIGatewayCustomAuthorizerResponse{}, ErrUnauthorized
}
