package hedera

import "github.com/hashgraph/hedera-sdk-go/proto"

type TransactionReceiptQuery struct {
	QueryBuilder
	pb *proto.TransactionGetReceiptQuery
}

func NewTransactionReceiptQuery() *TransactionReceiptQuery {
	pb := &proto.TransactionGetReceiptQuery{Header: &proto.QueryHeader{}}

	inner := newQueryBuilder(pb.Header)
	inner.pb.Query = &proto.Query_TransactionGetReceipt{TransactionGetReceipt: pb}

	return &TransactionReceiptQuery{inner, pb}
}

// SetTransactionID sets the TransactionID for which to request the TransactionReceipt.
func (transaction *TransactionReceiptQuery) SetTransactionID(id TransactionID) *TransactionReceiptQuery {
	transaction.pb.TransactionID = id.toProto()
	return transaction
}

func (transaction *TransactionReceiptQuery) Execute(client *Client) (TransactionReceipt, error) {
	resp, err := transaction.execute(client)
	if err != nil {
		return TransactionReceipt{}, err
	}

	return transactionReceiptFromResponse(resp), nil
}

//
// The following _3_ must be copy-pasted at the bottom of **every** _query.go file
// We override the embedded fluent setter methods to return the outer type
//

// SetMaxQueryPayment sets the maximum payment allowed for this Query.
func (transaction *TransactionReceiptQuery) SetMaxQueryPayment(maxPayment Hbar) *TransactionReceiptQuery {
	return &TransactionReceiptQuery{*transaction.QueryBuilder.SetMaxQueryPayment(maxPayment), transaction.pb}
}

// SetQueryPayment sets the payment amount for this Query.
func (transaction *TransactionReceiptQuery) SetQueryPayment(paymentAmount Hbar) *TransactionReceiptQuery {
	return &TransactionReceiptQuery{*transaction.QueryBuilder.SetQueryPayment(paymentAmount), transaction.pb}
}

// SetQueryPaymentTransaction sets the payment Transaction for this Query.
func (transaction *TransactionReceiptQuery) SetQueryPaymentTransaction(tx Transaction) *TransactionReceiptQuery {
	return &TransactionReceiptQuery{*transaction.QueryBuilder.SetQueryPaymentTransaction(tx), transaction.pb}
}
