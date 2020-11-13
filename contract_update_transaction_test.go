package hedera

import (
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
	"time"
)

func TestSerializeContractUpdateTransaction(t *testing.T) {
	mockClient, err := newMockClient()
	assert.NoError(t, err)

	privateKey, err := PrivateKeyFromString(mockPrivateKey)
	assert.NoError(t, err)

	tx, err := NewContractUpdateTransaction().
		SetContractID(ContractID{Contract: 3}).
		SetAdminKey(privateKey.PublicKey()).
		SetBytecodeFileID(FileID{File: 5}).
		SetExpirationTime(time.Unix(1569375111277, 0)).
		SetProxyAccountID(AccountID{Account: 3}).
		SetAutoRenewPeriod(60 * 60 * 24 * 14 * time.Second).
		SetMaxTransactionFee(HbarFromTinybar(1e6)).
		SetTransactionID(testTransactionID).
		SetNodeAccountIDs([]AccountID{{Account: 3}}).
		FreezeWith(mockClient)

	assert.NoError(t, err)

	tx.Sign(privateKey)
	assert.NoError(t, err)

	assert.Equal(t, `bodyBytes:"\n\016\n\010\010\334\311\007\020\333\237\t\022\002\030\003\022\002\030\003\030\300\204=\"\002\010xJ?\n\002\030\003\022\007\010\355\220\257\260\326-\032\"\022\344\361\300\353L}\315\303\347\353\021p\263\010\212=\022\242\227\364\243\353\342\362\205\003\375g5F\355\2162\002\030\003:\004\010\200\352IB\002\030\005"sigMap:<sigPair:<pubKeyPrefix:"\344\361\300\353L}\315\303\347\353\021p\263\010\212=\022\242\227\364\243\353\342\362\205\003\375g5F\355\216"ed25519:"\023i\223\311Iem\320^\202}I\2127\327\270'\240#\354\206X\2053\340!\2010#j\275r\335s\216O\371\226\211\247\242\356\004\374L)\323\352\251\252C+\303L\346~k\222G\216\212\325L\n">>transactionID:<transactionValidStart:<seconds:124124nanos:151515>accountID:<accountNum:3>>nodeAccountID:<accountNum:3>transactionFee:1000000transactionValidDuration:<seconds:120>contractUpdateInstance:<contractID:<contractNum:3>expirationTime:<seconds:1569375111277>adminKey:<ed25519:"\344\361\300\353L}\315\303\347\353\021p\263\010\212=\022\242\227\364\243\353\342\362\205\003\375g5F\355\216">proxyAccountID:<accountNum:3>autoRenewPeriod:<seconds:1209600>fileID:<fileNum:5>>`, strings.ReplaceAll(strings.ReplaceAll(tx.String(), " ", ""), "\n", ""))
}

func TestContractUpdateTransaction_Execute(t *testing.T) {
	client := newTestClient(t)

	// Note: this is the bytecode for the contract found in the example for ./examples/create_simple_contract
	testContractByteCode := []byte(`608060405234801561001057600080fd5b50336000806101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055506101cb806100606000396000f3fe608060405260043610610046576000357c01000000000000000000000000000000000000000000000000000000009004806341c0e1b51461004b578063cfae321714610062575b600080fd5b34801561005757600080fd5b506100606100f2565b005b34801561006e57600080fd5b50610077610162565b6040518080602001828103825283818151815260200191508051906020019080838360005b838110156100b757808201518184015260208101905061009c565b50505050905090810190601f1680156100e45780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b6000809054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff161415610160573373ffffffffffffffffffffffffffffffffffffffff16ff5b565b60606040805190810160405280600d81526020017f48656c6c6f2c20776f726c64210000000000000000000000000000000000000081525090509056fea165627a7a72305820ae96fb3af7cde9c0abfe365272441894ab717f816f07f41f07b1cbede54e256e0029`)

	resp, err := NewFileCreateTransaction().
		SetKeys(client.GetOperatorPublicKey()).
		SetContents(testContractByteCode).
		SetMaxTransactionFee(NewHbar(5)).
		Execute(client)

	assert.NoError(t, err)

	receipt, err := resp.GetReceipt(client)
	assert.NoError(t, err)

	fileID := *receipt.FileID
	assert.NotNil(t, fileID)

	resp, err = NewContractCreateTransaction().
		SetAdminKey(client.GetOperatorPublicKey()).
		SetGas(2000).
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		SetConstructorParameters(NewContractFunctionParameters().AddString("hello from hedera")).
		SetBytecodeFileID(fileID).
		SetContractMemo("[e2e::ContractCreateTransaction]").
		SetMaxTransactionFee(NewHbar(20)).
		Execute(client)
	assert.NoError(t, err)

	receipt, err = resp.GetReceipt(client)
	assert.NoError(t, err)

	contractID := *receipt.ContractID
	assert.NotNil(t, contractID)

	info, err := NewContractInfoQuery().
		SetContractID(contractID).
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		SetQueryPayment(NewHbar(1)).
		Execute(client)

	assert.NotNil(t, info, info.Storage)
	assert.Equal(t, info.ContractID, contractID)
	assert.NotNil(t, info.AccountID)
	assert.Equal(t, info.AccountID.String(), contractID.String())
	assert.NotNil(t, info.AdminKey)
	assert.Equal(t, info.AdminKey.String(), client.GetOperatorPublicKey().String())
	//assert.Equal(t, info.Storage, uint64(926))
	assert.Equal(t, info.ContractMemo, "[e2e::ContractCreateTransaction]")

	resp, err = NewContractUpdateTransaction().
		SetContractID(contractID).
		SetNodeAccountIDs([]AccountID{resp.NodeID}).SetContractMemo("[e2e::ContractUpdateTransaction]").
		SetMaxTransactionFee(NewHbar(5)).
		Execute(client)

	_, err = resp.GetReceipt(client)
	assert.NoError(t, err)

	info, err = NewContractInfoQuery().
		SetContractID(contractID).
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		SetQueryPayment(NewHbar(5)).
		Execute(client)

	assert.NotNil(t, info)
	assert.Equal(t, info.ContractID, contractID)
	assert.NotNil(t, info.AccountID)
	assert.Equal(t, info.AccountID.String(), contractID.String())
	assert.NotNil(t, info.AdminKey)
	assert.Equal(t, info.AdminKey.String(), client.GetOperatorPublicKey().String())
	//assert.Equal(t, info.Storage, uint64(926))
	assert.Equal(t, info.ContractMemo, "[e2e::ContractUpdateTransaction]")

	resp, err = NewContractDeleteTransaction().
		SetContractID(contractID).
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		Execute(client)
	assert.NoError(t, err)

	_, err = resp.GetReceipt(client)
	assert.NoError(t, err)

	resp, err = NewFileDeleteTransaction().
		SetFileID(fileID).
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		Execute(client)
	assert.NoError(t, err)

	_, err = resp.GetReceipt(client)
	assert.NoError(t, err)
}
