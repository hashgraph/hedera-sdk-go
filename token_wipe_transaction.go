package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"github.com/hiero-ledger/hiero-sdk-go/v2/proto/services"
)

// TokenWipeTransaction
// Wipes the provided amount of tokens from the specified Account. Must be signed by the Token's Wipe key.
// If the provided account is not found, the transaction will resolve to INVALID_ACCOUNT_ID.
// If the provided account has been deleted, the transaction will resolve to ACCOUNT_DELETED.
// If the provided token is not found, the transaction will resolve to INVALID_TOKEN_ID.
// If the provided token has been deleted, the transaction will resolve to TOKEN_WAS_DELETED.
// If an Association between the provided token and account is not found, the transaction will
// resolve to TOKEN_NOT_ASSOCIATED_TO_ACCOUNT.
// If Wipe Key is not present in the Token, transaction results in TOKEN_HAS_NO_WIPE_KEY.
// If the provided account is the Token's Treasury Account, transaction results in
// CANNOT_WIPE_TOKEN_TREASURY_ACCOUNT
// On success, tokens are removed from the account and the total supply of the token is decreased
// by the wiped amount.
//
// The amount provided is in the lowest denomination possible. Example:
// Token A has 2 decimals. In order to wipe 100 tokens from account, one must provide amount of
// 10000. In order to wipe 100.55 tokens, one must provide amount of 10055.
type TokenWipeTransaction struct {
	*Transaction[*TokenWipeTransaction]
	tokenID   *TokenID
	accountID *AccountID
	amount    uint64
	serial    []int64
}

// NewTokenWipeTransaction creates TokenWipeTransaction which
// wipes the provided amount of tokens from the specified Account. Must be signed by the Token's Wipe key.
// If the provided account is not found, the transaction will resolve to INVALID_ACCOUNT_ID.
// If the provided account has been deleted, the transaction will resolve to ACCOUNT_DELETED.
// If the provided token is not found, the transaction will resolve to INVALID_TOKEN_ID.
// If the provided token has been deleted, the transaction will resolve to TOKEN_WAS_DELETED.
// If an Association between the provided token and account is not found, the transaction will
// resolve to TOKEN_NOT_ASSOCIATED_TO_ACCOUNT.
// If Wipe Key is not present in the Token, transaction results in TOKEN_HAS_NO_WIPE_KEY.
// If the provided account is the Token's Treasury Account, transaction results in
// CANNOT_WIPE_TOKEN_TREASURY_ACCOUNT
// On success, tokens are removed from the account and the total supply of the token is decreased
// by the wiped amount.
//
// The amount provided is in the lowest denomination possible. Example:
// Token A has 2 decimals. In order to wipe 100 tokens from account, one must provide amount of
// 10000. In order to wipe 100.55 tokens, one must provide amount of 10055.
func NewTokenWipeTransaction() *TokenWipeTransaction {
	tx := &TokenWipeTransaction{}
	tx.Transaction = _NewTransaction(tx)

	tx._SetDefaultMaxTransactionFee(NewHbar(30))

	return tx
}

func _TokenWipeTransactionFromProtobuf(tx Transaction[*TokenWipeTransaction], pb *services.TransactionBody) TokenWipeTransaction {
	tokenWipeTransaction := TokenWipeTransaction{
		tokenID:   _TokenIDFromProtobuf(pb.GetTokenWipe().GetToken()),
		accountID: _AccountIDFromProtobuf(pb.GetTokenWipe().GetAccount()),
		amount:    pb.GetTokenWipe().Amount,
		serial:    pb.GetTokenWipe().GetSerialNumbers(),
	}

	tx.childTransaction = &tokenWipeTransaction
	tokenWipeTransaction.Transaction = &tx
	return tokenWipeTransaction
}

// SetTokenID Sets the token for which the account will be wiped. If token does not exist, transaction results in
// INVALID_TOKEN_ID
func (tx *TokenWipeTransaction) SetTokenID(tokenID TokenID) *TokenWipeTransaction {
	tx._RequireNotFrozen()
	tx.tokenID = &tokenID
	return tx
}

