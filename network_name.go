package hiero

// SPDX-License-Identifier: Apache-2.0

type NetworkName string

const (
	NetworkNameMainnet    NetworkName = "mainnet"
	NetworkNameTestnet    NetworkName = "testnet"
	NetworkNamePreviewnet NetworkName = "previewnet"
	NetworkNameOther      NetworkName = "other"
)

// Deprecated
func (networkName NetworkName) String() string { //nolint
	switch networkName {
	case NetworkNameMainnet:
		return "mainnet"
	case NetworkNameTestnet:
		return "testnet"
	case NetworkNamePreviewnet:
		return "previewnet"
	case NetworkNameOther:
		return "other"
	}

	panic("unreachable: NetworkName.String() switch statement is non-exhaustive.")
}

// Deprecated
func NetworkNameFromString(s string) NetworkName { //nolint
	switch s {
	case "mainnet":
		return NetworkNameMainnet
	case "testnet":
		return NetworkNameTestnet
	case "previewnet":
		return NetworkNamePreviewnet
	case "other":
		return NetworkNameOther
	}

	panic("unreachable: NetworkName.String() switch statement is non-exhaustive.")
}
