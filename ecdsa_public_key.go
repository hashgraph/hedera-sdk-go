package hedera

/*-
 *
 * Hedera Go SDK
 *
 * Copyright (C) 2020 - 2023 Hedera Hashgraph, LLC
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
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"encoding/asn1"
	"encoding/hex"
	"math/big"
	"strings"

	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/hashgraph/hedera-protobufs-go/services"
	"github.com/pkg/errors"
)

type _ECDSAPublicKey struct {
	*ecdsa.PublicKey
}

const _LegacyECDSAPubKeyPrefix = "302d300706052b8104000a032200"

func _ECDSAPublicKeyFromBytes(byt []byte) (*_ECDSAPublicKey, error) {
	length := len(byt)
	switch length {
	case 33:
		return _ECDSAPublicKeyFromBytesRaw(byt)
	case 47:
		return _LegacyECDSAPublicKeyFromBytesDer(byt)
	case 56:
		return _ECDSAPublicKeyFromBytesDer(byt)
	default:
		return &_ECDSAPublicKey{}, _NewErrBadKeyf("invalid compressed ECDSA public key length: %v bytes", len(byt))
	}
}

func _ECDSAPublicKeyFromBytesRaw(byt []byte) (*_ECDSAPublicKey, error) {
	if byt == nil {
		return &_ECDSAPublicKey{}, errByteArrayNull
	}
	if len(byt) != 33 {
		return &_ECDSAPublicKey{}, _NewErrBadKeyf("invalid public key length: %v bytes", len(byt))
	}

	key, err := crypto.DecompressPubkey(byt)
	if err != nil {
		return &_ECDSAPublicKey{}, err
	}

	return &_ECDSAPublicKey{
		key,
	}, nil
}

func _LegacyECDSAPublicKeyFromBytesDer(byt []byte) (*_ECDSAPublicKey, error) {
	if byt == nil {
		return &_ECDSAPublicKey{}, errByteArrayNull
	}

	given := hex.EncodeToString(byt)
	result := strings.ReplaceAll(given, _LegacyECDSAPubKeyPrefix, "")
	decoded, err := hex.DecodeString(result)
	if err != nil {
		return &_ECDSAPublicKey{}, err
	}

	if len(decoded) != 33 {
		return &_ECDSAPublicKey{}, _NewErrBadKeyf("invalid public key length: %v bytes", len(byt))
	}

	key, err := crypto.DecompressPubkey(decoded)
	if err != nil {
		return &_ECDSAPublicKey{}, err
	}

	return &_ECDSAPublicKey{
		key,
	}, nil
}
func _ECDSAPublicKeyFromBytesDer(byt []byte) (*_ECDSAPublicKey, error) {
	if byt == nil {
		return &_ECDSAPublicKey{}, errByteArrayNull
	}

	type AlgorithmIdentifier struct {
		Algorithm  asn1.ObjectIdentifier
		Parameters asn1.ObjectIdentifier
	}

	type PublicKeyInfo struct {
		AlgorithmIdentifier AlgorithmIdentifier
		PublicKey           asn1.BitString
	}

	key := &PublicKeyInfo{}
	_, err := asn1.Unmarshal(byt, key)
	if err != nil {
		return nil, err
	}

	// Check if the parsed key uses ECDSA public key algorithm
	ecdsaPublicKeyAlgorithmOID := asn1.ObjectIdentifier{1, 2, 840, 10045, 2, 1}
	if !key.AlgorithmIdentifier.Algorithm.Equal(ecdsaPublicKeyAlgorithmOID) {
		return nil, errors.New("public key is not an ECDSA public key")
	}

	// Check if the parsed key uses secp256k1 curve
	secp256k1OID := asn1.ObjectIdentifier{1, 3, 132, 0, 10}
	if !key.AlgorithmIdentifier.Parameters.Equal(secp256k1OID) {
		return nil, errors.New("public key is not a secp256k1 public key")
	}

	// Check if the public key is compressed and decompress it if necessary
	var pubKeyBytes []byte
	if key.PublicKey.Bytes[0] == 0x02 || key.PublicKey.Bytes[0] == 0x03 {
		// Decompress the public key
		pubKey, err := btcec.ParsePubKey(key.PublicKey.Bytes)
		if err != nil {
			return nil, err
		}
		pubKeyBytes = pubKey.SerializeUncompressed()
	} else {
		pubKeyBytes = key.PublicKey.Bytes
	}

	if len(pubKeyBytes) != 65 {
		return nil, errors.New("invalid public key length")
	}

	x := new(big.Int).SetBytes(pubKeyBytes[1:33])
	y := new(big.Int).SetBytes(pubKeyBytes[33:])

	pubKey := &ecdsa.PublicKey{
		Curve: btcec.S256(),
		X:     x,
		Y:     y,
	}

	// Validate the public key
	if !pubKey.Curve.IsOnCurve(pubKey.X, pubKey.Y) {
		return nil, errors.New("public key is not on the curve")
	}

	return &_ECDSAPublicKey{
		PublicKey: pubKey,
	}, nil
}

func _ECDSAPublicKeyFromString(s string) (*_ECDSAPublicKey, error) {
	byt, err := hex.DecodeString(s)
	if err != nil {
		return &_ECDSAPublicKey{}, err
	}

	return _ECDSAPublicKeyFromBytes(byt)
}

func (pk _ECDSAPublicKey) _BytesRaw() []byte {
	return crypto.CompressPubkey(pk.PublicKey)
}

func (pk _ECDSAPublicKey) _BytesDer() []byte {
	// Marshal the public key
	publicKeyBytes := pk._BytesRaw()

	// Define the public key structure
	publicKeyInfo := struct {
		Algorithm struct {
			Algorithm  asn1.ObjectIdentifier
			Parameters asn1.ObjectIdentifier
		}
		PublicKey asn1.BitString
	}{
		Algorithm: struct {
			Algorithm  asn1.ObjectIdentifier
			Parameters asn1.ObjectIdentifier
		}{
			Algorithm:  asn1.ObjectIdentifier{1, 2, 840, 10045, 2, 1}, // id-ecPublicKey
			Parameters: asn1.ObjectIdentifier{1, 3, 132, 0, 10},       // secp256k1
		},
		PublicKey: asn1.BitString{
			Bytes:     publicKeyBytes,
			BitLength: 8 * len(publicKeyBytes),
		},
	}

	// Marshal the public key info into DER format
	derBytes, err := asn1.Marshal(publicKeyInfo)
	if err != nil {
		return nil
	}

	return derBytes
}

func (pk _ECDSAPublicKey) _StringRaw() string {
	return hex.EncodeToString(pk._BytesRaw())
}
func (pk _ECDSAPublicKey) _StringDer() string {
	return hex.EncodeToString(pk._BytesDer())
}

func (pk _ECDSAPublicKey) _ToProtoKey() *services.Key {
	b := crypto.CompressPubkey(pk.PublicKey)
	return &services.Key{Key: &services.Key_ECDSASecp256K1{ECDSASecp256K1: b}}
}

func (pk _ECDSAPublicKey) _ToSignaturePairProtobuf(signature []byte) *services.SignaturePair {
	return &services.SignaturePair{
		PubKeyPrefix: pk._BytesRaw(),
		Signature: &services.SignaturePair_ECDSASecp256K1{
			ECDSASecp256K1: signature,
		},
	}
}

func (pk _ECDSAPublicKey) _Verify(message []byte, signature []byte) bool {
	return crypto.VerifySignature(pk._BytesRaw(), message, signature)
}

func (pk _ECDSAPublicKey) _VerifyTransaction(transaction Transaction) bool {
	if transaction.signedTransactions._Length() == 0 {
		return false
	}

	_, _ = transaction._BuildAllTransactions()

	for _, value := range transaction.signedTransactions.slice {
		tx := value.(*services.SignedTransaction)
		found := false
		for _, sigPair := range tx.SigMap.GetSigPair() {
			if bytes.Equal(sigPair.GetPubKeyPrefix(), pk._BytesRaw()) {
				found = true
				if !pk._Verify(tx.BodyBytes, sigPair.GetECDSASecp256K1()) {
					return false
				}
			}
		}

		if !found {
			return false
		}
	}

	return true
}

func (pk _ECDSAPublicKey) _ToFullKey() []byte {
	return elliptic.Marshal(crypto.S256(), pk.X, pk.Y)
}

func (pk _ECDSAPublicKey) _ToEthereumAddress() string {
	temp := pk._ToFullKey()[1:]
	hash := crypto.Keccak256(temp)
	return hex.EncodeToString(hash[12:])
}
