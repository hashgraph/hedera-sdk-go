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

// TransactionReceiptQuery
// Get the receipt of a transaction, given its transaction ID. Once a transaction reaches consensus,
// then information about whether it succeeded or failed will be available until the end of the
// receipt period.  Before and after the receipt period, and for a transaction that was never
// submitted, the receipt is unknown.  This query is free (the payment field is left empty). No
// State proof is available for this response
type TransactionReceiptQuery struct {
	Query
	transactionID *TransactionID
	childReceipts *bool
	duplicates    *bool
	timestamp     time.Time
}

// NewTransactionReceiptQuery creates TransactionReceiptQuery which
// gets the receipt of a transaction, given its transaction ID. Once a transaction reaches consensus,
// then information about whether it succeeded or failed will be available until the end of the
// receipt period.  Before and after the receipt period, and for a transaction that was never
// submitted, the receipt is unknown.  This query is free (the payment field is left empty). No
// State proof is available for this response
func NewTransactionReceiptQuery() *TransactionReceiptQuery {
	header := services.QueryHeader{}
	return &TransactionReceiptQuery{
		Query: _NewQuery(false, &header),
	}
}

// When execution is attempted, a single attempt will timeout when this deadline is reached. (The SDK may subsequently retry the execution.)
func (query *TransactionReceiptQuery) SetGrpcDeadline(deadline *time.Duration) *TransactionReceiptQuery {
	query.Query.SetGrpcDeadline(deadline)
	return query
}

// SetIncludeChildren Sets whether the response should include the receipts of any child transactions spawned by the
// top-level transaction with the given transactionID.
func (query *TransactionReceiptQuery) SetIncludeChildren(includeChildReceipts bool) *TransactionReceiptQuery {
	query.childReceipts = &includeChildReceipts
	return query
}

// GetIncludeChildren returns whether the response should include the receipts of any child transactions spawned by the
// top-level transaction with the given transactionID.
func (query *TransactionReceiptQuery) GetIncludeChildren() bool {
	if query.childReceipts != nil {
		return *query.childReceipts
	}

	return false
}

// SetIncludeDuplicates Sets whether receipts of processing duplicate transactions should be returned along with the
// receipt of processing the first consensus transaction with the given id whose status was
// neither INVALID_NODE_ACCOUNT nor INVALID_PAYER_SIGNATURE; or, if no
// such receipt exists, the receipt of processing the first transaction to reach consensus with
// the given transaction id.
func (query *TransactionReceiptQuery) SetIncludeDuplicates(includeDuplicates bool) *TransactionReceiptQuery {
	query.duplicates = &includeDuplicates
	return query
}

// GetIncludeDuplicates returns whether receipts of processing duplicate transactions should be returned along with the
// receipt of processing the first consensus transaction with the given id
func (query *TransactionReceiptQuery) GetIncludeDuplicates() bool {
	if query.duplicates != nil {
		return *query.duplicates
	}

	return false
}

func (query *TransactionReceiptQuery) _ValidateNetworkOnIDs(client *Client) error {
	if client == nil || !client.autoValidateChecksums {
		return nil
	}

	if err := query.transactionID.AccountID.ValidateChecksum(client); err != nil {
		return err
	}

	return nil
}

func (query *TransactionReceiptQuery) _Build() *services.Query_TransactionGetReceipt {
	body := &services.TransactionGetReceiptQuery{
		Header: &services.QueryHeader{},
	}

	if query.transactionID.AccountID != nil {
		body.TransactionID = query.transactionID._ToProtobuf()
	}

	if query.duplicates != nil {
		body.IncludeDuplicates = *query.duplicates
	}

	if query.childReceipts != nil {
		body.IncludeChildReceipts = query.GetIncludeChildren()
	}

	return &services.Query_TransactionGetReceipt{
		TransactionGetReceipt: body,
	}
}

