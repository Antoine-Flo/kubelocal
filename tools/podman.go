package tools

import (
	"fmt"

	"github.com/Antoine-Flo/kubelocal/internal/command"
	"github.com/Antoine-Flo/kubelocal/internal/logger"
)

func InstallPodman(log *logger.Logger) error {
	log.Info("Installing Podman...")

	log.Info("Updating package list...")
	if err := command.Run(log, "sudo", "apt-get", "update"); err != nil {
		return fmt.Errorf("failed to update packages: %w", err)
	}

	log.Info("Installing Podman package...")
	if err := command.Run(log, "sudo", "apt-get", "-y", "install", "podman"); err != nil {
		return fmt.Errorf("failed to install Podman: %w", err)
	}

	log.Info("Podman installed successfully")
	return nil
}
