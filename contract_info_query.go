package hedera

import (
	"github.com/hashgraph/hedera-sdk-go/proto"
	"time"
)

// ContractInfoQuery retrieves information about a smart contract instance. This includes the account that it uses, the
// file containing its bytecode, and the time when it will expire.
type ContractInfoQuery struct {
	QueryBuilder
	pb *proto.ContractGetInfoQuery
}

// ContractInfo is the information about the contract instance returned by a ContractInfoQuery
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

// NewContractInfoQuery creates a ContractInfoQuery builder which can be used to construct and execute a
// Contract Get Info Query.
func NewContractInfoQuery() *ContractInfoQuery {
	pb := &proto.ContractGetInfoQuery{Header: &proto.QueryHeader{}}

	inner := newQueryBuilder(pb.Header)
	inner.pb.Query = &proto.Query_ContractGetInfo{ContractGetInfo: pb}

	return &ContractInfoQuery{inner, pb}
}

// SetContractID sets the contract for which information is requested
func (builder *ContractInfoQuery) SetContractID(id ContractID) *ContractInfoQuery {
	builder.pb.ContractID = id.toProto()
	return builder
}

// Execute executes the ContractInfoQuery using the provided client
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

// Cost is a wrapper around the standard Cost function for a query. It must exist because deleted files return a
// COST_ANSWER of zero which triggers an INSUFFICIENT_TX_FEE response Status if set as the query payment. However,
// 25 tinybar seems to be enough to get FILE_DELETED back instead, so that is used instead.
func (builder *ContractInfoQuery) Cost(client *Client) (Hbar, error) {
	cost, err := builder.QueryBuilder.GetCost(client)
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
