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
	query
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

	result := ContractCallQuery{
		query: query,
	}

	result.e = &result
	return &result
}

// When execution is attempted, a single attempt will timeout when this deadline is reached. (The SDK may subsequently retry the execution.)
func (this *ContractCallQuery) SetGrpcDeadline(deadline *time.Duration) *ContractCallQuery {
	this.query.SetGrpcDeadline(deadline)
	return this
}

// SetContractID sets the contract instance to call
func (this *ContractCallQuery) SetContractID(contractID ContractID) *ContractCallQuery {
	this.contractID = &contractID
	return this
}

// GetContractID returns the contract instance to call
func (this *ContractCallQuery) GetContractID() ContractID {
	if this.contractID == nil {
		return ContractID{}
	}

	return *this.contractID
}

// SetSenderID
// The account that is the "sender." If not present it is the accountId from the transactionId.
// Typically a different value than specified in the transactionId requires a valid signature
// over either the hedera transaction or foreign transaction data.
func (this *ContractCallQuery) SetSenderID(id AccountID) *ContractCallQuery {
	this.senderID = &id
	return this
}

// GetSenderID returns the AccountID that is the "sender."
func (this *ContractCallQuery) GetSenderID() AccountID {
	if this.senderID == nil {
		return AccountID{}
	}

	return *this.senderID
}

// SetGas sets the amount of gas to use for the call. All of the gas offered will be charged for.
func (this *ContractCallQuery) SetGas(gas uint64) *ContractCallQuery {
	this.gas = gas
	return this
}

// GetGas returns the amount of gas to use for the call.
func (this *ContractCallQuery) GetGas() uint64 {
	return this.gas
}

// Deprecated
func (this *ContractCallQuery) SetMaxResultSize(size uint64) *ContractCallQuery {
	this.maxResultSize = size
	return this
}

// SetFunction sets which function to call, and the ContractFunctionParams to pass to the function
func (this *ContractCallQuery) SetFunction(name string, params *ContractFunctionParameters) *ContractCallQuery {
	if params == nil {
		params = NewContractFunctionParameters()
	}

	this.functionParameters = params._Build(&name)
	return this
}

// SetFunctionParameters sets the function parameters as their raw bytes.
func (this *ContractCallQuery) SetFunctionParameters(byteArray []byte) *ContractCallQuery {
	this.functionParameters = byteArray
	return this
}

// GetFunctionParameters returns the function parameters as their raw bytes.
func (this *ContractCallQuery) GetFunctionParameters() []byte {
	return this.functionParameters
}

// GetCost returns the fee that would be charged to get the requested information (if a cost was requested).
func (this *ContractCallQuery) GetCost(client *Client) (Hbar, error) {
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
	pb.ContractCallLocal.Header = this.pbHeader

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

	cost := int64(resp.(*services.Response).GetContractCallLocal().Header.Cost)
	return HbarFromTinybar(cost), nil
}

// Execute executes the Query with the provided client
func (this *ContractCallQuery) Execute(client *Client) (ContractFunctionResult, error) {
	if client == nil || client.operator == nil {
		return ContractFunctionResult{}, errNoClientProvided
	}

	var err error

	err = this.validateNetworkOnIDs(client)
	if err != nil {
		return ContractFunctionResult{}, err
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

	this.paymentTransactions = make([]*services.Transaction, 0)

	if this.nodeAccountIDs.locked {
		err = this._QueryGeneratePayments(client, cost)
		if err != nil {
			return ContractFunctionResult{}, err
		}
	} else {
		paymentTransaction, err := _QueryMakePaymentTransaction(this.paymentTransactionIDs._GetCurrent().(TransactionID), AccountID{}, client.operator, cost)
		if err != nil {
			return ContractFunctionResult{}, err
		}
		this.paymentTransactions = append(this.paymentTransactions, paymentTransaction)
	}

	pb := this.build()
	pb.ContractCallLocal.Header = this.pbHeader
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
		return ContractFunctionResult{}, err
	}

	return _ContractFunctionResultFromProtobuf(resp.(*services.Response).GetContractCallLocal().FunctionResult), nil
}

// SetMaxQueryPayment sets the maximum payment allowed for this Query.
func (this *ContractCallQuery) SetMaxQueryPayment(maxPayment Hbar) *ContractCallQuery {
	this.query.SetMaxQueryPayment(maxPayment)
	return this
}

