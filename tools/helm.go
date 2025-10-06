package tools

import (
	"fmt"

	"github.com/Antoine-Flo/kubelocal/internal/command"
	"github.com/Antoine-Flo/kubelocal/internal/logger"
)

func InstallHelm(log *logger.Logger) error {
	log.Info("Installing Helm using official script...")

	if err := command.Run(log, "bash", "-c", "curl https://raw.githubusercontent.com/helm/helm/main/scripts/get-helm-3 | bash"); err != nil {
		return fmt.Errorf("failed to install Helm: %w", err)
	}

	log.Info("Helm installed successfully")
	return nil
}
