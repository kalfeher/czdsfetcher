package fetcher

import (
	"context"
	"czdsfetch/configs"
	"czdsfetch/internal/czds"
	"czdsfetch/public/paramstore"
	"czdsfetch/public/s3"
	"fmt"
	"net/http"
	"os"
	"sync"

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
	bucket := paramstore.GetParameterStoreValue(configs.Bucket, false, client)

	authtoken := czds.GetAuthToken(myAuthHost, myuser, mypass, client)

	downloadlinks := czds.GetDownloadLinks(server, authtoken, client)

	DownloadFiles(downloadlinks, authtoken, client, bucket)
}

// A function to iterate over downloadlinks and download the files
func DownloadFiles(downloadlinks []string, authtoken string, client *http.Client, bucket string) {
	uploader := s3.Uploader()
	wg := sync.WaitGroup{}
	fetchers := make(chan struct{}, configs.Fetchers)

	for _, link := range downloadlinks {
		wg.Add(1)
		fetchers <- struct{}{} // Acquire a fetcher
		go func(link string) {

			zoneFile := czds.GetZoneFile(link, configs.LocalDirectory, authtoken, client)
			res := s3.UploadToBucket(uploader, bucket, zoneFile)
			fmt.Print(res)
			os.Remove(zoneFile)
			wg.Done()
			<-fetchers // Release the fetcher
		}(link)

	}
}
