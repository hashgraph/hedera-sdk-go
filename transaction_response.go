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

type TransactionResponse struct {
	TransactionID          TransactionID
	ScheduledTransactionId TransactionID // nolint
	NodeID                 AccountID
	Hash                   []byte
	ValidateStatus         bool
}

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

func (response TransactionResponse) GetReceiptQuery() *TransactionReceiptQuery {
	return NewTransactionReceiptQuery().
		SetTransactionID(response.TransactionID).
		SetNodeAccountIDs([]AccountID{response.NodeID})
}

func (response TransactionResponse) GetRecordQuery() *TransactionRecordQuery {
	return NewTransactionRecordQuery().
		SetTransactionID(response.TransactionID).
		SetNodeAccountIDs([]AccountID{response.NodeID})
}

func (response TransactionResponse) SetValidateStatus(validate bool) *TransactionResponse {
	response.ValidateStatus = validate
	return &response
}

func (response TransactionResponse) GetValidateStatus() bool {
	return response.ValidateStatus
}
