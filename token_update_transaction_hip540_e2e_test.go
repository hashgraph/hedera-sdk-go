//go:build all || e2e
// +build all e2e

package hiero

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// SPDX-License-Identifier: Apache-2.0

// HIP-540 tests

// KeyType represents different types of keys
type KeyType int

// Define enum values for different key types
const (
	WIPE_KEY KeyType = iota
	KYC_KEY
	SUPPLY_KEY
	FREEZE_KEY
	FEE_SCHEDULE_KEY
	PAUSE_KEY
	METADATA_KEY
	ADMIN_KEY
	ALL
	LOWER_PRIVILEGE
	NONE
)

func TestIntegrationTokenUpdateTransactionUpdateKeysToEmptyKeyListMakesTokenImmutable(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)

	validNewKey, err := PrivateKeyGenerateEd25519()
	require.NoError(t, err)

	// Make token immutable
	resp, tokenID, err := createTokenWithKeysAndUpdateTokenKeyHelper(t, ALL, ALL, env.Client, NewKeyList(), env.OperatorKey, env.OperatorKey, NO_VALIDATION)
	require.NoError(t, err)
	_, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	// Verify token is immutable
	resp, err = updateTokenKeysHelper(t, tokenID, ALL, env.Client, validNewKey, env.OperatorKey, NO_VALIDATION)
	require.NoError(t, err)
	_, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.ErrorContains(t, err, "exceptional receipt status: TOKEN_IS_IMMUTABLE")

	tokenInfo, err := NewTokenInfoQuery().SetTokenID(tokenID).Execute(env.Client)
	verifyAdminKeyNil(t, tokenInfo)
	verifyLowerPrivilegeKeysNil(t, tokenInfo)
}

func TestIntegrationTokenUpdateTransactionUpdateKeysToZeroKeyFails(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)

	zeroNewKey, err := ZeroKey()
	require.NoError(t, err)

	// Make token immutable
	resp, _, err := createTokenWithKeysAndUpdateTokenKeyHelper(t, ALL, ALL, env.Client, zeroNewKey, env.OperatorKey, env.OperatorKey, NO_VALIDATION)
	require.NoError(t, err)
	_, err = resp.SetValidateStatus(true).GetReceipt(env.Client)

	require.ErrorContains(t, err, "INVALID_SIGNATURE")
}

func TestIntegrationTokenUpdateTransactionUpdateLowerPrivilegeKeysWithAdminKeyFullValidation(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)

	// Update lower privilege keys to zero key list with admin key
	zeroKey, err := ZeroKey()
	require.NoError(t, err)

	resp, tokenID, err := createTokenWithKeysAndUpdateTokenKeyHelper(t, ALL, LOWER_PRIVILEGE, env.Client, zeroKey, env.OperatorKey, env.OperatorKey, FULL_VALIDATION)
	require.NoError(t, err)
	_, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)
	// Verify keys are set to zero
	tokenInfo, err := NewTokenInfoQuery().SetTokenID(tokenID).Execute(env.Client)
	verifyLowerPrivilegeKeys(t, tokenInfo, zeroKey)

	// Update lower privilege keys to valid key list with admin key
	validKey, err := PrivateKeyGenerateEd25519()
	require.NoError(t, err)
	resp, err = updateTokenKeysHelper(t, tokenID, LOWER_PRIVILEGE, env.Client, validKey, env.OperatorKey, FULL_VALIDATION)
	require.NoError(t, err)
	_, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)
	// Verify keys are set to valid key
	tokenInfo, err = NewTokenInfoQuery().SetTokenID(tokenID).Execute(env.Client)
	verifyLowerPrivilegeKeys(t, tokenInfo, validKey.PublicKey())

	// Update lower privilege keys to empty key list with admin key
	resp, err = updateTokenKeysHelper(t, tokenID, LOWER_PRIVILEGE, env.Client, NewKeyList(), env.OperatorKey, FULL_VALIDATION)
	require.NoError(t, err)
	_, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)
	// Verify keys are set to empty
	tokenInfo, err = NewTokenInfoQuery().SetTokenID(tokenID).Execute(env.Client)
	verifyLowerPrivilegeKeysNil(t, tokenInfo)
}

