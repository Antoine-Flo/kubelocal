package tools

import (
	"fmt"
	"os/exec"
	"os/user"

	"github.com/Antoine-Flo/kubelocal/pkg/logger"
	"go.uber.org/zap"
)

// InstallDocker installs Docker
func InstallDocker(log *logger.Logger) error {
	log.Info("Downloading Docker installation script...")

	// Download the installation script
	downloadCmd := exec.Command("curl", "-fsSL", "https://get.docker.com", "-o", "get-docker.sh")
	if err := log.LogCommand(downloadCmd); err != nil {
		return fmt.Errorf("failed to download Docker script: %w", err)
	}

	log.Info("Executing Docker installation script...")
	// Execute the script
	installCmd := exec.Command("sudo", "sh", "get-docker.sh")
	if err := log.LogCommand(installCmd); err != nil {
		return fmt.Errorf("failed to install Docker: %w", err)
	}

	// Clean up
	log.Debug("Cleaning up Docker installation files...")
	cleanupCmd := exec.Command("rm", "-f", "get-docker.sh")
	cleanupCmd.Run()

	// Configure Docker permissions
	log.Info("Configuring Docker permissions...")
	if err := setupDockerPermissions(log); err != nil {
		log.Warn("Docker installed but permission setup failed", zap.Error(err))
		return err
	}

	log.Info("User added to docker group")

	// Activate Docker permissions
	log.Info("Activating Docker permissions...")
	activateCmd := exec.Command("bash", "-c", "newgrp docker -c 'echo \"Docker permissions activated successfully\"'")
	activateCmd.Run() // This might fail in some environments, so we don't check error

	log.Info("Docker installed successfully")
	return nil
}

// setupDockerPermissions configures Docker permissions for the current user
func setupDockerPermissions(log *logger.Logger) error {
	// Get current user
	currentUser, err := user.Current()
	if err != nil {
		return fmt.Errorf("failed to get current user: %w", err)
	}

	log.Debug("Adding user to docker group", zap.String("username", currentUser.Username))

	// Add user to docker group
	addUserCmd := exec.Command("sudo", "usermod", "-aG", "docker", currentUser.Username)
	if err := log.LogCommand(addUserCmd); err != nil {
		return fmt.Errorf("failed to add user to docker group: %w", err)
	}

	return nil
}
