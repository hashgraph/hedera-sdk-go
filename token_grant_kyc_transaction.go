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

// TokenGrantKycTransaction
// Grants KYC to the account for the given token. Must be signed by the Token's kycKey.
// If the provided account is not found, the transaction will resolve to INVALID_ACCOUNT_ID.
// If the provided account has been deleted, the transaction will resolve to ACCOUNT_DELETED.
// If the provided token is not found, the transaction will resolve to INVALID_TOKEN_ID.
// If the provided token has been deleted, the transaction will resolve to TOKEN_WAS_DELETED.
// If an Association between the provided token and account is not found, the transaction will
// resolve to TOKEN_NOT_ASSOCIATED_TO_ACCOUNT.
// If no KYC Key is defined, the transaction will resolve to TOKEN_HAS_NO_KYC_KEY.
// Once executed the Account is marked as KYC Granted.
type TokenGrantKycTransaction struct {
	Transaction
	tokenID   *TokenID
	accountID *AccountID
}

// NewTokenGrantKycTransaction creates TokenGrantKycTransaction which
// grants KYC to the account for the given token. Must be signed by the Token's kycKey.
// If the provided account is not found, the transaction will resolve to INVALID_ACCOUNT_ID.
// If the provided account has been deleted, the transaction will resolve to ACCOUNT_DELETED.
// If the provided token is not found, the transaction will resolve to INVALID_TOKEN_ID.
// If the provided token has been deleted, the transaction will resolve to TOKEN_WAS_DELETED.
// If an Association between the provided token and account is not found, the transaction will
// resolve to TOKEN_NOT_ASSOCIATED_TO_ACCOUNT.
// If no KYC Key is defined, the transaction will resolve to TOKEN_HAS_NO_KYC_KEY.
// Once executed the Account is marked as KYC Granted.
func NewTokenGrantKycTransaction() *TokenGrantKycTransaction {
	tx := TokenGrantKycTransaction{
		Transaction: _NewTransaction(),
	}

	tx._SetDefaultMaxTransactionFee(NewHbar(30))

	return &tx
}

func _TokenGrantKycTransactionFromProtobuf(tx Transaction, pb *services.TransactionBody) *TokenGrantKycTransaction {
	return &TokenGrantKycTransaction{
		Transaction: tx,
		tokenID:     _TokenIDFromProtobuf(pb.GetTokenGrantKyc().GetToken()),
		accountID:   _AccountIDFromProtobuf(pb.GetTokenGrantKyc().GetAccount()),
	}
}

// SetTokenID Sets the token for which this account will be granted KYC.
// If token does not exist, transaction results in INVALID_TOKEN_ID
func (tx *TokenGrantKycTransaction) SetTokenID(tokenID TokenID) *TokenGrantKycTransaction {
	tx._RequireNotFrozen()
	tx.tokenID = &tokenID
	return tx
}

// GetTokenID returns the token for which this account will be granted KYC.
func (tx *TokenGrantKycTransaction) GetTokenID() TokenID {
	if tx.tokenID == nil {
		return TokenID{}
	}

	return *tx.tokenID
}

// SetAccountID Sets the account to be KYCed
func (tx *TokenGrantKycTransaction) SetAccountID(accountID AccountID) *TokenGrantKycTransaction {
	tx._RequireNotFrozen()
	tx.accountID = &accountID
	return tx
}

// GetAccountID returns the AccountID that is being KYCed
func (tx *TokenGrantKycTransaction) GetAccountID() AccountID {
	if tx.accountID == nil {
		return AccountID{}
	}

	return *tx.accountID
}

// ---- Required Interfaces ---- //

// Sign uses the provided privateKey to sign the transaction.
func (tx *TokenGrantKycTransaction) Sign(privateKey PrivateKey) *TokenGrantKycTransaction {
	tx.Transaction.Sign(privateKey)
	return tx
}

// SignWithOperator signs the transaction with client's operator privateKey.
func (tx *TokenGrantKycTransaction) SignWithOperator(client *Client) (*TokenGrantKycTransaction, error) {
	_, err := tx.Transaction.signWithOperator(client, tx)
	if err != nil {
		return nil, err
	}
	return tx, nil
}

// SignWith executes the TransactionSigner and adds the resulting signature data to the Transaction's signature map
// with the publicKey as the map key.
func (tx *TokenGrantKycTransaction) SignWith(
	publicKey PublicKey,
	signer TransactionSigner,
) *TokenGrantKycTransaction {
	tx.Transaction.SignWith(publicKey, signer)
	return tx
}

// AddSignature adds a signature to the transaction.
func (tx *TokenGrantKycTransaction) AddSignature(publicKey PublicKey, signature []byte) *TokenGrantKycTransaction {
	tx.Transaction.AddSignature(publicKey, signature)
	return tx
}

// When execution is attempted, a single attempt will timeout when this deadline is reached. (The SDK may subsequently retry the execution.)
func (tx *TokenGrantKycTransaction) SetGrpcDeadline(deadline *time.Duration) *TokenGrantKycTransaction {
	tx.Transaction.SetGrpcDeadline(deadline)
	return tx
}

func (tx *TokenGrantKycTransaction) Freeze() (*TokenGrantKycTransaction, error) {
	return tx.FreezeWith(nil)
}

