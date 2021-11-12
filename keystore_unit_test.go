//+build all unit

package hedera

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const passphrase string = "HelloHashgraph!"

const testKeystoreKeyString string = "45e512c479dc40bc47561c507daa962aa17b6133f8433c459adb481e06cbdafb"

// generated with the JS SDK
const testKeystore string = `{"version":1,"crypto":{"ciphertext":"9dfa728ba59e50d745f76a5fb4cab6918d4096869ab436a879731095f621456f5add51715ce81383a7996f3a75359e8216102285238a3ad8fac4dfa17894a6aa","cipherparams":{"iv":"c3198c3529fef9c5e2886f19c479683e"},"cipher":"aes-128-ctr","kdf":"pbkdf2","kdfparams":{"dkLen":32,"salt":"c87aba46b7db247694763cff3f2ec18bad1006590c6bb9befc14f05b2b2af479","c":262144,"prf":"hmac-sha256"},"mac":"f6f7a1552b3618209073feebe0109f57a3df57d8b11f07ff44aad82bafffbd15636263b96cfd2328b122d2851771c7b4"}}`

func TestUnitDecryptKeyStore(t *testing.T) {
	privateKey, err := PrivateKeyFromString(testKeystoreKeyString)
	assert.NoError(t, err)

	ksPrivateKey, err := _ParseKeystore([]byte(testKeystore), passphrase)
	assert.NoError(t, err)

	assert.Equal(t, privateKey.keyData, ksPrivateKey.keyData)
}

func TestUnitEncryptAndDecryptKeyStore(t *testing.T) {
	privateKey, err := PrivateKeyFromString(testPrivateKeyStr)
	assert.NoError(t, err)

	keyStore, err := _NewKeystore(privateKey.Bytes(), passphrase)
	assert.NoError(t, err)

	ksPrivateKey, err := _ParseKeystore(keyStore, passphrase)
	assert.NoError(t, err)

	assert.Equal(t, privateKey.keyData, ksPrivateKey.keyData)
}
