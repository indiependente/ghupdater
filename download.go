package main

import (
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
	var contentType string
	switch archive {
	case "zip":
		contentType = "zip"
	case "tar":
		contentType = "x-tar"
	case "tar.gz":
		contentType = "gzip"
	case "tar.bz2":
		contentType = "bzip2"
	case "tar.xz":
		contentType = "xz"
	default:
		return nil, errors.New("unsupported archive type")
	}
	for _, asset := range release.Assets {
		if asset.ContentType == fmt.Sprintf("application/%s", contentType) &&
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
