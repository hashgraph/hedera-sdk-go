package hiero

// SPDX-License-Identifier: Apache-2.0

// MirrorNodeContractEstimateGasQuery returns a result from EVM gas estimation of read-write operations.
type MirrorNodeContractEstimateGasQuery struct {
	mirrorNodeContractQuery
}

func NewMirrorNodeContractEstimateGasQuery() *MirrorNodeContractEstimateGasQuery {
	mirrorNodeEstimateGasQuery := new(MirrorNodeContractEstimateGasQuery)
	return mirrorNodeEstimateGasQuery
}

// SetContractID sets the contract instance to call.
func (mirrorNodeEstimateGasQuery *MirrorNodeContractEstimateGasQuery) SetContractID(contractID ContractID) *MirrorNodeContractEstimateGasQuery {
	mirrorNodeEstimateGasQuery.contractID = &contractID
	return mirrorNodeEstimateGasQuery
}

// SetContractEvmAddress sets the 20-byte EVM address of the contract to call.
func (mirrorNodeEstimateGasQuery *MirrorNodeContractEstimateGasQuery) SetContractEvmAddress(contractEvmAddress string) *MirrorNodeContractEstimateGasQuery {
	mirrorNodeEstimateGasQuery.contractEvmAddress = &contractEvmAddress
	mirrorNodeEstimateGasQuery.contractID = nil
	return mirrorNodeEstimateGasQuery
}

// SetSender sets the sender of the transaction simulation.
func (mirrorNodeEstimateGasQuery *MirrorNodeContractEstimateGasQuery) SetSender(sender AccountID) *MirrorNodeContractEstimateGasQuery {
	mirrorNodeEstimateGasQuery.sender = &sender
	return mirrorNodeEstimateGasQuery
}

// SetSenderEvmAddress sets the 20-byte EVM address of the sender of the transaction simulation.
func (mirrorNodeEstimateGasQuery *MirrorNodeContractEstimateGasQuery) SetSenderEvmAddress(senderEvmAddress string) *MirrorNodeContractEstimateGasQuery {
	mirrorNodeEstimateGasQuery.senderEvmAddress = &senderEvmAddress
	mirrorNodeEstimateGasQuery.sender = nil
	return mirrorNodeEstimateGasQuery
}

// SetFunction sets the function parameters as their raw bytes.
func (mirrorNodeEstimateGasQuery *MirrorNodeContractEstimateGasQuery) SetFunction(name string, params *ContractFunctionParameters) *MirrorNodeContractEstimateGasQuery {
	mirrorNodeEstimateGasQuery.setFunction(name, params)
	return mirrorNodeEstimateGasQuery
}

// SetFunction sets the function parameters as their raw bytes.
func (mirrorNodeEstimateGasQuery *MirrorNodeContractEstimateGasQuery) SetFunctionParameters(byteArray []byte) *MirrorNodeContractEstimateGasQuery {
	mirrorNodeEstimateGasQuery.callData = byteArray
	return mirrorNodeEstimateGasQuery
}

// SetValue sets the amount of value (in tinybars or wei) to be sent to the contract in the transaction.
func (mirrorNodeEstimateGasQuery *MirrorNodeContractEstimateGasQuery) SetValue(value int64) *MirrorNodeContractEstimateGasQuery {
	mirrorNodeEstimateGasQuery.value = &value
	return mirrorNodeEstimateGasQuery
}

// SetGasLimit sets the gas limit for the contract call.
// This specifies the maximum amount of gas that the transaction can consume.
func (mirrorNodeEstimateGasQuery *MirrorNodeContractEstimateGasQuery) SetGasLimit(gasLimit int64) *MirrorNodeContractEstimateGasQuery {
	mirrorNodeEstimateGasQuery.gasLimit = &gasLimit
	return mirrorNodeEstimateGasQuery
}

// SetGasPrice sets the gas price to be used for the contract call.
// This specifies the price of each unit of gas used in the transaction.
func (mirrorNodeEstimateGasQuery *MirrorNodeContractEstimateGasQuery) SetGasPrice(gasPrice int64) *MirrorNodeContractEstimateGasQuery {
	mirrorNodeEstimateGasQuery.gasPrice = &gasPrice
	return mirrorNodeEstimateGasQuery
}

// Returns gas estimation for the EVM execution
func (mirrorNodeEstimateGasQuery *MirrorNodeContractEstimateGasQuery) Execute(client *Client) (uint64, error) {
	return mirrorNodeEstimateGasQuery.estimateGas(client)
}
