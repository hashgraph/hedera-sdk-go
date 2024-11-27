package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"fmt"
	"time"

	"github.com/hiero-ledger/hiero-sdk-go/v2/proto/services"
	"github.com/pkg/errors"
	protobuf "google.golang.org/protobuf/proto"
)

// Query is the struct used to build queries.
type Query struct {
	executable
	client                *Client
	pb                    *services.Query
	pbHeader              *services.QueryHeader //nolint
	paymentTransactionIDs *_LockableSlice

	paymentTransactions []*services.Transaction
	maxQueryPayment     Hbar
	queryPayment        Hbar
	timestamp           time.Time

	isPaymentRequired bool
}

type queryResponse interface {
	GetHeader() *services.ResponseHeader
}

type QueryInterface interface {
	Executable

	buildQuery() *services.Query
	getQueryResponse(response *services.Response) queryResponse
}

// -------- Executable functions ----------

func _NewQuery(isPaymentRequired bool, header *services.QueryHeader) Query {
	minBackoff := 250 * time.Millisecond
	maxBackoff := 8 * time.Second
	return Query{
		pb:                    &services.Query{},
		pbHeader:              header,
		paymentTransactionIDs: _NewLockableSlice(),
		paymentTransactions:   make([]*services.Transaction, 0),
		isPaymentRequired:     isPaymentRequired,
		maxQueryPayment:       NewHbar(0),
		queryPayment:          NewHbar(0),
		executable: executable{
			nodeAccountIDs: _NewLockableSlice(),
			maxBackoff:     &maxBackoff,
			minBackoff:     &minBackoff,
			maxRetry:       10,
		},
	}
}

// SetMaxQueryPayment sets the maximum payment allowed for this query.
func (q *Query) SetMaxQueryPayment(maxPayment Hbar) *Query {
	q.maxQueryPayment = maxPayment
	return q
}

// SetQueryPayment sets the payment amount for this query.
func (q *Query) SetQueryPayment(paymentAmount Hbar) *Query {
	q.queryPayment = paymentAmount
	return q
}

// GetMaxQueryPayment returns the maximum payment allowed for this query.
func (q *Query) GetMaxQueryPayment() Hbar {
	return q.maxQueryPayment
}

// GetQueryPayment returns the payment amount for this query.
func (q *Query) GetQueryPayment() Hbar {
	return q.queryPayment
}

// GetCost returns the fee that would be charged to get the requested information (if a cost was requested).
func (q *Query) getCost(client *Client, e QueryInterface) (Hbar, error) {
	if client == nil || client.operator == nil {
		return Hbar{}, errNoClientProvided
	}

	var err error

	err = e.validateNetworkOnIDs(client)
	if err != nil {
		return Hbar{}, err
	}
	q.paymentTransactions = make([]*services.Transaction, 0)
	if !q.nodeAccountIDs.locked {
		q.SetNodeAccountIDs([]AccountID{client.network._GetNode().accountID})
	}

	q.pb = e.buildQuery()

	if q.isPaymentRequired && len(q.paymentTransactions) > 0 {
		q.paymentTransactionIDs._Advance()
	}

	q.pbHeader.ResponseType = services.ResponseType_COST_ANSWER
	q.paymentTransactionIDs._Advance()
	resp, err := _Execute(client, e)

	if err != nil {
		return Hbar{}, err
	}

	queryResp := e.getQueryResponse(resp.(*services.Response))
	cost := int64(queryResp.GetHeader().Cost)

	return HbarFromTinybar(cost), nil
}

func _QueryMakePaymentTransaction(transactionID TransactionID, nodeAccountID AccountID, operator *_Operator, cost Hbar) (*services.Transaction, error) {
	accountAmounts := make([]*services.AccountAmount, 0)
	accountAmounts = append(accountAmounts, &services.AccountAmount{
		AccountID: nodeAccountID._ToProtobuf(),
		Amount:    cost.tinybar,
	})
	accountAmounts = append(accountAmounts, &services.AccountAmount{
		AccountID: operator.accountID._ToProtobuf(),
		Amount:    -cost.tinybar,
	})

	body := services.TransactionBody{
		TransactionID:  transactionID._ToProtobuf(),
		NodeAccountID:  nodeAccountID._ToProtobuf(),
		TransactionFee: uint64(NewHbar(1).tinybar),
		TransactionValidDuration: &services.Duration{
			Seconds: 120,
		},
		Data: &services.TransactionBody_CryptoTransfer{
			CryptoTransfer: &services.CryptoTransferTransactionBody{
				Transfers: &services.TransferList{
					AccountAmounts: accountAmounts,
				},
			},
		},
	}

	bodyBytes, err := protobuf.Marshal(&body)
	if err != nil {
		return nil, errors.Wrap(err, "error serializing Query body")
	}

	signature := operator.signer(bodyBytes)
	if len(signature) == 65 {
		signature = signature[1:]
	}
	sigPairs := make([]*services.SignaturePair, 0)
	sigPairs = append(sigPairs, operator.publicKey._ToSignaturePairProtobuf(signature))

	return &services.Transaction{
		BodyBytes: bodyBytes,
		SigMap: &services.SignatureMap{
			SigPair: sigPairs,
		},
	}, nil
}

