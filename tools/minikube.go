package tools

import (
	"fmt"
	"os/exec"
)

// InstallMinikube installs Minikube
func InstallMinikube() error {
	// Download Minikube
	if err := exec.Command("curl", "-LO", "https://github.com/kubernetes/minikube/releases/latest/download/minikube-linux-amd64").Run(); err != nil {
		return fmt.Errorf("failed to download Minikube: %w", err)
	}

	// Install Minikube
	if err := exec.Command("sudo", "install", "minikube-linux-amd64", "/usr/local/bin/minikube").Run(); err != nil {
		return fmt.Errorf("failed to install Minikube: %w", err)
	}

	// Clean up
	exec.Command("rm", "-f", "minikube-linux-amd64").Run()

	fmt.Println("✅ Minikube installed successfully")
	return nil
}
