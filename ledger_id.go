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

import (
	"encoding/hex"

	"github.com/pkg/errors"
)

type LedgerID struct {
	_LedgerIDBytes []byte
}

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

func LedgerIDFromBytes(byt []byte) *LedgerID {
	return &LedgerID{
		_LedgerIDBytes: byt,
	}
}

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

func NewLedgerIDMainnet() *LedgerID {
	return &LedgerID{
		_LedgerIDBytes: []byte{0},
	}
}

func NewLedgerIDTestnet() *LedgerID {
	return &LedgerID{
		_LedgerIDBytes: []byte{1},
	}
}

func NewLedgerIDPreviewnet() *LedgerID {
	return &LedgerID{
		_LedgerIDBytes: []byte{2},
	}
}

func (id *LedgerID) IsMainnet() bool {
	return hex.EncodeToString(id._LedgerIDBytes) == "00"
}

func (id *LedgerID) IsTestnet() bool {
	return hex.EncodeToString(id._LedgerIDBytes) == "01"
}

func (id *LedgerID) IsPreviewnet() bool {
	return hex.EncodeToString(id._LedgerIDBytes) == "02"
}

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

func (id *LedgerID) ToBytes() []byte {
	return id._LedgerIDBytes
}

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
