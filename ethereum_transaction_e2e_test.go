//go:build all || e2e
// +build all e2e

package hedera

/*-
 *
 * Hedera Go SDK
 *
 * Copyright (C) 2020 - 2024 Hedera Hashgraph, LLC
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

import (
	"encoding/hex"
	"io"

	"testing"

	"github.com/ethereum/go-ethereum/rlp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type EIP1559RLP struct {
	chainId        []byte
	nonce          []byte
	maxPriorityGas []byte
	maxGas         []byte
	gasLimit       []byte
	to             []byte
	value          []byte
	callData       []byte
	accessList     [][]byte
	recId          []byte
	r              []byte
	s              []byte
}

// EncodeRLP writes l as RLP list [ethereumfield1, ethereumfield2...] Omits r,s,v values on first encode
func (l *EIP1559RLP) EncodeRLP(w io.Writer) (err error) {
	fields := []interface{}{
		l.chainId,
		l.nonce,
		l.maxPriorityGas,
		l.maxGas,
		l.gasLimit,
		l.to,
		l.value,
		l.callData,
		l.accessList,
	}
	if len(l.recId) > 0 && len(l.r) > 0 && len(l.s) > 0 {
		fields = append(fields, l.recId, l.r, l.s)
	}

	return rlp.Encode(w, fields)
}

// decodeHex is a helper function that decodes a hex string and fails the test if an error occurs.
func decodeHex(t *testing.T, s string) []byte {
	bytes, err := hex.DecodeString(s)
	if err != nil {
		t.Fatalf("Failed to decode string %s: %v", s, err)
	}
	return bytes
}

// Testing the signer nonce defined in HIP-844
func TestIntegrationEthereumTransaction(t *testing.T) {
	// Skip this test because it is flaky with newest version of Local Node
	t.Skip()
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)

	smartContractBytecode := []byte("608060405234801561001057600080fd5b506040516104d73803806104d78339818101604052602081101561003357600080fd5b810190808051604051939291908464010000000082111561005357600080fd5b90830190602082018581111561006857600080fd5b825164010000000081118282018810171561008257600080fd5b82525081516020918201929091019080838360005b838110156100af578181015183820152602001610097565b50505050905090810190601f1680156100dc5780820380516001836020036101000a031916815260200191505b506040525050600080546001600160a01b0319163317905550805161010890600190602084019061010f565b50506101aa565b828054600181600116156101000203166002900490600052602060002090601f016020900481019282601f1061015057805160ff191683800117855561017d565b8280016001018555821561017d579182015b8281111561017d578251825591602001919060010190610162565b5061018992915061018d565b5090565b6101a791905b808211156101895760008155600101610193565b90565b61031e806101b96000396000f3fe608060405234801561001057600080fd5b50600436106100415760003560e01c8063368b87721461004657806341c0e1b5146100ee578063ce6d41de146100f6575b600080fd5b6100ec6004803603602081101561005c57600080fd5b81019060208101813564010000000081111561007757600080fd5b82018360208201111561008957600080fd5b803590602001918460018302840111640100000000831117156100ab57600080fd5b91908080601f016020809104026020016040519081016040528093929190818152602001838380828437600092019190915250929550610173945050505050565b005b6100ec6101a2565b6100fe6101ba565b6040805160208082528351818301528351919283929083019185019080838360005b83811015610138578181015183820152602001610120565b50505050905090810190601f1680156101655780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b6000546001600160a01b0316331461018a5761019f565b805161019d906001906020840190610250565b505b50565b6000546001600160a01b03163314156101b85733ff5b565b60018054604080516020601f600260001961010087891615020190951694909404938401819004810282018101909252828152606093909290918301828280156102455780601f1061021a57610100808354040283529160200191610245565b820191906000526020600020905b81548152906001019060200180831161022857829003601f168201915b505050505090505b90565b828054600181600116156101000203166002900490600052602060002090601f016020900481019282601f1061029157805160ff19168380011785556102be565b828001600101855582156102be579182015b828111156102be5782518255916020019190600101906102a3565b506102ca9291506102ce565b5090565b61024d91905b808211156102ca57600081556001016102d456fea264697066735822122084964d4c3f6bc912a9d20e14e449721012d625aa3c8a12de41ae5519752fc89064736f6c63430006000033")
	ecdsaPrivateKey, _ := PrivateKeyFromStringECDSA("30540201010420ac318ea8ff8d991ab2f16172b4738e74dc35a56681199cfb1c0cb2e7cb560ffda00706052b8104000aa124032200036843f5cb338bbb4cdb21b0da4ea739d910951d6e8a5f703d313efe31afe788f4")
	aliasAccountId := ecdsaPrivateKey.ToAccountID(0, 0)

	// Create a shallow account for the ECDSA key
	resp, err := NewTransferTransaction().
		AddHbarTransfer(env.Client.GetOperatorAccountID(), NewHbar(-2)).
		AddHbarTransfer(*aliasAccountId, NewHbar(2)).
		Execute(env.Client)

	// Create file with the contract bytecode
	resp, err = NewFileCreateTransaction().
		SetNodeAccountIDs(env.NodeAccountIDs).
		SetKeys(env.OperatorKey.PublicKey()).
		SetContents(smartContractBytecode).
		Execute(env.Client)
	require.NoError(t, err)

	receipt, err := resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	fileID := *receipt.FileID
	assert.NotNil(t, fileID)

	// Create contract to be called by EthereumTransaction
	resp, err = NewContractCreateTransaction().
		SetAdminKey(env.Client.GetOperatorPublicKey()).
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		SetGas(1000000).
		SetConstructorParameters(NewContractFunctionParameters().AddString("hello from hedera")).
		SetBytecodeFileID(fileID).
		SetContractMemo("hedera-sdk-go::TestContractCreateTransaction_Execute").
		Execute(env.Client)
	require.NoError(t, err)

	receipt, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	assert.NotNil(t, receipt.ContractID)
	contractID := *receipt.ContractID

	// Call data for the smart contract
	contractMsg := "setMessage"
	msgPointer := &contractMsg

	// build the RLP list that should be signed with the test ECDSA private key
	list := &EIP1559RLP{
		chainId:        decodeHex(t, "012a"),
		nonce:          []byte{},
		maxPriorityGas: decodeHex(t, "00"),
		maxGas:         decodeHex(t, "d1385c7bf0"),
		gasLimit:       decodeHex(t, "0249f0"),
		to:             decodeHex(t, contractID.ToSolidityAddress()),
		value:          []byte{},
		callData:       NewContractFunctionParameters().AddString("new message")._Build(msgPointer),
		accessList:     [][]byte{},
	}

	bytes, _ := rlp.EncodeToBytes(list)

	// 02 is the type of the transaction EIP1559 and should be concatenated to the RLP by service requirement
	bytesToSign := append(decodeHex(t, "02"), bytes...)
	signedBytes := ecdsaPrivateKey.Sign(bytesToSign)

	// Add signature data to the RLP list for EthereumTransaction submition
	list.recId = decodeHex(t, "01")
	list.r = signedBytes[:32]
	list.s = signedBytes[len(signedBytes)-32:]

	ethereumTransactionData, _ := rlp.EncodeToBytes(list)
	// 02 is the type of the transaction EIP1559 and should be concatenated to the RLP by service requirement
	resp, err = NewEthereumTransaction().SetEthereumData(append(decodeHex(t, "02"), ethereumTransactionData...)).Execute(env.Client)

	require.NoError(t, err)

	record, _ := resp.GetRecord(env.Client)

	assert.Equal(t, int64(1), record.CallResult.SignerNonce)

	resp, err = NewContractDeleteTransaction().
		SetContractID(contractID).
		SetTransferAccountID(env.Client.GetOperatorAccountID()).
		Execute(env.Client)
	require.NoError(t, err)

	_, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	resp, err = NewFileDeleteTransaction().
		SetFileID(fileID).
		Execute(env.Client)
	require.NoError(t, err)

	_, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

}
