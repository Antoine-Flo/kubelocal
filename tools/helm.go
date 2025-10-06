package tools

import (
	"fmt"
	"os/exec"

	"github.com/Antoine-Flo/kubelocal/internal/logger"
)

// InstallHelm installs Helm using the official installation script
func InstallHelm(log *logger.Logger) error {
	log.Info("Installing Helm using official script...")

	// Download and execute Helm installation script
	installCmd := exec.Command("bash", "-c", "curl https://raw.githubusercontent.com/helm/helm/main/scripts/get-helm-3 | bash")
	if err := log.LogCommand(installCmd); err != nil {
		return fmt.Errorf("failed to install Helm: %w", err)
	}

	log.Info("Helm installed successfully")
	return nil
}
