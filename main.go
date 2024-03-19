package main

import (
	"context"
	_ "embed"
	"fmt"
	"os"

	"log/slog"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
)

var (
	signingKey = []byte(os.Getenv("SIGN_KEY"))
)

func init() {
	opts := &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}

	handler := slog.NewTextHandler(os.Stdout, opts)
	log := slog.New(handler)
	slog.SetDefault(log)
}

func validateToken(tokenString string) (bool, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return signingKey, nil
	})
	if err != nil {
		return false, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		userId := claims["sub"]
		if userId == nil {
			return false, fmt.Errorf("user id not found in token")
		}

		if _, err := uuid.Parse(userId.(string)); err != nil {
			return false, fmt.Errorf("invalid user id '%v' in token: %w", userId, err)
		}

		return true, nil
	}

	return false, fmt.Errorf("invalid token")
}

func handleRequest(ctx context.Context, request events.APIGatewayCustomAuthorizerRequestTypeRequest) (events.APIGatewayV2CustomAuthorizerSimpleResponse, error) {
	slog.Info("validating token")

	tokenString := request.Headers["authorization"]

	if tokenString == "" {
		slog.Error("no token found in request")
		return events.APIGatewayV2CustomAuthorizerSimpleResponse{
			IsAuthorized: false,
		}, nil
	}

	isValid, err := validateToken(tokenString)
	if err != nil {
		slog.Error("error validating token: %v", err)
		return events.APIGatewayV2CustomAuthorizerSimpleResponse{
			IsAuthorized: false,
		}, nil
	}

	if isValid {
		slog.Info("token is valid")
		return events.APIGatewayV2CustomAuthorizerSimpleResponse{
			IsAuthorized: true,
		}, nil
	}

	slog.Error("token is not valid")
	return events.APIGatewayV2CustomAuthorizerSimpleResponse{
		IsAuthorized: false,
	}, nil
}

func main() {
	lambda.Start(handleRequest)
}
