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
