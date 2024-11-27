package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"github.com/hiero-ledger/hiero-sdk-go/v2/proto/services"
)

// TokenDeleteTransaction
// Marks a token as deleted, though it will remain in the ledger.
// The operation must be signed by the specified Admin Key of the Token. If
// admin key is not set, transaction will result in TOKEN_IS_IMMUTABlE.
// Once deleted update, mint, burn, wipe, freeze, unfreeze, grant kyc, revoke
// kyc and token transfer transactions will resolve to TOKEN_WAS_DELETED.
type TokenDeleteTransaction struct {
	*Transaction[*TokenDeleteTransaction]
	tokenID *TokenID
}

// NewTokenDeleteTransaction creates TokenDeleteTransaction which marks a token as deleted,
// though it will remain in the ledger.
// The operation must be signed by the specified Admin Key of the Token. If
// admin key is not set, Transaction will result in TOKEN_IS_IMMUTABlE.
// Once deleted update, mint, burn, wipe, freeze, unfreeze, grant kyc, revoke
// kyc and token transfer transactions will resolve to TOKEN_WAS_DELETED.
func NewTokenDeleteTransaction() *TokenDeleteTransaction {
	tx := &TokenDeleteTransaction{}
	tx.Transaction = _NewTransaction(tx)

	tx._SetDefaultMaxTransactionFee(NewHbar(30))

	return tx
}

func _TokenDeleteTransactionFromProtobuf(tx Transaction[*TokenDeleteTransaction], pb *services.TransactionBody) TokenDeleteTransaction {
	tokenDeleteTransaction := TokenDeleteTransaction{
		tokenID: _TokenIDFromProtobuf(pb.GetTokenDeletion().GetToken()),
	}

	tx.childTransaction = &tokenDeleteTransaction
	tokenDeleteTransaction.Transaction = &tx
	return tokenDeleteTransaction
}

// SetTokenID Sets the Token to be deleted
func (tx *TokenDeleteTransaction) SetTokenID(tokenID TokenID) *TokenDeleteTransaction {
	tx._RequireNotFrozen()
	tx.tokenID = &tokenID
	return tx
}

// GetTokenID returns the TokenID of the token to be deleted
func (tx *TokenDeleteTransaction) GetTokenID() TokenID {
	if tx.tokenID == nil {
		return TokenID{}
	}

	return *tx.tokenID
}

// ----------- Overridden functions ----------------

func (tx TokenDeleteTransaction) getName() string {
	return "TokenDeleteTransaction"
}

func (tx TokenDeleteTransaction) validateNetworkOnIDs(client *Client) error {
	if client == nil || !client.autoValidateChecksums {
		return nil
	}

	if tx.tokenID != nil {
		if err := tx.tokenID.ValidateChecksum(client); err != nil {
			return err
		}
	}

	return nil
}

func (tx TokenDeleteTransaction) build() *services.TransactionBody {
	return &services.TransactionBody{
		TransactionFee:           tx.transactionFee,
		Memo:                     tx.Transaction.memo,
		TransactionValidDuration: _DurationToProtobuf(tx.GetTransactionValidDuration()),
		TransactionID:            tx.transactionID._ToProtobuf(),
		Data: &services.TransactionBody_TokenDeletion{
			TokenDeletion: tx.buildProtoBody(),
		},
	}
}

func (tx TokenDeleteTransaction) buildScheduled() (*services.SchedulableTransactionBody, error) {
	return &services.SchedulableTransactionBody{
		TransactionFee: tx.transactionFee,
		Memo:           tx.Transaction.memo,
		Data: &services.SchedulableTransactionBody_TokenDeletion{
			TokenDeletion: tx.buildProtoBody(),
		},
	}, nil
}

func (tx TokenDeleteTransaction) buildProtoBody() *services.TokenDeleteTransactionBody {
	body := &services.TokenDeleteTransactionBody{}
	if tx.tokenID != nil {
		body.Token = tx.tokenID._ToProtobuf()
	}

	return body
}

func (tx TokenDeleteTransaction) getMethod(channel *_Channel) _Method {
	return _Method{
		transaction: channel._GetToken().DeleteToken,
	}
}

func (tx TokenDeleteTransaction) constructScheduleProtobuf() (*services.SchedulableTransactionBody, error) {
	return tx.buildScheduled()
}

func (tx TokenDeleteTransaction) getBaseTransaction() *Transaction[TransactionInterface] {
	return castFromConcreteToBaseTransaction(tx.Transaction, &tx)
}
