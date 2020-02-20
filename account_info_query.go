package hedera

import (
	"github.com/hashgraph/hedera-sdk-go/proto"
	"time"
)

// AccountInfoQuery gets all the information about an account excluding account records.
// This includes the  balance.
type AccountInfoQuery struct {
	QueryBuilder
	pb *proto.CryptoGetInfoQuery
}

// AccountInfo is info about the account returned from an AccountInfoQuery
type AccountInfo struct {
	AccountID                      AccountID
	ContractAccountID              string
	Deleted                        bool
	ProxyAccountID                 AccountID
	ProxyReceived                  Hbar
	Key                            PublicKey
	Balance                        Hbar
	GenerateSendRecordThreshold    Hbar
	GenerateReceiveRecordThreshold Hbar
	ReceiverSigRequired            bool
	ExpirationTime                 time.Time
	AutoRenewPeriod                time.Duration
}

// NewAccountInfoQuery creates an AccountInfoQuery builder which can be used to construct and execute
// an AccountInfoQuery.
//
// It is recommended that you use this for creating new instances of an AccountInfoQuery
// instead of manually creating an instance of the struct.
func NewAccountInfoQuery() *AccountInfoQuery {
	pb := &proto.CryptoGetInfoQuery{Header: &proto.QueryHeader{}}

	inner := newQueryBuilder(pb.Header)
	inner.pb.Query = &proto.Query_CryptoGetInfo{CryptoGetInfo: pb}

	return &AccountInfoQuery{inner, pb}
}

// SetAccountID sets the account ID for which information is requested
func (builder *AccountInfoQuery) SetAccountID(id AccountID) *AccountInfoQuery {
	builder.pb.AccountID = id.toProto()
	return builder
}

// Execute executes the AccountInfoQuery using the provided client
func (builder *AccountInfoQuery) Execute(client *Client) (AccountInfo, error) {
	resp, err := builder.execute(client)
	if err != nil {
		return AccountInfo{}, err
	}

	pubKey, err := publicKeyFromProto(resp.GetCryptoGetInfo().AccountInfo.Key)
	if err != nil {
		return AccountInfo{}, err
	}

	return AccountInfo{
		AccountID:                      accountIDFromProto(resp.GetCryptoGetInfo().AccountInfo.AccountID),
		ContractAccountID:              resp.GetCryptoGetInfo().AccountInfo.ContractAccountID,
		Deleted:                        resp.GetCryptoGetInfo().AccountInfo.Deleted,
		ProxyAccountID:                 accountIDFromProto(resp.GetCryptoGetInfo().AccountInfo.ProxyAccountID),
		ProxyReceived:                  HbarFromTinybar(resp.GetCryptoGetInfo().AccountInfo.ProxyReceived),
		Key:                            pubKey,
		Balance:                        HbarFromTinybar(int64(resp.GetCryptoGetInfo().AccountInfo.Balance)),
		GenerateSendRecordThreshold:    HbarFromTinybar(int64(resp.GetCryptoGetInfo().AccountInfo.GenerateSendRecordThreshold)),
		GenerateReceiveRecordThreshold: HbarFromTinybar(int64(resp.GetCryptoGetInfo().AccountInfo.GenerateReceiveRecordThreshold)),
		ReceiverSigRequired:            resp.GetCryptoGetInfo().AccountInfo.ReceiverSigRequired,
		ExpirationTime:                 timeFromProto(resp.GetCryptoGetInfo().AccountInfo.ExpirationTime),
	}, nil
}

// Cost is a wrapper around the standard Cost function for a query. It must exist because the cost returned by the
// standard Cost() and the Hedera Network doesn't work for any accounnts that have been deleted. In that case the
// minimum cost should be ~25 Tinybar which seems to succeed most of the time.
func (builder *AccountInfoQuery) Cost(client *Client) (Hbar, error) {
	// deleted files return a COST_ANSWER of zero which triggers `INSUFFICIENT_TX_FEE`
	// if you set that as the query payment; 25 tinybar seems to be enough to get
	// `ACCOUNT_DELETED` back instead.
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

func (builder *AccountInfoQuery) SetMaxQueryPayment(maxPayment Hbar) *AccountInfoQuery {
	return &AccountInfoQuery{*builder.QueryBuilder.SetMaxQueryPayment(maxPayment), builder.pb}
}

func (builder *AccountInfoQuery) SetQueryPayment(paymentAmount Hbar) *AccountInfoQuery {
	return &AccountInfoQuery{*builder.QueryBuilder.SetQueryPayment(paymentAmount), builder.pb}
}

func (builder *AccountInfoQuery) SetQueryPaymentTransaction(tx Transaction) *AccountInfoQuery {
	return &AccountInfoQuery{*builder.QueryBuilder.SetQueryPaymentTransaction(tx), builder.pb}
}
