package hedera

/*-
 *
 * Hedera Go SDK
 *
 * Copyright (C) 2020 - 2023 Hedera Hashgraph, LLC
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
	"encoding/json"

	"github.com/pkg/errors"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
)

// Represents the data of an Ethereum transaction.
type EthereumTransactionData struct {
	eip1559 *types.DynamicFeeTx
	legacy  *types.LegacyTx
}

// EthereumTransactionDataFromBytes constructs an EthereumTransactionData from a raw byte array.
func EthereumTransactionDataFromBytes(b []byte) (*EthereumTransactionData, error) {
	var transactionData EthereumTransactionData
	if b[0] == 2 {
		byt := b
		byt = append(byt[:0], byt[0+1:]...)
		err := rlp.DecodeBytes(byt, &transactionData.eip1559)
		if err != nil {
			return nil, err
		}

		return &transactionData, nil
	}

	err := rlp.DecodeBytes(b, &transactionData.legacy)
	if err != nil {
		return nil, err
	}

	return &transactionData, nil
}

// ToBytes returns the raw bytes of the Ethereum transaction.
func (ethereumTxData *EthereumTransactionData) ToBytes() ([]byte, error) {
	var byt []byte
	var err error
	if ethereumTxData.eip1559 != nil {
		byt, err = rlp.EncodeToBytes(ethereumTxData.eip1559)
		if err != nil {
			return []byte{}, err
		}
		byt = append([]byte{2}, byt...)

		return byt, nil
	}

	byt, err = rlp.EncodeToBytes(ethereumTxData.legacy)
	if err != nil {
		return []byte{}, err
	}

	return byt, nil
}

func (ethereumTxData *EthereumTransactionData) _GetData() []byte {
	if ethereumTxData.eip1559 != nil {
		return ethereumTxData.eip1559.Data
	}

	return ethereumTxData.legacy.Data
}

func (ethereumTxData *EthereumTransactionData) _SetData(data []byte) *EthereumTransactionData {
	if ethereumTxData.eip1559 != nil {
		ethereumTxData.eip1559.Data = data
		return ethereumTxData
	}

	ethereumTxData.legacy.Data = data
	return ethereumTxData
}

// ToJson returns a JSON representation of the Ethereum transaction.
func (ethereumTxData *EthereumTransactionData) ToJson() ([]byte, error) {
	var byt []byte
	var err error
	if ethereumTxData.eip1559 != nil {
		byt, err = json.Marshal(ethereumTxData.eip1559)
		if err != nil {
			return []byte{}, err
		}

		return byt, nil
	}

	byt, err = json.Marshal(ethereumTxData.legacy)
	if err != nil {
		return []byte{}, err
	}

	return byt, nil
}

// EthereumTransactionDataFromJson constructs an EthereumTransactionData from a JSON string.
func EthereumTransactionDataFromJson(b []byte) (*EthereumTransactionData, error) {
	var eip1559 types.DynamicFeeTx
	var leg types.LegacyTx
	err := json.Unmarshal(b, &eip1559)
	if err != nil {
		err = json.Unmarshal(b, &leg)
		if err != nil {
			return nil, errors.New("Json bytes are neither eip1559 or legacy format")
		}

		return &EthereumTransactionData{
			legacy: &leg,
		}, nil
	}

	return &EthereumTransactionData{
		eip1559: &eip1559,
	}, nil
}
