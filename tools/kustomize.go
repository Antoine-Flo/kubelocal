package tools

import (
	"fmt"

	"github.com/Antoine-Flo/kubelocal/internal/command"
	"github.com/Antoine-Flo/kubelocal/internal/logger"
)

func InstallKustomize(log *logger.Logger) error {
	log.Info("Installing Kustomize using official script...")

	if err := command.Run(log, "bash", "-c", "curl -s \"https://raw.githubusercontent.com/kubernetes-sigs/kustomize/master/hack/install_kustomize.sh\" | bash"); err != nil {
		return fmt.Errorf("failed to install Kustomize: %w", err)
	}

	log.Info("Kustomize installed successfully")
	return nil
}
