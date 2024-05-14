package hedera

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTokenRelationshipFromJson(t *testing.T) {
	// Mock token JSON object
	tokenJSON := map[string]interface{}{
		"token_id":              "0.0.12345",
		"balance":               100.0,
		"kyc_status":            "GRANTED",
		"freeze_status":         "FROZEN",
		"decimals":              8.0,
		"automatic_association": true,
	}

	// Call the function
	tokenRelationship, err := TokenRelationshipFromJson(tokenJSON)

	// Assert that there's no error
	assert.NoError(t, err)

	// Assert that the returned token relationship is correct
	assert.Equal(t, "0.0.12345", tokenRelationship.TokenID.String())
	assert.Equal(t, 100.0, tokenRelationship.Balance)
	assert.NotNil(t, tokenRelationship.KycStatus)
	assert.True(t, *tokenRelationship.KycStatus)
	assert.NotNil(t, tokenRelationship.FreezeStatus)
	assert.True(t, *tokenRelationship.FreezeStatus)
	assert.Equal(t, 8.0, tokenRelationship.Decimals)
	assert.True(t, tokenRelationship.AutomaticAssociation)
}

func TestTokenRelationshipFromJsonInvalidTokenObject(t *testing.T) {
	// Mock invalid token JSON object
	tokenJSON := "invalid_token_object"

	// Call the function
	_, err := TokenRelationshipFromJson(tokenJSON)

	// Assert that an error is returned
	assert.Error(t, err)
}

func TestTokenRelationshipFromJsonInvalidTokenID(t *testing.T) {
	// Mock token JSON object with invalid token ID
	tokenJSON := map[string]interface{}{
		"token_id":              "invalid_token_id",
		"balance":               100.0,
		"kyc_status":            "GRANTED",
		"freeze_status":         "FROZEN",
		"decimals":              8.0,
		"automatic_association": true,
	}

	// Call the function
	_, err := TokenRelationshipFromJson(tokenJSON)

	// Assert that an error is returned
	assert.Error(t, err)
}

func TestTokenRelationshipFromJsonInvalidKycStatus(t *testing.T) {
	// Mock token JSON object with invalid KYC status
	tokenJSON := map[string]interface{}{
		"token_id":              "0.0.12345",
		"balance":               100.0,
		"kyc_status":            "INVALID_KYC_STATUS",
		"freeze_status":         "FROZEN",
		"decimals":              8.0,
		"automatic_association": true,
	}

	// Call the function
	tokenRelationship, err := TokenRelationshipFromJson(tokenJSON)

	// Assert that an error is not returned
	assert.NoError(t, err)
	assert.Nil(t, tokenRelationship.KycStatus)
}

func TestTokenRelationshipFromJsonInvalidFreezeStatus(t *testing.T) {
	// Mock token JSON object with invalid freeze status
	tokenJSON := map[string]interface{}{
		"token_id":              "0.0.12345",
		"balance":               100.0,
		"kyc_status":            "GRANTED",
		"freeze_status":         "INVALID_FREEZE_STATUS",
		"decimals":              8.0,
		"automatic_association": true,
	}

	// Call the function
	tokenRelationship, err := TokenRelationshipFromJson(tokenJSON)

	// Assert that an error is not returned
	assert.NoError(t, err)
	assert.Nil(t, tokenRelationship.FreezeStatus)
}
