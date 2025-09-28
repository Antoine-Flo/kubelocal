package tools

import (
	"fmt"
	"os/exec"

	"github.com/Antoine-Flo/kubelocal/pkg/logger"
)

// InstallKind installs Kind
func InstallKind(log *logger.Logger) error {
	log.Info("Downloading Kind (AMD64 Linux)...")

	// Download Kind (AMD64 Linux)
	downloadCmd := exec.Command("curl", "-Lo", "./kind", "https://kind.sigs.k8s.io/dl/v0.30.0/kind-linux-amd64")
	if err := log.LogCommand(downloadCmd); err != nil {
		return fmt.Errorf("failed to download Kind: %w", err)
	}

	log.Info("Making Kind executable...")
	// Make executable
	chmodCmd := exec.Command("chmod", "+x", "./kind")
	if err := log.LogCommand(chmodCmd); err != nil {
		return fmt.Errorf("failed to make Kind executable: %w", err)
	}

	log.Info("Installing Kind to /usr/local/bin...")
	// Move to /usr/local/bin
	installCmd := exec.Command("sudo", "mv", "./kind", "/usr/local/bin/kind")
	if err := log.LogCommand(installCmd); err != nil {
		return fmt.Errorf("failed to install Kind: %w", err)
	}

	log.Info("Kind installed successfully")
	return nil
}
