package actions

import (
	"fmt"
	"net/url"

	"github.com/go-resty/resty"
)

//RepoExists checks if repo exists
func RepoExists(owner, repo string) (bool, error) {
	u := url.URL{
		Host:   "github.com",
		Scheme: "https",
		Path:   fmt.Sprintf("%s/%s", owner, repo),
	}

	resp, err := resty.New().R().Get(u.String())
	if err != nil {
		return false, err
	}

	if resp.IsError() {
		return false, fmt.Errorf("received status %d from %s", resp.StatusCode(), u.String())
	}

	return true, nil
}
