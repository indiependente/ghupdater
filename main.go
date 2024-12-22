package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
)

const (
	// GitHub API URL
	ghAPIURL = "https://api.github.com/repos/%s/%s/releases/latest"
)

var (
	owner   = flag.String("owner", "", "owner")
	repo    = flag.String("repo", "", "repository")
	archive = flag.String("archive", "zip", "preferred archive type")
)

func main() {
	flag.Parse()
	err := validateFlags(*owner, *repo)
	if err != nil {
		panic(err)
	}
	if err := run(*owner, *repo, *archive); err != nil {
		panic(err)
	}
}

func run(owner, repo, archive string) error {
	release, err := getRelease(owner, repo)
	if err != nil {
		return fmt.Errorf("failed to get release: %w", err)
	}
	err = downloadAsset(release, archive)
	if err != nil {
		return fmt.Errorf("failed to download asset: %w", err)
	}

	return nil
}

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

func downloadAsset(release *Release, archive string) error {
	for _, asset := range release.Assets {
		if asset.ContentType != fmt.Sprintf("application/%s", archive) {
			continue
		}
		resp, err := http.Get(asset.BrowserDownloadURL)
		if err != nil {
			return fmt.Errorf("failed to download asset: %w", err)
		}
		defer resp.Body.Close()

		out, err := os.Create(asset.Name)
		if err != nil {
			return fmt.Errorf("failed to create file: %w", err)
		}
		defer out.Close()

		written, err := io.Copy(out, resp.Body)
		if err != nil {
			return fmt.Errorf("failed to write asset to file: %w", err)
		}

		if written != asset.Size {
			return errors.New("failed to download asset: size mismatch")
		}

		return nil
	}

	return nil
}

func validateFlags(owner, repo string) error {
	if owner == "" {
		return fmt.Errorf("owner is required")
	}
	if repo == "" {
		return fmt.Errorf("repo is required")
	}

	return nil
}
