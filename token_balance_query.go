package hedera

import "github.com/hashgraph/hedera-sdk-go/proto"

// FileContentsQuery retrieves the contents of a file.
type TokenBalanceQuery struct {
	QueryBuilder
	pb *proto.CryptoGetAccountBalanceQuery
}

// NewFileContentsQuery creates a FileContentsQuery builder which can be used to construct and execute a
// File Get Contents Query.
func NewTokenBalanceQuery() *TokenBalanceQuery {
	pb := &proto.CryptoGetAccountBalanceQuery{}

	inner := newQueryBuilder(pb.Header)
	inner.pb.Query = &proto.Query_CryptogetAccountBalance{CryptogetAccountBalance: pb}

	return &TokenBalanceQuery{inner, pb}
}

func (builder *TokenBalanceQuery) SetAccountID(id AccountID) *TokenBalanceQuery {
	builder.pb.BalanceSource = &proto.CryptoGetAccountBalanceQuery_AccountID{
		AccountID: id.toProto(),
	}

	return builder
}

func (builder *TokenBalanceQuery) SetContractID(id ContractID) *TokenBalanceQuery {
	builder.pb.BalanceSource = &proto.CryptoGetAccountBalanceQuery_ContractID{
		ContractID: id.toProto(),
	}

	return builder
}

// Execute executes the AccountBalanceQuery using the provided client
func (builder *TokenBalanceQuery) Execute(client *Client) (map[TokenID]uint64, error) {
	resp, err := builder.execute(client)
	if err != nil {
		return make(map[TokenID]uint64), err
	}

	tokenBalances := make(map[TokenID]uint64, len(resp.GetCryptogetAccountBalance().TokenBalances))
	for _, token := range resp.GetCryptogetAccountBalance().TokenBalances {
		tokenBalances[tokenIDFromProto(token.TokenId)] = token.Balance
	}

	return tokenBalances, err
}

// SetMaxQueryPayment sets the maximum payment allowed for this Query.
func (builder *TokenBalanceQuery) SetMaxQueryPayment(maxPayment Hbar) *TokenBalanceQuery {
	return &TokenBalanceQuery{*builder.QueryBuilder.SetMaxQueryPayment(maxPayment), builder.pb}
}

// SetQueryPayment sets the payment amount for this Query.
func (builder *TokenBalanceQuery) SetQueryPayment(paymentAmount Hbar) *TokenBalanceQuery {
	return &TokenBalanceQuery{*builder.QueryBuilder.SetQueryPayment(paymentAmount), builder.pb}
}

// SetQueryPaymentTransaction sets the payment Transaction for this Query.
func (builder *TokenBalanceQuery) SetQueryPaymentTransaction(tx Transaction) *TokenBalanceQuery {
	return &TokenBalanceQuery{*builder.QueryBuilder.SetQueryPaymentTransaction(tx), builder.pb}
}


