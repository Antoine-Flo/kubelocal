package tools

import (
	"fmt"
	"os/exec"
)

// InstallKubectl downloads and installs kubectl for linux amd64 with checksum validation
func InstallKubectl() error {
	// Get the latest stable version
	getVersionCmd := exec.Command("curl", "-L", "-s", "https://dl.k8s.io/release/stable.txt")
	versionBytes, err := getVersionCmd.Output()
	if err != nil {
		return fmt.Errorf("failed to get kubectl version: %w", err)
	}
	version := string(versionBytes)

	// Download kubectl binary
	downloadURL := fmt.Sprintf("https://dl.k8s.io/release/%s/bin/linux/amd64/kubectl", version)
	downloadCmd := exec.Command("curl", "-LO", downloadURL)
	if err := downloadCmd.Run(); err != nil {
		return fmt.Errorf("failed to download kubectl: %w", err)
	}

	// Download kubectl checksum
	checksumURL := fmt.Sprintf("https://dl.k8s.io/release/%s/bin/linux/amd64/kubectl.sha256", version)
	checksumCmd := exec.Command("curl", "-LO", checksumURL)
	if err := checksumCmd.Run(); err != nil {
		return fmt.Errorf("failed to download kubectl checksum: %w", err)
	}

	// Validate the kubectl binary against the checksum file
	validateCmd := exec.Command("bash", "-c", "echo \"$(cat kubectl.sha256)  kubectl\" | sha256sum --check")
	if err := validateCmd.Run(); err != nil {
		return fmt.Errorf("kubectl checksum validation failed: %w", err)
	}

	// Install kubectl to /usr/local/bin
	installCmd := exec.Command("sudo", "install", "-o", "root", "-g", "root", "-m", "0755", "kubectl", "/usr/local/bin/kubectl")
	if err := installCmd.Run(); err != nil {
		return fmt.Errorf("failed to install kubectl: %w", err)
	}

	// Clean up downloaded files
	cleanupCmd := exec.Command("rm", "-f", "kubectl", "kubectl.sha256")
	cleanupCmd.Run() // Ignore error as files might not exist

	return nil
}

// SetupKubectlAlias adds the k=kubectl alias to bashrc
func SetupKubectlAlias() error {
	// Add alias to bashrc
	addCmd := exec.Command("bash", "-c", "echo 'alias k=kubectl' >> ~/.bashrc")
	if err := addCmd.Run(); err != nil {
		return fmt.Errorf("failed to add kubectl alias: %w", err)
	}

	return nil
}
