package main

import (
	"context"
	"fmt"
	"os"

	xtrct "github.com/codeclysm/extract/v4"
)

func extract(assetName string, extractPath string) error {
	file, err := os.Open(assetName)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	return xtrct.Archive(context.Background(), file, extractPath, nil)
}
