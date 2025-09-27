package types

// Runtime represents a container runtime option
type Runtime struct {
	Key  string
	Name string
}

// Distribution represents a Kubernetes distribution option
type Distribution struct {
	Key  string
	Name string
}

// InstallStep represents a single installation step
type InstallStep struct {
	Name string
}

// Config holds the installation configuration
type Config struct {
	Runtime      Runtime
	Distribution Distribution
	Steps        []InstallStep
}
