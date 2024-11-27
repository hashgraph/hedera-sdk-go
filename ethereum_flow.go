package hiero

// SPDX-License-Identifier: Apache-2.0

import "github.com/pkg/errors"

// Execute an Ethereum transaction on Hiero
type EthereumFlow struct {
	ethereumData    *EthereumTransactionData
	callDataFileID  *FileID
	maxGasAllowance *Hbar
	nodeAccountIDs  []AccountID
}

// Execute an Ethereum transaction on Hiero
func NewEthereumFlow() *EthereumFlow {
	tx := EthereumFlow{}

	return &tx
}

// SetEthereumData sets the raw Ethereum transaction.
func (transaction *EthereumFlow) SetEthereumData(data *EthereumTransactionData) *EthereumFlow {
	transaction.ethereumData = data
	return transaction
}

// SetEthereumDataBytes sets the raw Ethereum transaction.
func (transaction *EthereumFlow) SetEthereumDataBytes(data []byte) *EthereumFlow {
	temp, err := EthereumTransactionDataFromBytes(data)
	if err != nil {
		panic(err)
	}
	transaction.ethereumData = temp
	return transaction
}

// GetEthreumData  returns the data of the Ethereum transaction
func (transaction *EthereumFlow) GetEthereumData() *EthereumTransactionData {
	return transaction.ethereumData
}

// SetCallDataFileID sets the file ID containing the call data.
func (transaction *EthereumFlow) SetCallDataFileID(callData FileID) *EthereumFlow {
	transaction.callDataFileID = &callData
	return transaction
}

// GetCallDataFileID returns the file ID containing the call data.
func (transaction *EthereumFlow) GetCallDataFileID() FileID {
	if transaction.callDataFileID == nil {
		return FileID{}
	}

	return *transaction.callDataFileID
}

// SetMaxGasAllowance sets the maximum gas allowance for the transaction.
func (transaction *EthereumFlow) SetMaxGasAllowance(max Hbar) *EthereumFlow {
	transaction.maxGasAllowance = &max
	return transaction
}

// GetMaxGasAllowance returns the maximum gas allowance for the transaction.
func (transaction *EthereumFlow) GetMaxGasAllowance() Hbar {
	if transaction.maxGasAllowance == nil {
		return Hbar{}
	}

	return *transaction.maxGasAllowance
}

// SetNodeAccountIDs sets the node account IDs for this Ethereum transaction.
func (transaction *EthereumFlow) SetNodeAccountIDs(nodes []AccountID) *EthereumFlow {
	transaction.nodeAccountIDs = nodes
	return transaction
}

// GetNodeAccountIDs returns the node account IDs for this Ethereum transaction.
func (transaction *EthereumFlow) GetNodeAccountIDs() []AccountID {
	return transaction.nodeAccountIDs
}

func (transaction *EthereumFlow) _CreateFile(callData []byte, client *Client) (FileID, error) {
	fileCreate := NewFileCreateTransaction()
	if len(transaction.nodeAccountIDs) > 0 {
		fileCreate.SetNodeAccountIDs(transaction.nodeAccountIDs)
	}

	if len(callData) < 4097 {
		resp, err := fileCreate.
			SetContents(callData).
			Execute(client)
		if err != nil {
			return FileID{}, err
		}

		receipt, err := resp.GetReceipt(client)
		if err != nil {
			return FileID{}, err
		}

		return *receipt.FileID, nil
	}

	resp, err := fileCreate.
		SetContents(callData[:4097]).
		Execute(client)
	if err != nil {
		return FileID{}, err
	}

	receipt, err := resp.GetReceipt(client)
	if err != nil {
		return FileID{}, err
	}

	fileID := *receipt.FileID

	resp, err = NewFileAppendTransaction().
		SetFileID(fileID).
		SetContents(callData[4097:]).
		Execute(client)
	if err != nil {
		return FileID{}, err
	}

	_, err = resp.GetReceipt(client)
	if err != nil {
		return FileID{}, err
	}

	return fileID, nil
}

// Execute executes the Transaction with the provided client
func (transaction *EthereumFlow) Execute(client *Client) (TransactionResponse, error) {
	if transaction.ethereumData == nil {
		return TransactionResponse{}, errors.New("cannot submit ethereum transaction with no ethereum data")
	}

	ethereumTransaction := NewEthereumTransaction()
	if len(transaction.nodeAccountIDs) > 0 {
		ethereumTransaction.SetNodeAccountIDs(transaction.nodeAccountIDs)
	}
	dataBytes, err := transaction.ethereumData.ToBytes()
	if err != nil {
		return TransactionResponse{}, err
	}

	if transaction.maxGasAllowance != nil {
		ethereumTransaction.SetMaxGasAllowanceHbar(*transaction.maxGasAllowance)
	}

	if transaction.callDataFileID != nil { //nolint
		if len(transaction.ethereumData._GetData()) != 0 {
			return TransactionResponse{}, errors.New("call data file ID provided, but ethereum data already contains call data")
		}

		ethereumTransaction.
			SetEthereumData(dataBytes).
			SetCallDataFileID(*transaction.callDataFileID)
	} else if len(dataBytes) <= 5120 {
		ethereumTransaction.
			SetEthereumData(dataBytes)
	} else {
		fileID, err := transaction.
			_CreateFile(dataBytes, client)
		if err != nil {
			return TransactionResponse{}, err
		}

		transaction.ethereumData._SetData([]byte{})

		ethereumTransaction.
			SetEthereumData(dataBytes).
			SetCallDataFileID(fileID)
	}

	resp, err := ethereumTransaction.
		Execute(client)
	if err != nil {
		return TransactionResponse{}, err
	}

	_, err = resp.GetReceipt(client)
	if err != nil {
		return TransactionResponse{}, err
	}

	return resp, nil
}
