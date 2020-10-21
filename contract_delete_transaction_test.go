package hedera

import (
	"github.com/stretchr/testify/assert"
	"os"

	"testing"
)

// func TestSerializeContractDeleteTransaction(t *testing.T) {
// 	mockClient, err := newMockClient()
// 	assert.NoError(t, err)

// 	privateKey, err := PrivateKeyFromString(mockPrivateKey)
// 	assert.NoError(t, err)

// 	tx, err := NewContractDeleteTransaction().
// 		SetContractID(ContractID{Contract: 5}).
// 		SetMaxTransactionFee(HbarFromTinybar(1e6)).
// 		SetTransactionID(testTransactionID).
// 		SetNodeAccountID(AccountID{Account: 3}).
// 		SetMaxTransactionFee(HbarFromTinybar(1e6)).
// 		FreezeWith(mockClient)

// 	assert.NoError(t, err)

// 	tx.Sign(privateKey)

// 	assert.Equal(t, `bodyBytes:"\n\016\n\010\010\334\311\007\020\333\237\t\022\002\030\003\022\002\030\003\030\300\204=\"\002\010x\262\001\004\n\002\030\005"sigMap:<sigPair:<pubKeyPrefix:"\344\361\300\353L}\315\303\347\353\021p\263\010\212=\022\242\227\364\243\353\342\362\205\003\375g5F\355\216"ed25519:"\002\037\252\273\3554\227\240V\217\231\347~S\204\227.\222\036\033reSJ\315?\240\224\341\272\271X\"\307\366\235\211k\360\264i<\224\313\220\343\022_\301w\201~e\376\203\227\2522|kg\202w\005">>transactionID:<transactionValidStart:<seconds:124124nanos:151515>accountID:<accountNum:3>>nodeAccountID:<accountNum:3>transactionFee:1000000transactionValidDuration:<seconds:120>contractDeleteInstance:<contractID:<contractNum:5>>`, strings.ReplaceAll(strings.ReplaceAll(tx.String(), " ", ""), "\n", ""))
// }

// func TestSerializeContractDeleteTransaction_WithAccountIDObtainer(t *testing.T) {
// 	mockClient, err := newMockClient()
// 	assert.NoError(t, err)

// 	privateKey, err := PrivateKeyFromString(mockPrivateKey)
// 	assert.NoError(t, err)

// 	tx, err := NewContractDeleteTransaction().
// 		SetContractID(ContractID{Contract: 5}).
// 		SetMaxTransactionFee(HbarFromTinybar(1e6)).
// 		SetTransferAccountID(AccountID{Account: 3}).
// 		SetTransactionID(testTransactionID).
// 		SetNodeAccountID(AccountID{Account: 3}).
// 		FreezeWith(mockClient)

// 	assert.NoError(t, err)

// 	tx.Sign(privateKey)

// 	assert.Equal(t, `bodyBytes:"\n\016\n\010\010\334\311\007\020\333\237\t\022\002\030\003\022\002\030\003\030\300\204=\"\002\010x\262\001\010\n\002\030\005\022\002\030\003"sigMap:<sigPair:<pubKeyPrefix:"\344\361\300\353L}\315\303\347\353\021p\263\010\212=\022\242\227\364\243\353\342\362\205\003\375g5F\355\216"ed25519:"\001\221y\266\365\355\330O\004\373&\004\227;\034)\027\320\23010\240\343?\240|\004\315\326\300\317\342-\322\325\354\027\332\374\005\t\331\320\361\262K=Vr'zb\014\347Z\342\374\0356B(\336\003\017">>transactionID:<transactionValidStart:<seconds:124124nanos:151515>accountID:<accountNum:3>>nodeAccountID:<accountNum:3>transactionFee:1000000transactionValidDuration:<seconds:120>contractDeleteInstance:<contractID:<contractNum:5>transferAccountID:<accountNum:3>>`, strings.ReplaceAll(strings.ReplaceAll(tx.String(), " ", ""), "\n", ""))
// }

