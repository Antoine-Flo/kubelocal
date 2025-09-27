package main

import (
	"fmt"
	"os/exec"
)

// installDocker installe Docker
func installDocker() error {
	fmt.Println("🐋 Installing Docker...")

	// Télécharger le script d'installation
	if err := exec.Command("curl", "-fsSL", "https://get.docker.com", "-o", "get-docker.sh").Run(); err != nil {
		return fmt.Errorf("failed to download Docker script: %w", err)
	}

	// Exécuter le script
	if err := exec.Command("sudo", "sh", "get-docker.sh").Run(); err != nil {
		return fmt.Errorf("failed to install Docker: %w", err)
	}

	// Nettoyer
	exec.Command("rm", "-f", "get-docker.sh").Run()

	fmt.Println("✅ Docker installed successfully")
	return nil
}

// installPodman installe Podman
func installPodman() error {
	fmt.Println("🚀 Installing Podman...")

	// Mettre à jour les paquets
	if err := exec.Command("sudo", "apt-get", "update").Run(); err != nil {
		return fmt.Errorf("failed to update packages: %w", err)
	}

	// Installer Podman
	if err := exec.Command("sudo", "apt-get", "-y", "install", "podman").Run(); err != nil {
		return fmt.Errorf("failed to install Podman: %w", err)
	}

	fmt.Println("✅ Podman installed successfully")
	return nil
}
