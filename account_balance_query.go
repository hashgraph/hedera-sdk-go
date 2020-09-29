package hedera

import "github.com/hashgraph/hedera-sdk-go/proto"

// AccountBalanceQuery gets the balance of a CryptoCurrency account. This returns only the balance, so it is a smaller
// and faster reply than AccountInfoQuery, which returns the balance plus additional information.
type AccountBalanceQuery struct {
	QueryBuilder
	pb *proto.CryptoGetAccountBalanceQuery
}

// NewAccountBalanceQuery creates an AccountBalanceQuery transaction which can be used to construct and execute
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
func (transaction *AccountBalanceQuery) SetAccountID(id AccountID) *AccountBalanceQuery {
	transaction.pb.BalanceSource = &proto.CryptoGetAccountBalanceQuery_AccountID{
		AccountID: id.toProto(),
	}

	return transaction
}

// SetContractID sets the ContractID for which you wish to query the balance.
//
// Note: you can only query an Account or Contract but not both -- if a Contract ID or Account ID has already been set,
// it will be overwritten by this method.
func (transaction *AccountBalanceQuery) SetContractID(id ContractID) *AccountBalanceQuery {
	transaction.pb.BalanceSource = &proto.CryptoGetAccountBalanceQuery_ContractID{
		ContractID: id.toProto(),
	}

	return transaction
}

// Execute executes the AccountBalanceQuery using the provided client
func (transaction *AccountBalanceQuery) Execute(client *Client) (Hbar, error) {
	resp, err := transaction.execute(client)
	if err != nil {
		return ZeroHbar, err
	}

	return HbarFromTinybar(int64(resp.GetCryptogetAccountBalance().Balance)), nil
}

//
// The following _3_ must be copy-pasted at the bottom of **every** _query.go file
// We override the embedded fluent setter methods to return the outer type
//

// SetMaxQueryPayment sets the maximum payment allowed for this Query.
func (transaction *AccountBalanceQuery) SetMaxQueryPayment(maxPayment Hbar) *AccountBalanceQuery {
	return &AccountBalanceQuery{*transaction.QueryBuilder.SetMaxQueryPayment(maxPayment), transaction.pb}
}

// SetQueryPayment sets the payment amount for this Query.
func (transaction *AccountBalanceQuery) SetQueryPayment(paymentAmount Hbar) *AccountBalanceQuery {
	return &AccountBalanceQuery{*transaction.QueryBuilder.SetQueryPayment(paymentAmount), transaction.pb}
}

// SetQueryPaymentTransaction sets the payment Transaction for this Query.
func (transaction *AccountBalanceQuery) SetQueryPaymentTransaction(tx Transaction) *AccountBalanceQuery {
	return &AccountBalanceQuery{*transaction.QueryBuilder.SetQueryPaymentTransaction(tx), transaction.pb}
}
