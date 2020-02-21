package hedera

import "github.com/hashgraph/hedera-sdk-go/proto"

// ContractRecordsQuery retrieves all of the records for a smart contract instance, for any function call
// (or the constructor call) during the last 25 hours, for which a Record was requested.
type ContractRecordsQuery struct {
	QueryBuilder
	pb *proto.ContractGetRecordsQuery
}

// NewContractRecordsQuery creates a ContractRecordsQuery builder which can be used to construct and execute a
// Contract Get Records Query
func NewContractRecordsQuery() *ContractRecordsQuery {
	pb := &proto.ContractGetRecordsQuery{Header: &proto.QueryHeader{}}

	inner := newQueryBuilder(pb.Header)
	inner.pb.Query = &proto.Query_ContractGetRecords{ContractGetRecords: pb}

	return &ContractRecordsQuery{inner, pb}
}

// SetContractID sets the smart contract instance for which the records should be retrieved
func (builder *ContractRecordsQuery) SetContractID(id ContractID) *ContractRecordsQuery {
	builder.pb.ContractID = id.toProto()
	return builder
}

// Execute executes the ContractRecordsQuery using the provided client.
func (builder *ContractRecordsQuery) Execute(client *Client) ([]TransactionRecord, error) {
	resp, err := builder.execute(client)
	if err != nil {
		return nil, err
	}

	rawRecords := resp.GetContractGetRecordsResponse().Records
	records := make([]TransactionRecord, len(rawRecords))

	for i, element := range resp.GetContractGetRecordsResponse().Records {
		records[i] = transactionRecordFromProto(element)
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
