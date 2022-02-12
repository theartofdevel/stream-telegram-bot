package imgur

import (
	"context"
	"encoding/base64"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

type client struct {
	accessToken, url, clientID string
	httpClient                 *http.Client
}

type Client interface {
	UploadImage(ctx context.Context, image []byte) (response *http.Response, err error)
}

func NewClient(url, accessToken, clientID string, httpClient *http.Client) *client {
	return &client{url: url, accessToken: accessToken, clientID: clientID, httpClient: httpClient}
}

func (c *client) UploadImage(ctx context.Context, image []byte) (response *http.Response, err error) {
	vals := url.Values{}
	vals.Set("image", base64.StdEncoding.EncodeToString(image))
	vals.Set("type", "base64")

	uri, err := url.ParseRequestURI(fmt.Sprintf("%s/upload", c.url))
	if err != nil {
		return nil, err
	}

	request, err := http.NewRequestWithContext(ctx, http.MethodPost, uri.String(), strings.NewReader(vals.Encode()))
	if err != nil {
		return nil, err
	}

	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	request.Header.Set("Authorization", fmt.Sprintf("Client-ID %s", c.clientID))
	request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.accessToken))

	return c.httpClient.Do(request)
}
