package hiero

// SPDX-License-Identifier: Apache-2.0

// AccountInfoFlowVerifySignature Verifies signature using AccountInfoQuery
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

// AccountInfoFlowVerifyTransaction Verifies transaction using AccountInfoQuery
func AccountInfoFlowVerifyTransaction(client *Client, accountID AccountID, tx TransactionInterface, _ []byte) (bool, error) {
	info, err := NewAccountInfoQuery().
		SetAccountID(accountID).
		Execute(client)

	if err != nil {
		return false, err
	}

	if key, ok := info.Key.(PublicKey); ok {
		return key.VerifyTransaction(tx), nil
	}

	return false, nil
}
