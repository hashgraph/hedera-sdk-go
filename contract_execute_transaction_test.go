package hedera

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestSerializeContractExecuteTransaction(t *testing.T) {
	mockClient, err := newMockClient()
	assert.NoError(t, err)

	privateKey, err := Ed25519PrivateKeyFromString(mockPrivateKey)
	assert.NoError(t, err)

	parameters := NewContractFunctionParams().
		AddBytes([]byte{24, 43, 11})

	tx, err := NewContractExecuteTransaction().
		SetContractID(ContractID{Contract: 5}).
		SetGas(141).
		SetPayableAmount(HbarFromTinybar(10000)).
		SetFunction("someFunction", parameters).
		SetMaxTransactionFee(HbarFromTinybar(1e6)).
		SetTransactionID(testTransactionID).
		Build(mockClient)

	assert.NoError(t, err)

	tx.Sign(privateKey)

	assert.Equal(t, `bodyBytes:"\n\016\n\010\010\334\311\007\020\333\237\t\022\002\030\003\022\002\030\003\030\300\204=\"\002\010x:p\n\002\030\005\020\215\001\030\220N\"d7Q\372\261\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\003\030+\013\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000"sigMap:<sigPair:<pubKeyPrefix:"\344\361\300\353L}\315\303\347\353\021p\263\010\212=\022\242\227\364\243\353\342\362\205\003\375g5F\355\216"ed25519:"pHk\026\274\217\017\216\235\324\220\275\366\236\304a\371\021No\030@R\360\225\035^\325a\377\024\246Y\373\023\204\330\3079M$:\252\277\341\264\013\037Q\037\311\310\214^\n\326\003\320O\001\215R\217\017">>transactionID:<transactionValidStart:<seconds:124124nanos:151515>accountID:<accountNum:3>>nodeAccountID:<accountNum:3>transactionFee:1000000transactionValidDuration:<seconds:120>contractCall:<contractID:<contractNum:5>gas:141amount:10000functionParameters:"7Q\372\261\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\003\030+\013\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000">`, strings.ReplaceAll(strings.ReplaceAll(tx.String(), " ", ""), "\n", ""))
}

func TestContractExecuteTransaction_Execute(t *testing.T) {
	type contract struct {
		Abi string `json:"abi"`
		Bin string `json:"bin"`
	}

	type contracts struct {
		Contracts map[string]contract `json:"contracts"`
		Version   string              `json:"version"`
	}
	var smartContract contracts = contracts{}

	// Note: this is the same contract used for the example in ./examples/create_stateful_contract
	rawContract := `{"contracts":{"stateful.sol:StatefulContract":{"abi":"[{\"inputs\":[{\"internalType\":\"string\",\"name\":\"message_\",\"type\":\"string\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"getMessage\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"kill\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"message_\",\"type\":\"string\"}],\"name\":\"setMessage\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]","bin":"608060405234801561001057600080fd5b506040516104d73803806104d78339818101604052602081101561003357600080fd5b810190808051604051939291908464010000000082111561005357600080fd5b90830190602082018581111561006857600080fd5b825164010000000081118282018810171561008257600080fd5b82525081516020918201929091019080838360005b838110156100af578181015183820152602001610097565b50505050905090810190601f1680156100dc5780820380516001836020036101000a031916815260200191505b506040525050600080546001600160a01b0319163317905550805161010890600190602084019061010f565b50506101aa565b828054600181600116156101000203166002900490600052602060002090601f016020900481019282601f1061015057805160ff191683800117855561017d565b8280016001018555821561017d579182015b8281111561017d578251825591602001919060010190610162565b5061018992915061018d565b5090565b6101a791905b808211156101895760008155600101610193565b90565b61031e806101b96000396000f3fe608060405234801561001057600080fd5b50600436106100415760003560e01c8063368b87721461004657806341c0e1b5146100ee578063ce6d41de146100f6575b600080fd5b6100ec6004803603602081101561005c57600080fd5b81019060208101813564010000000081111561007757600080fd5b82018360208201111561008957600080fd5b803590602001918460018302840111640100000000831117156100ab57600080fd5b91908080601f016020809104026020016040519081016040528093929190818152602001838380828437600092019190915250929550610173945050505050565b005b6100ec6101a2565b6100fe6101ba565b6040805160208082528351818301528351919283929083019185019080838360005b83811015610138578181015183820152602001610120565b50505050905090810190601f1680156101655780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b6000546001600160a01b0316331461018a5761019f565b805161019d906001906020840190610250565b505b50565b6000546001600160a01b03163314156101b85733ff5b565b60018054604080516020601f600260001961010087891615020190951694909404938401819004810282018101909252828152606093909290918301828280156102455780601f1061021a57610100808354040283529160200191610245565b820191906000526020600020905b81548152906001019060200180831161022857829003601f168201915b505050505090505b90565b828054600181600116156101000203166002900490600052602060002090601f016020900481019282601f1061029157805160ff19168380011785556102be565b828001600101855582156102be579182015b828111156102be5782518255916020019190600101906102a3565b506102ca9291506102ce565b5090565b61024d91905b808211156102ca57600081556001016102d456fea264697066735822122084964d4c3f6bc912a9d20e14e449721012d625aa3c8a12de41ae5519752fc89064736f6c63430006000033"}},"version":"0.6.0+commit.26b70077.Linux.g++"}`

	err := json.Unmarshal([]byte(rawContract), &smartContract)
	assert.NoError(t, err)

	testContractByteCode := []byte(smartContract.Contracts["stateful.sol:StatefulContract"].Bin)

	client := newTestClient(t)

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
		SetContractMemo("hedera-sdk-go::TestContractExecuteTransaction_Execute").
		SetMaxTransactionFee(NewHbar(20)).
		Execute(client)
	assert.NoError(t, err)

	receipt, err = txID.GetReceipt(client)
	assert.NoError(t, err)

	contractID := receipt.GetContractID()
	assert.NotNil(t, contractID)

	contractFunctionParams := NewContractFunctionParams().
		AddString("Hello from Hedera again!")

	txID, err = NewContractExecuteTransaction().
		SetContractID(contractID).
		SetMaxTransactionFee(NewHbar(5)).
		SetGas(7000).
		SetFunction("setMessage", contractFunctionParams).
		SetTransactionMemo("hedera-sdk-go::TestContractExecuteTransaction_Execute").
		Execute(client)
	assert.NoError(t, err)

	record, err := txID.GetRecord(client)
	assert.NoError(t, err)

	result, err := record.GetContractExecuteResult()
	assert.NoError(t, err)

	assert.NotZero(t, result.GasUsed)

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
