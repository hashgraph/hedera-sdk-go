package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"time"

	"github.com/hiero-ledger/hiero-sdk-go/v2/proto/services"
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
func (q *TransactionReceiptQuery) SetGrpcDeadline(deadline *time.Duration) *TransactionReceiptQuery {
	q.Query.SetGrpcDeadline(deadline)
	return q
}

// SetIncludeChildren Sets whether the response should include the receipts of any child transactions spawned by the
// top-level transaction with the given transactionID.
func (q *TransactionReceiptQuery) SetIncludeChildren(includeChildReceipts bool) *TransactionReceiptQuery {
	q.childReceipts = &includeChildReceipts
	return q
}

// GetIncludeChildren returns whether the response should include the receipts of any child transactions spawned by the
// top-level transaction with the given transactionID.
func (q *TransactionReceiptQuery) GetIncludeChildren() bool {
	if q.childReceipts != nil {
		return *q.childReceipts
	}

	return false
}

// SetIncludeDuplicates Sets whether receipts of processing duplicate transactions should be returned along with the
// receipt of processing the first consensus transaction with the given id whose status was
// neither INVALID_NODE_ACCOUNT nor INVALID_PAYER_SIGNATURE; or, if no
// such receipt exists, the receipt of processing the first transaction to reach consensus with
// the given transaction id.
func (q *TransactionReceiptQuery) SetIncludeDuplicates(includeDuplicates bool) *TransactionReceiptQuery {
	q.duplicates = &includeDuplicates
	return q
}

// GetIncludeDuplicates returns whether receipts of processing duplicate transactions should be returned along with the
// receipt of processing the first consensus transaction with the given id
func (q *TransactionReceiptQuery) GetIncludeDuplicates() bool {
	if q.duplicates != nil {
		return *q.duplicates
	}

	return false
}

func (q *TransactionReceiptQuery) GetCost(client *Client) (Hbar, error) {
	return q.Query.getCost(client, q)
}

// Execute executes the Query with the provided client
func (q *TransactionReceiptQuery) Execute(client *Client) (TransactionReceipt, error) {
	resp, err := q.Query.execute(client, q)

	if err, ok := err.(ErrHederaPreCheckStatus); ok {
		if resp.GetTransactionGetReceipt() != nil {
			return _TransactionReceiptFromProtobuf(resp.GetTransactionGetReceipt(), q.transactionID), err
		}
		// Manually add the receipt status, because an empty receipt has no status and no status defaults to 0, which means success
		return TransactionReceipt{Status: err.Status}, err
	}

	return _TransactionReceiptFromProtobuf(resp.GetTransactionGetReceipt(), q.transactionID), nil
}

// SetTransactionID sets the TransactionID for which the receipt is being requested.
func (q *TransactionReceiptQuery) SetTransactionID(transactionID TransactionID) *TransactionReceiptQuery {
	q.transactionID = &transactionID
	return q
}

// GetTransactionID returns the TransactionID for which the receipt is being requested.
func (q *TransactionReceiptQuery) GetTransactionID() TransactionID {
	if q.transactionID == nil {
		return TransactionID{}
	}

	return *q.transactionID
}

// SetNodeAccountIDs sets the _Node AccountID for this TransactionReceiptQuery.
func (q *TransactionReceiptQuery) SetNodeAccountIDs(accountID []AccountID) *TransactionReceiptQuery {
	q.Query.SetNodeAccountIDs(accountID)
	return q
}

// SetQueryPayment sets the Hbar payment to pay the _Node a fee for handling this query
func (q *TransactionReceiptQuery) SetQueryPayment(queryPayment Hbar) *TransactionReceiptQuery {
	q.Query.SetQueryPayment(queryPayment)
	return q
}

// SetMaxQueryPayment sets the maximum payment allowed for this Query.
func (q *TransactionReceiptQuery) SetMaxQueryPayment(queryMaxPayment Hbar) *TransactionReceiptQuery {
	q.Query.SetMaxQueryPayment(queryMaxPayment)
	return q
}

