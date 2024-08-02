package methods

import (
	"context"
	"testing"

	"github.com/hashgraph/hedera-sdk-go/tck/param"
	"github.com/stretchr/testify/assert"
)

func TestGenerateKeyWithInvalidFromKey(t *testing.T) {
	// Given
	params := param.KeyParams{
		Type:      param.ED25519_PRIVATE_KEY,
		FromKey:   "someKey",
		Threshold: 0,
		Keys:      nil,
	}

	// When
	_, err := GenerateKey(context.Background(), params)

	// Then
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid parameters: fromKey should only be provided for ed25519PublicKey, ecdsaSecp256k1PublicKey, or evmAddress types")
}

func TestGenerateKeyWithInvalidThreshold(t *testing.T) {
	// Given
	params := param.KeyParams{
		Type:      param.ED25519_PUBLIC_KEY,
		FromKey:   "",
		Threshold: 1,
		Keys:      nil,
	}

	// When
	_, err := GenerateKey(context.Background(), params)

	// Then
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid parameters: threshold should only be provided for thresholdKey types")
}

func TestGenerateKeyWithInvalidKeys(t *testing.T) {
	// Given
	params := param.KeyParams{
		Type:      param.ED25519_PUBLIC_KEY,
		FromKey:   "",
		Threshold: 0,
		Keys:      []param.KeyParams{},
	}

	// When
	_, err := GenerateKey(context.Background(), params)

	// Then
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid parameters: keys should only be provided for keyList or thresholdKey types")
}

func TestGenerateKeyWithMissingKeysForKeyList(t *testing.T) {
	// Given
	params := param.KeyParams{
		Type:      param.LIST_KEY,
		FromKey:   "",
		Threshold: 0,
		Keys:      nil,
	}

	// When
	_, err := GenerateKey(context.Background(), params)

	// Then
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid request: keys list is required for generating a KeyList type")
}

func TestGenerateKeyWithMissingThresholdForThresholdKey(t *testing.T) {
	// Given
	params := param.KeyParams{
		Type: param.THRESHOLD_KEY,
		Keys: []param.KeyParams{
			{
				Type: param.ED25519_PUBLIC_KEY,
			},
		},
	}

	// When
	_, err := GenerateKey(context.Background(), params)

	// Then
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid request: threshold is required for generating a ThresholdKey type")
}

func TestGenerateKeyWithValidEd25519PrivateKey(t *testing.T) {
	// Given
	params := param.KeyParams{
		Type: param.ED25519_PRIVATE_KEY,
	}

	// When
	response, err := GenerateKey(context.Background(), params)

	// Then
	assert.NoError(t, err)
	assert.NotNil(t, response.Key)
	assert.Contains(t, response.Key, "302e020100300506032b657004220420")
}

func TestGenerateKeyWithValidEd25519PublicKey(t *testing.T) {
	// Given
	params := param.KeyParams{
		Type: param.ED25519_PUBLIC_KEY,
	}

	// When
	response, err := GenerateKey(context.Background(), params)

	// Then
	assert.NoError(t, err)
	assert.NotNil(t, response.Key)
	assert.Contains(t, response.Key, "302a300506032b6570032100")
}

func TestGenerateKeyWithValidThresholdKey(t *testing.T) {
	// Given
	params := param.KeyParams{
		Type:      param.THRESHOLD_KEY,
		Threshold: 2,
		Keys: []param.KeyParams{
			{
				Type: param.ED25519_PUBLIC_KEY,
			},
		},
	}

	// When
	response, err := GenerateKey(context.Background(), params)

	// Then
	assert.NoError(t, err)
	assert.NotNil(t, response.Key)
	assert.NotEmpty(t, response.PrivateKeys)
}

func TestGenerateKeyWithValidListKey(t *testing.T) {
	// Given
	params := param.KeyParams{
		Type: param.LIST_KEY,
		Keys: []param.KeyParams{
			{
				Type: param.ED25519_PUBLIC_KEY,
			},
		},
	}

	// When
	response, err := GenerateKey(context.Background(), params)

	// Then
	assert.NoError(t, err)
	assert.NotNil(t, response.Key)
	assert.NotEmpty(t, response.PrivateKeys)
}

func TestGenerateKeyWithValidEvmAddressKey(t *testing.T) {
	// Given
	params := param.KeyParams{
		Type:    param.EVM_ADDRESS_KEY,
		FromKey: "3054020101042056b071002a75ab207a44bb2c18320286062bc26969fcb98240301e4afbe9ee2ea00706052b8104000aa124032200038ef0b62d60b1415f8cfb460303c498fbf09cb2ef2d2ff19fad33982228ef86fd",
	}

	// When
	response, err := GenerateKey(context.Background(), params)

	// Then
	assert.NoError(t, err)
	assert.NotNil(t, response.Key)
}
