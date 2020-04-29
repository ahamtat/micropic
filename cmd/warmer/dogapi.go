package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/pkg/errors"
)

type dogAPI struct {
	Message string `json:"message"`
	Status  string `json:"status"`
}

func getRandomDogURL() (string, error) {
	// Send request to dog API
	resp, err := http.Get("https://dog.ceo/api/breeds/image/random")
	if err != nil {
		return "", errors.Wrap(err, "error getting response from dog API")
	}
	defer resp.Body.Close()

	// Read response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", errors.Wrap(err, "error reading HTTP response")
	}

	// Convert from JSON
	message := &dogAPI{}
	err = json.Unmarshal(body, &message)
	if err != nil {
		return "", errors.Wrap(err, "failed unmarshalling JSON message")
	}

	return message.Message, nil
}
