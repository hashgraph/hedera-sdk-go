package hedera

type TransactionResponse struct {
	TransactionID TransactionID
	NodeID        AccountID
	Hash          []byte
}
