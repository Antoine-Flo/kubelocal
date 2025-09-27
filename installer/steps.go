package installer

import (
	"fmt"

	"github.com/Antoine-Flo/kubelocal/types"
)

// getInstallSteps returns the installation steps for given runtime and distribution
func getInstallSteps(runtime, distribution string) []types.InstallStep {
	return []types.InstallStep{
		{Name: "Updating package lists"},
		{Name: "Installing prerequisites"},
		{Name: fmt.Sprintf("Installing %s", runtime)},
		{Name: "Installing kubectl (latest version)"},
		{Name: fmt.Sprintf("Installing %s", distribution)},
		{Name: "Configuring cluster"},
		{Name: "Setting up local storage"},
	}
}

// ExecuteStep simulates execution of an installation step
func ExecuteStep(step types.InstallStep) error {
	// Mock execution - in real implementation, this would contain actual logic
	return nil
}
