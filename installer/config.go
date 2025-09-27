package installer

import "github.com/Antoine-Flo/kubelocal/types"

// GetRuntimeOptions returns available container runtime options
func GetRuntimeOptions() map[string]string {
	return map[string]string{
		"d": "Docker",
		"p": "Podman",
	}
}

// GetDistributionOptions returns available Kubernetes distribution options
func GetDistributionOptions() map[string]string {
	return map[string]string{
		"k": "Kind",
		"m": "Minikube",
	}
}

// CreateConfig creates installation configuration
func CreateConfig(runtimeKey, distributionKey string) types.Config {
	runtimeOpts := GetRuntimeOptions()
	distOpts := GetDistributionOptions()

	return types.Config{
		Runtime:      types.Runtime{Key: runtimeKey, Name: runtimeOpts[runtimeKey]},
		Distribution: types.Distribution{Key: distributionKey, Name: distOpts[distributionKey]},
		Steps:        getInstallSteps(runtimeOpts[runtimeKey], distOpts[distributionKey]),
	}
}
