package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/Antoine-Flo/kubelocal/pkg/logger"
	"github.com/Antoine-Flo/kubelocal/tools"
	"github.com/charmbracelet/huh"
	"go.uber.org/zap"
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

	// Initialize logger
	log, err := logger.NewLogger()
	if err != nil {
		fmt.Println("❌ Failed to initialize logging system.")
		fmt.Printf("   Error: %v\n", err)
		os.Exit(1)
	}
	defer log.Close()

	showWelcome()
	runSetup(log)
}

func showWelcome() {
	fmt.Print(`
╦╔═╦ ╦╔╗ ╔═╗╦  ╔═╗╔═╗╔═╗╦  
╠╩╗║ ║╠╩╗║╣ ║  ║ ║║  ╠═╣║  
╩ ╩╚═╝╚═╝╚═╝╩═╝╚═╝╚═╝╩ ╩╩═╝
Kubernetes Local Development Setup 🚀

`)
	fmt.Print("Welcome! Let's set up your local Kubernetes environment.\n")
}

// installComponent handles the common pattern of printing progress, logging, and error handling
func installComponent(log *logger.Logger, component string, installFunc func(*logger.Logger) error) {
	fmt.Printf("Installing %s...\n", component)
	log.Info("Installing component", zap.String("component", component))
	if err := installFunc(log); err != nil {
		handleInstallationError(log, component, err)
	}
}

// handleInstallationError logs the error and shows a simple message to the user
func handleInstallationError(log *logger.Logger, component string, err error) {
	log.Error("Installation failed",
		zap.String("component", component),
		zap.Error(err))
	fmt.Printf("❌ Something went wrong during %s installation.\n", component)
	fmt.Println("   Find the full log in ~/.local/share/kubelocal/logs/install.log or /var/log/kubelocal/install.log")
	os.Exit(1)
}

func runSetup(log *logger.Logger) {
	var (
		containerRuntime string
		kubernetesLocal  string
		cliTools         []string
	)

	log.Info("Starting kubelocal installation")

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
		log.Error("Form execution failed", zap.Error(err))
		fmt.Println("❌ Something went wrong during setup.")
		fmt.Println("   Find the full log in ~/.local/share/kubelocal/logs/install.log or /var/log/kubelocal/install.log")
		os.Exit(1)
	}

	// Log user selections
	log.Debug("User selections",
		zap.String("container_runtime", containerRuntime),
		zap.String("kubernetes_local", kubernetesLocal),
		zap.Strings("cli_tools", cliTools))

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
	log.Info("Starting installation process")

	// Install kubectl first (required for Kubernetes solutions)
	installComponent(log, "kubectl", tools.InstallKubectl)

	// Install container runtime
	switch containerRuntime {
	case "docker":
		installComponent(log, "Docker", tools.InstallDocker)
	case "podman":
		installComponent(log, "Podman", tools.InstallPodman)
	}

	// Install local Kubernetes solution
	switch kubernetesLocal {
	case "kind":
		installComponent(log, "Kind", tools.InstallKind)
	case "minikube":
		installComponent(log, "Minikube", tools.InstallMinikube)
	}

	// Install CLI tools
	for _, tool := range cliTools {
		switch tool {
		case "helm":
			installComponent(log, "Helm", tools.InstallHelm)
		case "kustomize":
			installComponent(log, "Kustomize", tools.InstallKustomize)
		}
	}

	// Configure kubectl alias in background
	log.Info("Setting up kubectl alias")
	if err := tools.SetupKubectlAlias(log); err != nil {
		// kubectl alias setup is not critical, just warn the user
		log.Warn("kubectl alias setup failed", zap.Error(err))
		fmt.Println("⚠️  Warning: kubectl alias setup failed.")
		fmt.Println("   You can manually add 'alias k=kubectl' to your ~/.bashrc")
	}

	log.Info("Installation completed successfully")

	// Show getting started instructions
	showGettingStarted(containerRuntime, kubernetesLocal, cliTools)
}

// showGettingStarted displays customized instructions based on installation
func showGettingStarted(containerRuntime, kubernetesLocal string, cliTools []string) {
	fmt.Println("\n🎉 Installation completed successfully!")
	fmt.Println("\n📝 We installed kubectl and created an alias 'k' for you!")

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
