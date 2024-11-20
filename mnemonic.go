package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"crypto/sha256"
	"fmt"
	"math/big"
	"regexp"
	"strconv"
	"strings"

	"crypto/sha512"

	"github.com/pkg/errors"
	"github.com/tyler-smith/go-bip39"
	"golang.org/x/crypto/pbkdf2"
	"golang.org/x/text/unicode/norm"
)

type Mnemonic struct {
	words string
}

// Deprecated
func (m Mnemonic) ToPrivateKey(passPhrase string) (PrivateKey, error) {
	return PrivateKeyFromMnemonic(m, passPhrase)
}

// GenerateMnemonic generates a random 24-word mnemonic
func GenerateMnemonic24() (Mnemonic, error) {
	entropy, err := bip39.NewEntropy(256)

	if err != nil {
		// It is only possible for there to be an error if the operating
		// system's rng is unreadable
		return Mnemonic{}, fmt.Errorf("could not retrieve random bytes from the operating system")
	}

	mnemonic, err := bip39.NewMnemonic(entropy)

	// Note that this should never actually fail since it is being provided by library generated mnemonic
	if err != nil {
		return Mnemonic{}, err
	}

	return Mnemonic{mnemonic}, nil
}

// GenerateMnemonic12 generates a random 12-word mnemonic
func GenerateMnemonic12() (Mnemonic, error) {
	entropy, err := bip39.NewEntropy(128)

	if err != nil {
		// It is only possible for there to be an error if the operating
		// system's rng is unreadable
		return Mnemonic{}, fmt.Errorf("could not retrieve random bytes from the operating system")
	}

	mnemonic, err := bip39.NewMnemonic(entropy)

	// Note that this should never actually fail since it is being provided by library generated mnemonic
	if err != nil {
		return Mnemonic{}, err
	}

	return Mnemonic{mnemonic}, nil
}

// MnemonicFromString creates a mnemonic from a string of 24 words separated by spaces
//
// Keys are lazily generated
func MnemonicFromString(s string) (Mnemonic, error) {
	return NewMnemonic(strings.Split(s, " "))
}

// String returns the mnemonic as a string.
func (m Mnemonic) String() string {
	return m.words
}

// Words returns the mnemonic as a slice of strings
func (m Mnemonic) Words() []string {
	return strings.Split(m.words, " ")
}

// NewMnemonic Creates a mnemonic from a slice of 24 strings
//
// Keys are lazily generated
func NewMnemonic(words []string) (Mnemonic, error) {
	joinedString := strings.Join(words, " ")

	if len(words) == 24 || len(words) == 12 || len(words) == 22 {
		if len(words) == 22 { //nolint
			return Mnemonic{
				words: joinedString,
			}._LegacyValidate()
		} else if bip39.IsMnemonicValid(joinedString) {
			return Mnemonic{
				words: joinedString,
			}, nil
		} else {
			return Mnemonic{}, fmt.Errorf("invalid mnemonic composition")
		}
	} else {
		return Mnemonic{}, fmt.Errorf("invalid mnemonic string")
	}
}

func (m Mnemonic) _LegacyValidate() (Mnemonic, error) {
	if len(strings.Split(m.words, " ")) != 22 {
		return Mnemonic{}, fmt.Errorf("not a legacy mnemonic")
	}

	indices, err := m._Indices()
	if err != nil {
		return Mnemonic{}, err
	}

	entropy, checksum := m._ToLegacyEntropy(indices)
	newchecksum := _Crc8(entropy)

	if checksum != newchecksum {
		return Mnemonic{}, fmt.Errorf("legacy mnemonic checksum mismatch")
	}

	return m, nil
}

func (m Mnemonic) _Indices() ([]int, error) {
	var indices []int
	var check bool
	temp := strings.Split(m.words, " ")
	if len(temp) == 22 { // nolint
		for _, mnemonicString := range strings.Split(m.words, " ") {
			check = false
			for i, stringCheck := range legacy {
				if mnemonicString == stringCheck {
					check = true
					indices = append(indices, i)
				}
			}
			if !check {
				return make([]int, 0), fmt.Errorf("word is not in the legacy word list")
			}
		}
	} else if len(temp) == 24 {
		for _, mnemonicString := range strings.Split(m.words, " ") {
			t, check := bip39.GetWordIndex(mnemonicString)
			if !check {
				return make([]int, 0), bip39.ErrInvalidMnemonic
			}
			indices = append(indices, t)
		}
	} else {
		return make([]int, 0), errors.New("not a 22 word or a 24 mnemonic")
	}

	return indices, nil
}

