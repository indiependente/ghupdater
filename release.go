package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const (
	// GitHub API URL
	ghAPIURL = "https://api.github.com/repos/%s/%s/releases/latest"
)

func getRelease(ctx context.Context, owner, repo, token string) (*Release, error) {
	url := fmt.Sprintf(ghAPIURL, owner, repo)
	resp, err := doGET(ctx, url, token)
	if err != nil {
		return nil, fmt.Errorf("failed to get metadata: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read metadata: %w", err)
	}

	if resp.StatusCode == http.StatusNotFound && token == "" {
		return nil, fmt.Errorf("release not found (404). If this repo is private, set GITHUB_TOKEN or use -token")
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status %d: %s", resp.StatusCode, string(data))
	}

	var release Release
	err = json.Unmarshal(data, &release)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal metadata: %w", err)
	}

	return &release, nil
}
