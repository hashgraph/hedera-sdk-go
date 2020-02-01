![](https://img.shields.io/github/v/tag/hashgraph/hedera-sdk-go)

[![Actions Status](https://github.com/hashgraph/hedera-sdk-go/workflows/Go/badge.svg)](https://github.com/hashgraph/hedera-sdk-Go/actions?query=workflow%3AGo)

![](https://img.shields.io/github/go-mod/go-version/hashgraph/hedera-sdk-go)

[![](https://godoc.org/github.com/hashgraph/hedera-sdk-go?status.svg)](http://godoc.org/github.com/hashgraph/hedera-sdk-go)

# Hedera™ Hashgraph Go SDK

> The Go SDK for interacting with [Hedera Hashgraph]: the official distributed consensus
> platform built using the hashgraph consensus algorithm for fast, fair and secure
> transactions. Hedera enables and empowers developers to build an entirely new
> class of decentralized applications.

[Hedera Hashgraph]: https://hedera.com/

Hedera Hashgraph communicates using [gRPC]; the Protobufs definitions for the protocol are
available in the [hashgraph/hedera-protobuf] repository.

[gRPC]: https://grpc.io
[hashgraph/hedera-protobuf]: https://github.com/hashgraph/hedera-protobuf

## Install

```sh
$ go get github.com/hashgraph/hedera-sdk-go
```

## Contributing to this Project

We welcome participation from all developers!
For instructions on how to contribute to this repo, please
review the [Contributing Guide](CONTRIBUTING.md).

## Running Tests
If a change is made that requires a snapshot test to be updated you can update the snapshots using this command:

`env UPDATE_SNAPSHOTS=true go test`

## License Information

Licensed under Apache License,
Version 2.0 – see [LICENSE](LICENSE) in this repo
or [apache.org/licenses/LICENSE-2.0](http://www.apache.org/licenses/LICENSE-2.0)
