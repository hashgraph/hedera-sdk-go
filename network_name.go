package hedera

import "fmt"

type NetworkName string

const (
	NetworkNameMainnet    NetworkName = "mainnet"
	NetworkNameTestnet    NetworkName = "testnet"
	NetworkNamePreviewnet NetworkName = "previewnet"
)

func (networkName NetworkName) _LedgerID() string {
	switch networkName {
	case NetworkNameMainnet:
		return "0"
	case NetworkNameTestnet:
		return "1"
	case NetworkNamePreviewnet:
		return "2"
	}

	panic(fmt.Sprintf("unreachable: NetworkName._LedgerID() switch statement is non-exhaustive. NetworkName: %s", networkName))
}

func (networkName NetworkName) String() string {
	switch networkName {
	case NetworkNameMainnet:
		return "mainnet"
	case NetworkNameTestnet:
		return "testnet"
	case NetworkNamePreviewnet:
		return "previewnet"
	}

	panic("unreachable: NetworkName.String() switch statement is non-exhaustive.")
}
