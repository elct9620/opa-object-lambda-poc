package main

import "github.com/aws/aws-sdk-go-v2/aws"

func internalError(err error) ObjectLambdaResponse {
	return ObjectLambdaResponse{
		StatusCode:   500,
		ErrorMessage: aws.String(err.Error()),
	}
}

func accessDenied(reason string) ObjectLambdaResponse {
	return ObjectLambdaResponse{
		StatusCode:   403,
		ErrorCode:    aws.String("AccessDenied"),
		ErrorMessage: aws.String(reason),
	}
}
