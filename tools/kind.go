package tools

import (
	"fmt"

	"github.com/Antoine-Flo/kubelocal/internal/command"
	"github.com/Antoine-Flo/kubelocal/internal/logger"
)

func InstallKind(log *logger.Logger) error {
	log.Info("Downloading Kind (AMD64 Linux)...")

	if err := command.Run(log, "curl", "-Lo", "./kind", "https://kind.sigs.k8s.io/dl/v0.30.0/kind-linux-amd64"); err != nil {
		return fmt.Errorf("failed to download Kind: %w", err)
	}

	log.Info("Making Kind executable...")
	if err := command.Run(log, "chmod", "+x", "./kind"); err != nil {
		return fmt.Errorf("failed to make Kind executable: %w", err)
	}

	log.Info("Installing Kind to /usr/local/bin...")
	if err := command.Run(log, "sudo", "mv", "./kind", "/usr/local/bin/kind"); err != nil {
		return fmt.Errorf("failed to install Kind: %w", err)
	}

	log.Info("Kind installed successfully")
	return nil
}
