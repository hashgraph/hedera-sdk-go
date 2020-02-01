package hedera

import "github.com/hashgraph/hedera-sdk-go/proto"

type AccountBalanceQuery struct {
	QueryBuilder
	pb *proto.CryptoGetAccountBalanceQuery
}

func NewAccountBalanceQuery() *AccountBalanceQuery {
	pb := &proto.CryptoGetAccountBalanceQuery{Header: &proto.QueryHeader{}}

	inner := newQueryBuilder(pb.Header)
	inner.pb.Query = &proto.Query_CryptogetAccountBalance{CryptogetAccountBalance: pb}

	return &AccountBalanceQuery{inner, pb}
}

func (builder *AccountBalanceQuery) SetAccountID(id AccountID) *AccountBalanceQuery {
	builder.pb.AccountID = id.toProto()
	return builder
}

func (builder *AccountBalanceQuery) Execute(client *Client) (Hbar, error) {
	resp, err := builder.execute(client)
	if err != nil {
		return ZeroHbar, err
	}

	return HbarFromTinybar(int64(resp.GetCryptogetAccountBalance().Balance)), nil
}

//
// The following _3_ must be copy-pasted at the bottom of **every** _query.go file
// We override the embedded fluent setter methods to return the outer type
//

func (builder *AccountBalanceQuery) SetMaxQueryPayment(maxPayment Hbar) *AccountBalanceQuery {
	return &AccountBalanceQuery{*builder.QueryBuilder.SetMaxQueryPayment(maxPayment), builder.pb}
}

func (builder *AccountBalanceQuery) SetQueryPayment(paymentAmount Hbar) *AccountBalanceQuery {
	return &AccountBalanceQuery{*builder.QueryBuilder.SetQueryPayment(paymentAmount), builder.pb}
}

func (builder *AccountBalanceQuery) SetQueryPaymentTransaction(tx Transaction) *AccountBalanceQuery {
	return &AccountBalanceQuery{*builder.QueryBuilder.SetQueryPaymentTransaction(tx), builder.pb}
}
