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

#### Note

google.golang.org/protobuf v1.27.1 Breaks the SDK as it contains multiple protobuf files
with the same name. Make sure to use v1.26.1 instead. The follow snippet can be used in 
`go.mod` to force the project to use v1.26.1

```
replace (
	google.golang.org/protobuf v1.27.1 => google.golang.org/protobuf v1.26.1-0.20210525005349-febffdd88e85
)
```

## Running Integration Tests
```bash
$ env CONFIG_FILE="<your_config_file>" go test -v Integration -timeout 9999s ```

or

```bash
$ env CONFIG_FILE="<your_config_file>" OPERATOR_KEY="<key>" OPERATOR_ID="<id>" go test -v Integration -timeout 9999s
```

or

```bash
$ env OPERATOR_KEY="<key>" OPERATOR_ID="<id>" go test -v Integration -timeout 9999s
```

The config file _can_ contain both the network and the operator, but you can also
use environment variables `OPERATOR_KEY` and `OPERATOR_ID`. If both are provided
the network is used from the config file, but for the operator the environment variables
take precedence. If the config file is not provided then the network will default to testnet
and `OPERATOR_KEY` and `OPERATOR_ID` **must** be provided.

[Example Config File](./client-config-with-operator.json)

## Support

If you have a question on how to use the product, please see our
[support guide](https://github.com/hashgraph/.github/blob/main/SUPPORT.md).

## Contributing

Contributions are welcome. Please see the
[contributing guide](https://github.com/hashgraph/.github/blob/main/CONTRIBUTING.md)
to see how you can get involved.

## Code of Conduct

This project is governed by the
[Contributor Covenant Code of Conduct](https://github.com/hashgraph/.github/blob/main/CODE_OF_CONDUCT.md). By
participating, you are expected to uphold this code of conduct. Please report unacceptable behavior
to [oss@hedera.com](mailto:oss@hedera.com).

## License

[Apache License 2.0](LICENSE)
