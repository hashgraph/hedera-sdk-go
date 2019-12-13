package hedera

import "github.com/hashgraph/hedera-sdk-go/proto"

type TransactionReceiptQuery struct {
	QueryBuilder
	pb *proto.TransactionGetReceiptQuery
}

func NewTransactionReceiptQuery() *TransactionReceiptQuery {
	pb := &proto.TransactionGetReceiptQuery{Header: &proto.QueryHeader{}}

	inner := newQueryBuilder(pb.Header)
	inner.pb.Query = &proto.Query_TransactionGetReceipt{pb}

	return &TransactionReceiptQuery{inner, pb}
}

func (builder *TransactionReceiptQuery) SetTransactionID(id TransactionID) *TransactionReceiptQuery {
	builder.pb.TransactionID = id.toProto()
	return builder
}

func (builder *TransactionReceiptQuery) Execute(client *Client) (TransactionReceipt, error) {
	resp, err := builder.execute(client)
	if err != nil {
		return TransactionReceipt{}, err
	}

	return transactionReceiptFromResponse(resp), nil
}
