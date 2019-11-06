package hedera

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

const testKeyStr = "302e020100300506032b657004220420db484b828e64b2d8f12ce3c0a0e93a0b8cce7af1bb8f39c97732394482538e10"

func TestEd25519PrivateKeyGenerate(t *testing.T) {
	key, err := GenerateEd25519PrivateKey()

	assert.NoError(t, err)
	assert.True(t, strings.HasPrefix(key.String(), ed25519PrivKeyPrefix))
}

func TestEd25519PrivateKeyExternalSerialization(t *testing.T) {
	key, err := Ed25519PrivateKeyFromString(testKeyStr)

	assert.NoError(t, err)
	assert.Equal(t, testKeyStr, key.String())
}

func TestEd25519PrivateKeyExternalSerializationForConcatenatedHex(t *testing.T) {
	keyStr := "db484b828e64b2d8f12ce3c0a0e93a0b8cce7af1bb8f39c97732394482538e10e0c8ec2758a5879ffac226a13c0c516b799e72e35141a0dd828f94d37988a4b7"
	key, err := Ed25519PrivateKeyFromString(keyStr)

	assert.NoError(t, err)
	assert.Equal(t, testKeyStr, key.String())
}

func TestEd25519PrivateKeyExternalSerializationForRawHex(t *testing.T) {
	keyStr := "db484b828e64b2d8f12ce3c0a0e93a0b8cce7af1bb8f39c97732394482538e10"
	key, err := Ed25519PrivateKeyFromString(keyStr)

	assert.NoError(t, err)
	assert.Equal(t, testKeyStr, key.String())
}
