package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"github.com/hiero-ledger/hiero-sdk-go/v2/proto/services"
)

// TokenUnfreezeTransaction
// Unfreezes transfers of the specified token for the account. Must be signed by the Token's freezeKey.
// If the provided account is not found, the transaction will resolve to INVALID_ACCOUNT_ID.
// If the provided account has been deleted, the transaction will resolve to ACCOUNT_DELETED.
// If the provided token is not found, the transaction will resolve to INVALID_TOKEN_ID.
// If the provided token has been deleted, the transaction will resolve to TOKEN_WAS_DELETED.
// If an Association between the provided token and account is not found, the transaction will
// resolve to TOKEN_NOT_ASSOCIATED_TO_ACCOUNT.
// If no Freeze Key is defined, the transaction will resolve to TOKEN_HAS_NO_FREEZE_KEY.
// Once executed the Account is marked as Unfrozen and will be able to receive or send tokens. The
// operation is idempotent.
type TokenUnfreezeTransaction struct {
	*Transaction[*TokenUnfreezeTransaction]
	tokenID   *TokenID
	accountID *AccountID
}

// NewTokenUnfreezeTransaction creates TokenUnfreezeTransaction which
// unfreezes transfers of the specified token for the account. Must be signed by the Token's freezeKey.
// If the provided account is not found, the transaction will resolve to INVALID_ACCOUNT_ID.
// If the provided account has been deleted, the transaction will resolve to ACCOUNT_DELETED.
// If the provided token is not found, the transaction will resolve to INVALID_TOKEN_ID.
// If the provided token has been deleted, the transaction will resolve to TOKEN_WAS_DELETED.
// If an Association between the provided token and account is not found, the transaction will
// resolve to TOKEN_NOT_ASSOCIATED_TO_ACCOUNT.
// If no Freeze Key is defined, the transaction will resolve to TOKEN_HAS_NO_FREEZE_KEY.
// Once executed the Account is marked as Unfrozen and will be able to receive or send tokens. The
// operation is idempotent.
func NewTokenUnfreezeTransaction() *TokenUnfreezeTransaction {
	tx := &TokenUnfreezeTransaction{}
	tx.Transaction = _NewTransaction(tx)

	tx._SetDefaultMaxTransactionFee(NewHbar(30))

	return tx
}

func _TokenUnfreezeTransactionFromProtobuf(tx Transaction[*TokenUnfreezeTransaction], pb *services.TransactionBody) TokenUnfreezeTransaction {
	tokenUnfreezeTransaction := TokenUnfreezeTransaction{
		Transaction: &tx,
		tokenID:     _TokenIDFromProtobuf(pb.GetTokenUnfreeze().GetToken()),
		accountID:   _AccountIDFromProtobuf(pb.GetTokenUnfreeze().GetAccount()),
	}

	tx.childTransaction = &tokenUnfreezeTransaction
	tokenUnfreezeTransaction.Transaction = &tx
	return tokenUnfreezeTransaction
}

// SetTokenID Sets the token for which this account will be unfrozen.
// If token does not exist, transaction results in INVALID_TOKEN_ID
func (tx *TokenUnfreezeTransaction) SetTokenID(tokenID TokenID) *TokenUnfreezeTransaction {
	tx._RequireNotFrozen()
	tx.tokenID = &tokenID
	return tx
}

// GetTokenID returns the token for which this account will be unfrozen.
func (tx *TokenUnfreezeTransaction) GetTokenID() TokenID {
	if tx.tokenID == nil {
		return TokenID{}
	}

	return *tx.tokenID
}

// SetAccountID Sets the account to be unfrozen
func (tx *TokenUnfreezeTransaction) SetAccountID(accountID AccountID) *TokenUnfreezeTransaction {
	tx._RequireNotFrozen()
	tx.accountID = &accountID
	return tx
}

// GetAccountID returns the account to be unfrozen
func (tx *TokenUnfreezeTransaction) GetAccountID() AccountID {
	if tx.accountID == nil {
		return AccountID{}
	}

	return *tx.accountID
}

// ----------- Overridden functions ----------------

func (tx TokenUnfreezeTransaction) getName() string {
	return "TokenUnfreezeTransaction"
}

func (tx TokenUnfreezeTransaction) validateNetworkOnIDs(client *Client) error {
	if client == nil || !client.autoValidateChecksums {
		return nil
	}

	if tx.tokenID != nil {
		if err := tx.tokenID.ValidateChecksum(client); err != nil {
			return err
		}
	}

	if tx.accountID != nil {
		if err := tx.accountID.ValidateChecksum(client); err != nil {
			return err
		}
	}

	return nil
}

func (tx TokenUnfreezeTransaction) build() *services.TransactionBody {
	return &services.TransactionBody{
		TransactionFee:           tx.transactionFee,
		Memo:                     tx.Transaction.memo,
		TransactionValidDuration: _DurationToProtobuf(tx.GetTransactionValidDuration()),
		TransactionID:            tx.transactionID._ToProtobuf(),
		Data: &services.TransactionBody_TokenUnfreeze{
			TokenUnfreeze: tx.buildProtoBody(),
		},
	}
}

func (tx TokenUnfreezeTransaction) buildScheduled() (*services.SchedulableTransactionBody, error) {
	return &services.SchedulableTransactionBody{
		TransactionFee: tx.transactionFee,
		Memo:           tx.Transaction.memo,
		Data: &services.SchedulableTransactionBody_TokenUnfreeze{
			TokenUnfreeze: tx.buildProtoBody(),
		},
	}, nil
}

func (tx TokenUnfreezeTransaction) buildProtoBody() *services.TokenUnfreezeAccountTransactionBody {
	body := &services.TokenUnfreezeAccountTransactionBody{}
	if tx.tokenID != nil {
		body.Token = tx.tokenID._ToProtobuf()
	}

	if tx.accountID != nil {
		body.Account = tx.accountID._ToProtobuf()
	}

	return body
}

func (tx TokenUnfreezeTransaction) getMethod(channel *_Channel) _Method {
	return _Method{
		transaction: channel._GetToken().UnfreezeTokenAccount,
	}
}

func (tx TokenUnfreezeTransaction) constructScheduleProtobuf() (*services.SchedulableTransactionBody, error) {
	return tx.buildScheduled()
}

func (tx TokenUnfreezeTransaction) getBaseTransaction() *Transaction[TransactionInterface] {
	return castFromConcreteToBaseTransaction(tx.Transaction, &tx)
}
