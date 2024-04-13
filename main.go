package main

import (
	_ "embed"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/jsfelipearaujo/lambda-authorizer/src/handler"
)

func main() {
	lambda.Start(handler.HandleRequest)
}
