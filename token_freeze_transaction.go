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
	transaction
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
	tx := TokenFreezeTransaction{
		transaction: _NewTransaction(),
	}

	tx.e = &tx
	tx._SetDefaultMaxTransactionFee(NewHbar(30))

	return &tx
}

func _TokenFreezeTransactionFromProtobuf(tx transaction, pb *services.TransactionBody) *TokenFreezeTransaction {
	resultTx := &TokenFreezeTransaction{
		transaction: tx,
		tokenID:     _TokenIDFromProtobuf(pb.GetTokenFreeze().GetToken()),
		accountID:   _AccountIDFromProtobuf(pb.GetTokenFreeze().GetAccount()),
	}
	resultTx.e = resultTx
	return resultTx
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

// ---- Required Interfaces ---- //

// Sign uses the provided privateKey to sign the transaction.
func (tx *TokenFreezeTransaction) Sign(privateKey PrivateKey) *TokenFreezeTransaction {
	tx.transaction.Sign(privateKey)
	return tx
}

// SignWithOperator signs the transaction with client's operator privateKey.
func (tx *TokenFreezeTransaction) SignWithOperator(client *Client) (*TokenFreezeTransaction, error) {
	_, err := tx.transaction.SignWithOperator(client)
	return tx, err
}

// SignWith executes the TransactionSigner and adds the resulting signature data to the transaction's signature map
// with the publicKey as the map key.
func (tx *TokenFreezeTransaction) SignWith(
	publicKey PublicKey,
	signer TransactionSigner,
) *TokenFreezeTransaction {
	tx.transaction.SignWith(publicKey, signer)
	return tx
}

// AddSignature adds a signature to the transaction.
func (tx *TokenFreezeTransaction) AddSignature(publicKey PublicKey, signature []byte) *TokenFreezeTransaction {
	tx.transaction.AddSignature(publicKey, signature)
	return tx
}

// When execution is attempted, a single attempt will timeout when this deadline is reached. (The SDK may subsequently retry the execution.)
func (tx *TokenFreezeTransaction) SetGrpcDeadline(deadline *time.Duration) *TokenFreezeTransaction {
	tx.transaction.SetGrpcDeadline(deadline)
	return tx
}

func (tx *TokenFreezeTransaction) Freeze() (*TokenFreezeTransaction, error) {
	return tx.FreezeWith(nil)
}

func (tx *TokenFreezeTransaction) FreezeWith(client *Client) (*TokenFreezeTransaction, error) {
	_, err := tx.transaction.FreezeWith(client)
	return tx, err
}

// SetMaxTransactionFee sets the max transaction fee for this TokenFreezeTransaction.
func (tx *TokenFreezeTransaction) SetMaxTransactionFee(fee Hbar) *TokenFreezeTransaction {
	tx.transaction.SetMaxTransactionFee(fee)
	return tx
}

// SetRegenerateTransactionID sets if transaction IDs should be regenerated when `TRANSACTION_EXPIRED` is received
func (tx *TokenFreezeTransaction) SetRegenerateTransactionID(regenerateTransactionID bool) *TokenFreezeTransaction {
	tx.transaction.SetRegenerateTransactionID(regenerateTransactionID)
	return tx
}

// SetTransactionMemo sets the memo for this TokenFreezeTransaction.
func (tx *TokenFreezeTransaction) SetTransactionMemo(memo string) *TokenFreezeTransaction {
	tx.transaction.SetTransactionMemo(memo)
	return tx
}

// SetTransactionValidDuration sets the valid duration for this TokenFreezeTransaction.
func (tx *TokenFreezeTransaction) SetTransactionValidDuration(duration time.Duration) *TokenFreezeTransaction {
	tx.transaction.SetTransactionValidDuration(duration)
	return tx
}

// SetTransactionID sets the TransactionID for this TokenFreezeTransaction.
func (tx *TokenFreezeTransaction) SetTransactionID(transactionID TransactionID) *TokenFreezeTransaction {
	tx.transaction.SetTransactionID(transactionID)
	return tx
}

// SetNodeAccountIDs sets the _Node AccountID for this TokenFreezeTransaction.
func (tx *TokenFreezeTransaction) SetNodeAccountIDs(nodeID []AccountID) *TokenFreezeTransaction {
	tx.transaction.SetNodeAccountIDs(nodeID)
	return tx
}

// SetMaxRetry sets the max number of errors before execution will fail.
func (tx *TokenFreezeTransaction) SetMaxRetry(count int) *TokenFreezeTransaction {
	tx.transaction.SetMaxRetry(count)
	return tx
}

// SetMaxBackoff The maximum amount of time to wait between retries.
// Every retry attempt will increase the wait time exponentially until it reaches this time.
func (tx *TokenFreezeTransaction) SetMaxBackoff(max time.Duration) *TokenFreezeTransaction {
	tx.transaction.SetMaxBackoff(max)
	return tx
}

// SetMinBackoff sets the minimum amount of time to wait between retries.
func (tx *TokenFreezeTransaction) SetMinBackoff(min time.Duration) *TokenFreezeTransaction {
	tx.transaction.SetMinBackoff(min)
	return tx
}

func (tx *TokenFreezeTransaction) SetLogLevel(level LogLevel) *TokenFreezeTransaction {
	tx.transaction.SetLogLevel(level)
	return tx
}

// ----------- overriden functions ----------------

func (tx *TokenFreezeTransaction) getName() string {
	return "TokenFreezeTransaction"
}

func (tx *TokenFreezeTransaction) validateNetworkOnIDs(client *Client) error {
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

func (tx *TokenFreezeTransaction) build() *services.TransactionBody {
	return &services.TransactionBody{
		TransactionFee:           tx.transactionFee,
		Memo:                     tx.transaction.memo,
		TransactionValidDuration: _DurationToProtobuf(tx.GetTransactionValidDuration()),
		TransactionID:            tx.transactionID._ToProtobuf(),
		Data: &services.TransactionBody_TokenFreeze{
			TokenFreeze: tx.buildProtoBody(),
		},
	}
}

func (tx *TokenFreezeTransaction) buildScheduled() (*services.SchedulableTransactionBody, error) {
	return &services.SchedulableTransactionBody{
		TransactionFee: tx.transactionFee,
		Memo:           tx.transaction.memo,
		Data: &services.SchedulableTransactionBody_TokenFreeze{
			TokenFreeze: tx.buildProtoBody(),
		},
	}, nil
}

func (tx *TokenFreezeTransaction) buildProtoBody() *services.TokenFreezeAccountTransactionBody {
	body := &services.TokenFreezeAccountTransactionBody{}
	if tx.tokenID != nil {
		body.Token = tx.tokenID._ToProtobuf()
	}

	if tx.accountID != nil {
		body.Account = tx.accountID._ToProtobuf()
	}

	return body
}

func (tx *TokenFreezeTransaction) getMethod(channel *_Channel) _Method {
	return _Method{
		transaction: channel._GetToken().FreezeTokenAccount,
	}
}
