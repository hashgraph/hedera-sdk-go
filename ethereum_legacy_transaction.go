package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"encoding/hex"
	"fmt"

	"github.com/pkg/errors"
)

// EthereumLegacyTransaction represents the legacy Ethereum transaction data.
type EthereumLegacyTransaction struct {
	Nonce    []byte
	GasPrice []byte
	GasLimit []byte
	To       []byte
	Value    []byte
	CallData []byte
	V        []byte
	R        []byte
	S        []byte
}

// nolint
// NewEthereumLegacyTransaction creates a new EthereumLegacyTransaction with the provided fields.
func NewEthereumLegacyTransaction(nonce, gasPrice, gasLimit, to, value, callData, v, r, s []byte) *EthereumLegacyTransaction {
	return &EthereumLegacyTransaction{
		Nonce:    nonce,
		GasPrice: gasPrice,
		GasLimit: gasLimit,
		To:       to,
		Value:    value,
		CallData: callData,
		V:        v,
		R:        r,
		S:        s,
	}
}

// FromBytes decodes the RLP encoded bytes into an EthereumLegacyTransaction.
func EthereumLegacyTransactionFromBytes(bytes []byte) (*EthereumLegacyTransaction, error) {
	item := NewRLPItem(LIST_TYPE)
	if err := item.Read(bytes); err != nil {
		return nil, errors.Wrap(err, "failed to read RLP data")
	}

	if item.itemType != LIST_TYPE {
		return nil, errors.New("input byte array does not represent a list of RLP-encoded elements")
	}

	if len(item.childItems) != 9 {
		return nil, errors.New("input byte array does not contain 9 RLP-encoded elements")
	}

	// Extract the values from the RLP item
	return NewEthereumLegacyTransaction(
		item.childItems[0].itemValue,
		item.childItems[1].itemValue,
		item.childItems[2].itemValue,
		item.childItems[3].itemValue,
		item.childItems[4].itemValue,
		item.childItems[5].itemValue,
		item.childItems[6].itemValue,
		item.childItems[7].itemValue,
		item.childItems[8].itemValue,
	), nil
}

// ToBytes encodes the EthereumLegacyTransaction into RLP format.
func (txn *EthereumLegacyTransaction) ToBytes() ([]byte, error) {
	item := NewRLPItem(LIST_TYPE)
	item.PushBack(NewRLPItem(VALUE_TYPE).AssignValue(txn.Nonce))
	item.PushBack(NewRLPItem(VALUE_TYPE).AssignValue(txn.GasPrice))
	item.PushBack(NewRLPItem(VALUE_TYPE).AssignValue(txn.GasLimit))
	item.PushBack(NewRLPItem(VALUE_TYPE).AssignValue(txn.To))
	item.PushBack(NewRLPItem(VALUE_TYPE).AssignValue(txn.Value))
	item.PushBack(NewRLPItem(VALUE_TYPE).AssignValue(txn.CallData))
	item.PushBack(NewRLPItem(VALUE_TYPE).AssignValue(txn.V))
	item.PushBack(NewRLPItem(VALUE_TYPE).AssignValue(txn.R))
	item.PushBack(NewRLPItem(VALUE_TYPE).AssignValue(txn.S))

	return item.Write()
}

// String returns a string representation of the EthereumLegacyTransaction.
func (txn *EthereumLegacyTransaction) String() string {
	return fmt.Sprintf("Nonce: %s\nGasPrice: %s\nGasLimit: %s\nTo: %s\nValue: %s\nCallData: %s\nV: %s\nR: %s\nS: %s",
		hex.EncodeToString(txn.Nonce),
		hex.EncodeToString(txn.GasPrice),
		hex.EncodeToString(txn.GasLimit),
		hex.EncodeToString(txn.To),
		hex.EncodeToString(txn.Value),
		hex.EncodeToString(txn.CallData),
		hex.EncodeToString(txn.V),
		hex.EncodeToString(txn.R),
		hex.EncodeToString(txn.S),
	)
}
