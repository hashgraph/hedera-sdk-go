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
	"time"

	"github.com/hashgraph/hedera-protobufs-go/services"
)

// TransactionRecordQuery
// Get the record for a transaction. If the transaction requested a record, then the record lasts
// for one hour, and a state proof is available for it. If the transaction created an account, file,
// or smart contract instance, then the record will contain the ID for what it created. If the
// transaction called a smart contract function, then the record contains the result of that call.
// If the transaction was a cryptocurrency transfer, then the record includes the TransferList which
// gives the details of that transfer. If the transaction didn't return anything that should be in
// the record, then the results field will be set to nothing.
type TransactionRecordQuery struct {
	query
	transactionID       *TransactionID
	includeChildRecords *bool
	duplicates          *bool
}

// NewTransactionRecordQuery creates TransactionRecordQuery which
// gets the record for a transaction. If the transaction requested a record, then the record lasts
// for one hour, and a state proof is available for it. If the transaction created an account, file,
// or smart contract instance, then the record will contain the ID for what it created. If the
// transaction called a smart contract function, then the record contains the result of that call.
// If the transaction was a cryptocurrency transfer, then the record includes the TransferList which
// gives the details of that transfer. If the transaction didn't return anything that should be in
// the record, then the results field will be set to nothing.
func NewTransactionRecordQuery() *TransactionRecordQuery {
	header := services.QueryHeader{}
	result := TransactionRecordQuery{
		query: _NewQuery(true, &header),
	}

	result.e = &result
	return &result
}

// When execution is attempted, a single attempt will timeout when this deadline is reached. (The SDK may subsequently retry the execution.)
func (q *TransactionRecordQuery) SetGrpcDeadline(deadline *time.Duration) *TransactionRecordQuery {
	q.query.SetGrpcDeadline(deadline)
	return q
}

// SetIncludeChildren Sets whether the response should include the records of any child transactions spawned by the
// top-level transaction with the given transactionID.
func (q *TransactionRecordQuery) SetIncludeChildren(includeChildRecords bool) *TransactionRecordQuery {
	q.includeChildRecords = &includeChildRecords
	return q
}

// GetIncludeChildren returns whether the response should include the records of any child transactions spawned by the
// top-level transaction with the given transactionID.
func (q *TransactionRecordQuery) GetIncludeChildren() bool {
	if q.includeChildRecords != nil {
		return *q.includeChildRecords
	}

	return false
}

// SetIncludeDuplicates Sets whether records of processing duplicate transactions should be returned along with the record
// of processing the first consensus transaction with the given id whose status was neither
// INVALID_NODE_ACCOUNT nor <tt>INVALID_PAYER_SIGNATURE; or, if no such
// record exists, the record of processing the first transaction to reach consensus with the
// given transaction id..
func (q *TransactionRecordQuery) SetIncludeDuplicates(includeDuplicates bool) *TransactionRecordQuery {
	q.duplicates = &includeDuplicates
	return q
}

// GetIncludeDuplicates returns whether records of processing duplicate transactions should be returned along with the record
// of processing the first consensus transaction with the given id.
func (q *TransactionRecordQuery) GetIncludeDuplicates() bool {
	if q.duplicates != nil {
		return *q.duplicates
	}

	return false
}

// Execute executes the Query with the provided client
func (q *TransactionRecordQuery) Execute(client *Client) (TransactionRecord, error) {
	resp, err := q.query.execute(client)

	if err != nil {
		if precheckErr, ok := err.(ErrHederaPreCheckStatus); ok {
			return TransactionRecord{}, _NewErrHederaReceiptStatus(precheckErr.TxID, precheckErr.Status)
		}
		return TransactionRecord{}, err
	}

	return _TransactionRecordFromProtobuf(resp.GetTransactionGetRecord(), q.transactionID), nil
}

// SetTransactionID sets the TransactionID for this TransactionRecordQuery.
func (q *TransactionRecordQuery) SetTransactionID(transactionID TransactionID) *TransactionRecordQuery {
	q.transactionID = &transactionID
	return q
}

// GetTransactionID returns the TransactionID for which the receipt is being requested.
func (q *TransactionRecordQuery) GetTransactionID() TransactionID {
	if q.transactionID == nil {
		return TransactionID{}
	}

	return *q.transactionID
}

// SetNodeAccountIDs sets the _Node AccountID for this TransactionRecordQuery.
func (q *TransactionRecordQuery) SetNodeAccountIDs(accountID []AccountID) *TransactionRecordQuery {
	q.query.SetNodeAccountIDs(accountID)
	return q
}

// SetQueryPayment sets the Hbar payment to pay the _Node a fee for handling this query
func (q *TransactionRecordQuery) SetQueryPayment(queryPayment Hbar) *TransactionRecordQuery {
	q.query.SetQueryPayment(queryPayment)
	return q
}

