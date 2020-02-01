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
	inner.pb.Query = &proto.Query_ContractGetBytecode{ContractGetBytecode: pb}

	return &ContractBytecodeQuery{inner, pb}
}

func (builder *ContractBytecodeQuery) SetContractID(id ContractID) *ContractBytecodeQuery {
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

//
// The following _3_ must be copy-pasted at the bottom of **every** _query.go file
// We override the embedded fluent setter methods to return the outer type
//

func (builder *ContractBytecodeQuery) SetMaxQueryPayment(maxPayment Hbar) *ContractBytecodeQuery {
	return &ContractBytecodeQuery{*builder.QueryBuilder.SetMaxQueryPayment(maxPayment), builder.pb}
}

func (builder *ContractBytecodeQuery) SetQueryPayment(paymentAmount Hbar) *ContractBytecodeQuery {
	return &ContractBytecodeQuery{*builder.QueryBuilder.SetQueryPayment(paymentAmount), builder.pb}
}

func (builder *ContractBytecodeQuery) SetQueryPaymentTransaction(tx Transaction) *ContractBytecodeQuery {
	return &ContractBytecodeQuery{*builder.QueryBuilder.SetQueryPaymentTransaction(tx), builder.pb}
}
