package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"time"

	"github.com/hiero-ledger/hiero-sdk-go/v2/proto/services"
)

// ContractCallQuery calls a function of the given smart contract instance, giving it ContractFunctionParameters as its
// inputs. It will consume the entire given amount of gas.
//
// This is performed locally on the particular _Node that the client is communicating with. It cannot change the state of
// the contract instance (and so, cannot spend anything from the instance's Hiero account). It will not have a
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
func (q *ContractCallQuery) SetGrpcDeadline(deadline *time.Duration) *ContractCallQuery {
	q.Query.SetGrpcDeadline(deadline)
	return q
}

// SetContractID sets the contract instance to call
func (q *ContractCallQuery) SetContractID(contractID ContractID) *ContractCallQuery {
	q.contractID = &contractID
	return q
}

// GetContractID returns the contract instance to call
func (q *ContractCallQuery) GetContractID() ContractID {
	if q.contractID == nil {
		return ContractID{}
	}

	return *q.contractID
}

// SetSenderID
// The account that is the "sender." If not present it is the accountId from the transactionId.
// Typically a different value than specified in the transactionId requires a valid signature
// over either the hiero transaction or foreign transaction data.
func (q *ContractCallQuery) SetSenderID(id AccountID) *ContractCallQuery {
	q.senderID = &id
	return q
}

// GetSenderID returns the AccountID that is the "sender."
func (q *ContractCallQuery) GetSenderID() AccountID {
	if q.senderID == nil {
		return AccountID{}
	}

	return *q.senderID
}

// SetGas sets the amount of gas to use for the call. All of the gas offered will be charged for.
func (q *ContractCallQuery) SetGas(gas uint64) *ContractCallQuery {
	q.gas = gas
	return q
}

// GetGas returns the amount of gas to use for the call.
func (q *ContractCallQuery) GetGas() uint64 {
	return q.gas
}

// Deprecated
func (q *ContractCallQuery) SetMaxResultSize(size uint64) *ContractCallQuery {
	q.maxResultSize = size
	return q
}

// SetFunction sets which function to call, and the ContractFunctionParams to pass to the function
func (q *ContractCallQuery) SetFunction(name string, params *ContractFunctionParameters) *ContractCallQuery {
	if params == nil {
		params = NewContractFunctionParameters()
	}

	q.functionParameters = params._Build(&name)
	return q
}

// SetFunctionParameters sets the function parameters as their raw bytes.
func (q *ContractCallQuery) SetFunctionParameters(byteArray []byte) *ContractCallQuery {
	q.functionParameters = byteArray
	return q
}

// GetFunctionParameters returns the function parameters as their raw bytes.
func (q *ContractCallQuery) GetFunctionParameters() []byte {
	return q.functionParameters
}

func (q *ContractCallQuery) GetCost(client *Client) (Hbar, error) {
	return q.Query.getCost(client, q)
}

// Execute executes the Query with the provided client
func (q *ContractCallQuery) Execute(client *Client) (ContractFunctionResult, error) {
	resp, err := q.Query.execute(client, q)

	if err != nil {
		return ContractFunctionResult{}, err
	}

	return _ContractFunctionResultFromProtobuf(resp.GetContractCallLocal().FunctionResult), nil
}

// SetMaxQueryPayment sets the maximum payment allowed for this Query.
func (q *ContractCallQuery) SetMaxQueryPayment(maxPayment Hbar) *ContractCallQuery {
	q.Query.SetMaxQueryPayment(maxPayment)
	return q
}

// SetQueryPayment sets the payment amount for this Query.
func (q *ContractCallQuery) SetQueryPayment(paymentAmount Hbar) *ContractCallQuery {
	q.Query.SetQueryPayment(paymentAmount)
	return q
}

// SetNodeAccountIDs sets the _Node AccountID for this ContractCallQuery.
func (q *ContractCallQuery) SetNodeAccountIDs(accountID []AccountID) *ContractCallQuery {
	q.Query.SetNodeAccountIDs(accountID)
	return q
}

// SetMaxRetry sets the max number of errors before execution will fail.
func (q *ContractCallQuery) SetMaxRetry(count int) *ContractCallQuery {
	q.Query.SetMaxRetry(count)
	return q
}

// SetMaxBackoff The maximum amount of time to wait between retries.
// Every retry attempt will increase the wait time exponentially until it reaches this time.
func (q *ContractCallQuery) SetMaxBackoff(max time.Duration) *ContractCallQuery {
	q.Query.SetMaxBackoff(max)
	return q
}

// SetMinBackoff sets the minimum amount of time to wait between retries.
func (q *ContractCallQuery) SetMinBackoff(min time.Duration) *ContractCallQuery {
	q.Query.SetMinBackoff(min)
	return q
}

// SetPaymentTransactionID assigns the payment transaction id.
func (q *ContractCallQuery) SetPaymentTransactionID(transactionID TransactionID) *ContractCallQuery {
	q.paymentTransactionIDs._Clear()._Push(transactionID)._SetLocked(true)
	return q
}

func (q *ContractCallQuery) SetLogLevel(level LogLevel) *ContractCallQuery {
	q.Query.SetLogLevel(level)
	return q
}

// ---------- Parent functions specific implementation ----------

func (q *ContractCallQuery) getMethod(channel *_Channel) _Method {
	return _Method{
		query: channel._GetContract().ContractCallLocalMethod,
	}
}

func (q *ContractCallQuery) getName() string {
	return "ContractCallQuery"
}

func (q *ContractCallQuery) buildQuery() *services.Query {
	pb := services.Query_ContractCallLocal{
		ContractCallLocal: &services.ContractCallLocalQuery{
			Header: q.pbHeader,
			Gas:    int64(q.gas),
		},
	}

	if q.contractID != nil {
		pb.ContractCallLocal.ContractID = q.contractID._ToProtobuf()
	}

	if q.senderID != nil {
		pb.ContractCallLocal.SenderId = q.senderID._ToProtobuf()
	}

	if len(q.functionParameters) > 0 {
		pb.ContractCallLocal.FunctionParameters = q.functionParameters
	}

	return &services.Query{
		Query: &pb,
	}
}

func (q *ContractCallQuery) validateNetworkOnIDs(client *Client) error {
	if client == nil || !client.autoValidateChecksums {
		return nil
	}

	if q.contractID != nil {
		if err := q.contractID.ValidateChecksum(client); err != nil {
			return err
		}
	}

	if q.senderID != nil {
		if err := q.senderID.ValidateChecksum(client); err != nil {
			return err
		}
	}

	return nil
}

func (q *ContractCallQuery) getQueryResponse(response *services.Response) queryResponse {
	return response.GetContractCallLocal()
}
