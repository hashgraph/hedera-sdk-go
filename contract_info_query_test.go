package hedera

import (
	"os"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewContractInfoQuery(t *testing.T) {
	mockTransaction, err := newMockTransaction()
	assert.NoError(t, err)

	query := NewContractInfoQuery().
		SetContractID(ContractID{Contract: 3}).
		SetQueryPaymentTransaction(mockTransaction)

	assert.Equal(t, `contractGetInfo:{header:{payment:{bodyBytes:"\n\x0e\n\x08\x08\xdc\xc9\x07\x10۟\t\x12\x02\x18\x03\x12\x02\x18\x03\x18\x80\xc2\xd7/\"\x02\x08xr\x14\n\x12\n\x07\n\x02\x18\x02\x10\xc7\x01\n\x07\n\x02\x18\x03\x10\xc8\x01"sigMap:{sigPair:{pubKeyPrefix:"\xe4\xf1\xc0\xebL}\xcd\xc3\xe7\xeb\x11p\xb3\x08\x8a=\x12\xa2\x97\xf4\xa3\xeb\xe2\xf2\x85\x03\xfdg5F\xed\x8e"ed25519:"\x12&5\x96\xfb\xb4\x1c]P\xbb%\xecP\x9bk͙\x0b߼\xac)\xa6+\xd2<\x97+\xbb\x8c\x8af\xcb\xdai\x17T4{\xf7\xf3UYn\n\x8f\xabep\x04\xf6\x83\x0f\xbaFUP\xa3\xd1/\x1d\x9d\x1a\x0b"}}}}contractID:{contractNum:3}}`, strings.ReplaceAll(query.QueryBuilder.pb.String(), " ", ""))
}

func TestContractInfoQuery_Execute(t *testing.T) {
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

	time.Sleep(5 * time.Second)

	txID, err = NewContractCreateTransaction().
		SetAdminKey(client.GetOperatorKey()).
		SetGas(2000).
		SetConstructorParams(NewContractFunctionParams().AddString("hello from hedera")).
		SetBytecodeFileID(fileID).
		SetContractMemo("hedera-sdk-go::TestContractInfoQuery_Execute").
		SetMaxTransactionFee(NewHbar(20)).
		Execute(client)
	assert.NoError(t, err)

	time.Sleep(5 * time.Second)

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
