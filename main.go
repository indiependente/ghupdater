package main

import (
	"flag"
	"fmt"
	"os"
)

var (
	owner       = flag.String("owner", "", "owner (mandatory)")
	repo        = flag.String("repo", "", "repository (mandatory)")
	archive     = flag.String("archive", "", "archive type")
	osType      = flag.String("os", "", "operating system")
	arch        = flag.String("arch", "", "architecture")
	list        = flag.Bool("list", false, "list available assets")
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
	list        bool
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
		list:        *list,
	}); err != nil {
		panic(err)
	}
}

func run(c runConfig) error {
	fmt.Printf("Get %s/%s latest release\n", c.owner, c.repo)
	release, err := getRelease(c.owner, c.repo)
	if err != nil {
		return fmt.Errorf("failed to get release: %w", err)
	}
	if c.list {
		assets := listAssets(release, c.archive, c.osType, c.arch)
		if len(assets) == 0 {
			fmt.Printf("No assets found with filters archive=%s os=%s arch=%s\n", c.archive, c.osType, c.arch)
		} else {
			fmt.Printf("Available assets with filters archive=%s os=%s arch=%s\n", c.archive, c.osType, c.arch)
			for _, asset := range assets {
				fmt.Println(asset)
			}
		}

		return nil
	}
	fmt.Printf("Selected tag %s\n", release.TagName)
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

	fmt.Printf("Extracting asset %s to %s\n", asset.Name, extractPath)

	err = extract(asset.Name, extractPath)
	if err != nil {
		return fmt.Errorf("failed to extract asset: %w", err)
	}
	fmt.Println("Asset extracted")

	if c.restart != "" {
		fmt.Printf("Restarting systemd service %s\n", c.restart)
		err = restartSystemDService(c.restart)
		if err != nil {
			return fmt.Errorf("failed to restart systemd service: %w", err)
		}
		fmt.Println("Systemd service restarted")
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
