package hedera

import (
	"fmt"
	"github.com/tyler-smith/go-bip39"
	"strings"
)

type Mnemonic struct {
	words string
}

func (m Mnemonic) ToPrivateKey(passPhrase string) (Ed25519PrivateKey, error) {
	return Ed25519PrivateKeyFromMnemonic(m, passPhrase)
}

// Generate a random 24-word mnemonic
func GenerateMnemonic() (Mnemonic, error) {
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

	return Mnemonic{mnemonic }, nil
}

// Create a mnemonic from a string of 24 words separated by spaces
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

// Create a mnemonic from a list of 24 strings
// Keys are lazily generated
func NewMnemonic(words []string) (Mnemonic, error){
	if len(words) != 24 {
		return Mnemonic{}, fmt.Errorf("invalid mnemonic string")
	}

	return Mnemonic{
		words: strings.Join(words, " "),
	}, nil
}