// ToLegacyPrivateKey converts a mnemonic to a legacy private key
func (m Mnemonic) ToLegacyPrivateKey() (PrivateKey, error) {
	indices, err := m._Indices()
	if err != nil {
		return PrivateKey{}, err
	}

	var entropy []byte
	if len(indices) == 22 { // nolint
		entropy, _ = m._ToLegacyEntropy(indices)
	} else if len(indices) == 24 {
		entropy, err = m._ToLegacyEntropy2()
		if err != nil {
			return PrivateKey{}, err
		}
	} else {
		return PrivateKey{}, errors.New("not a legacy key")
	}

	return PrivateKeyFromBytesEd25519(entropy)
}

func bytesToBits(dat []uint8) []bool {
	bits := make([]bool, len(dat)*8)

	for i := range bits {
		bits[i] = false
	}

	for i := 0; i < len(dat); i++ {
		for j := 0; j < 8; j++ {
			bits[(i*8)+j] = (dat[i] & (1 << (7 - j))) != 0
		}
	}

	return bits
}

func (m Mnemonic) _ToLegacyEntropy(indices []int) ([]byte, uint8) {
	data := _ConvertRadix(indices, len(legacy), 256, 33)

	checksum := data[len(data)-1]
	result := make([]uint8, len(data)-1)

	for i := 0; i < len(data)-1; i++ {
		result[i] = data[i] ^ checksum
	}

	return result, checksum
}

func (m Mnemonic) _ToLegacyEntropy2() ([]byte, error) {
	indices := strings.Split(m.words, " ")
	concatBitsLen := len(indices) * 11
	concatBits := make([]bool, concatBitsLen)

	for i := range concatBits {
		concatBits[i] = false
	}

	for index, word := range indices {
		nds, check := bip39.GetWordIndex(word)
		if !check {
			return make([]byte, 0), bip39.ErrInvalidMnemonic
		}

		for i := 0; i < 11; i++ {
			concatBits[(index*11)+i] = (nds & (1 << (10 - i))) != 0
		}
	}

	checksumBitsLen := concatBitsLen / 33
	entropyBitsLen := concatBitsLen - checksumBitsLen

	entropy := make([]uint8, entropyBitsLen/8)

	for i := 0; i < len(entropy); i++ {
		for j := 0; j < 8; j++ {
			if concatBits[(i*8)+j] {
				entropy[i] |= 1 << (7 - j)
			}
		}
	}

	hash := sha256.New()
	if _, err := hash.Write(entropy); err != nil {
		return nil, err
	}

	hashbits := bytesToBits(hash.Sum(nil))

	for i := 0; i < checksumBitsLen; i++ {
		if concatBits[entropyBitsLen+i] != hashbits[i] {
			return make([]uint8, 0), errors.New("checksum mismatch")
		}
	}

	return entropy, nil
}

func (m Mnemonic) _ToSeed(passPhrase string) []byte {
	passPhraseNFKD := norm.NFKD.String(passPhrase)
	salt := []byte("mnemonic" + passPhraseNFKD)
	seed := pbkdf2.Key([]byte(m.String()), salt, 2048, 64, sha512.New)
	return seed
}

// ToStandardEd25519PrivateKey converts a mnemonic to a standard ed25519 private key
func (m Mnemonic) ToStandardEd25519PrivateKey(passPhrase string, index uint32) (PrivateKey, error) {
	seed := m._ToSeed(passPhrase)
	derivedKey, err := _Ed25519PrivateKeyFromSeed(seed)
	if err != nil {
		return PrivateKey{}, err
	}

	keyBytes, chainCode := derivedKey.keyData, derivedKey.chainCode
	for _, i := range []uint32{44, 3030, 0, 0, index} {
		keyBytes, chainCode, err = _DeriveEd25519ChildKey(keyBytes, chainCode, i)
		if err != nil {
			return PrivateKey{}, err
		}
	}

	privateKey, err := _Ed25519PrivateKeyFromBytes(keyBytes)
	if err != nil {
		return PrivateKey{}, err
	}

	privateKey.chainCode = chainCode

	return PrivateKey{
		ed25519PrivateKey: privateKey,
	}, nil
}

