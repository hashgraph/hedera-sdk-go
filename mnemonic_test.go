//go:build all || unit
// +build all unit

package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"encoding/hex"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var mnemonicLegacyV1String = "jolly kidnap tom lawn drunk chick optic lust mutter mole bride galley dense member sage neural widow decide curb aboard margin manure"

var mnemonicLegacyV2String = "obvious favorite remain caution remove laptop base vacant increase video erase pass sniff sausage knock grid argue salt romance way alone fever slush dune"

var mnemonic24WordString = "inmate flip alley wear offer often piece magnet surge toddler submit right radio absent pear floor belt raven price stove replace reduce plate home"

var mnemonic12WordString = "finish furnace tomorrow wine mass goose festival air palm easy region guilt"

var passPhrase = "some pass"

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

func TestLegacyV1Mnemonic(t *testing.T) {
	test1PrivateKey := "00c2f59212cb3417f0ee0d38e7bd876810d04f2dd2cb5c2d8f26ff406573f2bd"
	test1PublicKey := "0c5bb4624df6b64c2f07a8cb8753945dd42d4b9a2ed4c0bf98e87ef154f473e9"
	test2PrivateKey := "fae0002d2716ea3a60c9cd05ee3c4bb88723b196341b68a02d20975f9d049dc6"
	test2PublicKey := "f40f9fdb1f161c31ed656794ada7af8025e8b5c70e538f38a4dfb46a0a6b0392"
	test3PrivateKey := "882a565ad8cb45643892b5366c1ee1c1ef4a730c5ce821a219ff49b6bf173ddf"
	test3PublicKey := "53c6b451e695d6abc52168a269316a0d20deee2331f612d4fb8b2b379e5c6854"
	test4PrivateKey := "6890dc311754ce9d3fc36bdf83301aa1c8f2556e035a6d0d13c2cccdbbab1242"
	test4PublicKey := "45f3a673984a0b4ee404a1f4404ed058475ecd177729daa042e437702f7791e9"

	mnemonic, err := MnemonicFromString(mnemonicLegacyV1String)
	require.NoError(t, err)

	key1, err := mnemonic.ToLegacyPrivateKey()
	require.NoError(t, err)

	assert.Equal(t, key1.StringRaw(), test1PrivateKey)
	assert.Equal(t, key1.PublicKey().StringRaw(), test1PublicKey)

	key2, err := key1.LegacyDerive(0)
	require.NoError(t, err)

	assert.Equal(t, key2.StringRaw(), test2PrivateKey)
	assert.Equal(t, key2.PublicKey().StringRaw(), test2PublicKey)

	key3, err := key1.LegacyDerive(-1)
	require.NoError(t, err)

	assert.Equal(t, key3.StringRaw(), test3PrivateKey)
	assert.Equal(t, key3.PublicKey().StringRaw(), test3PublicKey)

	key4, err := key1.LegacyDerive(1099511627775)
	require.NoError(t, err)

	assert.Equal(t, key4.StringRaw(), test4PrivateKey)
	assert.Equal(t, key4.PublicKey().StringRaw(), test4PublicKey)
}

func TestLegacyV2Mnemonic(t *testing.T) {
	test1PrivateKey := "98aa82d6125b5efa04bf8372be7931d05cd77f5ef3330b97d6ee7c006eaaf312"
	test1PublicKey := "e0ce688d614f22f96d9d213ca513d58a7d03d954fe45790006e6e86b25456465"
	test2PrivateKey := "2b7345f302a10c2a6d55bf8b7af40f125ec41d780957826006d30776f0c441fb"
	test2PublicKey := "0e19f99800b007cc7c82f9d85b73e0f6e48799469450caf43f253b48c4d0d91a"
	test3PrivateKey := "caffc03fdb9853e6a91a5b3c57a5c0031d164ce1c464dea88f3114786b5199e5"
	test3PublicKey := "9fe11da3fcfba5d28a6645ecb611a9a43dbe6014b102279ba1d34506ea86974b"

	mnemonic, err := MnemonicFromString(mnemonicLegacyV2String)
	require.NoError(t, err)

	key1, err := mnemonic.ToLegacyPrivateKey()
	require.NoError(t, err)

	assert.Equal(t, key1.StringRaw(), test1PrivateKey)
	assert.Equal(t, key1.PublicKey().StringRaw(), test1PublicKey)

	key2, err := key1.LegacyDerive(0)
	require.NoError(t, err)

	assert.Equal(t, key2.StringRaw(), test2PrivateKey)
	assert.Equal(t, key2.PublicKey().StringRaw(), test2PublicKey)

	key3, err := key1.LegacyDerive(-1)
	require.NoError(t, err)

	assert.Equal(t, key3.StringRaw(), test3PrivateKey)
	assert.Equal(t, key3.PublicKey().StringRaw(), test3PublicKey)
}

