package fetcher

import (
	"context"
	"czdsfetch/configs"
	"czdsfetch/internal/czds"
	"czdsfetch/public/paramstore"
	"fmt"
	"strings"

	"github.com/aws/aws-lambda-go/events"
)

func HandleRequest(ctx context.Context, s3Event events.S3Event) {

	authtoken := czds.GetAuthToken(configs.CZDSserver, configs.CZDSuser, configs.CZDSpassword)
	println(authtoken)
	println("************************")
	myuser := paramstore.GetParameterStoreValue(configs.CZDSuser, false)
	println(myuser)
	println("########################")
	mypass := paramstore.GetParameterStoreValue(configs.CZDSpassword, true)
	println(mypass)
	println("########################")
	tlds := paramstore.GetParameterStoreValue(configs.CZDStlds, false)
	println(tlds)
	mytlds := strings.Split(tlds, ",")
	fmt.Println("these are the tlds:", mytlds) // regular println invocation won't display

}
