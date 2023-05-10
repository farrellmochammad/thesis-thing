package middleware

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

func JkdPost(url string, payload interface{}) {

	client := &http.Client{}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		// handle error
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		// handle error
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		// handle error
	}
	defer resp.Body.Close()

	_, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		// handle error
	}

}

func JkdPut(url string, payload interface{}) {

	client := &http.Client{}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		// handle error
	}

	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(jsonData))
	if err != nil {
		// handle error
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		// handle error
	}
	defer resp.Body.Close()

	_, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		// handle error
	}

}

func JkdPutFile(url string, jsonbytes []byte) {
	// Create a new PUT request with the JSON byte slice as the request body
	req, err := http.NewRequest(http.MethodPut, url, bytes.NewReader(jsonbytes))
	if err != nil {
		panic(err)
	}

	// Set the Content-Type header to "application/json"
	req.Header.Set("Content-Type", "application/json")

	// Send the request and get the response
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
}