// GetCost returns the fee that would be charged to get the requested information (if a cost was requested).
func (query *TransactionReceiptQuery) GetCost(client *Client) (Hbar, error) {
	if client == nil {
		return Hbar{}, errNoClientProvided
	}

	var err error

	err = query._ValidateNetworkOnIDs(client)
	if err != nil {
		return Hbar{}, err
	}

	query.timestamp = time.Now()

	for range query.nodeAccountIDs.slice {
		paymentTransaction, err := _QueryMakePaymentTransaction(TransactionID{}, AccountID{}, client.operator, Hbar{})
		if err != nil {
			return Hbar{}, err
		}
		query.paymentTransactions = append(query.paymentTransactions, paymentTransaction)
	}

	pb := query._Build()
	pb.TransactionGetReceipt.Header = query.pbHeader

	query.pb = &services.Query{
		Query: pb,
	}

	resp, err := _Execute(
		client,
		&query.Query,
		_TransactionReceiptQueryShouldRetry,
		_CostQueryMakeRequest,
		_CostQueryAdvanceRequest,
		_QueryGetNodeAccountID,
		_TransactionReceiptQueryGetMethod,
		_TransactionReceiptQueryMapStatusError,
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

	cost := int64(resp.(*services.Response).GetTransactionGetReceipt().Header.Cost)
	return HbarFromTinybar(cost), nil
}

func _TransactionReceiptQueryShouldRetry(request interface{}, response interface{}) _ExecutionState {
	status := Status(response.(*services.Response).GetTransactionGetReceipt().GetHeader().GetNodeTransactionPrecheckCode())

	switch status {
	case StatusPlatformTransactionNotCreated, StatusBusy, StatusUnknown, StatusReceiptNotFound, StatusRecordNotFound:
		return executionStateRetry
	case StatusOk:
		break
	default:
		return executionStateError
	}

	status = Status(response.(*services.Response).GetTransactionGetReceipt().GetReceipt().GetStatus())

	switch status {
	case StatusBusy, StatusUnknown, StatusOk, StatusReceiptNotFound, StatusRecordNotFound:
		return executionStateRetry
	default:
		return executionStateFinished
	}
}

func _TransactionReceiptQueryMapStatusError(request interface{}, response interface{}) error {
	switch Status(response.(*services.Response).GetTransactionGetReceipt().GetHeader().GetNodeTransactionPrecheckCode()) {
	case StatusPlatformTransactionNotCreated, StatusBusy, StatusUnknown, StatusOk:
		break
	default:
		return ErrHederaPreCheckStatus{
			Status: Status(response.(*services.Response).GetTransactionGetReceipt().GetHeader().GetNodeTransactionPrecheckCode()),
		}
	}

	return ErrHederaPreCheckStatus{
		Status: Status(response.(*services.Response).GetTransactionGetReceipt().GetReceipt().GetStatus()),
	}
}

func _TransactionReceiptQueryGetMethod(_ interface{}, channel *_Channel) _Method {
	return _Method{
		query: channel._GetCrypto().GetTransactionReceipts,
	}
}

// SetTransactionID sets the TransactionID for which the receipt is being requested.
func (query *TransactionReceiptQuery) SetTransactionID(transactionID TransactionID) *TransactionReceiptQuery {
	query.transactionID = &transactionID
	return query
}

// GetTransactionID returns the TransactionID for which the receipt is being requested.
func (query *TransactionReceiptQuery) GetTransactionID() TransactionID {
	if query.transactionID == nil {
		return TransactionID{}
	}

	return *query.transactionID
}

// SetNodeAccountIDs sets the _Node AccountID for this TransactionReceiptQuery.
func (query *TransactionReceiptQuery) SetNodeAccountIDs(accountID []AccountID) *TransactionReceiptQuery {
	query.Query.SetNodeAccountIDs(accountID)
	return query
}

// SetQueryPayment sets the Hbar payment to pay the _Node a fee for handling this query
func (query *TransactionReceiptQuery) SetQueryPayment(queryPayment Hbar) *TransactionReceiptQuery {
	query.queryPayment = queryPayment
	return query
}

// SetMaxQueryPayment sets the maximum payment allowed for this Query.
func (query *TransactionReceiptQuery) SetMaxQueryPayment(queryMaxPayment Hbar) *TransactionReceiptQuery {
	query.maxQueryPayment = queryMaxPayment
	return query
}

// SetMaxRetry sets the max number of errors before execution will fail.
func (query *TransactionReceiptQuery) SetMaxRetry(count int) *TransactionReceiptQuery {
	query.Query.SetMaxRetry(count)
	return query
}

// SetMaxBackoff The maximum amount of time to wait between retries.
// Every retry attempt will increase the wait time exponentially until it reaches this time.
func (query *TransactionReceiptQuery) SetMaxBackoff(max time.Duration) *TransactionReceiptQuery {
	if max.Nanoseconds() < 0 {
		panic("maxBackoff must be a positive duration")
	} else if max.Nanoseconds() < query.minBackoff.Nanoseconds() {
		panic("maxBackoff must be greater than or equal to minBackoff")
	}
	query.maxBackoff = &max
	return query
}

// GetMaxBackoff returns the maximum amount of time to wait between retries.
func (query *TransactionReceiptQuery) GetMaxBackoff() time.Duration {
	if query.maxBackoff != nil {
		return *query.maxBackoff
	}

	return 8 * time.Second
}

// SetMinBackoff sets the minimum amount of time to wait between retries.
func (query *TransactionReceiptQuery) SetMinBackoff(min time.Duration) *TransactionReceiptQuery {
	if min.Nanoseconds() < 0 {
		panic("minBackoff must be a positive duration")
	} else if query.maxBackoff.Nanoseconds() < min.Nanoseconds() {
		panic("minBackoff must be less than or equal to maxBackoff")
	}
	query.minBackoff = &min
	return query
}

// GetMinBackoff returns the minimum amount of time to wait between retries.
func (query *TransactionReceiptQuery) GetMinBackoff() time.Duration {
	if query.minBackoff != nil {
		return *query.minBackoff
	}

	return 250 * time.Millisecond
}

// Execute executes the Query with the provided client
func (query *TransactionReceiptQuery) Execute(client *Client) (TransactionReceipt, error) {
	if client == nil {
		return TransactionReceipt{}, errNoClientProvided
	}

	var err error

	err = query._ValidateNetworkOnIDs(client)
	if err != nil {
		return TransactionReceipt{}, err
	}

	query.timestamp = time.Now()

	query.paymentTransactions = make([]*services.Transaction, 0)

	pb := query._Build()
	pb.TransactionGetReceipt.Header = query.pbHeader
	query.pb = &services.Query{
		Query: pb,
	}

	resp, err := _Execute(
		client,
		&query.Query,
		_TransactionReceiptQueryShouldRetry,
		_QueryMakeRequest,
		_QueryAdvanceRequest,
		_QueryGetNodeAccountID,
		_TransactionReceiptQueryGetMethod,
		_TransactionReceiptQueryMapStatusError,
		_QueryMapResponse,
		query._GetLogID(),
		query.grpcDeadline,
		query.maxBackoff,
		query.minBackoff,
		query.maxRetry,
	)

	if err, ok := err.(ErrHederaPreCheckStatus); ok {
		if resp.(*services.Response).GetTransactionGetReceipt() != nil {
			return _TransactionReceiptFromProtobuf(resp.(*services.Response).GetTransactionGetReceipt(), query.transactionID), err
		}
		// Manually add the receipt status, because an empty receipt has no status and no status defaults to 0, which means success
		return TransactionReceipt{Status: err.Status}, err
	}

	return _TransactionReceiptFromProtobuf(resp.(*services.Response).GetTransactionGetReceipt(), query.transactionID), nil
}

func (query *TransactionReceiptQuery) _GetLogID() string {
	timestamp := query.timestamp.UnixNano()
	return fmt.Sprintf("TransactionReceiptQuery:%d", timestamp)
}

// SetPaymentTransactionID assigns the payment transaction id.
func (query *TransactionReceiptQuery) SetPaymentTransactionID(transactionID TransactionID) *TransactionReceiptQuery {
	query.paymentTransactionIDs._Clear()._Push(transactionID)._SetLocked(true)
	return query
}

func (query *TransactionReceiptQuery) SetLogLevel(level LogLevel) *TransactionReceiptQuery {
	query.Query.SetLogLevel(level)
	return query
}
