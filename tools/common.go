package tools

import (
	"fmt"
	"os/exec"
)

// SetupKubectlAlias adds the k=kubectl alias to bashrc
func SetupKubectlAlias() error {
	// Add alias to bashrc
	addCmd := exec.Command("bash", "-c", "echo 'alias k=kubectl' >> ~/.bashrc")
	if err := addCmd.Run(); err != nil {
		return fmt.Errorf("failed to add kubectl alias: %w", err)
	}

	return nil
}
