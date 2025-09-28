package tools

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/Antoine-Flo/kubelocal/pkg/logger"
	"go.uber.org/zap"
)

// InstallKubectl downloads and installs kubectl for linux amd64 with checksum validation
func InstallKubectl(log *logger.Logger) error {
	log.Info("Getting latest kubectl version...")

	// Get the latest stable version
	getVersionCmd := exec.Command("curl", "-L", "-s", "https://dl.k8s.io/release/stable.txt")
	versionBytes, err := getVersionCmd.Output()
	if err != nil {
		log.Error("Failed to get kubectl version", zap.Error(err))
		return fmt.Errorf("failed to get kubectl version: %w", err)
	}
	version := strings.TrimSpace(string(versionBytes))
	log.Info("Latest kubectl version", zap.String("version", version))

	// Download kubectl binary
	log.Info("Downloading kubectl binary...")
	downloadURL := fmt.Sprintf("https://dl.k8s.io/release/%s/bin/linux/amd64/kubectl", version)
	downloadCmd := exec.Command("curl", "-LO", downloadURL)
	if err := log.LogCommand(downloadCmd); err != nil {
		return fmt.Errorf("failed to download kubectl: %w", err)
	}

	// Download kubectl checksum
	log.Info("Downloading kubectl checksum...")
	checksumURL := fmt.Sprintf("https://dl.k8s.io/release/%s/bin/linux/amd64/kubectl.sha256", version)
	checksumCmd := exec.Command("curl", "-LO", checksumURL)
	if err := log.LogCommand(checksumCmd); err != nil {
		return fmt.Errorf("failed to download kubectl checksum: %w", err)
	}

	// Validate the kubectl binary against the checksum file
	log.Info("Validating kubectl checksum...")
	validateCmd := exec.Command("bash", "-c", "echo \"$(cat kubectl.sha256)  kubectl\" | sha256sum --check")
	if err := log.LogCommand(validateCmd); err != nil {
		return fmt.Errorf("kubectl checksum validation failed: %w", err)
	}

	// Install kubectl to /usr/local/bin
	log.Info("Installing kubectl to /usr/local/bin...")
	installCmd := exec.Command("sudo", "install", "-o", "root", "-g", "root", "-m", "0755", "kubectl", "/usr/local/bin/kubectl")
	if err := log.LogCommand(installCmd); err != nil {
		return fmt.Errorf("failed to install kubectl: %w", err)
	}

	// Clean up downloaded files
	log.Debug("Cleaning up downloaded files...")
	cleanupCmd := exec.Command("rm", "-f", "kubectl", "kubectl.sha256")
	cleanupCmd.Run() // Ignore error as files might not exist

	log.Info("kubectl installed successfully")
	return nil
}

// SetupKubectlAlias adds the k=kubectl alias to bashrc
func SetupKubectlAlias(log *logger.Logger) error {
	log.Info("Adding kubectl alias 'k' to bashrc...")

	// Add alias to bashrc
	addCmd := exec.Command("bash", "-c", "echo 'alias k=kubectl' >> ~/.bashrc")
	if err := log.LogCommand(addCmd); err != nil {
		return fmt.Errorf("failed to add kubectl alias: %w", err)
	}

	log.Info("kubectl alias 'k' added successfully")
	return nil
}
