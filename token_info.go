package hedera

import (
	protobuf "github.com/golang/protobuf/proto"
	"github.com/hashgraph/hedera-sdk-go/v2/proto"
	"time"
)

type TokenInfo struct {
	TokenID             TokenID
	Name                string
	Symbol              string
	Decimals            uint32
	TotalSupply         uint64
	Treasury            AccountID
	AdminKey            Key
	KycKey              Key
	FreezeKey           Key
	WipeKey             Key
	SupplyKey           Key
	DefaultFreezeStatus *bool
	DefaultKycStatus    *bool
	Deleted             bool
	AutoRenewPeriod     *time.Duration
	AutoRenewAccountID  AccountID
	ExpirationTime      *time.Time
	TokenMemo           string
	TokenType           TokenType
	SupplyType          TokenSupplyType
	MaxSupply           int64
	CustomFees          []Fee
}

func freezeStatusFromProtobuf(pb proto.TokenFreezeStatus) *bool {
	var freezeStatus bool
	switch pb.Number() {
	case 1:
		freezeStatus = true
		break
	case 2:
		freezeStatus = false
		break
	default:
		return nil
	}

	return &freezeStatus
}

func kycStatusFromProtobuf(pb proto.TokenKycStatus) *bool {
	var kycStatus bool
	switch pb.Number() {
	case 1:
		kycStatus = true
		break
	case 2:
		kycStatus = false
		break
	default:
		return nil
	}
	return &kycStatus
}

func (tokenInfo *TokenInfo) FreezeStatusToProtobuf() *proto.TokenFreezeStatus {
	var freezeStatus proto.TokenFreezeStatus

	if tokenInfo.DefaultFreezeStatus == nil {
		return nil
	}

	switch *tokenInfo.DefaultFreezeStatus {
	case true:
		freezeStatus = proto.TokenFreezeStatus_Frozen
		break
	case false:
		freezeStatus = proto.TokenFreezeStatus_Unfrozen
		break
	default:
		freezeStatus = proto.TokenFreezeStatus_FreezeNotApplicable
	}

	return &freezeStatus
}

func (tokenInfo *TokenInfo) KycStatusToProtobuf() *proto.TokenKycStatus {
	var kycStatus proto.TokenKycStatus

	if tokenInfo.DefaultKycStatus == nil {
		return nil
	}

	switch *tokenInfo.DefaultKycStatus {
	case true:
		kycStatus = proto.TokenKycStatus_Granted
		break
	case false:
		kycStatus = proto.TokenKycStatus_Revoked
		break
	default:
		kycStatus = proto.TokenKycStatus_KycNotApplicable
	}

	return &kycStatus
}

func tokenInfoFromProtobuf(pb *proto.TokenInfo, networkName *NetworkName) TokenInfo {
	if pb == nil {
		return TokenInfo{}
	}

	var adminKey Key
	if pb.AdminKey != nil {
		adminKey, _ = keyFromProtobuf(pb.AdminKey, networkName)
	}

	var kycKey Key
	if pb.KycKey != nil {
		kycKey, _ = keyFromProtobuf(pb.KycKey, networkName)
	}

	var freezeKey Key
	if pb.FreezeKey != nil {
		freezeKey, _ = keyFromProtobuf(pb.FreezeKey, networkName)
	}

	var wipeKey Key
	if pb.WipeKey != nil {
		wipeKey, _ = keyFromProtobuf(pb.WipeKey, networkName)
	}

	var supplyKey Key
	if pb.SupplyKey != nil {
		supplyKey, _ = keyFromProtobuf(pb.SupplyKey, networkName)
	}

	var autoRenewPeriod time.Duration
	if pb.AutoRenewPeriod != nil {
		autoRenewPeriod = time.Duration(pb.GetAutoRenewPeriod().Seconds * time.Second.Nanoseconds())
	}

	var expirationTime time.Time
	if pb.Expiry != nil {
		expirationTime = time.Unix(pb.GetExpiry().Seconds, int64(pb.GetExpiry().Nanos))
	}

	var accountID AccountID
	if pb.AutoRenewAccount != nil {
		accountID = accountIDFromProtobuf(pb.AutoRenewAccount, networkName)
	}

	customFees := make([]Fee, 0)
	if pb.CustomFees != nil {
		for _, custom := range pb.CustomFees {
			customFees = append(customFees, customFeeFromProtobuf(custom, networkName))
		}
	}

	return TokenInfo{
		TokenID:             tokenIDFromProtobuf(pb.TokenId, networkName),
		Name:                pb.Name,
		Symbol:              pb.Symbol,
		Decimals:            pb.Decimals,
		TotalSupply:         pb.TotalSupply,
		Treasury:            accountIDFromProtobuf(pb.Treasury, networkName),
		AdminKey:            adminKey,
		KycKey:              kycKey,
		FreezeKey:           freezeKey,
		WipeKey:             wipeKey,
		SupplyKey:           supplyKey,
		DefaultFreezeStatus: freezeStatusFromProtobuf(pb.DefaultFreezeStatus),
		DefaultKycStatus:    kycStatusFromProtobuf(pb.DefaultKycStatus),
		Deleted:             pb.Deleted,
		AutoRenewPeriod:     &autoRenewPeriod,
		AutoRenewAccountID:  accountID,
		ExpirationTime:      &expirationTime,
		TokenMemo:           pb.Memo,
		TokenType:           TokenType(pb.TokenType),
		SupplyType:          TokenSupplyType(pb.SupplyType),
		MaxSupply:           pb.MaxSupply,
		CustomFees:          customFees,
	}
}