func TestIntegrationTokenUpdateTransactionUpdateLowerPrivilegeKeysWithAdminKeyNoValidation(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)

	// Update lower privilege keys to zero key list with admin key
	zeroKey, err := ZeroKey()
	require.NoError(t, err)

	resp, tokenID, err := createTokenWithKeysAndUpdateTokenKeyHelper(t, ALL, LOWER_PRIVILEGE, env.Client, zeroKey, env.OperatorKey, env.OperatorKey, NO_VALIDATION)
	require.NoError(t, err)
	_, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)
	// Verify keys are set to zero
	tokenInfo, err := NewTokenInfoQuery().SetTokenID(tokenID).Execute(env.Client)
	verifyLowerPrivilegeKeys(t, tokenInfo, zeroKey)

	// Update lower privilege keys to valid key list with admin key
	validKey, err := PrivateKeyGenerateEd25519()
	require.NoError(t, err)
	resp, err = updateTokenKeysHelper(t, tokenID, LOWER_PRIVILEGE, env.Client, validKey, env.OperatorKey, NO_VALIDATION)
	require.NoError(t, err)
	_, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)
	// Verify keys are set to valid key
	tokenInfo, err = NewTokenInfoQuery().SetTokenID(tokenID).Execute(env.Client)
	verifyLowerPrivilegeKeys(t, tokenInfo, validKey.PublicKey())

	// Update lower privilege keys to empty key list with admin key
	resp, err = updateTokenKeysHelper(t, tokenID, LOWER_PRIVILEGE, env.Client, NewKeyList(), env.OperatorKey, NO_VALIDATION)
	require.NoError(t, err)
	_, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)
	// Verify keys are set to empty
	tokenInfo, err = NewTokenInfoQuery().SetTokenID(tokenID).Execute(env.Client)
	verifyLowerPrivilegeKeysNil(t, tokenInfo)
}

func TestIntegrationTokenUpdateTransactionUpdateLowerPrivilegeKeysWithInvalidKeyFails(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)

	adminKey, err := GeneratePrivateKey()
	require.NoError(t, err)

	nonAdminKey, err := GeneratePrivateKey()
	require.NoError(t, err)

	someKey, err := GeneratePrivateKey()
	require.NoError(t, err)

	// create the token
	resp, tokenID, err := createTokenWithKeysAndUpdateTokenKeyHelper(t, ALL, LOWER_PRIVILEGE, env.Client, nonAdminKey, adminKey, adminKey, NO_VALIDATION)
	require.NoError(t, err)
	_, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	// Cannot remove tokens without admin key
	resp, err = updateTokenKeysHelper(t, tokenID, LOWER_PRIVILEGE, env.Client, NewKeyList(), nonAdminKey, NO_VALIDATION)
	require.NoError(t, err)
	_, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.ErrorContains(t, err, "INVALID_SIGNATURE")

	// Cannot upadate token with some random key
	resp, err = updateTokenKeysHelper(t, tokenID, LOWER_PRIVILEGE, env.Client, NewKeyList(), someKey, NO_VALIDATION)
	require.NoError(t, err)
	_, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.ErrorContains(t, err, "INVALID_SIGNATURE")

	zeroKey, err := ZeroKey()
	require.NoError(t, err)

	// Cannot upadate token with some random key
	resp, err = updateTokenKeysHelper(t, tokenID, LOWER_PRIVILEGE, env.Client, zeroKey, someKey, NO_VALIDATION)
	require.NoError(t, err)
	_, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.ErrorContains(t, err, "INVALID_SIGNATURE")
}

func TestIntegrationTokenUpdateTransactionUpdateAdminKeyWithoutAlreadySetKeyFails(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)

	someKey, err := GeneratePrivateKey()
	require.NoError(t, err)

	// Create token without keys and fail updating them
	for _, keyType := range []KeyType{WIPE_KEY, KYC_KEY, SUPPLY_KEY, FREEZE_KEY, FEE_SCHEDULE_KEY, PAUSE_KEY, METADATA_KEY, ADMIN_KEY} {
		resp, _, err := createTokenWithKeysAndUpdateTokenKeyHelper(t, NONE, keyType, env.Client, someKey, env.OperatorKey, env.OperatorKey, NO_VALIDATION)
		require.NoError(t, err)
		_, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
		require.ErrorContains(t, err, "TOKEN_IS_IMMUTABLE")
	}
}

