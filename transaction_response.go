package hedera

type TransactionResponse struct {
	TransactionID          TransactionID
	ScheduledTransactionId TransactionID // nolint
	NodeID                 AccountID
	Hash                   []byte
}

func (response TransactionResponse) GetReceipt(client *Client) (TransactionReceipt, error) {
	receipt, err := NewTransactionReceiptQuery().
		SetTransactionID(response.TransactionID).
		SetNodeAccountIDs([]AccountID{response.NodeID}).
		Execute(client)

	if err != nil {
		return receipt, err
	}

	if receipt.Status != StatusSuccess {
		return receipt, ErrHederaReceiptStatus{
			TxID:    response.TransactionID,
			Status:  receipt.Status,
			Receipt: receipt,
		}
	}

	return receipt, nil
}

func (response TransactionResponse) GetRecord(client *Client) (TransactionRecord, error) {
	_, err := NewTransactionReceiptQuery().
		SetTransactionID(response.TransactionID).
		SetNodeAccountIDs([]AccountID{response.NodeID}).
		Execute(client)

	if err != nil {
		return TransactionRecord{}, err
	}

	return NewTransactionRecordQuery().
		SetTransactionID(response.TransactionID).
		SetNodeAccountIDs([]AccountID{response.NodeID}).
		Execute(client)
}
