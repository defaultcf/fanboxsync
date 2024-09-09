package iframely_test

import (
	"testing"

	. "github.com/defaultcf/fanboxsync/iframely"
	"github.com/stretchr/testify/assert"
)

func TestClient_GetReal(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		iframelyUrl string
		want        string
	}{
		{
			name:        "JSON をパースし、URL を取得できる",
			iframelyUrl: "https://cdn.iframe.ly/123.json",
			want:        "https://example.com/",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			//setup
			fakeClient := NewFakeHttpClient()
			iframelyClient := NewIframelyClient(fakeClient)

			// execute
			url, err := iframelyClient.GetRealUrl(tt.iframelyUrl)

			// verify
			assert.NoError(t, err)
			assert.Equal(t, tt.want, url)
		})
	}
}
