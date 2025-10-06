package tools

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/Antoine-Flo/kubelocal/internal/command"
	"github.com/Antoine-Flo/kubelocal/internal/logger"
	"go.uber.org/zap"
)

// InstallIstio installs Istio by downloading istioctl and setting up the environment
func InstallIstio(log *logger.Logger) error {
	log.Info("Installing Istio...")

	// Download and install istioctl
	if err := installIstioctl(log); err != nil {
		return fmt.Errorf("failed to install istioctl: %w", err)
	}

	// Setup istioctl in PATH
	if err := setupIstioctlPath(log); err != nil {
		return fmt.Errorf("failed to setup istioctl PATH: %w", err)
	}

	log.Info("Istio installed successfully")
	return nil
}

// installIstioctl downloads and installs the istioctl binary
func installIstioctl(log *logger.Logger) error {
	log.Info("Downloading Istio installation script...")

	// Check current directory before download
	log.Info("Current directory before download:")
	command.Run(log, "pwd")
	command.Run(log, "ls", "-la")

	// Download the Istio installation script
	log.Info("Running curl command...")
	if err := command.Run(log, "curl", "-L", "-o", "downloadIstio.sh", "https://istio.io/downloadIstio"); err != nil {
		return fmt.Errorf("failed to download Istio script: %w", err)
	}

	// Check if script was downloaded
	log.Info("Checking if script was downloaded...")
	command.Run(log, "ls", "-la", "downloadIstio.sh")
	command.Run(log, "file", "downloadIstio.sh")

	// Show first few lines of the script
	log.Info("First 20 lines of downloaded script:")
	command.Run(log, "head", "-20", "downloadIstio.sh")

	// Execute the installation script
	log.Info("Executing Istio installation script...")
	if err := command.Run(log, "bash", "-x", "downloadIstio.sh"); err != nil {
		log.Error("Script execution failed", zap.Error(err))
		return fmt.Errorf("failed to execute Istio installation: %w", err)
	}

	log.Info("Script execution completed successfully")

	// Check what's in current directory after script execution
	log.Info("Current directory contents after script execution:")
	command.Run(log, "pwd")
	command.Run(log, "ls", "-la")

	// Clean up the script file
	log.Info("Cleaning up script file...")
	command.Run(log, "rm", "downloadIstio.sh") // Ignore error as file might not exist

	// Find the downloaded Istio directory
	istioDir, err := findIstioDirectory(log)
	if err != nil {
		return fmt.Errorf("failed to find Istio directory: %w", err)
	}

	log.Info("Found Istio directory", zap.String("directory", istioDir))

	// Install istioctl to /usr/local/bin
	istioctlPath := filepath.Join(istioDir, "bin", "istioctl")
	if err := command.Run(log, "sudo", "install", "-o", "root", "-g", "root", "-m", "0755", istioctlPath, "/usr/local/bin/istioctl"); err != nil {
		return fmt.Errorf("failed to install istioctl: %w", err)
	}

	// Clean up the downloaded directory
	log.Debug("Cleaning up Istio installation directory...")
	command.Run(log, "rm", "-rf", istioDir) // Ignore error as directory might not exist

	return nil
}

// findIstioDirectory finds the downloaded Istio directory
func findIstioDirectory(log *logger.Logger) (string, error) {
	// Use find to locate istio directories
	output, err := command.RunWithOutput(log, "find", ".", "-maxdepth", "1", "-type", "d", "-name", "istio-*")
	if err != nil {
		return "", fmt.Errorf("failed to find Istio directory: %w", err)
	}

	lines := strings.Split(strings.TrimSpace(output), "\n")
	if len(lines) == 0 || lines[0] == "" {
		return "", fmt.Errorf("no Istio directory found")
	}

	// Get the first directory and remove "./" prefix if present
	dir := strings.TrimPrefix(lines[0], "./")

	return dir, nil
}

// setupIstioctlPath adds istioctl to PATH in bashrc
func setupIstioctlPath(log *logger.Logger) error {
	log.Info("Adding istioctl to PATH in bashrc...")

	// Add istioctl to PATH in bashrc
	if err := command.Run(log, "bash", "-c", "echo 'export PATH=/usr/local/bin:$PATH' >> ~/.bashrc"); err != nil {
		return fmt.Errorf("failed to add istioctl to PATH: %w", err)
	}

	// Source bashrc to make PATH immediately available
	if err := command.Run(log, "bash", "-c", "source ~/.bashrc"); err != nil {
		log.Warn("Failed to source bashrc, istioctl will be available after shell restart", zap.Error(err))
	}

	return nil
}

