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
	query
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
	result := TransactionReceiptQuery{
		query: _NewQuery(false, &header),
	}
	result.e = &result
	return &result
}

// When execution is attempted, a single attempt will timeout when this deadline is reached. (The SDK may subsequently retry the execution.)
func (this *TransactionReceiptQuery) SetGrpcDeadline(deadline *time.Duration) *TransactionReceiptQuery {
	this.query.SetGrpcDeadline(deadline)
	return this
}

// SetIncludeChildren Sets whether the response should include the receipts of any child transactions spawned by the
// top-level transaction with the given transactionID.
func (this *TransactionReceiptQuery) SetIncludeChildren(includeChildReceipts bool) *TransactionReceiptQuery {
	this.childReceipts = &includeChildReceipts
	return this
}

// GetIncludeChildren returns whether the response should include the receipts of any child transactions spawned by the
// top-level transaction with the given transactionID.
func (this *TransactionReceiptQuery) GetIncludeChildren() bool {
	if this.childReceipts != nil {
		return *this.childReceipts
	}

	return false
}

// SetIncludeDuplicates Sets whether receipts of processing duplicate transactions should be returned along with the
// receipt of processing the first consensus transaction with the given id whose status was
// neither INVALID_NODE_ACCOUNT nor INVALID_PAYER_SIGNATURE; or, if no
// such receipt exists, the receipt of processing the first transaction to reach consensus with
// the given transaction id.
func (this *TransactionReceiptQuery) SetIncludeDuplicates(includeDuplicates bool) *TransactionReceiptQuery {
	this.duplicates = &includeDuplicates
	return this
}

// GetIncludeDuplicates returns whether receipts of processing duplicate transactions should be returned along with the
// receipt of processing the first consensus transaction with the given id
func (this *TransactionReceiptQuery) GetIncludeDuplicates() bool {
	if this.duplicates != nil {
		return *this.duplicates
	}

	return false
}
// GetCost returns the fee that would be charged to get the requested information (if a cost was requested).
func (this *TransactionReceiptQuery) GetCost(client *Client) (Hbar, error) {
	if client == nil {
		return Hbar{}, errNoClientProvided
	}

	var err error

	err = this.validateNetworkOnIDs(client)
	if err != nil {
		return Hbar{}, err
	}

	this.timestamp = time.Now()

	for range this.nodeAccountIDs.slice {
		paymentTransaction, err := _QueryMakePaymentTransaction(TransactionID{}, AccountID{}, client.operator, Hbar{})
		if err != nil {
			return Hbar{}, err
		}
		this.paymentTransactions = append(this.paymentTransactions, paymentTransaction)
	}

	pb := this.build()
	pb.TransactionGetReceipt.Header = this.pbHeader

	this.pb = &services.Query{
		Query: pb,
	}

	this.pbHeader.ResponseType = services.ResponseType_COST_ANSWER
	this.paymentTransactionIDs._Advance()

	resp, err := _Execute(
		client,
		this.e,
	)

	if err != nil {
		return Hbar{}, err
	}

	cost := int64(resp.(*services.Response).GetTransactionGetReceipt().Header.Cost)
	return HbarFromTinybar(cost), nil
}

// Execute executes the Query with the provided client
func (this *TransactionReceiptQuery) Execute(client *Client) (TransactionReceipt, error) {
	if client == nil {
		return TransactionReceipt{}, errNoClientProvided
	}

	var err error

	err = this.validateNetworkOnIDs(client)
	if err != nil {
		return TransactionReceipt{}, err
	}

	this.timestamp = time.Now()

	this.paymentTransactions = make([]*services.Transaction, 0)

	pb := this.build()
	pb.TransactionGetReceipt.Header = this.pbHeader
	this.pb = &services.Query{
		Query: pb,
	}

	if this.isPaymentRequired && len(this.paymentTransactions) > 0 {
		this.paymentTransactionIDs._Advance()
	}
	this.pbHeader.ResponseType = services.ResponseType_ANSWER_ONLY

	resp, err := _Execute(
		client,
		this.e,
	)

	if err, ok := err.(ErrHederaPreCheckStatus); ok {
		if resp.(*services.Response).GetTransactionGetReceipt() != nil {
			return _TransactionReceiptFromProtobuf(resp.(*services.Response).GetTransactionGetReceipt(), this.transactionID), err
		}
		// Manually add the receipt status, because an empty receipt has no status and no status defaults to 0, which means success
		return TransactionReceipt{Status: err.Status}, err
	}

	return _TransactionReceiptFromProtobuf(resp.(*services.Response).GetTransactionGetReceipt(), this.transactionID), nil
}


// SetTransactionID sets the TransactionID for which the receipt is being requested.
func (this *TransactionReceiptQuery) SetTransactionID(transactionID TransactionID) *TransactionReceiptQuery {
	this.transactionID = &transactionID
	return this
}