// func TestSerializeContractDeleteTransaction_WithContractIDObtainer(t *testing.T) {
// 	mockClient, err := newMockClient()
// 	assert.NoError(t, err)

// 	privateKey, err := PrivateKeyFromString(mockPrivateKey)
// 	assert.NoError(t, err)

// 	tx, err := NewContractDeleteTransaction().
// 		SetContractID(ContractID{Contract: 5}).
// 		SetMaxTransactionFee(HbarFromTinybar(1e6)).
// 		SetTransferContractID(ContractID{Contract: 3}).
// 		SetTransactionID(testTransactionID).
// 		SetNodeAccountID(AccountID{Account: 3}).
// 		FreezeWith(mockClient)

// 	assert.NoError(t, err)

// 	tx.Sign(privateKey)

// 	assert.Equal(t, `bodyBytes:"\n\016\n\010\010\334\311\007\020\333\237\t\022\002\030\003\022\002\030\003\030\300\204=\"\002\010x\262\001\010\n\002\030\005\032\002\030\003"sigMap:<sigPair:<pubKeyPrefix:"\344\361\300\353L}\315\303\347\353\021p\263\010\212=\022\242\227\364\243\353\342\362\205\003\375g5F\355\216"ed25519:"\265\353ah\312\304mn\206ul\234\341F[pJ\"\342\352\220&wl\315\310UD\352:$GQ\326U\204\003\177\204\215\315k\277\342\376W]\377\312\037\237D\230aa\032\370>t\203\345\310c\017">>transactionID:<transactionValidStart:<seconds:124124nanos:151515>accountID:<accountNum:3>>nodeAccountID:<accountNum:3>transactionFee:1000000transactionValidDuration:<seconds:120>contractDeleteInstance:<contractID:<contractNum:5>transferContractID:<contractNum:3>>`, strings.ReplaceAll(strings.ReplaceAll(tx.String(), " ", ""), "\n", ""))
// }

func TestContractDeleteTransaction_Execute(t *testing.T) {
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

	nodeIDs := make([]AccountID, 1)
	nodeIDs[0] = resp.NodeID

	resp, err = NewContractCreateTransaction().
		SetAdminKey(client.GetOperatorKey()).
		SetGas(2000).
		SetNodeAccountIDs(nodeIDs).
		SetConstructorParameters(NewContractFunctionParameters().AddString("hello from hedera")).
		SetBytecodeFileID(fileID).
		SetContractMemo("hedera-sdk-go::TestContractDeleteTransaction_Execute").
		SetMaxTransactionFee(NewHbar(20)).
		Execute(client)
	assert.NoError(t, err)

	receipt, err = resp.GetReceipt(client)
	assert.NoError(t, err)

	contractID := *receipt.ContractID
	assert.NotNil(t, contractID)

	resp, err = NewContractDeleteTransaction().
		SetContractID(contractID).
		SetNodeAccountIDs(nodeIDs).
		SetMaxTransactionFee(NewHbar(5)).
		Execute(client)
	assert.NoError(t, err)

	_, err = resp.GetReceipt(client)
	assert.NoError(t, err)

	_, err = NewContractInfoQuery().
		SetContractID(contractID).
		SetNodeAccountIDs(nodeIDs).
		SetMaxQueryPayment(NewHbar(2)).
		Execute(client)
	// an error should occur if the contract was properly deleted
	assert.Error(t, err)

	status := err.(ErrHederaPreCheckStatus).Status
	assert.Equal(t, status, StatusContractDeleted)

	resp, err = NewFileDeleteTransaction().
		SetFileID(fileID).
		SetNodeAccountIDs(nodeIDs).
		SetMaxTransactionFee(NewHbar(5)).
		Execute(client)
	assert.NoError(t, err)

	_, err = resp.GetReceipt(client)
	assert.NoError(t, err)
}
