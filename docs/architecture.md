TODO

- ksync has three parts: a client (`ksync`), a server (`radar`) and a server to handle file syncing (`mirror`).
- Radar and mirror run inside of your Kubernetes cluster as a DaemonSet on every node.
- Radar provides an API to discover what the remote container's filepath is and manage the container lifecycle.
- [Mirror][mirror] is a real-time, two-way sync. The server operates on your Kubernetes cluster and the client is managed by `ksync` locally.

## Workflow

1. `ksync init` starts the DaemonSet on the remote cluster.

1. `ksync watch` starts up to manage the lifecycle of syncs.

1. `ksync create` adds a spec to the config. This contains everything required to locate a remote container.

1. `ksync watch` sees a change to the config and looks up the remote container.

    - If it does not exist, `watch` will continue to monitor the Kubernetes cluster for something that matches. When that happens, it will move to the next step.

1. `ksync watch` finds a remote container and creates a tunnel to the `radar` server running on the node that the remote container is running on.

1. The remote `radar` server inspects the remote container and returns the file path that contains the container's mounted filesystem.

1. `watch` starts a docker container in the background. This has the correct host path mounted into it.

1. The docker container runs `ksync run`. This creates a tunnel to the `mirror` server running on the node that the remote container is running on. It then executes the `mirror` client with the local host path and the remote container path.
