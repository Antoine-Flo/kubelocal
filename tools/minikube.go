package tools

import (
	"fmt"

	"github.com/Antoine-Flo/kubelocal/internal/command"
	"github.com/Antoine-Flo/kubelocal/internal/logger"
)

func InstallMinikube(log *logger.Logger) error {
	log.Info("Downloading Minikube...")

	if err := command.Run(log, "curl", "-LO", "https://github.com/kubernetes/minikube/releases/latest/download/minikube-linux-amd64"); err != nil {
		return fmt.Errorf("failed to download Minikube: %w", err)
	}

	log.Info("Installing Minikube to /usr/local/bin...")
	if err := command.Run(log, "sudo", "install", "minikube-linux-amd64", "/usr/local/bin/minikube"); err != nil {
		return fmt.Errorf("failed to install Minikube: %w", err)
	}

	log.Debug("Cleaning up Minikube installation files...")
	command.Run(log, "rm", "-f", "minikube-linux-amd64")

	log.Info("Minikube installed successfully")
	return nil
}
