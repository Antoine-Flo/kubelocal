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

func InstallJq(log *logger.Logger) error {
	log.Info("Installing jq...")

	if err := command.Run(log, "sudo", "apt", "update"); err != nil {
		return fmt.Errorf("failed to update package list: %w", err)
	}

	if err := command.Run(log, "sudo", "apt", "install", "-y", "jq"); err != nil {
		return fmt.Errorf("failed to install jq: %w", err)
	}

	log.Info("jq installed successfully")
	return nil
}

func SetupAliases(log *logger.Logger) error {
	log.Info("Setting up aliases...")

	aliases := []string{
		"alias k=kubectl",
		"alias logs=\"cat ~/.local/share/kubelocal/logs/install.log | jq\"",
	}

	for _, alias := range aliases {
		if err := command.Run(log, "bash", "-c", fmt.Sprintf("echo '%s' >> ~/.bashrc", alias)); err != nil {
			return fmt.Errorf("failed to add alias: %w", err)
		}
	}

	log.Info("Sourcing bashrc to activate aliases...")
	if err := command.Run(log, "bash", "-c", "source ~/.bashrc"); err != nil {
		log.Warn("Failed to source bashrc, aliases will be available after shell restart", zap.Error(err))
	}

	log.Info("All aliases added and activated successfully")
	return nil
}
