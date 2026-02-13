package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/schollz/progressbar/v3"
)

func listAssets(release *Release, archive, osType, arch string) []string {
	assets := make([]string, 0)
	for _, asset := range release.Assets {
		if strings.Contains(asset.Name, osType) &&
			strings.Contains(asset.Name, arch) &&
			strings.Contains(asset.Name, archive) {
			assets = append(assets, asset.Name)
		}
	}

	return assets
}

func selectAsset(release *Release, archive, osType, arch string) (*Asset, error) {
	if archive == "" {
		return nil, errors.New("archive type is required")
	}
	suffix := "." + archive

	for _, asset := range release.Assets {
		if strings.HasSuffix(asset.Name, suffix) &&
			strings.Contains(asset.Name, osType) &&
			strings.Contains(asset.Name, arch) {
			return &asset, nil
		}
	}

	return nil, errors.New("failed to select asset")
}

func downloadAsset(ctx context.Context, release *Release, archive, osType, arch, token string) (*Asset, error) {
	asset, err := selectAsset(release, archive, osType, arch)
	if err != nil {
		return nil, fmt.Errorf("failed to find asset: %w", err)
	}

	resp, err := doGETAsset(ctx, asset.BrowserDownloadURL, asset.URL, token)
	if err != nil {
		return nil, fmt.Errorf("failed to download asset: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("download failed: status %d: %s", resp.StatusCode, string(body))
	}

	out, err := os.Create(asset.Name)
	if err != nil {
		return nil, fmt.Errorf("failed to create file: %w", err)
	}
	defer func() { _ = out.Close() }()

	bar := progressbar.DefaultBytes(
		asset.Size,
		asset.Name,
	)

	written, err := io.Copy(io.MultiWriter(out, bar), resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to write asset to file: %w", err)
	}

	if written != asset.Size {
		return nil, errors.New("failed to download asset: size mismatch")
	}

	return asset, nil
}
