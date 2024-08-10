package main

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/rs/zerolog/log"
)

func onGetObject(ctx context.Context, event events.S3ObjectLambdaEvent) (ObjectLambdaResponse, error) {
	uri, err := url.Parse(event.UserRequest.URL)
	if err != nil {
		return internalError(err), err
	}
	log.Info().Str("key", uri.Path).Msg("Received GetObject request")

	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return internalError(err), err
	}

	s3api := s3.NewFromConfig(cfg)

	policyQuery := GetPolicy(ctx)
	isAllowed, err := policyQuery.Allow(ctx, event.UserIdentity.Type, event.UserIdentity.ARN)
	if err != nil {
		return internalError(err), err
	}

	if !isAllowed {
		_, err := s3api.WriteGetObjectResponse(ctx, &s3.WriteGetObjectResponseInput{
			RequestRoute: aws.String(event.GetObjectContext.OutputRoute),
			RequestToken: aws.String(event.GetObjectContext.OutputToken),
			StatusCode:   aws.Int32(403),
			ErrorCode:    aws.String("AccessDenied"),
			ErrorMessage: aws.String(fmt.Sprintf("UserIdentity Type %s is not allowed", event.UserIdentity.Type)),
		})

		if err != nil {
			return internalError(err), err
		}

		return ObjectLambdaResponse{
			StatusCode: 200,
		}, nil
	}

	res, err := http.Get(event.GetObjectContext.InputS3URL)
	if err != nil {
		return internalError(err), err
	}

	_, err = s3api.WriteGetObjectResponse(ctx, &s3.WriteGetObjectResponseInput{
		RequestRoute:  aws.String(event.GetObjectContext.OutputRoute),
		RequestToken:  aws.String(event.GetObjectContext.OutputToken),
		ContentType:   aws.String(res.Header.Get("Content-Type")),
		ContentLength: aws.Int64(res.ContentLength),
		StatusCode:    aws.Int32(int32(res.StatusCode)),
		Body:          res.Body,
	})

	if err != nil {
		return internalError(err), err
	}

	return ObjectLambdaResponse{
		StatusCode: 200,
	}, nil
}
