package tools

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/Antoine-Flo/kubelocal/internal/command"
	"github.com/Antoine-Flo/kubelocal/internal/logger"
	"go.uber.org/zap"
)

func InstallKubectl(log *logger.Logger) error {
	log.Info("Getting latest kubectl version...")

	getVersionCmd := exec.Command("curl", "-L", "-s", "https://dl.k8s.io/release/stable.txt")
	versionBytes, err := getVersionCmd.Output()
	if err != nil {
		return fmt.Errorf("failed to get kubectl version: %w", err)
	}
	version := strings.TrimSpace(string(versionBytes))
	log.Info("Latest kubectl version", zap.String("version", version))

	log.Info("Downloading kubectl binary...")
	downloadURL := fmt.Sprintf("https://dl.k8s.io/release/%s/bin/linux/amd64/kubectl", version)
	if err := command.Run(log, "curl", "-LO", downloadURL); err != nil {
		return fmt.Errorf("failed to download kubectl: %w", err)
	}

	log.Info("Installing kubectl to /usr/local/bin...")
	if err := command.Run(log, "sudo", "install", "-o", "root", "-g", "root", "-m", "0755", "kubectl", "/usr/local/bin/kubectl"); err != nil {
		return fmt.Errorf("failed to install kubectl: %w", err)
	}

	log.Debug("Cleaning up downloaded files...")
	command.Run(log, "rm", "-f", "kubectl") // Ignore error as files might not exist

	log.Info("kubectl installed successfully")
	return nil
}

func SetupKubectlAlias(log *logger.Logger) error {
	log.Info("Adding kubectl alias 'k' to bashrc...")

	if err := command.Run(log, "bash", "-c", "echo 'alias k=kubectl' >> ~/.bashrc"); err != nil {
		return fmt.Errorf("failed to add kubectl alias: %w", err)
	}

	log.Info("Sourcing bashrc to activate alias...")
	if err := command.Run(log, "bash", "-c", "source ~/.bashrc"); err != nil {
		log.Warn("Failed to source bashrc, alias will be available after shell restart", zap.Error(err))
	}

	log.Info("kubectl alias 'k' added and activated successfully")
	return nil
}
