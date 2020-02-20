package hedera

import (
	"github.com/hashgraph/hedera-sdk-go/proto"
)

// ContractCallQuery calls a function of the given smart contract instance, giving it ContractFunctionParameters as its
// inputs. It will consume the entire given amount of gas.
//
// This is performed locally on the particular node that the client is communicating with. It cannot change the state of
// the contract instance (and so, cannot spend anything from the instance's Hedera account). It will not have a
// consensus timestamp. It cannot generate a record or a receipt. This is useful for calling getter functions, which
// purely read the state and don't change it. It is faster and cheaper than a ContractExecuteTransaction, because it is
// purely local to a single  node.
type ContractCallQuery struct {
	QueryBuilder
	pb *proto.ContractCallLocalQuery
}

// NewContractCallQuery creates a ContractCallQuery builder which can be used to construct and execute a
// Contract Call Local Query.
func NewContractCallQuery() *ContractCallQuery {
	pb := &proto.ContractCallLocalQuery{Header: &proto.QueryHeader{}}

	inner := newQueryBuilder(pb.Header)
	inner.pb.Query = &proto.Query_ContractCallLocal{ContractCallLocal: pb}

	return &ContractCallQuery{inner, pb}
}

// SetContractID sets the contract instance to call
func (builder *ContractCallQuery) SetContractID(id ContractID) *ContractCallQuery {
	builder.pb.ContractID = id.toProto()
	return builder
}

// SetGas sets the amount of gas to use for the call. All of the gas offered will be charged for.
func (builder *ContractCallQuery) SetGas(gas uint64) *ContractCallQuery {
	builder.pb.Gas = int64(gas)
	return builder
}

// SetMaxResultSize sets the max number of bytes that the result might include. The run will fail if it would have
// returned more than this number of bytes.
func (builder *ContractCallQuery) SetMaxResultSize(size uint64) *ContractCallQuery {
	builder.pb.MaxResultSize = int64(size)
	return builder
}

// SetFunction sets which function to call, and the ContractFunctionParams to pass to the function
func (builder *ContractCallQuery) SetFunction(name string, params *ContractFunctionParams) *ContractCallQuery {
	if params == nil {
		params = NewContractFunctionParams()
	}

	builder.pb.FunctionParameters = params.build(&name)
	return builder
}

// Execute executes the ContractCallQuery using the provided client. The returned ContractFunctionResult will contain
// the output returned by the function call.
func (builder *ContractCallQuery) Execute(client *Client) (ContractFunctionResult, error) {
	resp, err := builder.execute(client)
	if err != nil {
		return ContractFunctionResult{}, err
	}

	return contractFunctionResultFromProto(resp.GetContractCallLocal().FunctionResult), nil
}

// Cost is a wrapper around the standard Cost function for a query. It must exist because the cost returned the standard
// QueryBuilder.Cost() function and therein the Hedera Network doesn't work for ContractCallQueries. However, if the
// value returned by the network is increased by 10% then most contractCallQueries will complete fine.
func (builder *ContractCallQuery) Cost(client *Client) (Hbar, error) {
	cost, err := builder.QueryBuilder.GetCost(client)
	if err != nil {
		return ZeroHbar, err
	}

	return HbarFromTinybar(int64(float64(cost.AsTinybar()) * 1.1)), nil
}

//
// The following _3_ must be copy-pasted at the bottom of **every** _query.go file
// We override the embedded fluent setter methods to return the outer type
//

func (builder *ContractCallQuery) SetMaxQueryPayment(maxPayment Hbar) *ContractCallQuery {
	return &ContractCallQuery{*builder.QueryBuilder.SetMaxQueryPayment(maxPayment), builder.pb}
}

func (builder *ContractCallQuery) SetQueryPayment(paymentAmount Hbar) *ContractCallQuery {
	return &ContractCallQuery{*builder.QueryBuilder.SetQueryPayment(paymentAmount), builder.pb}
}

func (builder *ContractCallQuery) SetQueryPaymentTransaction(tx Transaction) *ContractCallQuery {
	return &ContractCallQuery{*builder.QueryBuilder.SetQueryPaymentTransaction(tx), builder.pb}
}
