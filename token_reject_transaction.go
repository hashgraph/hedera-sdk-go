package hiero

import (
	"github.com/hiero-ledger/hiero-sdk-go/v2/proto/services"
)

// SPDX-License-Identifier: Apache-2.0

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

func _TokenRejectTransactionFromProtobuf(tx Transaction[*TokenRejectTransaction], pb *services.TransactionBody) TokenRejectTransaction {
	rejectTransaction := TokenRejectTransaction{
		ownerID: _AccountIDFromProtobuf(pb.GetTokenReject().Owner),
	}

	for _, rejection := range pb.GetTokenReject().Rejections {
		if rejection.GetFungibleToken() != nil {
			rejectTransaction.tokenIDs = append(rejectTransaction.tokenIDs, *_TokenIDFromProtobuf(rejection.GetFungibleToken()))
		} else if rejection.GetNft() != nil {
			rejectTransaction.nftIDs = append(rejectTransaction.nftIDs, _NftIDFromProtobuf(rejection.GetNft()))
		}
	}

	tx.childTransaction = &rejectTransaction
	rejectTransaction.Transaction = &tx
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

func (tx TokenRejectTransaction) getName() string {
	return "TokenRejectTransaction"
}

func (tx TokenRejectTransaction) validateNetworkOnIDs(client *Client) error {
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

func (tx TokenRejectTransaction) build() *services.TransactionBody {
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

func (tx TokenRejectTransaction) buildScheduled() (*services.SchedulableTransactionBody, error) {
	return &services.SchedulableTransactionBody{
		TransactionFee: tx.transactionFee,
		Memo:           tx.Transaction.memo,
		Data: &services.SchedulableTransactionBody_TokenReject{
			TokenReject: tx.buildProtoBody(),
		},
	}, nil
}

func (tx TokenRejectTransaction) buildProtoBody() *services.TokenRejectTransactionBody {
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

func (tx TokenRejectTransaction) getMethod(channel *_Channel) _Method {
	return _Method{
		transaction: channel._GetToken().RejectToken,
	}
}

func (tx TokenRejectTransaction) constructScheduleProtobuf() (*services.SchedulableTransactionBody, error) {
	return tx.buildScheduled()
}

func (tx TokenRejectTransaction) getBaseTransaction() *Transaction[TransactionInterface] {
	return castFromConcreteToBaseTransaction(tx.Transaction, &tx)
}