func TestUnitMnemonicBreaksWithBadLength(t *testing.T) {
	// note this mnemonic is probably invalid and is only used to test breakage based on length
	shortMnemonic := "inmate flip alley wear offer often piece magnet surge toddler submit right business"

	_, err := MnemonicFromString(shortMnemonic)
	assert.Error(t, err)

	_, err = NewMnemonic(strings.Split(shortMnemonic, " "))
	assert.Error(t, err)
}

func TestBIP39NFKD(t *testing.T) {
	passphrase := "\u03B4\u03BF\u03BA\u03B9\u03BC\u03AE"
	expectedPrivateKey := "302e020100300506032b6570042204203fefe1000db9485372851d542453b07e7970de4e2ecede7187d733ac037f4d2c"
	mnemonic, err := MnemonicFromString(mnemonic24WordString)
	assert.NoError(t, err)
	key, err := mnemonic.ToPrivateKey(passphrase)
	assert.NoError(t, err)
	assert.Equal(t, key.String(), expectedPrivateKey)
}
func TestBIP39Vector(t *testing.T) {
	passPhrase := "TREZOR"

	// the commented out tests are 18 words mnemonic, which are not supported.
	tests := [][]string{
		{
			"abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon about",
			"c55257c360c07c72029aebc1b53c05ed0362ada38ead3e3e9efa3708e53495531f09a6987599d18264c1e1c92f2cf141630c7a3c4ab7c81b2f001698e7463b04",
		},
		{
			"legal winner thank year wave sausage worth useful legal winner thank yellow",
			"2e8905819b8723fe2c1d161860e5ee1830318dbf49a83bd451cfb8440c28bd6fa457fe1296106559a3c80937a1c1069be3a3a5bd381ee6260e8d9739fce1f607",
		},
		{
			"letter advice cage absurd amount doctor acoustic avoid letter advice cage above",
			"d71de856f81a8acc65e6fc851a38d4d7ec216fd0796d0a6827a3ad6ed5511a30fa280f12eb2e47ed2ac03b5c462a0358d18d69fe4f985ec81778c1b370b652a8",
		},
		{
			"zoo zoo zoo zoo zoo zoo zoo zoo zoo zoo zoo wrong",
			"ac27495480225222079d7be181583751e86f571027b0497b5b5d11218e0a8a13332572917f0f8e5a589620c6f15b11c61dee327651a14c34e18231052e48c069",
		},
		// {
		// 	"abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon agent",
		// 	"035895f2f481b1b0f01fcf8c289c794660b289981a78f8106447707fdd9666ca06da5a9a565181599b79f53b844d8a71dd9f439c52a3d7b3e8a79c906ac845fa",
		// },
		// {
		// 	"legal winner thank year wave sausage worth useful legal winner thank year wave sausage worth useful legal will",
		// 	"f2b94508732bcbacbcc020faefecfc89feafa6649a5491b8c952cede496c214a0c7b3c392d168748f2d4a612bada0753b52a1c7ac53c1e93abd5c6320b9e95dd",
		// },
		// {
		// 	"letter advice cage absurd amount doctor acoustic avoid letter advice cage absurd amount doctor acoustic avoid letter always",
		// 	"107d7c02a5aa6f38c58083ff74f04c607c2d2c0ecc55501dadd72d025b751bc27fe913ffb796f841c49b1d33b610cf0e91d3aa239027f5e99fe4ce9e5088cd65",
		// },
		// {
		// 	"zoo zoo zoo zoo zoo zoo zoo zoo zoo zoo zoo zoo zoo zoo zoo zoo zoo when",
		// 	"0cd6e5d827bb62eb8fc1e262254223817fd068a74b5b449cc2f667c3f1f985a76379b43348d952e2265b4cd129090758b3e3c2c49103b5051aac2eaeb890a528",
		// },
		{
			"abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon art",
			"bda85446c68413707090a52022edd26a1c9462295029f2e60cd7c4f2bbd3097170af7a4d73245cafa9c3cca8d561a7c3de6f5d4a10be8ed2a5e608d68f92fcc8",
		},
		{
			"legal winner thank year wave sausage worth useful legal winner thank year wave sausage worth useful legal winner thank year wave sausage worth title",
			"bc09fca1804f7e69da93c2f2028eb238c227f2e9dda30cd63699232578480a4021b146ad717fbb7e451ce9eb835f43620bf5c514db0f8add49f5d121449d3e87",
		},
		{
			"letter advice cage absurd amount doctor acoustic avoid letter advice cage absurd amount doctor acoustic avoid letter advice cage absurd amount doctor acoustic bless",
			"c0c519bd0e91a2ed54357d9d1ebef6f5af218a153624cf4f2da911a0ed8f7a09e2ef61af0aca007096df430022f7a2b6fb91661a9589097069720d015e4e982f",
		},
		{
			"zoo zoo zoo zoo zoo zoo zoo zoo zoo zoo zoo zoo zoo zoo zoo zoo zoo zoo zoo zoo zoo zoo zoo vote",
			"dd48c104698c30cfe2b6142103248622fb7bb0ff692eebb00089b32d22484e1613912f0a5b694407be899ffd31ed3992c456cdf60f5d4564b8ba3f05a69890ad",
		},
		{
			"ozone drill grab fiber curtain grace pudding thank cruise elder eight picnic",
			"274ddc525802f7c828d8ef7ddbcdc5304e87ac3535913611fbbfa986d0c9e5476c91689f9c8a54fd55bd38606aa6a8595ad213d4c9c9f9aca3fb217069a41028",
		},
		// {
		// 	"gravity machine north sort system female filter attitude volume fold club stay feature office ecology stable narrow fog",
		// 	"628c3827a8823298ee685db84f55caa34b5cc195a778e52d45f59bcf75aba68e4d7590e101dc414bc1bbd5737666fbbef35d1f1903953b66624f910feef245ac",
		// },
		{
			"hamster diagram private dutch cause delay private meat slide toddler razor book happy fancy gospel tennis maple dilemma loan word shrug inflict delay length",
			"64c87cde7e12ecf6704ab95bb1408bef047c22db4cc7491c4271d170a1b213d20b385bc1588d9c7b38f1b39d415665b8a9030c9ec653d75e65f847d8fc1fc440",
		},
		{
			"scheme spot photo card baby mountain device kick cradle pact join borrow",
			"ea725895aaae8d4c1cf682c1bfd2d358d52ed9f0f0591131b559e2724bb234fca05aa9c02c57407e04ee9dc3b454aa63fbff483a8b11de949624b9f1831a9612",
		},
		// {
		// 	"horn tenant knee talent sponsor spell gate clip pulse soap slush warm silver nephew swap uncle crack brave",
		// 	"fd579828af3da1d32544ce4db5c73d53fc8acc4ddb1e3b251a31179cdb71e853c56d2fcb11aed39898ce6c34b10b5382772db8796e52837b54468aeb312cfc3d",
		// },
		{
			"panda eyebrow bullet gorilla call smoke muffin taste mesh discover soft ostrich alcohol speed nation flash devote level hobby quick inner drive ghost inside",
			"72be8e052fc4919d2adf28d5306b5474b0069df35b02303de8c1729c9538dbb6fc2d731d5f832193cd9fb6aeecbc469594a70e3dd50811b5067f3b88b28c3e8d",
		},
		{
			"cat swing flag economy stadium alone churn speed unique patch report train",
			"deb5f45449e615feff5640f2e49f933ff51895de3b4381832b3139941c57b59205a42480c52175b6efcffaa58a2503887c1e8b363a707256bdd2b587b46541f5",
		},
		// {
		// 	"light rule cinnamon wrap drastic word pride squirrel upgrade then income fatal apart sustain crack supply proud access",
		// 	"4cbdff1ca2db800fd61cae72a57475fdc6bab03e441fd63f96dabd1f183ef5b782925f00105f318309a7e9c3ea6967c7801e46c8a58082674c860a37b93eda02",
		// },
		{
			"all hour make first leader extend hole alien behind guard gospel lava path output census museum junior mass reopen famous sing advance salt reform",
			"26e975ec644423f4a4c4f4215ef09b4bd7ef924e85d1d17c4cf3f136c2863cf6df0a475045652c57eb5fb41513ca2a2d67722b77e954b4b3fc11f7590449191d",
		},
		{
			"vessel ladder alter error federal sibling chat ability sun glass valve picture",
			"2aaa9242daafcee6aa9d7269f17d4efe271e1b9a529178d7dc139cd18747090bf9d60295d0ce74309a78852a9caadf0af48aae1c6253839624076224374bc63f",
		},
		// {
		// 	"scissors invite lock maple supreme raw rapid void congress muscle digital elegant little brisk hair mango congress clump",
		// 	"7b4a10be9d98e6cba265566db7f136718e1398c71cb581e1b2f464cac1ceedf4f3e274dc270003c670ad8d02c4558b2f8e39edea2775c9e232c7cb798b069e88",
		// },
		{
			"void come effort suffer camp survey warrior heavy shoot primary clutch crush open amazing screen patrol group space point ten exist slush involve unfold",
			"01f5bced59dec48e362f2c45b5de68b9fd6c92c6634f44d6d40aab69056506f0e35524a518034ddc1192e1dacd32c1ed3eaa3c3b131c88ed8e7e54c49a5d0998",
		},
	}

	for _, test := range tests {
		test1Mnemonic, err := MnemonicFromString(test[0])
		assert.NoError(t, err)
		assert.Equal(t, hex.EncodeToString(test1Mnemonic._ToSeed(passPhrase)), test[1], "failed for mnemonic: ", test[0])
	}
}

