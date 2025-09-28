package tools

import (
	"fmt"
	"os/exec"
)

// InstallKind installs Kind
func InstallKind() error {
	// Download Kind (AMD64 Linux)
	if err := exec.Command("curl", "-Lo", "./kind", "https://kind.sigs.k8s.io/dl/v0.30.0/kind-linux-amd64").Run(); err != nil {
		return fmt.Errorf("failed to download Kind: %w", err)
	}

	// Make executable
	if err := exec.Command("chmod", "+x", "./kind").Run(); err != nil {
		return fmt.Errorf("failed to make Kind executable: %w", err)
	}

	// Move to /usr/local/bin
	if err := exec.Command("sudo", "mv", "./kind", "/usr/local/bin/kind").Run(); err != nil {
		return fmt.Errorf("failed to install Kind: %w", err)
	}

	fmt.Println("✅ Kind installed successfully")
	return nil
}
