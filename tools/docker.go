package tools

import (
	"fmt"
	"os/exec"
	"os/user"
)

// InstallDocker installs Docker
func InstallDocker() error {
	// Download the installation script
	if err := exec.Command("curl", "-fsSL", "https://get.docker.com", "-o", "get-docker.sh").Run(); err != nil {
		return fmt.Errorf("failed to download Docker script: %w", err)
	}

	// Execute the script
	if err := exec.Command("sudo", "sh", "get-docker.sh").Run(); err != nil {
		return fmt.Errorf("failed to install Docker: %w", err)
	}

	// Clean up
	exec.Command("rm", "-f", "get-docker.sh").Run()

	// Configure Docker permissions
	if err := setupDockerPermissions(); err != nil {
		fmt.Printf("⚠️  Docker installed but permission setup failed: %v\n", err)
		fmt.Println("💡 You may need to run: sudo usermod -aG docker $USER && newgrp docker")
	}

	fmt.Println("✅ Docker installed successfully")
	return nil
}

// setupDockerPermissions configures Docker permissions for the current user
func setupDockerPermissions() error {
	// Get current user
	currentUser, err := user.Current()
	if err != nil {
		return fmt.Errorf("failed to get current user: %w", err)
	}

	// Add user to docker group
	if err := exec.Command("sudo", "usermod", "-aG", "docker", currentUser.Username).Run(); err != nil {
		return fmt.Errorf("failed to add user to docker group: %w", err)
	}

	return nil
}
