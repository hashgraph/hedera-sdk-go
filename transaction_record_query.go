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
	"fmt"
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
	Query
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
	return &TransactionRecordQuery{
		Query: _NewQuery(true, &header),
	}
}

// When execution is attempted, a single attempt will timeout when this deadline is reached. (The SDK may subsequently retry the execution.)
func (query *TransactionRecordQuery) SetGrpcDeadline(deadline *time.Duration) *TransactionRecordQuery {
	query.Query.SetGrpcDeadline(deadline)
	return query
}

// SetIncludeChildren Sets whether the response should include the records of any child transactions spawned by the
// top-level transaction with the given transactionID.
func (query *TransactionRecordQuery) SetIncludeChildren(includeChildRecords bool) *TransactionRecordQuery {
	query.includeChildRecords = &includeChildRecords
	return query
}

// GetIncludeChildren returns whether the response should include the records of any child transactions spawned by the
// top-level transaction with the given transactionID.
func (query *TransactionRecordQuery) GetIncludeChildren() bool {
	if query.includeChildRecords != nil {
		return *query.includeChildRecords
	}

	return false
}

// SetIncludeDuplicates Sets whether records of processing duplicate transactions should be returned along with the record
// of processing the first consensus transaction with the given id whose status was neither
// INVALID_NODE_ACCOUNT nor <tt>INVALID_PAYER_SIGNATURE; or, if no such
// record exists, the record of processing the first transaction to reach consensus with the
// given transaction id..
func (query *TransactionRecordQuery) SetIncludeDuplicates(includeDuplicates bool) *TransactionRecordQuery {
	query.duplicates = &includeDuplicates
	return query
}

// GetIncludeDuplicates returns whether records of processing duplicate transactions should be returned along with the record
// of processing the first consensus transaction with the given id.
func (query *TransactionRecordQuery) GetIncludeDuplicates() bool {
	if query.duplicates != nil {
		return *query.duplicates
	}

	return false
}

func (query *TransactionRecordQuery) _ValidateNetworkOnIDs(client *Client) error {
	if client == nil || !client.autoValidateChecksums {
		return nil
	}

	if err := query.transactionID.AccountID.ValidateChecksum(client); err != nil {
		return err
	}

	return nil
}

func (query *TransactionRecordQuery) _Build() *services.Query_TransactionGetRecord {
	body := &services.TransactionGetRecordQuery{
		Header: &services.QueryHeader{},
	}

	if query.includeChildRecords != nil {
		body.IncludeChildRecords = query.GetIncludeChildren()
	}

	if query.duplicates != nil {
		body.IncludeDuplicates = query.GetIncludeDuplicates()
	}

	if query.transactionID.AccountID != nil {
		body.TransactionID = query.transactionID._ToProtobuf()
	}

	return &services.Query_TransactionGetRecord{
		TransactionGetRecord: body,
	}
}

// GetCost returns the fee that would be charged to get the requested information (if a cost was requested).
func (query *TransactionRecordQuery) GetCost(client *Client) (Hbar, error) {
	if client == nil || client.operator == nil {
		return Hbar{}, errNoClientProvided
	}

	var err error

	err = query._ValidateNetworkOnIDs(client)
	if err != nil {
		return Hbar{}, err
	}

	for range query.nodeAccountIDs.slice {
		paymentTransaction, err := _QueryMakePaymentTransaction(TransactionID{}, AccountID{}, client.operator, Hbar{})
		if err != nil {
			return Hbar{}, err
		}
		query.paymentTransactions = append(query.paymentTransactions, paymentTransaction)
	}

	pb := query._Build()
	pb.TransactionGetRecord.Header = query.pbHeader

	query.pb = &services.Query{
		Query: pb,
	}

	resp, err := _Execute(
		client,
		&query.Query,
		_TransactionRecordQueryShouldRetry,
		_CostQueryMakeRequest,
		_CostQueryAdvanceRequest,
		_QueryGetNodeAccountID,
		_TransactionRecordQueryGetMethod,
		_TransactionRecordQueryMapStatusError,
		_QueryMapResponse,
		query._GetLogID(),
		query.grpcDeadline,
		query.maxBackoff,
		query.minBackoff,
		query.maxRetry,
	)

	if err != nil {
		return Hbar{}, err
	}

	cost := int64(resp.(*services.Response).GetTransactionGetRecord().Header.Cost)
	return HbarFromTinybar(cost), nil
}

