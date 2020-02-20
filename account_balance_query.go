package hedera

import "github.com/hashgraph/hedera-sdk-go/proto"

// AccountBalanceQuery gets the balance of a CryptoCurrency account. This returns only the balance, so it is a smaller
// and faster reply than AccountInfoQuery, which returns the balance plus additional information.
type AccountBalanceQuery struct {
	QueryBuilder
	pb *proto.CryptoGetAccountBalanceQuery
}

// NewAccountBalanceQuery creates an AccountBalanceQuery builder which can be used to construct and execute
// an AccountBalanceQuery.
// It is recommended that you use this for creating new instances of an AccountBalanceQuery
// instead of manually creating an instance of the struct.
func NewAccountBalanceQuery() *AccountBalanceQuery {
	pb := &proto.CryptoGetAccountBalanceQuery{Header: &proto.QueryHeader{}}

	inner := newQueryBuilder(pb.Header)
	inner.pb.Query = &proto.Query_CryptogetAccountBalance{CryptogetAccountBalance: pb}

	return &AccountBalanceQuery{inner, pb}
}

// SetAccountID sets the AccountID for which you wish to query the balance.
//
// Note: you can only query an Account or Contract but not both -- if a Contract ID or Account ID has already been set,
// it will be overwritten by this method.
func (builder *AccountBalanceQuery) SetAccountID(id AccountID) *AccountBalanceQuery {
	builder.pb.BalanceSource = &proto.CryptoGetAccountBalanceQuery_AccountID{
		AccountID: id.toProto(),
	}

	return builder
}

// SetContractID sets the ContractID for which you wish to query the balance.
//
// Note: you can only query an Account or Contract but not both -- if a Contract ID or Account ID has already been set,
// it will be overwritten by this method.
func (builder *AccountBalanceQuery) SetContractID(id ContractID) *AccountBalanceQuery {
	builder.pb.BalanceSource = &proto.CryptoGetAccountBalanceQuery_ContractID{
		ContractID: id.toProto(),
	}

	return builder
}

// Execute executes the AccountBalanceQuery using the provided client
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
