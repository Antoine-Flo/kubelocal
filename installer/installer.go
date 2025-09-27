package installer

import (
	"fmt"

	"github.com/Antoine-Flo/kubelocal/types"
	"github.com/Antoine-Flo/kubelocal/ui"
)

// Run executes the complete installation process
func Run(config types.Config) error {
	printInstallationStart()

	if err := executeSteps(config.Steps); err != nil {
		return fmt.Errorf("installation failed: %w", err)
	}

	printInstallationComplete(config)
	return nil
}

// printInstallationStart displays installation start message
func printInstallationStart() {
	fmt.Printf("%s🔧 Starting installation process...%s\n", ui.ColorYellow, ui.ColorReset)
	fmt.Println()
}

// executeSteps runs all installation steps
func executeSteps(steps []types.InstallStep) error {
	for i, step := range steps {
		ui.PrintStep(i+1, len(steps), step.Name)
		if err := ExecuteStep(step); err != nil {
			return err
		}
		ui.PrintSuccess("Done")
	}
	return nil
}

// printInstallationComplete displays completion message
func printInstallationComplete(config types.Config) {
	fmt.Println()
	fmt.Printf("%s🎉 Installation completed successfully!%s\n", ui.ColorGreen, ui.ColorReset)
	fmt.Printf("Your local Kubernetes environment is ready with %s and %s\n",
		config.Runtime.Name, config.Distribution.Name)
}
