package hedera

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestNewContractInfoQuery(t *testing.T) {
	mockTransaction, err := newMockTransaction()
	assert.NoError(t, err)

	query := NewContractInfoQuery().
		SetContractID(ContractID{Contract: 3}).
		SetQueryPaymentTransaction(mockTransaction)

	assert.Equal(t, `contractGetInfo:<header:<payment:<bodyBytes:"\n\016\n\010\010\334\311\007\020\333\237\t\022\002\030\003\022\002\030\003\030\200\302\327/\"\002\010xr\024\n\022\n\007\n\002\030\002\020\307\001\n\007\n\002\030\003\020\310\001" sigMap:<sigPair:<pubKeyPrefix:"\344\361\300\353L}\315\303\347\353\021p\263\010\212=\022\242\227\364\243\353\342\362\205\003\375g5F\355\216" ed25519:"\022&5\226\373\264\034]P\273%\354P\233k\315\231\013\337\274\254)\246+\322<\227+\273\214\212f\313\332i\027T4{\367\363UYn\n\217\253ep\004\366\203\017\272FUP\243\321/\035\235\032\013" > > > > contractID:<contractNum:3 > > `, query.QueryBuilder.pb.String())
}

func TestContractInfoQuery_Execute(t *testing.T) {
	// Note: this is the bytecode for the contract found in the example for ./examples/create_simple_contract
	testContractByteCode := []byte(`608060405234801561001057600080fd5b50336000806101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055506101cb806100606000396000f3fe608060405260043610610046576000357c01000000000000000000000000000000000000000000000000000000009004806341c0e1b51461004b578063cfae321714610062575b600080fd5b34801561005757600080fd5b506100606100f2565b005b34801561006e57600080fd5b50610077610162565b6040518080602001828103825283818151815260200191508051906020019080838360005b838110156100b757808201518184015260208101905061009c565b50505050905090810190601f1680156100e45780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b6000809054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff161415610160573373ffffffffffffffffffffffffffffffffffffffff16ff5b565b60606040805190810160405280600d81526020017f48656c6c6f2c20776f726c64210000000000000000000000000000000000000081525090509056fea165627a7a72305820ae96fb3af7cde9c0abfe365272441894ab717f816f07f41f07b1cbede54e256e0029`)

	client, err := ClientFromFile(os.Getenv("CONFIG_FILE"))

	if err != nil {
		client = ClientForTestnet()
	}

	configOperatorID := os.Getenv("OPERATOR_ID")
	configOperatorKey := os.Getenv("OPERATOR_KEY")

	if configOperatorID != "" && configOperatorKey != "" {
		operatorAccountID, err := AccountIDFromString(configOperatorID)
		assert.NoError(t, err)

		operatorKey, err := Ed25519PrivateKeyFromString(configOperatorKey)
		assert.NoError(t, err)

		client.SetOperator(operatorAccountID, operatorKey)
	}

	txID, err := NewFileCreateTransaction().
		AddKey(client.GetOperatorKey()).
		SetContents(testContractByteCode).
		SetMaxTransactionFee(NewHbar(3)).
		Execute(client)
	assert.NoError(t, err)

	receipt, err := txID.GetReceipt(client)
	assert.NoError(t, err)

	fileID := receipt.GetFileID()
	assert.NotNil(t, fileID)

	txID, err = NewContractCreateTransaction().
		SetAdminKey(client.GetOperatorKey()).
		SetGas(2000).
		SetConstructorParams(NewContractFunctionParams().AddString("hello from hedera")).
		SetBytecodeFileID(fileID).
		SetContractMemo("hedera-sdk-go::TestContractInfoQuery_Execute").
		SetMaxTransactionFee(NewHbar(20)).
		Execute(client)
	assert.NoError(t, err)

	receipt, err = txID.GetReceipt(client)
	assert.NoError(t, err)

	contractID := receipt.GetContractID()
	assert.NotNil(t, contractID)

	info, err := NewContractInfoQuery().
		SetContractID(contractID).
		SetMaxQueryPayment(NewHbar(2)).
		Execute(client)

	assert.Equal(t, contractID, info.ContractID)
	assert.Equal(t, client.GetOperatorKey(), info.AdminKey)
	assert.Equal(t, "hedera-sdk-go::TestContractInfoQuery_Execute", info.ContractMemo)

	_, err = NewContractDeleteTransaction().
		SetContractID(contractID).
		SetMaxTransactionFee(NewHbar(5)).
		Execute(client)
	assert.NoError(t, err)

	_, err = NewFileDeleteTransaction().
		SetFileID(fileID).
		SetMaxTransactionFee(NewHbar(5)).
		Execute(client)
	assert.NoError(t, err)
}
