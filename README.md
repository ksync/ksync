
# Getting started

TODO

# Config

TODO

# Architecture

TODO

## Commands

- `ksync init`

  Sets the cluster up by starting the radar daemonset.

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

- `ksync background`

  Install the watcher into the local process manager and run it in the background.

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

# Outstanding

- Tests
- Add an option to resolve symlinks?

  This might work for a one time sync for remote -> local, but it would break things during any kind of general work.

- TLS for mirror
- Put some details into the task bar (see https://github.com/cratonica/trayhost)
  - Active syncs
  - Files being updated

[protoc]: https://github.com/golang/protobuf/
[protoc-gen-go]: https://github.com/golang/protobuf/
[dep]: https://github.com/golang/dep/
