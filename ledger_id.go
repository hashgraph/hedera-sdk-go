package hedera

import (
	"encoding/hex"
	"fmt"

	"github.com/pkg/errors"
)

type LedgerID struct {
	LedgerID []byte
}

func LedgerIDFromString(id string) (*LedgerID, error) {
	switch id {
	case "mainnet": //nolint
		temp, err := hex.DecodeString("00")
		if err != nil {
			return &LedgerID{}, err
		}
		return &LedgerID{
			LedgerID: temp,
		}, nil
	case "testnet": //nolint
		temp, err := hex.DecodeString("01")
		if err != nil {
			return &LedgerID{}, err
		}
		return &LedgerID{
			LedgerID: temp,
		}, nil
	case "previewnet": //nolint
		temp, err := hex.DecodeString("02")
		if err != nil {
			return &LedgerID{}, err
		}
		return &LedgerID{
			LedgerID: temp,
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
		temp, err := hex.DecodeString("00")
		if err != nil {
			return &LedgerID{}, err
		}
		return &LedgerID{
			LedgerID: temp,
		}, nil
	case "testnet": //nolint
		temp, err := hex.DecodeString("01")
		if err != nil {
			return &LedgerID{}, err
		}
		return &LedgerID{
			LedgerID: temp,
		}, nil
	case "previewnet": //nolint
		temp, err := hex.DecodeString("02")
		if err != nil {
			return &LedgerID{}, err
		}
		return &LedgerID{
			LedgerID: temp,
		}, nil
	default:
		return &LedgerID{}, errors.New("unknown network in network name")
	}
}

func NewLedgerIDMainnet() *LedgerID {
	temp, _ := hex.DecodeString("00")
	return &LedgerID{
		LedgerID: temp,
	}
}

func NewLedgerIDTestnet() *LedgerID {
	temp, _ := hex.DecodeString("01")
	return &LedgerID{
		LedgerID: temp,
	}
}

func NewLedgerIDPreviewnet() *LedgerID {
	temp, _ := hex.DecodeString("02")
	return &LedgerID{
		LedgerID: temp,
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
	switch hex.EncodeToString(id.LedgerID) {
	case "00":
		return "mainnet"
	case "01":
		return "testnet"
	case "02":
		return "previewnet"
	default:
		return hex.EncodeToString(id.LedgerID)
	}
}

func (id *LedgerID) _ForChecksum() string {
	h := hex.EncodeToString(id.LedgerID)
	switch h {
	case "00":
		return "0"
	case "01":
		return "1"
	case "02":
		return "2"
	default:
		return h
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
	}

	panic(fmt.Sprintf("unreachable: LedgerID.ToNetworkName() switch statement is non-exhaustive. LederID: %s", hex.EncodeToString(id.LedgerID)))
}