// SetQueryPayment sets the payment amount for this Query.
func (this *ContractCallQuery) SetQueryPayment(paymentAmount Hbar) *ContractCallQuery {
	this.query.SetQueryPayment(paymentAmount)
	return this
}

// SetNodeAccountIDs sets the _Node AccountID for this ContractCallQuery.
func (this *ContractCallQuery) SetNodeAccountIDs(accountID []AccountID) *ContractCallQuery {
	this.query.SetNodeAccountIDs(accountID)
	return this
}

// SetMaxRetry sets the max number of errors before execution will fail.
func (this *ContractCallQuery) SetMaxRetry(count int) *ContractCallQuery {
	this.query.SetMaxRetry(count)
	return this
}

// SetMaxBackoff The maximum amount of time to wait between retries.
// Every retry attempt will increase the wait time exponentially until it reaches this time.
func (this *ContractCallQuery) SetMaxBackoff(max time.Duration) *ContractCallQuery {
	this.query.SetMaxBackoff(max)
	return this
}

// GetMaxBackoff returns the maximum amount of time to wait between retries.
func (this *ContractCallQuery) GetMaxBackoff() time.Duration {
	return this.query.GetMaxBackoff()
}

// SetMinBackoff sets the minimum amount of time to wait between retries.
func (this *ContractCallQuery) SetMinBackoff(min time.Duration) *ContractCallQuery {
	this.query.SetMinBackoff(min)
	return this
}

// GetMinBackoff returns the minimum amount of time to wait between retries.
func (this *ContractCallQuery) GetMinBackoff() time.Duration {
	return this.query.GetMinBackoff()
}

func (this *ContractCallQuery) _GetLogID() string {
	timestamp := this.timestamp.UnixNano()
	if this.paymentTransactionIDs._Length() > 0 && this.paymentTransactionIDs._GetCurrent().(TransactionID).ValidStart != nil {
		timestamp = this.paymentTransactionIDs._GetCurrent().(TransactionID).ValidStart.UnixNano()
	}
	return fmt.Sprintf("ContractCallQuery:%d", timestamp)
}

// SetPaymentTransactionID assigns the payment transaction id.
func (this *ContractCallQuery) SetPaymentTransactionID(transactionID TransactionID) *ContractCallQuery {
	this.paymentTransactionIDs._Clear()._Push(transactionID)._SetLocked(true)
	return this
}

func (this *ContractCallQuery) SetLogLevel(level LogLevel) *ContractCallQuery {
	this.query.SetLogLevel(level)
	return this
}

// ---------- Parent functions specific implementation ----------

func (this *ContractCallQuery) getMethod(channel *_Channel) _Method {
	return _Method{
		query: channel._GetContract().ContractCallLocalMethod,
	}
}

func (this *ContractCallQuery) mapStatusError(_ interface{}, response interface{}) error {
	return ErrHederaPreCheckStatus{
		Status: Status(response.(*services.Response).GetContractCallLocal().Header.NodeTransactionPrecheckCode),
	}
}

func (this *ContractCallQuery) getName() string {
	return "ContractCallQuery"
}

func (this *ContractCallQuery) build() *services.Query_ContractCallLocal {
	pb := services.Query_ContractCallLocal{
		ContractCallLocal: &services.ContractCallLocalQuery{
			Header: &services.QueryHeader{},
			Gas:    int64(this.gas),
		},
	}

	if this.contractID != nil {
		pb.ContractCallLocal.ContractID = this.contractID._ToProtobuf()
	}

	if this.senderID != nil {
		pb.ContractCallLocal.SenderId = this.senderID._ToProtobuf()
	}

	if len(this.functionParameters) > 0 {
		pb.ContractCallLocal.FunctionParameters = this.functionParameters
	}

	return &pb
}

func (this *ContractCallQuery) validateNetworkOnIDs(client *Client) error {
	if client == nil || !client.autoValidateChecksums {
		return nil
	}

	if this.contractID != nil {
		if err := this.contractID.ValidateChecksum(client); err != nil {
			return err
		}
	}

	if this.senderID != nil {
		if err := this.senderID.ValidateChecksum(client); err != nil {
			return err
		}
	}

	return nil
}

func (this *ContractCallQuery) getQueryStatus(response interface{}) Status {
	return Status(response.(*services.Response).GetContractCallLocal().Header.NodeTransactionPrecheckCode)
}
