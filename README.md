# czdsfetcher
Fetch CZDS files
## Secrets and Parameters
The code assumes that an `AWS-Parameters-and-Secrets-Lambda-Extension` layer has been added to the lambda function. This will allow fetching secrets and parameters via **SSM:Parameter Store**. Parameters and secrets will be cached until execution is completed. No cache expiry is set.

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

### Parameters
Set the following parameters in **SSM: Parameter Store**. You should 
|**Parameter Name** | **Value**|
| --- | --- |
|/dev/czda/authhost | https://account-api.icann.org/api/authenticate/|
|/dev/czda/server | https://czds-api.icann.org|
|/dev/czda/tlds | example, test _(replace with those TLDs for which you have download permissions)_|
|/dev/czda/user | myuser@example.example _(replace with your username)_|

### Secrets
Put your CZDA password into **Secrets Manager**. Use the name _dev/czda_. Do not use a leading '/' in the secret name.
_Note that ICANN's systems appear to be unable to correctly parse any password that contains escape chars like '\'. I recommend limiting special chars in your CZDS account password strictly to those displayed on their password reset page._

