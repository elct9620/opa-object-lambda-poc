package main

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go-v2/aws"
)

func onListObjects(ctx context.Context, event events.S3ObjectLambdaEvent) (ObjectLambdaResponse, error) {
	policyQuery := GetPolicy(ctx)
	isAllowed, err := policyQuery.Allow(ctx, event.UserIdentity.Type, event.UserIdentity.ARN)
	if err != nil {
		return internalError(err), err
	}

	if !isAllowed {
		return accessDenied(fmt.Sprintf("UserIdentity Type %s is not allowed", event.UserIdentity.Type)), nil
	}

	res, err := http.Get(event.ListObjectsContext.InputS3URL)
	if err != nil {
		return internalError(err), err
	}

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return internalError(err), err
	}

	return ObjectLambdaResponse{
		StatusCode:    200,
		ListResultXml: aws.String(string((resBody))),
	}, nil
}

func onListObjectsV2(ctx context.Context, event events.S3ObjectLambdaEvent) (ObjectLambdaResponse, error) {
	policyQuery := GetPolicy(ctx)
	isAllowed, err := policyQuery.Allow(ctx, event.UserIdentity.Type, event.UserIdentity.ARN)
	if err != nil {
		return internalError(err), err
	}

	if !isAllowed {
		return accessDenied(fmt.Sprintf("UserIdentity Type %s is not allowed", event.UserIdentity.Type)), nil
	}

	res, err := http.Get(event.ListObjectsV2Context.InputS3URL)
	if err != nil {
		return internalError(err), err
	}

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return internalError(err), err
	}

	return ObjectLambdaResponse{
		StatusCode:    200,
		ListResultXml: aws.String(string((resBody))),
	}, nil
}
