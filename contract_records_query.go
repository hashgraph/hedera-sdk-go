package hedera

import "github.com/hashgraph/hedera-sdk-go/proto"

type ContractRecordsQuery struct {
	QueryBuilder
	pb *proto.ContractGetRecordsQuery
}

func NewContractRecordsQuery() *ContractRecordsQuery {
	pb := &proto.ContractGetRecordsQuery{Header: &proto.QueryHeader{}}

	inner := newQueryBuilder(pb.Header)
	inner.pb.Query = &proto.Query_ContractGetRecords{ContractGetRecords: pb}

	return &ContractRecordsQuery{inner, pb}
}

func (builder *ContractRecordsQuery) SetContractID(id ContractID) *ContractRecordsQuery {
	builder.pb.ContractID = id.toProto()
	return builder
}

func (builder *ContractRecordsQuery) Execute(client *Client) ([]TransactionRecord, error) {
	var records = []TransactionRecord{}

	resp, err := builder.execute(client)
	if err != nil {
		return records, err
	}

	for _, element := range resp.GetContractGetRecordsResponse().Records {
		records = append(records, transactionRecordFromProto(element))
	}

	return records, nil
}

//
// The following _3_ must be copy-pasted at the bottom of **every** _query.go file
// We override the embedded fluent setter methods to return the outer type
//

func (builder *ContractRecordsQuery) SetMaxQueryPayment(maxPayment Hbar) *ContractRecordsQuery {
	return &ContractRecordsQuery{*builder.QueryBuilder.SetMaxQueryPayment(maxPayment), builder.pb}
}

func (builder *ContractRecordsQuery) SetQueryPayment(paymentAmount Hbar) *ContractRecordsQuery {
	return &ContractRecordsQuery{*builder.QueryBuilder.SetQueryPayment(paymentAmount), builder.pb}
}

func (builder *ContractRecordsQuery) SetQueryPaymentTransaction(tx Transaction) *ContractRecordsQuery {
	return &ContractRecordsQuery{*builder.QueryBuilder.SetQueryPaymentTransaction(tx), builder.pb}
}
