package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"crypto/aes"
	cipher2 "crypto/cipher"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"crypto/sha512"
	"crypto/subtle"
	"encoding/hex"
	"encoding/json"

	"io"

	"golang.org/x/crypto/pbkdf2"
)

type _Keystore struct {
	Version uint8       `json:"version"`
	Crypto  _CryptoData `json:"crypto"`
}

// internal struct used for cipher parameters
type _CipherParams struct {
	// hex-encoded initialization vector
	IV string `json:"iv"`
}

// internal struct used for kdf parameters
type _KdfParams struct {
	// derived key length
	DKLength int `json:"dklength"`
	// hex-encoded salt
	Salt string `json:"salt"`
	// iteration count
	Count int `json:"c"`
	// hash function
	PRF string `json:"prf"`
}

// internal type used in _Keystore to represent the crypto data
type _CryptoData struct {
	// hex-encoded ciphertext
	CipherText   string        `json:"ciphertext"`
	CipherParams _CipherParams `json:"cipherparams"`
	// Cipher being used
	Cipher string `json:"cipher"`
	// key derivation function being used
	KDF string `json:"kdf"`
	// parameters for key derivation function
	KDFParams _KdfParams `json:"kdfparams"`
	// hex-encoded HMAC-SHA384
	Mac string `json:"mac"`
}

const Aes128Ctr = "aes-128-ctr"
const HmacSha256 = "hmac-sha256"

// all values taken from https://github.com/ethereumjs/ethereumjs-wallet/blob/de3a92e752673ada1d78f95cf80bc56ae1f59775/src/index.ts#L25
const dkLen int = 32
const c int = 262144
const saltLen uint = 32

func _RandomBytes(n uint) ([]byte, error) {
	// based on https://github.com/gophercon/2016-talks/tree/master/GeorgeTankersley-CryptoForGoDevelopers
	b := make([]byte, n)
	_, err := io.ReadFull(rand.Reader, b)
	if err != nil {
		return nil, err
	}

	return b, nil
}

func _NewKeystore(privateKey []byte, passphrase string) ([]byte, error) {
	salt, err := _RandomBytes(saltLen)
	if err != nil {
		return nil, err
	}

	key := pbkdf2.Key([]byte(passphrase), salt, c, dkLen, sha256.New)

	iv, err := _RandomBytes(16)
	if err != nil {
		return nil, err
	}

	// AES-128-CTR with the first half of the derived key and a random IV
	block, err := aes.NewCipher(key[0:16])
	if err != nil {
		return nil, err
	}

	cipher := cipher2.NewCTR(block, iv)
	cipherText := make([]byte, len(privateKey))
	cipher.XORKeyStream(cipherText, privateKey)

	h := hmac.New(sha512.New384, key[16:])

	if _, err = h.Write(cipherText); err != nil {
		return nil, err
	}

	mac := h.Sum(nil)

	keystore := _Keystore{
		Version: 1,
		Crypto: _CryptoData{
			CipherText: hex.EncodeToString(cipherText),
			CipherParams: _CipherParams{
				IV: hex.EncodeToString(iv),
			},
			Cipher: Aes128Ctr,
			KDF:    "pbkdf2",
			KDFParams: _KdfParams{
				DKLength: dkLen,
				Salt:     hex.EncodeToString(salt),
				Count:    c,
				PRF:      HmacSha256,
			},
			Mac: hex.EncodeToString(mac),
		},
	}

	return json.Marshal(keystore)
}

func _ParseKeystore(keystoreBytes []byte, passphrase string) (PrivateKey, error) {
	keyStore := _Keystore{}

	err := json.Unmarshal(keystoreBytes, &keyStore)

	if err != nil {
		return PrivateKey{}, err
	}

	if keyStore.Version != 1 {
		// todo: change to a switch and handle differently if future _Keystore versions are added
		return PrivateKey{}, _NewErrBadKeyf("unsupported _Keystore version: %v", keyStore.Version)
	}

	if keyStore.Crypto.KDF != "pbkdf2" {
		return PrivateKey{}, _NewErrBadKeyf("unsupported KDF: %v", keyStore.Crypto.KDF)
	}

	if keyStore.Crypto.Cipher != Aes128Ctr {
		return PrivateKey{}, _NewErrBadKeyf("unsupported _Keystore cipher: %v", keyStore.Crypto.Cipher)
	}

	if keyStore.Crypto.KDFParams.PRF != HmacSha256 {
		return PrivateKey{}, _NewErrBadKeyf(
			"unsupported PRF: %v",
			keyStore.Crypto.KDFParams.PRF)
	}

	salt, err := hex.DecodeString(keyStore.Crypto.KDFParams.Salt)

	if err != nil {
		return PrivateKey{}, err
	}

	iv, err := hex.DecodeString(keyStore.Crypto.CipherParams.IV)

	if err != nil {
		return PrivateKey{}, err
	}

	cipherBytes, err := hex.DecodeString(keyStore.Crypto.CipherText)

	if err != nil {
		return PrivateKey{}, err
	}

	key := pbkdf2.Key([]byte(passphrase), salt, keyStore.Crypto.KDFParams.Count, dkLen, sha256.New)

	mac, err := hex.DecodeString(keyStore.Crypto.Mac)

	if err != nil {
		return PrivateKey{}, err
	}

	h := hmac.New(sha512.New384, key[16:])

	_, err = h.Write(cipherBytes)

	if err != nil {
		return PrivateKey{}, err
	}

	verifyMac := h.Sum(nil)

	if subtle.ConstantTimeCompare(mac, verifyMac) == 0 {
		return PrivateKey{}, _NewErrBadKeyf("hmac mismatch; passphrase is incorrect")
	}

	block, err := aes.NewCipher(key[:16])
	if err != nil {
		return PrivateKey{}, err
	}

	decipher := cipher2.NewCTR(block, iv)
	pkBytes := make([]byte, len(cipherBytes))

	decipher.XORKeyStream(pkBytes, cipherBytes)

	return PrivateKeyFromBytesEd25519(pkBytes)
}
