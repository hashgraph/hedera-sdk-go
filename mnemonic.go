package hedera

import (
	"crypto/sha256"
	"fmt"
	"math/big"
	"strings"

	"github.com/pkg/errors"
	"github.com/tyler-smith/go-bip39"
)

type Mnemonic struct {
	words string
}

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

func (m Mnemonic) String() string {
	return m.words
}

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
			}.legacyValidate()
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

func (m Mnemonic) legacyValidate() (Mnemonic, error) {
	if len(strings.Split(m.words, " ")) != 22 {
		return Mnemonic{}, fmt.Errorf("not a legacy mnemonic")
	}

	indices, err := m.indices()
	if err != nil {
		return Mnemonic{}, err
	}

	entropy, checksum := m.toLegacyEntropy(indices)
	newchecksum := crc8(entropy)

	if checksum != newchecksum {
		return Mnemonic{}, fmt.Errorf("legacy mnemonic checksum mismatch")
	}

	return m, nil
}

func (m Mnemonic) indices() ([]int, error) {
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

func (m Mnemonic) ToLegacyPrivateKey() (PrivateKey, error) {
	indices, err := m.indices()
	if err != nil {
		return PrivateKey{}, err
	}

	var entropy []byte
	if len(indices) == 22 { // nolint
		entropy, _ = m.toLegacyEntropy(indices)
	} else if len(indices) == 24 {
		entropy, err = m.toLegacyEntropy2()
		if err != nil {
			return PrivateKey{}, err
		}
	} else {
		return PrivateKey{}, errors.New("not a legacy key")
	}

	return PrivateKeyFromBytes(entropy)
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

func (m Mnemonic) toLegacyEntropy(indices []int) ([]byte, uint8) {
	data := convertRadix(indices, len(legacy), 256, 33)

	checksum := data[len(data)-1]
	result := make([]uint8, len(data)-1)

	for i := 0; i < len(data)-1; i++ {
		result[i] = data[i] ^ checksum
	}

	return result, checksum
}

func (m Mnemonic) toLegacyEntropy2() ([]byte, error) {
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

func convertRadix(nums []int, fromRadix int, toRadix int, toLength int) []uint8 {
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

func crc8(data []uint8) uint8 {
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