func TestToStandardED25519PrivateKey(t *testing.T) {
	test1PrivateKey := "f8dcc99a1ced1cc59bc2fee161c26ca6d6af657da9aa654da724441343ecd16f"
	test1PublicKey := "2e42c9f5a5cdbde64afa65ce3dbaf013d5f9ff8d177f6ef4eb89fbe8c084ec0d"
	test1ChainCode := "404914563637c92d688deb9d41f3f25cbe8d6659d859cc743712fcfac72d7eda"
	test2PrivateKey := "e978a6407b74a0730f7aeb722ad64ab449b308e56006c8bff9aad070b9b66ddf"
	test2PublicKey := "c4b33dca1f83509f17b69b2686ee46b8556143f79f4b9df7fe7ed3864c0c64d0"
	test2ChainCode := "9c2b0073ac934696cd0b52c6c521b9bd1902aac134380a737282fdfe29014bf1"
	test3PrivateKey := "abeca64d2337db386e289482a252334c68c7536daaefff55dc169ddb77fbae28"
	test3PublicKey := "fd311925a7a04b38f7508931c6ae6a93e5dc4394d83dafda49b051c0017d3380"
	test3ChainCode := "699344acc5e07c77eb63b154b4c5c3d33cab8bf85ee21bea4cc29ab7f0502259"
	test4PrivateKey := "9a601db3e24b199912cec6573e6a3d01ffd3600d50524f998b8169c105165ae5"
	test4PublicKey := "cf525500706faa7752dca65a086c9381d30d72cc67f23bf334f330579074a890"
	test4ChainCode := "e5af7c95043a912af57a6e031ddcad191677c265d75c39954152a2733c750a3b"

	mnemonic, err := MnemonicFromString(mnemonic24WordString)
	require.NoError(t, err)

	// Chain m/44'/3030'/0'/0'/0'
	key1, err := mnemonic.ToStandardEd25519PrivateKey("", 0)
	require.NoError(t, err)

	assert.Equal(t, key1.StringRaw(), test1PrivateKey)
	assert.Equal(t, key1.PublicKey().StringRaw(), test1PublicKey)
	assert.Equal(t, hex.EncodeToString(key1.ed25519PrivateKey.chainCode), test1ChainCode)

	// Chain m/44'/3030'/0'/0'/2147483647'
	key2, err := mnemonic.ToStandardEd25519PrivateKey("", 2147483647)
	require.NoError(t, err)

	assert.Equal(t, key2.StringRaw(), test2PrivateKey)
	assert.Equal(t, key2.PublicKey().StringRaw(), test2PublicKey)
	assert.Equal(t, hex.EncodeToString(key2.ed25519PrivateKey.chainCode), test2ChainCode)

	// Chain m/44'/3030'/0'/0'/0'; Passphrase: "some pass"
	key3, err := mnemonic.ToStandardEd25519PrivateKey(passPhrase, 0)
	require.NoError(t, err)

	assert.Equal(t, key3.StringRaw(), test3PrivateKey)
	assert.Equal(t, key3.PublicKey().StringRaw(), test3PublicKey)
	assert.Equal(t, hex.EncodeToString(key3.ed25519PrivateKey.chainCode), test3ChainCode)

	// Chain m/44'/3030'/0'/0'/2147483647'; Passphrase: "some pass"
	key4, err := mnemonic.ToStandardEd25519PrivateKey(passPhrase, 2147483647)
	require.NoError(t, err)

	assert.Equal(t, key4.StringRaw(), test4PrivateKey)
	assert.Equal(t, key4.PublicKey().StringRaw(), test4PublicKey)
	assert.Equal(t, hex.EncodeToString(key4.ed25519PrivateKey.chainCode), test4ChainCode)
}

