package main

import (
    "bytes"
    "context"
    "encoding/json"
    "os"

    "github.com/aws/aws-lambda-go/events"
    "github.com/aws/aws-lambda-go/lambda"
)

// Response is of type APIGatewayProxyResponse since we're leveraging the
// AWS Lambda Proxy Request functionality (default behavior)
//
// https://serverless.com/framework/docs/providers/aws/events/apigateway/#lambda-proxy-integration
type Response events.APIGatewayProxyResponse

// Handler is our lambda handler invoked by the `lambda.Start` function call
func Handler(ctx context.Context) (Response, error) {
    var buf bytes.Buffer

    body, err := json.Marshal(map[string]interface{}{
        "application": "hagowebapp",
        "body": ctx, // TODO: get actual body
        "environment": os.Getenv("STAGE"),
        "method": ctx, // TODO: get the actual method
        "message": "bar function executed successfully!", // TODO: get function name
        "provisioner": "Created with Serverless v1.0",
        "runtime": "Go (go1.x)",
        "version": "0.0.0",
    })
    if err != nil {
        return Response{StatusCode: 404}, err
    }
    json.HTMLEscape(&buf, body)

    resp := Response{
        StatusCode: 200,
        IsBase64Encoded: false,
        Body: buf.String(),
        Headers: map[string]string{
            "Content-Type": "application/json",
        },
    }

    return resp, nil
}

func main() {
    lambda.Start(Handler)
}
