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
	inner.pb.Query = &proto.Query_ContractCallLocal{pb}

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

func (builder *ContractCallQuery) SetFunctionParameters(params ContractFunctionParams) (*ContractCallQuery, error) {
	function, err := params.build(nil)
	if err != nil {
		return builder, err
	}

	builder.pb.FunctionParameters = function
	return builder, nil
}

func (builder *ContractCallQuery) Execute(client *Client) (ContractFunctionResult, error) {
	resp, err := builder.execute(client)
	if err != nil {
		return ContractFunctionResult{}, err
	}

	return contractFunctionResultFromProto(resp.GetContractCallLocal().FunctionResult), nil
}
