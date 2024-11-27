package hiero

// SPDX-License-Identifier: Apache-2.0

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
