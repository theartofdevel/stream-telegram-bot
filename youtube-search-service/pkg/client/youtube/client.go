package youtube

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
)

type client struct {
	accessToken, url string
	httpClient       *http.Client
}

type Client interface {
	SearchTrack(ctx context.Context, trackName string) (response *http.Response, err error)
}

func NewClient(url, accessToken string, httpClient *http.Client) *client {
	return &client{accessToken: accessToken, url: url, httpClient: httpClient}
}

func (c *client) SearchTrack(ctx context.Context, trackName string) (response *http.Response, err error) {
	params := map[string]string{
		"part":       "snippet",
		"maxResults": "1",
		"q":          trackName,
		"type":       "video",
	}

	uri, err := url.ParseRequestURI(fmt.Sprintf("%s/search", c.url))
	if err != nil {
		return nil, err
	}
	query := uri.Query()
	for k, v := range params {
		query.Set(k, v)
	}
	uri.RawQuery = query.Encode()

	request, err := http.NewRequestWithContext(ctx, http.MethodGet, uri.String(), nil)
	if err != nil {
		return nil, err
	}

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.accessToken))

	return c.httpClient.Do(request)
}
