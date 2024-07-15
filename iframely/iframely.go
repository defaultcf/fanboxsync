package iframely

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"regexp"
)

type httpClient interface {
	Get(url string) (*http.Response, error)
}

type IframelyClient struct {
	HttpClient httpClient
}

type iframelyApi struct {
	Id  string
	Url string
}

func NewIframelyClient(httpClient httpClient) *IframelyClient {
	return &IframelyClient{
		HttpClient: httpClient,
	}
}

// https://iframely.com/docs/iframely-api

func (c *IframelyClient) GetRealUrl(iframelyUrl string) (string, error) {
	re := regexp.MustCompile(`^https:\/\/cdn\.iframe\.ly\/(\w+)`)
	matches := re.FindAllStringSubmatch(iframelyUrl, -1)
	if len(matches) != 1 {
		return "", errors.New("unexpected iframely url")
	}
	iframelyId := matches[0][1]

	response, err := c.HttpClient.Get(fmt.Sprintf("https://cdn.iframe.ly/%s.json", iframelyId))
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	return c.parseIframely(response.Body)
}

func (c *IframelyClient) parseIframely(data io.Reader) (string, error) {
	iframely := &iframelyApi{}
	bytes, err := io.ReadAll(data)
	if err != nil {
		return "", err
	}
	err = json.Unmarshal(bytes, iframely)
	if err != nil {
		return "", err
	}
	return iframely.Url, nil
}
