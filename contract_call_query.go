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

// ContractCallQuery calls a function of the given smart contract instance, giving it ContractFunctionParameters as its
// inputs. It will consume the entire given amount of gas.
//
// This is performed locally on the particular _Node that the client is communicating with. It cannot change the state of
// the contract instance (and so, cannot spend anything from the instance's Hedera account). It will not have a
// consensus timestamp. It cannot generate a record or a receipt. This is useful for calling getter functions, which
// purely read the state and don't change it. It is faster and cheaper than a ContractExecuteTransaction, because it is
// purely local to a single  _Node.
type ContractCallQuery struct {
	Query
	contractID         *ContractID
	gas                uint64
	maxResultSize      uint64
	functionParameters []byte
	senderID           *AccountID
}

// NewContractCallQuery creates a ContractCallQuery query which can be used to construct and execute a
// Contract Call Local Query.
func NewContractCallQuery() *ContractCallQuery {
	header := services.QueryHeader{}
	query := _NewQuery(true, &header)

	return &ContractCallQuery{
		Query: query,
	}
}

// When execution is attempted, a single attempt will timeout when this deadline is reached. (The SDK may subsequently retry the execution.)
func (query *ContractCallQuery) SetGrpcDeadline(deadline *time.Duration) *ContractCallQuery {
	query.Query.SetGrpcDeadline(deadline)
	return query
}

// SetContractID sets the contract instance to call
func (query *ContractCallQuery) SetContractID(contractID ContractID) *ContractCallQuery {
	query.contractID = &contractID
	return query
}

// GetContractID returns the contract instance to call
func (query *ContractCallQuery) GetContractID() ContractID {
	if query.contractID == nil {
		return ContractID{}
	}

	return *query.contractID
}

// SetSenderID
// The account that is the "sender." If not present it is the accountId from the transactionId.
// Typically a different value than specified in the transactionId requires a valid signature
// over either the hedera transaction or foreign transaction data.
func (query *ContractCallQuery) SetSenderID(id AccountID) *ContractCallQuery {
	query.senderID = &id
	return query
}

// GetSenderID returns the AccountID that is the "sender."
func (query *ContractCallQuery) GetSenderID() AccountID {
	if query.senderID == nil {
		return AccountID{}
	}

	return *query.senderID
}

// SetGas sets the amount of gas to use for the call. All of the gas offered will be charged for.
func (query *ContractCallQuery) SetGas(gas uint64) *ContractCallQuery {
	query.gas = gas
	return query
}

// GetGas returns the amount of gas to use for the call.
func (query *ContractCallQuery) GetGas() uint64 {
	return query.gas
}

// Deprecated
func (query *ContractCallQuery) SetMaxResultSize(size uint64) *ContractCallQuery {
	query.maxResultSize = size
	return query
}

// SetFunction sets which function to call, and the ContractFunctionParams to pass to the function
func (query *ContractCallQuery) SetFunction(name string, params *ContractFunctionParameters) *ContractCallQuery {
	if params == nil {
		params = NewContractFunctionParameters()
	}

	query.functionParameters = params._Build(&name)
	return query
}

// SetFunctionParameters sets the function parameters as their raw bytes.
func (query *ContractCallQuery) SetFunctionParameters(byteArray []byte) *ContractCallQuery {
	query.functionParameters = byteArray
	return query
}

// GetFunctionParameters returns the function parameters as their raw bytes.
func (query *ContractCallQuery) GetFunctionParameters() []byte {
	return query.functionParameters
}

func (query *ContractCallQuery) _ValidateNetworkOnIDs(client *Client) error {
	if client == nil || !client.autoValidateChecksums {
		return nil
	}

	if query.contractID != nil {
		if err := query.contractID.ValidateChecksum(client); err != nil {
			return err
		}
	}

	if query.senderID != nil {
		if err := query.senderID.ValidateChecksum(client); err != nil {
			return err
		}
	}

	return nil
}

