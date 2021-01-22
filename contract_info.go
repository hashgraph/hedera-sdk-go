package hedera

import (
	protobuf "github.com/golang/protobuf/proto"
	"github.com/hashgraph/hedera-sdk-go/v2/proto"
	"time"
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

func newContractInfo(accountID AccountID, contractID ContractID, contractAccountID string, adminKey Key, expirationTime time.Time,
	autoRenewPeriod time.Duration, storage uint64, ContractMemo string) ContractInfo {
	return ContractInfo{
		AccountID:         accountID,
		ContractID:        contractID,
		ContractAccountID: contractAccountID,
		AdminKey:          adminKey,
		ExpirationTime:    expirationTime,
		AutoRenewPeriod:   autoRenewPeriod,
		Storage:           storage,
		ContractMemo:      ContractMemo,
	}
}

func contractInfoFromProtobuf(contractInfo *proto.ContractGetInfoResponse_ContractInfo) (ContractInfo, error) {
	adminKey, err := keyFromProtobuf(contractInfo.GetAdminKey())
	if err != nil {
		return ContractInfo{}, err
	}

	return ContractInfo{
		AccountID:         accountIDFromProtobuf(contractInfo.AccountID),
		ContractID:        contractIDFromProtobuf(contractInfo.ContractID),
		ContractAccountID: contractInfo.ContractAccountID,
		AdminKey:          adminKey,
		ExpirationTime:    timeFromProtobuf(contractInfo.ExpirationTime),
		AutoRenewPeriod:   durationFromProtobuf(contractInfo.AutoRenewPeriod),
		Storage:           uint64(contractInfo.Storage),
		ContractMemo:      contractInfo.Memo,
		Balance:           contractInfo.Balance,
	}, nil
}

func (contractInfo *ContractInfo) toProtobuf() *proto.ContractGetInfoResponse_ContractInfo {
	return &proto.ContractGetInfoResponse_ContractInfo{
		ContractID:        contractInfo.ContractID.toProtobuf(),
		AccountID:         contractInfo.AccountID.toProtobuf(),
		ContractAccountID: contractInfo.ContractAccountID,
		AdminKey:          contractInfo.AdminKey.toProtoKey(),
		ExpirationTime:    timeToProtobuf(contractInfo.ExpirationTime),
		AutoRenewPeriod:   durationToProtobuf(contractInfo.AutoRenewPeriod),
		Storage:           int64(contractInfo.Storage),
		Memo:              contractInfo.ContractMemo,
		Balance:           contractInfo.Balance,
	}
}

func (contractInfo ContractInfo) ToBytes() []byte {
	data, err := protobuf.Marshal(contractInfo.toProtobuf())
	if err != nil {
		return make([]byte, 0)
	}

	return data
}

func ContractInfoFromBytes(data []byte) (ContractInfo, error) {
	pb := proto.ContractGetInfoResponse_ContractInfo{}
	err := protobuf.Unmarshal(data, &pb)
	if err != nil {
		return ContractInfo{}, err
	}

	info, err := contractInfoFromProtobuf(&pb)
	if err != nil {
		return ContractInfo{}, err
	}

	return info, nil
}
