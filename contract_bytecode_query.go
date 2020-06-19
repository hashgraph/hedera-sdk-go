package hedera

import (
	"github.com/hashgraph/hedera-sdk-go/proto"
)

// ContractBytecodeQuery retrieves the bytecode for a smart contract instance
type ContractBytecodeQuery struct {
	QueryBuilder
	pb *proto.ContractGetBytecodeQuery
}

// NewContractBytecodeQuery creates a ContractBytecodeQuery builder which can be used to construct and execute a
// Contract Get Bytecode Query.
func NewContractBytecodeQuery() *ContractBytecodeQuery {
	pb := &proto.ContractGetBytecodeQuery{Header: &proto.QueryHeader{}}

	inner := newQueryBuilder(pb.Header)
	inner.pb.Query = &proto.Query_ContractGetBytecode{ContractGetBytecode: pb}

	return &ContractBytecodeQuery{inner, pb}
}

// SetContractID sets the contract for which the bytecode is requested
func (builder *ContractBytecodeQuery) SetContractID(id ContractID) *ContractBytecodeQuery {
	builder.pb.ContractID = id.toProto()
	return builder
}

// Execute executes the ContractByteCodeQuery using the provided client
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

// SetMaxQueryPayment sets the maximum payment allowed for this Query.
func (builder *ContractBytecodeQuery) SetMaxQueryPayment(maxPayment Hbar) *ContractBytecodeQuery {
	return &ContractBytecodeQuery{*builder.QueryBuilder.SetMaxQueryPayment(maxPayment), builder.pb}
}

// SetQueryPayment sets the payment amount for this Query.
func (builder *ContractBytecodeQuery) SetQueryPayment(paymentAmount Hbar) *ContractBytecodeQuery {
	return &ContractBytecodeQuery{*builder.QueryBuilder.SetQueryPayment(paymentAmount), builder.pb}
}

// SetQueryPaymentTransaction sets the payment Transaction for this Query.
func (builder *ContractBytecodeQuery) SetQueryPaymentTransaction(tx Transaction) *ContractBytecodeQuery {
	return &ContractBytecodeQuery{*builder.QueryBuilder.SetQueryPaymentTransaction(tx), builder.pb}
}
