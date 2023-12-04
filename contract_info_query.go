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

// ContractInfoQuery retrieves information about a smart contract instance. This includes the account that it uses, the
// file containing its bytecode, and the time when it will expire.
type ContractInfoQuery struct {
	query
	contractID *ContractID
}

// NewContractInfoQuery creates a ContractInfoQuery query which can be used to construct and execute a
// Contract Get Info Query.
func NewContractInfoQuery() *ContractInfoQuery {
	header := services.QueryHeader{}
	query := _NewQuery(true, &header)

	result := ContractInfoQuery{
		query: query,
	}

	//	result.e = &result
	return &result
}

// When execution is attempted, a single attempt will timeout when this deadline is reached. (The SDK may subsequently retry the execution.)
func (q *ContractInfoQuery) SetGrpcDeadline(deadline *time.Duration) *ContractInfoQuery {
	q.query.SetGrpcDeadline(deadline)
	return q
}

// SetContractID sets the contract for which information is requested
func (q *ContractInfoQuery) SetContractID(contractID ContractID) *ContractInfoQuery {
	q.contractID = &contractID
	return q
}

func (q *ContractInfoQuery) GetContractID() ContractID {
	if q.contractID == nil {
		return ContractID{}
	}

	return *q.contractID
}

// GetCost returns the fee that would be charged to get the requested information (if a cost was requested).
func (q *ContractInfoQuery) GetCost(client *Client) (Hbar, error) {
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
	pb.ContractGetInfo.Header = q.pbHeader

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

	cost := int64(resp.(*services.Response).GetContractGetInfo().Header.Cost)

	return HbarFromTinybar(cost), nil
}

// Execute executes the Query with the provided client
func (q *ContractInfoQuery) Execute(client *Client) (ContractInfo, error) {
	if client == nil || client.operator == nil {
		return ContractInfo{}, errNoClientProvided
	}

	var err error

	err = q.validateNetworkOnIDs(client)
	if err != nil {
		return ContractInfo{}, err
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
			return ContractInfo{}, err
		}

		if cost.tinybar < actualCost.tinybar {
			return ContractInfo{}, ErrMaxQueryPaymentExceeded{
				QueryCost:       actualCost,
				MaxQueryPayment: cost,
				query:           "ContractInfoQuery",
			}
		}

		cost = actualCost
	}

	q.paymentTransactions = make([]*services.Transaction, 0)

	if q.nodeAccountIDs.locked {
		err = q._QueryGeneratePayments(client, cost)
		if err != nil {
			if err != nil {
				return ContractInfo{}, err
			}
		}
	} else {
		paymentTransaction, err := _QueryMakePaymentTransaction(q.paymentTransactionIDs._GetCurrent().(TransactionID), AccountID{}, client.operator, cost)
		if err != nil {
			if err != nil {
				return ContractInfo{}, err
			}
		}
		q.paymentTransactions = append(q.paymentTransactions, paymentTransaction)
	}

	pb := q.build()
	pb.ContractGetInfo.Header = q.pbHeader
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
		return ContractInfo{}, err
	}

	info, err := _ContractInfoFromProtobuf(resp.(*services.Response).GetContractGetInfo().ContractInfo)
	if err != nil {
		return ContractInfo{}, err
	}

	return info, nil
}

// SetMaxQueryPayment sets the maximum payment allowed for this Query.
func (q *ContractInfoQuery) SetMaxQueryPayment(maxPayment Hbar) *ContractInfoQuery {
	q.query.SetMaxQueryPayment(maxPayment)
	return q
}

// SetQueryPayment sets the payment amount for this Query.
func (q *ContractInfoQuery) SetQueryPayment(paymentAmount Hbar) *ContractInfoQuery {
	q.query.SetQueryPayment(paymentAmount)
	return q
}

// SetNodeAccountIDs sets the _Node AccountID for this ContractInfoQuery.
func (q *ContractInfoQuery) SetNodeAccountIDs(accountID []AccountID) *ContractInfoQuery {
	q.query.SetNodeAccountIDs(accountID)
	return q
}

// SetMaxRetry sets the max number of errors before execution will fail.
func (q *ContractInfoQuery) SetMaxRetry(count int) *ContractInfoQuery {
	q.query.SetMaxRetry(count)
	return q
}

// SetMaxBackoff The maximum amount of time to wait between retries.
// Every retry attempt will increase the wait time exponentially until it reaches this time.
func (q *ContractInfoQuery) SetMaxBackoff(max time.Duration) *ContractInfoQuery {
	q.query.SetMaxBackoff(max)
	return q
}

// SetMinBackoff sets the minimum amount of time to wait between retries.
func (q *ContractInfoQuery) SetMinBackoff(min time.Duration) *ContractInfoQuery {
	q.query.SetMinBackoff(min)
	return q
}

func (q *ContractInfoQuery) SetPaymentTransactionID(transactionID TransactionID) *ContractInfoQuery {
	q.paymentTransactionIDs._Clear()._Push(transactionID)._SetLocked(true)
	return q
}

func (q *ContractInfoQuery) SetLogLevel(level LogLevel) *ContractInfoQuery {
	q.query.SetLogLevel(level)
	return q
}

// ---------- Parent functions specific implementation ----------

func (q *ContractInfoQuery) getMethod(channel *_Channel) _Method {
	return _Method{
		query: channel._GetContract().GetContractInfo,
	}
}

func (q *ContractInfoQuery) mapStatusError(_ interface{}, response interface{}) error {
	return ErrHederaPreCheckStatus{
		Status: Status(response.(*services.Response).GetContractGetInfo().Header.NodeTransactionPrecheckCode),
	}
}

func (q *ContractInfoQuery) getName() string {
	return "ContractInfoQuery"
}

func (q *ContractInfoQuery) build() *services.Query_ContractGetInfo {
	pb := services.Query_ContractGetInfo{
		ContractGetInfo: &services.ContractGetInfoQuery{
			Header: &services.QueryHeader{},
		},
	}

	if q.contractID != nil {
		pb.ContractGetInfo.ContractID = q.contractID._ToProtobuf()
	}

	return &pb
}

func (q *ContractInfoQuery) validateNetworkOnIDs(client *Client) error {
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

func (q *ContractInfoQuery) getQueryStatus(response interface{}) Status {
	return Status(response.(*services.Response).GetContractGetInfo().Header.NodeTransactionPrecheckCode)
}
