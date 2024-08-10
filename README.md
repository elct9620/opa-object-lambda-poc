OpenPolicyAgent Object Lambda PoC
===

Building a custom data access control solution using AWS Lambda and Amazon S3 Object Lambda.

For future work, we can put audit events to CloudTrail to monitor the access to the object lambda.
> https://docs.aws.amazon.com/awscloudtraildata/latest/APIReference/API_PutAuditEvents.html

## Requirements

* Go 1.2x
* AWS CLI
* AWS SAM CLI

## Usage

Deploy to your AWS account using the following command:

```bash
AWS_PROFILE=<your-profile> make deploy
```

Find `S3ObjectLambdaAccessPointAlias` in the output of the command above, and use it to access the object lambda.

## Testing

To verify only IAM user is allowed to access the object lambda, you can create a AssumeRole profile in your `~/.aws/config` file:

```bash
[profile <your-profile>]
role_arn = arn:aws:iam::<account-id>:role/<role-name>
source_profile = <source-profile>
```

### ListObjects

```bash
AWS_PROFILE=<your-profile> aws s3api list-objects --bucket <bucket-alias>
```

### ListObjectsV2

```bash
AWS_PROFILE=<your-profile> aws s3api list-objects-v2 --bucket <bucket-alias>
```

For higher-level APIs, you can use the following commands:
```bash
AWS_PROFILE=<your-profile> aws s3 ls s3://<bucket-alias>
```

### GetObject

```bash
AWS_PROFILE=<your-profile> aws s3api get-object --bucket <bucket-alias> --key <object-key> <output-file>
```

## References

* https://aws.amazon.com/tw/blogs/storage/managing-access-to-your-amazon-s3-objects-with-a-custom-authorizer/
* https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-s3objectlambda-accesspoint.html
* https://docs.aws.amazon.com/AmazonS3/latest/userguide/olap-writing-lambda.html