// GetTokenID returns the TokenID that is being wiped
func (tx *TokenWipeTransaction) GetTokenID() TokenID {
	if tx.tokenID == nil {
		return TokenID{}
	}

	return *tx.tokenID
}

// SetAccountID Sets the account to be wiped
func (tx *TokenWipeTransaction) SetAccountID(accountID AccountID) *TokenWipeTransaction {
	tx._RequireNotFrozen()
	tx.accountID = &accountID
	return tx
}

// GetAccountID returns the AccountID that is being wiped
func (tx *TokenWipeTransaction) GetAccountID() AccountID {
	if tx.accountID == nil {
		return AccountID{}
	}

	return *tx.accountID
}

// SetAmount Sets the amount of tokens to wipe from the specified account. Amount must be a positive non-zero
// number in the lowest denomination possible, not bigger than the token balance of the account
// (0; balance]
func (tx *TokenWipeTransaction) SetAmount(amount uint64) *TokenWipeTransaction {
	tx._RequireNotFrozen()
	tx.amount = amount
	return tx
}

// GetAmount returns the amount of tokens to be wiped from the specified account
func (tx *TokenWipeTransaction) GetAmount() uint64 {
	return tx.amount
}

// GetSerialNumbers returns the list of serial numbers to be wiped.
func (tx *TokenWipeTransaction) GetSerialNumbers() []int64 {
	return tx.serial
}

// SetSerialNumbers
// Sets applicable to tokens of type NON_FUNGIBLE_UNIQUE. The list of serial numbers to be wiped.
func (tx *TokenWipeTransaction) SetSerialNumbers(serial []int64) *TokenWipeTransaction {
	tx._RequireNotFrozen()
	tx.serial = serial
	return tx
}

// ----------- Overridden functions ----------------

func (tx TokenWipeTransaction) getName() string {
	return "TokenWipeTransaction"
}

func (tx TokenWipeTransaction) validateNetworkOnIDs(client *Client) error {
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

func (tx TokenWipeTransaction) build() *services.TransactionBody {
	return &services.TransactionBody{
		TransactionFee:           tx.transactionFee,
		Memo:                     tx.Transaction.memo,
		TransactionValidDuration: _DurationToProtobuf(tx.GetTransactionValidDuration()),
		TransactionID:            tx.transactionID._ToProtobuf(),
		Data: &services.TransactionBody_TokenWipe{
			TokenWipe: tx.buildProtoBody(),
		},
	}
}

func (tx TokenWipeTransaction) buildScheduled() (*services.SchedulableTransactionBody, error) {
	return &services.SchedulableTransactionBody{
		TransactionFee: tx.transactionFee,
		Memo:           tx.Transaction.memo,
		Data: &services.SchedulableTransactionBody_TokenWipe{
			TokenWipe: tx.buildProtoBody(),
		},
	}, nil
}

func (tx TokenWipeTransaction) buildProtoBody() *services.TokenWipeAccountTransactionBody {
	body := &services.TokenWipeAccountTransactionBody{
		Amount: tx.amount,
	}

	if len(tx.serial) > 0 {
		body.SerialNumbers = tx.serial
	}

	if tx.tokenID != nil {
		body.Token = tx.tokenID._ToProtobuf()
	}

	if tx.accountID != nil {
		body.Account = tx.accountID._ToProtobuf()
	}

	return body
}

func (tx TokenWipeTransaction) getMethod(channel *_Channel) _Method {
	return _Method{
		transaction: channel._GetToken().WipeTokenAccount,
	}
}

func (tx TokenWipeTransaction) constructScheduleProtobuf() (*services.SchedulableTransactionBody, error) {
	return tx.buildScheduled()
}

func (tx TokenWipeTransaction) getBaseTransaction() *Transaction[TransactionInterface] {
	return castFromConcreteToBaseTransaction(tx.Transaction, &tx)
}