func TestIntegrationTokenUpdateTransactionUpdateKeysLowerPrivKeysUpdateThemselvesNoValidation(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)

	zeroNewKey, err := ZeroKey()
	require.NoError(t, err)

	validNewKey, err := PrivateKeyGenerateEd25519()
	require.NoError(t, err)

	initialKey := env.OperatorKey

	// Update to valid key
	for _, keyType := range []KeyType{WIPE_KEY, KYC_KEY, SUPPLY_KEY, FREEZE_KEY, FEE_SCHEDULE_KEY, PAUSE_KEY, METADATA_KEY} {
		resp, _, err := createTokenWithKeysAndUpdateTokenKeyHelper(t, keyType, keyType, env.Client, validNewKey, initialKey, initialKey, NO_VALIDATION)
		_, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
		require.NoError(t, err)
	}

	// Zero keys
	for _, keyType := range []KeyType{WIPE_KEY, KYC_KEY, SUPPLY_KEY, FREEZE_KEY, FEE_SCHEDULE_KEY, PAUSE_KEY, METADATA_KEY} {
		resp, _, err := createTokenWithKeysAndUpdateTokenKeyHelper(t, keyType, keyType, env.Client, zeroNewKey, initialKey, initialKey, NO_VALIDATION)
		_, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
		require.NoError(t, err)
	}
}

func TestIntegrationTokenUpdateTransactionUpdateKeysLowerPrivilegeKeysUpdateThemselvesFullValidation(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)

	validNewKey, err := PrivateKeyGenerateEd25519()
	require.NoError(t, err)

	initialKey := env.OperatorKey

	// Update to valid key
	for _, keyType := range []KeyType{WIPE_KEY, KYC_KEY, SUPPLY_KEY, FREEZE_KEY, FEE_SCHEDULE_KEY, PAUSE_KEY, METADATA_KEY} {
		resp, _, err := createTokenWithKeysAndUpdateTokenKeyHelper(t, keyType, keyType, env.Client, validNewKey, initialKey, initialKey, FULL_VALIDATION)
		_, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
		require.NoError(t, err)
	}

}

func TestIntegrationTokenUpdateTransactionUpdateKeysLowerPrivilegeKeysUpdateFullValidationFails(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)

	zeroNewKey, err := ZeroKey()
	require.NoError(t, err)

	initialKey := env.OperatorKey

	// Zero out keys
	for _, keyType := range []KeyType{WIPE_KEY, KYC_KEY, SUPPLY_KEY, FREEZE_KEY, FEE_SCHEDULE_KEY, PAUSE_KEY, METADATA_KEY} {
		resp, _, err := createTokenWithKeysAndUpdateTokenKeyHelper(t, keyType, keyType, env.Client, zeroNewKey, initialKey, initialKey, FULL_VALIDATION)
		_, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
		require.ErrorContains(t, err, "exceptional receipt status: INVALID_SIGNATURE")
	}
}

func TestIntegrationTokenUpdateTransactionRemoveKeysWithoutAdminKeyFails(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)

	initialKey := env.OperatorKey

	// Cannot remove keys without admin key
	for _, keyType := range []KeyType{WIPE_KEY, KYC_KEY, SUPPLY_KEY, FREEZE_KEY, FEE_SCHEDULE_KEY, PAUSE_KEY, METADATA_KEY} {
		resp, _, err := createTokenWithKeysAndUpdateTokenKeyHelper(t, keyType, keyType, env.Client, NewKeyList(), initialKey, initialKey, NO_VALIDATION)
		_, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
		require.ErrorContains(t, err, "TOKEN_IS_IMMUTABLE")
	}
}

func TestIntegrationTokenUpdateTransactionRemoveKeysWithoutAdminKeySignFails(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)

	newKey, err := PrivateKeyGenerateEd25519()
	require.NoError(t, err)

	initialKey, err := PrivateKeyGenerateEd25519()
	require.NoError(t, err)

	wipeKey, err := PrivateKeyGenerateEd25519()
	require.NoError(t, err)

	// Create NFT with wipe and supply keys
	tokenID, err := createNft(&env, func(transaction *TokenCreateTransaction) {
		transaction.
			SetWipeKey(initialKey).
			SetKycKey(initialKey).
			SetSupplyKey(initialKey).
			SetFreezeKey(initialKey).
			SetFeeScheduleKey(initialKey).
			SetPauseKey(initialKey).
			SetMetadataKey(initialKey).
			SetAdminKey(initialKey).
			FreezeWith(env.Client)
		transaction.Sign(initialKey)
	})
	require.NoError(t, err)

	// Update supply key
	tx1, err := NewTokenUpdateTransaction().
		SetTokenID(tokenID).
		SetWipeKey(newKey).
		SetKycKey(newKey).
		SetSupplyKey(newKey).
		SetFreezeKey(newKey).
		SetFeeScheduleKey(newKey).
		SetPauseKey(newKey).
		SetMetadataKey(newKey).
		SetAdminKey(newKey).
		SetKeyVerificationMode(NO_VALIDATION).
		FreezeWith(env.Client)
	assert.NoError(t, err)

	// Cannot remove keys without admin key signing
	resp, err := tx1.Sign(wipeKey).Execute(env.Client)
	assert.NoError(t, err)

	_, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.ErrorContains(t, err, "exceptional receipt status: INVALID_SIGNATURE")
}