func TestToStandardED25519PrivateKey2(t *testing.T) {
	test1PrivateKey := "020487611f3167a68482b0f4aacdeb02cc30c52e53852af7b73779f67eeca3c5"
	test1PublicKey := "2d047ff02a2091f860633f849ea2024b23e7803cfd628c9bdd635010cbd782d3"
	test1ChainCode := "48c89d67e9920e443f09d2b14525213ff83b245c8b98d63747ea0801e6d0ff3f"
	test2PrivateKey := "d0c4484480944db698dd51936b7ecc81b0b87e8eafc3d5563c76339338f9611a"
	test2PublicKey := "a1a2573c2c45bd57b0fd054865b5b3d8f492a6e1572bf04b44471e07e2f589b2"
	test2ChainCode := "c0bcdbd9df6d8a4f214f20f3e5c7856415b68be34a1f406398c04690818bea16"
	test3PrivateKey := "d06630d6e4c17942155819bbbe0db8306cd989ba7baf3c29985c8455fbefc37f"
	test3PublicKey := "6bd0a51e0ca6fcc8b13cf25efd0b4814978bcaca7d1cf7dbedf538eb02969acb"
	test3ChainCode := "998a156855ab5398afcde06164b63c5523ff2c8900db53962cc2af191df59e1c"
	test4PrivateKey := "a095ef77ee88da28f373246e9ae143f76e5839f680746c3f921e90bf76c81b08"
	test4PublicKey := "35be6a2a37ff6bbb142e9f4d9b558308f4f75d7c51d5632c6a084257455e1461"
	test4ChainCode := "19d99506a5ce2dc0080092068d278fe29b85ffb8d9c26f8956bfca876307c79c"

	mnemonic, err := MnemonicFromString(mnemonic12WordString)
	require.NoError(t, err)

	// Chain m/44'/3030'/0'/0'/0'
	key1, err := mnemonic.ToStandardEd25519PrivateKey("", 0)
	require.NoError(t, err)

	assert.Equal(t, key1.StringRaw(), test1PrivateKey)
	assert.Equal(t, key1.PublicKey().StringRaw(), test1PublicKey)
	assert.Equal(t, hex.EncodeToString(key1.ed25519PrivateKey.chainCode), test1ChainCode)

	// Chain m/44'/3030'/0'/0'/2147483647'
	key2, err := mnemonic.ToStandardEd25519PrivateKey("", 2147483647)
	require.NoError(t, err)

	assert.Equal(t, key2.StringRaw(), test2PrivateKey)
	assert.Equal(t, key2.PublicKey().StringRaw(), test2PublicKey)
	assert.Equal(t, hex.EncodeToString(key2.ed25519PrivateKey.chainCode), test2ChainCode)

	// Chain m/44'/3030'/0'/0'/0'; Passphrase: "some pass"
	key3, err := mnemonic.ToStandardEd25519PrivateKey(passPhrase, 0)
	require.NoError(t, err)

	assert.Equal(t, key3.StringRaw(), test3PrivateKey)
	assert.Equal(t, key3.PublicKey().StringRaw(), test3PublicKey)
	assert.Equal(t, hex.EncodeToString(key3.ed25519PrivateKey.chainCode), test3ChainCode)

	// Chain m/44'/3030'/0'/0'/2147483647'; Passphrase: "some pass"
	key4, err := mnemonic.ToStandardEd25519PrivateKey(passPhrase, 2147483647)
	require.NoError(t, err)

	assert.Equal(t, key4.StringRaw(), test4PrivateKey)
	assert.Equal(t, key4.PublicKey().StringRaw(), test4PublicKey)
	assert.Equal(t, hex.EncodeToString(key4.ed25519PrivateKey.chainCode), test4ChainCode)
}

