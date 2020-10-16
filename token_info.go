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
	DefaultFreezeStatus bool
	DefaultKycStatus    bool
	IsDelete            bool
	AutoRenewPeriod     uint64
	ExpirationTime      uint64
}

func freezeStatusFromProtobuf(freezeStatus TokenFreezeStatus) bool {
	if freezeStatus.Number() == 1 {
		return true
	}
	return false
}

func kycStatusFromProtobuf(kycStatus TokenKycStatus) bool {
	if kycStatus.Number() == 1 {
		return true
	}
	return false
}

func tokenInfoFromProtobuf(tokenInfoResponse *proto.TokenGetInfoResponse) TokenInfo {
	tokenInfo := tokenInfoResponse.TokenInfo
	return TokenInfo{
		TokenID:     tokenIDFromProtobuf(tokenInfo.TokenId),
		Name:        tokenInfo.Name,
		Symbol:      tokenInfo.Symbol,
		Decimals:    tokenInfo.Decimals,
		TotalSupply: tokenInfo.TotalSupply,
		Treasury:    accountIDFromProtobuf(tokenInfo.Treasury),
		AdminKey: PublicKey{
			keyData: tokenInfo.AdminKey.GetEd25519(),
		},
		KycKey: PublicKey{
			keyData: tokenInfo.KycKey.GetEd25519(),
		},
		FreezeKey: PublicKey{
			keyData: tokenInfo.FreezeKey.GetEd25519(),
		},
		WipeKey: PublicKey{
			keyData: tokenInfo.WipeKey.GetEd25519(),
		},
		SupplyKey: PublicKey{
			keyData: tokenInfo.SupplyKey.GetEd25519(),
		},
		DefaultFreezeStatus: freezeStatusFromProtobuf(tokenInfo.DefaultFreezeStatus),
		DefaultKycStatus:    kycStatusFromProtobuf(tokenInfo.DefaultKycStatus),
		IsDelete:            tokenInfo.IsDeleted,
		AutoRenewPeriod:     tokenInfo.AutoRenewPeriod,
		ExpirationTime:      tokenInfo.Expiry,
	}
}
