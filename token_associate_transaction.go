package hedera

/*-
 *
 * Hedera Go SDK
 *
 * Copyright (C) 2020 - 2023 Hedera Hashgraph, LLC
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

import (
	"time"

	"github.com/hashgraph/hedera-protobufs-go/services"
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
	transaction
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
	tx := TokenAssociateTransaction{
		transaction: _NewTransaction(),
	}

	tx.e = &tx
	tx._SetDefaultMaxTransactionFee(NewHbar(5))

	return &tx
}

func _TokenAssociateTransactionFromProtobuf(tx transaction, pb *services.TransactionBody) *TokenAssociateTransaction {
	tokens := make([]TokenID, 0)
	for _, token := range pb.GetTokenAssociate().Tokens {
		if tokenID := _TokenIDFromProtobuf(token); tokenID != nil {
			tokens = append(tokens, *tokenID)
		}
	}

	resultTx := &TokenAssociateTransaction{
		transaction: tx,
		accountID:   _AccountIDFromProtobuf(pb.GetTokenAssociate().GetAccount()),
		tokens:      tokens,
	}
	resultTx.e = resultTx
	return resultTx
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

// ---- Required Interfaces ---- //

// Sign uses the provided privateKey to sign the transaction.
func (tx *TokenAssociateTransaction) Sign(privateKey PrivateKey) *TokenAssociateTransaction {
	tx.transaction.Sign(privateKey)
	return tx
}

// SignWithOperator signs the transaction with client's operator privateKey.
func (tx *TokenAssociateTransaction) SignWithOperator(client *Client) (*TokenAssociateTransaction, error) {
	_, err := tx.transaction.SignWithOperator(client)
	return tx, err
}

// SignWith executes the TransactionSigner and adds the resulting signature data to the transaction's signature map
// with the publicKey as the map key.
func (tx *TokenAssociateTransaction) SignWith(
	publicKey PublicKey,
	signer TransactionSigner,
) *TokenAssociateTransaction {
	tx.transaction.SignWith(publicKey, signer)
	return tx
}

// AddSignature adds a signature to the transaction.
func (tx *TokenAssociateTransaction) AddSignature(publicKey PublicKey, signature []byte) *TokenAssociateTransaction {
	tx.transaction.AddSignature(publicKey, signature)
	return tx
}

// When execution is attempted, a single attempt will timeout when this deadline is reached. (The SDK may subsequently retry the execution.)
func (tx *TokenAssociateTransaction) SetGrpcDeadline(deadline *time.Duration) *TokenAssociateTransaction {
	tx.transaction.SetGrpcDeadline(deadline)
	return tx
}

func (tx *TokenAssociateTransaction) Freeze() (*TokenAssociateTransaction, error) {
	return tx.FreezeWith(nil)
}

func (tx *TokenAssociateTransaction) FreezeWith(client *Client) (*TokenAssociateTransaction, error) {
	_, err := tx.transaction.FreezeWith(client)
	return tx, err
}

// SetMaxTransactionFee sets the max transaction fee for this TokenAssociateTransaction.
func (tx *TokenAssociateTransaction) SetMaxTransactionFee(fee Hbar) *TokenAssociateTransaction {
	tx.transaction.SetMaxTransactionFee(fee)
	return tx
}

// SetRegenerateTransactionID sets if transaction IDs should be regenerated when `TRANSACTION_EXPIRED` is received
func (tx *TokenAssociateTransaction) SetRegenerateTransactionID(regenerateTransactionID bool) *TokenAssociateTransaction {
	tx.transaction.SetRegenerateTransactionID(regenerateTransactionID)
	return tx
}

// SetTransactionMemo sets the memo for this TokenAssociateTransaction.
func (tx *TokenAssociateTransaction) SetTransactionMemo(memo string) *TokenAssociateTransaction {
	tx.transaction.SetTransactionMemo(memo)
	return tx
}

// SetTransactionValidDuration sets the valid duration for this TokenAssociateTransaction.
func (tx *TokenAssociateTransaction) SetTransactionValidDuration(duration time.Duration) *TokenAssociateTransaction {
	tx.transaction.SetTransactionValidDuration(duration)
	return tx
}

// SetTransactionID sets the TransactionID for this TokenAssociateTransaction.
func (tx *TokenAssociateTransaction) SetTransactionID(transactionID TransactionID) *TokenAssociateTransaction {
	tx.transaction.SetTransactionID(transactionID)
	return tx
}

// SetNodeAccountIDs sets the _Node AccountID for this TokenAssociateTransaction.
func (tx *TokenAssociateTransaction) SetNodeAccountIDs(nodeID []AccountID) *TokenAssociateTransaction {
	tx.transaction.SetNodeAccountIDs(nodeID)
	return tx
}

// SetMaxRetry sets the max number of errors before execution will fail.
func (tx *TokenAssociateTransaction) SetMaxRetry(count int) *TokenAssociateTransaction {
	tx.transaction.SetMaxRetry(count)
	return tx
}

// SetMaxBackoff The maximum amount of time to wait between retries.
// Every retry attempt will increase the wait time exponentially until it reaches this time.
func (tx *TokenAssociateTransaction) SetMaxBackoff(max time.Duration) *TokenAssociateTransaction {
	tx.transaction.SetMaxBackoff(max)
	return tx
}

// SetMinBackoff sets the minimum amount of time to wait between retries.
func (tx *TokenAssociateTransaction) SetMinBackoff(min time.Duration) *TokenAssociateTransaction {
	tx.transaction.SetMinBackoff(min)
	return tx
}

func (tx *TokenAssociateTransaction) SetLogLevel(level LogLevel) *TokenAssociateTransaction {
	tx.transaction.SetLogLevel(level)
	return tx
}

// ----------- overriden functions ----------------

func (tx *TokenAssociateTransaction) getName() string {
	return "TokenAssociateTransaction"
}

func (tx *TokenAssociateTransaction) validateNetworkOnIDs(client *Client) error {
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

func (tx *TokenAssociateTransaction) build() *services.TransactionBody {
	body := tx.buildProtoBody()

	return &services.TransactionBody{
		TransactionFee:           tx.transactionFee,
		Memo:                     tx.transaction.memo,
		TransactionValidDuration: _DurationToProtobuf(tx.GetTransactionValidDuration()),
		TransactionID:            tx.transactionID._ToProtobuf(),
		Data: &services.TransactionBody_TokenAssociate{
			TokenAssociate: body,
		},
	}
}

func (tx *TokenAssociateTransaction) buildScheduled() (*services.SchedulableTransactionBody, error) {
	return &services.SchedulableTransactionBody{
		TransactionFee: tx.transactionFee,
		Memo:           tx.transaction.memo,
		Data: &services.SchedulableTransactionBody_TokenAssociate{
			TokenAssociate: tx.buildProtoBody(),
		},
	}, nil
}

func (tx *TokenAssociateTransaction) buildProtoBody() *services.TokenAssociateTransactionBody {
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

func (tx *TokenAssociateTransaction) getMethod(channel *_Channel) _Method {
	return _Method{
		transaction: channel._GetToken().AssociateTokens,
	}
}
func (tx *TokenAssociateTransaction) _ConstructScheduleProtobuf() (*services.SchedulableTransactionBody, error) {
	return tx.buildScheduled()
}
