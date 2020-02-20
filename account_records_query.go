package hedera

import "github.com/hashgraph/hedera-sdk-go/proto"

// AccountRecordsQuery gets all of the records for an account for any transfers into it and out of
// it, that were above the threshold, during the last 25 hours.
type AccountRecordsQuery struct {
	QueryBuilder
	pb *proto.CryptoGetAccountRecordsQuery
}

// NewAccountRecordsQuery creates an AccountRecordsQuery builder which can be used to construct and execute
// an AccountRecordsQuery.
//
// It is recommended that you use this for creating new instances of an AccountRecordQuery
// instead of manually creating an instance of the struct.
func NewAccountRecordsQuery() *AccountRecordsQuery {
	pb := &proto.CryptoGetAccountRecordsQuery{Header: &proto.QueryHeader{}}

	inner := newQueryBuilder(pb.Header)
	inner.pb.Query = &proto.Query_CryptoGetAccountRecords{CryptoGetAccountRecords: pb}

	return &AccountRecordsQuery{inner, pb}
}

// SetAccountID sets the account ID for which the records should be retrieved.
func (builder *AccountRecordsQuery) SetAccountID(id AccountID) *AccountRecordsQuery {
	builder.pb.AccountID = id.toProto()
	return builder
}

// Execute executes the AccountRecordsQuery using the provided client
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
