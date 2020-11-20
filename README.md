![](https://img.shields.io/github/v/tag/hashgraph/hedera-sdk-go)
[![Actions Status](https://github.com/hashgraph/hedera-sdk-go/v2/workflows/Go/badge.svg)](https://github.com/hashgraph/hedera-sdk-Go/actions?query=workflow%3AGo)
![](https://img.shields.io/github/go-mod/go-version/hashgraph/hedera-sdk-go)
[![](https://godoc.org/github.com/hashgraph/hedera-sdk-go/v2?status.svg)](http://godoc.org/github.com/hashgraph/hedera-sdk-go/v2)

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
$ go get github.com/hashgraph/hedera-sdk-go/v2
```

## Running Integration Tests
```bash
$ env CONFIG_FILE="<your_config_file>" go test -v _Execute
```

or

```bash
$ env CONFIG_FILE="<your_config_file>" OPERATOR_KEY="<key>" OPERATOR_ID="<id>" go test -v _Execute
```

or

```bash
$ env OPERATOR_KEY="<key>" OPERATOR_ID="<id>" go test -v _Execute
```

The config file _can_ contain both the network and the operator, but you can also
use environment variables `OPERATOR_KEY` and `OPERATOR_ID`. If both are provided
the network is used from the config file, but for the operator the environment variables
take precedence. If the config file is not provided then the network will default to testnet
and `OPERATOR_KEY` and `OPERATOR_ID` **must** be provided.

[Example Config File](./client-config-with-operator.json)

## Contributing to this Project

We welcome participation from all developers!
For instructions on how to contribute to this repo, please
review the [Contributing Guide](CONTRIBUTING.md).

## License Information

Licensed under Apache License,
Version 2.0 – see [LICENSE](LICENSE) in this repo
or [apache.org/licenses/LICENSE-2.0](http://www.apache.org/licenses/LICENSE-2.0)
