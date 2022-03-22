package hedera

import (
	"bytes"
	"crypto/ecdsa"
	"encoding/hex"
	"strings"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/hashgraph/hedera-protobufs-go/services"
)

type _ECDSAPublicKey struct {
	*ecdsa.PublicKey
}

func _ECDSAPublicKeyFromBytes(byt []byte) (*_ECDSAPublicKey, error) {
	length := len(byt)
	switch length {
	case 33:
		return _ECDSAPublicKeyFromBytesRaw(byt)
	case 49:
		return _ECDSAPublicKeyFromBytesDer(byt)
	default:
		return &_ECDSAPublicKey{}, _NewErrBadKeyf("invalid compressed ecsda public key length: %v bytes", len(byt))
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

func _ECDSAPublicKeyFromBytesDer(byt []byte) (*_ECDSAPublicKey, error) {
	if byt == nil {
		return &_ECDSAPublicKey{}, errByteArrayNull
	}

	given := hex.EncodeToString(byt)
	result := strings.ReplaceAll(given, _ECDSAPubKeyPrefix, "")
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
	decoded, _ := hex.DecodeString(_ECDSAPubKeyPrefix)
	return append(decoded, pk._BytesRaw()...)
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
