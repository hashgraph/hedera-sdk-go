package hedera

/*-
 *
 * Hedera Go SDK
 *
 * Copyright (C) 2020 - 2024 Hedera Hashgraph, LLC
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

// Represents the data of an Ethereum transaction.
type EthereumTransactionData struct {
	eip1559 *EthereumEIP1559Transaction
	legacy  *EthereumLegacyTransaction
}

// EthereumTransactionDataFromBytes constructs an EthereumTransactionData from a raw byte array.
func EthereumTransactionDataFromBytes(b []byte) (*EthereumTransactionData, error) {
	var transactionData EthereumTransactionData
	if b[0] == 0x02 {
		eip1559, err := EthereumEIP1559TransactionFromBytes(b)
		if err != nil {
			return nil, err
		}

		transactionData.eip1559 = eip1559
		return &transactionData, nil
	}

	legacy, err := EthereumLegacyTransactionFromBytes(b)
	if err != nil {
		return nil, err
	}

	transactionData.legacy = legacy
	return &transactionData, nil
}

// ToBytes returns the raw bytes of the Ethereum transaction.
func (txData *EthereumTransactionData) ToBytes() ([]byte, error) {
	if txData.eip1559 != nil {
		return txData.eip1559.ToBytes()
	}

	if txData.legacy != nil {
		return txData.legacy.ToBytes()
	}

	return nil, nil
}

func (ethereumTxData *EthereumTransactionData) _GetData() []byte {
	if ethereumTxData.eip1559 != nil {
		return ethereumTxData.eip1559.CallData
	}

	return ethereumTxData.legacy.CallData
}

func (ethereumTxData *EthereumTransactionData) _SetData(data []byte) *EthereumTransactionData {
	if ethereumTxData.eip1559 != nil {
		ethereumTxData.eip1559.CallData = data
		return ethereumTxData
	}

	ethereumTxData.legacy.CallData = data
	return ethereumTxData
}
