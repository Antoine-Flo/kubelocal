package main

import (
	"fmt"
	"os/exec"
)

// installDocker installe Docker
func installDocker() error {

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

// installKind installe Kind
func installKind() error {
	// Télécharger Kind (AMD64 Linux)
	if err := exec.Command("curl", "-Lo", "./kind", "https://kind.sigs.k8s.io/dl/v0.30.0/kind-linux-amd64").Run(); err != nil {
		return fmt.Errorf("failed to download Kind: %w", err)
	}

	// Rendre exécutable
	if err := exec.Command("chmod", "+x", "./kind").Run(); err != nil {
		return fmt.Errorf("failed to make Kind executable: %w", err)
	}

	// Déplacer vers /usr/local/bin
	if err := exec.Command("sudo", "mv", "./kind", "/usr/local/bin/kind").Run(); err != nil {
		return fmt.Errorf("failed to install Kind: %w", err)
	}

	fmt.Println("✅ Kind installed successfully")
	return nil
}

// installMinikube installe Minikube
func installMinikube() error {
	// Télécharger Minikube
	if err := exec.Command("curl", "-LO", "https://github.com/kubernetes/minikube/releases/latest/download/minikube-linux-amd64").Run(); err != nil {
		return fmt.Errorf("failed to download Minikube: %w", err)
	}

	// Installer Minikube
	if err := exec.Command("sudo", "install", "minikube-linux-amd64", "/usr/local/bin/minikube").Run(); err != nil {
		return fmt.Errorf("failed to install Minikube: %w", err)
	}

	// Nettoyer
	exec.Command("rm", "-f", "minikube-linux-amd64").Run()

	fmt.Println("✅ Minikube installed successfully")
	return nil
}

// setupKubectlAlias ajoute l'alias k=kubectl au bashrc
func setupKubectlAlias() error {
	// Ajouter l'alias au bashrc
	addCmd := exec.Command("bash", "-c", "echo 'alias k=kubectl' >> ~/.bashrc")
	if err := addCmd.Run(); err != nil {
		return fmt.Errorf("failed to add kubectl alias: %w", err)
	}

	return nil
}
