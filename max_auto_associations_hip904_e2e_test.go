//go:build all || e2e
// +build all e2e

package hedera

import (
	"testing"
)

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

// Limited max auto association tests

func TestLimitedMaxAutoAssociationsFungibleTokensFlow(t *testing.T) {
	t.Parallel()

	// create token1 with 1 mil supply
	// create token2 with 1 mil supply

	// account create with 1 max auto associations
	// account update with 1 max auto associations
	// contract create with 1 max auto associations
	// contract update with 1 max auto associations
	// contract create flow with 1 max auto associations

	// transfer token1 to all some tokens

	// transfer token2 to all should fail with NO_REMAINING_AUTOMATIC_ASSOCIATIONS
}

func TestLimitedMaxAutoAssociationsNFTsFlow(t *testing.T) {
	t.Parallel()

	// create 2 NFT collections and mint 10 NFTs for each collection

	// account create with 1 max auto associations
	// account update with 1 max auto associations
	// contract create with 1 max auto associations
	// contract update with 1 max auto associations
	// contract create flow with 1 max auto associations

	// transfer nft1 to all, 2 for each

	// transfer nft2 to all should fail with NO_REMAINING_AUTOMATIC_ASSOCIATIONS
}

// HIP-904 Unlimited max auto association tests

func TestUnlimitedMaxAutoAssociationsExecutes(t *testing.T) {
	t.Parallel()
	// account create
	// account update
	// contract create
	// contract update
	// contract create flow

	// all should execute
}

func TestUnlimitedMaxAutoAssociationsAllowsToTransferFungibleTokens(t *testing.T) {
	t.Parallel()
	// create token1 with 1 mil supply
	// create token2 with 1 mil supply

	// account create
	// account update
	// contract create
	// contract update
	// contract create flow

	// transfer to all some tokens
	// transfer to all some tokens

	// all should execute
}

func TestUnlimitedMaxAutoAssociationsAllowsToTransferFungibleTokensWithDecimals(t *testing.T) {
	t.Parallel()

	// same as the above test but with the transfers are with decimals `AddTokenTransferWithDecimals`
}

func TestUnlimitedMaxAutoAssociationsAllowsToTransferFromFungibleTokens(t *testing.T) {
	t.Parallel()
	// create token1 with 1 mil supply
	// create token2 with 1 mil supply

	// account create
	// account update
	// contract create
	// contract update
	// contract create flow

	// approve to all

	// transferFrom token1 to all some tokens
	// transferFrom token2 to all some tokens

	// all should execute
}

func TestUnlimitedMaxAutoAssociationsAllowsToTransferNFTs(t *testing.T) {
	t.Parallel()
	// create 2 NFT collections and mint 10 NFTs for each collection

	// account create
	// account update
	// contract create
	// contract update
	// contract create flow

	// transfer nft1 to all, 2 for each

	// transfer nft2 to all, 2 for each

	// all should execute
}

func TestUnlimitedMaxAutoAssociationsAllowsToTransferFromNFTs(t *testing.T) {
	t.Parallel()
	// create 2 NFT collections and mint 10 NFTs for each collection

	// account create
	// account update
	// contract create
	// contract update
	// contract create flow

	// approve nft1 to all
	// approve nft2 to all

	// transferFrom nft1 to all, 2 for each

	// transferFrom nft2 to all, 2 for each

	// all should execute
}

func TestUnlimitedMaxAutoAssociationsFailsWithInvalid(t *testing.T) {
	t.Parallel()

	// account create with -2 and with -1000
	// account update with - 2 and with -1000
	// contract create with -2 and with -1000
	// contract update with -2 and with -1000
	// contract create flow with -2 and with -1000

	// all fails with invalid max auto associations
}
