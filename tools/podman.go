package tools

import (
	"fmt"
	"os/exec"
)

// InstallPodman installs Podman
func InstallPodman() error {
	fmt.Println("🚀 Installing Podman...")

	// Update packages
	if err := exec.Command("sudo", "apt-get", "update").Run(); err != nil {
		return fmt.Errorf("failed to update packages: %w", err)
	}

	// Install Podman
	if err := exec.Command("sudo", "apt-get", "-y", "install", "podman").Run(); err != nil {
		return fmt.Errorf("failed to install Podman: %w", err)
	}

	fmt.Println("✅ Podman installed successfully")
	return nil
}
