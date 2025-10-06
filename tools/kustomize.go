package tools

import (
	"fmt"
	"os/exec"

	"github.com/Antoine-Flo/kubelocal/internal/logger"
)

// InstallKustomize installs Kustomize using the official installation script
func InstallKustomize(log *logger.Logger) error {
	log.Info("Installing Kustomize using official script...")

	// Download and execute Kustomize installation script
	installCmd := exec.Command("bash", "-c", "curl -s \"https://raw.githubusercontent.com/kubernetes-sigs/kustomize/master/hack/install_kustomize.sh\" | bash")
	if err := log.LogCommand(installCmd); err != nil {
		return fmt.Errorf("failed to install Kustomize: %w", err)
	}

	log.Info("Kustomize installed successfully")
	return nil
}
