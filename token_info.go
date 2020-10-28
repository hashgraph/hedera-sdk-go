package hedera

import "github.com/hashgraph/hedera-sdk-go/proto"

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
	IsDelete            bool
	AutoRenewPeriod     uint64
	ExpirationTime      uint64
}

func freezeStatusFromProtobuf(pb *proto.TokenFreezeStatus) *bool {
	var freezeStatus bool
	switch pb.Number() {
	case 1:
		freezeStatus = true
	case 2:
		freezeStatus = false
	default:
		return nil
	}

	return &freezeStatus
}

func kycStatusFromProtobuf(pb *proto.TokenKycStatus) *bool {
	var kycStatus bool
	switch pb.Number() {
	case 1:
		kycStatus = true
	case 2:
		kycStatus = false
	default:
		return nil
	}
	return &kycStatus
}

func tokenInfoFromProtobuf(tokenInfo *proto.TokenInfo) TokenInfo {
	var adminKey PublicKey
	if tokenInfo.AdminKey != nil {
		adminKey = PublicKey{keyData: tokenInfo.AdminKey.GetEd25519()}
	}

	var kycKey PublicKey
	if tokenInfo.KycKey != nil {
		kycKey = PublicKey{keyData: tokenInfo.KycKey.GetEd25519()}
	}

	var freezeKey PublicKey
	if tokenInfo.FreezeKey != nil {
		freezeKey = PublicKey{keyData: tokenInfo.FreezeKey.GetEd25519()}
	}

	var wipeKey PublicKey
	if tokenInfo.WipeKey != nil {
		wipeKey = PublicKey{keyData: tokenInfo.WipeKey.GetEd25519()}
	}

	var supplyKey PublicKey
	if tokenInfo.SupplyKey != nil {
		supplyKey = PublicKey{keyData: tokenInfo.SupplyKey.GetEd25519()}
	}

	return TokenInfo{
		TokenID:             tokenIDFromProtobuf(tokenInfo.TokenId),
		Name:                tokenInfo.Name,
		Symbol:              tokenInfo.Symbol,
		Decimals:            tokenInfo.Decimals,
		TotalSupply:         tokenInfo.TotalSupply,
		Treasury:            accountIDFromProtobuf(tokenInfo.Treasury),
		AdminKey:            adminKey,
		KycKey:              kycKey,
		FreezeKey:           freezeKey,
		WipeKey:             wipeKey,
		SupplyKey:           supplyKey,
		DefaultFreezeStatus: freezeStatusFromProtobuf(&tokenInfo.DefaultFreezeStatus),
		DefaultKycStatus:    kycStatusFromProtobuf(&tokenInfo.DefaultKycStatus),
		IsDelete:            tokenInfo.IsDeleted,
		AutoRenewPeriod:     tokenInfo.AutoRenewPeriod,
		ExpirationTime:      tokenInfo.Expiry,
	}
}
