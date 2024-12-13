package hiero

// SPDX-License-Identifier: Apache-2.0

// MirrorNodeContractCallQuery extends MirrorNodeContractQuery
type MirrorNodeContractCallQuery struct {
	mirrorNodeContractQuery
}

func NewMirrorNodeContractCallQuery() *MirrorNodeContractCallQuery {
	query := new(MirrorNodeContractCallQuery)
	return query
}

// SetContractID sets the contract instance to call.
func (query *MirrorNodeContractCallQuery) SetContractID(contractID ContractID) *MirrorNodeContractCallQuery {
	query.setContractID(contractID)
	return query
}

// SetContractEvmAddress sets the 20-byte EVM address of the contract to call.
func (query *MirrorNodeContractCallQuery) SetContractEvmAddress(contractEvmAddress string) *MirrorNodeContractCallQuery {
	query.setContractEvmAddress(contractEvmAddress)
	return query
}

// SetSender sets the sender of the transaction simulation.
func (query *MirrorNodeContractCallQuery) SetSender(sender AccountID) *MirrorNodeContractCallQuery {
	query.setSender(sender)
	return query
}

// SetSenderEvmAddress sets the 20-byte EVM address of the sender of the transaction simulation.
func (query *MirrorNodeContractCallQuery) SetSenderEvmAddress(senderEvmAddress string) *MirrorNodeContractCallQuery {
	query.setSenderEvmAddress(senderEvmAddress)
	return query
}

// SetFunction sets the function parameters as their raw bytes.
func (query *MirrorNodeContractCallQuery) SetFunction(name string, params *ContractFunctionParameters) *MirrorNodeContractCallQuery {
	query.setFunction(name, params)
	return query
}

// SetValue sets the amount of value (in tinybars or wei) to be sent to the contract in the transaction.
func (query *MirrorNodeContractCallQuery) SetValue(value int64) *MirrorNodeContractCallQuery {
	query.setValue(value)
	return query
}

// SetGasLimit sets the gas limit for the contract call.
// This specifies the maximum amount of gas that the transaction can consume.
func (query *MirrorNodeContractCallQuery) SetGasLimit(gasLimit int64) *MirrorNodeContractCallQuery {
	query.setGasLimit(gasLimit)
	return query
}

// SetGasPrice sets the gas price to be used for the contract call.
// This specifies the price of each unit of gas used in the transaction.
func (query *MirrorNodeContractCallQuery) SetGasPrice(gasPrice int64) *MirrorNodeContractCallQuery {
	query.setGasPrice(gasPrice)
	return query
}

// SetBlockNumber sets the block number for the simulation of the contract call.
// The block number determines the context of the contract call simulation within the blockchain.
func (query *MirrorNodeContractCallQuery) SetBlockNumber(blockNumber int64) *MirrorNodeContractCallQuery {
	query.setBlockNumber(blockNumber)
	return query
}

// Returns gas estimation for the EVM execution
func (query *MirrorNodeContractCallQuery) Execute(client *Client) (string, error) {
	return query.call(client)
}
