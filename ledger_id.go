package hedera

import (
	"bytes"
	"encoding/hex"

	"github.com/pkg/errors"
)

type LedgerID struct {
	LedgerID []byte
}

func LedgerIDFromString(id string) (*LedgerID, error) {
	switch id {
	case "mainnet": //nolint
		return &LedgerID{
			LedgerID: []byte{0},
		}, nil
	case "testnet": //nolint
		return &LedgerID{
			LedgerID: []byte{1},
		}, nil
	case "previewnet": //nolint
		return &LedgerID{
			LedgerID: []byte{2},
		}, nil
	default:
		temp, err := hex.DecodeString(id)
		if err != nil {
			return &LedgerID{}, err
		}
		return &LedgerID{
			LedgerID: temp,
		}, nil
	}
}

func LedgerIDFromBytes(byt []byte) *LedgerID {
	return &LedgerID{
		LedgerID: byt,
	}
}

func LedgerIDFromNetworkName(network NetworkName) (*LedgerID, error) {
	switch network.String() {
	case "mainnet": //nolint
		return &LedgerID{
			LedgerID: []byte{0},
		}, nil
	case "testnet": //nolint
		return &LedgerID{
			LedgerID: []byte{1},
		}, nil
	case "previewnet": //nolint
		return &LedgerID{
			LedgerID: []byte{2},
		}, nil
	default:
		return &LedgerID{}, errors.New("unknown network in network name")
	}
}

func NewLedgerIDMainnet() *LedgerID {
	return &LedgerID{
		LedgerID: []byte{0},
	}
}

func NewLedgerIDTestnet() *LedgerID {
	return &LedgerID{
		LedgerID: []byte{1},
	}
}

func NewLedgerIDPreviewnet() *LedgerID {
	return &LedgerID{
		LedgerID: []byte{2},
	}
}

func (id *LedgerID) IsMainnet() bool {
	return hex.EncodeToString(id.LedgerID) == "00"
}

func (id *LedgerID) IsTestnet() bool {
	return hex.EncodeToString(id.LedgerID) == "01"
}

func (id *LedgerID) IsPreviewnet() bool {
	return hex.EncodeToString(id.LedgerID) == "02"
}

func (id *LedgerID) String() string {
	h := hex.EncodeToString(id.LedgerID)
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

func (id *LedgerID) _ForChecksum() string {
	if bytes.Equal(id.LedgerID, []byte{0}) { //nolint
		return "0"
	} else if bytes.Equal(id.LedgerID, []byte{1}) {
		return "1"
	} else if bytes.Equal(id.LedgerID, []byte{2}) {
		return "2"
	} else {
		return hex.EncodeToString(id.LedgerID)
	}
}

func (id *LedgerID) ToBytes() []byte {
	return id.LedgerID
}

func (id *LedgerID) ToNetworkName() (NetworkName, error) {
	switch hex.EncodeToString(id.LedgerID) {
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
