package tools

import (
	"fmt"

	"github.com/Antoine-Flo/kubelocal/internal/command"
	"github.com/Antoine-Flo/kubelocal/internal/logger"
)

func InstallK3s(log *logger.Logger) error {
	log.Info("Installing K3s...")

	// Install K3s using the official installer
	if err := command.Run(log, "curl", "-sfL", "https://get.k3s.io", "|", "sh", "-"); err != nil {
		return fmt.Errorf("failed to install K3s: %w", err)
	}

	// Enable and start K3s service
	log.Info("Enabling K3s service...")
	if err := command.Run(log, "sudo", "systemctl", "enable", "k3s"); err != nil {
		return fmt.Errorf("failed to enable K3s service: %w", err)
	}

	log.Info("Starting K3s service...")
	if err := command.Run(log, "sudo", "systemctl", "start", "k3s"); err != nil {
		return fmt.Errorf("failed to start K3s service: %w", err)
	}

	// Set up kubectl access for the current user
	log.Info("Setting up kubectl access...")
	if err := command.Run(log, "mkdir", "-p", "~/.kube"); err != nil {
		return fmt.Errorf("failed to create .kube directory: %w", err)
	}

	if err := command.Run(log, "sudo", "cp", "/etc/rancher/k3s/k3s.yaml", "~/.kube/config"); err != nil {
		return fmt.Errorf("failed to copy K3s config: %w", err)
	}

	if err := command.Run(log, "sudo", "chown", "$USER", "~/.kube/config"); err != nil {
		return fmt.Errorf("failed to change ownership of kubeconfig: %w", err)
	}

	log.Info("K3s installed and configured successfully")
	return nil
}
