package czds

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

// get jwt token from czds host
func GetAuthToken(server string, username string, password string) string {
	jsonData := []byte(`{"username":"` + username + `","password":"` + password + `"}`)

	resp, err := http.Post(server, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	// The lines below will need to be cleaned up
	fmt.Println("response Status:", resp.Status)
	fmt.Println("response Headers:", resp.Header)
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("response Body:", string(body))
	return string(body)
}
