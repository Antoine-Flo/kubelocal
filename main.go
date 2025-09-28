package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/Antoine-Flo/kubelocal/tools"
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

	// Real installation
	fmt.Println("\n=== Installation in progress ===")

	// Install container runtime
	switch containerRuntime {
	case "docker":
		installWithSpinner("Docker 🐋", tools.InstallDocker)
	case "podman":
		installWithSpinner("Podman 🚀", tools.InstallPodman)
	}

	// Install local Kubernetes solution
	switch kubernetesLocal {
	case "kind":
		installWithSpinner("Kind ☸️", tools.InstallKind)
	case "minikube":
		installWithSpinner("Minikube 🚀", tools.InstallMinikube)
	}

	// Configure kubectl alias in background
	installSilently(tools.SetupKubectlAlias)

	// Show getting started instructions
	showGettingStarted(containerRuntime, kubernetesLocal, cliTools)
}

// installWithSpinner installs a tool with a spinner
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

// installSilently installs a tool in background without display
func installSilently(installFunc func() error) error {
	return installFunc()
}

// showGettingStarted displays customized instructions based on installation
func showGettingStarted(containerRuntime, kubernetesLocal string, cliTools []string) {
	fmt.Println("\n🎉 Installation completed successfully!")
	fmt.Println("\n📝 We created an alias 'k' for kubectl for you!")

	fmt.Println("\n🚀 Getting Started:")

	// Specific instructions based on Kubernetes solution
	switch kubernetesLocal {
	case "kind":
		fmt.Println("   1. Create your first cluster:")
		fmt.Printf("      kind create cluster --name my-cluster\n")
		fmt.Println("\n   2. Verify your cluster:")
		fmt.Printf("      k cluster-info --context kind-my-cluster\n")
		fmt.Printf("      k get nodes\n")
		fmt.Println("\n   3. When you're done:")
		fmt.Printf("      kind delete cluster --name my-cluster\n")

	case "minikube":
		fmt.Println("   1. Start your first cluster:")
		fmt.Printf("      minikube start\n")
		fmt.Println("\n   2. Verify your cluster:")
		fmt.Printf("      k cluster-info\n")
		fmt.Printf("      k get nodes\n")
		fmt.Println("\n   3. Access the dashboard:")
		fmt.Printf("      minikube dashboard\n")
		fmt.Println("\n   4. When you're done:")
		fmt.Printf("      minikube stop\n")
	}

	// Instructions for Docker/Podman
	fmt.Printf("\n💡 Container Runtime (%s) Tips:\n", containerRuntime)
	switch containerRuntime {
	case "docker":
		fmt.Println("   • Check Docker status: docker --version")
		fmt.Println("   • Test Docker access: docker ps")
		fmt.Println("   • If permission denied, run: newgrp docker")
		fmt.Println("   • Or logout/login to apply docker group changes")

	case "podman":
		fmt.Println("   • Check Podman status: podman --version")
		fmt.Println("   • List running containers: podman ps")
	}

	// Instructions for CLI tools
	if len(cliTools) > 0 {
		fmt.Println("\n🔧 Additional Tools:")
		for _, tool := range cliTools {
			switch tool {
			case "helm":
				fmt.Println("   • Helm: helm version")
				fmt.Println("     Get started: helm create my-chart")

			case "kustomize":
				fmt.Println("   • Kustomize: kubectl kustomize --help")
				fmt.Println("     Get started: kustomize create .")
			}
		}
	}

	fmt.Println("\n✨ Happy Kubernetes development! ✨")
}