// SetMaxRetry sets the max number of errors before execution will fail.
func (q *TransactionReceiptQuery) SetMaxRetry(count int) *TransactionReceiptQuery {
	q.Query.SetMaxRetry(count)
	return q
}

// SetMaxBackoff The maximum amount of time to wait between retries.
// Every retry attempt will increase the wait time exponentially until it reaches this time.
func (q *TransactionReceiptQuery) SetMaxBackoff(max time.Duration) *TransactionReceiptQuery {
	q.Query.SetMaxBackoff(max)
	return q
}

// SetMinBackoff sets the minimum amount of time to wait between retries.
func (q *TransactionReceiptQuery) SetMinBackoff(min time.Duration) *TransactionReceiptQuery {
	q.Query.SetMinBackoff(min)
	return q
}

// SetPaymentTransactionID assigns the payment transaction id.
func (q *TransactionReceiptQuery) SetPaymentTransactionID(transactionID TransactionID) *TransactionReceiptQuery {
	q.Query.SetPaymentTransactionID(transactionID)
	return q
}

func (q *TransactionReceiptQuery) SetLogLevel(level LogLevel) *TransactionReceiptQuery {
	q.Query.SetLogLevel(level)
	return q
}

// ---------- Parent functions specific implementation ----------

func (q *TransactionReceiptQuery) getMethod(channel *_Channel) _Method {
	return _Method{
		query: channel._GetCrypto().GetTransactionReceipts,
	}
}

func (q *TransactionReceiptQuery) mapStatusError(_ Executable, response interface{}) error {
	status := Status(response.(*services.Response).GetTransactionGetReceipt().GetHeader().GetNodeTransactionPrecheckCode())
	switch status {
	case StatusPlatformTransactionNotCreated, StatusBusy, StatusUnknown, StatusOk:
		break
	default:
		return ErrHederaPreCheckStatus{
			Status: status,
		}
	}

	return ErrHederaPreCheckStatus{
		Status: Status(response.(*services.Response).GetTransactionGetReceipt().GetReceipt().GetStatus()),
	}
}

func (q *TransactionReceiptQuery) getName() string {
	return "TransactionReceiptQuery"
}

func (q *TransactionReceiptQuery) buildQuery() *services.Query {
	body := &services.TransactionGetReceiptQuery{
		Header: q.pbHeader,
	}

	if q.transactionID.AccountID != nil {
		body.TransactionID = q.transactionID._ToProtobuf()
	}

	if q.duplicates != nil {
		body.IncludeDuplicates = *q.duplicates
	}

	if q.childReceipts != nil {
		body.IncludeChildReceipts = q.GetIncludeChildren()
	}

	return &services.Query{
		Query: &services.Query_TransactionGetReceipt{
			TransactionGetReceipt: body,
		},
	}
}

func (q *TransactionReceiptQuery) validateNetworkOnIDs(client *Client) error {
	if client == nil || !client.autoValidateChecksums {
		return nil
	}

	if err := q.transactionID.AccountID.ValidateChecksum(client); err != nil {
		return err
	}

	return nil
}

func (q *TransactionReceiptQuery) shouldRetry(_ Executable, response interface{}) _ExecutionState {
	status := Status(response.(*services.Response).GetTransactionGetReceipt().GetHeader().GetNodeTransactionPrecheckCode())

	switch status {
	case StatusPlatformTransactionNotCreated, StatusBusy, StatusUnknown, StatusReceiptNotFound, StatusRecordNotFound, StatusPlatformNotActive:
		return executionStateRetry
	case StatusOk:
		break
	default:
		return executionStateError
	}

	status = Status(response.(*services.Response).GetTransactionGetReceipt().GetReceipt().GetStatus())

	switch status {
	case StatusBusy, StatusUnknown, StatusOk, StatusReceiptNotFound, StatusRecordNotFound, StatusPlatformNotActive:
		return executionStateRetry
	default:
		return executionStateFinished
	}
}

func (q *TransactionReceiptQuery) getQueryResponse(response *services.Response) queryResponse {
	return response.GetTransactionGetReceipt()
}
