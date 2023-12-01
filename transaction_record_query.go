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
func (this *TransactionRecordQuery) SetGrpcDeadline(deadline *time.Duration) *TransactionRecordQuery {
	this.query.SetGrpcDeadline(deadline)
	return this
}

// SetIncludeChildren Sets whether the response should include the records of any child transactions spawned by the
// top-level transaction with the given transactionID.
func (this *TransactionRecordQuery) SetIncludeChildren(includeChildRecords bool) *TransactionRecordQuery {
	this.includeChildRecords = &includeChildRecords
	return this
}

// GetIncludeChildren returns whether the response should include the records of any child transactions spawned by the
// top-level transaction with the given transactionID.
func (this *TransactionRecordQuery) GetIncludeChildren() bool {
	if this.includeChildRecords != nil {
		return *this.includeChildRecords
	}

	return false
}

// SetIncludeDuplicates Sets whether records of processing duplicate transactions should be returned along with the record
// of processing the first consensus transaction with the given id whose status was neither
// INVALID_NODE_ACCOUNT nor <tt>INVALID_PAYER_SIGNATURE; or, if no such
// record exists, the record of processing the first transaction to reach consensus with the
// given transaction id..
func (this *TransactionRecordQuery) SetIncludeDuplicates(includeDuplicates bool) *TransactionRecordQuery {
	this.duplicates = &includeDuplicates
	return this
}

// GetIncludeDuplicates returns whether records of processing duplicate transactions should be returned along with the record
// of processing the first consensus transaction with the given id.
func (this *TransactionRecordQuery) GetIncludeDuplicates() bool {
	if this.duplicates != nil {
		return *this.duplicates
	}

	return false
}

// GetCost returns the fee that would be charged to get the requested information (if a cost was requested).
func (this *TransactionRecordQuery) GetCost(client *Client) (Hbar, error) {
	if client == nil || client.operator == nil {
		return Hbar{}, errNoClientProvided
	}

	var err error

	err = this.validateNetworkOnIDs(client)
	if err != nil {
		return Hbar{}, err
	}

	for range this.nodeAccountIDs.slice {
		paymentTransaction, err := _QueryMakePaymentTransaction(TransactionID{}, AccountID{}, client.operator, Hbar{})
		if err != nil {
			return Hbar{}, err
		}
		this.paymentTransactions = append(this.paymentTransactions, paymentTransaction)
	}

	pb := this.build()
	pb.TransactionGetRecord.Header = this.pbHeader

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

	cost := int64(resp.(*services.Response).GetTransactionGetRecord().Header.Cost)
	return HbarFromTinybar(cost), nil
}

