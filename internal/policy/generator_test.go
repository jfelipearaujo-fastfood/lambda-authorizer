package policy

import (
	"reflect"
	"testing"

	"github.com/aws/aws-lambda-go/events"
)

func Test_generator(t *testing.T) {
	type args struct {
		principalId string
		effect      string
	}
	tests := []struct {
		name string
		args args
		want events.APIGatewayCustomAuthorizerResponse
	}{
		{
			name: "Generate Allow Policy",
			args: args{
				principalId: "user",
				effect:      "allow",
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
		},
		{
			name: "Generate Deny Policy",
			args: args{
				principalId: "user",
				effect:      "deny",
			},
			want: events.APIGatewayCustomAuthorizerResponse{
				PrincipalID: "user",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := generator(tt.args.principalId, tt.args.effect); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("generator() = %v, want %v", got, tt.want)
			}
		})
	}
}
