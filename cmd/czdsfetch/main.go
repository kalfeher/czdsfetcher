package main

import (
	"czdsfetch/internal/fetcher"

	"github.com/aws/aws-lambda-go/lambda"
)

func main() {
	lambda.Start(fetcher.HandleRequest)
}
