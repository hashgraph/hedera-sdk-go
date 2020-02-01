package hedera

import (
	"github.com/hashgraph/hedera-sdk-go/proto"
)

type TransactionRecordQuery struct {
	QueryBuilder
	pb *proto.TransactionGetRecordQuery
}

func NewTransactionRecordQuery() *TransactionRecordQuery {
	pb := &proto.TransactionGetRecordQuery{Header: &proto.QueryHeader{}}

	inner := newQueryBuilder(pb.Header)
	inner.pb.Query = &proto.Query_TransactionGetRecord{TransactionGetRecord: pb}

	return &TransactionRecordQuery{inner, pb}
}

func (builder *TransactionRecordQuery) SetTransactionID(id TransactionID) *TransactionRecordQuery {
	builder.pb.TransactionID = id.toProto()
	return builder
}

func (builder *TransactionRecordQuery) Execute(client *Client) (TransactionRecord, error) {
	resp, err := builder.execute(client)
	if err != nil {
		return TransactionRecord{}, err
	}

	return transactionRecordFromProto(resp.GetTransactionGetRecord().TransactionRecord), nil
}

//
// The following _3_ must be copy-pasted at the bottom of **every** _query.go file
// We override the embedded fluent setter methods to return the outer type
//

func (builder *TransactionRecordQuery) SetMaxQueryPayment(maxPayment Hbar) *TransactionRecordQuery {
	return &TransactionRecordQuery{*builder.QueryBuilder.SetMaxQueryPayment(maxPayment), builder.pb}
}

func (builder *TransactionRecordQuery) SetQueryPayment(paymentAmount Hbar) *TransactionRecordQuery {
	return &TransactionRecordQuery{*builder.QueryBuilder.SetQueryPayment(paymentAmount), builder.pb}
}

func (builder *TransactionRecordQuery) SetQueryPaymentTransaction(tx Transaction) *TransactionRecordQuery {
	return &TransactionRecordQuery{*builder.QueryBuilder.SetQueryPaymentTransaction(tx), builder.pb}
}