func TestIntegrationTokenUpdateTransactionUpdateSupplyKeyFailsWhenSignWithWipeKey(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)

	newSupplyKey, err := PrivateKeyGenerateEd25519()
	require.NoError(t, err)

	supplyKey, err := PrivateKeyGenerateEd25519()
	require.NoError(t, err)

	wipeKey, err := PrivateKeyGenerateEd25519()
	require.NoError(t, err)

	// Create NFT with wipe and supply keys
	tokenID, err := createNft(&env, func(transaction *TokenCreateTransaction) {
		transaction.
			SetAdminKey(nil).
			SetWipeKey(wipeKey).
			SetSupplyKey(supplyKey)
	})
	require.NoError(t, err)

	// Update supply key
	tx, err := NewTokenUpdateTransaction().
		SetTokenID(tokenID).
		SetSupplyKey(newSupplyKey).
		SetKeyVerificationMode(NO_VALIDATION).
		FreezeWith(env.Client)
	assert.NoError(t, err)

	// Sign with the supply key
	resp, err := tx.Sign(wipeKey).Execute(env.Client)
	assert.NoError(t, err)

	_, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.ErrorContains(t, err, "exceptional receipt status: INVALID_SIGNATURE")
}

func TestIntegrationTokenUpdateTransactionUpdateSupplyKeyToEmptyKeyAndVerifyItsImmutable(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)

	newSupplyKey, err := PrivateKeyGenerateEd25519()
	require.NoError(t, err)

	supplyKey, err := PrivateKeyGenerateEd25519()
	require.NoError(t, err)

	// Create NFT with admin and supply keys
	tokenID, err := createNft(&env, func(transaction *TokenCreateTransaction) {
		transaction.SetSupplyKey(supplyKey).
			SetAdminKey(env.OperatorKey)
	})
	require.NoError(t, err)

	// Update supply key
	tx, err := NewTokenUpdateTransaction().
		SetTokenID(tokenID).
		SetSupplyKey(NewKeyList()).
		SetKeyVerificationMode(NO_VALIDATION).
		FreezeWith(env.Client)
	assert.NoError(t, err)

	// Sign with the supply key
	resp, err := tx.Sign(supplyKey).Execute(env.Client)
	assert.NoError(t, err)

	_, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	tx, err = NewTokenUpdateTransaction().
		SetTokenID(tokenID).
		SetSupplyKey(newSupplyKey).
		SetKeyVerificationMode(NO_VALIDATION).
		FreezeWith(env.Client)
	assert.NoError(t, err)

	// Sign with the admin key
	resp, err = tx.Sign(env.OperatorKey).Execute(env.Client)
	assert.NoError(t, err)

	_, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.ErrorContains(t, err, "exceptional receipt status: TOKEN_HAS_NO_SUPPLY_KEY")
}

func TestIntegrationTokenUpdateTransactionUpdateSupplyKeyFullValidationFails(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)

	newSupplyKey, err := PrivateKeyGenerateEd25519()
	require.NoError(t, err)

	supplyKey, err := PrivateKeyGenerateEd25519()
	require.NoError(t, err)

	// Create NFT with supply key
	tokenID, err := createNft(&env, func(transaction *TokenCreateTransaction) {
		transaction.
			SetSupplyKey(supplyKey).
			SetAdminKey(nil)
	})
	require.NoError(t, err)

	// Update supply key
	tx, err := NewTokenUpdateTransaction().
		SetTokenID(tokenID).
		SetSupplyKey(newSupplyKey).
		SetKeyVerificationMode(FULL_VALIDATION).
		FreezeWith(env.Client)
	assert.NoError(t, err)

	_, err = supplyKey.SignTransaction(tx)
	assert.NoError(t, err)

	// Sign with the old supply key, should fail
	resp, err := tx.Execute(env.Client)
	assert.NoError(t, err)

	_, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.ErrorContains(t, err, "exceptional receipt status: INVALID_SIGNATURE")

	// Update supply key
	tx2, err := NewTokenUpdateTransaction().
		SetTokenID(tokenID).
		SetSupplyKey(newSupplyKey).
		SetKeyVerificationMode(FULL_VALIDATION).
		FreezeWith(env.Client)
	assert.NoError(t, err)

	// Sign with only the new supply key, should fail
	_, err = newSupplyKey.SignTransaction(tx2)
	assert.NoError(t, err)

	resp, err = tx2.Execute(env.Client)
	assert.NoError(t, err)

	_, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.ErrorContains(t, err, "exceptional receipt status: INVALID_SIGNATURE")
}