func (tx *TokenGrantKycTransaction) FreezeWith(client *Client) (*TokenGrantKycTransaction, error) {
	_, err := tx.Transaction.freezeWith(client, tx)
	return tx, err
}

// SetMaxTransactionFee sets the max transaction fee for this TokenGrantKycTransaction.
func (tx *TokenGrantKycTransaction) SetMaxTransactionFee(fee Hbar) *TokenGrantKycTransaction {
	tx.Transaction.SetMaxTransactionFee(fee)
	return tx
}

// SetRegenerateTransactionID sets if transaction IDs should be regenerated when `TRANSACTION_EXPIRED` is received
func (tx *TokenGrantKycTransaction) SetRegenerateTransactionID(regenerateTransactionID bool) *TokenGrantKycTransaction {
	tx.Transaction.SetRegenerateTransactionID(regenerateTransactionID)
	return tx
}

// SetTransactionMemo sets the memo for this TokenGrantKycTransaction.
func (tx *TokenGrantKycTransaction) SetTransactionMemo(memo string) *TokenGrantKycTransaction {
	tx.Transaction.SetTransactionMemo(memo)
	return tx
}

// SetTransactionValidDuration sets the valid duration for this TokenGrantKycTransaction.
func (tx *TokenGrantKycTransaction) SetTransactionValidDuration(duration time.Duration) *TokenGrantKycTransaction {
	tx.Transaction.SetTransactionValidDuration(duration)
	return tx
}

// SetTransactionID sets the TransactionID for this TokenGrantKycTransaction.
func (tx *TokenGrantKycTransaction) SetTransactionID(transactionID TransactionID) *TokenGrantKycTransaction {
	tx.Transaction.SetTransactionID(transactionID)
	return tx
}

// SetNodeAccountIDs sets the _Node AccountID for this TokenGrantKycTransaction.
func (tx *TokenGrantKycTransaction) SetNodeAccountIDs(nodeID []AccountID) *TokenGrantKycTransaction {
	tx.Transaction.SetNodeAccountIDs(nodeID)
	return tx
}

// SetMaxRetry sets the max number of errors before execution will fail.
func (tx *TokenGrantKycTransaction) SetMaxRetry(count int) *TokenGrantKycTransaction {
	tx.Transaction.SetMaxRetry(count)
	return tx
}

// SetMaxBackoff The maximum amount of time to wait between retries.
// Every retry attempt will increase the wait time exponentially until it reaches this time.
func (tx *TokenGrantKycTransaction) SetMaxBackoff(max time.Duration) *TokenGrantKycTransaction {
	tx.Transaction.SetMaxBackoff(max)
	return tx
}

// SetMinBackoff sets the minimum amount of time to wait between retries.
func (tx *TokenGrantKycTransaction) SetMinBackoff(min time.Duration) *TokenGrantKycTransaction {
	tx.Transaction.SetMinBackoff(min)
	return tx
}

func (tx *TokenGrantKycTransaction) SetLogLevel(level LogLevel) *TokenGrantKycTransaction {
	tx.Transaction.SetLogLevel(level)
	return tx
}

func (tx *TokenGrantKycTransaction) Execute(client *Client) (TransactionResponse, error) {
	return tx.Transaction.execute(client, tx)
}

func (tx *TokenGrantKycTransaction) Schedule() (*ScheduleCreateTransaction, error) {
	return tx.Transaction.schedule(tx)
}

// ----------- Overridden functions ----------------

func (tx *TokenGrantKycTransaction) getName() string {
	return "TokenGrantKycTransaction"
}

func (tx *TokenGrantKycTransaction) validateNetworkOnIDs(client *Client) error {
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

func (tx *TokenGrantKycTransaction) build() *services.TransactionBody {
	return &services.TransactionBody{
		TransactionFee:           tx.transactionFee,
		Memo:                     tx.Transaction.memo,
		TransactionValidDuration: _DurationToProtobuf(tx.GetTransactionValidDuration()),
		TransactionID:            tx.transactionID._ToProtobuf(),
		Data: &services.TransactionBody_TokenGrantKyc{
			TokenGrantKyc: tx.buildProtoBody(),
		},
	}
}

func (tx *TokenGrantKycTransaction) buildScheduled() (*services.SchedulableTransactionBody, error) {
	return &services.SchedulableTransactionBody{
		TransactionFee: tx.transactionFee,
		Memo:           tx.Transaction.memo,
		Data: &services.SchedulableTransactionBody_TokenGrantKyc{
			TokenGrantKyc: tx.buildProtoBody(),
		},
	}, nil
}

func (tx *TokenGrantKycTransaction) buildProtoBody() *services.TokenGrantKycTransactionBody {
	body := &services.TokenGrantKycTransactionBody{}
	if tx.tokenID != nil {
		body.Token = tx.tokenID._ToProtobuf()
	}

	if tx.accountID != nil {
		body.Account = tx.accountID._ToProtobuf()
	}

	return body
}

func (tx *TokenGrantKycTransaction) getMethod(channel *_Channel) _Method {
	return _Method{
		transaction: channel._GetToken().GrantKycToTokenAccount,
	}
}
func (tx *TokenGrantKycTransaction) _ConstructScheduleProtobuf() (*services.SchedulableTransactionBody, error) {
	return tx.buildScheduled()
}
