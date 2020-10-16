package hedera

type TransactionResponse struct {
	TransactionID TransactionID
	NodeID        AccountID
	Hash          []byte
}

func (response TransactionResponse) GetReceipt(client *Client) (TransactionReceipt, error) {
	return NewTransactionReceiptQuery().
		SetTransactionID(response.TransactionID).
		SetNodeAccountID(response.NodeID).
		Execute(client)
}
