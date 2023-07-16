package fetcher

import (
	"context"
	"czdsfetch/configs"
	"czdsfetch/internal/czds"
	"czdsfetch/public/paramstore"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
)

var (
	client = &http.Client{}
)

func HandleRequest(ctx context.Context, s3Event events.S3Event) {
	myuser := paramstore.GetParameterStoreValue(configs.CZDSuser, false, client)
	mypass := paramstore.GetParameterStoreValue(configs.CZDSpassword, true, client)
	tlds := paramstore.GetParameterStoreValue(configs.CZDStlds, false, client)
	println(tlds)
	// mytlds := strings.Split(tlds, ",")
	// fmt.Println("these are the tlds:", mytlds) // regular println invocation won't display
	myAuthHost := paramstore.GetParameterStoreValue(configs.CZDSAuthHost, false, client)
	println("************************")
	authtoken := czds.GetAuthToken(myAuthHost, myuser, mypass, client)
	println("########################")
	println("full response:", authtoken)
	println("************************")
}
