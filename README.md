
# Getting started

TODO

# Config

TODO

# Architecture

TODO

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
go get -u github.com/golang/dep/cmd/dep
```

# Outstanding

- Tests
- Add an option to resolve symlinks?

  This might work for a one time sync for remote -> local, but it would break things during any kind of general work.

- TLS for mirror

[protoc]: https://github.com/golang/protobuf/
[protoc-gen-go]: https://github.com/golang/protobuf/
[dep]: https://github.com/golang/dep/
