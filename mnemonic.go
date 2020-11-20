package hedera

import (
	"crypto/sha512"
	"fmt"
	"github.com/tyler-smith/go-bip39"
	"golang.org/x/crypto/pbkdf2"
	"math/big"
	"strings"
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
	if len(words) == 24 || len(words) == 12 || len(words) == 22 {
		if len(words) == 22 {
			return Mnemonic{
				words: strings.Join(words, " "),
			}.legacyValidate()
		}
		return Mnemonic{
			words: strings.Join(words, " "),
		}, nil
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
	newChecksum := crc8(entropy)

	if checksum != newChecksum {
		return Mnemonic{}, fmt.Errorf("legacy mnemonic checksum mismatch")
	}

	return m, nil
}

func (m Mnemonic) indices() ([]int, error) {
	var indices []int
	var check bool
	for _, mnemonicString := range strings.Split(m.words, " ") {
		check = false
		for i, stringCheck := range legacy {
			if mnemonicString == stringCheck {
				check = true
				indices = append(indices, int(i))
			}
		}
		if !check {
			return make([]int, 0), fmt.Errorf("word is not in the legacy word list")
		}
	}

	return indices, nil
}

func (m Mnemonic) ToLegacyPrivateKey() (PrivateKey, error) {
	indices, err := m.indices()
	if err != nil {
		return PrivateKey{}, err
	}

	entropy, _ := m.toLegacyEntropy(indices)
	password := make([]uint8, len(entropy)+8)
	for i, number := range entropy {
		password[i] = number
	}
	for i := len(entropy); i<len(password); i++ {
		password[i] = 0xFF
	}
	salt := []byte{0xFF}

	keyData := pbkdf2.Key(password, salt, 2048, 32, sha512.New)

	return PrivateKeyFromBytes(keyData)
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
