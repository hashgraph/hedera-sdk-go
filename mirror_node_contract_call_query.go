package hiero

// SPDX-License-Identifier: Apache-2.0

// MirrorNodeContractCallQuery returns a result from EVM transient simulation of read-write operations.
type MirrorNodeContractCallQuery struct {
	mirrorNodeContractQuery
}

func NewMirrorNodeContractCallQuery() *MirrorNodeContractCallQuery {
	mirrorNodeContractCallQuery := new(MirrorNodeContractCallQuery)
	return mirrorNodeContractCallQuery
}

// SetContractID sets the contract instance to call.
func (mirrorNodeContractCallQuery *MirrorNodeContractCallQuery) SetContractID(contractID ContractID) *MirrorNodeContractCallQuery {
	mirrorNodeContractCallQuery.contractID = &contractID
	return mirrorNodeContractCallQuery
}

// SetContractEvmAddress sets the 20-byte EVM address of the contract to call.
func (mirrorNodeContractCallQuery *MirrorNodeContractCallQuery) SetContractEvmAddress(contractEvmAddress string) *MirrorNodeContractCallQuery {
	mirrorNodeContractCallQuery.contractEvmAddress = &contractEvmAddress
	mirrorNodeContractCallQuery.contractID = nil
	return mirrorNodeContractCallQuery
}

// SetSender sets the sender of the transaction simulation.
func (mirrorNodeContractCallQuery *MirrorNodeContractCallQuery) SetSender(sender AccountID) *MirrorNodeContractCallQuery {
	mirrorNodeContractCallQuery.sender = &sender
	return mirrorNodeContractCallQuery
}

// SetSenderEvmAddress sets the 20-byte EVM address of the sender of the transaction simulation.
func (mirrorNodeContractCallQuery *MirrorNodeContractCallQuery) SetSenderEvmAddress(senderEvmAddress string) *MirrorNodeContractCallQuery {
	mirrorNodeContractCallQuery.senderEvmAddress = &senderEvmAddress
	mirrorNodeContractCallQuery.sender = nil
	return mirrorNodeContractCallQuery
}

// SetFunction sets the function parameters as their raw bytes.
func (mirrorNodeContractCallQuery *MirrorNodeContractCallQuery) SetFunction(name string, params *ContractFunctionParameters) *MirrorNodeContractCallQuery {
	mirrorNodeContractCallQuery.setFunction(name, params)
	return mirrorNodeContractCallQuery
}

// SetFunction sets the function parameters as their raw bytes.
func (mirrorNodeContractCallQuery *MirrorNodeContractCallQuery) SetFunctionParameters(byteArray []byte) *MirrorNodeContractCallQuery {
	mirrorNodeContractCallQuery.callData = byteArray
	return mirrorNodeContractCallQuery
}

// SetValue sets the amount of value (in tinybars or wei) to be sent to the contract in the transaction.
func (mirrorNodeContractCallQuery *MirrorNodeContractCallQuery) SetValue(value int64) *MirrorNodeContractCallQuery {
	mirrorNodeContractCallQuery.value = &value
	return mirrorNodeContractCallQuery
}

// SetGasLimit sets the gas limit for the contract call.
// This specifies the maximum amount of gas that the transaction can consume.
func (mirrorNodeContractCallQuery *MirrorNodeContractCallQuery) SetGasLimit(gasLimit int64) *MirrorNodeContractCallQuery {
	mirrorNodeContractCallQuery.gasLimit = &gasLimit
	return mirrorNodeContractCallQuery
}

// SetGasPrice sets the gas price to be used for the contract call.
// This specifies the price of each unit of gas used in the transaction.
func (mirrorNodeContractCallQuery *MirrorNodeContractCallQuery) SetGasPrice(gasPrice int64) *MirrorNodeContractCallQuery {
	mirrorNodeContractCallQuery.gasPrice = &gasPrice
	return mirrorNodeContractCallQuery
}

// SetBlockNumber sets the block number for the simulation of the contract call.
// The block number determines the context of the contract call simulation within the blockchain.
func (mirrorNodeContractCallQuery *MirrorNodeContractCallQuery) SetBlockNumber(blockNumber int64) *MirrorNodeContractCallQuery {
	mirrorNodeContractCallQuery.blockNumber = &blockNumber
	return mirrorNodeContractCallQuery
}

// Does transient simulation of read-write operations and returns the result in hexadecimal string format.
func (mirrorNodeContractCallQuery *MirrorNodeContractCallQuery) Execute(client *Client) (string, error) {
	return mirrorNodeContractCallQuery.call(client)
}
