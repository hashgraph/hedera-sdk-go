package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"time"

	"github.com/hiero-ledger/hiero-sdk-go/v2/proto/services"
	protobuf "google.golang.org/protobuf/proto"
)

// TokenInfo is the information about a token
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
	FeeScheduleKey      Key
	CustomFees          []Fee
	PauseKey            Key
	PauseStatus         *bool
	MetadataKey         Key
	Metadata            []byte
	LedgerID            LedgerID
}

func _FreezeStatusFromProtobuf(pb services.TokenFreezeStatus) *bool {
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

func _KycStatusFromProtobuf(pb services.TokenKycStatus) *bool {
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

func _PauseStatusFromProtobuf(pb services.TokenPauseStatus) *bool {
	var pauseStatus bool
	switch pb.Number() {
	case 1:
		pauseStatus = true
	case 2:
		pauseStatus = false
	default:
		return nil
	}
	return &pauseStatus
}

// FreezeStatusToProtobuf converts the TokenInfo's DefaultFreezeStatus to a protobuf TokenFreezeStatus
func (tokenInfo *TokenInfo) FreezeStatusToProtobuf() *services.TokenFreezeStatus {
	var freezeStatus services.TokenFreezeStatus

	if tokenInfo.DefaultFreezeStatus == nil {
		freezeStatus = services.TokenFreezeStatus_FreezeNotApplicable
		return &freezeStatus
	}

	switch *tokenInfo.DefaultFreezeStatus {
	case true:
		freezeStatus = services.TokenFreezeStatus_Frozen
	case false:
		freezeStatus = services.TokenFreezeStatus_Unfrozen
	default:
		freezeStatus = services.TokenFreezeStatus_FreezeNotApplicable
	}

	return &freezeStatus
}

// KycStatusToProtobuf converts the TokenInfo's DefaultKycStatus to a protobuf TokenKycStatus
func (tokenInfo *TokenInfo) KycStatusToProtobuf() *services.TokenKycStatus {
	var kycStatus services.TokenKycStatus

	if tokenInfo.DefaultKycStatus == nil {
		kycStatus = services.TokenKycStatus_KycNotApplicable
		return &kycStatus
	}

	switch *tokenInfo.DefaultKycStatus {
	case true:
		kycStatus = services.TokenKycStatus_Granted
	case false:
		kycStatus = services.TokenKycStatus_Revoked
	default:
		kycStatus = services.TokenKycStatus_KycNotApplicable
	}

	return &kycStatus
}

// PauseStatusToProtobuf converts the TokenInfo's PauseStatus to a protobuf TokenPauseStatus
func (tokenInfo *TokenInfo) PauseStatusToProtobuf() *services.TokenPauseStatus {
	var pauseStatus services.TokenPauseStatus

	if tokenInfo.PauseStatus == nil {
		pauseStatus = services.TokenPauseStatus_PauseNotApplicable
		return &pauseStatus
	}

	switch *tokenInfo.PauseStatus {
	case true:
		pauseStatus = services.TokenPauseStatus_Paused
	case false:
		pauseStatus = services.TokenPauseStatus_Unpaused
	default:
		pauseStatus = services.TokenPauseStatus_PauseNotApplicable
	}

	return &pauseStatus
}

func _TokenInfoFromProtobuf(pb *services.TokenInfo) TokenInfo {
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

	var pauseKey Key
	if pb.PauseKey != nil {
		pauseKey, _ = _KeyFromProtobuf(pb.PauseKey)
	}

	var metadataKey Key
	if pb.MetadataKey != nil {
		metadataKey, _ = _KeyFromProtobuf(pb.MetadataKey)
	}

	var feeScheduleKey Key
	if pb.FeeScheduleKey != nil {
		feeScheduleKey, _ = _KeyFromProtobuf(pb.FeeScheduleKey)
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
	if pb.Treasury != nil {
		treasury = *_AccountIDFromProtobuf(pb.Treasury)
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

	var metadata []byte
	if pb.Metadata != nil {
		metadata = pb.Metadata
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
		FeeScheduleKey:      feeScheduleKey,
		CustomFees:          customFees,
		PauseKey:            pauseKey,
		MetadataKey:         metadataKey,
		Metadata:            metadata,
		PauseStatus:         _PauseStatusFromProtobuf(pb.PauseStatus),
		LedgerID:            LedgerID{pb.LedgerId},
	}
}

func (tokenInfo *TokenInfo) _ToProtobuf() *services.TokenInfo {
	var adminKey *services.Key
	if tokenInfo.AdminKey != nil {
		adminKey = tokenInfo.AdminKey._ToProtoKey()
	}

	var kycKey *services.Key
	if tokenInfo.KycKey != nil {
		kycKey = tokenInfo.KycKey._ToProtoKey()
	}

	var freezeKey *services.Key
	if tokenInfo.FreezeKey != nil {
		freezeKey = tokenInfo.FreezeKey._ToProtoKey()
	}

	var wipeKey *services.Key
	if tokenInfo.WipeKey != nil {
		wipeKey = tokenInfo.WipeKey._ToProtoKey()
	}

	var supplyKey *services.Key
	if tokenInfo.SupplyKey != nil {
		supplyKey = tokenInfo.SupplyKey._ToProtoKey()
	}

	var pauseKey *services.Key
	if tokenInfo.PauseKey != nil {
		pauseKey = tokenInfo.PauseKey._ToProtoKey()
	}

	var metadataKey *services.Key
	if tokenInfo.MetadataKey != nil {
		metadataKey = tokenInfo.MetadataKey._ToProtoKey()
	}

	var feeScheduleKey *services.Key
	if tokenInfo.FeeScheduleKey != nil {
		feeScheduleKey = tokenInfo.FeeScheduleKey._ToProtoKey()
	}

	var autoRenewPeriod *services.Duration
	if tokenInfo.AutoRenewPeriod != nil {
		autoRenewPeriod = _DurationToProtobuf(*tokenInfo.AutoRenewPeriod)
	}

	var expirationTime *services.Timestamp
	if tokenInfo.ExpirationTime != nil {
		expirationTime = _TimeToProtobuf(*tokenInfo.ExpirationTime)
	}

	customFees := make([]*services.CustomFee, 0)
	if tokenInfo.CustomFees != nil {
		for _, customFee := range tokenInfo.CustomFees {
			customFees = append(customFees, customFee._ToProtobuf())
		}
	}

	var metadata []byte
	if tokenInfo.Metadata != nil {
		metadata = tokenInfo.Metadata
	}

	return &services.TokenInfo{
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
		TokenType:           services.TokenType(tokenInfo.TokenType),
		SupplyType:          services.TokenSupplyType(tokenInfo.SupplyType),
		MaxSupply:           tokenInfo.MaxSupply,
		FeeScheduleKey:      feeScheduleKey,
		CustomFees:          customFees,
		PauseKey:            pauseKey,
		MetadataKey:         metadataKey,
		Metadata:            metadata,
		PauseStatus:         *tokenInfo.PauseStatusToProtobuf(),
		LedgerId:            tokenInfo.LedgerID.ToBytes(),
	}
}

// ToBytes returns the byte representation of the TokenInfo
func (tokenInfo TokenInfo) ToBytes() []byte {
	data, err := protobuf.Marshal(tokenInfo._ToProtobuf())
	if err != nil {
		return make([]byte, 0)
	}

	return data
}

// TokenInfoFromBytes returns a TokenInfo struct from a raw protobuf byte array
func TokenInfoFromBytes(data []byte) (TokenInfo, error) {
	if data == nil {
		return TokenInfo{}, errByteArrayNull
	}
	pb := services.TokenInfo{}
	err := protobuf.Unmarshal(data, &pb)
	if err != nil {
		return TokenInfo{}, err
	}

	return _TokenInfoFromProtobuf(&pb), nil
}
