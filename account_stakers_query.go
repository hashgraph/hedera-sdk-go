package hedera

import "github.com/hashgraph/hedera-sdk-go/proto"

// AccountStakersQuery gets all of the accounts that are proxy staking to this account. For each of  them, the amount
// currently staked will be given. This is not yet implemented, but will be in a future version of the API.
type AccountStakersQuery struct {
	QueryBuilder
	pb *proto.CryptoGetStakersQuery
}

// NewAccountStakersQuery creates an AccountStakersQuery builder which can be used to construct and execute
// an AccountStakersQuery.
//
// It is recommended that you use this for creating new instances of an AccountStakersQuery
// instead of manually creating an instance of the struct.
func NewAccountStakersQuery() *AccountStakersQuery {
	pb := &proto.CryptoGetStakersQuery{Header: &proto.QueryHeader{}}

	inner := newQueryBuilder(pb.Header)
	inner.pb.Query = &proto.Query_CryptoGetProxyStakers{CryptoGetProxyStakers: pb}

	return &AccountStakersQuery{inner, pb}
}

// SetAccountID sets the Account ID for which the stakers should be retrieved
func (builder *AccountStakersQuery) SetAccountID(id AccountID) *AccountStakersQuery {
	builder.pb.AccountID = id.toProto()
	return builder
}

// Execute executes the AccountStakersQuery using the provided client.
func (builder *AccountStakersQuery) Execute(client *Client) ([]Transfer, error) {
	resp, err := builder.execute(client)
	if err != nil {
		return []Transfer{}, err
	}

	var stakers = make([]Transfer, len(resp.GetCryptoGetProxyStakers().Stakers.ProxyStaker))

	for i, element := range resp.GetCryptoGetProxyStakers().Stakers.ProxyStaker {
		stakers[i] = Transfer{
			AccountID: accountIDFromProto(element.AccountID),
			Amount:    HbarFromTinybar(element.Amount),
		}
	}

	return stakers, nil
}

//
// The following _3_ must be copy-pasted at the bottom of **every** _query.go file
// We override the embedded fluent setter methods to return the outer type
//

func (builder *AccountStakersQuery) SetMaxQueryPayment(maxPayment Hbar) *AccountStakersQuery {
	return &AccountStakersQuery{*builder.QueryBuilder.SetMaxQueryPayment(maxPayment), builder.pb}
}

func (builder *AccountStakersQuery) SetQueryPayment(paymentAmount Hbar) *AccountStakersQuery {
	return &AccountStakersQuery{*builder.QueryBuilder.SetQueryPayment(paymentAmount), builder.pb}
}

func (builder *AccountStakersQuery) SetQueryPaymentTransaction(tx Transaction) *AccountStakersQuery {
	return &AccountStakersQuery{*builder.QueryBuilder.SetQueryPaymentTransaction(tx), builder.pb}
}
