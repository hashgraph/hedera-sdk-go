package hedera

/*-
 *
 * Hedera Go SDK
 *
 * Copyright (C) 2020 - 2022 Hedera Hashgraph, LLC
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

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
