package tools

import (
	"fmt"
	"os/exec"

	"github.com/Antoine-Flo/kubelocal/pkg/logger"
)

// InstallPodman installs Podman
func InstallPodman(log *logger.Logger) error {
	log.Info("Installing Podman...")

	// Update packages
	log.Info("Updating package list...")
	updateCmd := exec.Command("sudo", "apt-get", "update")
	if err := log.LogCommand(updateCmd); err != nil {
		return fmt.Errorf("failed to update packages: %w", err)
	}

	// Install Podman
	log.Info("Installing Podman package...")
	installCmd := exec.Command("sudo", "apt-get", "-y", "install", "podman")
	if err := log.LogCommand(installCmd); err != nil {
		return fmt.Errorf("failed to install Podman: %w", err)
	}

	log.Info("Podman installed successfully")
	return nil
}
