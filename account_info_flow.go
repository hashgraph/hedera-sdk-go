package hedera

/*-
 *
 * Hedera Go SDK
 *
 * Copyright (C) 2020 - 2023 Hedera Hashgraph, LLC
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

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
