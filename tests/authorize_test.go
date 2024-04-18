package tests

import (
	"context"
	"errors"
	"flag"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/cucumber/godog"
	"github.com/cucumber/godog/colors"
	"github.com/golang-jwt/jwt"
	"github.com/jfelipearaujo-org/lambda-authorizer/src/handler"
)

var opts = godog.Options{
	Format:      "pretty",
	Paths:       []string{"features"},
	Output:      colors.Colored(os.Stdout),
	Concurrency: 4,
}

func init() {
	godog.BindFlags("godog.", flag.CommandLine, &opts)
}

func TestFeatures(t *testing.T) {
	t.Setenv("SIGN_KEY", "key")

	o := opts
	o.TestingT = t

	status := godog.TestSuite{
		TestSuiteInitializer: InitializeTestSuite,
		ScenarioInitializer:  InitializeScenario,
		Options:              &o,
	}.Run()

	if status == 2 {
		t.SkipNow()
	}

	if status != 0 {
		t.Fatalf("zero status code expected, %d received", status)
	}
}

func generateToken(t *testing.T,
	method jwt.SigningMethod,
	signingKey string,
	claims jwt.MapClaims,
) string {
	token := jwt.NewWithClaims(method, claims)

	tokenString, err := token.SignedString([]byte(signingKey))
	if err != nil {
		t.Fatal(err)
	}

	return tokenString
}

type tokenCtxKey struct{}
type tokenResponseCtxKey struct{}

func tokenToContext(ctx context.Context, token string) context.Context {
	return context.WithValue(ctx, tokenCtxKey{}, token)
}

func tokenFromContext(ctx context.Context) string {
	token, _ := ctx.Value(tokenCtxKey{}).(string)
	return token
}

func responseToContext(ctx context.Context, response events.APIGatewayCustomAuthorizerResponse) context.Context {
	return context.WithValue(ctx, tokenResponseCtxKey{}, response)
}

func responseFromContext(ctx context.Context) events.APIGatewayCustomAuthorizerResponse {
	response, _ := ctx.Value(tokenResponseCtxKey{}).(events.APIGatewayCustomAuthorizerResponse)
	return response
}

// Steps

func iHaveAValidToken(ctx context.Context) (context.Context, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": "123e4567-e89b-12d3-a456-426614174000",
		"exp": time.Now().Add(time.Hour).Unix(),
	})

	tokenString, err := token.SignedString([]byte("key"))
	if err != nil {
		return ctx, err
	}

	return tokenToContext(ctx, tokenString), nil
}

func iAuthorizeTheRequest(ctx context.Context) (context.Context, error) {
	token := tokenFromContext(ctx)

	resp, err := handler.HandleRequest(ctx, events.APIGatewayCustomAuthorizerRequestTypeRequest{
		Headers: map[string]string{
			"Authorization": token,
		},
	})

	return responseToContext(ctx, resp), err
}

func theRequestShouldBeAuthorized(ctx context.Context) (context.Context, error) {
	resp := responseFromContext(ctx)

	expected := events.APIGatewayCustomAuthorizerResponse{
		PrincipalID: "user",
		PolicyDocument: events.APIGatewayCustomAuthorizerPolicy{
			Version: "2012-10-17",
			Statement: []events.IAMPolicyStatement{
				{
					Action:   []string{"execute-api:Invoke"},
					Effect:   "allow",
					Resource: []string{"*"},
				},
			},
		},
	}

	if !reflect.DeepEqual(resp, expected) {
		return ctx, errors.New("response does not match expected")
	}

	return ctx, nil
}

func InitializeTestSuite(ctx *godog.TestSuiteContext) {
	ctx.BeforeSuite(func() {
		// do nothing for now
	})
}

func InitializeScenario(ctx *godog.ScenarioContext) {
	// ctx.Before(func(ctx context.Context, sc *godog.Scenario) (context.Context, error) {
	// 	return tokenToContext(ctx, "123"), nil
	// })

	ctx.Step(`^I have a valid token$`, iHaveAValidToken)
	ctx.Step(`^I authorize the request$`, iAuthorizeTheRequest)
	ctx.Step(`^the request should be authorized$`, theRequestShouldBeAuthorized)
}
