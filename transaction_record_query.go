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

// SetTransactionID sets the TransactionID for which to request the TransactionRecord.
func (transaction *TransactionRecordQuery) SetTransactionID(id TransactionID) *TransactionRecordQuery {
	transaction.pb.TransactionID = id.toProto()
	return transaction
}

func (transaction *TransactionRecordQuery) Execute(client *Client) (TransactionRecord, error) {
	resp, err := transaction.execute(client)
	if err != nil {
		return TransactionRecord{}, err
	}

	return transactionRecordFromProto(resp.GetTransactionGetRecord().TransactionRecord), nil
}

//
// The following _3_ must be copy-pasted at the bottom of **every** _query.go file
// We override the embedded fluent setter methods to return the outer type
//

// SetMaxQueryPayment sets the maximum payment allowed for this Query.
func (transaction *TransactionRecordQuery) SetMaxQueryPayment(maxPayment Hbar) *TransactionRecordQuery {
	return &TransactionRecordQuery{*transaction.QueryBuilder.SetMaxQueryPayment(maxPayment), transaction.pb}
}

// SetQueryPayment sets the payment amount for this Query.
func (transaction *TransactionRecordQuery) SetQueryPayment(paymentAmount Hbar) *TransactionRecordQuery {
	return &TransactionRecordQuery{*transaction.QueryBuilder.SetQueryPayment(paymentAmount), transaction.pb}
}

// SetQueryPaymentTransaction sets the payment Transaction for this Query.
func (transaction *TransactionRecordQuery) SetQueryPaymentTransaction(tx Transaction) *TransactionRecordQuery {
	return &TransactionRecordQuery{*transaction.QueryBuilder.SetQueryPaymentTransaction(tx), transaction.pb}
}