// InstallIstioOnCluster installs Istio on the specified cluster type
func InstallIstioOnCluster(log *logger.Logger, clusterType string) error {
	log.Info("Installing Istio on cluster", zap.String("cluster_type", clusterType))

	// Prepare cluster for Istio installation
	if err := prepareClusterForIstio(log, clusterType); err != nil {
		return fmt.Errorf("failed to prepare cluster for Istio: %w", err)
	}

	// Install Istio using istioctl
	if err := installIstioMesh(log, clusterType); err != nil {
		return fmt.Errorf("failed to install Istio mesh: %w", err)
	}

	// Verify installation
	if err := verifyIstioInstallation(log); err != nil {
		log.Warn("Istio installation verification failed", zap.Error(err))
	}

	log.Info("Istio successfully installed on cluster")
	return nil
}

// prepareClusterForIstio prepares the cluster for Istio installation
func prepareClusterForIstio(log *logger.Logger, clusterType string) error {
	log.Info("Preparing cluster for Istio installation", zap.String("cluster_type", clusterType))

	switch clusterType {
	case "kind":
		return prepareKindForIstio(log)
	case "minikube":
		return prepareMinikubeForIstio(log)
	default:
		return fmt.Errorf("unsupported cluster type: %s", clusterType)
	}
}

// prepareKindForIstio prepares a Kind cluster for Istio
func prepareKindForIstio(log *logger.Logger) error {
	log.Info("Preparing Kind cluster for Istio...")

	// Kind clusters are ready for Istio by default
	// No specific preparation needed
	log.Info("Kind cluster is ready for Istio installation")
	return nil
}

// prepareMinikubeForIstio prepares a Minikube cluster for Istio
func prepareMinikubeForIstio(log *logger.Logger) error {
	log.Info("Preparing Minikube cluster for Istio...")

	// Check if minikube is running
	if err := command.Run(log, "minikube", "status"); err != nil {
		return fmt.Errorf("minikube is not running, please start it first")
	}

	// Check if minikube has sufficient resources
	log.Info("Checking Minikube resources...")
	// Note: We assume the user has already started minikube with sufficient resources
	// as mentioned in the installation instructions

	log.Info("Minikube cluster is ready for Istio installation")
	return nil
}

// installIstioMesh installs the Istio mesh on the cluster
func installIstioMesh(log *logger.Logger, clusterType string) error {
	log.Info("Installing Istio mesh on cluster...")

	// Install Istio with default profile
	if err := command.Run(log, "istioctl", "install", "--set", "values.defaultRevision=default"); err != nil {
		return fmt.Errorf("failed to install Istio mesh: %w", err)
	}

	// Enable Istio sidecar injection for default namespace
	log.Info("Enabling Istio sidecar injection for default namespace...")
	if err := command.Run(log, "kubectl", "label", "namespace", "default", "istio-injection=enabled", "--overwrite"); err != nil {
		log.Warn("Failed to enable sidecar injection for default namespace", zap.Error(err))
	}

	return nil
}

// verifyIstioInstallation verifies that Istio is properly installed
func verifyIstioInstallation(log *logger.Logger) error {
	log.Info("Verifying Istio installation...")

	// Check if Istio pods are running
	if err := command.Run(log, "kubectl", "get", "pods", "-n", "istio-system"); err != nil {
		return fmt.Errorf("failed to check Istio pods: %w", err)
	}

	// Verify istioctl installation
	if err := command.Run(log, "istioctl", "version"); err != nil {
		return fmt.Errorf("failed to verify istioctl: %w", err)
	}

	log.Info("Istio installation verified successfully")
	return nil
}

// GetIstioInstructions returns instructions for using Istio based on cluster type
func GetIstioInstructions(clusterType string) []string {
	instructions := []string{
		"🔧 Istio Service Mesh:",
		"   • Check Istio status: istioctl version",
		"   • View Istio pods: k get pods -n istio-system",
		"   • Enable sidecar injection: k label namespace <namespace> istio-injection=enabled",
	}

	switch clusterType {
	case "kind":
		instructions = append(instructions,
			"   • Deploy sample app: kubectl apply -f https://raw.githubusercontent.com/istio/istio/release-1.27/samples/bookinfo/platform/kube/bookinfo.yaml",
			"   • Access via port-forward: k port-forward svc/productpage 9080:9080",
		)
	case "minikube":
		instructions = append(instructions,
			"   • Deploy sample app: kubectl apply -f https://raw.githubusercontent.com/istio/istio/release-1.27/samples/bookinfo/platform/kube/bookinfo.yaml",
			"   • Access via minikube tunnel: minikube tunnel (in separate terminal)",
			"   • Then: k port-forward svc/productpage 9080:9080",
		)
	}

	instructions = append(instructions,
		"   • Access Bookinfo: http://localhost:9080/productpage",
		"   • Uninstall Istio: istioctl uninstall --purge",
	)

	return instructions
}
