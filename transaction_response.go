package hedera

import (
	"encoding/hex"

	jsoniter "github.com/json-iterator/go"
)

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

// When the client sends the node a transaction of any kind, the node replies with this, which
// simply says that the transaction passed the precheck (so the node will submit it to the network)
// or it failed (so it won't). If the fee offered was insufficient, this will also contain the
// amount of the required fee. To learn the consensus result, the client should later obtain a
// receipt (free), or can buy a more detailed record (not free).
type TransactionResponse struct {
	TransactionID          TransactionID
	ScheduledTransactionId TransactionID // nolint
	NodeID                 AccountID
	Hash                   []byte
	ValidateStatus         bool
}

// MarshalJSON returns the JSON representation of the TransactionResponse.
// This should yield the same result in all SDK's.
func (response TransactionResponse) MarshalJSON() ([]byte, error) {
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	obj := make(map[string]interface{})
	obj["nodeID"] = response.NodeID.String()
	obj["hash"] = hex.EncodeToString(response.Hash)
	obj["transactionID"] = response.TransactionID.String()
	return json.Marshal(obj)
}

// GetReceipt retrieves the receipt for the transaction
func (response TransactionResponse) GetReceipt(client *Client) (TransactionReceipt, error) {
	receipt, err := NewTransactionReceiptQuery().
		SetTransactionID(response.TransactionID).
		SetNodeAccountIDs([]AccountID{response.NodeID}).
		Execute(client)

	if err != nil {
		return receipt, err
	}

	return receipt, receipt.ValidateStatus(response.ValidateStatus)
}

// GetRecord retrieves the record for the transaction
func (response TransactionResponse) GetRecord(client *Client) (TransactionRecord, error) {
	receipt, err := NewTransactionReceiptQuery().
		SetTransactionID(response.TransactionID).
		SetNodeAccountIDs([]AccountID{response.NodeID}).
		Execute(client)

	if err != nil {
		// Manually add the receipt, because an empty TransactionRecord will have an empty receipt and empty receipt has no status and no status defaults to 0, which means success
		return TransactionRecord{Receipt: receipt}, err
	}

	return NewTransactionRecordQuery().
		SetTransactionID(response.TransactionID).
		SetNodeAccountIDs([]AccountID{response.NodeID}).
		Execute(client)
}

// GetReceiptQuery retrieves the receipt query for the transaction
func (response TransactionResponse) GetReceiptQuery() *TransactionReceiptQuery {
	return NewTransactionReceiptQuery().
		SetTransactionID(response.TransactionID).
		SetNodeAccountIDs([]AccountID{response.NodeID})
}

// GetRecordQuery retrieves the record query for the transaction
func (response TransactionResponse) GetRecordQuery() *TransactionRecordQuery {
	return NewTransactionRecordQuery().
		SetTransactionID(response.TransactionID).
		SetNodeAccountIDs([]AccountID{response.NodeID})
}

// SetValidateStatus sets the validate status for the transaction
func (response TransactionResponse) SetValidateStatus(validate bool) *TransactionResponse {
	response.ValidateStatus = validate
	return &response
}

// GetValidateStatus returns the validate status for the transaction
func (response TransactionResponse) GetValidateStatus() bool {
	return response.ValidateStatus
}
