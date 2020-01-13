package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-lambda-go/lambdacontext"
)

// Response is of type APIGatewayProxyResponse since we're leveraging the
// AWS Lambda Proxy Request functionality (default behavior)
//
// https://serverless.com/framework/docs/providers/aws/events/apigateway/#lambda-proxy-integration
type Response events.APIGatewayProxyResponse

// Handler function Using AWS Lambda Proxy Request (invoked by the `lambda.Start`)
func Handler(ctx context.Context, request events.APIGatewayProxyRequest) (Response, error) {
	if request.HTTPMethod != "GET" {
		return Response{StatusCode: 405}, nil
	}

	var buf bytes.Buffer

	lambda := fmt.Sprintf("%s (%s)", lambdacontext.FunctionName, lambdacontext.FunctionVersion)
	fn := lambdacontext.FunctionName

	body, err := json.Marshal(map[string]interface{}{
		"application": os.Getenv("APPLICATION"),
		"body":        request.Body,
		"environment": os.Getenv("ENVIRONMENT"),
		"host":        request.Headers["Host"],
		"lambda":      lambda,
		"message":     fmt.Sprintf("%s function executed successfully!", fn),
		"method":      request.HTTPMethod,
		"provisioner": "Created with Serverless v1.0",
		"runtime":     "Go (go1.x)",
		"version":     "0.0.3",
	})
	if err != nil {
		return Response{StatusCode: 404}, err
	}
	json.HTMLEscape(&buf, body)

	resp := Response{
		StatusCode:      200,
		IsBase64Encoded: false,
		Body:            buf.String(),
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	}

	return resp, nil
}

func main() {
	lambda.Start(Handler)
}
