package krew

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"github.com/go-resty/resty"
)

func getSha256ForAsset(uri string) (string, error) {
	client := resty.New()

	file, err := ioutil.TempFile("", "")
	if err != nil {
		return "", err
	}
	defer os.Remove(file.Name())

	resp, err := client.R().SetOutput(file.Name()).Get(uri)
	if err != nil {
		return "", err
	}

	if resp.IsError() {
		return "", fmt.Errorf("received response-code %d from %s", resp.StatusCode(), uri)
	}

	sha256, err := getSha256(file.Name())
	if err != nil {
		return "", err
	}

	return sha256, nil
}

func getSha256(filename string) (string, error) {
	f, err := os.Open(filename)
	if err != nil {
		return "", err
	}
	defer f.Close()

	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}

	return hex.EncodeToString(h.Sum(nil)), nil
}
