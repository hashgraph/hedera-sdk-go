package hedera

import (
	"github.com/hashgraph/hedera-sdk-go/proto"
)

type ContractBytecodeQuery struct {
	QueryBuilder
	pb *proto.ContractGetBytecodeQuery
}

func NewContractBytecodeQuery() *ContractBytecodeQuery {
	pb := &proto.ContractGetBytecodeQuery{Header: &proto.QueryHeader{}}

	inner := newQueryBuilder(pb.Header)
	inner.pb.Query = &proto.Query_ContractGetBytecode{pb}

	return &ContractBytecodeQuery{inner, pb}
}

func (builder *ContractBytecodeQuery) SetContractId(id ContractID) *ContractBytecodeQuery {
	builder.pb.ContractID = id.toProto()
	return builder
}

func (builder *ContractBytecodeQuery) Execute(client *Client) ([]byte, error) {
	resp, err := builder.execute(client)
	if err != nil {
		return []byte{}, err
	}

	return resp.GetContractGetBytecodeResponse().Bytecode, nil
}
