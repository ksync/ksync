package doctor

var (
	KubernetesRange = ">=1.7.0"
	DockerAPIRange  = ">=1.25.0"

	// This is really for the error message
	DockerRange = ">=1.13.0"

	DockerDriver = map[string]bool{
		"overlay": true,
	}
)
