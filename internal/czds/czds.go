package czds

import (
	"bytes"
	"io/ioutil"
	"log"
	"net/http"
)

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
	body, _ := ioutil.ReadAll(resp.Body)
	//fmt.Println("response Body:", string(body))
	return string(body)
}
