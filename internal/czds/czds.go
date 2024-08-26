package czds

import (
	"bytes"
	"czdsfetch/configs"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"
)

var (
	successMessage = "Authentication Successful"
	listDownloads  = "/czds/downloads/links"
)

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
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	downloadList := strings.Split(strings.Trim(string(body), "[]"), ",")

	return downloadList

}

// get zone file
func GetZoneFile(url string, localDirectory string, token string, client *http.Client) (string, error) {
	url = strings.Trim(url, "\"")
	for i := 0; i < configs.Retries; {
		fmt.Println("Downloading ", url, " to ", localDirectory, "attempt ", i+1)
		// Create the file
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			fmt.Println("Error: ", err.Error(), " Zone: ", path.Base(req.URL.Path))
			continue
		}
		req.Header.Add("Authorization", "Bearer "+token)
		zoneFile := path.Base(req.URL.Path)
		out, err := os.Create(localDirectory + "/" + zoneFile)
		if err != nil {
			fmt.Println("Error: ", err.Error(), " Zone: ", zoneFile)
			continue
		}
		defer out.Close()
		resp, err := client.Do(req)
		if err != nil {
			fmt.Println("Error: ", err.Error(), " Zone: ", zoneFile)
			continue
		}
		defer resp.Body.Close()
		if resp.StatusCode != 200 {
			fmt.Println("HTTP status: ", resp.Status, " Zone: ", zoneFile)
			continue
		}
		var buf []byte
		if zoneFile == "com.zone" {
			buf = make([]byte, 50*1024*1024) //50MB
		} else {
			buf = make([]byte, 5*1024*1024) //5MB
		}

		_, err = io.CopyBuffer(out, resp.Body, buf)
		if err != nil {
			fmt.Println("Error: ", err.Error(), " Zone: ", zoneFile)
			continue
		}

		return out.Name(), nil
	}
	//if we reach here, we failed to download the zone file after all retries
	e := errors.New("Error: Failed to download zone file after " + strconv.Itoa(configs.Retries) + " attempts")
	return "", e
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

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

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
