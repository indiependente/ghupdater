package main

import (
	"context"
	"fmt"
	"time"

	"github.com/coreos/go-systemd/v22/dbus"
)

const (
	restartMode = "replace"
	timeout     = 30 * time.Second
)

func restartSystemDService(ctx context.Context, targetSystemdUnit string) error {
	systemdConnection, err := dbus.NewSystemConnectionContext(ctx)
	if err != nil {
		return fmt.Errorf("failed to connect to systemd: %w", err)
	}
	defer systemdConnection.Close()

	found, err := findUnitByName(ctx, targetSystemdUnit, systemdConnection)
	if err != nil {
		return fmt.Errorf("failed to find unit: %w", err)
	}
	if !found {
		return fmt.Errorf("unit %s not found", targetSystemdUnit)
	}

	err = restartUnit(ctx, targetSystemdUnit, systemdConnection)
	if err != nil {
		return fmt.Errorf("failed to restart unit: %w", err)
	}

	return nil
}

func findUnitByName(ctx context.Context, targetSystemdUnit string, conn *dbus.Conn) (bool, error) {
	units, err := conn.ListUnitsContext(ctx)
	if err != nil {
		return false, fmt.Errorf("failed to list units: %w", err)
	}

	for _, unit := range units {
		if unit.Name == targetSystemdUnit {
			return true, nil
		}
	}

	return false, nil
}

func restartUnit(ctx context.Context, targetSystemdUnit string, conn *dbus.Conn) error {
	completedRestartCh := make(chan string)
	_, err := conn.RestartUnitContext(
		ctx,
		targetSystemdUnit,
		restartMode,
		completedRestartCh,
	)
	if err != nil {
		return fmt.Errorf("failed to restart unit: %w", err)
	}

	// Wait for the restart to complete
	select {
	case res := <-completedRestartCh:
		if res != "done" {
			return fmt.Errorf("restart job failed: %s", res)
		}
		fmt.Printf("Restart job completed for unit: %s\n", targetSystemdUnit)
	case <-time.After(timeout):
		return fmt.Errorf("timed out waiting for restart job to complete for unit: %s", targetSystemdUnit)
	case <-ctx.Done():
		return ctx.Err()
	}
	return nil
}
