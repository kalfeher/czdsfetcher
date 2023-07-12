# czdsfetcher
Fetch CZDS files
## Secrets and Parameters
The code assumes that an `AWS-Parameters-and-Secrets-Lambda-Extension` layer has been added to the lambda function. This will allow fetching secrets and parameters via SSM:Parameter Store. Parameters and secrets will be cached until execution is completed. No cache expiry is set.

Secrets are fetched by referencing AWS Secrets Manager secrets from Parameter Store parameters. In practice this involves placing a prefix to the secret-id before calling the Parameter Store cache.
The lambda execution role needs permission to access the secret as well as permission to read the secret prefix via the parameter store.
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