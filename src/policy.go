package main

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws/arn"
	"github.com/open-policy-agent/opa/rego"
	"github.com/rs/zerolog/log"
)

type policyCtxKey struct{}

const accessPolicy = `
package me.aotoki.s3.object
import rego.v1

default allow := false
allow if {
	input.type == "IAMUser"
}

# Example: Allow role name contains github
#
# allow if {
# 	input.type == "AssumedRole"
# 	contains(input.identity, "github")
# }
`

type Policy interface {
	Allow(ctx context.Context, identityType, identityArn string) (bool, error)
}

func GetPolicy(ctx context.Context) Policy {
	return ctx.Value(policyCtxKey{}).(Policy)
}

type OpenPolicy struct {
	query rego.PreparedEvalQuery
}

func NewOpenPolicy() (*OpenPolicy, error) {
	query, err := rego.New(
		rego.Query("data.me.aotoki.s3.object.allow"),
		rego.Module("me.aotoki.s3", accessPolicy),
	).PrepareForEval(context.Background())
	if err != nil {
		return nil, err
	}

	return &OpenPolicy{query: query}, nil
}

func (p *OpenPolicy) Allow(ctx context.Context, identityType, identityArn string) (bool, error) {
	arn, err := arn.Parse(identityArn)
	if err != nil {
		return false, err
	}

	input := map[string]interface{}{
		"type":     identityType,
		"identity": arn.Resource,
	}

	results, err := p.query.Eval(ctx, rego.EvalInput(input))
	if err != nil {
		return false, err
	}

	log.Info().Bool("allowed", results.Allowed()).Str("identityType", identityType).Str("identityArn", identityArn).Msg("Auditing access...")

	return results.Allowed(), nil
}
