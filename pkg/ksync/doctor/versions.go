package doctor

var (
	// KubernetesRange is the acceptible cluster versions
	KubernetesRange = ">=1.7.0"
	// DockerAPIRange is the acceptible range of API versions.
	DockerAPIRange = ">=1.25.0"

	// DockerRange is the user friendly version of DockerAPIRange. It is really
	// just for the error.
	DockerRange = ">=1.13.0"

	// DockerDriver is all the compatible storage drivers.
	DockerDriver = map[string]bool{
		"overlay": true,
	}
)
