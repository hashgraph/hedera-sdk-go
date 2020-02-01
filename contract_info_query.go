package hedera

import (
	"github.com/hashgraph/hedera-sdk-go/proto"
	"time"
)

type ContractInfoQuery struct {
	QueryBuilder
	pb *proto.ContractGetInfoQuery
}

type ContractInfo struct {
	AccountID         AccountID
	ContractID        ContractID
	ContractAccountID string
	AdminKey          PublicKey
	ExpirationTime    time.Time
	AutoRenewPeriod   time.Duration
	Storage           uint64
	ContractMemo      string
}

func NewContractInfoQuery() *ContractInfoQuery {
	pb := &proto.ContractGetInfoQuery{Header: &proto.QueryHeader{}}

	inner := newQueryBuilder(pb.Header)
	inner.pb.Query = &proto.Query_ContractGetInfo{ContractGetInfo: pb}

	return &ContractInfoQuery{inner, pb}
}

func (builder *ContractInfoQuery) SetContractID(id ContractID) *ContractInfoQuery {
	builder.pb.ContractID = id.toProto()
	return builder
}

func (builder *ContractInfoQuery) Execute(client *Client) (ContractInfo, error) {
	resp, err := builder.execute(client)
	if err != nil {
		return ContractInfo{}, err
	}

	adminKey, err := publicKeyFromProto(resp.GetContractGetInfo().GetContractInfo().GetAdminKey())
	if err != nil {
		return ContractInfo{}, err
	}

	return ContractInfo{
		AccountID:         accountIDFromProto(resp.GetContractGetInfo().ContractInfo.AccountID),
		ContractID:        contractIDFromProto(resp.GetContractGetInfo().ContractInfo.ContractID),
		ContractAccountID: resp.GetContractGetInfo().ContractInfo.ContractAccountID,
		AdminKey:          adminKey,
		ExpirationTime:    timeFromProto(resp.GetContractGetInfo().ContractInfo.ExpirationTime),
		AutoRenewPeriod:   durationFromProto(resp.GetContractGetInfo().ContractInfo.AutoRenewPeriod),
		Storage:           uint64(resp.GetContractGetInfo().ContractInfo.Storage),
		ContractMemo:      resp.GetContractGetInfo().ContractInfo.Memo,
	}, nil
}

func (builder *ContractInfoQuery) Cost(client *Client) (Hbar, error) {
	// deleted files return a COST_ANSWER of zero which triggers `INSUFFICIENT_TX_FEE`
	// if you set that as the query payment; 25 tinybar seems to be enough to get
	// `FILE_DELETED` back instead.
	cost, err := builder.QueryBuilder.Cost(client)
	if err != nil {
		return ZeroHbar, err
	}

	// math.Min requires float64 and returns float64
	if cost.AsTinybar() > 25 {
		return cost, nil
	}

	return HbarFromTinybar(25), nil
}

//
// The following _3_ must be copy-pasted at the bottom of **every** _query.go file
// We override the embedded fluent setter methods to return the outer type
//

func (builder *ContractInfoQuery) SetMaxQueryPayment(maxPayment Hbar) *ContractInfoQuery {
	return &ContractInfoQuery{*builder.QueryBuilder.SetMaxQueryPayment(maxPayment), builder.pb}
}

func (builder *ContractInfoQuery) SetQueryPayment(paymentAmount Hbar) *ContractInfoQuery {
	return &ContractInfoQuery{*builder.QueryBuilder.SetQueryPayment(paymentAmount), builder.pb}
}

func (builder *ContractInfoQuery) SetQueryPaymentTransaction(tx Transaction) *ContractInfoQuery {
	return &ContractInfoQuery{*builder.QueryBuilder.SetQueryPaymentTransaction(tx), builder.pb}
}
