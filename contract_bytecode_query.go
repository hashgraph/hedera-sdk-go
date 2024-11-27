package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"time"

	"github.com/hiero-ledger/hiero-sdk-go/v2/proto/services"
)

// ContractBytecodeQuery retrieves the bytecode for a smart contract instance
type ContractBytecodeQuery struct {
	Query
	contractID *ContractID
}

// NewContractBytecodeQuery creates a ContractBytecodeQuery query which can be used to construct and execute a
// Contract Get Bytecode Query.
func NewContractBytecodeQuery() *ContractBytecodeQuery {
	header := services.QueryHeader{}
	return &ContractBytecodeQuery{
		Query: _NewQuery(true, &header),
	}
}

// When execution is attempted, a single attempt will timeout when this deadline is reached. (The SDK may subsequently retry the execution.)
func (q *ContractBytecodeQuery) SetGrpcDeadline(deadline *time.Duration) *ContractBytecodeQuery {
	q.Query.SetGrpcDeadline(deadline)
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
	return q.Query.getCost(client, q)
}

// Execute executes the Query with the provided client
func (q *ContractBytecodeQuery) Execute(client *Client) ([]byte, error) {
	resp, err := q.Query.execute(client, q)

	if err != nil {
		return []byte{}, err
	}

	return resp.GetContractGetBytecodeResponse().Bytecode, nil
}

// SetMaxQueryPayment sets the maximum payment allowed for this Query.
func (q *ContractBytecodeQuery) SetMaxQueryPayment(maxPayment Hbar) *ContractBytecodeQuery {
	q.Query.SetMaxQueryPayment(maxPayment)
	return q
}

// SetQueryPayment sets the payment amount for this Query.
func (q *ContractBytecodeQuery) SetQueryPayment(paymentAmount Hbar) *ContractBytecodeQuery {
	q.Query.SetQueryPayment(paymentAmount)
	return q
}

// SetNodeAccountIDs sets the _Node AccountID for this ContractBytecodeQuery.
func (q *ContractBytecodeQuery) SetNodeAccountIDs(accountID []AccountID) *ContractBytecodeQuery {
	q.Query.SetNodeAccountIDs(accountID)
	return q
}

// SetMaxRetry sets the max number of errors before execution will fail.
func (q *ContractBytecodeQuery) SetMaxRetry(count int) *ContractBytecodeQuery {
	q.Query.SetMaxRetry(count)
	return q
}

// SetMaxBackoff The maximum amount of time to wait between retries.
// Every retry attempt will increase the wait time exponentially until it reaches this time.
func (q *ContractBytecodeQuery) SetMaxBackoff(max time.Duration) *ContractBytecodeQuery {
	q.Query.SetMaxBackoff(max)
	return q
}

// SetMinBackoff sets the minimum amount of time to wait between retries.
func (q *ContractBytecodeQuery) SetMinBackoff(min time.Duration) *ContractBytecodeQuery {
	q.Query.SetMinBackoff(min)
	return q
}

// SetPaymentTransactionID assigns the payment transaction id.
func (q *ContractBytecodeQuery) SetPaymentTransactionID(transactionID TransactionID) *ContractBytecodeQuery {
	q.paymentTransactionIDs._Clear()._Push(transactionID)._SetLocked(true)
	return q
}

func (q *ContractBytecodeQuery) SetLogLevel(level LogLevel) *ContractBytecodeQuery {
	q.Query.SetLogLevel(level)
	return q
}

// ---------- Parent functions specific implementation ----------

func (q *ContractBytecodeQuery) getMethod(channel *_Channel) _Method {
	return _Method{
		query: channel._GetContract().ContractGetBytecode}
}

func (q *ContractBytecodeQuery) getName() string {
	return "ContractBytecodeQuery"
}

func (q *ContractBytecodeQuery) buildQuery() *services.Query {
	pb := services.Query_ContractGetBytecode{
		ContractGetBytecode: &services.ContractGetBytecodeQuery{
			Header: q.pbHeader,
		},
	}

	if q.contractID != nil {
		pb.ContractGetBytecode.ContractID = q.contractID._ToProtobuf()
	}

	return &services.Query{
		Query: &pb,
	}
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

func (q *ContractBytecodeQuery) getQueryResponse(response *services.Response) queryResponse {
	return response.GetContractGetBytecodeResponse()
}
