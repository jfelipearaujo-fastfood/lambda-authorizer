package policy

import "github.com/aws/aws-lambda-go/events"

func GenerateAllowPolicy() events.APIGatewayCustomAuthorizerResponse {
	return generator("user", "allow")
}

func generator(principalId, effect string) events.APIGatewayCustomAuthorizerResponse {
	authResponse := events.APIGatewayCustomAuthorizerResponse{
		PrincipalID: principalId,
	}

	if effect == "allow" {
		authResponse.PolicyDocument = events.APIGatewayCustomAuthorizerPolicy{
			Version: "2012-10-17",
			Statement: []events.IAMPolicyStatement{
				{
					Action:   []string{"execute-api:Invoke"},
					Effect:   effect,
					Resource: []string{"*"},
				},
			},
		}
	}

	return authResponse
}
