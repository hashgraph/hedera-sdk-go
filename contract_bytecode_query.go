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

// ContractBytecodeQuery retrieves the bytecode for a smart contract instance
type ContractBytecodeQuery struct {
	query
	contractID *ContractID
}

// NewContractBytecodeQuery creates a ContractBytecodeQuery query which can be used to construct and execute a
// Contract Get Bytecode Query.
func NewContractBytecodeQuery() *ContractBytecodeQuery {
	header := services.QueryHeader{}
	result := ContractBytecodeQuery{
		query: _NewQuery(true, &header),
	}
	//	result.e = &result
	return &result
}

// When execution is attempted, a single attempt will timeout when this deadline is reached. (The SDK may subsequently retry the execution.)
func (q *ContractBytecodeQuery) SetGrpcDeadline(deadline *time.Duration) *ContractBytecodeQuery {
	q.query.SetGrpcDeadline(deadline)
	return q
}

// SetContractID sets the contract for which the bytecode is requested
func (q *ContractBytecodeQuery) SetContractID(contractID ContractID) *ContractBytecodeQuery {
	q.contractID = &contractID
	return q
}

// GetContractID returns the contract for which the bytecode is requested
func (q *ContractBytecodeQuery) GetContractID() ContractID {
	if q.contractID == nil {
		return ContractID{}
	}

	return *q.contractID
}

func (q *ContractBytecodeQuery) GetCost(client *Client) (Hbar, error) {
	if client == nil || client.operator == nil {
		return Hbar{}, errNoClientProvided
	}

	var err error

	err = q.validateNetworkOnIDs(client)
	if err != nil {
		return Hbar{}, err
	}

	for range q.nodeAccountIDs.slice {
		paymentTransaction, err := _QueryMakePaymentTransaction(TransactionID{}, AccountID{}, client.operator, Hbar{})
		if err != nil {
			return Hbar{}, err
		}
		q.paymentTransactions = append(q.paymentTransactions, paymentTransaction)
	}

	pb := q.build()
	pb.ContractGetBytecode.Header = q.pbHeader

	q.pb = &services.Query{
		Query: pb,
	}

	q.pbHeader.ResponseType = services.ResponseType_COST_ANSWER
	q.paymentTransactionIDs._Advance()

	resp, err := _Execute(
		client,
		q.e,
	)

	if err != nil {
		return Hbar{}, err
	}

	cost := int64(resp.(*services.Response).GetContractGetBytecodeResponse().Header.Cost)
	return HbarFromTinybar(cost), nil
}

// Execute executes the Query with the provided client
func (q *ContractBytecodeQuery) Execute(client *Client) ([]byte, error) {
	if client == nil || client.operator == nil {
		return make([]byte, 0), errNoClientProvided
	}

	var err error

	err = q.validateNetworkOnIDs(client)
	if err != nil {
		return []byte{}, err
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
			return []byte{}, err
		}

		if cost.tinybar < actualCost.tinybar {
			return []byte{}, ErrMaxQueryPaymentExceeded{
				QueryCost:       actualCost,
				MaxQueryPayment: cost,
				query:           "ContractBytecodeQuery",
			}
		}

		cost = actualCost
	}

	q.paymentTransactions = make([]*services.Transaction, 0)

	if q.nodeAccountIDs.locked {
		err = q._QueryGeneratePayments(client, cost)
		if err != nil {
			if err != nil {
				return []byte{}, err
			}
		}
	} else {
		paymentTransaction, err := _QueryMakePaymentTransaction(q.paymentTransactionIDs._GetCurrent().(TransactionID), AccountID{}, client.operator, cost)
		if err != nil {
			return []byte{}, err
		}
		q.paymentTransactions = append(q.paymentTransactions, paymentTransaction)
	}

	pb := q.build()
	pb.ContractGetBytecode.Header = q.pbHeader
	q.pb = &services.Query{
		Query: pb,
	}

	if q.isPaymentRequired && len(q.paymentTransactions) > 0 {
		q.paymentTransactionIDs._Advance()
	}
	q.pbHeader.ResponseType = services.ResponseType_ANSWER_ONLY

	resp, err := _Execute(
		client,
		q.e,
	)

	if err != nil {
		return []byte{}, err
	}

	return resp.(*services.Response).GetContractGetBytecodeResponse().Bytecode, nil
}

