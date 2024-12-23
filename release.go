package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func getRelease(owner, repo string) (*Release, error) {
	url := fmt.Sprintf(ghAPIURL, owner, repo)
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to get metadata: %w", err)
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read metadata: %w", err)
	}

	var release Release
	err = json.Unmarshal(data, &release)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal metadata: %w", err)
	}

	return &release, nil
}
