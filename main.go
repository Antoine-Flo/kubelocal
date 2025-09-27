package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/Antoine-Flo/kubelocal/installer"
	"github.com/Antoine-Flo/kubelocal/ui"
)

func main() {
	setupSignalHandler()

	if err := runInstall(); err != nil {
		log.Printf("Installation failed: %v", err)
		os.Exit(1)
	}
}

// setupSignalHandler handles Ctrl+C gracefully
func setupSignalHandler() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-c
		fmt.Printf("\n%sInstallation cancelled%s\n", ui.ColorRed, ui.ColorReset)
		os.Exit(0)
	}()
}

// runInstall orchestrates the installation process
func runInstall() error {
	ui.PrintHeader("🚀 Kubernetes Local Environment Installer")

	runtime, err := chooseContainerRuntime()
	if err != nil {
		return fmt.Errorf("failed to choose runtime: %w", err)
	}

	distribution, err := chooseKubernetesDistribution()
	if err != nil {
		return fmt.Errorf("failed to choose distribution: %w", err)
	}

	return executeInstallation(runtime, distribution)
}

// chooseContainerRuntime prompts user to select container runtime
func chooseContainerRuntime() (string, error) {
	options := installer.GetRuntimeOptions()
	choice, err := ui.AskChoiceWithDefault("Install Docker (d) or Podman (p) : ", options, "d")
	if err != nil {
		return "", err
	}
	return choice, nil
}

// chooseKubernetesDistribution prompts user to select Kubernetes distribution
func chooseKubernetesDistribution() (string, error) {
	options := installer.GetDistributionOptions()
	choice, err := ui.AskChoiceWithDefault("Install Kind (k) or Minikube (m) : ", options, "k")
	if err != nil {
		return "", err
	}
	return choice, nil
}

// executeInstallation creates config and runs installation
func executeInstallation(runtimeKey, distributionKey string) error {
	config := installer.CreateConfig(runtimeKey, distributionKey)
	return installer.Run(config)
}
