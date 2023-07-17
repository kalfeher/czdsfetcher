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
	//tlds := paramstore.GetParameterStoreValue(configs.CZDStlds, false, client)

	myAuthHost := paramstore.GetParameterStoreValue(configs.CZDSAuthHost, false, client)
	server := paramstore.GetParameterStoreValue(configs.CZDSserver, false, client)
	//bucket := paramstore.GetParameterStoreValue(configs.Bucket, false, client)
	println("************************")
	authtoken := czds.GetAuthToken(myAuthHost, myuser, mypass, client)
	println("########################")
	downloadlinks := czds.GetDownloadLinks(server, authtoken, client)
	for i, link := range downloadlinks {
		println(i, ":", link)

	}
}