// Execute executes the Query with the provided client
func (this *TransactionRecordQuery) Execute(client *Client) (TransactionRecord, error) {
	if client == nil || client.operator == nil {
		return TransactionRecord{}, errNoClientProvided
	}

	var err error

	err = this.validateNetworkOnIDs(client)
	if err != nil {
		return TransactionRecord{}, err
	}

	if !this.paymentTransactionIDs.locked {
		this.paymentTransactionIDs._Clear()._Push(TransactionIDGenerate(client.operator.accountID))
	}

	var cost Hbar
	if this.queryPayment.tinybar != 0 {
		cost = this.queryPayment
	} else {
		if this.maxQueryPayment.tinybar == 0 {
			cost = client.GetDefaultMaxQueryPayment()
		} else {
			cost = this.maxQueryPayment
		}

		actualCost, err := this.GetCost(client)
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

	this.paymentTransactions = make([]*services.Transaction, 0)

	if this.nodeAccountIDs.locked {
		err = this._QueryGeneratePayments(client, cost)
		if err != nil {
			return TransactionRecord{}, err
		}
	} else {
		paymentTransaction, err := _QueryMakePaymentTransaction(this.paymentTransactionIDs._GetCurrent().(TransactionID), AccountID{}, client.operator, cost)
		if err != nil {
			return TransactionRecord{}, err
		}
		this.paymentTransactions = append(this.paymentTransactions, paymentTransaction)
	}

	pb := this.build()
	pb.TransactionGetRecord.Header = this.pbHeader
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

	if err != nil {
		if precheckErr, ok := err.(ErrHederaPreCheckStatus); ok {
			return TransactionRecord{}, _NewErrHederaReceiptStatus(precheckErr.TxID, precheckErr.Status)
		}
		return TransactionRecord{}, err
	}

	return _TransactionRecordFromProtobuf(resp.(*services.Response).GetTransactionGetRecord(), this.transactionID), nil
}

// SetTransactionID sets the TransactionID for this TransactionRecordQuery.
func (this *TransactionRecordQuery) SetTransactionID(transactionID TransactionID) *TransactionRecordQuery {
	this.transactionID = &transactionID
	return this
}

// GetTransactionID gets the TransactionID for this TransactionRecordQuery.
func (this *TransactionRecordQuery) GetTransactionID() TransactionID {
	if this.transactionID == nil {
		return TransactionID{}
	}

	return *this.transactionID
}

// SetNodeAccountIDs sets the _Node AccountID for this TransactionRecordQuery.
func (this *TransactionRecordQuery) SetNodeAccountIDs(accountID []AccountID) *TransactionRecordQuery {
	this.query.SetNodeAccountIDs(accountID)
	return this
}

// SetQueryPayment sets the Hbar payment to pay the _Node a fee for handling this query
func (this *TransactionRecordQuery) SetQueryPayment(queryPayment Hbar) *TransactionRecordQuery {
	this.queryPayment = queryPayment
	return this
}

// SetMaxQueryPayment sets the maximum payment allowed for this Query.
func (this *TransactionRecordQuery) SetMaxQueryPayment(queryMaxPayment Hbar) *TransactionRecordQuery {
	this.maxQueryPayment = queryMaxPayment
	return this
}

// SetMaxRetry sets the max number of errors before execution will fail.
func (this *TransactionRecordQuery) SetMaxRetry(count int) *TransactionRecordQuery {
	this.query.SetMaxRetry(count)
	return this
}

// SetMaxBackoff The maximum amount of time to wait between retries.
// Every retry attempt will increase the wait time exponentially until it reaches this time.
func (this *TransactionRecordQuery) SetMaxBackoff(max time.Duration) *TransactionRecordQuery {
	this.query.SetMaxBackoff(max)
	return this
}

// GetMaxBackoff returns the maximum amount of time to wait between retries.
func (this *TransactionRecordQuery) GetMaxBackoff() time.Duration {
	return this.query.GetMaxBackoff()
}

// SetMinBackoff sets the minimum amount of time to wait between retries.
func (this *TransactionRecordQuery) SetMinBackoff(min time.Duration) *TransactionRecordQuery {
	this.query.SetMinBackoff(min)
	return this
}

// GetMinBackoff returns the minimum amount of time to wait between retries.
func (this *TransactionRecordQuery) GetMinBackoff() time.Duration {
	return this.query.GetMinBackoff()
}

func (this *TransactionRecordQuery) _GetLogID() string {
	timestamp := this.timestamp.UnixNano()
	if this.paymentTransactionIDs._Length() > 0 && this.paymentTransactionIDs._GetCurrent().(TransactionID).ValidStart != nil {
		timestamp = this.paymentTransactionIDs._GetCurrent().(TransactionID).ValidStart.UnixNano()
	}
	return fmt.Sprintf("TransactionRecordQuery:%d", timestamp)
}

// SetPaymentTransactionID assigns the payment transaction id.
func (this *TransactionRecordQuery) SetPaymentTransactionID(transactionID TransactionID) *TransactionRecordQuery {
	this.paymentTransactionIDs._Clear()._Push(transactionID)._SetLocked(true)
	return this
}

func (this *TransactionRecordQuery) SetLogLevel(level LogLevel) *TransactionRecordQuery {
	this.query.SetLogLevel(level)
	return this
}

// ---------- Parent functions specific implementation ----------

func (this *TransactionRecordQuery) getMethod(channel *_Channel) _Method {
	return _Method{
		query: channel._GetCrypto().GetTxRecordByTxID,
	}
}

func (this *TransactionRecordQuery) mapStatusError(_ interface{}, response interface{}) error {
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

func (this *TransactionRecordQuery) getName() string {
	return "TransactionRecordQuery"
}

func (this *TransactionRecordQuery) build() *services.Query_TransactionGetRecord {
	body := &services.TransactionGetRecordQuery{
		Header: &services.QueryHeader{},
	}

	if this.includeChildRecords != nil {
		body.IncludeChildRecords = this.GetIncludeChildren()
	}

	if this.duplicates != nil {
		body.IncludeDuplicates = this.GetIncludeDuplicates()
	}

	if this.transactionID.AccountID != nil {
		body.TransactionID = this.transactionID._ToProtobuf()
	}

	return &services.Query_TransactionGetRecord{
		TransactionGetRecord: body,
	}
}

func (this *TransactionRecordQuery) validateNetworkOnIDs(client *Client) error {
	if client == nil || !client.autoValidateChecksums {
		return nil
	}

	if err := this.transactionID.AccountID.ValidateChecksum(client); err != nil {
		return err
	}

	return nil
}

func (q *ContractInfoQuery) shouldRetry(_ interface{}, response interface{}) _ExecutionState {
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
