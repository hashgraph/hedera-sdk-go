package hedera

type NetworkName string

const (
	NetworkNameMainnet    NetworkName = "mainnet"
	NetworkNameTestnet    NetworkName = "testnet"
	NetworkNamePreviewnet NetworkName = "previewnet"
)

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

func NetworkNameFromString(s string) NetworkName { //nolint
	switch s {
	case "mainnet":
		return NetworkNameMainnet
	case "testnet":
		return NetworkNameTestnet
	case "previewnet":
		return NetworkNamePreviewnet
	}

	panic("unreachable: NetworkName.String() switch statement is non-exhaustive.")
}
