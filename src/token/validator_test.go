package token

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt"
)

func generateToken(t *testing.T,
	method jwt.SigningMethod,
	signingKey string,
	claims jwt.MapClaims,
) (string, error) {
	tokenData := jwt.NewWithClaims(method, claims)

	return tokenData.SignedString([]byte(signingKey))
}

func TestValidator(t *testing.T) {
	type args struct {
		signingKey string
		method     jwt.SigningMethod
		claims     jwt.MapClaims
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "Valid Token",
			args: args{
				signingKey: "key",
				method:     jwt.SigningMethodHS256,
				claims: jwt.MapClaims{
					"sub": "123e4567-e89b-12d3-a456-426614174000",
					"exp": time.Now().Add(time.Hour).Unix(),
				},
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "Invalid Token (wrong method)",
			args: args{
				signingKey: "key",
				method:     jwt.SigningMethodHS384,
				claims: jwt.MapClaims{
					"sub": "123e4567-e89b-12d3-a456-426614174000",
					"exp": time.Now().Add(time.Hour).Unix(),
				},
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "Invalid Token (missing sub)",
			args: args{
				signingKey: "key",
				method:     jwt.SigningMethodHS256,
				claims: jwt.MapClaims{
					"exp": time.Now().Add(time.Hour).Unix(),
				},
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "Invalid sub",
			args: args{
				signingKey: "key",
				method:     jwt.SigningMethodHS256,
				claims: jwt.MapClaims{
					"sub": "123456",
					"exp": time.Now().Add(time.Hour).Unix(),
				},
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "Expired Token",
			args: args{
				signingKey: "key",
				method:     jwt.SigningMethodHS384,
				claims: jwt.MapClaims{
					"sub": "123e4567-e89b-12d3-a456-426614174000",
					"exp": time.Now().Add(time.Hour * -1).Unix(),
				},
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "Missing Claims",
			args: args{
				signingKey: "key",
				method:     jwt.SigningMethodHS384,
				claims:     jwt.MapClaims{},
			},
			want:    false,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Setenv("SIGN_KEY", tt.args.signingKey)

			tokenString, err := generateToken(t, tt.args.method, tt.args.signingKey, tt.args.claims)
			if err != nil {
				t.Errorf("generateToken() error = %v", err)
				return
			}

			got, err := Validator(tokenString)
			if (err != nil) != tt.wantErr {
				t.Errorf("Validator() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Validator() = %v, want %v", got, tt.want)
			}
		})
	}
}
