package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"github.com/hiero-ledger/hiero-sdk-go/v2/proto/services"
)

// TokenAssociateTransaction Associates the provided account with the provided tokens. Must be signed by the provided Account's key.
// If the provided account is not found, the transaction will resolve to
// INVALID_ACCOUNT_ID.
// If the provided account has been deleted, the transaction will resolve to
// ACCOUNT_DELETED.
// If any of the provided tokens is not found, the transaction will resolve to
// INVALID_TOKEN_REF.
// If any of the provided tokens has been deleted, the transaction will resolve to
// TOKEN_WAS_DELETED.
// If an association between the provided account and any of the tokens already exists, the
// transaction will resolve to
// TOKEN_ALREADY_ASSOCIATED_TO_ACCOUNT.
// If the provided account's associations count exceed the constraint of maximum token
// associations per account, the transaction will resolve to
// TOKENS_PER_ACCOUNT_LIMIT_EXCEEDED.
// On success, associations between the provided account and tokens are made and the account is
// ready to interact with the tokens.
type TokenAssociateTransaction struct {
	*Transaction[*TokenAssociateTransaction]
	accountID *AccountID
	tokens    []TokenID
}

// NewTokenAssociateTransaction creates TokenAssociateTransaction which associates the provided account with the provided tokens.
// Must be signed by the provided Account's key.
// If the provided account is not found, the transaction will resolve to
// INVALID_ACCOUNT_ID.
// If the provided account has been deleted, the transaction will resolve to
// ACCOUNT_DELETED.
// If any of the provided tokens is not found, the transaction will resolve to
// INVALID_TOKEN_REF.
// If any of the provided tokens has been deleted, the transaction will resolve to
// TOKEN_WAS_DELETED.
// If an association between the provided account and any of the tokens already exists, the
// transaction will resolve to
// TOKEN_ALREADY_ASSOCIATED_TO_ACCOUNT.
// If the provided account's associations count exceed the constraint of maximum token
// associations per account, the transaction will resolve to
// TOKENS_PER_ACCOUNT_LIMIT_EXCEEDED.
// On success, associations between the provided account and tokens are made and the account is
// ready to interact with the tokens.
func NewTokenAssociateTransaction() *TokenAssociateTransaction {
	tx := &TokenAssociateTransaction{}
	tx.Transaction = _NewTransaction(tx)

	tx._SetDefaultMaxTransactionFee(NewHbar(5))

	return tx
}

func _TokenAssociateTransactionFromProtobuf(tx Transaction[*TokenAssociateTransaction], pb *services.TransactionBody) TokenAssociateTransaction {
	tokens := make([]TokenID, 0)
	for _, token := range pb.GetTokenAssociate().Tokens {
		if tokenID := _TokenIDFromProtobuf(token); tokenID != nil {
			tokens = append(tokens, *tokenID)
		}
	}

	tokenAssociateTransaction := TokenAssociateTransaction{
		accountID: _AccountIDFromProtobuf(pb.GetTokenAssociate().GetAccount()),
		tokens:    tokens,
	}

	tx.childTransaction = &tokenAssociateTransaction
	tokenAssociateTransaction.Transaction = &tx
	return tokenAssociateTransaction
}

// SetAccountID Sets the account to be associated with the provided tokens
func (tx *TokenAssociateTransaction) SetAccountID(accountID AccountID) *TokenAssociateTransaction {
	tx._RequireNotFrozen()
	tx.accountID = &accountID
	return tx
}

// GetAccountID returns the account to be associated with the provided tokens
func (tx *TokenAssociateTransaction) GetAccountID() AccountID {
	if tx.accountID == nil {
		return AccountID{}
	}

	return *tx.accountID
}

// SetTokenIDs Sets the tokens to be associated with the provided account
func (tx *TokenAssociateTransaction) SetTokenIDs(ids ...TokenID) *TokenAssociateTransaction {
	tx._RequireNotFrozen()
	tx.tokens = make([]TokenID, len(ids))
	copy(tx.tokens, ids)

	return tx
}

// AddTokenID Adds the token to a token list to be associated with the provided account
func (tx *TokenAssociateTransaction) AddTokenID(id TokenID) *TokenAssociateTransaction {
	tx._RequireNotFrozen()
	if tx.tokens == nil {
		tx.tokens = make([]TokenID, 0)
	}

	tx.tokens = append(tx.tokens, id)

	return tx
}

// GetTokenIDs returns the tokens to be associated with the provided account
func (tx *TokenAssociateTransaction) GetTokenIDs() []TokenID {
	tokenIDs := make([]TokenID, len(tx.tokens))
	copy(tokenIDs, tx.tokens)

	return tokenIDs
}

// ----------- Overridden functions ----------------

func (tx TokenAssociateTransaction) getName() string {
	return "TokenAssociateTransaction"
}

func (tx TokenAssociateTransaction) validateNetworkOnIDs(client *Client) error {
	if client == nil || !client.autoValidateChecksums {
		return nil
	}

	if tx.accountID != nil {
		if err := tx.accountID.ValidateChecksum(client); err != nil {
			return err
		}
	}

	for _, tokenID := range tx.tokens {
		if err := tokenID.ValidateChecksum(client); err != nil {
			return err
		}
	}

	return nil
}

func (tx TokenAssociateTransaction) build() *services.TransactionBody {
	body := tx.buildProtoBody()

	return &services.TransactionBody{
		TransactionFee:           tx.transactionFee,
		Memo:                     tx.Transaction.memo,
		TransactionValidDuration: _DurationToProtobuf(tx.GetTransactionValidDuration()),
		TransactionID:            tx.transactionID._ToProtobuf(),
		Data: &services.TransactionBody_TokenAssociate{
			TokenAssociate: body,
		},
	}
}

func (tx TokenAssociateTransaction) buildScheduled() (*services.SchedulableTransactionBody, error) {
	return &services.SchedulableTransactionBody{
		TransactionFee: tx.transactionFee,
		Memo:           tx.Transaction.memo,
		Data: &services.SchedulableTransactionBody_TokenAssociate{
			TokenAssociate: tx.buildProtoBody(),
		},
	}, nil
}

func (tx TokenAssociateTransaction) buildProtoBody() *services.TokenAssociateTransactionBody {
	body := &services.TokenAssociateTransactionBody{}
	if tx.accountID != nil {
		body.Account = tx.accountID._ToProtobuf()
	}

	if len(tx.tokens) > 0 {
		for _, tokenID := range tx.tokens {
			if body.Tokens == nil {
				body.Tokens = make([]*services.TokenID, 0)
			}
			body.Tokens = append(body.Tokens, tokenID._ToProtobuf())
		}
	}
	return body
}

func (tx TokenAssociateTransaction) getMethod(channel *_Channel) _Method {
	return _Method{
		transaction: channel._GetToken().AssociateTokens,
	}
}

func (tx TokenAssociateTransaction) constructScheduleProtobuf() (*services.SchedulableTransactionBody, error) {
	return tx.buildScheduled()
}

func (tx TokenAssociateTransaction) getBaseTransaction() *Transaction[TransactionInterface] {
	return castFromConcreteToBaseTransaction(tx.Transaction, &tx)
}