// calculateDerivationPathValues converts a derivation path string to an array of integers
func calculateDerivationPathValues(derivationPath string) ([]uint32, error) {
	re := regexp.MustCompile(`m/(\d+'?)/(\d+'?)/(\d+'?)/(\d+'?)/(\d+'?)`)
	matches := re.FindStringSubmatch(derivationPath)
	if len(matches) != 6 {
		return nil, fmt.Errorf("invalid derivation path format")
	}

	values := make([]uint32, 5)
	for i, match := range matches[1:] {
		if strings.HasSuffix(match, "'") {
			match = strings.TrimSuffix(match, "'")
			value, err := strconv.Atoi(match)
			if err != nil {
				return nil, err
			}
			values[i] = ToHardenedIndex(uint32(value))
		} else {
			value, err := strconv.Atoi(match)
			if err != nil {
				return nil, err
			}
			values[i] = uint32(value)
		}
	}

	return values, nil
}

func (m Mnemonic) toStandardECDSAsecp256k1PrivateKeyImpl(passPhrase string, derivationPathValues []uint32) (PrivateKey, error) {
	seed := m._ToSeed(passPhrase)
	derivedKey, err := _ECDSAPrivateKeyFromSeed(seed)
	if err != nil {
		return PrivateKey{}, err
	}

	keyBytes, chainCode := derivedKey.keyData.ToECDSA().D.Bytes(), derivedKey.chainCode
	for _, i := range derivationPathValues {
		keyBytes, chainCode, err = _DeriveECDSAChildKey(keyBytes, chainCode, i)
		if err != nil {
			return PrivateKey{}, err
		}
	}

	privateKey, err := _ECDSAPrivateKeyFromBytes(keyBytes)
	if err != nil {
		return PrivateKey{}, err
	}

	privateKey.chainCode = chainCode

	return PrivateKey{
		ecdsaPrivateKey: privateKey,
	}, nil
}

// ToStandardECDSAsecp256k1PrivateKey converts a mnemonic to a standard ecdsa secp256k1 private key
func (m Mnemonic) ToStandardECDSAsecp256k1PrivateKeyCustomDerivationPath(passPhrase string, derivationPath string) (PrivateKey, error) {
	derivationPathValues, err := calculateDerivationPathValues(derivationPath)
	if err != nil {
		return PrivateKey{}, err
	}

	return m.toStandardECDSAsecp256k1PrivateKeyImpl(passPhrase, derivationPathValues)
}

// ToStandardECDSAsecp256k1PrivateKey converts a mnemonic to a standard ecdsa secp256k1 private key
// Uses the default derivation path of `m/44'/3030'/0'/0/${index}`
func (m Mnemonic) ToStandardECDSAsecp256k1PrivateKey(passPhrase string, index uint32) (PrivateKey, error) {
	seed := m._ToSeed(passPhrase)
	derivedKey, err := _ECDSAPrivateKeyFromSeed(seed)
	if err != nil {
		return PrivateKey{}, err
	}

	keyBytes, chainCode := derivedKey.keyData.ToECDSA().D.Bytes(), derivedKey.chainCode
	for _, i := range []uint32{
		ToHardenedIndex(44),
		ToHardenedIndex(3030),
		ToHardenedIndex(0),
		0,
		index} {
		keyBytes, chainCode, err = _DeriveECDSAChildKey(keyBytes, chainCode, i)
		if err != nil {
			return PrivateKey{}, err
		}
	}

	privateKey, err := _ECDSAPrivateKeyFromBytes(keyBytes)
	if err != nil {
		return PrivateKey{}, err
	}

	privateKey.chainCode = chainCode

	return PrivateKey{
		ecdsaPrivateKey: privateKey,
	}, nil
}

func _ConvertRadix(nums []int, fromRadix int, toRadix int, toLength int) []uint8 {
	num := big.NewInt(0)

	for _, element := range nums {
		num = num.Mul(num, big.NewInt(int64(fromRadix)))
		num = num.Add(num, big.NewInt(int64(element)))
	}

	result := make([]uint8, toLength)

	for i := toLength - 1; i >= 0; i-- {
		tem := new(big.Int).Div(num, big.NewInt(int64(toRadix)))
		rem := new(big.Int).Mod(num, big.NewInt(int64(toRadix)))
		num = num.Set(tem)
		result[i] = uint8(rem.Uint64())
	}

	return result
}

func _Crc8(data []uint8) uint8 {
	var crc uint8
	crc = 0xff

	for i := 0; i < len(data)-1; i++ {
		crc ^= data[i]
		for j := 0; j < 8; j++ {
			var temp uint8
			if crc&1 == 0 {
				temp = 0
			} else {
				temp = 0xb2
			}
			crc = crc>>1 ^ temp
		}
	}

	return crc ^ 0xff
}
