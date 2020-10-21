package hedera

import (
	"github.com/stretchr/testify/assert"
	"os"

	"testing"
)

// func TestSerializeContractCreateTransaction(t *testing.T) {
// 	mockClient, err := newMockClient()
// 	assert.NoError(t, err)

// 	privateKey, err := PrivateKeyFromString(mockPrivateKey)
// 	assert.NoError(t, err)

// 	tx, err := NewContractCreateTransaction().
// 		SetAdminKey(privateKey.PublicKey()).
// 		SetInitialBalance(HbarFromTinybar(1e3)).
// 		SetBytecodeFileID(FileID{File: 4}).
// 		SetGas(100).
// 		SetProxyAccountID(AccountID{Account: 3}).
// 		SetAutoRenewPeriod(60 * 60 * 24 * 14 * time.Second).
// 		SetMaxTransactionFee(HbarFromTinybar(1e6)).
// 		SetTransactionID(testTransactionID).
// 		FreezeWith(mockClient)

// 	assert.NoError(t, err)

// 	tx.Sign(privateKey)

// 	assert.Equal(t, `bodyBytes:"\n\016\n\010\010\334\311\007\020\333\237\t\022\002\030\003\022\002\030\003\030\300\204=\"\002\010xB7\n\002\030\004\032\"\022\344\361\300\353L}\315\303\347\353\021p\263\010\212=\022\242\227\364\243\353\342\362\205\003\375g5F\355\216d(\350\0072\002\030\003B\004\010\200\352I"sigMap:<sigPair:<pubKeyPrefix:"\344\361\300\353L}\315\303\347\353\021p\263\010\212=\022\242\227\364\243\353\342\362\205\003\375g5F\355\216"ed25519:"\335\035\364\306\364\013\204\232\376\343\320b\270\310\017I\020\223\306K\316\264\277\363x\251\3418\004>\352\272\264\202\343\260\344\264\344<l\004H\200\253m\306\305)\304t\245\227z\305\267\270w\353\326\255\340\374\005">>transactionID:<transactionValidStart:<seconds:124124nanos:151515>accountID:<accountNum:3>>nodeAccountID:<accountNum:3>transactionFee:1000000transactionValidDuration:<seconds:120>contractCreateInstance:<fileID:<fileNum:4>adminKey:<ed25519:"\344\361\300\353L}\315\303\347\353\021p\263\010\212=\022\242\227\364\243\353\342\362\205\003\375g5F\355\216">gas:100initialBalance:1000proxyAccountID:<accountNum:3>autoRenewPeriod:<seconds:1209600>>`, strings.ReplaceAll(strings.ReplaceAll(tx.String(), " ", ""), "\n", ""))
// }

func TestContractCreateTransaction_Execute(t *testing.T) {
	// Note: this is the bytecode for the contract found in the example for ./examples/create_simple_contract
	testContractByteCode := []byte(`608060405234801561001057600080fd5b50336000806101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055506101cb806100606000396000f3fe608060405260043610610046576000357c01000000000000000000000000000000000000000000000000000000009004806341c0e1b51461004b578063cfae321714610062575b600080fd5b34801561005757600080fd5b506100606100f2565b005b34801561006e57600080fd5b50610077610162565b6040518080602001828103825283818151815260200191508051906020019080838360005b838110156100b757808201518184015260208101905061009c565b50505050905090810190601f1680156100e45780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b6000809054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff161415610160573373ffffffffffffffffffffffffffffffffffffffff16ff5b565b60606040805190810160405280600d81526020017f48656c6c6f2c20776f726c64210000000000000000000000000000000000000081525090509056fea165627a7a72305820ae96fb3af7cde9c0abfe365272441894ab717f816f07f41f07b1cbede54e256e0029`)

	client, err := ClientFromJsonFile(os.Getenv("CONFIG_FILE"))

	if err != nil {
		client = ClientForTestnet()
	}

	configOperatorID := os.Getenv("OPERATOR_ID")
	configOperatorKey := os.Getenv("OPERATOR_KEY")

	if configOperatorID != "" && configOperatorKey != "" {
		operatorAccountID, err := AccountIDFromString(configOperatorID)
		assert.NoError(t, err)

		operatorKey, err := PrivateKeyFromString(configOperatorKey)
		assert.NoError(t, err)

		client.SetOperator(operatorAccountID, operatorKey)
	}

	resp, err := NewFileCreateTransaction().
		SetKeys(client.GetOperatorKey()).
		SetContents(testContractByteCode).
		SetMaxTransactionFee(NewHbar(3)).
		Execute(client)

	assert.NoError(t, err)

	receipt, err := resp.GetReceipt(client)
	assert.NoError(t, err)

	fileID := *receipt.FileID
	assert.NotNil(t, fileID)

	resp, err = NewContractCreateTransaction().
		SetAdminKey(client.GetOperatorKey()).
		SetGas(2000).
		SetConstructorParameters(NewContractFunctionParameters().AddString("hello from hedera")).
		SetBytecodeFileID(fileID).
		SetContractMemo("hedera-sdk-go::TestContractCreateTransaction_Execute").
		SetMaxTransactionFee(NewHbar(20)).
		Execute(client)
	assert.NoError(t, err)

	receipt, err = resp.GetReceipt(client)
	assert.NoError(t, err)

	contractID := *receipt.ContractID
	assert.NotNil(t, contractID)

	_, err = NewContractDeleteTransaction().
		SetContractID(contractID).
		SetNodeAccountID(resp.NodeID).
		SetMaxTransactionFee(NewHbar(5)).
		Execute(client)
	assert.NoError(t, err)

	_, err = NewFileDeleteTransaction().
		SetFileID(fileID).
		SetNodeAccountID(resp.NodeID).
		SetMaxTransactionFee(NewHbar(5)).
		Execute(client)
	assert.NoError(t, err)
}