// SetMaxQueryPayment sets the maximum payment allowed for this Query.
func (q *ContractBytecodeQuery) SetMaxQueryPayment(maxPayment Hbar) *ContractBytecodeQuery {
	q.query.SetMaxQueryPayment(maxPayment)
	return q
}

// SetQueryPayment sets the payment amount for this Query.
func (q *ContractBytecodeQuery) SetQueryPayment(paymentAmount Hbar) *ContractBytecodeQuery {
	q.query.SetQueryPayment(paymentAmount)
	return q
}

// SetNodeAccountIDs sets the _Node AccountID for this ContractBytecodeQuery.
func (q *ContractBytecodeQuery) SetNodeAccountIDs(accountID []AccountID) *ContractBytecodeQuery {
	q.query.SetNodeAccountIDs(accountID)
	return q
}

// SetMaxRetry sets the max number of errors before execution will fail.
func (q *ContractBytecodeQuery) SetMaxRetry(count int) *ContractBytecodeQuery {
	q.query.SetMaxRetry(count)
	return q
}

// SetMaxBackoff The maximum amount of time to wait between retries.
// Every retry attempt will increase the wait time exponentially until it reaches this time.
func (q *ContractBytecodeQuery) SetMaxBackoff(max time.Duration) *ContractBytecodeQuery {
	q.query.SetMaxBackoff(max)
	return q
}

// SetMinBackoff sets the minimum amount of time to wait between retries.
func (q *ContractBytecodeQuery) SetMinBackoff(min time.Duration) *ContractBytecodeQuery {
	q.SetMinBackoff(min)
	return q
}

// SetPaymentTransactionID assigns the payment transaction id.
func (q *ContractBytecodeQuery) SetPaymentTransactionID(transactionID TransactionID) *ContractBytecodeQuery {
	q.paymentTransactionIDs._Clear()._Push(transactionID)._SetLocked(true)
	return q
}

func (q *ContractBytecodeQuery) SetLogLevel(level LogLevel) *ContractBytecodeQuery {
	q.query.SetLogLevel(level)
	return q
}

// ---------- Parent functions specific implementation ----------

func (q *ContractBytecodeQuery) getMethod(channel *_Channel) _Method {
	return _Method{
		query: channel._GetContract().ContractGetBytecode}
}

func (q *ContractBytecodeQuery) mapStatusError(_ interface{}, response interface{}) error {
	return ErrHederaPreCheckStatus{
		Status: Status(response.(*services.Response).GetContractGetBytecodeResponse().Header.NodeTransactionPrecheckCode),
	}
}

func (q *ContractBytecodeQuery) getName() string {
	return "ContractBytecodeQuery"
}

func (q *ContractBytecodeQuery) build() *services.Query_ContractGetBytecode {
	pb := services.Query_ContractGetBytecode{
		ContractGetBytecode: &services.ContractGetBytecodeQuery{
			Header: &services.QueryHeader{},
		},
	}

	if q.contractID != nil {
		pb.ContractGetBytecode.ContractID = q.contractID._ToProtobuf()
	}

	return &pb
}

func (q *ContractBytecodeQuery) validateNetworkOnIDs(client *Client) error {
	if client == nil || !client.autoValidateChecksums {
		return nil
	}

	if q.contractID != nil {
		if err := q.contractID.ValidateChecksum(client); err != nil {
			return err
		}
	}

	return nil
}

func (q *ContractBytecodeQuery) getQueryStatus(response interface{}) Status {
	return Status(response.(*services.Response).GetContractGetBytecodeResponse().Header.NodeTransactionPrecheckCode)
}
