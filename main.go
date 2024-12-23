package main

import (
	"flag"
	"fmt"
	"os"
)

const (
	// GitHub API URL
	ghAPIURL = "https://api.github.com/repos/%s/%s/releases/latest"
)

var (
	owner       = flag.String("owner", "", "owner (mandatory)")
	repo        = flag.String("repo", "", "repository (mandatory)")
	archive     = flag.String("archive", "zip", "preferred archive type")
	osType      = flag.String("os", "", "operating system")
	arch        = flag.String("arch", "", "architecture")
	extractPath = flag.String("extract", ".", "path to extract archive")
	restart     = flag.String("restart", "", "unit name to restart systemd service")
)

type runConfig struct {
	owner       string
	repo        string
	archive     string
	osType      string
	arch        string
	extractPath string
	restart     string
}

func main() {
	flag.Parse()
	err := validateFlags(*owner, *repo)
	if err != nil {
		panic(err)
	}

	if err := run(runConfig{
		owner:       *owner,
		repo:        *repo,
		archive:     *archive,
		osType:      *osType,
		arch:        *arch,
		extractPath: *extractPath,
		restart:     *restart,
	}); err != nil {
		panic(err)
	}
}

func run(c runConfig) error {
	release, err := getRelease(c.owner, c.repo)
	if err != nil {
		return fmt.Errorf("failed to get release: %w", err)
	}
	asset, err := downloadAsset(release, c.archive, c.osType, c.arch)
	if err != nil {
		return fmt.Errorf("failed to download asset: %w", err)
	}
	wd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get working directory: %w", err)
	}

	extractPath := c.extractPath
	if extractPath == "." {
		extractPath = fmt.Sprintf("%s%c", wd, os.PathSeparator)
	}

	err = extract(asset.Name, extractPath)
	if err != nil {
		return fmt.Errorf("failed to extract asset: %w", err)
	}

	if c.restart != "" {
		err = restartSystemDService(c.restart)
		if err != nil {
			return fmt.Errorf("failed to restart systemd service: %w", err)
		}
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
