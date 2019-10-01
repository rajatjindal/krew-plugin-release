package actions

import (
	"fmt"

	"github.com/google/go-github/github"
	"github.com/sirupsen/logrus"
)

//GetReleaseInfo gets the release info
func GetReleaseInfo() error {
	payload, err := GetPayload()
	if err != nil {
		return err
	}

	e, err := github.ParseWebHook("release", payload)
	if err != nil {
		return err
	}

	event, ok := e.(*github.ReleaseEvent)
	if !ok {
		return fmt.Errorf("invalid event data")
	}

	logrus.Info(event)
	return nil
}
