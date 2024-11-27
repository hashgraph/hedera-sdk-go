package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"encoding/hex"

	"github.com/pkg/errors"
)

type LedgerID struct {
	_LedgerIDBytes []byte
}

// LedgerIDFromString returns a LedgerID from a string representation of a ledger ID.
func LedgerIDFromString(id string) (*LedgerID, error) {
	switch id {
	case "mainnet": //nolint
		return &LedgerID{
			_LedgerIDBytes: []byte{0},
		}, nil
	case "testnet": //nolint
		return &LedgerID{
			_LedgerIDBytes: []byte{1},
		}, nil
	case "previewnet": //nolint
		return &LedgerID{
			_LedgerIDBytes: []byte{2},
		}, nil
	default:
		temp, err := hex.DecodeString(id)
		if err != nil {
			return &LedgerID{}, err
		}
		return &LedgerID{
			_LedgerIDBytes: temp,
		}, nil
	}
}

// LedgerIDFromBytes returns a LedgerID from a byte representation of a ledger ID.
func LedgerIDFromBytes(byt []byte) *LedgerID {
	return &LedgerID{
		_LedgerIDBytes: byt,
	}
}

// LedgerIDFromNetworkName returns a LedgerID from a NetworkName.
func LedgerIDFromNetworkName(network NetworkName) (*LedgerID, error) {
	switch network.String() {
	case "mainnet": //nolint
		return &LedgerID{
			_LedgerIDBytes: []byte{0},
		}, nil
	case "testnet": //nolint
		return &LedgerID{
			_LedgerIDBytes: []byte{1},
		}, nil
	case "previewnet": //nolint
		return &LedgerID{
			_LedgerIDBytes: []byte{2},
		}, nil
	default:
		return &LedgerID{}, errors.New("unknown network in network name")
	}
}

// LedgerIDMainnet returns a LedgerID for mainnet.
func NewLedgerIDMainnet() *LedgerID {
	return &LedgerID{
		_LedgerIDBytes: []byte{0},
	}
}

// LedgerIDTestnet returns a LedgerID for testnet.
func NewLedgerIDTestnet() *LedgerID {
	return &LedgerID{
		_LedgerIDBytes: []byte{1},
	}
}

// LedgerIDPreviewnet returns a LedgerID for previewnet.
func NewLedgerIDPreviewnet() *LedgerID {
	return &LedgerID{
		_LedgerIDBytes: []byte{2},
	}
}

// IsMainnet returns true if the LedgerID is for mainnet.
func (id *LedgerID) IsMainnet() bool {
	return hex.EncodeToString(id._LedgerIDBytes) == "00"
}

// IsTestnet returns true if the LedgerID is for testnet.
func (id *LedgerID) IsTestnet() bool {
	return hex.EncodeToString(id._LedgerIDBytes) == "01"
}

// IsPreviewnet returns true if the LedgerID is for previewnet.
func (id *LedgerID) IsPreviewnet() bool {
	return hex.EncodeToString(id._LedgerIDBytes) == "02"
}

// String returns a string representation of the LedgerID.
func (id *LedgerID) String() string {
	h := hex.EncodeToString(id._LedgerIDBytes)
	switch h {
	case "00":
		return "mainnet"
	case "01":
		return "testnet"
	case "02":
		return "previewnet"
	default:
		return h
	}
}

// ToBytes returns a byte representation of the LedgerID.
func (id *LedgerID) ToBytes() []byte {
	return id._LedgerIDBytes
}

// ToNetworkName returns a NetworkName from the LedgerID.
func (id *LedgerID) ToNetworkName() (NetworkName, error) {
	switch hex.EncodeToString(id._LedgerIDBytes) {
	case "00":
		return NetworkNameMainnet, nil
	case "01":
		return NetworkNameTestnet, nil
	case "02":
		return NetworkNamePreviewnet, nil
	default:
		return NetworkNameOther, nil
	}
}
