package czds

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"path"
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
func GetZoneFile(url string, localDirectory string, token string, client *http.Client) string {
	url = strings.Trim(url, "\"")

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Add("Authorization", "Bearer "+token)
	zoneFile := path.Base(req.URL.Path)
	out, err := os.Create(localDirectory + "/" + zoneFile)
	if err != nil {
		log.Fatal(err)
		return ""
	}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	defer out.Close()
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	return out.Name()
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
	body, err := io.ReadAll(resp.Body)
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
