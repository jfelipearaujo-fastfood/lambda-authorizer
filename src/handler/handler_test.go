package handler

import (
	"context"
	_ "embed"
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/golang-jwt/jwt"
)

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

func Test_getTokenFromRequestHeaders(t *testing.T) {
	type args struct {
		request events.APIGatewayCustomAuthorizerRequestTypeRequest
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Header [A]uthorization exists",
			args: args{
				request: events.APIGatewayCustomAuthorizerRequestTypeRequest{
					Headers: map[string]string{
						"Authorization": "token",
					},
				},
			},
			want: "token",
		},
		{
			name: "Header [a]uthorization exists",
			args: args{
				request: events.APIGatewayCustomAuthorizerRequestTypeRequest{
					Headers: map[string]string{
						"authorization": "token",
					},
				},
			},
			want: "token",
		},
		{
			name: "No Header exists",
			args: args{
				request: events.APIGatewayCustomAuthorizerRequestTypeRequest{
					Headers: map[string]string{
						"abc": "123",
					},
				},
			},
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getTokenFromRequestHeaders(tt.args.request); got != tt.want {
				t.Errorf("getTokenFromRequestHeaders() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_HandleRequest(t *testing.T) {
	type args struct {
		ctx     context.Context
		request events.APIGatewayCustomAuthorizerRequestTypeRequest
	}
	tests := []struct {
		name    string
		args    args
		want    events.APIGatewayCustomAuthorizerResponse
		wantErr error
	}{
		{
			name: "Token valid and policy generated",
			args: args{
				ctx: context.Background(),
				request: events.APIGatewayCustomAuthorizerRequestTypeRequest{
					Headers: map[string]string{
						"Authorization": generateToken(t, jwt.SigningMethodHS256, "key", jwt.MapClaims{
							"sub": "123e4567-e89b-12d3-a456-426614174000",
							"exp": time.Now().Add(time.Hour).Unix(),
						}),
					},
				},
			},
			want: events.APIGatewayCustomAuthorizerResponse{
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
			},
			wantErr: nil,
		},
		{
			name: "No token found in request",
			args: args{
				ctx: context.Background(),
				request: events.APIGatewayCustomAuthorizerRequestTypeRequest{
					Headers: map[string]string{},
				},
			},
			want:    events.APIGatewayCustomAuthorizerResponse{},
			wantErr: ErrUnauthorized,
		},
		{
			name: "Token invalid",
			args: args{
				ctx: context.Background(),
				request: events.APIGatewayCustomAuthorizerRequestTypeRequest{
					Headers: map[string]string{
						"Authorization": "abc",
					},
				},
			},
			want:    events.APIGatewayCustomAuthorizerResponse{},
			wantErr: ErrUnauthorized,
		},
		{
			name: "Token expired",
			args: args{
				ctx: context.Background(),
				request: events.APIGatewayCustomAuthorizerRequestTypeRequest{
					Headers: map[string]string{
						"Authorization": generateToken(t, jwt.SigningMethodHS256, "key", jwt.MapClaims{
							"sub": "123e4567-e89b-12d3-a456-426614174000",
							"exp": time.Now().Add(time.Hour * -1).Unix(),
						}),
					},
				},
			},
			want:    events.APIGatewayCustomAuthorizerResponse{},
			wantErr: ErrUnauthorized,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Setenv("SIGN_KEY", "key")

			got, err := HandleRequest(tt.args.ctx, tt.args.request)
			if err != nil && tt.wantErr == nil || !errors.Is(err, tt.wantErr) {
				t.Errorf("handleRequest() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("handleRequest() = %v, want %v", got, tt.want)
			}
		})
	}
}
