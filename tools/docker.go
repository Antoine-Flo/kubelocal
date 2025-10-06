package tools

import (
	"fmt"
	"os/exec"
	"os/user"
	"strings"

	"github.com/Antoine-Flo/kubelocal/internal/logger"
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

	log.Info("Starting Docker permissions setup", zap.String("username", currentUser.Username))

	// Check if user is already in docker group
	log.Debug("Checking if user is already in docker group")
	if isUserInDockerGroup(currentUser.Username) {
		log.Info("User is already in docker group")
	} else {
		log.Debug("User not in docker group, adding to docker group")
		addUserCmd := exec.Command("sudo", "usermod", "-aG", "docker", currentUser.Username)
		if err := log.LogCommand(addUserCmd); err != nil {
			log.Error("Failed to add user to docker group", zap.Error(err))
			return fmt.Errorf("failed to add user to docker group: %w", err)
		}
		log.Info("User successfully added to docker group")
	}

	// Test Docker access before trying solutions
	log.Debug("Testing Docker access before applying solutions")
	if canAccessDocker() {
		log.Info("Docker is already accessible, no additional setup needed")
		return nil
	}

	// Solution 1: Try using sg to activate docker group
	log.Info("Solution 1: Attempting to activate docker group with sg")
	sgCmd := exec.Command("sg", "docker", "-c", "docker version")
	if err := log.LogCommand(sgCmd); err != nil {
		log.Warn("Solution 1 failed: sg command failed", zap.Error(err))
	} else {
		log.Info("Solution 1 successful: Docker accessible via sg")
		// Test if we can still access docker after sg
		if canAccessDocker() {
			log.Info("Docker access confirmed after sg activation")
			return nil
		}
	}

	// Solution 2: Try changing docker.sock permissions
	log.Info("Solution 2: Attempting to change docker.sock permissions")
	sockCmd := exec.Command("sudo", "chmod", "666", "/var/run/docker.sock")
	if err := log.LogCommand(sockCmd); err != nil {
		log.Warn("Solution 2 failed: Could not change docker.sock permissions", zap.Error(err))
	} else {
		log.Info("Solution 2 applied: docker.sock permissions changed to 666")
		// Test Docker access after permission change
		if canAccessDocker() {
			log.Info("Solution 2 successful: Docker accessible after permission change")
			return nil
		} else {
			log.Warn("Solution 2 applied but Docker still not accessible")
		}
	}

	// Final test
	log.Debug("Running final Docker access test")
	if canAccessDocker() {
		log.Info("Docker permissions setup completed successfully")
		return nil
	}

	log.Warn("All solutions attempted but Docker still not accessible")
	return fmt.Errorf("failed to make Docker accessible after trying all solutions")
}

// isUserInDockerGroup checks if the user is already in the docker group
func isUserInDockerGroup(username string) bool {
	cmd := exec.Command("groups", username)
	output, err := cmd.Output()
	if err != nil {
		return false
	}
	return strings.Contains(string(output), "docker")
}

// canAccessDocker tests if Docker is accessible without sudo
func canAccessDocker() bool {
	cmd := exec.Command("docker", "version", "--format", "{{.Server.Version}}")
	err := cmd.Run()
	return err == nil
}
