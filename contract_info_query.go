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
	AdminKey          Ed25519PublicKey
	ExpirationTime    time.Time
	AutoRenewPeriod   time.Duration
	Storage           uint64
	Memo              string
}

func NewContractInfoQuery() *ContractInfoQuery {
	pb := &proto.ContractGetInfoQuery{Header: &proto.QueryHeader{}}

	inner := newQueryBuilder(pb.Header)
	inner.pb.Query = &proto.Query_ContractGetInfo{pb}

	return &ContractInfoQuery{inner, pb}
}

func (builder *ContractInfoQuery) SetContractID(id ContractID) *ContractInfoQuery {
	builder.pb.ContractID = id.toProto()
	return builder
}

func (builder *ContractInfoQuery) Execute(client *Client) (*ContractInfo, error) {
	resp, err := builder.execute(client)
	if err != nil {
		return nil, err
	}

	return &ContractInfo{
		AccountID:         accountIDFromProto(resp.GetContractGetInfo().ContractInfo.AccountID),
		ContractID:        contractIDFromProto(resp.GetContractGetInfo().ContractInfo.ContractID),
		ContractAccountID: resp.GetContractGetInfo().ContractInfo.ContractAccountID,
		AdminKey:          Ed25519PublicKey{keyData: resp.GetContractGetInfo().ContractInfo.AdminKey.GetEd25519()},
		ExpirationTime:    timeFromProto(resp.GetContractGetInfo().ContractInfo.ExpirationTime),
		AutoRenewPeriod:   durationFromProto(resp.GetContractGetInfo().ContractInfo.AutoRenewPeriod),
		Storage:           uint64(resp.GetContractGetInfo().ContractInfo.Storage),
		Memo:              resp.GetContractGetInfo().ContractInfo.Memo,
	}, nil
}