func _TransactionRecordQueryShouldRetry(request interface{}, response interface{}) _ExecutionState {
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

func _TransactionRecordQueryMapStatusError(request interface{}, response interface{}) error {
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

func _TransactionRecordQueryGetMethod(_ interface{}, channel *_Channel) _Method {
	return _Method{
		query: channel._GetCrypto().GetTxRecordByTxID,
	}
}

// SetTransactionID sets the TransactionID for this TransactionRecordQuery.
func (query *TransactionRecordQuery) SetTransactionID(transactionID TransactionID) *TransactionRecordQuery {
	query.transactionID = &transactionID
	return query
}

// GetTransactionID gets the TransactionID for this TransactionRecordQuery.
func (query *TransactionRecordQuery) GetTransactionID() TransactionID {
	if query.transactionID == nil {
		return TransactionID{}
	}

	return *query.transactionID
}

// SetNodeAccountIDs sets the _Node AccountID for this TransactionRecordQuery.
func (query *TransactionRecordQuery) SetNodeAccountIDs(accountID []AccountID) *TransactionRecordQuery {
	query.Query.SetNodeAccountIDs(accountID)
	return query
}

// SetQueryPayment sets the Hbar payment to pay the _Node a fee for handling this query
func (query *TransactionRecordQuery) SetQueryPayment(queryPayment Hbar) *TransactionRecordQuery {
	query.queryPayment = queryPayment
	return query
}

// SetMaxQueryPayment sets the maximum payment allowed for this Query.
func (query *TransactionRecordQuery) SetMaxQueryPayment(queryMaxPayment Hbar) *TransactionRecordQuery {
	query.maxQueryPayment = queryMaxPayment
	return query
}

// SetMaxRetry sets the max number of errors before execution will fail.
func (query *TransactionRecordQuery) SetMaxRetry(count int) *TransactionRecordQuery {
	query.Query.SetMaxRetry(count)
	return query
}

// SetMaxBackoff The maximum amount of time to wait between retries.
// Every retry attempt will increase the wait time exponentially until it reaches this time.
func (query *TransactionRecordQuery) SetMaxBackoff(max time.Duration) *TransactionRecordQuery {
	if max.Nanoseconds() < 0 {
		panic("maxBackoff must be a positive duration")
	} else if max.Nanoseconds() < query.minBackoff.Nanoseconds() {
		panic("maxBackoff must be greater than or equal to minBackoff")
	}
	query.maxBackoff = &max
	return query
}

// GetMaxBackoff returns the maximum amount of time to wait between retries.
func (query *TransactionRecordQuery) GetMaxBackoff() time.Duration {
	if query.maxBackoff != nil {
		return *query.maxBackoff
	}

	return 8 * time.Second
}

// SetMinBackoff sets the minimum amount of time to wait between retries.
func (query *TransactionRecordQuery) SetMinBackoff(min time.Duration) *TransactionRecordQuery {
	if min.Nanoseconds() < 0 {
		panic("minBackoff must be a positive duration")
	} else if query.maxBackoff.Nanoseconds() < min.Nanoseconds() {
		panic("minBackoff must be less than or equal to maxBackoff")
	}
	query.minBackoff = &min
	return query
}

// GetMinBackoff returns the minimum amount of time to wait between retries.
func (query *TransactionRecordQuery) GetMinBackoff() time.Duration {
	if query.minBackoff != nil {
		return *query.minBackoff
	}

	return 250 * time.Millisecond
}

// Execute executes the Query with the provided client
func (query *TransactionRecordQuery) Execute(client *Client) (TransactionRecord, error) {
	if client == nil || client.operator == nil {
		return TransactionRecord{}, errNoClientProvided
	}

	var err error

	err = query._ValidateNetworkOnIDs(client)
	if err != nil {
		return TransactionRecord{}, err
	}

	if !query.paymentTransactionIDs.locked {
		query.paymentTransactionIDs._Clear()._Push(TransactionIDGenerate(client.operator.accountID))
	}

	var cost Hbar
	if query.queryPayment.tinybar != 0 {
		cost = query.queryPayment
	} else {
		if query.maxQueryPayment.tinybar == 0 {
			cost = client.GetDefaultMaxQueryPayment()
		} else {
			cost = query.maxQueryPayment
		}

		actualCost, err := query.GetCost(client)
		if err != nil {
			return TransactionRecord{}, err
		}

		if cost.tinybar < actualCost.tinybar {
			return TransactionRecord{}, ErrMaxQueryPaymentExceeded{
				QueryCost:       actualCost,
				MaxQueryPayment: cost,
				query:           "TransactionRecordQuery",
			}
		}

		cost = actualCost
	}

	query.paymentTransactions = make([]*services.Transaction, 0)

	if query.nodeAccountIDs.locked {
		err = _QueryGeneratePayments(&query.Query, client, cost)
		if err != nil {
			return TransactionRecord{}, err
		}
	} else {
		paymentTransaction, err := _QueryMakePaymentTransaction(query.paymentTransactionIDs._GetCurrent().(TransactionID), AccountID{}, client.operator, cost)
		if err != nil {
			return TransactionRecord{}, err
		}
		query.paymentTransactions = append(query.paymentTransactions, paymentTransaction)
	}

	pb := query._Build()
	pb.TransactionGetRecord.Header = query.pbHeader
	query.pb = &services.Query{
		Query: pb,
	}

	resp, err := _Execute(
		client,
		&query.Query,
		_TransactionRecordQueryShouldRetry,
		_QueryMakeRequest,
		_QueryAdvanceRequest,
		_QueryGetNodeAccountID,
		_TransactionRecordQueryGetMethod,
		_TransactionRecordQueryMapStatusError,
		_QueryMapResponse,
		query._GetLogID(),
		query.grpcDeadline,
		query.maxBackoff,
		query.minBackoff,
		query.maxRetry,
	)

	if err != nil {
		if precheckErr, ok := err.(ErrHederaPreCheckStatus); ok {
			return TransactionRecord{}, _NewErrHederaReceiptStatus(precheckErr.TxID, precheckErr.Status)
		}
		return TransactionRecord{}, err
	}

	return _TransactionRecordFromProtobuf(resp.(*services.Response).GetTransactionGetRecord(), query.transactionID), nil
}

func (query *TransactionRecordQuery) _GetLogID() string {
	timestamp := query.timestamp.UnixNano()
	if query.paymentTransactionIDs._Length() > 0 && query.paymentTransactionIDs._GetCurrent().(TransactionID).ValidStart != nil {
		timestamp = query.paymentTransactionIDs._GetCurrent().(TransactionID).ValidStart.UnixNano()
	}
	return fmt.Sprintf("TransactionRecordQuery:%d", timestamp)
}

// SetPaymentTransactionID assigns the payment transaction id.
func (query *TransactionRecordQuery) SetPaymentTransactionID(transactionID TransactionID) *TransactionRecordQuery {
	query.paymentTransactionIDs._Clear()._Push(transactionID)._SetLocked(true)
	return query
}

func (query *TransactionRecordQuery) SetLogLevel(level LogLevel) *TransactionRecordQuery {
	query.Query.SetLogLevel(level)
	return query
}