func TestToStandardED25519PrivateKeyShouldFailWhenIndexIsPreHardened(t *testing.T) {
	mnemonic, err := MnemonicFromString(mnemonic24WordString)
	require.NoError(t, err)

	hardenedIndex := ToHardenedIndex(10)

	_, err = mnemonic.ToStandardEd25519PrivateKey("", hardenedIndex)
	assert.Error(t, err)
}

func TestToStandardECDSAsecp256k1PrivateKey(t *testing.T) {
	test1PrivateKey := "0fde7bfd57ae6ec310bdd8b95967d98e8762a2c02da6f694b152cf9860860ab8"
	test1PublicKey := "03b1c064b4d04d52e51f6c8e8bb1bff75d62fa7b1446412d5901d424f6aedd6fd4"
	test1ChainCode := "7717bc71194c257d4b233e16cf48c24adef630052f874a262d19aeb2b527620d"
	test2PrivateKey := "aab7d720a32c2d1ea6123f58b074c865bb07f6c621f14cb012f66c08e64996bb"
	test2PublicKey := "03a0ea31bb3562f8a309b1436bc4b2f537301778e8a5e12b68cec26052f567a235"
	test2ChainCode := "e333da4bd9e21b5dbd2b0f6d88bad02f0fa24cf4b70b2fb613368d0364cdf8af"
	test3PrivateKey := "6df5ed217cf6d5586fdf9c69d39c843eb9d152ca19d3e41f7bab483e62f6ac25"
	test3PublicKey := "0357d69bb36fee569838fe7b325c07ca511e8c1b222873cde93fc6bb541eb7ecea"
	test3ChainCode := "0ff552587f6baef1f0818136bacac0bb37236473f6ecb5a8c1cc68a716726ed1"
	test4PrivateKey := "80df01f79ee1b1f4e9ab80491c592c0ef912194ccca1e58346c3d35cb5b7c098"
	test4PublicKey := "039ebe79f85573baa065af5883d0509a5634245f7864ddead76a008c9e42aa758d"
	test4ChainCode := "3a5048e93aad88f1c42907163ba4dce914d3aaf2eea87b4dd247ca7da7530f0b"
	test5PrivateKey := "60cb2496a623e1201d4e0e7ce5da3833cd4ec7d6c2c06bce2bcbcbc9dfef22d6"
	test5PublicKey := "02b59f348a6b69bd97afa80115e2d5331749b3c89c61297255430c487d6677f404"
	test5ChainCode := "e54254940db58ef4913a377062ac6e411daebf435ad592d262d5a66d808a8b94"
	test6PrivateKey := "100477c333028c8849250035be2a0a166a347a5074a8a727bce1db1c65181a50"
	test6PublicKey := "03d10ebfa2d8ff2cd34aa96e5ef59ca2e69316b4c0996e6d5f54b6932fe51be560"
	test6ChainCode := "cb23165e9d2d798c85effddc901a248a1a273fab2a56fe7976df97b016e7bb77"

	mnemonic, err := MnemonicFromString(mnemonic24WordString)
	require.NoError(t, err)

	// Chain m/44'/3030'/0'/0/0
	key1, err := mnemonic.ToStandardECDSAsecp256k1PrivateKey("", 0)
	require.NoError(t, err)

	assert.Equal(t, key1.StringRaw(), test1PrivateKey)
	assert.Equal(t, key1.PublicKey().StringRaw(), test1PublicKey)
	assert.Equal(t, hex.EncodeToString(key1.ecdsaPrivateKey.chainCode), test1ChainCode)

	// Chain m/44'/3030'/0'/0/0'
	key2, err := mnemonic.ToStandardECDSAsecp256k1PrivateKey("", ToHardenedIndex(0))
	require.NoError(t, err)

	assert.Equal(t, key2.StringRaw(), test2PrivateKey)
	assert.Equal(t, key2.PublicKey().StringRaw(), test2PublicKey)
	assert.Equal(t, hex.EncodeToString(key2.ecdsaPrivateKey.chainCode), test2ChainCode)

	// Chain m/44'/3030'/0'/0/0; Passphrase "some pass"
	key3, err := mnemonic.ToStandardECDSAsecp256k1PrivateKey(passPhrase, 0)
	require.NoError(t, err)

	assert.Equal(t, key3.StringRaw(), test3PrivateKey)
	assert.Equal(t, key3.PublicKey().StringRaw(), test3PublicKey)
	assert.Equal(t, hex.EncodeToString(key3.ecdsaPrivateKey.chainCode), test3ChainCode)

	// Chain m/44'/3030'/0'/0/0'; Passphrase "some pass"
	key4, err := mnemonic.ToStandardECDSAsecp256k1PrivateKey(passPhrase, ToHardenedIndex(0))
	require.NoError(t, err)

	assert.Equal(t, key4.StringRaw(), test4PrivateKey)
	assert.Equal(t, key4.PublicKey().StringRaw(), test4PublicKey)
	assert.Equal(t, hex.EncodeToString(key4.ecdsaPrivateKey.chainCode), test4ChainCode)

	// Chain m/44'/3030'/0'/0/2147483647; Passphrase "some pass"
	key5, err := mnemonic.ToStandardECDSAsecp256k1PrivateKey(passPhrase, 2147483647)
	require.NoError(t, err)

	assert.Equal(t, key5.StringRaw(), test5PrivateKey)
	assert.Equal(t, key5.PublicKey().StringRaw(), test5PublicKey)
	assert.Equal(t, hex.EncodeToString(key5.ecdsaPrivateKey.chainCode), test5ChainCode)

	// Chain m/44'/3030'/0'/0/2147483647'; Passphrase "some pass"
	key6, err := mnemonic.ToStandardECDSAsecp256k1PrivateKey(passPhrase, ToHardenedIndex(2147483647))
	require.NoError(t, err)

	assert.Equal(t, key6.StringRaw(), test6PrivateKey)
	assert.Equal(t, key6.PublicKey().StringRaw(), test6PublicKey)
	assert.Equal(t, hex.EncodeToString(key6.ecdsaPrivateKey.chainCode), test6ChainCode)
}

