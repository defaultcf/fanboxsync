package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"regexp"
)

type iframelyApi struct {
	Id  string
	Url string
}

// https://iframely.com/docs/iframely-api

func GetRealUrl(iframelyUrl string) (string, error) {
	re := regexp.MustCompile(`^https:\/\/cdn\.iframe\.ly\/(\w+)`)
	matches := re.FindAllStringSubmatch(iframelyUrl, -1)
	if len(matches) != 1 {
		return "", errors.New("unexpected iframely url")
	}
	iframelyId := matches[0][1]

	response, err := http.Get(fmt.Sprintf("https://cdn.iframe.ly/%s.json", iframelyId))
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	parsedBody, err := parseIframely(response.Body)
	if err != nil {
		return "", err
	}

	return parsedBody.Url, nil
}

func parseIframely(data io.Reader) (*iframelyApi, error) {
	iframely := &iframelyApi{}
	bytes, err := io.ReadAll(data)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(bytes, iframely)
	if err != nil {
		return nil, err
	}
	return iframely, nil
}
