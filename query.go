package hedera

/*-
 *
 * Hedera Go SDK
 *
 * Copyright (C) 2020 - 2023 Hedera Hashgraph, LLC
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use q file except in compliance with the License.
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
	"github.com/pkg/errors"
	protobuf "google.golang.org/protobuf/proto"
)

// query is the struct used to build queries.
type query struct {
	executable
	pb       *services.Query
	pbHeader *services.QueryHeader //nolint

	paymentTransactionIDs *_LockableSlice
	maxQueryPayment       Hbar
	queryPayment          Hbar

	paymentTransactions []*services.Transaction

	isPaymentRequired bool
	timestamp         time.Time
}

type queryResponse interface {
	GetHeader() *services.ResponseHeader
}

type Query interface {
	Executable

	execute(client *Client) (*services.Response, error)
	buildQuery() *services.Query
	getQueryResponse(response *services.Response) queryResponse
}

// -------- Executable functions ----------

func _NewQuery(isPaymentRequired bool, header *services.QueryHeader) query {
	minBackoff := 250 * time.Millisecond
	maxBackoff := 8 * time.Second
	return query{
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

// SetMaxQueryPayment sets the maximum payment allowed for q Query.
func (q *query) SetMaxQueryPayment(maxPayment Hbar) *query {
	q.maxQueryPayment = maxPayment
	return q
}

// SetQueryPayment sets the payment amount for q Query.
func (q *query) SetQueryPayment(paymentAmount Hbar) *query {
	q.queryPayment = paymentAmount
	return q
}

// GetMaxQueryPayment returns the maximum payment allowed for q Query.
func (q *query) GetMaxQueryPayment() Hbar {
	return q.maxQueryPayment
}

// GetQueryPayment returns the payment amount for q Query.
func (q *query) GetQueryPayment() Hbar {
	return q.queryPayment
}

// GetCost returns the fee that would be charged to get the requested information (if a cost was requested).
func (q *query) GetCost(client *Client) (Hbar, error) {
	if client == nil || client.operator == nil {
		return Hbar{}, errNoClientProvided
	}

	var err error

	if !q.paymentTransactionIDs.locked {
		q.paymentTransactionIDs._Clear()._Push(TransactionIDGenerate(client.operator.accountID))
	}

	err = q.e.validateNetworkOnIDs(client)
	if err != nil {
		return Hbar{}, err
	}

	if !q.nodeAccountIDs.locked {
		q.SetNodeAccountIDs([]AccountID{client.network._GetNode().accountID})
	}

	err = q.generatePayments(client, Hbar{})
	if err != nil {
		return Hbar{}, err
	}

	q.pb = q.e.(Query).buildQuery()

	q.pbHeader.ResponseType = services.ResponseType_COST_ANSWER
	q.paymentTransactionIDs._Advance()
	resp, err := _Execute(
		client,
		q.e,
	)

	if err != nil {
		return Hbar{}, err
	}

	queryResp := q.e.(Query).getQueryResponse(resp.(*services.Response))
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
		return nil, errors.Wrap(err, "error serializing query body")
	}

	signature := operator.signer(bodyBytes)
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
func (q *query) GetPaymentTransactionID() TransactionID {
	if !q.paymentTransactionIDs._IsEmpty() {
		return q.paymentTransactionIDs._GetCurrent().(TransactionID)
	}

	return TransactionID{}
}

// GetMaxRetryCount returns the max number of errors before execution will fail.
func (q *query) GetMaxRetryCount() int {
	return q.GetMaxRetry()
}

// SetPaymentTransactionID assigns the payment transaction id.
func (q *query) SetPaymentTransactionID(transactionID TransactionID) *query {
	q.paymentTransactionIDs._Clear()._Push(transactionID)._SetLocked(true)
	return q
}

func (q *query) execute(client *Client) (*services.Response, error) {
	if client == nil || client.operator == nil {
		return nil, errNoClientProvided
	}

	var err error

	err = q.e.validateNetworkOnIDs(client)
	if err != nil {
		return nil, err
	}

	if !q.paymentTransactionIDs.locked {
		q.paymentTransactionIDs._Clear()._Push(TransactionIDGenerate(client.operator.accountID))
	}

	var cost Hbar
	if q.queryPayment.tinybar != 0 {
		cost = q.queryPayment
	} else {
		if q.maxQueryPayment.tinybar == 0 {
			cost = client.GetDefaultMaxQueryPayment()
		} else {
			cost = q.maxQueryPayment
		}

		actualCost, err := q.GetCost(client)
		if err != nil {
			return nil, err
		}

		if cost.tinybar < actualCost.tinybar {
			return nil, ErrMaxQueryPaymentExceeded{
				QueryCost:       actualCost,
				MaxQueryPayment: cost,
				query:           q.e.getName(),
			}
		}

		cost = actualCost
	}

	q.paymentTransactions = make([]*services.Transaction, 0)

	if !q.nodeAccountIDs.locked {
		q.SetNodeAccountIDs([]AccountID{client.network._GetNode().accountID})
	}

	if cost.tinybar > 0 {
		err = q.generatePayments(client, cost)

		if err != nil {
			return nil, err
		}
	}

	q.pb = q.e.(Query).buildQuery()
	q.pbHeader.ResponseType = services.ResponseType_ANSWER_ONLY

	if q.isPaymentRequired && len(q.paymentTransactions) > 0 {
		q.paymentTransactionIDs._Advance()
	}

	resp, err := _Execute(client, q.e)
	if err != nil {
		return nil, err
	}

	return resp.(*services.Response), nil
}

func (q *query) shouldRetry(response interface{}) _ExecutionState {
	queryResp := q.e.(Query).getQueryResponse(response.(*services.Response))
	status := Status(queryResp.GetHeader().NodeTransactionPrecheckCode)
	switch status {
	case StatusPlatformTransactionNotCreated, StatusPlatformNotActive, StatusBusy:
		return executionStateRetry
	case StatusOk:
		return executionStateFinished
	}

	return executionStateError
}

func (q *query) generatePayments(client *Client, cost Hbar) error {
	for _, nodeID := range q.nodeAccountIDs.slice {
		tx, err := _QueryMakePaymentTransaction(
			q.paymentTransactionIDs._GetCurrent().(TransactionID),
			nodeID.(AccountID),
			client.operator,
			cost,
		)
		if err != nil {
			return err
		}

		q.paymentTransactions = append(q.paymentTransactions, tx)
	}

	return nil
}

func (q *query) advanceRequest() {
	q.nodeAccountIDs._Advance()
}

func (q *query) makeRequest() interface{} {
	if q.isPaymentRequired && len(q.paymentTransactions) > 0 {
		q.pbHeader.Payment = q.paymentTransactions[q.paymentTransactionIDs.index]
	}

	return q.pb
}

func (q *query) mapResponse(response interface{}, _ AccountID, protoRequest interface{}) (interface{}, error) {
	return response.(*services.Response), nil
}

func (q *query) isTransaction() bool {
	return false
}

func (q *query) mapStatusError(response interface{}) error {
	queryResp := q.e.(Query).getQueryResponse(response.(*services.Response))
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
func (q *query) getMethod(*_Channel) _Method {
	return _Method{}
}

func (q *query) getName() string {
	return "Query"
}

// NOTE: Should be implemented in every inheritor.
func (q *query) buildQuery() *services.Query {
	return nil
}

func (q *query) buildScheduled() (*services.SchedulableTransactionBody, error) {
	return nil, errors.New("Not implemented")
}

// NOTE: Should be implemented in every inheritor.
func (q *query) validateNetworkOnIDs(client *Client) error {
	return errors.New("Not implemented")
}
