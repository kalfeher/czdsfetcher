package czds

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

var (
	successMessage = "Authentication Successful"
	listDownloads  = "/czds/downloads/links"
)

type Download struct {
	url string
}

// password struct
type AccessToken struct {
	AccessToken string `json:"accessToken"`
	Message     string `json:"message"`
}

// get download links from czds host
func GetDownloadLinks(server string, token string, client *http.Client) []string {
	req, err := http.NewRequest("GET", server+listDownloads, nil)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Add("Authorization", "Bearer "+token)
	req.Header.Add("Accept", `application/json`)
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	downloadList := strings.Split(strings.Trim(string(body), "[]"), ",")

	return downloadList

}

// get jwt token from czds host
func GetAuthToken(server string, username string, password string, client *http.Client) string {
	jsonData := []byte(`{"username":"` + username + `","password":"` + password + `"}`)
	req, err := http.NewRequest("POST", server, bytes.NewBuffer(jsonData))
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("Accept", `application/json`)

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	// The lines below will need to be cleaned up
	//fmt.Println("response Status:", resp.Status)
	//fmt.Println("response Headers:", resp.Header)
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	//fmt.Println("response Body:", string(body))
	var data AccessToken
	err = json.Unmarshal(body, &data)
	if err != nil {
		log.Fatal(err)
	}
	if data.Message != successMessage {
		log.Fatal(data.Message)
	}

	return string(data.AccessToken)
}