func (tokenInfo *TokenInfo) toProtobuf() *proto.TokenInfo {
	var adminKey *proto.Key
	if tokenInfo.AdminKey != nil {
		adminKey = tokenInfo.AdminKey.toProtoKey()
	}

	var kycKey *proto.Key
	if tokenInfo.KycKey != nil {
		kycKey = tokenInfo.KycKey.toProtoKey()
	}

	var freezeKey *proto.Key
	if tokenInfo.FreezeKey != nil {
		freezeKey = tokenInfo.FreezeKey.toProtoKey()
	}

	var wipeKey *proto.Key
	if tokenInfo.WipeKey != nil {
		wipeKey = tokenInfo.WipeKey.toProtoKey()
	}

	var supplyKey *proto.Key
	if tokenInfo.SupplyKey != nil {
		supplyKey = tokenInfo.SupplyKey.toProtoKey()
	}

	var autoRenewPeriod *proto.Duration
	if tokenInfo.AutoRenewPeriod != nil {
		autoRenewPeriod = durationToProtobuf(*tokenInfo.AutoRenewPeriod)
	}

	var expirationTime *proto.Timestamp
	if tokenInfo.ExpirationTime != nil {
		expirationTime = timeToProtobuf(*tokenInfo.ExpirationTime)
	}

	customFees := make([]*proto.CustomFee, 0)
	if tokenInfo.CustomFees != nil {
		for _, customFee := range tokenInfo.CustomFees {
			customFees = append(customFees, customFee.toProtobuf())
		}
	}

	return &proto.TokenInfo{
		TokenId:             tokenInfo.TokenID.toProtobuf(),
		Name:                tokenInfo.Name,
		Symbol:              tokenInfo.Symbol,
		Decimals:            tokenInfo.Decimals,
		TotalSupply:         tokenInfo.TotalSupply,
		Treasury:            tokenInfo.Treasury.toProtobuf(),
		AdminKey:            adminKey,
		KycKey:              kycKey,
		FreezeKey:           freezeKey,
		WipeKey:             wipeKey,
		SupplyKey:           supplyKey,
		DefaultFreezeStatus: *tokenInfo.FreezeStatusToProtobuf(),
		DefaultKycStatus:    *tokenInfo.KycStatusToProtobuf(),
		Deleted:             tokenInfo.Deleted,
		AutoRenewAccount:    tokenInfo.AutoRenewAccountID.toProtobuf(),
		AutoRenewPeriod:     autoRenewPeriod,
		Expiry:              expirationTime,
		Memo:                tokenInfo.TokenMemo,
		TokenType:           proto.TokenType(tokenInfo.TokenType),
		SupplyType:          proto.TokenSupplyType(tokenInfo.SupplyType),
		MaxSupply:           tokenInfo.MaxSupply,
		CustomFees:          customFees,
	}
}

func (tokenInfo TokenInfo) ToBytes() []byte {
	data, err := protobuf.Marshal(tokenInfo.toProtobuf())
	if err != nil {
		return make([]byte, 0)
	}

	return data
}

func TokenInfoFromBytes(data []byte) (TokenInfo, error) {
	if data == nil {
		return TokenInfo{}, errByteArrayNull
	}
	pb := proto.TokenInfo{}
	err := protobuf.Unmarshal(data, &pb)
	if err != nil {
		return TokenInfo{}, err
	}

	return tokenInfoFromProtobuf(&pb, nil), nil
}
