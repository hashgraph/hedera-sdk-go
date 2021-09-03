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

func contractInfoFromProtobuf(contractInfo *proto.ContractGetInfoResponse_ContractInfo) (ContractInfo, error) {
	if contractInfo == nil {
		return ContractInfo{}, errParameterNull
	}
	adminKey, err := keyFromProtobuf(contractInfo.GetAdminKey())
	if err != nil {
		return ContractInfo{}, err
	}

	accountID := AccountID{}
	if contractInfo.AccountID != nil {
		accountID = *accountIDFromProtobuf(contractInfo.AccountID)
	}

	contractID := ContractID{}
	if contractInfo.ContractID != nil {
		contractID = *contractIDFromProtobuf(contractInfo.ContractID)
	}

	return ContractInfo{
		AccountID:         accountID,
		ContractID:        contractID,
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
	if data == nil {
		return ContractInfo{}, errByteArrayNull
	}
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
