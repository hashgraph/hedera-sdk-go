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

type Query interface {
	Executable

	Execute(client *Client) (TransactionResponse, error)

	build() *services.TransactionBody
	getQueryStatus(response interface{}) Status
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

// When execution is attempted, a single attempt will timeout when this deadline is reached. (The SDK may subsequently retry the execution.)
func (q *query) SetGrpcDeadline(deadline *time.Duration) *query {
	q.grpcDeadline = deadline
	return q
}

// SetNodeAccountID sets the node account ID for this Query.
func (q *query) SetNodeAccountIDs(nodeAccountIDs []AccountID) *query {

	for _, nodeAccountID := range nodeAccountIDs {
		q.nodeAccountIDs._Push(nodeAccountID)
	}
	q.nodeAccountIDs._SetLocked(true)
	return q
}

// SetMaxQueryPayment sets the maximum payment allowed for this Query.
func (q *query) SetMaxQueryPayment(maxPayment Hbar) *query {
	q.maxQueryPayment = maxPayment
	return q
}

// SetQueryPayment sets the payment amount for this Query.
func (q *query) SetQueryPayment(paymentAmount Hbar) *query {
	q.queryPayment = paymentAmount
	return q
}

// GetMaxQueryPayment returns the maximum payment allowed for this Query.
func (q *query) GetMaxQueryPayment() Hbar {
	return q.maxQueryPayment
}

// GetQueryPayment returns the payment amount for this Query.
func (q *query) GetQueryPayment() Hbar {
	return q.queryPayment
}

// SetMaxRetry sets the max number of errors before execution will fail.
func (q *query) SetMaxRetry(count int) *query {
	q.maxRetry = count
	return q
}

func (q *query) shouldRetry(_ interface{}, response interface{}) _ExecutionState {
	status := q.getQueryStatus(response)
	switch status {
	case StatusPlatformTransactionNotCreated, StatusPlatformNotActive, StatusBusy:
		return executionStateRetry
	case StatusOk:
		return executionStateFinished
	}

	return executionStateError
}

func (q *query) _QueryGeneratePayments(client *Client, cost Hbar) error {
	for _, nodeID := range q.nodeAccountIDs.slice {
		transaction, err := _QueryMakePaymentTransaction(
			q.paymentTransactionIDs._GetCurrent().(TransactionID),
			nodeID.(AccountID),
			client.operator,
			cost,
		)
		if err != nil {
			return err
		}

		q.paymentTransactions = append(q.paymentTransactions, transaction)
	}

	return nil
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

func (q *query) SetLogLevel(level LogLevel) *query {
	q.logLevel = &level
	return q
}

func (q *query) advanceRequest(request interface{}) {
	query := request.(*query)
	query.nodeAccountIDs._Advance()
}
func (q *query) getNodeAccountID(request interface{}) AccountID {
	return request.(*query).nodeAccountIDs._GetCurrent().(AccountID)
}

func (q *query) makeRequest(request interface{}) interface{} {
	query := request.(*query)
	if query.isPaymentRequired && len(query.paymentTransactions) > 0 {
		query.pbHeader.Payment = query.paymentTransactions[query.paymentTransactionIDs.index]
	}

	return query.pb
}

func (q *query) mapResponse(request interface{}, response interface{}, _ AccountID, protoRequest interface{}) (interface{}, error) {
	return response.(*services.Response), nil
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

// NOTE: Should be implemented in every inheritor. Example:
//
//	return ErrHederaPreCheckStatus{
//		Status: Status(response.(*services.Response).GetCryptoGetInfo().Header.NodeTransactionPrecheckCode),
//	}
func (q *query) mapStatusError(interface{}, interface{}) error {
	return errors.New("Not implemented")
}
func (q *query) getName() string {
	return "Query"
}

// NOTE: Should be implemented in every inheritor.
func (q *query) build() *services.TransactionBody {
	return nil
}

// NOTE: Should be implemented in every inheritor.
func (q *query) validateNetworkOnIDs(client *Client) error {
	return errors.New("Not implemented")
}
func (q *query) getQueryStatus(response interface{}) Status {
	return Status(1)
}
