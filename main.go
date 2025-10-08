package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/Antoine-Flo/kubelocal/internal/logger"
	"github.com/Antoine-Flo/kubelocal/tools"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/huh/spinner"
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
	log.Info("Installing component", zap.String("component", component))

	var installErr error
	if err := spinner.New().
		Title(fmt.Sprintf("Installing %s...", component)).
		Action(func() {
			installErr = installFunc(log)
		}).
		Run(); err != nil {
		handleInstallationError(log, component, err)
		return
	}

	if installErr != nil {
		handleInstallationError(log, component, installErr)
	}
}

// handleInstallationError logs the error and shows a simple message to the user
func handleInstallationError(log *logger.Logger, component string, err error) {
	log.Error("Installation failed",
		zap.String("component", component),
		zap.Error(err))
	fmt.Printf("❌ Something went wrong during %s installation.\n", component)
	fmt.Println("   Find the full log in ~/.local/share/kubelocal/logs/install.log")
	os.Exit(1)
}

func runSetup(log *logger.Logger) {
	var (
		containerRuntime string
		kubernetesLocal  string
		cliTools         []string
	)

	log.Info("Starting kubelocal installation")

	// Kubernetes solution selection
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Choose your local Kubernetes solution:").
				Options(
					huh.NewOption("Kind", "kind"),
					huh.NewOption("Minikube", "minikube"),
					huh.NewOption("K3s", "k3s"),
				).
				Value(&kubernetesLocal),
		),
		huh.NewGroup(
			huh.NewMultiSelect[string]().
				Title("Choose your manifest management tools:").
				Options(
					huh.NewOption("Helm", "helm"),
					huh.NewOption("Kustomize", "kustomize"),
				).
				Value(&cliTools),
		),
		huh.NewGroup(
			huh.NewMultiSelect[string]().
				Title("Choose your service mesh:").
				Options(
					huh.NewOption("Istio", "istio"),
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

	// Determine if Docker is needed
	needsDocker := kubernetesLocal == "kind" || kubernetesLocal == "minikube"
	if needsDocker {
		containerRuntime = "docker"
	}

	// Log user selections
	log.Debug("User selections",
		zap.String("kubernetes_local", kubernetesLocal),
		zap.Bool("needs_docker", needsDocker),
		zap.Strings("cli_tools", cliTools))

	// Display choices and simulate installation
	fmt.Println("\n=== Configuration Summary ===")
	fmt.Printf("Kubernetes Local: %s\n", kubernetesLocal)
	if needsDocker {
		fmt.Printf("Container Runtime: Docker (required for %s)\n", kubernetesLocal)
	} else {
		fmt.Println("Container Runtime: None (K3s includes everything)")
	}
	fmt.Printf("Manifest Management: kubectl")
	manifestTools := []string{}
	serviceMesh := []string{}
	for _, tool := range cliTools {
		if tool == "istio" {
			serviceMesh = append(serviceMesh, tool)
		} else {
			manifestTools = append(manifestTools, tool)
		}
	}
	for _, tool := range manifestTools {
		fmt.Printf(", %s", tool)
	}
	if len(serviceMesh) > 0 {
		fmt.Printf("\nService Mesh: %s", strings.Join(serviceMesh, ", "))
	}
	fmt.Println()

	// Real installation
	fmt.Println("\n=== Installation in progress ===")
	log.Info("Starting installation process")

	// Install kubectl first (required for Kubernetes solutions)
	installComponent(log, "kubectl", tools.InstallKubectl)

	// Configure kubectl alias in background
	log.Info("Setting up aliases")
	if err := tools.SetupAliases(log); err != nil {
		// kubectl alias setup is not critical, just warn the user
		log.Warn("kubectl alias setup failed", zap.Error(err))
		fmt.Println("⚠️  Warning: kubectl alias setup failed.")
		fmt.Println("   You can manually add 'alias k=kubectl' to your ~/.bashrc")
	}

	// Install jq in background
	log.Info("Installing jq")
	if err := tools.InstallJq(log); err != nil {
		log.Warn("jq installation failed", zap.Error(err))
		fmt.Println("⚠️  Warning: jq installation failed.")
		fmt.Println("   You can manually install jq later with: sudo apt install jq")
	}
	// Install container runtime only if needed
	if needsDocker {
		installComponent(log, "Docker", tools.InstallDocker)
	}

	// Install local Kubernetes solution
	switch kubernetesLocal {
	case "kind":
		installComponent(log, "Kind", tools.InstallKind)
	case "minikube":
		installComponent(log, "Minikube", tools.InstallMinikube)
	case "k3s":
		installComponent(log, "K3s", tools.InstallK3s)
	}

	// Install CLI tools
	for _, tool := range cliTools {
		switch tool {
		case "helm":
			installComponent(log, "Helm", tools.InstallHelm)
		case "kustomize":
			installComponent(log, "Kustomize", tools.InstallKustomize)
		case "istio":
			installComponent(log, "Istio", tools.InstallIstio)
		}
	}

	// Install Istio on cluster if selected
	if contains(cliTools, "istio") {
		fmt.Println("\n=== Installing Istio on cluster ===")
		log.Info("Installing Istio on cluster")
		if err := tools.InstallIstioOnCluster(log, kubernetesLocal); err != nil {
			log.Warn("Istio cluster installation failed", zap.Error(err))
			fmt.Println("⚠️  Warning: Istio cluster installation failed.")
			fmt.Println("   You can manually install Istio later with: istioctl install")
		} else {
			fmt.Println("✅ Istio successfully installed on cluster")
		}
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

	case "k3s":
		fmt.Println("   1. Start your cluster:")
		fmt.Printf("      sudo systemctl start k3s\n")
		fmt.Println("\n   2. Verify your cluster:")
		fmt.Printf("      sudo k3s kubectl get nodes\n")
		fmt.Println("\n   3. Set up kubectl access:")
		fmt.Printf("      mkdir -p ~/.kube\n")
		fmt.Printf("      sudo cp /etc/rancher/k3s/k3s.yaml ~/.kube/config\n")
		fmt.Printf("      sudo chown $USER ~/.kube/config\n")
		fmt.Println("\n   4. When you're done:")
		fmt.Printf("      sudo systemctl stop k3s\n")
	}

	// Instructions for container runtime
	needsDocker := kubernetesLocal == "kind" || kubernetesLocal == "minikube"
	if needsDocker {
		fmt.Println("\n💡 Docker Tips:")
		fmt.Println("   • Check Docker status: docker --version")
		fmt.Println("   • Test Docker access: docker ps")
		fmt.Println("   • If permission denied, run: newgrp docker")
		fmt.Println("   • Or logout/login to apply docker group changes")
	} else {
		fmt.Println("\n💡 K3s Tips:")
		fmt.Println("   • Check K3s status: sudo systemctl status k3s")
		fmt.Println("   • View K3s logs: sudo journalctl -u k3s -f")
		fmt.Println("   • K3s includes everything: containerd, CNI, etc.")
	}

	// Instructions for manifest management tools
	manifestTools := []string{}
	serviceMesh := []string{}
	for _, tool := range cliTools {
		if tool == "istio" {
			serviceMesh = append(serviceMesh, tool)
		} else {
			manifestTools = append(manifestTools, tool)
		}
	}

	if len(manifestTools) > 0 {
		fmt.Println("\n🔧 Manifest Management Tools:")
		for _, tool := range manifestTools {
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

	// Instructions for service mesh
	if len(serviceMesh) > 0 {
		fmt.Println("\n🌐 Service Mesh:")
		for _, tool := range serviceMesh {
			switch tool {
			case "istio":
				fmt.Println("   • Istio Service Mesh:")
				istioInstructions := tools.GetIstioInstructions(kubernetesLocal)
				for _, instruction := range istioInstructions {
					fmt.Println(instruction)
				}
			}
		}
	}

	fmt.Println("\n✨ Happy Kubernetes development! ✨")
}

// contains checks if a string slice contains a specific string
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
