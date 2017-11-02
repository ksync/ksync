
ksync is a tool for syncing files between a local directory and arbitrary containers running remotely on a Kubernetes cluster. It does not require any changes to the remote containers and works transparently.

TODO - something about the watch flow instead of run.

Use ksync to:

- Develop applications remotely, while still using your favorite editor and local environment.
- TODO

# Demo

TODO

# Install

TODO

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

- `ksync list` TODO: current functionality should be renamed, as it should list the syncs not files.

    Lists the files for a specific set of containers (selector, pod name, container name)

- `ksync create`

    Add a pattern to sync. This gets watched and started/stopped automatically.

- `ksync delete`

    Remove a pattern to sync.

- `ksync run`

    Runs a specific sync for the lifetime of a pod.

- `ksync get` TODO: is this maybe a better sync list? can show running and waiting ones.

    Fetch the status of all current syncs

- `ksync watch`

    Watch for matching pods in the background (based off pod name and selector). Start syncs for any that come online.

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

# Outstanding

## Major

- Put some details into the task bar (see https://github.com/cratonica/trayhost)
  - Active syncs
  - Files being updated
- Hot reload (docker restart on file change)

## Minor

- Only allow syncs on directories (not single files), check in add
- Move watch to init
    - Start with docker
    - Refactor `Service` to work with multiple container types? This might be more work than required.
- There is a timing error between the mirror restart and the run container restart.
    - run restarts mirror
    - run tries to create tunnel to mirror container
    - mirror is still starting (not listening), tunnel doesn't connect (crashing run)
    - loop endlessly
- Add health and readiness checks to `ksync init` for both radar and watch. There should be a flag that disables the wait (but it should wait by default and output status).
- Reduce the number of restarts that `ksync run` goes through while trying to setup a sync for the first time. It is tough to understand whether it is working or not (as a lot of the restarts are expected) and even harder to debug when it isn't working.
- Verify that the configured container user can actually write to localPath
- IO timeout errors (when the remote cluster cannot be reached) take a long time. There should be a better experience here.
- Test coverage
- TLS for mirror (is it required?)
- TLS for radar (is it required?)
- Allow configuration of who the container runs as (default to current user/group)
    - Maybe in the spec itself?

[protoc]: https://github.com/golang/protobuf/
[protoc-gen-go]: https://github.com/golang/protobuf/
[dep]: https://github.com/golang/dep/
[mirror]: https://github.com/stephenh/mirror
