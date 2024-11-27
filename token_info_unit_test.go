//go:build all || unit
// +build all unit

package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestUnitTokenInfo_Protobuf(t *testing.T) {
	t.Parallel()

	tokenInfo := setupTokenInfo()
	pb := tokenInfo._ToProtobuf()
	actual := _TokenInfoFromProtobuf(pb)

	assertTokenInfo(t, tokenInfo, actual)
}

func TestUnitTokenInfo_Bytes(t *testing.T) {
	t.Parallel()

	tokenInfo := setupTokenInfo()
	pb := tokenInfo.ToBytes()
	actual, _ := TokenInfoFromBytes(pb)

	assertTokenInfo(t, tokenInfo, actual)
}

func TestUnitTokenInfo_ProtobufCoverage(t *testing.T) {
	t.Parallel()

	tokenInfo := setupTokenInfo()

	_true := true
	_false := false

	tokenInfo.DefaultKycStatus = &_false
	tokenInfo.DefaultFreezeStatus = &_false
	tokenInfo.PauseStatus = &_true

	pb := tokenInfo._ToProtobuf()
	actual := _TokenInfoFromProtobuf(pb)

	assertTokenInfo(t, tokenInfo, actual)
}

func setupTokenInfo() TokenInfo {
	adminPK, _ := PrivateKeyGenerate()
	adminPubK := adminPK.PublicKey()

	kycPK, _ := PrivateKeyGenerate()
	kycPubK := kycPK.PublicKey()

	freezePK, _ := PrivateKeyGenerate()
	freezePubK := freezePK.PublicKey()

	wipePK, _ := PrivateKeyGenerate()
	wipePubK := wipePK.PublicKey()

	supplyPK, _ := PrivateKeyGenerate()
	supplyPubK := supplyPK.PublicKey()

	pausePK, _ := PrivateKeyGenerate()
	pausePubK := pausePK.PublicKey()

	metadataPK, _ := PrivateKeyGenerate()
	metadataPubK := metadataPK.PublicKey()

	feeSchedulePK, _ := PrivateKeyGenerate()
	feeSchedulePubK := feeSchedulePK.PublicKey()

	accId, _ := AccountIDFromString("0.0.1111")

	_true := true
	_false := false
	ledgerId := NewLedgerIDTestnet()
	timeDuration := time.Duration(2230000) * time.Second

	timeTime := time.Unix(1230000, 0)
	tokenId, _ := TokenIDFromString("0.0.123")
	feeCollectorAccountId, _ := AccountIDFromString("0.0.123")

	customFees := []Fee{
		NewCustomFixedFee().
			SetAmount(1).
			SetDenominatingTokenID(tokenId).
			SetFeeCollectorAccountID(feeCollectorAccountId),
	}

	return TokenInfo{
		TokenID:             tokenId,
		Name:                "Test Token",
		Symbol:              "TST",
		Decimals:            8,
		TotalSupply:         10000,
		Treasury:            accId,
		AdminKey:            adminPubK,
		KycKey:              kycPubK,
		FreezeKey:           freezePubK,
		WipeKey:             wipePubK,
		SupplyKey:           supplyPubK,
		DefaultFreezeStatus: &_true,
		DefaultKycStatus:    &_true,
		Deleted:             false,
		AutoRenewPeriod:     &timeDuration,
		AutoRenewAccountID:  accId,
		ExpirationTime:      &timeTime,
		TokenMemo:           "test-memo",
		TokenType:           TokenTypeFungibleCommon,
		SupplyType:          TokenSupplyTypeInfinite,
		MaxSupply:           10000000,
		FeeScheduleKey:      feeSchedulePubK,
		CustomFees:          customFees,
		PauseKey:            pausePubK,
		MetadataKey:         metadataPubK,
		Metadata:            testMetadata,
		PauseStatus:         &_false,
		LedgerID:            *ledgerId,
	}
}

func assertTokenInfo(t assert.TestingT, tokenInfo TokenInfo, actual TokenInfo) {
	assert.Equal(t, tokenInfo.TokenID, actual.TokenID)
	assert.Equal(t, tokenInfo.Name, actual.Name)
	assert.Equal(t, tokenInfo.Symbol, actual.Symbol)
	assert.Equal(t, tokenInfo.Decimals, actual.Decimals)
	assert.Equal(t, tokenInfo.TotalSupply, actual.TotalSupply)
	assert.Equal(t, tokenInfo.Treasury, actual.Treasury)
	assert.Equal(t, tokenInfo.AdminKey, actual.AdminKey)
	assert.Equal(t, tokenInfo.KycKey, actual.KycKey)
	assert.Equal(t, tokenInfo.FreezeKey, actual.FreezeKey)
	assert.Equal(t, tokenInfo.WipeKey, actual.WipeKey)
	assert.Equal(t, tokenInfo.SupplyKey, actual.SupplyKey)
	assert.Equal(t, tokenInfo.DefaultFreezeStatus, actual.DefaultFreezeStatus)
	assert.Equal(t, tokenInfo.DefaultKycStatus, actual.DefaultKycStatus)
	assert.Equal(t, tokenInfo.Deleted, actual.Deleted)
	assert.Equal(t, tokenInfo.AutoRenewPeriod, actual.AutoRenewPeriod)
	assert.Equal(t, tokenInfo.AutoRenewAccountID, actual.AutoRenewAccountID)
	assert.Equal(t, tokenInfo.ExpirationTime, actual.ExpirationTime)
	assert.Equal(t, tokenInfo.TokenMemo, actual.TokenMemo)
	assert.Equal(t, tokenInfo.TokenType, actual.TokenType)
	assert.Equal(t, tokenInfo.SupplyType, actual.SupplyType)
	assert.Equal(t, tokenInfo.MaxSupply, actual.MaxSupply)
	assert.Equal(t, tokenInfo.FeeScheduleKey, actual.FeeScheduleKey)
	assert.Equal(t, tokenInfo.PauseKey, actual.PauseKey)
	assert.Equal(t, tokenInfo.MetadataKey, actual.MetadataKey)
	assert.Equal(t, tokenInfo.Metadata, actual.Metadata)
	assert.Equal(t, tokenInfo.PauseStatus, actual.PauseStatus)
	assert.Equal(t, tokenInfo.LedgerID, actual.LedgerID)
}
