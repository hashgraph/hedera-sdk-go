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

// Query is the struct used to build queries.
type Query struct {
	pb       *services.Query
	pbHeader *services.QueryHeader //nolint

	paymentTransactionIDs *_LockableSlice
	nodeAccountIDs        *_LockableSlice
	maxQueryPayment       Hbar
	queryPayment          Hbar
	maxRetry              int

	paymentTransactions []*services.Transaction

	isPaymentRequired bool

	maxBackoff   *time.Duration
	minBackoff   *time.Duration
	grpcDeadline *time.Duration
	timestamp    time.Time
	logLevel     *LogLevel
}

func _NewQuery(isPaymentRequired bool, header *services.QueryHeader) Query {
	minBackoff := 250 * time.Millisecond
	maxBackoff := 8 * time.Second
	return Query{
		pb:                    &services.Query{},
		pbHeader:              header,
		paymentTransactionIDs: _NewLockableSlice(),
		maxRetry:              10,
		nodeAccountIDs:        _NewLockableSlice(),
		paymentTransactions:   make([]*services.Transaction, 0),
		isPaymentRequired:     isPaymentRequired,
		maxQueryPayment:       NewHbar(0),
		queryPayment:          NewHbar(0),
		timestamp:             time.Now(),
		maxBackoff:            &maxBackoff,
		minBackoff:            &minBackoff,
	}
}

// When execution is attempted, a single attempt will timeout when this deadline is reached. (The SDK may subsequently retry the execution.)
func (this *Query) SetGrpcDeadline(deadline *time.Duration) *Query {
	this.grpcDeadline = deadline
	return this
}

// GetGrpcDeadline returns the grpc deadline.
func (this *Query) GetGrpcDeadline() *time.Duration {
	return this.grpcDeadline
}

// SetNodeAccountID sets the node account ID for this Query.
func (this *Query) SetNodeAccountIDs(nodeAccountIDs []AccountID) *Query {
	for _, nodeAccountID := range nodeAccountIDs {
		this.nodeAccountIDs._Push(nodeAccountID)
	}
	this.nodeAccountIDs._SetLocked(true)
	return this
}

// GetNodeAccountID returns the node account ID for this Query.
func (this *Query) GetNodeAccountIDs() (nodeAccountIDs []AccountID) {
	nodeAccountIDs = []AccountID{}

	for _, value := range this.nodeAccountIDs.slice {
		nodeAccountIDs = append(nodeAccountIDs, value.(AccountID))
	}

	return nodeAccountIDs
}

func _QueryGetNodeAccountID(request interface{}) AccountID {
	return request.(*Query).nodeAccountIDs._GetCurrent().(AccountID)
}

// SetMaxQueryPayment sets the maximum payment allowed for this Query.
func (this *Query) SetMaxQueryPayment(maxPayment Hbar) *Query {
	this.maxQueryPayment = maxPayment
	return this
}

// SetQueryPayment sets the payment amount for this Query.
func (this *Query) SetQueryPayment(paymentAmount Hbar) *Query {
	this.queryPayment = paymentAmount
	return this
}

// GetMaxQueryPayment returns the maximum payment allowed for this Query.
func (this *Query) GetMaxQueryPayment() Hbar {
	return this.maxQueryPayment
}

// GetQueryPayment returns the payment amount for this Query.
func (this *Query) GetQueryPayment() Hbar {
	return this.queryPayment
}

// GetMaxRetryCount returns the max number of errors before execution will fail.
func (this *Query) GetMaxRetryCount() int {
	return this.maxRetry
}

// SetMaxRetry sets the max number of errors before execution will fail.
func (this *Query) SetMaxRetry(count int) *Query {
	this.maxRetry = count
	return this
}

func _QueryShouldRetry(status Status) _ExecutionState {
	switch status {
	case StatusPlatformTransactionNotCreated, StatusPlatformNotActive, StatusBusy:
		return executionStateRetry
	case StatusOk:
		return executionStateFinished
	}

	return executionStateError
}

func _QueryMakeRequest(request interface{}) interface{} {
	query := request.(*Query)
	if query.isPaymentRequired && len(query.paymentTransactions) > 0 {
		query.pbHeader.Payment = query.paymentTransactions[query.paymentTransactionIDs.index]
	}
	query.pbHeader.ResponseType = services.ResponseType_ANSWER_ONLY

	return query.pb
}

func _CostQueryMakeRequest(request interface{}) interface{} {
	query := request.(*Query)
	if query.isPaymentRequired && len(query.paymentTransactions) > 0 {
		query.pbHeader.Payment = query.paymentTransactions[query.paymentTransactionIDs.index]
	}
	query.pbHeader.ResponseType = services.ResponseType_COST_ANSWER
	return query.pb
}

func _QueryAdvanceRequest(request interface{}) {
	query := request.(*Query)
	if query.isPaymentRequired && len(query.paymentTransactions) > 0 {
		query.paymentTransactionIDs._Advance()
	}
	query.nodeAccountIDs._Advance()
}

func _CostQueryAdvanceRequest(request interface{}) {
	query := request.(*Query)
	query.paymentTransactionIDs._Advance()
	query.nodeAccountIDs._Advance()
}

func _QueryMapResponse(request interface{}, response interface{}, _ AccountID, protoRequest interface{}) (interface{}, error) {
	return response.(*services.Response), nil
}

func _QueryGeneratePayments(query *Query, client *Client, cost Hbar) error {
	for _, nodeID := range query.nodeAccountIDs.slice {
		transaction, err := _QueryMakePaymentTransaction(
			query.paymentTransactionIDs._GetCurrent().(TransactionID),
			nodeID.(AccountID),
			client.operator,
			cost,
		)
		if err != nil {
			return err
		}

		query.paymentTransactions = append(query.paymentTransactions, transaction)
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
func (this *Query) GetPaymentTransactionID() TransactionID {
	if !this.paymentTransactionIDs._IsEmpty() {
		return this.paymentTransactionIDs._GetCurrent().(TransactionID)
	}

	return TransactionID{}
}

// SetPaymentTransactionID assigns the payment transaction id.
func (this *Query) SetPaymentTransactionID(transactionID TransactionID) *Query {
	this.paymentTransactionIDs._Clear()._Push(transactionID)._SetLocked(true)
	return this
}

func (query *Query) SetLogLevel(level LogLevel) *Query {
	query.logLevel = &level
	return query
}

func (query *Query) GetLogLevel() *LogLevel {
	return query.logLevel
}