// GetTransactionID returns the TransactionID for which the receipt is being requested.
func (this *TransactionReceiptQuery) GetTransactionID() TransactionID {
	if this.transactionID == nil {
		return TransactionID{}
	}

	return *this.transactionID
}

// SetNodeAccountIDs sets the _Node AccountID for this TransactionReceiptQuery.
func (this *TransactionReceiptQuery) SetNodeAccountIDs(accountID []AccountID) *TransactionReceiptQuery {
	this.query.SetNodeAccountIDs(accountID)
	return this
}

// SetQueryPayment sets the Hbar payment to pay the _Node a fee for handling this query
func (this *TransactionReceiptQuery) SetQueryPayment(queryPayment Hbar) *TransactionReceiptQuery {
	this.queryPayment = queryPayment
	return this
}

// SetMaxQueryPayment sets the maximum payment allowed for this Query.
func (this *TransactionReceiptQuery) SetMaxQueryPayment(queryMaxPayment Hbar) *TransactionReceiptQuery {
	this.maxQueryPayment = queryMaxPayment
	return this
}

// SetMaxRetry sets the max number of errors before execution will fail.
func (this *TransactionReceiptQuery) SetMaxRetry(count int) *TransactionReceiptQuery {
	this.query.SetMaxRetry(count)
	return this
}

// SetMaxBackoff The maximum amount of time to wait between retries.
// Every retry attempt will increase the wait time exponentially until it reaches this time.
func (this *TransactionReceiptQuery) SetMaxBackoff(max time.Duration) *TransactionReceiptQuery {
	this.query.SetMaxBackoff(max)
	return this
}

// GetMaxBackoff returns the maximum amount of time to wait between retries.
func (this *TransactionReceiptQuery) GetMaxBackoff() time.Duration {
	return this.query.GetMaxBackoff()
}

// SetMinBackoff sets the minimum amount of time to wait between retries.
func (this *TransactionReceiptQuery) SetMinBackoff(min time.Duration) *TransactionReceiptQuery {
	this.query.SetMinBackoff(min)
	return this
}

// GetMinBackoff returns the minimum amount of time to wait between retries.
func (this *TransactionReceiptQuery) GetMinBackoff() time.Duration {
	return this.query.GetMinBackoff()
}

func (this *TransactionReceiptQuery) _GetLogID() string {
	timestamp := this.timestamp.UnixNano()
	return fmt.Sprintf("TransactionReceiptQuery:%d", timestamp)
}

// SetPaymentTransactionID assigns the payment transaction id.
func (this *TransactionReceiptQuery) SetPaymentTransactionID(transactionID TransactionID) *TransactionReceiptQuery {
	this.paymentTransactionIDs._Clear()._Push(transactionID)._SetLocked(true)
	return this
}

func (this *TransactionReceiptQuery) SetLogLevel(level LogLevel) *TransactionReceiptQuery {
	this.query.SetLogLevel(level)
	return this
}
// ---------- Parent functions specific implementation ----------

func (this *TransactionReceiptQuery) getMethod(channel *_Channel) _Method {
	return _Method{
		query: channel._GetCrypto().GetTransactionReceipts,
	}
}

func (this *TransactionReceiptQuery) mapStatusError(_ interface{}, response interface{}) error {
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

func (this *TransactionReceiptQuery) getName() string {
	return "TransactionReceiptQuery"
}

func (this *TransactionReceiptQuery) build()  *services.Query_TransactionGetReceipt {
	body := &services.TransactionGetReceiptQuery{
		Header: &services.QueryHeader{},
	}

	if this.transactionID.AccountID != nil {
		body.TransactionID = this.transactionID._ToProtobuf()
	}

	if this.duplicates != nil {
		body.IncludeDuplicates = *this.duplicates
	}

	if this.childReceipts != nil {
		body.IncludeChildReceipts = this.GetIncludeChildren()
	}

	return &services.Query_TransactionGetReceipt{
		TransactionGetReceipt: body,
	}
}

func (this *TransactionReceiptQuery) validateNetworkOnIDs(client *Client) error {
	if client == nil || !client.autoValidateChecksums {
		return nil
	}

	if err := this.transactionID.AccountID.ValidateChecksum(client); err != nil {
		return err
	}

	return nil
}

func (this *TransactionReceiptQuery) getQueryStatus(response interface{}) Status {
	return  Status(response.(*services.Response).GetTransactionGetReceipt().GetHeader().GetNodeTransactionPrecheckCode())
}

func (this *TransactionReceiptQuery) shouldRetry(_ interface{}, response interface{}) _ExecutionState {
	status := this.getQueryStatus(response)

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