func TestToStandardECDSAsecp256k1PrivateKey2(t *testing.T) {
	test1PrivateKey := "f033824c20dd9949ad7a4440f67120ee02a826559ed5884077361d69b2ad51dd"
	test1PublicKey := "0294bf84a54806989a74ca4b76291d386914610b40b610d303162b9e495bc06416"
	test1ChainCode := "e76e0480faf2790e62dc1a7bac9dce51db1b3571fd74d8e264abc0d240a55d09"
	test2PrivateKey := "962f549dafe2d9c8091ac918cb4fc348ab0767353f37501067897efbc84e7651"
	test2PublicKey := "027123855357fd41d28130fbc59053192b771800d28ef47319ef277a1a032af78f"
	test2ChainCode := "60c39c6a77bd68c0aaabfe2f4711dc9c2247214c4f4dae15ad4cb76905f5f544"
	test3PrivateKey := "c139ebb363d7f441ccbdd7f58883809ec0cc3ee7a122ef67974eec8534de65e8"
	test3PublicKey := "0293bdb1507a26542ed9c1ec42afe959cf8b34f39daab4bf842cdac5fa36d50ef7"
	test3ChainCode := "911a1095b64b01f7f3a06198df3d618654e5ed65862b211997c67515e3167892"
	test4PrivateKey := "87c1d8d4bb0cebb4e230852f2a6d16f6847881294b14eb1d6058b729604afea0"
	test4PublicKey := "03358e7761a422ca1c577f145fe845c77563f164b2c93b5b34516a8fa13c2c0888"
	test4ChainCode := "64173f2dcb1d65e15e787ef882fa15f54db00209e2dab16fa1661244cd98e95c"
	test5PrivateKey := "2583170ee745191d2bb83474b1de41a1621c47f6e23db3f2bf413a1acb5709e4"
	test5PublicKey := "03f9eb27cc73f751e8e476dd1db79037a7df2c749fa75b6cc6951031370d2f95a5"
	test5ChainCode := "a7250c2b07b368a054f5c91e6a3dbe6ca3bbe01eb0489fe8778304bd0a19c711"
	test6PrivateKey := "f2d008cd7349bdab19ed85b523ba218048f35ca141a3ecbc66377ad50819e961"
	test6PublicKey := "027b653d04958d4bf83dd913a9379b4f9a1a1e64025a691830a67383bc3157c044"
	test6ChainCode := "66a1175e7690e3714d53ffce16ee6bb4eb02065516be2c2ad6bf6c9df81ec394"

	mnemonic, err := MnemonicFromString(mnemonic12WordString)
	require.NoError(t, err)

	// Chain m/44'/3030'/0'/0/0
	key1, err := mnemonic.ToStandardECDSAsecp256k1PrivateKey("", 0)
	require.NoError(t, err)

	assert.Equal(t, key1.StringRaw(), test1PrivateKey)
	assert.Equal(t, key1.PublicKey().StringRaw(), test1PublicKey)
	assert.Equal(t, hex.EncodeToString(key1.ecdsaPrivateKey.chainCode), test1ChainCode)

	// Chain m/44'/3030'/0'/0/0'
	key2, err := mnemonic.ToStandardECDSAsecp256k1PrivateKey("", ToHardenedIndex(0))
	require.NoError(t, err)

	assert.Equal(t, key2.StringRaw(), test2PrivateKey)
	assert.Equal(t, key2.PublicKey().StringRaw(), test2PublicKey)
	assert.Equal(t, hex.EncodeToString(key2.ecdsaPrivateKey.chainCode), test2ChainCode)

	// Chain m/44'/3030'/0'/0/0; Passphrase "some pass"
	key3, err := mnemonic.ToStandardECDSAsecp256k1PrivateKey(passPhrase, 0)
	require.NoError(t, err)

	assert.Equal(t, key3.StringRaw(), test3PrivateKey)
	assert.Equal(t, key3.PublicKey().StringRaw(), test3PublicKey)
	assert.Equal(t, hex.EncodeToString(key3.ecdsaPrivateKey.chainCode), test3ChainCode)

	// Chain m/44'/3030'/0'/0/0'; Passphrase "some pass"
	key4, err := mnemonic.ToStandardECDSAsecp256k1PrivateKey(passPhrase, ToHardenedIndex(0))
	require.NoError(t, err)

	assert.Equal(t, key4.StringRaw(), test4PrivateKey)
	assert.Equal(t, key4.PublicKey().StringRaw(), test4PublicKey)
	assert.Equal(t, hex.EncodeToString(key4.ecdsaPrivateKey.chainCode), test4ChainCode)

	// Chain m/44'/3030'/0'/0/2147483647; Passphrase "some pass"
	key5, err := mnemonic.ToStandardECDSAsecp256k1PrivateKey(passPhrase, 2147483647)
	require.NoError(t, err)

	assert.Equal(t, key5.StringRaw(), test5PrivateKey)
	assert.Equal(t, key5.PublicKey().StringRaw(), test5PublicKey)
	assert.Equal(t, hex.EncodeToString(key5.ecdsaPrivateKey.chainCode), test5ChainCode)

	// Chain m/44'/3030'/0'/0/2147483647'; Passphrase "some pass"
	key6, err := mnemonic.ToStandardECDSAsecp256k1PrivateKey(passPhrase, ToHardenedIndex(2147483647))
	require.NoError(t, err)

	assert.Equal(t, key6.StringRaw(), test6PrivateKey)
	assert.Equal(t, key6.PublicKey().StringRaw(), test6PublicKey)
	assert.Equal(t, hex.EncodeToString(key6.ecdsaPrivateKey.chainCode), test6ChainCode)
}

