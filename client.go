package main

import (
	"context"
	"net/http"
)

// doGET performs a GET request to url. If token is non-empty, adds
// Authorization: Bearer <token> so private GitHub repos can be accessed.
func doGET(ctx context.Context, url, token string) (*http.Response, error) {
	if token == "" {
		req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
		if err != nil {
			return nil, err
		}
		return http.DefaultClient.Do(req)
	}
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+token)
	return http.DefaultClient.Do(req)
}

// doGETAsset downloads a release asset. For public repos (token empty), GETs
// browser_download_url. For private repos (token set), GETs the asset API url
// with Accept: application/octet-stream and Bearer token, as required by
// GitHub (browser_download_url does not accept API token auth).
func doGETAsset(ctx context.Context, downloadURL, apiURL, token string) (*http.Response, error) {
	if token == "" {
		req, err := http.NewRequestWithContext(ctx, "GET", downloadURL, nil)
		if err != nil {
			return nil, err
		}
		return http.DefaultClient.Do(req)
	}
	req, err := http.NewRequestWithContext(ctx, "GET", apiURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Accept", "application/octet-stream")
	return http.DefaultClient.Do(req)
}
