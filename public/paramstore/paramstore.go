package paramstore

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"
)

var (
	SessionToken            = os.Getenv("AWS_SESSION_TOKEN")
	ParamSecretsExtHTTPPort = "2773"
	ParamStoreURL           = "http://localhost:" + ParamSecretsExtHTTPPort + "/systemsmanager/parameters/get?"
	LayerHeader             = "X-Aws-Parameters-Secrets-Token"

	AWSSecretsPrefix = "/aws/reference/secretsmanager/"
	Decrypt          = "withDecryption"
)

// structs to save unmarshalled json to
type AWSParamResponse struct {
	Parameter      *AWSParam         `json:"Parameter"`
	ResultMetadata *AWSParamMetadata `json:"ResultMetadata"`
}
type AWSParam struct {
	ARN              string
	DataType         string
	LastModifiedDate time.Time
	Name             string
	Selector         string
	SourceResult     string
	Type             string
	Value            string
	Version          int64
}
type AWSParamMetadata struct {
}

func fatal(err error) string {
	log.Fatal(err)
	return ""
}

// Function to get parameter store values from lambda layer
func GetParameterStoreValue(key string, isSecret bool, client *http.Client) string {

	params := url.Values{}

	if isSecret {
		params.Add(Decrypt, "true")
		key = AWSSecretsPrefix + key
	}
	params.Add("name", key)

	req, err := http.NewRequest("GET", ParamStoreURL+params.Encode(), nil)
	if err != nil {
		fatal(err)
	}
	req.Header.Add(LayerHeader, SessionToken)
	resp, err := client.Do(req)
	if err != nil {
		fatal(err)
	} // return err
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fatal(err)
	} // return err
	var data AWSParamResponse
	err = json.Unmarshal(body, &data)
	if err != nil {
		fatal(err)
	} // return err
	return data.Parameter.Value
}
