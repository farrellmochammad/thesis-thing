package middleware

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"

	MQTT "github.com/eclipse/paho.mqtt.golang"
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

func PublishMessage(client MQTT.Client, topic string, payload interface{}) {

	// encode the struct as JSON
	jsonData, err := json.Marshal(payload)
	if err != nil {
		panic(err)
	}

	token := client.Publish(topic, 0, false, jsonData)
	token.Wait()
}
