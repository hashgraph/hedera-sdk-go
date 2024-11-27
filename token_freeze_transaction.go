package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"github.com/hiero-ledger/hiero-sdk-go/v2/proto/services"
)

// TokenFreezeTransaction
// Freezes transfers of the specified token for the account. Must be signed by the Token's freezeKey.
// If the provided account is not found, the transaction will resolve to INVALID_ACCOUNT_ID.
// If the provided account has been deleted, the transaction will resolve to ACCOUNT_DELETED.
// If the provided token is not found, the transaction will resolve to INVALID_TOKEN_ID.
// If the provided token has been deleted, the transaction will resolve to TOKEN_WAS_DELETED.
// If an Association between the provided token and account is not found, the transaction will
// resolve to TOKEN_NOT_ASSOCIATED_TO_ACCOUNT.
// If no Freeze Key is defined, the transaction will resolve to TOKEN_HAS_NO_FREEZE_KEY.
// Once executed the Account is marked as Frozen and will not be able to receive or send tokens
// unless unfrozen. The operation is idempotent.
type TokenFreezeTransaction struct {
	*Transaction[*TokenFreezeTransaction]
	tokenID   *TokenID
	accountID *AccountID
}

// NewTokenFreezeTransaction creates TokenFreezeTransaction which
// freezes transfers of the specified token for the account. Must be signed by the Token's freezeKey.
// If the provided account is not found, the transaction will resolve to INVALID_ACCOUNT_ID.
// If the provided account has been deleted, the transaction will resolve to ACCOUNT_DELETED.
// If the provided token is not found, the transaction will resolve to INVALID_TOKEN_ID.
// If the provided token has been deleted, the transaction will resolve to TOKEN_WAS_DELETED.
// If an Association between the provided token and account is not found, the transaction will
// resolve to TOKEN_NOT_ASSOCIATED_TO_ACCOUNT.
// If no Freeze Key is defined, the transaction will resolve to TOKEN_HAS_NO_FREEZE_KEY.
// Once executed the Account is marked as Frozen and will not be able to receive or send tokens
// unless unfrozen. The operation is idempotent.
func NewTokenFreezeTransaction() *TokenFreezeTransaction {
	tx := &TokenFreezeTransaction{}
	tx.Transaction = _NewTransaction(tx)

	tx._SetDefaultMaxTransactionFee(NewHbar(30))

	return tx
}

func _TokenFreezeTransactionFromProtobuf(tx Transaction[*TokenFreezeTransaction], pb *services.TransactionBody) TokenFreezeTransaction {
	tokenFreezeTransaction := TokenFreezeTransaction{
		tokenID:   _TokenIDFromProtobuf(pb.GetTokenFreeze().GetToken()),
		accountID: _AccountIDFromProtobuf(pb.GetTokenFreeze().GetAccount()),
	}

	tx.childTransaction = &tokenFreezeTransaction
	tokenFreezeTransaction.Transaction = &tx
	return tokenFreezeTransaction
}

// SetTokenID Sets the token for which this account will be frozen. If token does not exist, transaction results
// in INVALID_TOKEN_ID
func (tx *TokenFreezeTransaction) SetTokenID(tokenID TokenID) *TokenFreezeTransaction {
	tx._RequireNotFrozen()
	tx.tokenID = &tokenID
	return tx
}

// GetTokenID returns the token for which this account will be frozen.
func (tx *TokenFreezeTransaction) GetTokenID() TokenID {
	if tx.tokenID == nil {
		return TokenID{}
	}

	return *tx.tokenID
}

// SetAccountID Sets the account to be frozen
func (tx *TokenFreezeTransaction) SetAccountID(accountID AccountID) *TokenFreezeTransaction {
	tx._RequireNotFrozen()
	tx.accountID = &accountID
	return tx
}

// GetAccountID returns the account to be frozen
func (tx *TokenFreezeTransaction) GetAccountID() AccountID {
	if tx.accountID == nil {
		return AccountID{}
	}

	return *tx.accountID
}

// ----------- Overridden functions ----------------

func (tx TokenFreezeTransaction) getName() string {
	return "TokenFreezeTransaction"
}

func (tx TokenFreezeTransaction) validateNetworkOnIDs(client *Client) error {
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

func (tx TokenFreezeTransaction) build() *services.TransactionBody {
	return &services.TransactionBody{
		TransactionFee:           tx.transactionFee,
		Memo:                     tx.Transaction.memo,
		TransactionValidDuration: _DurationToProtobuf(tx.GetTransactionValidDuration()),
		TransactionID:            tx.transactionID._ToProtobuf(),
		Data: &services.TransactionBody_TokenFreeze{
			TokenFreeze: tx.buildProtoBody(),
		},
	}
}

func (tx TokenFreezeTransaction) buildScheduled() (*services.SchedulableTransactionBody, error) {
	return &services.SchedulableTransactionBody{
		TransactionFee: tx.transactionFee,
		Memo:           tx.Transaction.memo,
		Data: &services.SchedulableTransactionBody_TokenFreeze{
			TokenFreeze: tx.buildProtoBody(),
		},
	}, nil
}

func (tx TokenFreezeTransaction) buildProtoBody() *services.TokenFreezeAccountTransactionBody {
	body := &services.TokenFreezeAccountTransactionBody{}
	if tx.tokenID != nil {
		body.Token = tx.tokenID._ToProtobuf()
	}

	if tx.accountID != nil {
		body.Account = tx.accountID._ToProtobuf()
	}

	return body
}

func (tx TokenFreezeTransaction) getMethod(channel *_Channel) _Method {
	return _Method{
		transaction: channel._GetToken().FreezeTokenAccount,
	}
}

func (tx TokenFreezeTransaction) constructScheduleProtobuf() (*services.SchedulableTransactionBody, error) {
	return tx.buildScheduled()
}

func (tx TokenFreezeTransaction) getBaseTransaction() *Transaction[TransactionInterface] {
	return castFromConcreteToBaseTransaction(tx.Transaction, &tx)
}
