package hedera

type TransactionResponse struct {
	TransactionID TransactionID
	NodeID        AccountID
	Hash          []byte
}

func (response TransactionResponse) GetReceipt(client *Client) (TransactionReceipt, error) {
	nodeIDs := make([]AccountID, 1)
	nodeIDs[0] = response.NodeID

	return NewTransactionReceiptQuery().
		SetTransactionID(response.TransactionID).
		SetNodeAccountIDs(nodeIDs).
		Execute(client)
}

func (response TransactionResponse) GetRecord(client *Client) (TransactionRecord, error) {
	nodeIDs := make([]AccountID, 1)
	nodeIDs[0] = response.NodeID

	_, err := NewTransactionReceiptQuery().
		SetTransactionID(response.TransactionID).
		SetNodeAccountIDs(nodeIDs).
		Execute(client)

	if err != nil {
		return TransactionRecord{}, err
	}

	return NewTransactionRecordQuery().
		SetTransactionID(response.TransactionID).
		SetNodeAccountIDs(nodeIDs).
		Execute(client)
}
