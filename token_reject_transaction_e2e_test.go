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

	// given

	// create fungible token with treasury
	// create receiver account with auto associations

	// when

	// transfer ft to the receiver
	// reject the token

	// then

	// verify the balance is 0
	// verify the auto associations are not decremented
}

func TestIntegrationTokenRejectTransactionCanExecuteForNFT(t *testing.T) {
	t.Parallel()
	//env := NewIntegrationTestEnv(t)

	// given

	// create nft with treasury
	// mint
	// create receiver account with auto associations

	// when

	// transfer nft to the receiver
	// reject the token

	// then

	// verify the balance is 0
	// verify the auto associations are not decremented
}

func TestIntegrationTokenRejectTransactionReceiverSigRequired(t *testing.T) {
	t.Parallel()
	//env := NewIntegrationTestEnv(t)

	// given

	// create nft with treasury with receiver sig required
	// mint
	// create receiver account with auto associations

	// when

	// transfer nft to the receiver
	// reject the token

	// then

	// verify the balance is 0
	// verify the auto associations are not decremented

	// same test for fungible token
}

func TestIntegrationTokenRejectTransactionReceiverTokenFrozen(t *testing.T) {
	t.Parallel()
	//env := NewIntegrationTestEnv(t)

	// given

	// create nft with treasury
	// mint
	// create receiver account with auto associations

	// when

	// transfer nft to the receiver
	// freeze the token

	// then

	// reject the token - should fail with TOKEN_IS_FROZEN

	// same test with fungible token
}

func TestIntegrationTokenRejectTransactionReceiverTokenPaused(t *testing.T) {
	t.Parallel()
	//env := NewIntegrationTestEnv(t)

	//given

	// create nft with treasury
	// mint
	// create receiver account with auto associations

	// when

	// transfer nft to the receiver
	// pause the token

	// then

	// reject the token - should fail with TOKEN_IS_PAUSED

	// same test with fungible token
}

func TestIntegrationTokenRejectTransactionRemovesAllowance(t *testing.T) {
	t.Parallel()
	//env := NewIntegrationTestEnv(t)

	// given

	// create fungible token with treasury
	// create receiver account with auto associations
	// create spender account to be approved

	// when

	// transfer ft to the receiver
	// approve allowance to the spender
	// verify the allowance with query
	// reject the token

	// then

	// verify the allowance - should be 0 , because the receiver is no longer the owner

	// same test for nft
}

func TestIntegrationTokenRejectTransactionFailsWithTokenReferenceRepeated(t *testing.T) {
	t.Parallel()
	//env := NewIntegrationTestEnv(t)

	// given

	// create fungible token with treasury
	// create receiver account with auto associations

	// when

	// transfer ft to the receiver
	// reject the token

	// then

	// duplicate the reject - should fail with TOKEN_REFERENCE_REPEATED

	// same test for nft
}

func TestIntegrationTokenRejectTransactionFailsWhenOwnerHasNoBalance(t *testing.T) {
	t.Parallel()
	//env := NewIntegrationTestEnv(t)

	// given

	// create fungible token with treasury
	// create receiver account with auto associations

	// when

	// skip the transfer
	// reject the token - should fail with INSUFFICIENT_TOKEN_BALANCE

	// same test for nft
}

func TestIntegrationTokenRejectTransactionFailsTreasuryRejects(t *testing.T) {
	t.Parallel()
	//env := NewIntegrationTestEnv(t)

	// given

	// create fungible token with treasury

	// when

	// skip the transfer
	// reject the token with the treasury - should fail with ACCOUNT_IS_TREASURY

	// same test for nft
}

func TestIntegrationTokenRejectTransactionFailsWithInvalidOwner(t *testing.T) {
	t.Parallel()
	//env := NewIntegrationTestEnv(t)

	// given

	// create fungible token with treasury

	// when

	// reject the token with invalid owner - should fail with INVALID_OWNER_ID

	// same test for nft
}

func TestIntegrationTokenRejectTransactionFailsWithInvalidToken(t *testing.T) {
	t.Parallel()
	//env := NewIntegrationTestEnv(t)

	// given

	// create receiver account with auto associations

	// when

	// reject the token with invalid token - should fail with INVALID_TOKEN_ID

	// same test for nft
}
