package hedera

import "github.com/hashgraph/hedera-sdk-go/proto"

type AccountStakersQuery struct {
	QueryBuilder
	pb *proto.CryptoGetStakersQuery
}

type ProxyStaker struct {
	AccountID AccountID
	Amount    int64
}

func NewAccountStakersQuery() *AccountStakersQuery {
	pb := &proto.CryptoGetStakersQuery{Header: &proto.QueryHeader{}}

	inner := newQueryBuilder(pb.Header)
	inner.pb.Query = &proto.Query_CryptoGetProxyStakers{pb}

	return &AccountStakersQuery{inner, pb}
}

func (builder *AccountStakersQuery) SetAccountID(id AccountID) *AccountStakersQuery {
	builder.pb.AccountID = id.toProto()
	return builder
}

func (builder *AccountStakersQuery) Execute(client *Client) ([]ProxyStaker, error) {
	var stakers = []ProxyStaker{}

	resp, err := builder.execute(client)
	if err != nil {
		return stakers, err
	}

	for _, element := range resp.GetCryptoGetProxyStakers().Stakers.ProxyStaker {
		stakers = append(stakers, ProxyStaker{
			AccountID: accountIDFromProto(element.AccountID),
			Amount:    element.Amount,
		})
	}

	return stakers, nil
}
