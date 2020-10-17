package hedera

import (
	"time"
)

var mockPrivateKey string = "302e020100300506032b6570042204203b054fade7a2b0869c6bd4a63b7017cbae7855d12acc357bea718e2c3e805962"

var testTransactionID TransactionID = TransactionID{
	AccountID{Account: 3},
	time.Unix(124124, 151515),
}

func newMockClient() (*Client, error) {
	privateKey, err := PrivateKeyFromString(mockPrivateKey)

	if err != nil {
		return nil, err
	}

	client := newClient(map[AccountID]string{
		{
			Shard:   0,
			Realm:   0,
			Account: 3,
		}: "nonexistent-testnet",
	}, []string{""})

	client.SetOperator(AccountID{Account: 2}, privateKey)

	return client, nil
}

func newMockTransaction() (Transaction, error) {
	privateKey, err := PrivateKeyFromString(mockPrivateKey)

	if err != nil {
		return Transaction{}, err
	}

	client, err := newMockClient()

	if err != nil {
		return Transaction{}, err
	}

	tx, err := NewCryptoTransferTransaction().
		AddSender(AccountID{Account: 2}, HbarFromTinybar(100)).
		AddRecipient(AccountID{Account: 3}, HbarFromTinybar(100)).
		SetMaxTransactionFee(HbarFrom(1, HbarUnits.Hbar)).
		SetTransactionID(testTransactionID).
		FreezeWith(client)

	if err != nil {
		return Transaction{}, err
	}

	tx.Sign(privateKey)

	return tx.Transaction, nil
}