func TestToStandardECDSAsecp256k1PrivateKeyCustomDPath(t *testing.T) {
	const (
		DPATH_1       = "m/44'/60'/0'/0/0"
		PASSPHRASE_1  = ""
		CHAIN_CODE_1  = "58a9ee31eaf7499abc01952b44dbf0a2a5d6447512367f09d99381c9605bf9e8"
		PRIVATE_KEY_1 = "78f9545e40025cf7da9126a4d6a861ae34031d1c74c3404df06110c9fde371ad"
		PUBLIC_KEY_1  = "02a8f4c22eea66617d4f119e3a951b93f584949bbfee90bd555305402da6c4e569"
		DPATH_2       = "m/44'/60'/0'/0/1"
		PASSPHRASE_2  = ""
		CHAIN_CODE_2  = "6dcfc7a4914bd0e75b94a2f38afee8c247b34810202a2c64fe599ee1b88afdc9"
		PRIVATE_KEY_2 = "77ca263661ebdd5a8b33c224aeff5e7bf67eedacee68a1699d97ee8929d7b130"
		PUBLIC_KEY_2  = "03e84c9be9be53ad722038cc1943e79df27e5c1d31088adb4f0e62444f4dece683"
		DPATH_3       = "m/44'/60'/0'/0/2"
		PASSPHRASE_3  = ""
		CHAIN_CODE_3  = "c8c798d2b3696be1e7a29d1cea205507eedc2057006b9ef1cde1b4e346089e17"
		PRIVATE_KEY_3 = "31c24292eac951279b659c335e44a2e812d0f1a228b1d4d87034874d376e605a"
		PUBLIC_KEY_3  = "0207ff3faf4055c1aa7a5ad94d6ff561fac35b9ae695ef486706243667d2b4d10e"
	)

	mnemonic, err := MnemonicFromString(mnemonic24WordString)
	require.NoError(t, err)

	// m/44'/60'/0'/0/0
	key1, err := mnemonic.ToStandardECDSAsecp256k1PrivateKeyCustomDerivationPath(PASSPHRASE_1, DPATH_1)
	require.NoError(t, err)
	assert.Equal(t, hex.EncodeToString(key1.ecdsaPrivateKey.chainCode), CHAIN_CODE_1)
	assert.Equal(t, key1.StringRaw(), PRIVATE_KEY_1)
	assert.Contains(t, key1.PublicKey().StringRaw(), PUBLIC_KEY_1)

	// m/44'/60'/0'/0/1
	key2, err := mnemonic.ToStandardECDSAsecp256k1PrivateKeyCustomDerivationPath(PASSPHRASE_2, DPATH_2)
	require.NoError(t, err)
	assert.Equal(t, hex.EncodeToString(key2.ecdsaPrivateKey.chainCode), CHAIN_CODE_2)
	assert.Equal(t, key2.StringRaw(), PRIVATE_KEY_2)
	assert.Contains(t, key2.PublicKey().StringRaw(), PUBLIC_KEY_2)

	// m/44'/60'/0'/0/2
	key3, err := mnemonic.ToStandardECDSAsecp256k1PrivateKeyCustomDerivationPath(PASSPHRASE_3, DPATH_3)
	require.NoError(t, err)
	assert.Equal(t, hex.EncodeToString(key3.ecdsaPrivateKey.chainCode), CHAIN_CODE_3)
	assert.Equal(t, key3.StringRaw(), PRIVATE_KEY_3)
	assert.Contains(t, key3.PublicKey().StringRaw(), PUBLIC_KEY_3)
}
