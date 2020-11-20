package hedera

import (
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
	AdminKey            *Key
	KycKey              *Key
	FreezeKey           *Key
	WipeKey             *Key
	SupplyKey           *Key
	DefaultFreezeStatus *bool
	DefaultKycStatus    *bool
	Deleted             bool
	AutoRenewPeriod     time.Duration
	ExpirationTime      time.Time
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

func tokenInfoFromProtobuf(pb *proto.TokenInfo) TokenInfo {
	var adminKey Key
	if pb.AdminKey != nil {
		adminKey = PublicKey{keyData: pb.AdminKey.GetEd25519()}
	}

	var kycKey Key
	if pb.KycKey != nil {
		kycKey = PublicKey{keyData: pb.KycKey.GetEd25519()}
	}

	var freezeKey Key
	if pb.FreezeKey != nil {
		freezeKey = PublicKey{keyData: pb.FreezeKey.GetEd25519()}
	}

	var wipeKey Key
	if pb.WipeKey != nil {
		wipeKey = PublicKey{keyData: pb.WipeKey.GetEd25519()}
	}

	var supplyKey Key
	if pb.SupplyKey != nil {
		supplyKey = PublicKey{keyData: pb.SupplyKey.GetEd25519()}
	}

	return TokenInfo{
		TokenID:             tokenIDFromProtobuf(pb.TokenId),
		Name:                pb.Name,
		Symbol:              pb.Symbol,
		Decimals:            pb.Decimals,
		TotalSupply:         pb.TotalSupply,
		Treasury:            accountIDFromProtobuf(pb.Treasury),
		AdminKey:            &adminKey,
		KycKey:              &kycKey,
		FreezeKey:           &freezeKey,
		WipeKey:             &wipeKey,
		SupplyKey:           &supplyKey,
		DefaultFreezeStatus: freezeStatusFromProtobuf(pb.DefaultFreezeStatus),
		DefaultKycStatus:    kycStatusFromProtobuf(pb.DefaultKycStatus),
		Deleted:             pb.Deleted,
		AutoRenewPeriod:     time.Duration(pb.GetAutoRenewPeriod().Seconds * time.Second.Nanoseconds()),
		ExpirationTime:      time.Unix(pb.GetExpiry().Seconds, int64(pb.GetExpiry().Nanos)),
	}
}

func (tokenInfo *TokenInfo) toProtobuf() *proto.TokenInfo {
	var adminKey Key
	if tokenInfo.AdminKey != nil {
		adminKey = *tokenInfo.AdminKey
	}

	var kycKey Key
	if tokenInfo.KycKey != nil {
		kycKey = *tokenInfo.KycKey
	}

	var freezeKey Key
	if tokenInfo.FreezeKey != nil {
		freezeKey = *tokenInfo.FreezeKey
	}

	var wipeKey Key
	if tokenInfo.WipeKey != nil {
		wipeKey = *tokenInfo.WipeKey
	}

	var supplyKey Key
	if tokenInfo.SupplyKey != nil {
		supplyKey = *tokenInfo.SupplyKey
	}

	return &proto.TokenInfo{
		TokenId:             tokenInfo.TokenID.toProtobuf(),
		Name:                tokenInfo.Name,
		Symbol:              tokenInfo.Symbol,
		Decimals:            tokenInfo.Decimals,
		TotalSupply:         tokenInfo.TotalSupply,
		Treasury:            tokenInfo.Treasury.toProtobuf(),
		AdminKey:            adminKey.toProtoKey(),
		KycKey:              kycKey.toProtoKey(),
		FreezeKey:           freezeKey.toProtoKey(),
		WipeKey:             wipeKey.toProtoKey(),
		SupplyKey:           supplyKey.toProtoKey(),
		DefaultFreezeStatus: *tokenInfo.FreezeStatusToProtobuf(),
		DefaultKycStatus:    *tokenInfo.KycStatusToProtobuf(),
		Deleted:             tokenInfo.Deleted,
		AutoRenewPeriod:     durationToProtobuf(tokenInfo.AutoRenewPeriod),
		Expiry:              timeToProtobuf(tokenInfo.ExpirationTime),
	}
}
