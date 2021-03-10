package hedera

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestSerializeContractInfoQuery(t *testing.T) {
	query := NewContractInfoQuery().
		SetContractID(ContractID{Contract: 3}).
		Query

	assert.Equal(t, `contractGetInfo:{header:{}contractID:{contractNum:3}}`, strings.ReplaceAll(query.pb.String(), " ", ""))
}

func TestContractInfoQuery_Execute(t *testing.T) {
	// Note: this is the bytecode for the contract found in the example for ./examples/create_simple_contract
	testContractByteCode := []byte(`608060405234801561001057600080fd5b50336000806101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055506101cb806100606000396000f3fe608060405260043610610046576000357c01000000000000000000000000000000000000000000000000000000009004806341c0e1b51461004b578063cfae321714610062575b600080fd5b34801561005757600080fd5b506100606100f2565b005b34801561006e57600080fd5b50610077610162565b6040518080602001828103825283818151815260200191508051906020019080838360005b838110156100b757808201518184015260208101905061009c565b50505050905090810190601f1680156100e45780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b6000809054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff161415610160573373ffffffffffffffffffffffffffffffffffffffff16ff5b565b60606040805190810160405280600d81526020017f48656c6c6f2c20776f726c64210000000000000000000000000000000000000081525090509056fea165627a7a72305820ae96fb3af7cde9c0abfe365272441894ab717f816f07f41f07b1cbede54e256e0029`)

	client := newTestClient(t, false)

	resp, err := NewFileCreateTransaction().
		SetKeys(client.GetOperatorPublicKey()).
		SetContents(testContractByteCode).
		Execute(client)
	assert.NoError(t, err)

	receipt, err := resp.GetReceipt(client)
	assert.NoError(t, err)

	fileID := *receipt.FileID
	assert.NotNil(t, fileID)

	_, err = resp.GetReceipt(client)
	assert.NoError(t, err)

	resp, err = NewContractCreateTransaction().
		SetAdminKey(client.GetOperatorPublicKey()).
		SetGas(2000).
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		SetConstructorParameters(NewContractFunctionParameters().AddString("hello from hedera")).
		SetBytecodeFileID(fileID).
		SetContractMemo("hedera-sdk-go::TestContractInfoQuery_Execute").
		Execute(client)
	assert.NoError(t, err)

	receipt, err = resp.GetReceipt(client)
	assert.NoError(t, err)

	assert.NotNil(t, receipt.ContractID)
	contractID := *receipt.ContractID

	info, err := NewContractInfoQuery().
		SetContractID(contractID).
		SetMaxQueryPayment(NewHbar(1)).
		SetQueryPayment(HbarFromTinybar(25)).
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		Execute(client)
	assert.NoError(t, err)

	assert.Equal(t, contractID, info.ContractID)
	assert.Equal(t, client.GetOperatorPublicKey(), info.AdminKey)
	assert.Equal(t, "hedera-sdk-go::TestContractInfoQuery_Execute", info.ContractMemo)

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

func TestContractInfoQueryCost_Execute(t *testing.T) {
	// Note: this is the bytecode for the contract found in the example for ./examples/create_simple_contract
	testContractByteCode := []byte(`608060405234801561001057600080fd5b50336000806101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055506101cb806100606000396000f3fe608060405260043610610046576000357c01000000000000000000000000000000000000000000000000000000009004806341c0e1b51461004b578063cfae321714610062575b600080fd5b34801561005757600080fd5b506100606100f2565b005b34801561006e57600080fd5b50610077610162565b6040518080602001828103825283818151815260200191508051906020019080838360005b838110156100b757808201518184015260208101905061009c565b50505050905090810190601f1680156100e45780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b6000809054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff161415610160573373ffffffffffffffffffffffffffffffffffffffff16ff5b565b60606040805190810160405280600d81526020017f48656c6c6f2c20776f726c64210000000000000000000000000000000000000081525090509056fea165627a7a72305820ae96fb3af7cde9c0abfe365272441894ab717f816f07f41f07b1cbede54e256e0029`)

	client := newTestClient(t, false)

	resp, err := NewFileCreateTransaction().
		SetKeys(client.GetOperatorPublicKey()).
		SetContents(testContractByteCode).
		Execute(client)
	assert.NoError(t, err)

	receipt, err := resp.GetReceipt(client)
	assert.NoError(t, err)

	fileID := *receipt.FileID
	assert.NotNil(t, fileID)

	_, err = resp.GetReceipt(client)
	assert.NoError(t, err)

	resp, err = NewContractCreateTransaction().
		SetAdminKey(client.GetOperatorPublicKey()).
		SetGas(2000).
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		SetConstructorParameters(NewContractFunctionParameters().AddString("hello from hedera")).
		SetBytecodeFileID(fileID).
		SetContractMemo("hedera-sdk-go::TestContractInfoQuery_Execute").
		Execute(client)
	assert.NoError(t, err)

	receipt, err = resp.GetReceipt(client)
	assert.NoError(t, err)

	assert.NotNil(t, receipt.ContractID)
	contractID := *receipt.ContractID

	contractInfo := NewContractInfoQuery().
		SetContractID(contractID).
		SetMaxQueryPayment(NewHbar(1)).
		SetNodeAccountIDs([]AccountID{resp.NodeID})

	cost, err := contractInfo.GetCost(client)
	assert.NoError(t, err)

	info, err := contractInfo.SetQueryPayment(cost).Execute(client)
	assert.NoError(t, err)

	assert.Equal(t, contractID, info.ContractID)
	assert.Equal(t, client.GetOperatorPublicKey(), info.AdminKey)
	assert.Equal(t, "hedera-sdk-go::TestContractInfoQuery_Execute", info.ContractMemo)

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

func TestContractInfoQueryCost_BigMax_Execute(t *testing.T) {
	// Note: this is the bytecode for the contract found in the example for ./examples/create_simple_contract
	testContractByteCode := []byte(`608060405234801561001057600080fd5b50336000806101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055506101cb806100606000396000f3fe608060405260043610610046576000357c01000000000000000000000000000000000000000000000000000000009004806341c0e1b51461004b578063cfae321714610062575b600080fd5b34801561005757600080fd5b506100606100f2565b005b34801561006e57600080fd5b50610077610162565b6040518080602001828103825283818151815260200191508051906020019080838360005b838110156100b757808201518184015260208101905061009c565b50505050905090810190601f1680156100e45780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b6000809054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff161415610160573373ffffffffffffffffffffffffffffffffffffffff16ff5b565b60606040805190810160405280600d81526020017f48656c6c6f2c20776f726c64210000000000000000000000000000000000000081525090509056fea165627a7a72305820ae96fb3af7cde9c0abfe365272441894ab717f816f07f41f07b1cbede54e256e0029`)

	client := newTestClient(t, false)

	resp, err := NewFileCreateTransaction().
		SetKeys(client.GetOperatorPublicKey()).
		SetContents(testContractByteCode).
		Execute(client)
	assert.NoError(t, err)

	receipt, err := resp.GetReceipt(client)
	assert.NoError(t, err)

	fileID := *receipt.FileID
	assert.NotNil(t, fileID)

	_, err = resp.GetReceipt(client)
	assert.NoError(t, err)

	resp, err = NewContractCreateTransaction().
		SetAdminKey(client.GetOperatorPublicKey()).
		SetGas(2000).
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		SetConstructorParameters(NewContractFunctionParameters().AddString("hello from hedera")).
		SetBytecodeFileID(fileID).
		SetContractMemo("hedera-sdk-go::TestContractInfoQuery_Execute").
		Execute(client)
	assert.NoError(t, err)

	receipt, err = resp.GetReceipt(client)
	assert.NoError(t, err)

	assert.NotNil(t, receipt.ContractID)
	contractID := *receipt.ContractID

	contractInfo := NewContractInfoQuery().
		SetContractID(contractID).
		SetMaxQueryPayment(NewHbar(100000)).
		SetNodeAccountIDs([]AccountID{resp.NodeID})

	cost, err := contractInfo.GetCost(client)
	assert.NoError(t, err)

	info, err := contractInfo.SetQueryPayment(cost).Execute(client)
	assert.NoError(t, err)

	assert.Equal(t, contractID, info.ContractID)
	assert.Equal(t, client.GetOperatorPublicKey(), info.AdminKey)
	assert.Equal(t, "hedera-sdk-go::TestContractInfoQuery_Execute", info.ContractMemo)

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

func TestContractInfoQueryCost_SmallMax_Execute(t *testing.T) {
	// Note: this is the bytecode for the contract found in the example for ./examples/create_simple_contract
	testContractByteCode := []byte(`608060405234801561001057600080fd5b50336000806101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055506101cb806100606000396000f3fe608060405260043610610046576000357c01000000000000000000000000000000000000000000000000000000009004806341c0e1b51461004b578063cfae321714610062575b600080fd5b34801561005757600080fd5b506100606100f2565b005b34801561006e57600080fd5b50610077610162565b6040518080602001828103825283818151815260200191508051906020019080838360005b838110156100b757808201518184015260208101905061009c565b50505050905090810190601f1680156100e45780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b6000809054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff161415610160573373ffffffffffffffffffffffffffffffffffffffff16ff5b565b60606040805190810160405280600d81526020017f48656c6c6f2c20776f726c64210000000000000000000000000000000000000081525090509056fea165627a7a72305820ae96fb3af7cde9c0abfe365272441894ab717f816f07f41f07b1cbede54e256e0029`)

	client := newTestClient(t, false)

	resp, err := NewFileCreateTransaction().
		SetKeys(client.GetOperatorPublicKey()).
		SetContents(testContractByteCode).
		Execute(client)
	assert.NoError(t, err)

	receipt, err := resp.GetReceipt(client)
	assert.NoError(t, err)

	fileID := *receipt.FileID
	assert.NotNil(t, fileID)

	_, err = resp.GetReceipt(client)
	assert.NoError(t, err)

	resp, err = NewContractCreateTransaction().
		SetAdminKey(client.GetOperatorPublicKey()).
		SetGas(2000).
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		SetConstructorParameters(NewContractFunctionParameters().AddString("hello from hedera")).
		SetBytecodeFileID(fileID).
		SetContractMemo("hedera-sdk-go::TestContractInfoQuery_Execute").
		Execute(client)
	assert.NoError(t, err)

	receipt, err = resp.GetReceipt(client)
	assert.NoError(t, err)

	assert.NotNil(t, receipt.ContractID)
	contractID := *receipt.ContractID

	contractInfo := NewContractInfoQuery().
		SetContractID(contractID).
		SetMaxQueryPayment(HbarFromTinybar(1)).
		SetNodeAccountIDs([]AccountID{resp.NodeID})

	cost, err := contractInfo.GetCost(client)
	assert.NoError(t, err)

	_, err = contractInfo.Execute(client)
	if err != nil {
		assert.Equal(t, fmt.Sprintf("cost of ContractInfoQuery ("+cost.String()+") without explicit payment is greater than the max query payment of 1 tħ"), err.Error())
	}

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

func TestContractInfoQueryCost_InsufficientFee_Execute(t *testing.T) {
	// Note: this is the bytecode for the contract found in the example for ./examples/create_simple_contract
	testContractByteCode := []byte(`608060405234801561001057600080fd5b50336000806101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055506101cb806100606000396000f3fe608060405260043610610046576000357c01000000000000000000000000000000000000000000000000000000009004806341c0e1b51461004b578063cfae321714610062575b600080fd5b34801561005757600080fd5b506100606100f2565b005b34801561006e57600080fd5b50610077610162565b6040518080602001828103825283818151815260200191508051906020019080838360005b838110156100b757808201518184015260208101905061009c565b50505050905090810190601f1680156100e45780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b6000809054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff161415610160573373ffffffffffffffffffffffffffffffffffffffff16ff5b565b60606040805190810160405280600d81526020017f48656c6c6f2c20776f726c64210000000000000000000000000000000000000081525090509056fea165627a7a72305820ae96fb3af7cde9c0abfe365272441894ab717f816f07f41f07b1cbede54e256e0029`)

	client := newTestClient(t, false)

	resp, err := NewFileCreateTransaction().
		SetKeys(client.GetOperatorPublicKey()).
		SetContents(testContractByteCode).
		Execute(client)
	assert.NoError(t, err)

	receipt, err := resp.GetReceipt(client)
	assert.NoError(t, err)

	fileID := *receipt.FileID
	assert.NotNil(t, fileID)

	_, err = resp.GetReceipt(client)
	assert.NoError(t, err)

	resp, err = NewContractCreateTransaction().
		SetAdminKey(client.GetOperatorPublicKey()).
		SetGas(2000).
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		SetConstructorParameters(NewContractFunctionParameters().AddString("hello from hedera")).
		SetBytecodeFileID(fileID).
		SetContractMemo("hedera-sdk-go::TestContractInfoQuery_Execute").
		Execute(client)
	assert.NoError(t, err)

	receipt, err = resp.GetReceipt(client)
	assert.NoError(t, err)

	assert.NotNil(t, receipt.ContractID)
	contractID := *receipt.ContractID

	contractInfo := NewContractInfoQuery().
		SetContractID(contractID).
		SetMaxQueryPayment(NewHbar(1)).
		SetNodeAccountIDs([]AccountID{resp.NodeID})

	_, err = contractInfo.GetCost(client)
	assert.NoError(t, err)

	_, err = contractInfo.SetQueryPayment(HbarFromTinybar(1)).Execute(client)
	if err != nil {
		assert.Equal(t, fmt.Sprintf("exceptional precheck status INSUFFICIENT_TX_FEE"), err.Error())
	}

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
}

func Test_ContractInfo_NoContractID(t *testing.T) {
	client := newTestClient(t, false)

	_, err := NewContractInfoQuery().
		SetMaxQueryPayment(NewHbar(2)).
		Execute(client)
	assert.Error(t, err)
	if err != nil {
		assert.Equal(t, fmt.Sprintf("exceptional precheck status INVALID_CONTRACT_ID"), err.Error())
	}
}