func (query *ContractCallQuery) _Build() *services.Query_ContractCallLocal {
	pb := services.Query_ContractCallLocal{
		ContractCallLocal: &services.ContractCallLocalQuery{
			Header: &services.QueryHeader{},
			Gas:    int64(query.gas),
		},
	}

	if query.contractID != nil {
		pb.ContractCallLocal.ContractID = query.contractID._ToProtobuf()
	}

	if query.senderID != nil {
		pb.ContractCallLocal.SenderId = query.senderID._ToProtobuf()
	}

	if len(query.functionParameters) > 0 {
		pb.ContractCallLocal.FunctionParameters = query.functionParameters
	}

	return &pb
}

// GetCost returns the fee that would be charged to get the requested information (if a cost was requested).
func (query *ContractCallQuery) GetCost(client *Client) (Hbar, error) {
	if client == nil || client.operator == nil {
		return Hbar{}, errNoClientProvided
	}

	var err error

	err = query._ValidateNetworkOnIDs(client)
	if err != nil {
		return Hbar{}, err
	}

	for range query.nodeAccountIDs.slice {
		paymentTransaction, err := _QueryMakePaymentTransaction(TransactionID{}, AccountID{}, client.operator, Hbar{})
		if err != nil {
			return Hbar{}, err
		}
		query.paymentTransactions = append(query.paymentTransactions, paymentTransaction)
	}

	pb := query._Build()
	pb.ContractCallLocal.Header = query.pbHeader

	query.pb = &services.Query{
		Query: pb,
	}

	resp, err := _Execute(
		client,
		&query.Query,
		_ContractCallQueryShouldRetry,
		_CostQueryMakeRequest,
		_CostQueryAdvanceRequest,
		_QueryGetNodeAccountID,
		_ContractCallQueryGetMethod,
		_ContractCallQueryMapStatusError,
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

	cost := int64(resp.(*services.Response).GetContractCallLocal().Header.Cost)
	return HbarFromTinybar(cost), nil
}

func _ContractCallQueryShouldRetry(_ interface{}, response interface{}) _ExecutionState {
	return _QueryShouldRetry(Status(response.(*services.Response).GetContractCallLocal().Header.NodeTransactionPrecheckCode))
}

func _ContractCallQueryMapStatusError(_ interface{}, response interface{}) error {
	return ErrHederaPreCheckStatus{
		Status: Status(response.(*services.Response).GetContractCallLocal().Header.NodeTransactionPrecheckCode),
	}
}

func _ContractCallQueryGetMethod(_ interface{}, channel *_Channel) _Method {
	return _Method{
		query: channel._GetContract().ContractCallLocalMethod,
	}
}

// Execute executes the Query with the provided client
func (query *ContractCallQuery) Execute(client *Client) (ContractFunctionResult, error) {
	if client == nil || client.operator == nil {
		return ContractFunctionResult{}, errNoClientProvided
	}

	var err error

	err = query._ValidateNetworkOnIDs(client)
	if err != nil {
		return ContractFunctionResult{}, err
	}

	if !query.paymentTransactionIDs.locked {
		query.paymentTransactionIDs._Clear()._Push(TransactionIDGenerate(client.operator.accountID))
	}

	var cost Hbar
	if query.queryPayment.tinybar != 0 {
		cost = query.queryPayment
	} else {
		if query.maxQueryPayment.tinybar == 0 {
			cost = client.GetDefaultMaxQueryPayment()
		} else {
			cost = query.maxQueryPayment
		}

		actualCost, err := query.GetCost(client)
		if err != nil {
			return ContractFunctionResult{}, err
		}

		if cost.tinybar < actualCost.tinybar {
			return ContractFunctionResult{}, ErrMaxQueryPaymentExceeded{
				QueryCost:       actualCost,
				MaxQueryPayment: cost,
				query:           "ContractFunctionResultQuery",
			}
		}

		cost = actualCost
	}

	query.paymentTransactions = make([]*services.Transaction, 0)

	if query.nodeAccountIDs.locked {
		err = _QueryGeneratePayments(&query.Query, client, cost)
		if err != nil {
			return ContractFunctionResult{}, err
		}
	} else {
		paymentTransaction, err := _QueryMakePaymentTransaction(query.paymentTransactionIDs._GetCurrent().(TransactionID), AccountID{}, client.operator, cost)
		if err != nil {
			return ContractFunctionResult{}, err
		}
		query.paymentTransactions = append(query.paymentTransactions, paymentTransaction)
	}

	pb := query._Build()
	pb.ContractCallLocal.Header = query.pbHeader
	query.pb = &services.Query{
		Query: pb,
	}

	resp, err := _Execute(
		client,
		&query.Query,
		_ContractCallQueryShouldRetry,
		_QueryMakeRequest,
		_QueryAdvanceRequest,
		_QueryGetNodeAccountID,
		_ContractCallQueryGetMethod,
		_ContractCallQueryMapStatusError,
		_QueryMapResponse,
		query._GetLogID(),
		query.grpcDeadline,
		query.maxBackoff,
		query.minBackoff,
		query.maxRetry,
	)

	if err != nil {
		return ContractFunctionResult{}, err
	}

	return _ContractFunctionResultFromProtobuf(resp.(*services.Response).GetContractCallLocal().FunctionResult), nil
}

// SetMaxQueryPayment sets the maximum payment allowed for this Query.
func (query *ContractCallQuery) SetMaxQueryPayment(maxPayment Hbar) *ContractCallQuery {
	query.Query.SetMaxQueryPayment(maxPayment)
	return query
}

// SetQueryPayment sets the payment amount for this Query.
func (query *ContractCallQuery) SetQueryPayment(paymentAmount Hbar) *ContractCallQuery {
	query.Query.SetQueryPayment(paymentAmount)
	return query
}

// SetNodeAccountIDs sets the _Node AccountID for this ContractCallQuery.
func (query *ContractCallQuery) SetNodeAccountIDs(accountID []AccountID) *ContractCallQuery {
	query.Query.SetNodeAccountIDs(accountID)
	return query
}

// SetMaxRetry sets the max number of errors before execution will fail.
func (query *ContractCallQuery) SetMaxRetry(count int) *ContractCallQuery {
	query.Query.SetMaxRetry(count)
	return query
}

// SetMaxBackoff The maximum amount of time to wait between retries.
// Every retry attempt will increase the wait time exponentially until it reaches this time.
func (query *ContractCallQuery) SetMaxBackoff(max time.Duration) *ContractCallQuery {
	if max.Nanoseconds() < 0 {
		panic("maxBackoff must be a positive duration")
	} else if max.Nanoseconds() < query.minBackoff.Nanoseconds() {
		panic("maxBackoff must be greater than or equal to minBackoff")
	}
	query.maxBackoff = &max
	return query
}

// GetMaxBackoff returns the maximum amount of time to wait between retries.
func (query *ContractCallQuery) GetMaxBackoff() time.Duration {
	if query.maxBackoff != nil {
		return *query.maxBackoff
	}

	return 8 * time.Second
}

// SetMinBackoff sets the minimum amount of time to wait between retries.
func (query *ContractCallQuery) SetMinBackoff(min time.Duration) *ContractCallQuery {
	if min.Nanoseconds() < 0 {
		panic("minBackoff must be a positive duration")
	} else if query.maxBackoff.Nanoseconds() < min.Nanoseconds() {
		panic("minBackoff must be less than or equal to maxBackoff")
	}
	query.minBackoff = &min
	return query
}

// GetMinBackoff returns the minimum amount of time to wait between retries.
func (query *ContractCallQuery) GetMinBackoff() time.Duration {
	if query.minBackoff != nil {
		return *query.minBackoff
	}

	return 250 * time.Millisecond
}

func (query *ContractCallQuery) _GetLogID() string {
	timestamp := query.timestamp.UnixNano()
	if query.paymentTransactionIDs._Length() > 0 && query.paymentTransactionIDs._GetCurrent().(TransactionID).ValidStart != nil {
		timestamp = query.paymentTransactionIDs._GetCurrent().(TransactionID).ValidStart.UnixNano()
	}
	return fmt.Sprintf("ContractCallQuery:%d", timestamp)
}

// SetPaymentTransactionID assigns the payment transaction id.
func (query *ContractCallQuery) SetPaymentTransactionID(transactionID TransactionID) *ContractCallQuery {
	query.paymentTransactionIDs._Clear()._Push(transactionID)._SetLocked(true)
	return query
}

func (query *ContractCallQuery) SetLogLevel(level LogLevel) *ContractCallQuery {
	query.Query.SetLogLevel(level)
	return query
}
