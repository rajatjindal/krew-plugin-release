package actions

import (
	"fmt"

	"github.com/google/go-github/github"
)

//GetReleaseInfo gets the release info
func GetReleaseInfo() (*github.RepositoryRelease, error) {
	payload, err := GetPayload()
	if err != nil {
		return nil, err
	}

	e, err := github.ParseWebHook("release", payload)
	if err != nil {
		return nil, err
	}

	event, ok := e.(*github.ReleaseEvent)
	if !ok {
		return nil, fmt.Errorf("invalid event data")
	}

	if len(event.Release.Assets) == 0 {
		return nil, fmt.Errorf("no assets found")
	}

	return event.Release, nil
}
