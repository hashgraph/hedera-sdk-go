package hedera

/*-
 *
 * Hedera Go SDK
 *
 * Copyright (C) 2020 - 2022 Hedera Hashgraph, LLC
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

import (
	"crypto/ecdsa"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha512"
	"crypto/x509"
	"encoding/asn1"
	"encoding/hex"
	"encoding/pem"
	"fmt"
	"io"
	"strings"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/hashgraph/hedera-protobufs-go/services"
	"github.com/pkg/errors"
)

// _ECDSAPrivateKey is an Key_ECDSASecp256K1 private key.
type _ECDSAPrivateKey struct {
	keyData   *ecdsa.PrivateKey
	chainCode []byte
}

const _LegacyECDSAPrivateKeyPrefix = "3030020100300706052b8104000a04220420"

func _GenerateECDSAPrivateKey() (*_ECDSAPrivateKey, error) {
	key, err := ecdsa.GenerateKey(crypto.S256(), rand.Reader)
	if err != nil {
		return &_ECDSAPrivateKey{}, err
	}

	return &_ECDSAPrivateKey{
		keyData: key,
	}, nil
}

func _ECDSAPrivateKeyFromBytes(byt []byte) (*_ECDSAPrivateKey, error) {
	length := len(byt)
	switch {
	case length == 32:
		return _ECDSAPrivateKeyFromBytesRaw(byt)
	case length > 32:
		return _ECDSAPrivateKeyFromBytesDer(byt)
	default:
		return &_ECDSAPrivateKey{}, _NewErrBadKeyf("invalid private key length: %v bytes", len(byt))
	}
}

func _ECDSAPrivateKeyFromBytesRaw(byt []byte) (*_ECDSAPrivateKey, error) {
	length := len(byt)
	if length != 32 {
		return &_ECDSAPrivateKey{}, _NewErrBadKeyf("invalid private key length: %v bytes", len(byt))
	}

	key, err := crypto.ToECDSA(byt)
	if err != nil {
		return nil, err
	}

	return &_ECDSAPrivateKey{
		keyData: key,
	}, nil
}

func _LegacyECDSAPrivateKeyFromBytesDer(byt []byte) (*_ECDSAPrivateKey, error) {
	given := hex.EncodeToString(byt)

	result := strings.ReplaceAll(given, _LegacyECDSAPrivateKeyPrefix, "")
	decoded, err := hex.DecodeString(result)
	if err != nil {
		return &_ECDSAPrivateKey{}, err
	}

	if len(decoded) != 32 {
		return &_ECDSAPrivateKey{}, _NewErrBadKeyf("invalid private key length: %v bytes", len(byt))
	}

	key, err := crypto.ToECDSA(decoded)
	if err != nil {
		return nil, err
	}

	return &_ECDSAPrivateKey{
		keyData: key,
	}, nil
}

func _ECDSAPrivateKeyFromBytesDer(data []byte) (*_ECDSAPrivateKey, error) {
	given := hex.EncodeToString(data)
	if strings.HasPrefix(given, _LegacyECDSAPrivateKeyPrefix) {
		return _LegacyECDSAPrivateKeyFromBytesDer(data)
	}

	type ECPrivateKey struct {
		Version       int
		PrivateKey    []byte
		NamedCurveOID asn1.ObjectIdentifier `asn1:"optional,explicit,tag:0"`
		PublicKey     asn1.BitString        `asn1:"optional,explicit,tag:1"`
	}
	var ecKey ECPrivateKey
	if rest, err := asn1.Unmarshal(data, &ecKey); err != nil {
		return nil, err
	} else if len(rest) != 0 {
		return nil, errors.New("x509: trailing data after ASN.1 of public-key")
	}
	key, err := crypto.ToECDSA(ecKey.PrivateKey)
	if err != nil {
		return nil, err
	}
	return &_ECDSAPrivateKey{keyData: key}, nil
}

func _ECDSAPrivateKeyFromSeed(seed []byte) (*_ECDSAPrivateKey, error) {
	h := hmac.New(sha512.New, []byte("Bitcoin seed"))

	_, err := h.Write(seed)
	if err != nil {
		return &_ECDSAPrivateKey{}, err
	}

	digest := h.Sum(nil)

	keyBytes := digest[0:32]
	chainCode := digest[32:]
	privateKey, err := _ECDSAPrivateKeyFromBytes(keyBytes)

	if err != nil {
		return &_ECDSAPrivateKey{}, err
	}
	privateKey.chainCode = chainCode

	return privateKey, nil
}

func _ECDSAPrivateKeyFromString(s string) (*_ECDSAPrivateKey, error) {
	b, err := hex.DecodeString(strings.ToLower(s))
	if err != nil {
		return &_ECDSAPrivateKey{}, err
	}

	return _ECDSAPrivateKeyFromBytes(b)
}

func (sk *_ECDSAPrivateKey) _PublicKey() *_ECDSAPublicKey {
	if sk.keyData.Y == nil && sk.keyData.X == nil {
		b := sk.keyData.D.Bytes()
		x, y := crypto.S256().ScalarBaseMult(b)
		sk.keyData.X = x
		sk.keyData.Y = y
		return &_ECDSAPublicKey{
			&ecdsa.PublicKey{
				Curve: crypto.S256(),
				X:     x,
				Y:     y,
			},
		}
	}

	return &_ECDSAPublicKey{
		&ecdsa.PublicKey{
			Curve: sk.keyData.Curve,
			X:     sk.keyData.X,
			Y:     sk.keyData.Y,
		},
	}
}

func _ECDSAPrivateKeyFromPem(bytes []byte, passphrase string) (*_ECDSAPrivateKey, error) {
	block, _ := pem.Decode(bytes)
	if block == nil {
		return nil, fmt.Errorf("failed to decode PEM block containing private key")
	}

	if block.Type != "EC PRIVATE KEY" {
		return nil, fmt.Errorf("unsupported PEM block type: %s", block.Type)
	}
	//nolint
	if x509.IsEncryptedPEMBlock(block) {
		der, err := x509.DecryptPEMBlock(block, []byte(passphrase))
		if err != nil {
			return nil, err
		}
		block.Bytes = der
	}

	key, err := _ECDSAPrivateKeyFromBytes(block.Bytes)
	if err != nil {
		return nil, err
	}
	return key, nil
}

func _ECDSAPrivateKeyReadPem(source io.Reader, passphrase string) (*_ECDSAPrivateKey, error) {
	// note: Passphrases are currently not supported, but included in the function definition to avoid breaking
	// changes in the future.

	pemFileBytes, err := io.ReadAll(source)
	if err != nil {
		return &_ECDSAPrivateKey{}, err
	}

	return _ECDSAPrivateKeyFromPem(pemFileBytes, passphrase)
}

func (sk _ECDSAPrivateKey) _Sign(message []byte) []byte {
	hash := crypto.Keccak256Hash(message)
	sig, err := crypto.Sign(hash.Bytes(), sk.keyData)
	if err != nil {
		panic(err)
	}

	// signature returned has a ecdsa recovery byte at the end,
	// need to remove it for verification to work.
	return sig[:len(sig)-1]
}

// SupportsDerivation returns true if the _ECDSAPrivateKey supports derivation.
func (sk _ECDSAPrivateKey) _SupportsDerivation() bool {
	return sk.chainCode != nil
}

// Derive a child key compatible with the iOS and Android wallets using a provided wallet/account index. Use index 0 for
// the default account.
//
// This will fail if the key does not support derivation which can be checked by calling SupportsDerivation()
func (sk _ECDSAPrivateKey) _Derive(index uint32) (*_ECDSAPrivateKey, error) {
	if !sk._SupportsDerivation() {
		return nil, _NewErrBadKeyf("child key cannot be derived from this key")
	}

	derivedKeyBytes, chainCode, err := _DeriveECDSAChildKey(sk._BytesRaw(), sk.chainCode, index)
	if err != nil {
		return nil, err
	}

	derivedKey, err := _ECDSAPrivateKeyFromBytes(derivedKeyBytes)
	if err != nil {
		return nil, err
	}

	derivedKey.chainCode = chainCode

	return derivedKey, nil
}

func (sk _ECDSAPrivateKey) _BytesRaw() []byte {
	privateKey := make([]byte, 32)
	temp := sk.keyData.D.Bytes()
	copy(privateKey[32-len(temp):], temp)

	return privateKey
}

func (sk _ECDSAPrivateKey) _BytesDer() []byte {
	type ECPrivateKey struct {
		Version       int
		PrivateKey    []byte
		NamedCurveOID asn1.ObjectIdentifier `asn1:"optional,explicit,tag:0"`
		PublicKey     asn1.BitString        `asn1:"optional,explicit,tag:1"`
	}

	secp256k1OID := asn1.ObjectIdentifier{1, 3, 132, 0, 10}

	ecPrivateKey := ECPrivateKey{
		Version:       1, // EC private keys have a version of 1
		PrivateKey:    sk._BytesRaw(),
		NamedCurveOID: secp256k1OID,
		PublicKey:     asn1.BitString{Bytes: sk._PublicKey()._BytesRaw()},
	}

	derBytes, err := asn1.Marshal(ecPrivateKey)
	if err != nil {
		return nil
	}

	return derBytes
}

func (sk _ECDSAPrivateKey) _StringDer() string {
	return fmt.Sprint(hex.EncodeToString(sk._BytesDer()))
}

func (sk _ECDSAPrivateKey) _StringRaw() string {
	return fmt.Sprint(hex.EncodeToString(sk._BytesRaw()))
}
func (sk _ECDSAPrivateKey) StringRaw2() string {
	return fmt.Sprint(hex.EncodeToString(sk._BytesRaw()))
}
func (sk _ECDSAPrivateKey) _ToProtoKey() *services.Key {
	return sk._PublicKey()._ToProtoKey()
}

func (sk _ECDSAPrivateKey) _SignTransaction(transaction *Transaction) ([]byte, error) {
	transaction._RequireOneNodeAccountID()

	if transaction.signedTransactions._Length() == 0 {
		return make([]byte, 0), errTransactionRequiresSingleNodeAccountID
	}

	signature := sk._Sign(transaction.signedTransactions._Get(0).(*services.SignedTransaction).GetBodyBytes())

	publicKey := sk._PublicKey()
	if publicKey == nil {
		return []byte{}, errors.New("public key is nil")
	}

	wrappedPublicKey := PublicKey{
		ecdsaPublicKey: publicKey,
	}

	if transaction._KeyAlreadySigned(wrappedPublicKey) {
		return []byte{}, nil
	}

	transaction.transactions = _NewLockableSlice()
	transaction.publicKeys = append(transaction.publicKeys, wrappedPublicKey)
	transaction.transactionSigners = append(transaction.transactionSigners, nil)
	transaction.transactionIDs.locked = true

	for index := 0; index < transaction.signedTransactions._Length(); index++ {
		temp := transaction.signedTransactions._Get(index).(*services.SignedTransaction)

		temp.SigMap.SigPair = append(
			temp.SigMap.SigPair,
			publicKey._ToSignaturePairProtobuf(signature),
		)
		transaction.signedTransactions._Set(index, temp)
	}

	return signature, nil
}
