package hedera

import "fmt"

type NetworkName string

const (
	Mainnet    NetworkName = "mainnet"
	Testnet    NetworkName = "testnet"
	Previewnet NetworkName = "previewnet"
)

//func (networkName NetworkName) String() string {
//	switch networkName {
//	case Mainnet:
//		return "mainnet"
//	case Testnet:
//		return "testnet"
//	case Previewnet:
//		return "previewnet"
//	}
//
//	panic(fmt.Sprintf("unreacahble: NetworkName.String() switch statement is non-exhaustive. NetworkName: %s", networkName))
//}

func (networkName NetworkName) Network() string {
	switch networkName {
	case Mainnet:
		return "0"
	case Testnet:
		return "1"
	case Previewnet:
		return "2"
	}

	panic(fmt.Sprintf("unreacahble: NetworkName.Network() switch statement is non-exhaustive. NetworkName: %s", networkName))
}