// GetPaymentTransactionID returns the payment transaction id.
func (q *Query) GetPaymentTransactionID() TransactionID {
	if !q.paymentTransactionIDs._IsEmpty() {
		return q.paymentTransactionIDs._GetCurrent().(TransactionID)
	}

	return TransactionID{}
}

// GetMaxRetryCount returns the max number of errors before execution will fail.
func (q *Query) GetMaxRetryCount() int {
	return q.GetMaxRetry()
}

// SetPaymentTransactionID assigns the payment transaction id.
func (q *Query) SetPaymentTransactionID(transactionID TransactionID) *Query {
	q.paymentTransactionIDs._Clear()._Push(transactionID)._SetLocked(true)
	return q
}

func (q *Query) execute(client *Client, e QueryInterface) (*services.Response, error) {
	q.client = client
	if client == nil {
		return nil, errNoClientProvided
	}

	var err error

	err = e.validateNetworkOnIDs(client)
	if err != nil {
		return nil, err
	}

	var cost Hbar
	if q.queryPayment.tinybar == 0 && q.isPaymentRequired {
		if q.maxQueryPayment.tinybar == 0 {
			cost = client.GetDefaultMaxQueryPayment()
		} else {
			cost = q.maxQueryPayment
		}

		actualCost, err := q.getCost(client, e)
		if err != nil {
			return nil, err
		}

		if cost.tinybar < actualCost.tinybar {
			return nil, ErrMaxQueryPaymentExceeded{
				QueryCost:       actualCost,
				MaxQueryPayment: cost,
				query:           e.getName(),
			}
		}

		q.queryPayment = actualCost
	}

	q.paymentTransactions = make([]*services.Transaction, 0)
	if !q.nodeAccountIDs.locked {
		q.SetNodeAccountIDs([]AccountID{client.network._GetNode().accountID})
	}

	q.pb = e.buildQuery()
	q.pbHeader.ResponseType = services.ResponseType_ANSWER_ONLY

	resp, err := _Execute(client, e)
	if err != nil {
		return nil, err
	}

	return resp.(*services.Response), nil
}

func (q *Query) shouldRetry(e Executable, response interface{}) _ExecutionState {
	queryResp := e.(QueryInterface).getQueryResponse(response.(*services.Response))

	status := Status(queryResp.GetHeader().NodeTransactionPrecheckCode)

	retryableStatuses := map[Status]bool{
		StatusPlatformTransactionNotCreated: true,
		StatusPlatformNotActive:             true,
		StatusBusy:                          true,
	}

	if retryableStatuses[status] {
		return executionStateRetry
	}

	if status == StatusOk {
		return executionStateFinished
	}

	return executionStateError
}

func (q *Query) generatePayments(client *Client, cost Hbar) (*services.Transaction, error) {
	var tx *services.Transaction
	var err error
	for _, nodeID := range q.nodeAccountIDs.slice {
		txnID := TransactionIDGenerate(client.operator.accountID)
		tx, err = _QueryMakePaymentTransaction(
			txnID,
			nodeID.(AccountID),
			client.operator,
			cost,
		)
		if err != nil {
			return nil, err
		}
		q.paymentTransactions = append(q.paymentTransactions, tx)
	}
	return tx, nil
}

func (q *Query) advanceRequest() {
	q.nodeAccountIDs._Advance()
}

func (q *Query) makeRequest() interface{} {
	if q.client != nil && q.isPaymentRequired {
		tx, err := q.generatePayments(q.client, q.queryPayment)
		if err != nil {
			return q.pb
		}
		q.pbHeader.Payment = tx
	}

	return q.pb
}

func (q *Query) mapResponse(response interface{}, _ AccountID, _ interface{}) (interface{}, error) { // nolint
	return response.(*services.Response), nil
}

func (q *Query) isTransaction() bool {
	return false
}

func (q *Query) mapStatusError(e Executable, response interface{}) error {
	queryResp := e.(QueryInterface).getQueryResponse(response.(*services.Response))
	return ErrHederaPreCheckStatus{
		Status: Status(queryResp.GetHeader().NodeTransactionPrecheckCode),
	}
}

// ----------- Next methods should be overridden in each subclass ---------------

// NOTE: Should be implemented in every inheritor. Example:
//
//	return ErrHederaPreCheckStatus{
//		Status: Status(response.(*services.Response).GetNetworkGetVersionInfo().Header.NodeTransactionPrecheckCode),
//	}
func (q *Query) getMethod(*_Channel) _Method {
	return _Method{}
}

func (q *Query) getName() string {
	return "QueryInterface"
}

func (q *Query) getLogID(queryInterface Executable) string {
	timestamp := q.timestamp.UnixNano()
	return fmt.Sprintf("%s:%d", queryInterface.getName(), timestamp)
}

//lint:ignore U1000
func (q *Query) buildQuery() *services.Query {
	return nil
}

//lint:ignore U1000
func (q *Query) buildScheduled() (*services.SchedulableTransactionBody, error) {
	return nil, errors.New("Not implemented")
}

// NOTE: Should be implemented in every inheritor.
func (q *Query) validateNetworkOnIDs(*Client) error {
	return errors.New("Not implemented")
}

func (q *Query) getTransactionIDAndMessage() (string, string) {
	txID := q.GetPaymentTransactionID().String()
	if txID == "" {
		txID = "None"
	}
	return txID, "QueryInterface status received"
}
