package iframely

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
)

type fakeClient struct {
	httpClient httpClient
}

func NewFakeHttpClient() *fakeClient {
	return &fakeClient{
		httpClient: nil,
	}
}

func (f *fakeClient) Get(url string) (*http.Response, error) {
	iframelyResponse, _ := json.Marshal(&iframelyApi{
		Id:  "123456",
		Url: "https://example.com/",
	})

	return &http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(bytes.NewReader(iframelyResponse)),
	}, nil
}
