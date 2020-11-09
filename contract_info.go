package hedera

import (
	"github.com/hashgraph/hedera-sdk-go/proto"
	"time"
)

type ContractInfo struct {
	AccountID         AccountID
	ContractID        ContractID
	ContractAccountID string
	AdminKey          PublicKey
	ExpirationTime    time.Time
	AutoRenewPeriod   time.Duration
	Storage           uint64
	ContractMemo      string
	Balance           uint64
	IsDeleted         bool
	TokenRelationships map[TokenID]TokenRelationship
}

func newContractInfo(accountID AccountID, contractID ContractID, contractAccountID string, adminKey PublicKey, expirationTime time.Time,
	autoRenewPeriod time.Duration, storage uint64, ContractMemo string, isDeleted bool, tokenRelationship map[TokenID]TokenRelationship) ContractInfo {
	return ContractInfo{
		AccountID:         accountID,
		ContractID:        contractID,
		ContractAccountID: contractAccountID,
		AdminKey:          adminKey,
		ExpirationTime:    expirationTime,
		AutoRenewPeriod:   autoRenewPeriod,
		Storage:           storage,
		ContractMemo:      ContractMemo,
		IsDeleted:         isDeleted,
		TokenRelationships: tokenRelationship,
	}
}

func contractInfoFromProtobuf(contractInfo *proto.ContractGetInfoResponse_ContractInfo) (ContractInfo, error) {
	adminKey, err := keyFromProtobuf(contractInfo.GetAdminKey())
	if err != nil {
		return ContractInfo{}, err
	}

	var tokenRelationship = make(map[TokenID]TokenRelationship, len(contractInfo.TokenRelationships))
	if contractInfo.TokenRelationships != nil {
		for _, relation := range contractInfo.TokenRelationships {
			tokenRelationship[tokenIDFromProtobuf(relation.TokenId)] = tokenRelationshipFromProtobuf(relation)
		}
	}

	return ContractInfo{
		AccountID:         accountIDFromProtobuf(contractInfo.AccountID),
		ContractID:        contractIDFromProtobuf(contractInfo.ContractID),
		ContractAccountID: contractInfo.ContractAccountID,
		AdminKey: PublicKey{
			keyData: adminKey.toProtoKey().GetEd25519(),
		},
		ExpirationTime:    timeFromProtobuf(contractInfo.ExpirationTime),
		AutoRenewPeriod:   durationFromProtobuf(contractInfo.AutoRenewPeriod),
		Storage:           uint64(contractInfo.Storage),
		ContractMemo:      contractInfo.Memo,
		Balance:           contractInfo.Balance,
		IsDeleted:         contractInfo.Deleted,
		TokenRelationships: tokenRelationship,
	}, nil
}

func (contractInfo *ContractInfo) toProtobuf() *proto.ContractGetInfoResponse_ContractInfo {
	var tokenRelationship = make([]*proto.TokenRelationship, len(contractInfo.TokenRelationships))
	count := 0
	if len(tokenRelationship) > 0 {
		for _, relation := range contractInfo.TokenRelationships {
			tokenRelationship[count] = relation.toProtobuf()
			count++
		}
	}

	return &proto.ContractGetInfoResponse_ContractInfo{
		ContractID:         contractInfo.ContractID.toProtobuf(),
		AccountID:          contractInfo.AccountID.toProtobuf(),
		ContractAccountID:  contractInfo.ContractAccountID,
		AdminKey:           contractInfo.AdminKey.toProtoKey(),
		ExpirationTime:     timeToProtobuf(contractInfo.ExpirationTime),
		AutoRenewPeriod:    durationToProtobuf(contractInfo.AutoRenewPeriod),
		Storage:            int64(contractInfo.Storage),
		Memo:               contractInfo.ContractMemo,
		Balance:            contractInfo.Balance,
		TokenRelationships: tokenRelationship,
	}
}