func TestIntegrationTokenUpdateTransactionUpdateSupplyKeyWithInvalidKey(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)
	var invalidKey _Ed25519PublicKey
	randomBytes := make([]byte, 32)
	keyData := [32]byte{
		0x01, 0x23, 0x45, 0x67,
		0x89, 0xab, 0xcd, 0xef,
		0xfe, 0xdc, 0xba, 0x98,
		0x76, 0x54, 0x32, 0x10,
		0x00, 0x11, 0x22, 0x33,
		0x44, 0x55, 0x66, 0x77,
		0x88, 0x99, 0xaa, 0xbb,
		0xcc, 0xdd, 0xee, 0xff,
	}
	randomBytes = keyData[:]
	copy(invalidKey.keyData[:], randomBytes)
	invalidSupplyKey := PublicKey{
		ed25519PublicKey: &invalidKey,
	}

	supplyKey, err := PrivateKeyGenerateEd25519()
	require.NoError(t, err)

	// Create NFT with supplyKey
	tokenID, err := createNft(&env, func(transaction *TokenCreateTransaction) {
		transaction.SetSupplyKey(supplyKey)
	})
	require.NoError(t, err)

	// Update supply key
	tx, err := NewTokenUpdateTransaction().
		SetTokenID(tokenID).
		SetSupplyKey(invalidSupplyKey).
		SetKeyVerificationMode(NO_VALIDATION).
		FreezeWith(env.Client)
	assert.NoError(t, err)

	_, err = supplyKey.SignTransaction(tx)
	assert.NoError(t, err)

	//Sign with the old supply key
	_, err = tx.Execute(env.Client)
	assert.ErrorContains(t, err, "exceptional precheck status INVALID_SUPPLY_KEY")

}

func createTokenWithKeysAndUpdateTokenKeyHelper(t *testing.T, createKeyType KeyType, updateKeyType KeyType, client *Client, newKey Key, initialKey Key, signerKey PrivateKey, verificationMode TokenKeyValidation) (TransactionResponse, TokenID, error) {
	// Create Fungible token with keys
	tx := NewTokenCreateTransaction().
		SetTokenName("ffff").
		SetTokenSymbol("F").
		SetTreasuryAccountID(client.GetOperatorAccountID()).
		SetFreezeDefault(false)

	switch createKeyType {
	case WIPE_KEY:
		tx.SetWipeKey(initialKey)
	case KYC_KEY:
		tx.SetKycKey(initialKey)
	case SUPPLY_KEY:
		tx.SetSupplyKey(initialKey)
	case FREEZE_KEY:
		tx.SetFreezeKey(initialKey)
	case FEE_SCHEDULE_KEY:
		tx.SetFeeScheduleKey(initialKey)
	case PAUSE_KEY:
		tx.SetPauseKey(initialKey)
	case METADATA_KEY:
		tx.SetMetadataKey(initialKey)
	case ADMIN_KEY:
		tx.SetAdminKey(initialKey)
	case LOWER_PRIVILEGE:
		tx.SetWipeKey(initialKey).
			SetKycKey(initialKey).
			SetSupplyKey(initialKey).
			SetFreezeKey(initialKey).
			SetFeeScheduleKey(initialKey).
			SetPauseKey(initialKey).
			SetMetadataKey(initialKey)
	case ALL:
		tx.SetWipeKey(initialKey).
			SetKycKey(initialKey).
			SetSupplyKey(initialKey).
			SetFreezeKey(initialKey).
			SetFeeScheduleKey(initialKey).
			SetPauseKey(initialKey).
			SetMetadataKey(initialKey).
			SetAdminKey(initialKey)
	}
	frozenTx, err := tx.FreezeWith(client)
	require.NoError(t, err)
	resp, err := frozenTx.Sign(signerKey).Execute(client)
	require.NoError(t, err)
	receipt, err := resp.SetValidateStatus(true).GetReceipt(client)
	require.NoError(t, err)

	resp, err = updateTokenKeysHelper(t, *receipt.TokenID, updateKeyType, client, newKey, signerKey, verificationMode)
	return resp, *receipt.TokenID, err
}

