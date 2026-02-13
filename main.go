package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"time"
)

var (
	owner       = flag.String("owner", "", "owner (mandatory)")
	repo        = flag.String("repo", "", "repository (mandatory)")
	token       = flag.String("token", "", "GitHub token for private repos (or set GITHUB_TOKEN)")
	archive     = flag.String("archive", "", "archive type")
	osType      = flag.String("os", "", "operating system")
	arch        = flag.String("arch", "", "architecture")
	list        = flag.Bool("list", false, "list available assets")
	extractPath = flag.String("extract", ".", "path to extract archive")
	restart     = flag.String("restart", "", "unit name to restart systemd service")
	cleanup     = flag.Bool("cleanup", false, "remove archive after extraction")
	dryRun      = flag.Bool("dry-run", false, "simulate actions without changes")
	showVersion = flag.Bool("version", false, "show version")
)

var version = "dev"

type runConfig struct {
	owner       string
	repo        string
	token       string
	archive     string
	osType      string
	arch        string
	extractPath string
	restart     string
	list        bool
	cleanup     bool
	dryRun      bool
}

func main() {
	flag.Parse()

	if *showVersion {
		fmt.Printf("ghupdater version %s\n", version)
		os.Exit(0)
	}

	err := validateFlags(*owner, *repo)
	if err != nil {
		panic(err)
	}

	tokenVal := *token
	if tokenVal == "" {
		tokenVal = os.Getenv("GITHUB_TOKEN")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()

	if err := run(ctx, runConfig{
		owner:       *owner,
		repo:        *repo,
		token:       tokenVal,
		archive:     *archive,
		osType:      *osType,
		arch:        *arch,
		extractPath: *extractPath,
		restart:     *restart,
		list:        *list,
		cleanup:     *cleanup,
		dryRun:      *dryRun,
	}); err != nil {
		panic(err)
	}
}

func run(ctx context.Context, c runConfig) error {
	fmt.Printf("Get %s/%s latest release\n", c.owner, c.repo)
	release, err := getRelease(ctx, c.owner, c.repo, c.token)
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

	if c.dryRun {
		asset, err := selectAsset(release, c.archive, c.osType, c.arch)
		if err != nil {
			return fmt.Errorf("failed to select asset: %w", err)
		}
		fmt.Printf("[Dry Run] Would download asset %s\n", asset.Name)
		fmt.Printf("[Dry Run] Would extract to %s\n", c.extractPath)
		if c.cleanup {
			fmt.Printf("[Dry Run] Would remove %s after extraction\n", asset.Name)
		}
		if c.restart != "" {
			fmt.Printf("[Dry Run] Would restart systemd service %s\n", c.restart)
		}
		return nil
	}

	asset, err := downloadAsset(ctx, release, c.archive, c.osType, c.arch, c.token)
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

	if c.cleanup {
		fmt.Printf("Removing asset %s\n", asset.Name)
		if err := os.Remove(asset.Name); err != nil {
			return fmt.Errorf("failed to remove asset: %w", err)
		}
		fmt.Println("Asset removed")
	}

	if c.restart != "" {
		fmt.Printf("Restarting systemd service %s\n", c.restart)
		err = restartSystemDService(ctx, c.restart)
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
