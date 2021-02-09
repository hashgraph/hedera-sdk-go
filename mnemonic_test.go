package hedera

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateMnemonic(t *testing.T) {
	mnemonic, err := GenerateMnemonic24()
	assert.NoError(t, err)

	assert.Equal(t, 24, len(mnemonic.Words()))
}

func TestMnemonicFromString(t *testing.T) {
	mnemonic, err := MnemonicFromString(testMnemonic)
	assert.NoError(t, err)

	assert.Equal(t, testMnemonic, mnemonic.String())
	assert.Equal(t, 24, len(mnemonic.Words()))
}

func TestNewMnemonicFromGeneratedMnemonic(t *testing.T) {
	generatedMnemonic, err := GenerateMnemonic24()
	assert.NoError(t, err)

	mnemonicFromSlice, err := NewMnemonic(generatedMnemonic.Words())
	assert.NoError(t, err)
	assert.Equal(t, generatedMnemonic.words, mnemonicFromSlice.words)

	mnemonicFromString, err := MnemonicFromString(generatedMnemonic.String())
	assert.NoError(t, err)
	assert.Equal(t, generatedMnemonic, mnemonicFromString)

	gKey, err := generatedMnemonic.ToPrivateKey(passphrase)
	assert.NoError(t, err)

	slKey, err := generatedMnemonic.ToPrivateKey(passphrase)
	assert.NoError(t, err)

	stKey, err := generatedMnemonic.ToPrivateKey(passphrase)
	assert.NoError(t, err)

	assert.Equal(t, gKey.keyData, slKey.keyData)
	assert.Equal(t, gKey.keyData, stKey.keyData)
}

func TestMnemonicBreaksWithBadLength(t *testing.T) {
	// note this mnemonic is probably invalid and is only used to test breakage based on length
	shortMnemonic := "inmate flip alley wear offer often piece magnet surge toddler submit right"

	_, err := MnemonicFromString(shortMnemonic)
	assert.Error(t, err)

	_, err = NewMnemonic(strings.Split(shortMnemonic, " "))
	assert.Error(t, err)
}

func TestNewMnemonic(t *testing.T) {
	legacyString := "obvious favorite remain caution remove laptop base vacant increase video erase pass sniff sausage knock grid argue salt romance way alone fever slush dune"

	mnemonicLegacy, err := NewMnemonic(strings.Split(legacyString, " "))
	assert.NoError(t, err)

	gKey, err := mnemonicLegacy.ToLegacyPrivateKey()
	assert.NoError(t, err)

	assert.Equal(t, "302e020100300506032b6570042204202b7345f302a10c2a6d55bf8b7af40f125ec41d780957826006d30776f0c441fb", gKey.String())
}

func TestLegacyMnemonic(t *testing.T) {
	legacyString := "jolly,kidnap,tom,lawn,drunk,chick,optic,lust,mutter,mole,bride,galley,dense,member,sage,neural,widow,decide,curb,aboard,margin,manure"

	mnemonicLegacy, err := NewMnemonic(strings.Split(legacyString, ","))
	assert.NoError(t, err)

	legacyWithSpaces := strings.Join(strings.Split(legacyString, ","), " ")

	mnemonicFromString, err := MnemonicFromString(legacyWithSpaces)
	assert.NoError(t, err)
	assert.Equal(t, mnemonicLegacy, mnemonicFromString)

	gKey, err := mnemonicLegacy.ToLegacyPrivateKey()
	assert.NoError(t, err)

	slKey, err := mnemonicLegacy.ToLegacyPrivateKey()
	assert.NoError(t, err)

	stKey, err := mnemonicLegacy.ToLegacyPrivateKey()
	assert.NoError(t, err)

	assert.Equal(t, gKey.keyData, slKey.keyData)
	assert.Equal(t, gKey.keyData, stKey.keyData)
	assert.Equal(t, gKey.String(), "302e020100300506032b657004220420882a565ad8cb45643892b5366c1ee1c1ef4a730c5ce821a219ff49b6bf173ddf")
}
