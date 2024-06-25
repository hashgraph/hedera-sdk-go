//go:build all || e2e
// +build all e2e

package hedera

import "testing"

/*-
 *
 * Hedera Go SDK
 *
 * Copyright (C) 2020 - 2024 Hedera Hashgraph, LLC
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

func TestIntegrationTokenRejectTransactionCanExecuteForFungible(t *testing.T) {
	t.Parallel()
	//env := NewIntegrationTestEnv(t)

	// create fungible token with treasury
	// create receiver account with auto associations

	// transfer ft to the receiver
	// reject the transfer
	// verify the balance is 0
	// verify the auto associations are not decremented
}

func TestIntegrationTokenRejectTransactionCanExecuteForNFT(t *testing.T) {
	t.Parallel()
	//env := NewIntegrationTestEnv(t)

	// create nft with treasury
	// mint
	// create receiver account with auto associations

	// transfer nft to the receiver
	// reject the transfer
	// verify the balance is 0
	// verify the auto associations are not decremented
}

func TestIntegrationTokenRejectTransactionReceiverSignRequired(t *testing.T) {
	t.Parallel()
	//env := NewIntegrationTestEnv(t)

	// create nft with treasury with receiver sign required
	// mint
	// create receiver account with auto associations

	// transfer nft to the receiver
	// reject the transfer
	// verify the balance is 0
	// verify the auto associations are not decremented

	// same test with fungible token
}

func TestIntegrationTokenRejectTransactionReceiverTokenFrozen(t *testing.T) {
	t.Parallel()
	//env := NewIntegrationTestEnv(t)

	// create nft with treasury
	// mint
	// create receiver account with auto associations

	// transfer nft to the receiver
	// freeze the token
	// reject the transfer
	// verify the balance is 0
	// verify the auto associations are not decremented

	// same test with fungible token
}

func TestIntegrationTokenRejectTransactionFailsWithTokenReferenceRepeated(t *testing.T) {
	t.Parallel()
	//env := NewIntegrationTestEnv(t)

	// create fungible token with treasury
	// create receiver account with auto associations

	// transfer ft to the receiver
	// reject the transfer

	// duplicate the reject - should fail with TOKEN_REFERENCE_REPEATED

	// same test for nft
}

func TestIntegrationTokenRejectTransactionFailsWithInvalidOwner(t *testing.T) {
	t.Parallel()
	//env := NewIntegrationTestEnv(t)

	// create fungible token with treasury

	// reject the transfer with invalid owner - should fail with INVALID_OWNER_ID

	// same test for nft
}

func TestIntegrationTokenRejectTransactionFailsWithInvalidToken(t *testing.T) {
	t.Parallel()
	//env := NewIntegrationTestEnv(t)

	// create receiver account with auto associations

	// reject the transfer with invalid token - should fail with INVALID_TOKEN_ID

	// same test for nft
}
