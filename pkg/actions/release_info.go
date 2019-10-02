package actions

import (
	"crypto/sha256"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"

	"github.com/go-resty/resty"
	"github.com/google/go-github/github"
	"github.com/sirupsen/logrus"
)

type ReleaseInfo struct {
	Tag    string
	Assets []Asset
}

type Asset struct {
	DownloadURL string
	Sha256      string
	Platform    string
	OS          string
}

//GetReleaseInfo gets the release info
func GetReleaseInfo() (*ReleaseInfo, error) {
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

	client := resty.New()
	for _, releaseAsset := range event.Release.Assets {
		_, err = getAssetInfo(client, releaseAsset)
		if err != nil {
			return nil, err
		}
	}

	logrus.Info(event)
	return nil, nil
}

func getAssetInfo(client *resty.Client, releaseAsset github.ReleaseAsset) (*Asset, error) {
	file, err := ioutil.TempFile("", "")
	if err != nil {
		return nil, err
	}
	defer os.Remove(file.Name())

	resp, err := client.R().SetOutput(file.Name()).Get(releaseAsset.GetBrowserDownloadURL())
	if err != nil {
		return nil, err
	}

	if resp.IsError() {
		return nil, fmt.Errorf("received response-code %d from %s", resp.StatusCode(), releaseAsset.GetBrowserDownloadURL())
	}

	sha256, err := getSha256(file.Name())
	if err != nil {
		return nil, err
	}

	platformOS, platformArch, err := getPlatformInfo(releaseAsset.GetBrowserDownloadURL())
	logrus.Infof("platformOS: %s, platformArch: %s, sha256: %s", platformOS, platformArch, string(sha256))
	return nil, nil
}

func stringValue(s *string) string {
	if s == nil {
		return ""
	}

	return *s
}

func getSha256(filename string) ([]byte, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return nil, err
	}

	return h.Sum(nil), nil
}

func getPlatformInfo(u string) (string, string, error) {
	platformOS := ""
	platformArch := ""

	switch {
	case strings.Contains(strings.ToLower(u), "darwin"):
		platformOS = "darwin"
	case strings.Contains(strings.ToLower(u), "windows"):
		platformOS = "windows"
	default:
		platformOS = "linux"
	}

	switch {
	case strings.Contains(strings.ToLower(u), "arm"):
		platformArch = "arm"
	default:
		platformArch = "x86_64"
	}

	return platformOS, platformArch, nil
}
