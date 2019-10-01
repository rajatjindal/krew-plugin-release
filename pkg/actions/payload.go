package actions

import (
	"io/ioutil"
	"os"
)

//GetPayload reads payload and returns it
func GetPayload() ([]byte, error) {
	eventJSONPath := os.Getenv("GITHUB_EVENT_PATH")
	data, err := ioutil.ReadFile(eventJSONPath)
	if err != nil {
		return nil, err
	}

	return data, nil
}
