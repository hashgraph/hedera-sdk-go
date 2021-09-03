package hedera

import (
	"time"

	"github.com/hashgraph/hedera-sdk-go/v2/proto"
	protobuf "google.golang.org/protobuf/proto"
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

func _FreezeStatusFromProtobuf(pb proto.TokenFreezeStatus) *bool {
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

func _KycStatusFromProtobuf(pb proto.TokenKycStatus) *bool {
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

func (tokenInfo *TokenInfo) FreezeStatusToProtobuf() *proto.TokenFreezeStatus {
	var freezeStatus proto.TokenFreezeStatus

	if tokenInfo.DefaultFreezeStatus == nil {
		return nil
	}

	switch *tokenInfo.DefaultFreezeStatus {
	case true:
		freezeStatus = proto.TokenFreezeStatus_Frozen
	case false:
		freezeStatus = proto.TokenFreezeStatus_Unfrozen
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
	case false:
		kycStatus = proto.TokenKycStatus_Revoked
	default:
		kycStatus = proto.TokenKycStatus_KycNotApplicable
	}

	return &kycStatus
}

func _TokenInfoFromProtobuf(pb *proto.TokenInfo) TokenInfo {
	if pb == nil {
		return TokenInfo{}
	}

	var adminKey Key
	if pb.AdminKey != nil {
		adminKey, _ = _KeyFromProtobuf(pb.AdminKey)
	}

	var kycKey Key
	if pb.KycKey != nil {
		kycKey, _ = _KeyFromProtobuf(pb.KycKey)
	}

	var freezeKey Key
	if pb.FreezeKey != nil {
		freezeKey, _ = _KeyFromProtobuf(pb.FreezeKey)
	}

	var wipeKey Key
	if pb.WipeKey != nil {
		wipeKey, _ = _KeyFromProtobuf(pb.WipeKey)
	}

	var supplyKey Key
	if pb.SupplyKey != nil {
		supplyKey, _ = _KeyFromProtobuf(pb.SupplyKey)
	}

	var autoRenewPeriod time.Duration
	if pb.AutoRenewPeriod != nil {
		autoRenewPeriod = time.Duration(pb.GetAutoRenewPeriod().Seconds * time.Second.Nanoseconds())
	}

	var expirationTime time.Time
	if pb.Expiry != nil {
		expirationTime = time.Unix(pb.GetExpiry().Seconds, int64(pb.GetExpiry().Nanos))
	}

	var autoRenewAccountID AccountID
	if pb.AutoRenewAccount != nil {
		autoRenewAccountID = *_AccountIDFromProtobuf(pb.AutoRenewAccount)
	}

	var treasury AccountID
	if pb.AutoRenewAccount != nil {
		treasury = *_AccountIDFromProtobuf(pb.AutoRenewAccount)
	}

	customFees := make([]Fee, 0)
	if pb.CustomFees != nil {
		for _, custom := range pb.CustomFees {
			customFees = append(customFees, _CustomFeeFromProtobuf(custom))
		}
	}

	tokenID := TokenID{}
	if pb.TokenId != nil {
		tokenID = *_TokenIDFromProtobuf(pb.TokenId)
	}

	return TokenInfo{
		TokenID:             tokenID,
		Name:                pb.Name,
		Symbol:              pb.Symbol,
		Decimals:            pb.Decimals,
		TotalSupply:         pb.TotalSupply,
		Treasury:            treasury,
		AdminKey:            adminKey,
		KycKey:              kycKey,
		FreezeKey:           freezeKey,
		WipeKey:             wipeKey,
		SupplyKey:           supplyKey,
		DefaultFreezeStatus: _FreezeStatusFromProtobuf(pb.DefaultFreezeStatus),
		DefaultKycStatus:    _KycStatusFromProtobuf(pb.DefaultKycStatus),
		Deleted:             pb.Deleted,
		AutoRenewPeriod:     &autoRenewPeriod,
		AutoRenewAccountID:  autoRenewAccountID,
		ExpirationTime:      &expirationTime,
		TokenMemo:           pb.Memo,
		TokenType:           TokenType(pb.TokenType),
		SupplyType:          TokenSupplyType(pb.SupplyType),
		MaxSupply:           pb.MaxSupply,
		CustomFees:          customFees,
	}
}

func (tokenInfo *TokenInfo) _ToProtobuf() *proto.TokenInfo {
	var adminKey *proto.Key
	if tokenInfo.AdminKey != nil {
		adminKey = tokenInfo.AdminKey._ToProtoKey()
	}

	var kycKey *proto.Key
	if tokenInfo.KycKey != nil {
		kycKey = tokenInfo.KycKey._ToProtoKey()
	}

	var freezeKey *proto.Key
	if tokenInfo.FreezeKey != nil {
		freezeKey = tokenInfo.FreezeKey._ToProtoKey()
	}

	var wipeKey *proto.Key
	if tokenInfo.WipeKey != nil {
		wipeKey = tokenInfo.WipeKey._ToProtoKey()
	}

	var supplyKey *proto.Key
	if tokenInfo.SupplyKey != nil {
		supplyKey = tokenInfo.SupplyKey._ToProtoKey()
	}

	var autoRenewPeriod *proto.Duration
	if tokenInfo.AutoRenewPeriod != nil {
		autoRenewPeriod = _DurationToProtobuf(*tokenInfo.AutoRenewPeriod)
	}

	var expirationTime *proto.Timestamp
	if tokenInfo.ExpirationTime != nil {
		expirationTime = _TimeToProtobuf(*tokenInfo.ExpirationTime)
	}

	customFees := make([]*proto.CustomFee, 0)
	if tokenInfo.CustomFees != nil {
		for _, customFee := range tokenInfo.CustomFees {
			customFees = append(customFees, customFee._ToProtobuf())
		}
	}

	return &proto.TokenInfo{
		TokenId:             tokenInfo.TokenID._ToProtobuf(),
		Name:                tokenInfo.Name,
		Symbol:              tokenInfo.Symbol,
		Decimals:            tokenInfo.Decimals,
		TotalSupply:         tokenInfo.TotalSupply,
		Treasury:            tokenInfo.Treasury._ToProtobuf(),
		AdminKey:            adminKey,
		KycKey:              kycKey,
		FreezeKey:           freezeKey,
		WipeKey:             wipeKey,
		SupplyKey:           supplyKey,
		DefaultFreezeStatus: *tokenInfo.FreezeStatusToProtobuf(),
		DefaultKycStatus:    *tokenInfo.KycStatusToProtobuf(),
		Deleted:             tokenInfo.Deleted,
		AutoRenewAccount:    tokenInfo.AutoRenewAccountID._ToProtobuf(),
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
	data, err := protobuf.Marshal(tokenInfo._ToProtobuf())
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

	return _TokenInfoFromProtobuf(&pb), nil
}
