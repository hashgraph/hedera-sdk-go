package hedera

import (
	"github.com/hashgraph/hedera-sdk-go/proto"
)

type ContractCallQuery struct {
	QueryBuilder
	pb *proto.ContractCallLocalQuery
}

func NewContractCallQuery() *ContractCallQuery {
	pb := &proto.ContractCallLocalQuery{Header: &proto.QueryHeader{}}

	inner := newQueryBuilder(pb.Header)
	inner.pb.Query = &proto.Query_ContractCallLocal{ContractCallLocal: pb}

	return &ContractCallQuery{inner, pb}
}

func (builder *ContractCallQuery) SetContractID(id ContractID) *ContractCallQuery {
	builder.pb.ContractID = id.toProto()
	return builder
}

func (builder *ContractCallQuery) SetGas(gas uint64) *ContractCallQuery {
	builder.pb.Gas = int64(gas)
	return builder
}

func (builder *ContractCallQuery) SetMaxResultSize(size uint64) *ContractCallQuery {
	builder.pb.MaxResultSize = int64(size)
	return builder
}

func (builder *ContractCallQuery) SetFunction(name string, params *ContractFunctionParams) *ContractCallQuery {
	if params == nil {
		params = NewContractFunctionParams()
	}

	builder.pb.FunctionParameters = params.build(&name)
	return builder
}

func (builder *ContractCallQuery) Execute(client *Client) (ContractFunctionResult, error) {
	resp, err := builder.execute(client)
	if err != nil {
		return ContractFunctionResult{}, err
	}

	return contractFunctionResultFromProto(resp.GetContractCallLocal().FunctionResult), nil
}

func (builder *ContractCallQuery) Cost(client *Client) (Hbar, error) {
	cost, err := builder.QueryBuilder.Cost(client)
	if err != nil {
		return ZeroHbar, err
	}

	// TODO: Document why
	return HbarFromTinybar(int64(float64(cost.AsTinybar()) * float64(1.1))), nil
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
