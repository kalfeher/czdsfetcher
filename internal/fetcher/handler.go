package fetcher

import (
	"context"
	"czdsfetch/configs"
	"czdsfetch/internal/czds"

	"github.com/aws/aws-lambda-go/events"
)

func HandleRequest(ctx context.Context, s3Event events.S3Event) {
	authtoken := czds.GetAuthToken(configs.CZDSserver, configs.CZDSuser, configs.CZDSpassword)
	println(authtoken)
}
