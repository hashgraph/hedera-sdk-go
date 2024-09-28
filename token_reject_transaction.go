package hedera

import (
	"github.com/hashgraph/hedera-protobufs-go/services"
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

/**
 * A transaction body to "reject" undesired tokens.<br/>
 * This transaction will transfer one or more tokens or token
 * balances held by the requesting account to the treasury
 * for each token type.
 * <p>
 * Each transfer MUST be one of the following:
 * <ul>
 *   <li>A single non-fungible/unique token.</li>
 *   <li>The full balance held for a fungible/common
 *       token type.</li>
 * </ul>
 * When complete, the requesting account SHALL NOT hold the
 * rejected tokens.<br/>
 * Custom fees and royalties defined for the tokens rejected
 * SHALL NOT be charged for this transaction.
 */
type TokenRejectTransaction struct {
	*Transaction[*TokenRejectTransaction]
	ownerID  *AccountID
	tokenIDs []TokenID
	nftIDs   []NftID
}

func NewTokenRejectTransaction() *TokenRejectTransaction {
	tx := &TokenRejectTransaction{}
	tx.Transaction = _NewTransaction(tx)
	return tx
}

func _TokenRejectTransactionFromProtobuf(pb *services.TransactionBody) *TokenRejectTransaction {
	rejectTransaction := &TokenRejectTransaction{
		ownerID: _AccountIDFromProtobuf(pb.GetTokenReject().Owner),
	}

	for _, rejection := range pb.GetTokenReject().Rejections {
		if rejection.GetFungibleToken() != nil {
			rejectTransaction.tokenIDs = append(rejectTransaction.tokenIDs, *_TokenIDFromProtobuf(rejection.GetFungibleToken()))
		} else if rejection.GetNft() != nil {
			rejectTransaction.nftIDs = append(rejectTransaction.nftIDs, _NftIDFromProtobuf(rejection.GetNft()))
		}
	}

	return rejectTransaction
}

// SetOwnerID Sets the account which owns the tokens to be rejected
func (tx *TokenRejectTransaction) SetOwnerID(ownerID AccountID) *TokenRejectTransaction {
	tx._RequireNotFrozen()
	tx.ownerID = &ownerID
	return tx
}

// GetOwnerID Gets the account which owns the tokens to be rejected
func (tx *TokenRejectTransaction) GetOwnerID() AccountID {
	if tx.ownerID == nil {
		return AccountID{}
	}
	return *tx.ownerID
}

// SetTokenIDs Sets the tokens to be rejected
func (tx *TokenRejectTransaction) SetTokenIDs(ids ...TokenID) *TokenRejectTransaction {
	tx._RequireNotFrozen()
	tx.tokenIDs = make([]TokenID, len(ids))
	copy(tx.tokenIDs, ids)

	return tx
}

// AddTokenID Adds a token to be rejected
func (tx *TokenRejectTransaction) AddTokenID(id TokenID) *TokenRejectTransaction {
	tx._RequireNotFrozen()
	tx.tokenIDs = append(tx.tokenIDs, id)
	return tx
}

// GetTokenIDs Gets the tokens to be rejected
func (tx *TokenRejectTransaction) GetTokenIDs() []TokenID {
	return tx.tokenIDs
}

// SetNftIDs Sets the NFTs to be rejected
func (tx *TokenRejectTransaction) SetNftIDs(ids ...NftID) *TokenRejectTransaction {
	tx._RequireNotFrozen()
	tx.nftIDs = make([]NftID, len(ids))
	copy(tx.nftIDs, ids)

	return tx
}

// AddNftID Adds an NFT to be rejected
func (tx *TokenRejectTransaction) AddNftID(id NftID) *TokenRejectTransaction {
	tx._RequireNotFrozen()
	tx.nftIDs = append(tx.nftIDs, id)
	return tx
}

// GetNftIDs Gets the NFTs to be rejected
func (tx *TokenRejectTransaction) GetNftIDs() []NftID {
	return tx.nftIDs
}

// ----------- Overridden functions ----------------

func (tx *TokenRejectTransaction) getName() string {
	return "TokenRejectTransaction"
}

func (tx *TokenRejectTransaction) validateNetworkOnIDs(client *Client) error {
	if client == nil || !client.autoValidateChecksums {
		return nil
	}

	if tx.ownerID != nil {
		if err := tx.ownerID.ValidateChecksum(client); err != nil {
			return err
		}
	}

	for _, tokenID := range tx.tokenIDs {
		if err := tokenID.ValidateChecksum(client); err != nil {
			return err
		}
	}

	for _, nftID := range tx.nftIDs {
		if err := nftID.TokenID.ValidateChecksum(client); err != nil {
			return err
		}
	}

	return nil
}

func (tx *TokenRejectTransaction) build() *services.TransactionBody {
	body := tx.buildProtoBody()

	return &services.TransactionBody{
		TransactionFee:           tx.transactionFee,
		Memo:                     tx.Transaction.memo,
		TransactionValidDuration: _DurationToProtobuf(tx.GetTransactionValidDuration()),
		TransactionID:            tx.transactionID._ToProtobuf(),
		Data: &services.TransactionBody_TokenReject{
			TokenReject: body,
		},
	}
}

func (tx *TokenRejectTransaction) buildScheduled() (*services.SchedulableTransactionBody, error) {
	return &services.SchedulableTransactionBody{
		TransactionFee: tx.transactionFee,
		Memo:           tx.Transaction.memo,
		Data: &services.SchedulableTransactionBody_TokenReject{
			TokenReject: tx.buildProtoBody(),
		},
	}, nil
}

func (tx *TokenRejectTransaction) buildProtoBody() *services.TokenRejectTransactionBody {
	body := &services.TokenRejectTransactionBody{}

	if tx.ownerID != nil {
		body.Owner = tx.ownerID._ToProtobuf()
	}

	for _, tokenID := range tx.tokenIDs {
		tokenReference := &services.TokenReference_FungibleToken{
			FungibleToken: tokenID._ToProtobuf(),
		}

		body.Rejections = append(body.Rejections, &services.TokenReference{
			TokenIdentifier: tokenReference,
		})
	}

	for _, nftID := range tx.nftIDs {
		tokenReference := &services.TokenReference_Nft{
			Nft: nftID._ToProtobuf(),
		}

		body.Rejections = append(body.Rejections, &services.TokenReference{
			TokenIdentifier: tokenReference,
		})
	}

	return body
}

func (tx *TokenRejectTransaction) getMethod(channel *_Channel) _Method {
	return _Method{
		transaction: channel._GetToken().RejectToken,
	}
}

func (tx *TokenRejectTransaction) constructScheduleProtobuf() (*services.SchedulableTransactionBody, error) {
	return tx.buildScheduled()
}

func (tx *TokenRejectTransaction) getBaseTransaction() *Transaction[TransactionInterface] {
	return castFromConcreteToBaseTransaction(tx.Transaction)
}

func (tx *TokenRejectTransaction) setBaseTransaction(baseTx Transaction[TransactionInterface]) {
	tx.Transaction = castFromBaseToConcreteTransaction[*TokenRejectTransaction](baseTx)
}
