package hedera

import (
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
)

type EthereumTransactionData struct {
	eip1559 *types.DynamicFeeTx
	legacy  *types.LegacyTx
}

func EthereumTransactionDataFromBytes(b []byte) (*EthereumTransactionData, error) {
	var transactionData EthereumTransactionData
	if b[0] == 2 {
		byt := append(b[:0], b[0+1:]...)
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
