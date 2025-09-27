package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/huh/spinner"
)

var (
	version = "dev"
)

func main() {
	showVersion := flag.Bool("version", false, "Display version")
	flag.Parse()

	if *showVersion {
		fmt.Printf("Version: %s\n", version)
		os.Exit(0)
	}

	showWelcome()
	runSetup()
}

func showWelcome() {
	fmt.Println(`
 ██╗  ██╗██╗   ██╗██████╗ ███████╗██╗      ██████╗  ██████╗ █████╗ ██╗     
 ██║ ██╔╝██║   ██║██╔══██╗██╔════╝██║     ██╔═══██╗██╔════╝██╔══██╗██║     
 █████╔╝ ██║   ██║██████╔╝█████╗  ██║     ██║   ██║██║     ███████║██║     
 ██╔═██╗ ██║   ██║██╔══██╗██╔══╝  ██║     ██║   ██║██║     ██╔══██║██║     
 ██║  ██╗╚██████╔╝██████╔╝███████╗███████╗╚██████╔╝╚██████╗██║  ██║███████╗
 ╚═╝  ╚═╝ ╚═════╝ ╚═════╝ ╚══════╝╚══════╝ ╚═════╝  ╚═════╝╚═╝  ╚═╝╚══════╝
                                                                            
          🚀 Quick Kubernetes Local Development Setup Tool 🚀
`)
	fmt.Println("Welcome! Let's set up your local Kubernetes environment.\n")
}

func runSetup() {
	var (
		containerRuntime string
		kubernetesLocal  string
		cliTools         []string
	)

	// Container runtime selection
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Choose your container runtime:").
				Options(
					huh.NewOption("Docker", "docker"),
					huh.NewOption("Podman", "podman"),
				).
				Value(&containerRuntime),
		),
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Choose your local Kubernetes solution:").
				Options(
					huh.NewOption("Kind", "kind"),
					huh.NewOption("Minikube", "minikube"),
				).
				Value(&kubernetesLocal),
		),
		huh.NewGroup(
			huh.NewMultiSelect[string]().
				Title("Choose your deployment tools:").
				Options(
					huh.NewOption("Helm", "helm"),
					huh.NewOption("Kustomize", "kustomize"),
				).
				Value(&cliTools),
		),
	).WithTheme(huh.ThemeBase16())

	err := form.Run()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	// Display choices and simulate installation
	fmt.Println("\n=== Configuration Summary ===")
	fmt.Printf("Container Runtime: %s\n", containerRuntime)
	fmt.Printf("Kubernetes Local: %s\n", kubernetesLocal)
	fmt.Printf("Deployment Tools: kubectl")
	for _, tool := range cliTools {
		fmt.Printf(", %s", tool)
	}
	fmt.Println()

	// Installation réelle
	fmt.Println("\n=== Installation in progress ===")

	// Installer le container runtime
	switch containerRuntime {
	case "docker":
		installWithSpinner("Docker", installDocker)
	case "podman":
		installWithSpinner("Podman", installPodman)
	}

	fmt.Println("\n✅ Installation completed!")
}

// installWithSpinner installe un outil avec un spinner
func installWithSpinner(name string, installFunc func() error) {
	var installErr error

	err := spinner.New().
		Title(fmt.Sprintf("Installing %s...", name)).
		Action(func() {
			installErr = installFunc()
		}).
		Run()

	if err != nil {
		fmt.Printf("❌ %s spinner error: %v\n", name, err)
	} else if installErr != nil {
		fmt.Printf("❌ %s installation failed: %v\n", name, installErr)
	}
}
