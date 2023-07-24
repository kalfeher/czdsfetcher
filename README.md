# czdsfetcher
Fetch CZDS files. AWS Lambda function written in go.
Your **Lambda execution role** will contain 3 _permissions policies_:
1. _AWSLambdaBasicExecutionRole_: This is the standard Lambda execution role
1. [LambdaGetCZDAParams](#lambdagetczdaparams)
1. [S3WriteCZDA](#s3writeczda)


This function makes use of several AWS services in order to remain lightweight and flexible. Most of the configuration resides in **SSM: Parameter Store**. Zone files are downloaded to an **S3** bucket, for processing by other functions.
## Secrets and Parameters
The code assumes that an `AWS-Parameters-and-Secrets-Lambda-Extension` layer has been added to the lambda function. This will allow fetching secrets and parameters via **SSM: Parameter Store**. Parameters and secrets will be cached until execution is completed. No cache expiry is set.

Secrets are fetched by referencing AWS Secrets Manager secrets from Parameter Store parameters. In practice this involves placing a prefix to the secret-id before calling the Parameter Store cache.
The lambda execution role needs permission to access the secret as well as permission to read the secret prefix via the parameter store.

### LambdaGetCZDAParams
```json
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Sid": "AllowLambdaReadParamsAndSecrets",
            "Effect": "Allow",
            "Action": [
                "secretsmanager:GetSecretValue",
                "ssm:GetParameter"
            ],
            "Resource": [
                "arn:aws:ssm:yourregion:123456789:parameter/dev/czda/*",
                "arn:aws:secretsmanager:yourregion:123456789:secret:dev/czda-randomdigits",
                "arn:aws:ssm:yourregion:123456789:parameter/aws/reference/secretsmanager/dev/czda"
            ]
        }
    ]
}
```

### Parameters
Set the following parameters in **SSM: Parameter Store**.
|**Parameter Name** | **Value**|
| --- | --- |
|/dev/czda/authhost | https://account-api.icann.org/api/authenticate/|
|/dev/czda/server | https://czds-api.icann.org|
|/dev/czda/tlds | example, test _(replace with those TLDs for which you have download permissions)_|
|/dev/czda/user | myuser@example.example _(replace with your username)_|
|/dev/czda/bucket | mybucket _(replace with your bucket)_|

### Secrets
Put your CZDA password into **Secrets Manager**. Use the name _dev/czda_. Do not use a leading '/' in the secret name.

_Note that ICANN's systems appear to be unable to correctly parse any password that contains escape chars like `\`. I recommend limiting special chars in your CZDS account password strictly to those displayed on their password reset page._
## S3 Bucket
In order to allow the lambda function to upload files to **S3** create the following policy and ensure it is added to your Lambda execution role.
### S3WriteCZDA

```json
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Sid": "allowUpload",
            "Effect": "Allow",
            "Action": [
                "s3:AbortMultipartUpload",
                "s3:DeleteObject",
                "s3:ListMultipartUploadParts",
                "s3:PutObject",
                "s3:GetObject"
            ],
            "Resource": [
                "arn:aws:s3:::my-czda-bucket/*"
            ]
        },
        {
            "Effect": "Allow",
            "Action": [
                "s3:GetBucketLocation",
                "s3:ListBucket",
                "s3:ListBucketMultipartUploads"
            ],
            "Resource": [
                "arn:aws:s3:::my-czda-bucket"
            ]
        }
    ]
}
```
## Lambda Timeout
The default `Timeout` value for a Lambda is 15secs which will be far too short if you are downloading several TLDs. You can extend the timeout or alter the code to split TLDs between several executions.