func updateTokenKeysHelper(t *testing.T, tokenID TokenID, updateKeyType KeyType, client *Client, newKey Key, signerKey PrivateKey, verificationMode TokenKeyValidation) (TransactionResponse, error) {
	privateKey, _ := newKey.(PrivateKey)
	// Update the key
	tx := NewTokenUpdateTransaction().
		SetTokenID(tokenID).
		SetKeyVerificationMode(verificationMode)

	switch updateKeyType {
	case WIPE_KEY:
		tx.SetWipeKey(newKey)
	case KYC_KEY:
		tx.SetKycKey(newKey)
	case SUPPLY_KEY:
		tx.SetSupplyKey(newKey)
	case FREEZE_KEY:
		tx.SetFreezeKey(newKey)
	case FEE_SCHEDULE_KEY:
		tx.SetFeeScheduleKey(newKey)
	case PAUSE_KEY:
		tx.SetPauseKey(newKey)
	case METADATA_KEY:
		tx.SetMetadataKey(newKey)
	case ADMIN_KEY:
		tx.SetAdminKey(newKey)
	case LOWER_PRIVILEGE:
		tx.SetWipeKey(newKey).
			SetKycKey(newKey).
			SetSupplyKey(newKey).
			SetFreezeKey(newKey).
			SetFeeScheduleKey(newKey).
			SetPauseKey(newKey).
			SetMetadataKey(newKey)
	case ALL:
		tx.SetWipeKey(newKey).
			SetKycKey(newKey).
			SetSupplyKey(newKey).
			SetFreezeKey(newKey).
			SetFeeScheduleKey(newKey).
			SetPauseKey(newKey).
			SetMetadataKey(newKey).
			SetAdminKey(newKey)
	}
	frozenTx, err := tx.FreezeWith(client)
	assert.NoError(t, err)

	if updateKeyType == ADMIN_KEY || updateKeyType == ALL || verificationMode == FULL_VALIDATION {
		privateKey.SignTransaction(tx)
		signerKey.SignTransaction(tx)
		resp, err := frozenTx.Execute(client)
		return resp, err

	} else {
		// Sign with the signer key
		resp, err := frozenTx.Sign(signerKey).Execute(client)
		return resp, err
	}
}

func verifyAdminKey(t *testing.T, tokenInfo TokenInfo, expectedKey Key) {
	assert.Equalf(t, expectedKey.String(), tokenInfo.AdminKey.String(), " admin key did not match")
}

func verifyAdminKeyNil(t *testing.T, tokenInfo TokenInfo) {
	assert.Nil(t, tokenInfo.AdminKey)
}

func verifyLowerPrivilegeKeysNil(t *testing.T, tokenInfo TokenInfo) {
	assert.Nil(t, tokenInfo.WipeKey)
	assert.Nil(t, tokenInfo.KycKey)
	assert.Nil(t, tokenInfo.SupplyKey)
	assert.Nil(t, tokenInfo.FreezeKey)
	assert.Nil(t, tokenInfo.FeeScheduleKey)
	assert.Nil(t, tokenInfo.PauseKey)
	assert.Nil(t, tokenInfo.MetadataKey)
}

func verifyLowerPrivilegeKeys(t *testing.T, tokenInfo TokenInfo, expectedKey Key) {
	assert.Equalf(t, expectedKey.String(), tokenInfo.WipeKey.String(), "wipe key did not match")
	assert.Equalf(t, expectedKey.String(), tokenInfo.KycKey.String(), "kyc key did not match")
	assert.Equalf(t, expectedKey.String(), tokenInfo.SupplyKey.String(), "supply key did not match")
	assert.Equalf(t, expectedKey.String(), tokenInfo.FreezeKey.String(), "freeze key did not match")
	assert.Equalf(t, expectedKey.String(), tokenInfo.FeeScheduleKey.String(), "fee schedule key did not match")
	assert.Equalf(t, expectedKey.String(), tokenInfo.PauseKey.String(), "pause key did not match")
	assert.Equalf(t, expectedKey.String(), tokenInfo.MetadataKey.String(), "metadata key did not match")
}
