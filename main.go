package main

import (
	"context"
	_ "embed"
	"errors"
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

func generatePolicy(principalId, effect, resource string) events.APIGatewayCustomAuthorizerResponse {
	authResponse := events.APIGatewayCustomAuthorizerResponse{PrincipalID: principalId}

	if effect != "" && resource != "" {
		authResponse.PolicyDocument = events.APIGatewayCustomAuthorizerPolicy{
			Version: "2012-10-17",
			Statement: []events.IAMPolicyStatement{
				{
					Action:   []string{"execute-api:Invoke"},
					Effect:   effect,
					Resource: []string{resource},
				},
			},
		}
	}

	return authResponse
}

func handleRequest(ctx context.Context, request events.APIGatewayCustomAuthorizerRequest) (events.APIGatewayCustomAuthorizerResponse, error) {
	slog.Info("validating token")

	tokenString := request.AuthorizationToken

	if tokenString == "" {
		slog.Error("no token found in request")
		return events.APIGatewayCustomAuthorizerResponse{}, errors.New("Unauthorized") // Return a 401 Unauthorized response
	}

	isValid, err := validateToken(tokenString)
	if err != nil {
		slog.Error("error validating token: %v", err)
		return events.APIGatewayCustomAuthorizerResponse{}, errors.New("Unauthorized") // Return a 401 Unauthorized response
	}

	if isValid {
		slog.Info("token is valid")
		return generatePolicy("user", "Allow", request.MethodArn), nil
	}

	slog.Error("token is not valid")
	return events.APIGatewayCustomAuthorizerResponse{}, errors.New("Unauthorized") // Return a 401 Unauthorized response
}

func main() {
	lambda.Start(handleRequest)
}