// SetMaxQueryPayment sets the maximum payment allowed for this Query.
func (q *TransactionRecordQuery) SetMaxQueryPayment(queryMaxPayment Hbar) *TransactionRecordQuery {
	q.query.SetMaxQueryPayment(queryMaxPayment)
	return q
}

// SetMaxRetry sets the max number of errors before execution will fail.
func (q *TransactionRecordQuery) SetMaxRetry(count int) *TransactionRecordQuery {
	q.query.SetMaxRetry(count)
	return q
}

// SetMaxBackoff The maximum amount of time to wait between retries.
// Every retry attempt will increase the wait time exponentially until it reaches this time.
func (q *TransactionRecordQuery) SetMaxBackoff(max time.Duration) *TransactionRecordQuery {
	q.query.SetMaxBackoff(max)
	return q
}

// SetMinBackoff sets the minimum amount of time to wait between retries.
func (q *TransactionRecordQuery) SetMinBackoff(min time.Duration) *TransactionRecordQuery {
	q.query.SetMinBackoff(min)
	return q
}

// SetPaymentTransactionID assigns the payment transaction id.
func (q *TransactionRecordQuery) SetPaymentTransactionID(transactionID TransactionID) *TransactionRecordQuery {
	q.paymentTransactionIDs._Clear()._Push(transactionID)._SetLocked(true)
	return q
}

func (q *TransactionRecordQuery) SetLogLevel(level LogLevel) *TransactionRecordQuery {
	q.query.SetLogLevel(level)
	return q
}

// ---------- Parent functions specific implementation ----------

func (q *TransactionRecordQuery) getMethod(channel *_Channel) _Method {
	return _Method{
		query: channel._GetCrypto().GetTxRecordByTxID,
	}
}

func (q *TransactionRecordQuery) mapStatusError(response interface{}) error {
	query := response.(*services.Response)
	switch Status(query.GetTransactionGetRecord().GetHeader().GetNodeTransactionPrecheckCode()) {
	case StatusPlatformTransactionNotCreated, StatusBusy, StatusUnknown, StatusReceiptNotFound, StatusRecordNotFound, StatusOk:
		break
	default:
		return ErrHederaPreCheckStatus{
			Status: Status(query.GetTransactionGetRecord().GetHeader().GetNodeTransactionPrecheckCode()),
		}
	}

	return ErrHederaReceiptStatus{
		Status: Status(query.GetTransactionGetRecord().GetTransactionRecord().GetReceipt().GetStatus()),
		// TxID:    _TransactionIDFromProtobuf(_Request.query.pb.GetTransactionGetRecord().TransactionID, networkName),
		Receipt: _TransactionReceiptFromProtobuf(query.GetTransactionGetReceipt(), nil),
	}
}

func (q *TransactionRecordQuery) getName() string {
	return "TransactionRecordQuery"
}

func (q *TransactionRecordQuery) buildQuery() *services.Query {
	body := &services.TransactionGetRecordQuery{
		Header: q.pbHeader,
	}

	if q.includeChildRecords != nil {
		body.IncludeChildRecords = q.GetIncludeChildren()
	}

	if q.duplicates != nil {
		body.IncludeDuplicates = q.GetIncludeDuplicates()
	}

	if q.transactionID.AccountID != nil {
		body.TransactionID = q.transactionID._ToProtobuf()
	}

	return &services.Query{
		Query: &services.Query_TransactionGetRecord{
			TransactionGetRecord: body,
		},
	}
}

func (q *TransactionRecordQuery) validateNetworkOnIDs(client *Client) error {
	if client == nil || !client.autoValidateChecksums {
		return nil
	}

	if err := q.transactionID.AccountID.ValidateChecksum(client); err != nil {
		return err
	}

	return nil
}

func (q *TransactionRecordQuery) shouldRetry(response interface{}) _ExecutionState {
	status := Status(response.(*services.Response).GetTransactionGetRecord().GetHeader().GetNodeTransactionPrecheckCode())

	switch status {
	case StatusPlatformTransactionNotCreated, StatusBusy, StatusUnknown, StatusReceiptNotFound, StatusRecordNotFound:
		return executionStateRetry
	case StatusOk:
		if response.(*services.Response).GetTransactionGetRecord().GetHeader().ResponseType == services.ResponseType_COST_ANSWER {
			return executionStateFinished
		}
	default:
		return executionStateError
	}

	status = Status(response.(*services.Response).GetTransactionGetRecord().GetTransactionRecord().GetReceipt().GetStatus())

	switch status {
	case StatusBusy, StatusUnknown, StatusOk, StatusReceiptNotFound, StatusRecordNotFound:
		return executionStateRetry
	case StatusSuccess:
		return executionStateFinished
	default:
		return executionStateError
	}
}

func (q *TransactionRecordQuery) getQueryResponse(response *services.Response) queryResponse {
	return response.GetTransactionGetRecord()
}
