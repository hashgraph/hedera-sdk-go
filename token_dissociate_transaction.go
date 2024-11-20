package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"github.com/hiero-ledger/hiero-sdk-go/v2/proto/services"
)

// TokenDissociateTransaction
// Dissociates the provided account with the provided tokens. Must be signed by the provided Account's key.
// If the provided account is not found, the transaction will resolve to INVALID_ACCOUNT_ID.
// If the provided account has been deleted, the transaction will resolve to ACCOUNT_DELETED.
// If any of the provided tokens is not found, the transaction will resolve to INVALID_TOKEN_REF.
// If any of the provided tokens has been deleted, the transaction will resolve to TOKEN_WAS_DELETED.
// If an association between the provided account and any of the tokens does not exist, the
// transaction will resolve to TOKEN_NOT_ASSOCIATED_TO_ACCOUNT.
// If a token has not been deleted and has not expired, and the user has a nonzero balance, the
// transaction will resolve to TRANSACTION_REQUIRES_ZERO_TOKEN_BALANCES.
// If a <b>fungible token</b> has expired, the user can disassociate even if their token balance is
// not zero.
// If a <b>non fungible token</b> has expired, the user can <b>not</b> disassociate if their token
// balance is not zero. The transaction will resolve to TRANSACTION_REQUIRED_ZERO_TOKEN_BALANCES.
// On success, associations between the provided account and tokens are removed.
type TokenDissociateTransaction struct {
	*Transaction[*TokenDissociateTransaction]
	accountID *AccountID
	tokens    []TokenID
}

// NewTokenDissociateTransaction creates TokenDissociateTransaction which
// dissociates the provided account with the provided tokens. Must be signed by the provided Account's key.
// If the provided account is not found, the transaction will resolve to INVALID_ACCOUNT_ID.
// If the provided account has been deleted, the transaction will resolve to ACCOUNT_DELETED.
// If any of the provided tokens is not found, the transaction will resolve to INVALID_TOKEN_REF.
// If any of the provided tokens has been deleted, the transaction will resolve to TOKEN_WAS_DELETED.
// If an association between the provided account and any of the tokens does not exist, the
// transaction will resolve to TOKEN_NOT_ASSOCIATED_TO_ACCOUNT.
// If a token has not been deleted and has not expired, and the user has a nonzero balance, the
// transaction will resolve to TRANSACTION_REQUIRES_ZERO_TOKEN_BALANCES.
// If a <b>fungible token</b> has expired, the user can disassociate even if their token balance is
// not zero.
// If a <b>non fungible token</b> has expired, the user can <b>not</b> disassociate if their token
// balance is not zero. The transaction will resolve to TRANSACTION_REQUIRED_ZERO_TOKEN_BALANCES.
// On success, associations between the provided account and tokens are removed.
func NewTokenDissociateTransaction() *TokenDissociateTransaction {
	tx := &TokenDissociateTransaction{}
	tx.Transaction = _NewTransaction(tx)

	tx._SetDefaultMaxTransactionFee(NewHbar(5))

	return tx
}

func _TokenDissociateTransactionFromProtobuf(tx Transaction[*TokenDissociateTransaction], pb *services.TransactionBody) TokenDissociateTransaction {
	tokens := make([]TokenID, 0)
	for _, token := range pb.GetTokenDissociate().Tokens {
		if tokenID := _TokenIDFromProtobuf(token); tokenID != nil {
			tokens = append(tokens, *tokenID)
		}
	}

	tokenDissociateTransaction := TokenDissociateTransaction{
		accountID: _AccountIDFromProtobuf(pb.GetTokenDissociate().GetAccount()),
		tokens:    tokens,
	}

	tx.childTransaction = &tokenDissociateTransaction
	tokenDissociateTransaction.Transaction = &tx
	return tokenDissociateTransaction
}

// SetAccountID Sets the account to be dissociated with the provided tokens
func (tx *TokenDissociateTransaction) SetAccountID(accountID AccountID) *TokenDissociateTransaction {
	tx._RequireNotFrozen()
	tx.accountID = &accountID
	return tx
}

func (tx *TokenDissociateTransaction) GetAccountID() AccountID {
	if tx.accountID == nil {
		return AccountID{}
	}

	return *tx.accountID
}

// SetTokenIDs Sets the tokens to be dissociated with the provided account
func (tx *TokenDissociateTransaction) SetTokenIDs(ids ...TokenID) *TokenDissociateTransaction {
	tx._RequireNotFrozen()
	tx.tokens = make([]TokenID, len(ids))
	copy(tx.tokens, ids)

	return tx
}

// AddTokenID Adds the token to the list of tokens to be dissociated.
func (tx *TokenDissociateTransaction) AddTokenID(id TokenID) *TokenDissociateTransaction {
	tx._RequireNotFrozen()
	if tx.tokens == nil {
		tx.tokens = make([]TokenID, 0)
	}

	tx.tokens = append(tx.tokens, id)

	return tx
}

// GetTokenIDs returns the tokens to be associated with the provided account
func (tx *TokenDissociateTransaction) GetTokenIDs() []TokenID {
	tokenIDs := make([]TokenID, len(tx.tokens))
	copy(tokenIDs, tx.tokens)

	return tokenIDs
}

// ----------- Overridden functions ----------------

func (tx TokenDissociateTransaction) getName() string {
	return "TokenDissociateTransaction"
}

func (tx TokenDissociateTransaction) validateNetworkOnIDs(client *Client) error {
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

func (tx TokenDissociateTransaction) build() *services.TransactionBody {
	return &services.TransactionBody{
		TransactionFee:           tx.transactionFee,
		Memo:                     tx.Transaction.memo,
		TransactionValidDuration: _DurationToProtobuf(tx.GetTransactionValidDuration()),
		TransactionID:            tx.transactionID._ToProtobuf(),
		Data: &services.TransactionBody_TokenDissociate{
			TokenDissociate: tx.buildProtoBody(),
		},
	}
}

func (tx TokenDissociateTransaction) buildScheduled() (*services.SchedulableTransactionBody, error) {
	return &services.SchedulableTransactionBody{
		TransactionFee: tx.transactionFee,
		Memo:           tx.Transaction.memo,
		Data: &services.SchedulableTransactionBody_TokenDissociate{
			TokenDissociate: tx.buildProtoBody(),
		},
	}, nil
}

func (tx TokenDissociateTransaction) buildProtoBody() *services.TokenDissociateTransactionBody {
	body := &services.TokenDissociateTransactionBody{}
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

func (tx TokenDissociateTransaction) getMethod(channel *_Channel) _Method {
	return _Method{
		transaction: channel._GetToken().DissociateTokens,
	}
}

func (tx TokenDissociateTransaction) constructScheduleProtobuf() (*services.SchedulableTransactionBody, error) {
	return tx.buildScheduled()
}

func (tx TokenDissociateTransaction) getBaseTransaction() *Transaction[TransactionInterface] {
	return castFromConcreteToBaseTransaction(tx.Transaction, &tx)
}
