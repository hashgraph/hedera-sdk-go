package hedera

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

import (
	"time"

	"github.com/hashgraph/hedera-protobufs-go/services"
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
	Transaction
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
	tx := TokenDissociateTransaction{
		Transaction: _NewTransaction(),
	}

	tx._SetDefaultMaxTransactionFee(NewHbar(5))

	return &tx
}

func _TokenDissociateTransactionFromProtobuf(tx Transaction, pb *services.TransactionBody) *TokenDissociateTransaction {
	tokens := make([]TokenID, 0)
	for _, token := range pb.GetTokenDissociate().Tokens {
		if tokenID := _TokenIDFromProtobuf(token); tokenID != nil {
			tokens = append(tokens, *tokenID)
		}
	}

	return &TokenDissociateTransaction{
		Transaction: tx,
		accountID:   _AccountIDFromProtobuf(pb.GetTokenDissociate().GetAccount()),
		tokens:      tokens,
	}
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

// ---- Required Interfaces ---- //

// Sign uses the provided privateKey to sign the transaction.
func (tx *TokenDissociateTransaction) Sign(privateKey PrivateKey) *TokenDissociateTransaction {
	tx.Transaction.Sign(privateKey)
	return tx
}

// SignWithOperator signs the transaction with client's operator privateKey.
func (tx *TokenDissociateTransaction) SignWithOperator(client *Client) (*TokenDissociateTransaction, error) {
	_, err := tx.Transaction.signWithOperator(client, tx)
	if err != nil {
		return nil, err
	}
	return tx, nil
}

// SignWith executes the TransactionSigner and adds the resulting signature data to the Transaction's signature map
// with the publicKey as the map key.
func (tx *TokenDissociateTransaction) SignWith(
	publicKey PublicKey,
	signer TransactionSigner,
) *TokenDissociateTransaction {
	tx.Transaction.SignWith(publicKey, signer)
	return tx
}

// AddSignature adds a signature to the transaction.
func (tx *TokenDissociateTransaction) AddSignature(publicKey PublicKey, signature []byte) *TokenDissociateTransaction {
	tx.Transaction.AddSignature(publicKey, signature)
	return tx
}

// When execution is attempted, a single attempt will timeout when this deadline is reached. (The SDK may subsequently retry the execution.)
func (tx *TokenDissociateTransaction) SetGrpcDeadline(deadline *time.Duration) *TokenDissociateTransaction {
	tx.Transaction.SetGrpcDeadline(deadline)
	return tx
}

func (tx *TokenDissociateTransaction) Freeze() (*TokenDissociateTransaction, error) {
	return tx.FreezeWith(nil)
}

func (tx *TokenDissociateTransaction) FreezeWith(client *Client) (*TokenDissociateTransaction, error) {
	_, err := tx.Transaction.freezeWith(client, tx)
	return tx, err
}

// SetMaxTransactionFee sets the max transaction fee for this TokenDissociateTransaction.
func (tx *TokenDissociateTransaction) SetMaxTransactionFee(fee Hbar) *TokenDissociateTransaction {
	tx.Transaction.SetMaxTransactionFee(fee)
	return tx
}

// SetRegenerateTransactionID sets if transaction IDs should be regenerated when `TRANSACTION_EXPIRED` is received
func (tx *TokenDissociateTransaction) SetRegenerateTransactionID(regenerateTransactionID bool) *TokenDissociateTransaction {
	tx.Transaction.SetRegenerateTransactionID(regenerateTransactionID)
	return tx
}

// SetTransactionMemo sets the memo for this TokenDissociateTransaction.
func (tx *TokenDissociateTransaction) SetTransactionMemo(memo string) *TokenDissociateTransaction {
	tx.Transaction.SetTransactionMemo(memo)
	return tx
}

// SetTransactionValidDuration sets the valid duration for this TokenDissociateTransaction.
func (tx *TokenDissociateTransaction) SetTransactionValidDuration(duration time.Duration) *TokenDissociateTransaction {
	tx.Transaction.SetTransactionValidDuration(duration)
	return tx
}

// SetTransactionID sets the TransactionID for this TokenDissociateTransaction.
func (tx *TokenDissociateTransaction) SetTransactionID(transactionID TransactionID) *TokenDissociateTransaction {
	tx.Transaction.SetTransactionID(transactionID)
	return tx
}

// SetNodeAccountIDs sets the _Node AccountID for this TokenDissociateTransaction.
func (tx *TokenDissociateTransaction) SetNodeAccountIDs(nodeID []AccountID) *TokenDissociateTransaction {
	tx.Transaction.SetNodeAccountIDs(nodeID)
	return tx
}

// SetMaxRetry sets the max number of errors before execution will fail.
func (tx *TokenDissociateTransaction) SetMaxRetry(count int) *TokenDissociateTransaction {
	tx.Transaction.SetMaxRetry(count)
	return tx
}

// SetMaxBackoff The maximum amount of time to wait between retries.
// Every retry attempt will increase the wait time exponentially until it reaches this time.
func (tx *TokenDissociateTransaction) SetMaxBackoff(max time.Duration) *TokenDissociateTransaction {
	tx.Transaction.SetMaxBackoff(max)
	return tx
}

// SetMinBackoff sets the minimum amount of time to wait between retries.
func (tx *TokenDissociateTransaction) SetMinBackoff(min time.Duration) *TokenDissociateTransaction {
	tx.Transaction.SetMinBackoff(min)
	return tx
}

func (tx *TokenDissociateTransaction) SetLogLevel(level LogLevel) *TokenDissociateTransaction {
	tx.Transaction.SetLogLevel(level)
	return tx
}

func (tx *TokenDissociateTransaction) Execute(client *Client) (TransactionResponse, error) {
	return tx.Transaction.execute(client, tx)
}

func (tx *TokenDissociateTransaction) Schedule() (*ScheduleCreateTransaction, error) {
	return tx.Transaction.schedule(tx)
}

// ----------- Overridden functions ----------------

func (tx *TokenDissociateTransaction) getName() string {
	return "TokenDissociateTransaction"
}

func (tx *TokenDissociateTransaction) validateNetworkOnIDs(client *Client) error {
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

func (tx *TokenDissociateTransaction) build() *services.TransactionBody {
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

func (tx *TokenDissociateTransaction) buildScheduled() (*services.SchedulableTransactionBody, error) {
	return &services.SchedulableTransactionBody{
		TransactionFee: tx.transactionFee,
		Memo:           tx.Transaction.memo,
		Data: &services.SchedulableTransactionBody_TokenDissociate{
			TokenDissociate: tx.buildProtoBody(),
		},
	}, nil
}

func (tx *TokenDissociateTransaction) buildProtoBody() *services.TokenDissociateTransactionBody {
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

func (tx *TokenDissociateTransaction) getMethod(channel *_Channel) _Method {
	return _Method{
		transaction: channel._GetToken().DissociateTokens,
	}
}

func (tx *TokenDissociateTransaction) _ConstructScheduleProtobuf() (*services.SchedulableTransactionBody, error) {
	return tx.buildScheduled()
}
