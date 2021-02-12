package hedera

import (
	"github.com/hashgraph/hedera-sdk-go/proto"
)

// TokenInfoQuery gets all the information about an token excluding token records.
// This includes the  balance.
type TokenInfoQuery struct {
	QueryBuilder
	pb *proto.TokenGetInfoQuery
}

// TokenInfo is info about the token returned from an TokenInfoQuery
type TokenInfo struct {
	TokenID           TokenID
	Name              string
	Symbol            string
	Decimals          uint32
	TotalSupply       uint64
	Treasury          AccountID
	AdminKey          PublicKey
	KycKey            PublicKey
	WipeKey           PublicKey
	FreezeKey         PublicKey
	SupplyKey         PublicKey
	TokenFreezeStatus *bool
	TokenKycStatus    *bool
	IsDeleted         bool
	AutoRenewAccount  AccountID
	AutoRenewPeriod   uint64
	ExpirationTime    uint64
}

// NewTokenInfoQuery creates an TokenInfoQuery builder which can be used to construct and execute
// an TokenInfoQuery.
//
// It is recommended that you use this for creating new instances of an TokenInfoQuery
// instead of manually creating an instance of the struct.
func NewTokenInfoQuery() *TokenInfoQuery {
	pb := &proto.TokenGetInfoQuery{Header: &proto.QueryHeader{}}

	inner := newQueryBuilder(pb.Header)
	inner.pb.Query = &proto.Query_TokenGetInfo{TokenGetInfo: pb}

	return &TokenInfoQuery{inner, pb}
}

// SetTokenID sets the token ID for which information is requested
func (builder *TokenInfoQuery) SetTokenID(id TokenID) *TokenInfoQuery {
	builder.pb.Token = id.toProto()
	return builder
}

// Execute executes the TokenInfoQuery using the provided client
func (builder *TokenInfoQuery) Execute(client *Client) (TokenInfo, error) {
	resp, err := builder.execute(client)
	if err != nil {
		return TokenInfo{}, err
	}

	adminKey, err := publicKeyFromProto(resp.GetTokenGetInfo().TokenInfo.AdminKey)
	if err != nil {
		return TokenInfo{}, err
	}

	kycKey, err := publicKeyFromProto(resp.GetTokenGetInfo().TokenInfo.KycKey)
	if err != nil {
		return TokenInfo{}, err
	}

	wipeKey, err := publicKeyFromProto(resp.GetTokenGetInfo().TokenInfo.WipeKey)
	if err != nil {
		return TokenInfo{}, err
	}

	freezeKey, err := publicKeyFromProto(resp.GetTokenGetInfo().TokenInfo.FreezeKey)
	if err != nil {
		return TokenInfo{}, err
	}

	supplyKey, err := publicKeyFromProto(resp.GetTokenGetInfo().TokenInfo.SupplyKey)
	if err != nil {
		return TokenInfo{}, err
	}

	var kycStatus *bool = nil
	if resp.GetTokenGetInfo().TokenInfo.DefaultKycStatus == 1 {
		status := false
		kycStatus = &status
	} else if resp.GetTokenGetInfo().TokenInfo.DefaultKycStatus == 2 {
		status := true
		kycStatus = &status
	}

	var freezeStatus *bool = nil
	if resp.GetTokenGetInfo().TokenInfo.DefaultFreezeStatus == 1 {
		status := false
		freezeStatus = &status
	} else if resp.GetTokenGetInfo().TokenInfo.DefaultFreezeStatus == 2 {
		status := true
		freezeStatus = &status
	}

	return TokenInfo{
		TokenID:           tokenIDFromProto(resp.GetTokenGetInfo().TokenInfo.TokenId),
		Name:              resp.GetTokenGetInfo().TokenInfo.Name,
		Symbol:            resp.GetTokenGetInfo().TokenInfo.Symbol,
		Decimals:          resp.GetTokenGetInfo().TokenInfo.Decimals,
		TotalSupply:       resp.GetTokenGetInfo().TokenInfo.TotalSupply,
		Treasury:          accountIDFromProto(resp.GetTokenGetInfo().TokenInfo.Treasury),
		AdminKey:          adminKey,
		KycKey:            kycKey,
		WipeKey:           wipeKey,
		FreezeKey:         freezeKey,
		SupplyKey:         supplyKey,
		TokenFreezeStatus: freezeStatus,
		TokenKycStatus:    kycStatus,
		IsDeleted:         resp.GetTokenGetInfo().TokenInfo.Deleted,
		AutoRenewAccount:  accountIDFromProto(resp.GetTokenGetInfo().TokenInfo.AutoRenewAccount),
		AutoRenewPeriod:   uint64(resp.GetTokenGetInfo().TokenInfo.AutoRenewPeriod.Seconds),
		ExpirationTime:    uint64(resp.GetTokenGetInfo().TokenInfo.Expiry.Seconds),
	}, nil
}

// Cost is a wrapper around the standard Cost function for a query. It must exist because the cost returned by the
// standard Cost() and the Hedera Network doesn't work for any accounnts that have been deleted. In that case the
// minimum cost should be ~25 Tinybar which seems to succeed most of the time.
func (builder *TokenInfoQuery) Cost(client *Client) (Hbar, error) {
	// deleted files return a COST_ANSWER of zero which triggers `INSUFFICIENT_TX_FEE`
	// if you set that as the query payment; 25 tinybar seems to be enough to get
	// `TOKEN_DELETED` back instead.
	cost, err := builder.QueryBuilder.GetCost(client)
	if err != nil {
		return ZeroHbar, err
	}

	// math.Max requires float64 and returns float64
	if cost.AsTinybar() > 25 {
		return cost, nil
	}

	return HbarFromTinybar(25), nil
}

//
// The following _3_ must be copy-pasted at the bottom of **every** _query.go file
// We override the embedded fluent setter methods to return the outer type
//

// SetMaxQueryPayment sets the maximum payment allowed for this Query.
func (builder *TokenInfoQuery) SetMaxQueryPayment(maxPayment Hbar) *TokenInfoQuery {
	return &TokenInfoQuery{*builder.QueryBuilder.SetMaxQueryPayment(maxPayment), builder.pb}
}

// SetQueryPayment sets the payment amount for this Query.
func (builder *TokenInfoQuery) SetQueryPayment(paymentAmount Hbar) *TokenInfoQuery {
	return &TokenInfoQuery{*builder.QueryBuilder.SetQueryPayment(paymentAmount), builder.pb}
}

// SetQueryPaymentTransaction sets the payment Transaction for this Query.
func (builder *TokenInfoQuery) SetQueryPaymentTransaction(tx Transaction) *TokenInfoQuery {
	return &TokenInfoQuery{*builder.QueryBuilder.SetQueryPaymentTransaction(tx), builder.pb}
}
