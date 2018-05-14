While `ksync` is primarily intended for scenarios with access to the internet, it can be used in a disconnected ("offline") state, where only access to the cluster is available.

# Offline

## Setup

In order to download required components, `ksync` does require internet access on initial setup. There are two ways to get around this, either you can download the components individually and place them in the correct locations, or you can run `init` when connected to the internet, then move offline (with one [exception](#remote)).

### Components

Here are the components retrieved during `init`:

#### Local

1. The local binary downloads the `syncthing` binary to it's config directory, `~/.ksync/bin/syncthing` (you should be able to use the latest stable from [the project repo](https://github.com/syncthing/syncthing) )
2. `ksync update` checks this repo for the latest release and updates it if found.

#### Remote

1. The remote server portion of `ksync` (called `radar`) uses a docker image. The nodes in the remote cluster attempt to pull the docker image `vaporio/ksync` from [Docker Hub](https://hub.docker.com/r/vaporio/ksync/). The tag matches the release version. You can pull this image into your cluster's image registry manually.

2. You can also specify a different image to use when running `init` by supplying the flag `--image=my-repo/ksync:tag`.

## Usage

Once all components are in place `ksync` should function normally. Some commands (such as `update`) will not function as they require internet access.
