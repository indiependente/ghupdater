package main

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

func selectAsset(release *Release, archive, osType, arch string) (*Asset, error) {
	for _, asset := range release.Assets {
		if asset.ContentType == fmt.Sprintf("application/%s", archive) &&
			strings.Contains(asset.Name, osType) &&
			strings.Contains(asset.Name, arch) {
			return &asset, nil
		}
	}

	return nil, errors.New("failed to select asset")
}

func downloadAsset(release *Release, archive, osType, arch string) (*Asset, error) {
	asset, err := selectAsset(release, archive, osType, arch)
	if err != nil {
		return nil, fmt.Errorf("failed to find asset: %w", err)
	}

	resp, err := http.Get(asset.BrowserDownloadURL)
	if err != nil {
		return nil, fmt.Errorf("failed to download asset: %w", err)
	}
	defer resp.Body.Close()

	out, err := os.Create(asset.Name)
	if err != nil {
		return nil, fmt.Errorf("failed to create file: %w", err)
	}
	defer out.Close()

	written, err := io.Copy(out, resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to write asset to file: %w", err)
	}

	if written != asset.Size {
		return nil, errors.New("failed to download asset: size mismatch")
	}

	return asset, nil
}
