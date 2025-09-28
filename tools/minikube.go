package tools

import (
	"fmt"
	"os/exec"

	"github.com/Antoine-Flo/kubelocal/pkg/logger"
)

// InstallMinikube installs Minikube
func InstallMinikube(log logger.Logger) error {
	log.Info("Downloading Minikube...")

	// Download Minikube
	downloadCmd := exec.Command("curl", "-LO", "https://github.com/kubernetes/minikube/releases/latest/download/minikube-linux-amd64")
	if err := log.LogCommand(downloadCmd); err != nil {
		return fmt.Errorf("failed to download Minikube: %w", err)
	}

	log.Info("Installing Minikube to /usr/local/bin...")
	// Install Minikube
	installCmd := exec.Command("sudo", "install", "minikube-linux-amd64", "/usr/local/bin/minikube")
	if err := log.LogCommand(installCmd); err != nil {
		return fmt.Errorf("failed to install Minikube: %w", err)
	}

	// Clean up
	log.Debug("Cleaning up Minikube installation files...")
	cleanupCmd := exec.Command("rm", "-f", "minikube-linux-amd64")
	cleanupCmd.Run()

	log.Info("Minikube installed successfully")
	return nil
}
