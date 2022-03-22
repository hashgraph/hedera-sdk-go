package hedera

func AccountInfoFlowVerifySignature(client *Client, accountID AccountID, message []byte, signature []byte) (bool, error) {
	info, err := NewAccountInfoQuery().
		SetAccountID(accountID).
		Execute(client)

	if err != nil {
		return false, err
	}

	if key, ok := info.Key.(PublicKey); ok {
		return key.Verify(message, signature), nil
	}

	return false, nil
}

func AccountInfoFlowVerifyTransaction(client *Client, accountID AccountID, transaction Transaction, signature []byte) (bool, error) {
	info, err := NewAccountInfoQuery().
		SetAccountID(accountID).
		Execute(client)

	if err != nil {
		return false, err
	}

	if key, ok := info.Key.(PublicKey); ok {
		return key.VerifyTransaction(transaction), nil
	}

	return false, nil
}
