package hedera

import "github.com/hashgraph/hedera-sdk-go/proto"

type AccountRecordsQuery struct {
	QueryBuilder
	pb *proto.CryptoGetAccountRecordsQuery
}

func NewAccountRecordsQuery() *AccountRecordsQuery {
	pb := &proto.CryptoGetAccountRecordsQuery{Header: &proto.QueryHeader{}}

	inner := newQueryBuilder(pb.Header)
	inner.pb.Query = &proto.Query_CryptoGetAccountRecords{CryptoGetAccountRecords: pb}

	return &AccountRecordsQuery{inner, pb}
}

func (builder *AccountRecordsQuery) SetAccountID(id AccountID) *AccountRecordsQuery {
	builder.pb.AccountID = id.toProto()
	return builder
}

func (builder *AccountRecordsQuery) Execute(client *Client) ([]TransactionRecord, error) {
	var records = []TransactionRecord{}

	resp, err := builder.execute(client)
	if err != nil {
		return records, err
	}

	for _, element := range resp.GetCryptoGetAccountRecords().Records {
		records = append(records, transactionRecordFromProto(element))
	}

	return records, nil
}

//
// The following _3_ must be copy-pasted at the bottom of **every** _query.go file
// We override the embedded fluent setter methods to return the outer type
//

func (builder *AccountRecordsQuery) SetMaxQueryPayment(maxPayment Hbar) *AccountRecordsQuery {
	return &AccountRecordsQuery{*builder.QueryBuilder.SetMaxQueryPayment(maxPayment), builder.pb}
}

func (builder *AccountRecordsQuery) SetQueryPayment(paymentAmount Hbar) *AccountRecordsQuery {
	return &AccountRecordsQuery{*builder.QueryBuilder.SetQueryPayment(paymentAmount), builder.pb}
}

func (builder *AccountRecordsQuery) SetQueryPaymentTransaction(tx Transaction) *AccountRecordsQuery {
	return &AccountRecordsQuery{*builder.QueryBuilder.SetQueryPaymentTransaction(tx), builder.pb}
}
