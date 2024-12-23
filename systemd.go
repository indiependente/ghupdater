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

func restartSystemDService(targetSystemdUnit string) error {
	ctx := context.Background()
	systemdConnection, err := dbus.NewSystemConnectionContext(ctx)
	if err != nil {
		return fmt.Errorf("Failed to connect to systemd: %w", err)
	}
	defer systemdConnection.Close()

	found, err := findUnitByName(ctx, targetSystemdUnit, systemdConnection)
	if err != nil {
		return fmt.Errorf("Failed to find unit: %w", err)
	}
	if !found {
		return fmt.Errorf("Unit %s not found", targetSystemdUnit)
	}

	err = restartUnit(ctx, targetSystemdUnit, systemdConnection)
	if err != nil {
		return fmt.Errorf("Failed to restart unit: %w", err)
	}

	return nil
}

func findUnitByName(ctx context.Context, targetSystemdUnit string, conn *dbus.Conn) (bool, error) {
	units, err := conn.ListUnitsContext(ctx)
	if err != nil {
		return false, fmt.Errorf("Failed to list units: %w", err)
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
		return fmt.Errorf("Failed to restart unit: %w", err)
	}

	// Wait for the restart to complete
	select {
	case <-completedRestartCh:
		fmt.Printf("Restart job completed for unit: %s\n", targetSystemdUnit)
	case <-time.After(timeout):
		fmt.Printf("Timed out waiting for restart job to complete for unit: %s\n", targetSystemdUnit)
	}
	return nil
}
