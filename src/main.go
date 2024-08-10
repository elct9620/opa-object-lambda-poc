package main

import (
	"context"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/rs/zerolog/log"
)

type ObjectLambdaResponse struct {
	// Shared Field
	StatusCode int `json:"statusCode"`
	// Error Status
	ErrorCode    *string `json:"errorCode,omitempty"`
	ErrorMessage *string `json:"errorMessage,omitempty"`
	// HeadObject
	Headers map[string]string `json:"headers,omitempty"`
	// ListObjects
	ListResultXml *string `json:"listResultXml,omitempty"`
}

func handleDataAccess(ctx context.Context, event events.S3ObjectLambdaEvent) (ObjectLambdaResponse, error) {
	log.Debug().Interface("event", event).Msg("Received event")

	if event.ListObjectsContext != nil {
		return onListObjects(ctx, event)
	}

	if event.ListObjectsV2Context != nil {
		return onListObjectsV2(ctx, event)
	}

	if event.GetObjectContext != nil {
		return onGetObject(ctx, event)
	}

	log.Warn().Msg("Unsupported operation")

	return ObjectLambdaResponse{
		StatusCode:   500,
		ErrorCode:    aws.String("InvalidRequest"),
		ErrorMessage: aws.String("Unsupported operation"),
	}, nil
}

func main() {
	ctx := context.Background()

	policy, err := NewOpenPolicy()
	if err != nil {
		panic(err)
	}

	lambda.StartWithOptions(
		handleDataAccess,
		lambda.WithContext(context.WithValue(ctx, policyCtxKey{}, policy)),
	)
}
