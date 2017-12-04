[![CircleCI](https://circleci.com/gh/vapor-ware/ksync.svg?style=svg&circle-token=429269824f09028301b6e65310bd0cea8031d292)](https://circleci.com/gh/vapor-ware/ksync)

ksync is a tool for syncing files between a local directory and arbitrary containers running remotely on a Kubernetes cluster. It does not require any changes to the remote containers and works transparently.

TODO - something about the watch flow instead of run.

Use ksync to:

- Develop applications remotely, while still using your favorite editor and local environment.
- TODO

# Demo

TODO

# Install

## Quick
Grab the [latest release](https://github.com/vapor-ware/ksync/releases/latest) from the [releases](https://github.com/vapor-ware/ksync/releases) page. Place the binary in your `PATH` and make executable. Alternatively you can run the following command to do this for you (places binary in `/usr/local/bin`).

```shell
curl https://github.com/vapor-ware/ksync/releases/download/latest/ksync_$(go env GOHOSTOS)_$(go env GOHOSTARCH) -o /usr/local/bin/ksync && chmod +x /usr/local/bin/ksync
```

## Development
You can also get the code and compile it yourself. If you have `go` installed you can run the following.

```shell
go get github.com/vapor-ware/ksync
cd ${GOPATH}/src/github.com/vapor-ware/ksync
go install cmd/*
```

**Note**: If you compile the binaries yourself the output of `ksync version` may not be correct. Only the binaries on the [releases](https://github.com/vapor-ware/ksync/releases) page are stamped with this information.

# Getting started

1. Initialize ksync and install radar.

    ```bash
    ksync init
    ```

1. Startup watch in the background.

    ```bash
    ksync watch
    ```

1. Create a new spec.

    ```bash
    ksync create --selector=app=demo /tmp/demo /demo-files
    ```

1. Start the demo container on your cluster.

    ```bash
    kubectl apply -f TODO
    ```

1. Look at the status of your specs.

    ```bash
    ksync get
    ```

1. See the local files that were written to your local environment.

    ```bash
    ls /tmp/demo
    ```

1. Edit something locally (to see it sync to the remote container).

    ```bash
    touch /tmp/demo/foobar
    ```

1. Verify that it updated.

    ```bash
    kubectl exec -it \
        $(kubectl get po --selector=app=volume \
            | tail -n1 \
            | awk '{ print $1 }') \
        -- ls -la /tmp/demo
    ```

# Config

TODO

# Architecture

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

# Commands

- `ksync init`

    - Sets the cluster up by starting the radar daemonset.
    - Starts `ksync watch` in the background. TODO

- `ksync create`

    Add a spec to sync. This gets watched and started/stopped automatically.

- `ksync delete`

    Remove a spec to sync.

- `ksync get` TODO: is this maybe a better sync list? can show running and waiting ones.

    Fetch the status of all current specs

- `ksync watch`

    Watch for matching pods in the background (based off pod name and selector). Start syncs for any that come online.

- `ksync doctor`

    Debug what's happening under the covers and look for any possible issues with the system.

- `ksync version`

    Print out version information for the local binary. If the server binary is reachable and healthy, print information for that as well.

# Development

## Dependencies

- [protoc][protoc]

```bash
brew install protobuf
```

- [protoc-gen-go][protoc-gen-go]

```bash
go get -u github.com/golang/protobuf/protoc-gen-go
```

- [dep][dep]

```bash
go install -u github.com/golang/dep/cmd/dep
```

# Troubleshooting

- ntp issues

[protoc]: https://github.com/golang/protobuf/
[protoc-gen-go]: https://github.com/golang/protobuf/
[dep]: https://github.com/golang/dep/
[mirror]: https://github.com/stephenh/mirror
