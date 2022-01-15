//go:build all || unit
// +build all unit

package hedera

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/stretchr/testify/require"
)

func TestUnitGenerate24WordMnemonic(t *testing.T) {
	mnemonic, err := GenerateMnemonic24()
	require.NoError(t, err)

	assert.Equal(t, 24, len(mnemonic.Words()))
}

func TestUnitGenerate12WordMnemonic(t *testing.T) {
	mnemonic, err := GenerateMnemonic12()
	require.NoError(t, err)

	assert.Equal(t, 12, len(mnemonic.Words()))
}

func TestUnitMnemonicFromString(t *testing.T) {
	mnemonic, err := MnemonicFromString(testMnemonic)
	require.NoError(t, err)

	assert.Equal(t, testMnemonic, mnemonic.String())
	assert.Equal(t, 24, len(mnemonic.Words()))
}

func TestUnitNew24MnemonicFromGeneratedMnemonic(t *testing.T) {
	generatedMnemonic, err := GenerateMnemonic24()
	require.NoError(t, err)

	mnemonicFromSlice, err := NewMnemonic(generatedMnemonic.Words())
	require.NoError(t, err)
	assert.Equal(t, generatedMnemonic.words, mnemonicFromSlice.words)

	mnemonicFromString, err := MnemonicFromString(generatedMnemonic.String())
	require.NoError(t, err)
	assert.Equal(t, generatedMnemonic, mnemonicFromString)

	gKey, err := generatedMnemonic.ToPrivateKey(passphrase)
	require.NoError(t, err)

	slKey, err := generatedMnemonic.ToPrivateKey(passphrase)
	require.NoError(t, err)

	stKey, err := generatedMnemonic.ToPrivateKey(passphrase)
	require.NoError(t, err)

	assert.Equal(t, gKey.ed25519PrivateKey.keyData, slKey.ed25519PrivateKey.keyData)
	assert.Equal(t, gKey.ed25519PrivateKey.keyData, stKey.ed25519PrivateKey.keyData)
}

func TestUnitNew12MnemonicFromGeneratedMnemonic(t *testing.T) {
	generatedMnemonic, err := GenerateMnemonic12()
	require.NoError(t, err)

	mnemonicFromSlice, err := NewMnemonic(generatedMnemonic.Words())
	require.NoError(t, err)
	assert.Equal(t, generatedMnemonic.words, mnemonicFromSlice.words)

	mnemonicFromString, err := MnemonicFromString(generatedMnemonic.String())
	require.NoError(t, err)
	assert.Equal(t, generatedMnemonic, mnemonicFromString)

	gKey, err := generatedMnemonic.ToPrivateKey(passphrase)
	require.NoError(t, err)

	slKey, err := mnemonicFromSlice.ToPrivateKey(passphrase)
	require.NoError(t, err)

	stKey, err := mnemonicFromString.ToPrivateKey(passphrase)
	require.NoError(t, err)

	assert.Equal(t, gKey.ed25519PrivateKey.keyData, slKey.ed25519PrivateKey.keyData)
	assert.Equal(t, gKey.ed25519PrivateKey.keyData, stKey.ed25519PrivateKey.keyData)
}

func TestUnitNewMnemonic(t *testing.T) {
	legacyString := "obvious favorite remain caution remove laptop base vacant increase video erase pass sniff sausage knock grid argue salt romance way alone fever slush dune"

	mnemonicLegacy, err := NewMnemonic(strings.Split(legacyString, " "))
	require.NoError(t, err)

	gKey, err := mnemonicLegacy.ToLegacyPrivateKey()
	require.NoError(t, err)

	assert.Equal(t, "302e020100300506032b65700422042098aa82d6125b5efa04bf8372be7931d05cd77f5ef3330b97d6ee7c006eaaf312", gKey.String())
}

func TestUnitLegacyMnemonic(t *testing.T) {
	legacyString := "jolly,kidnap,tom,lawn,drunk,chick,optic,lust,mutter,mole,bride,galley,dense,member,sage,neural,widow,decide,curb,aboard,margin,manure"

	mnemonicLegacy, err := NewMnemonic(strings.Split(legacyString, ","))
	require.NoError(t, err)

	legacyWithSpaces := strings.Join(strings.Split(legacyString, ","), " ")

	mnemonicFromString, err := MnemonicFromString(legacyWithSpaces)
	require.NoError(t, err)
	assert.Equal(t, mnemonicLegacy, mnemonicFromString)

	gKey, err := mnemonicLegacy.ToLegacyPrivateKey()
	require.NoError(t, err)

	slKey, err := mnemonicLegacy.ToLegacyPrivateKey()
	require.NoError(t, err)

	stKey, err := mnemonicLegacy.ToLegacyPrivateKey()
	require.NoError(t, err)

	assert.Equal(t, gKey.ed25519PrivateKey.keyData, slKey.ed25519PrivateKey.keyData)
	assert.Equal(t, gKey.ed25519PrivateKey.keyData, stKey.ed25519PrivateKey.keyData)
	assert.Equal(t, gKey.String(), "302e020100300506032b65700422042000c2f59212cb3417f0ee0d38e7bd876810d04f2dd2cb5c2d8f26ff406573f2bd")
}

func TestUnitMnemonicBreaksWithBadLength(t *testing.T) {
	// note this mnemonic is probably invalid and is only used to test breakage based on length
	shortMnemonic := "inmate flip alley wear offer often piece magnet surge toddler submit right business"

	_, err := MnemonicFromString(shortMnemonic)
	assert.Error(t, err)

	_, err = NewMnemonic(strings.Split(shortMnemonic, " "))
	assert.Error(t, err)
}
