package hedera

import "fmt"

type NetworkName string

const (
	NetworkNameMainnet    NetworkName = "mainnet"
	NetworkNameTestnet    NetworkName = "testnet"
	NetworkNamePreviewnet NetworkName = "previewnet"
)

func (networkName NetworkName) ledgerID() string {
	switch networkName {
	case NetworkNameMainnet:
		return "0"
	case NetworkNameTestnet:
		return "1"
	case NetworkNamePreviewnet:
		return "2"
	}

	panic(fmt.Sprintf("unreacahble: NetworkName.ledgerID() switch statement is non-exhaustive. NetworkName: %s", networkName))
}
