package main

import (
	"net/http"
)

// doGET performs a GET request to url. If token is non-empty, adds
// Authorization: Bearer <token> so private GitHub repos can be accessed.
func doGET(url, token string) (*http.Response, error) {
	if token == "" {
		return http.Get(url)
	}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+token)
	return http.DefaultClient.Do(req)
}
