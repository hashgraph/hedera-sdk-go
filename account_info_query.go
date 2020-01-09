package hedera

import (
	"github.com/hashgraph/hedera-sdk-go/proto"
	"time"
)

type AccountInfoQuery struct {
	QueryBuilder
	pb *proto.CryptoGetInfoQuery
}

type AccountInfo struct {
	AccountID                      AccountID
	ContractAccountID              string
	Deleted                        bool
	ProxyAccountID                 AccountID
	ProxyReceived                  int64
	Key                            Ed25519PublicKey
	Balance                        uint64
	GenerateSendRecordThreshold    uint64
	GenerateReceiveRecordThreshold uint64
	ReceiverSigRequired            bool
	ExpirationTime                 time.Time
	AutoRenewPeriod                time.Duration
}

func NewAccountInfoQuery() *AccountInfoQuery {
	pb := &proto.CryptoGetInfoQuery{Header: &proto.QueryHeader{}}

	inner := newQueryBuilder(pb.Header)
	inner.pb.Query = &proto.Query_CryptoGetInfo{pb}

	return &AccountInfoQuery{inner, pb}
}

func (builder *AccountInfoQuery) SetAccountID(id AccountID) *AccountInfoQuery {
	builder.pb.AccountID = id.toProto()
	return builder
}

func (builder *AccountInfoQuery) Execute(client *Client) (AccountInfo, error) {
	resp, err := builder.execute(client)
	if err != nil {
		return AccountInfo{}, err
	}

	return AccountInfo{
		AccountID:                      accountIDFromProto(resp.GetCryptoGetInfo().AccountInfo.AccountID),
		ContractAccountID:              resp.GetCryptoGetInfo().AccountInfo.ContractAccountID,
		Deleted:                        resp.GetCryptoGetInfo().AccountInfo.Deleted,
		ProxyAccountID:                 accountIDFromProto(resp.GetCryptoGetInfo().AccountInfo.ProxyAccountID),
		ProxyReceived:                  resp.GetCryptoGetInfo().AccountInfo.ProxyReceived,
		Key:                            Ed25519PublicKey{keyData: resp.GetCryptoGetInfo().AccountInfo.Key.GetEd25519()},
		Balance:                        resp.GetCryptoGetInfo().AccountInfo.Balance,
		GenerateSendRecordThreshold:    resp.GetCryptoGetInfo().AccountInfo.GenerateSendRecordThreshold,
		GenerateReceiveRecordThreshold: resp.GetCryptoGetInfo().AccountInfo.GenerateReceiveRecordThreshold,
		ReceiverSigRequired:            resp.GetCryptoGetInfo().AccountInfo.ReceiverSigRequired,
		ExpirationTime:                 timeFromProto(resp.GetCryptoGetInfo().AccountInfo.ExpirationTime),
	}, nil
}

func (builder *AccountInfoQuery) Cost(client *Client) (uint64, error) {
	// deleted files return a COST_ANSWER of zero which triggers `INSUFFICIENT_TX_FEE`
	// if you set that as the query payment; 25 tinybar seems to be enough to get
	// `ACCOUNT_DELETED` back instead.
	cost, err := builder.QueryBuilder.Cost(client)
	if err != nil {
		return 0, err
	}

	// math.Min requires float64 and returns float64
	if cost > 25 {
		return cost, nil
	} else {
		return 25, nil
	}

}
