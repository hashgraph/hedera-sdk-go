package hedera

import (
	"time"

	"github.com/hashgraph/hedera-sdk-go/v2/proto"
	protobuf "google.golang.org/protobuf/proto"
)

type ContractInfo struct {
	AccountID         AccountID
	ContractID        ContractID
	ContractAccountID string
	AdminKey          Key
	ExpirationTime    time.Time
	AutoRenewPeriod   time.Duration
	Storage           uint64
	ContractMemo      string
	Balance           uint64
}

func _ContractInfoFromProtobuf(contractInfo *proto.ContractGetInfoResponse_ContractInfo) (ContractInfo, error) {
	if contractInfo == nil {
		return ContractInfo{}, errParameterNull
	}
	adminKey, err := _KeyFromProtobuf(contractInfo.GetAdminKey())
	if err != nil {
		return ContractInfo{}, err
	}

	accountID := AccountID{}
	if contractInfo.AccountID != nil {
		accountID = *_AccountIDFromProtobuf(contractInfo.AccountID)
	}

	contractID := ContractID{}
	if contractInfo.ContractID != nil {
		contractID = *_ContractIDFromProtobuf(contractInfo.ContractID)
	}

	return ContractInfo{
		AccountID:         accountID,
		ContractID:        contractID,
		ContractAccountID: contractInfo.ContractAccountID,
		AdminKey:          adminKey,
		ExpirationTime:    _TimeFromProtobuf(contractInfo.ExpirationTime),
		AutoRenewPeriod:   _DurationFromProtobuf(contractInfo.AutoRenewPeriod),
		Storage:           uint64(contractInfo.Storage),
		ContractMemo:      contractInfo.Memo,
		Balance:           contractInfo.Balance,
	}, nil
}

func (contractInfo *ContractInfo) _ToProtobuf() *proto.ContractGetInfoResponse_ContractInfo {
	return &proto.ContractGetInfoResponse_ContractInfo{
		ContractID:        contractInfo.ContractID._ToProtobuf(),
		AccountID:         contractInfo.AccountID._ToProtobuf(),
		ContractAccountID: contractInfo.ContractAccountID,
		AdminKey:          contractInfo.AdminKey._ToProtoKey(),
		ExpirationTime:    _TimeToProtobuf(contractInfo.ExpirationTime),
		AutoRenewPeriod:   _DurationToProtobuf(contractInfo.AutoRenewPeriod),
		Storage:           int64(contractInfo.Storage),
		Memo:              contractInfo.ContractMemo,
		Balance:           contractInfo.Balance,
	}
}

func (contractInfo ContractInfo) ToBytes() []byte {
	data, err := protobuf.Marshal(contractInfo._ToProtobuf())
	if err != nil {
		return make([]byte, 0)
	}

	return data
}

func ContractInfoFromBytes(data []byte) (ContractInfo, error) {
	if data == nil {
		return ContractInfo{}, errByteArrayNull
	}
	pb := proto.ContractGetInfoResponse_ContractInfo{}
	err := protobuf.Unmarshal(data, &pb)
	if err != nil {
		return ContractInfo{}, err
	}

	info, err := _ContractInfoFromProtobuf(&pb)
	if err != nil {
		return ContractInfo{}, err
	}

	return info, nil
}
