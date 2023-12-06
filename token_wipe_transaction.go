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
	Transaction
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
	tx := TokenWipeTransaction{
		Transaction: _NewTransaction(),
	}

	tx.e = &tx
	tx._SetDefaultMaxTransactionFee(NewHbar(30))

	return &tx
}

func _TokenWipeTransactionFromProtobuf(tx Transaction, pb *services.TransactionBody) *TokenWipeTransaction {
	resultTx := &TokenWipeTransaction{
		Transaction: tx,
		tokenID:     _TokenIDFromProtobuf(pb.GetTokenWipe().GetToken()),
		accountID:   _AccountIDFromProtobuf(pb.GetTokenWipe().GetAccount()),
		amount:      pb.GetTokenWipe().Amount,
		serial:      pb.GetTokenWipe().GetSerialNumbers(),
	}
	resultTx.e = resultTx
	return resultTx
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

// ---- Required Interfaces ---- //

// When execution is attempted, a single attempt will timeout when this deadline is reached. (The SDK may subsequently retry the execution.)
func (tx *TokenWipeTransaction) SetGrpcDeadline(deadline *time.Duration) *TokenWipeTransaction {
	tx.Transaction.SetGrpcDeadline(deadline)
	return tx
}

// Sign uses the provided privateKey to sign the transaction.
func (tx *TokenWipeTransaction) Sign(privateKey PrivateKey) *TokenWipeTransaction {
	tx.Transaction.Sign(privateKey)
	return tx
}

// SignWithOperator signs the transaction with client's operator privateKey.
func (tx *TokenWipeTransaction) SignWithOperator(client *Client) (*TokenWipeTransaction, error) {
	_, err := tx.Transaction.SignWithOperator(client)
	if err != nil {
		return nil, err
	}
	return tx, nil
}

// SignWith executes the TransactionSigner and adds the resulting signature data to the transaction's signature map
// with the publicKey as the map key.
func (tx *TokenWipeTransaction) SignWith(
	publicKey PublicKey,
	signer TransactionSigner,
) *TokenWipeTransaction {
	tx.Transaction.SignWith(publicKey, signer)
	return tx
}

// AddSignature adds a signature to the transaction.
func (tx *TokenWipeTransaction) AddSignature(publicKey PublicKey, signature []byte) *TokenWipeTransaction {
	tx.Transaction.AddSignature(publicKey, signature)
	return tx
}

func (tx *TokenWipeTransaction) Freeze() (*TokenWipeTransaction, error) {
	return tx.FreezeWith(nil)
}

func (tx *TokenWipeTransaction) FreezeWith(client *Client) (*TokenWipeTransaction, error) {
	_, err := tx.Transaction.FreezeWith(client)
	return tx, err
}

// SetMaxTransactionFee sets the max transaction fee for this TokenWipeTransaction.
func (tx *TokenWipeTransaction) SetMaxTransactionFee(fee Hbar) *TokenWipeTransaction {
	tx.Transaction.SetMaxTransactionFee(fee)
	return tx
}

// SetRegenerateTransactionID sets if transaction IDs should be regenerated when `TRANSACTION_EXPIRED` is received
func (tx *TokenWipeTransaction) SetRegenerateTransactionID(regenerateTransactionID bool) *TokenWipeTransaction {
	tx.Transaction.SetRegenerateTransactionID(regenerateTransactionID)
	return tx
}

// SetTransactionMemo sets the memo for this TokenWipeTransaction.
func (tx *TokenWipeTransaction) SetTransactionMemo(memo string) *TokenWipeTransaction {
	tx.Transaction.SetTransactionMemo(memo)
	return tx
}

// SetTransactionValidDuration sets the valid duration for this TokenWipeTransaction.
func (tx *TokenWipeTransaction) SetTransactionValidDuration(duration time.Duration) *TokenWipeTransaction {
	tx.Transaction.SetTransactionValidDuration(duration)
	return tx
}

// SetTransactionID sets the TransactionID for this TokenWipeTransaction.
func (tx *TokenWipeTransaction) SetTransactionID(transactionID TransactionID) *TokenWipeTransaction {
	tx.Transaction.SetTransactionID(transactionID)
	return tx
}

// SetNodeAccountIDs sets the _Node AccountID for this TokenWipeTransaction.
func (tx *TokenWipeTransaction) SetNodeAccountIDs(nodeID []AccountID) *TokenWipeTransaction {
	tx.Transaction.SetNodeAccountIDs(nodeID)
	return tx
}

// SetMaxRetry sets the max number of errors before execution will fail.
func (tx *TokenWipeTransaction) SetMaxRetry(count int) *TokenWipeTransaction {
	tx.Transaction.SetMaxRetry(count)
	return tx
}

// SetMaxBackoff The maximum amount of time to wait between retries.
// Every retry attempt will increase the wait time exponentially until it reaches this time.
func (tx *TokenWipeTransaction) SetMaxBackoff(max time.Duration) *TokenWipeTransaction {
	tx.Transaction.SetMaxBackoff(max)
	return tx
}

// SetMinBackoff sets the minimum amount of time to wait between retries.
func (tx *TokenWipeTransaction) SetMinBackoff(min time.Duration) *TokenWipeTransaction {
	tx.Transaction.SetMinBackoff(min)
	return tx
}

func (tx *TokenWipeTransaction) SetLogLevel(level LogLevel) *TokenWipeTransaction {
	tx.Transaction.SetLogLevel(level)
	return tx
}

// ----------- overriden functions ----------------

func (tx *TokenWipeTransaction) getName() string {
	return "TokenWipeTransaction"
}

func (tx *TokenWipeTransaction) validateNetworkOnIDs(client *Client) error {
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

func (tx *TokenWipeTransaction) build() *services.TransactionBody {
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

func (tx *TokenWipeTransaction) buildScheduled() (*services.SchedulableTransactionBody, error) {
	return &services.SchedulableTransactionBody{
		TransactionFee: tx.transactionFee,
		Memo:           tx.Transaction.memo,
		Data: &services.SchedulableTransactionBody_TokenWipe{
			TokenWipe: tx.buildProtoBody(),
		},
	}, nil
}

func (tx *TokenWipeTransaction) buildProtoBody() *services.TokenWipeAccountTransactionBody {
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

func (tx *TokenWipeTransaction) getMethod(channel *_Channel) _Method {
	return _Method{
		transaction: channel._GetToken().WipeTokenAccount,
	}
}
func (tx *TokenWipeTransaction) _ConstructScheduleProtobuf() (*services.SchedulableTransactionBody, error) {
	return tx.buildScheduled()
}
