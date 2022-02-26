package youtube

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"path"
)

type Client struct {
	url, accessToken, refreshToken, clientID, clientSecret, refreshTokenURL, authRedirectUri, accountsUri string
	httpClient                                                                                            *http.Client
}

func New(url, accessToken, refreshToken, clientID, clientSecret, refreshTokenURL, authRedirectUri, accountsUri string, client *http.Client) Client {
	if client == nil {
		client = &http.Client{}
	}
	return Client{
		url:             url,
		accessToken:     accessToken,
		refreshToken:    refreshToken,
		clientID:        clientID,
		clientSecret:    clientSecret,
		refreshTokenURL: refreshTokenURL,
		authRedirectUri: authRedirectUri,
		accountsUri:     accountsUri,
		httpClient:      client,
	}
}

func (c *Client) SetAccessToken(val string) {
	c.accessToken = val
}

func (c *Client) SetRefreshToken(val string) {
	c.refreshToken = val
}

func (c *Client) SearchTrack(ctx context.Context, trackName string) (response *http.Response, err error) {
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

func (c *Client) UpdateAccessToken(ctx context.Context) (response *http.Response, err error) {
	parsedURL, err := url.ParseRequestURI(c.refreshTokenURL)
	if err != nil {
		return response, fmt.Errorf("failed to refresh token URL. error: %w", err)
	}
	query := parsedURL.Query()
	parsedURL.RawQuery = query.Encode()
	uri := parsedURL.String()

	data := map[string]string{
		"client_id":     c.clientID,
		"client_secret": c.clientSecret,
		"refresh_token": c.refreshToken,
		"grant_type":    "refresh_token",
	}

	dataBytes, err := json.Marshal(data)
	if err != nil {
		return response, fmt.Errorf("failed to marshal dto")
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, uri, bytes.NewBuffer(dataBytes))
	if err != nil {
		return response, fmt.Errorf("failed to create new request due to error: %v", err)
	}

	response, err = c.httpClient.Do(req)
	if err != nil {
		return response, fmt.Errorf("failed to send request due to error: %v", err)
	}

	return response, nil
}

func (c *Client) BuildURL(resource string, parameters map[string]string, baseUrl string) (string, error) {
	var resultURL string
	parsedURL, err := url.ParseRequestURI(baseUrl)
	if err != nil {
		return resultURL, fmt.Errorf("failed to parse base URL. error: %w", err)
	}
	parsedURL.Path = path.Join(parsedURL.Path, resource)

	query := parsedURL.Query()
	for k, v := range parameters {
		query.Set(k, v)
	}
	parsedURL.RawQuery = query.Encode()

	return parsedURL.String(), nil
}
