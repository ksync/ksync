You can also get the code and compile it yourself. If you have `go` installed you can run the following.

```shell
go get github.com/vapor-ware/ksync
cd ${GOPATH}/src/github.com/vapor-ware/ksync
go install cmd/*
```

**Note**: If you compile the binaries yourself the output of `ksync version` may not be correct. Only the binaries on the [releases](https://github.com/vapor-ware/ksync/releases) page are stamped with this information.

The makefile has some convenient helpers:

- `watch` - Watches the directory for changes, rebuilds the binary and starts `ksync watch`. This requires `ag` and `entr`.
- `docker-binary docker-build docker-push` - Build and refresh the cluster docker pieces.
- `lint` - Run the linter.
- `radar-logs` - Streams the logs from the cluster locally to help with debugging. Requires `stern`.

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

[protoc]: https://github.com/golang/protobuf/
[protoc-gen-go]: https://github.com/golang/protobuf/
[dep]: https://github.com/golang/dep/
[mirror]: https://github.com/stephenh/mirror
