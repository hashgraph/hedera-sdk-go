package hiero

// SPDX-License-Identifier: Apache-2.0

// MirrorNodeContractCallQuery returns a result from EVM gas estimation of read-write operations.
type MirrorNodeContractEstimateGasQuery struct {
	mirrorNodeContractQuery
}

func NewMirrorNodeContractEstimateGasQuery() *MirrorNodeContractEstimateGasQuery {
	query := new(MirrorNodeContractEstimateGasQuery)
	return query
}

// SetContractID sets the contract instance to call.
func (query *MirrorNodeContractEstimateGasQuery) SetContractID(contractID ContractID) *MirrorNodeContractEstimateGasQuery {
	query.setContractID(contractID)
	return query
}

// SetContractEvmAddress sets the 20-byte EVM address of the contract to call.
func (query *MirrorNodeContractEstimateGasQuery) SetContractEvmAddress(contractEvmAddress string) *MirrorNodeContractEstimateGasQuery {
	query.setContractEvmAddress(contractEvmAddress)
	return query
}

// SetSender sets the sender of the transaction simulation.
func (query *MirrorNodeContractEstimateGasQuery) SetSender(sender AccountID) *MirrorNodeContractEstimateGasQuery {
	query.setSender(sender)
	return query
}

// SetSenderEvmAddress sets the 20-byte EVM address of the sender of the transaction simulation.
func (query *MirrorNodeContractEstimateGasQuery) SetSenderEvmAddress(senderEvmAddress string) *MirrorNodeContractEstimateGasQuery {
	query.setSenderEvmAddress(senderEvmAddress)
	return query
}

// SetFunction sets the function parameters as their raw bytes.
func (query *MirrorNodeContractEstimateGasQuery) SetFunction(name string, params *ContractFunctionParameters) *MirrorNodeContractEstimateGasQuery {
	query.setFunction(name, params)
	return query
}

// SetValue sets the amount of value (in tinybars or wei) to be sent to the contract in the transaction.
func (query *MirrorNodeContractEstimateGasQuery) SetValue(value int64) *MirrorNodeContractEstimateGasQuery {
	query.setValue(value)
	return query
}

// SetGasLimit sets the gas limit for the contract call.
// This specifies the maximum amount of gas that the transaction can consume.
func (query *MirrorNodeContractEstimateGasQuery) SetGasLimit(gasLimit int64) *MirrorNodeContractEstimateGasQuery {
	query.setGasLimit(gasLimit)
	return query
}

// SetGasPrice sets the gas price to be used for the contract call.
// This specifies the price of each unit of gas used in the transaction.
func (query *MirrorNodeContractEstimateGasQuery) SetGasPrice(gasPrice int64) *MirrorNodeContractEstimateGasQuery {
	query.setGasPrice(gasPrice)
	return query
}

// Returns gas estimation for the EVM execution
func (query *MirrorNodeContractEstimateGasQuery) Execute(client *Client) (uint64, error) {
	return query.estimateGas(client)
}
