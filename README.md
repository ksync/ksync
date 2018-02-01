<img align="left" src="logos/ksync_logo_color.png">

<img align="right" src="https://goreportcard.com/badge/github.com/vapor-ware/ksync">
<img align="right" src="https://circleci.com/gh/vapor-ware/ksync.svg?style=svg&circle-token=429269824f09028301b6e65310bd0cea8031d292">

<br clear="all" />

------------

ksync speeds up developers who build applications for Kubernetes. It syncs files between a local directory and arbitrary containers running remotely. You do not need to change your existing workflow to develop directly on a Kubernetes cluster.

If you've been wanting to do something like `docker run -v /foo:/bar` with Kubernetes, ksync is for you!

Using ksync is as simple as:

1. `ksync init` to run the server component.
1. `ksync create --pod=my-pod local_directory remote_directory` to configure a folder you'd like to sync between the cluster and your local system.
1. `ksync watch` to monitor the kubernetes API and sync.
1. Use your favorite editor, like [Atom][atom] or [Sublime Text][st3] to modify the application. It will auto-reload for you remotely, in seconds.

# Installation

```bash
curl https://vapor-ware.github.io/gimme-that/gimme.sh | bash
```

You can also download the [latest release][latest-release] and install it yourself.

## Updating

To update to (or check for) a newer version of `ksync`, you can simply call the built in updater.

```shell
ksync update
```

This will check GitHub for the [latest official release][latest-release] and download it if newer. You can also follow the [installation](#installation) instructions or compile the binary yourself.

Once a newer `ksync` binary has been downloaded, the cluster portion can be updated with `ksync init`.

```shell
ksync init --upgrade
```

This will deploy the cluster component matching your `ksync` version to the target cluster. You can check the versions of both `ksync` and `radar` by running `ksync version`.

```shell
ksync version
```

# Prerequisites

- Kubernetes cluster. Take a look at the [docs][k8s-setup] for instructions on how to do it.

    A couple fast and easy solutions:

    - To keep it all local, check out [minikube][minikube].
    - If you'd like something remote, [GKE][GKE] can create a cluster fast.

- `kubectl` configured to talk to your cluster.

# Getting Started

1. Install ksync. This will fetch the binary and put it at `/usr/local/bin`. Feel free to just download the release binary for your platform and install it yourself.

    ```bash
    curl https://vapor-ware.github.io/gimme-that/gimme.sh | bash
    ```

1. Initialize ksync and install the server component on your cluster. The server component is a DaemonSet that provides access to each node's filesystem.

    ```bash
    ksync init
    ```

1. Startup the local client. It watches your local config to start new jobs and the kubernetes API to react when things change there. This will just put it into the background. Feel free to run in a separate terminal or add as a service to your host.

    ```bash
    ksync watch &
    ```

1. Add the [demo app][demo-app] to your cluster. This is a simple python app made with flask. Because ksync moves files around, it would work for any kind of data you'd like to move between your local system and the cluster.

    ```bash
    kubectl apply -f https://vapor-ware.github.io/ksync/example/app/app.yaml
    ```

1. Make sure that the app is ready and running.

    ```bash
    kubectl get po --selector=app=app
    ```

1. Create a new spec that describes a folder to sync between a local directory and a directory inside a running container on the remote cluster. The local directory is empty and that is okay. Because ksync is bi-directional, it will move all the files from the running container locally. This is just a convenient way to get the code from the container and skip a couple steps. If you're working with a local copy already, only the most recently updated files will be transfered between the container and your local machine.

    ```bash
    mkdir -p $(pwd)/ksync
    ksync create --selector=app=app $(pwd)/ksync /code
    ```

1. Check on the status.

    ```bash
    ksync get
    ```

1. Forward the remote port to your local system.

    ```bash
    kubectl get po --selector=app=app -o=custom-columns=:metadata.name --no-headers | \
        xargs -IPOD kubectl port-forward POD 8080:80 &
    ```

1. Take a look at what the app's response is now. You'll see all the files in the remote container, their modification times and when the container was last restarted.

    ```bash
    curl localhost:8080
    ```

1. Open up the code in your favorite editor. For demo purposes, this assumes you've configured `EDITOR`. You really can open it however you'd like though.

    ```bash
    open ksync/server.py
    ```

1. Add a new key to the JSON response by editing the return value.

    ```python
    return jsonify({
        "ksync": True,
        "restart": LAST_RESTART,
        "pod": os.environ.get('POD_NAME'),
        "files": file_list,
    })
    ```

1. Take a look at the status now. It should be reloading the remote container.

    ```bash
    ksync get
    ```

1. After about 10 seconds, hit the container again and you should see your new response.

    ```bash
    curl localhost:8080
    ```

## Further Exploration

- Modify the number of replicas and see what happens.

    ```bash
    kubectl scale deployment/app --replicas=2
    ```

- Startup the [visualization][frontend] so you can see updates in real time. Save some files and change the replica count of app to see the updates.

    ```bash
    kubectl apply -f https://vapor-ware.github.io/ksync/example/frontend/frontend.yaml
    kubectl get po \
        --selector=app=frontend \
        -o=custom-columns=:metadata.name \
        --no-headers \
        | xargs -IPOD kubectl port-forward POD 8081:80 &
    python -mwebbrowser http://localhost:8081
    ```

![visualizer](docs/visualizer.png)

# Tested Configurations

## Cluster

- Minikube
    - v0.23.*
    - v0.24.*

- GKE
    - v1.7.*
    - v1.8.*

- Docker for Mac (Kubernetes)
    - 17.12-ce

## Docker

- Docker
    - 1.13.*
    - 17.*-ce

## Filesystem

- OverlayFS (overlay2)

# Troubleshooting

- `ERROR Path ... does not exist on the server`

    There's likely something in your configuration that we're not able to handle yet.

- `client is newer than server (client API version: ..., server API version: ...)`

    You're using an older version of docker than we support.

# Documentation

More detailed documentation can be found in the [docs](docs) directory.

- [Architecture](docs/architecture.md)
- [Development](docs/development.md)
- [Releasing](docs/releasing.md)

[atom]: https://atom.io/
[st3]: https://www.sublimetext.com/
[latest-release]: https://github.com/vapor-ware/ksync/releases
[k8s-setup]: https://kubernetes.io/docs/setup/pick-right-solution/
[GKE]: https://cloud.google.com/kubernetes-engine/docs/quickstart
[minikube]: https://github.com/kubernetes/minikube
[demo-app]: https://vapor-ware.github.io/ksync/example/app/app.yaml
[frontend]: https://vapor-ware.github.io/ksync/example/frontend/frontend.yaml
