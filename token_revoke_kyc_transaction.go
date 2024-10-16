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

	"github.com/hashgraph/hedera-sdk-go/v2/proto/services"
)

// TokenRevokeKycTransaction
// Revokes KYC to the account for the given token. Must be signed by the Token's kycKey.
// If the provided account is not found, the transaction will resolve to INVALID_ACCOUNT_ID.
// If the provided account has been deleted, the transaction will resolve to ACCOUNT_DELETED.
// If the provided token is not found, the transaction will resolve to INVALID_TOKEN_ID.
// If the provided token has been deleted, the transaction will resolve to TOKEN_WAS_DELETED.
// If an Association between the provided token and account is not found, the transaction will
// resolve to TOKEN_NOT_ASSOCIATED_TO_ACCOUNT.
// If no KYC Key is defined, the transaction will resolve to TOKEN_HAS_NO_KYC_KEY.
// Once executed the Account is marked as KYC Revoked
type TokenRevokeKycTransaction struct {
	Transaction
	tokenID   *TokenID
	accountID *AccountID
}

// NewTokenRevokeKycTransaction creates TokenRevokeKycTransaction which
// revokes KYC to the account for the given token. Must be signed by the Token's kycKey.
// If the provided account is not found, the transaction will resolve to INVALID_ACCOUNT_ID.
// If the provided account has been deleted, the transaction will resolve to ACCOUNT_DELETED.
// If the provided token is not found, the transaction will resolve to INVALID_TOKEN_ID.
// If the provided token has been deleted, the transaction will resolve to TOKEN_WAS_DELETED.
// If an Association between the provided token and account is not found, the transaction will
// resolve to TOKEN_NOT_ASSOCIATED_TO_ACCOUNT.
// If no KYC Key is defined, the transaction will resolve to TOKEN_HAS_NO_KYC_KEY.
// Once executed the Account is marked as KYC Revoked
func NewTokenRevokeKycTransaction() *TokenRevokeKycTransaction {
	tx := TokenRevokeKycTransaction{
		Transaction: _NewTransaction(),
	}

	tx._SetDefaultMaxTransactionFee(NewHbar(30))

	return &tx
}

func _TokenRevokeKycTransactionFromProtobuf(transaction Transaction, pb *services.TransactionBody) *TokenRevokeKycTransaction {
	return &TokenRevokeKycTransaction{
		Transaction: transaction,
		tokenID:     _TokenIDFromProtobuf(pb.GetTokenRevokeKyc().GetToken()),
		accountID:   _AccountIDFromProtobuf(pb.GetTokenRevokeKyc().GetAccount()),
	}
}

// SetTokenID Sets the token for which this account will get his KYC revoked.
// If token does not exist, transaction results in INVALID_TOKEN_ID
func (tx *TokenRevokeKycTransaction) SetTokenID(tokenID TokenID) *TokenRevokeKycTransaction {
	tx._RequireNotFrozen()
	tx.tokenID = &tokenID
	return tx
}

// GetTokenID returns the token for which this account will get his KYC revoked.
func (tx *TokenRevokeKycTransaction) GetTokenID() TokenID {
	if tx.tokenID == nil {
		return TokenID{}
	}

	return *tx.tokenID
}

// SetAccountID Sets the account to be KYC Revoked
func (tx *TokenRevokeKycTransaction) SetAccountID(accountID AccountID) *TokenRevokeKycTransaction {
	tx._RequireNotFrozen()
	tx.accountID = &accountID
	return tx
}

// GetAccountID returns the AccountID that is being KYC Revoked
func (tx *TokenRevokeKycTransaction) GetAccountID() AccountID {
	if tx.accountID == nil {
		return AccountID{}
	}

	return *tx.accountID
}

// ---- Required Interfaces ---- //

// Sign uses the provided privateKey to sign the transaction.
func (tx *TokenRevokeKycTransaction) Sign(privateKey PrivateKey) *TokenRevokeKycTransaction {
	tx.Transaction.Sign(privateKey)
	return tx
}

// SignWithOperator signs the transaction with client's operator privateKey.
func (tx *TokenRevokeKycTransaction) SignWithOperator(client *Client) (*TokenRevokeKycTransaction, error) {
	_, err := tx.Transaction.signWithOperator(client, tx)
	if err != nil {
		return nil, err
	}
	return tx, nil
}

// SignWith executes the TransactionSigner and adds the resulting signature data to the Transaction's signature map
// with the publicKey as the map key.
func (tx *TokenRevokeKycTransaction) SignWith(
	publicKey PublicKey,
	signer TransactionSigner,
) *TokenRevokeKycTransaction {
	tx.Transaction.SignWith(publicKey, signer)
	return tx
}

// AddSignature adds a signature to the transaction.
func (tx *TokenRevokeKycTransaction) AddSignature(publicKey PublicKey, signature []byte) *TokenRevokeKycTransaction {
	tx.Transaction.AddSignature(publicKey, signature)
	return tx
}

// When execution is attempted, a single attempt will timeout when this deadline is reached. (The SDK may subsequently retry the execution.)
func (tx *TokenRevokeKycTransaction) SetGrpcDeadline(deadline *time.Duration) *TokenRevokeKycTransaction {
	tx.Transaction.SetGrpcDeadline(deadline)
	return tx
}

func (tx *TokenRevokeKycTransaction) Freeze() (*TokenRevokeKycTransaction, error) {
	return tx.FreezeWith(nil)
}

func (tx *TokenRevokeKycTransaction) FreezeWith(client *Client) (*TokenRevokeKycTransaction, error) {
	_, err := tx.Transaction.freezeWith(client, tx)
	return tx, err
}

// SetMaxTransactionFee sets the max transaction fee for this TokenRevokeKycTransaction.
func (tx *TokenRevokeKycTransaction) SetMaxTransactionFee(fee Hbar) *TokenRevokeKycTransaction {
	tx.Transaction.SetMaxTransactionFee(fee)
	return tx
}

// SetRegenerateTransactionID sets if transaction IDs should be regenerated when `TRANSACTION_EXPIRED` is received
func (tx *TokenRevokeKycTransaction) SetRegenerateTransactionID(regenerateTransactionID bool) *TokenRevokeKycTransaction {
	tx.Transaction.SetRegenerateTransactionID(regenerateTransactionID)
	return tx
}

// SetTransactionMemo sets the memo for this TokenRevokeKycTransaction.
func (tx *TokenRevokeKycTransaction) SetTransactionMemo(memo string) *TokenRevokeKycTransaction {
	tx.Transaction.SetTransactionMemo(memo)
	return tx
}

// SetTransactionValidDuration sets the valid duration for this TokenRevokeKycTransaction.
func (tx *TokenRevokeKycTransaction) SetTransactionValidDuration(duration time.Duration) *TokenRevokeKycTransaction {
	tx.Transaction.SetTransactionValidDuration(duration)
	return tx
}

// ToBytes serialise the tx to bytes, no matter if it is signed (locked), or not
func (tx *TokenRevokeKycTransaction) ToBytes() ([]byte, error) {
	bytes, err := tx.Transaction.toBytes(tx)
	if err != nil {
		return nil, err
	}
	return bytes, nil
}

// SetTransactionID sets the TransactionID for this TokenRevokeKycTransaction.
func (tx *TokenRevokeKycTransaction) SetTransactionID(transactionID TransactionID) *TokenRevokeKycTransaction {
	tx.Transaction.SetTransactionID(transactionID)
	return tx
}

// SetNodeAccountIDs sets the _Node AccountID for this TokenRevokeKycTransaction.
func (tx *TokenRevokeKycTransaction) SetNodeAccountIDs(nodeID []AccountID) *TokenRevokeKycTransaction {
	tx.Transaction.SetNodeAccountIDs(nodeID)
	return tx
}

// SetMaxRetry sets the max number of errors before execution will fail.
func (tx *TokenRevokeKycTransaction) SetMaxRetry(count int) *TokenRevokeKycTransaction {
	tx.Transaction.SetMaxRetry(count)
	return tx
}

// SetMaxBackoff The maximum amount of time to wait between retries.
// Every retry attempt will increase the wait time exponentially until it reaches this time.
func (tx *TokenRevokeKycTransaction) SetMaxBackoff(max time.Duration) *TokenRevokeKycTransaction {
	tx.Transaction.SetMaxBackoff(max)
	return tx
}

// SetMinBackoff sets the minimum amount of time to wait between retries.
func (tx *TokenRevokeKycTransaction) SetMinBackoff(min time.Duration) *TokenRevokeKycTransaction {
	tx.Transaction.SetMinBackoff(min)
	return tx
}

func (tx *TokenRevokeKycTransaction) SetLogLevel(level LogLevel) *TokenRevokeKycTransaction {
	tx.Transaction.SetLogLevel(level)
	return tx
}

func (tx *TokenRevokeKycTransaction) Execute(client *Client) (TransactionResponse, error) {
	return tx.Transaction.execute(client, tx)
}

func (tx *TokenRevokeKycTransaction) Schedule() (*ScheduleCreateTransaction, error) {
	return tx.Transaction.schedule(tx)
}

// ----------- Overridden functions ----------------

func (tx *TokenRevokeKycTransaction) getName() string {
	return "TokenRevokeKycTransaction"
}

func (tx *TokenRevokeKycTransaction) validateNetworkOnIDs(client *Client) error {
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

func (tx *TokenRevokeKycTransaction) build() *services.TransactionBody {
	return &services.TransactionBody{
		TransactionFee:           tx.transactionFee,
		Memo:                     tx.Transaction.memo,
		TransactionValidDuration: _DurationToProtobuf(tx.GetTransactionValidDuration()),
		TransactionID:            tx.transactionID._ToProtobuf(),
		Data: &services.TransactionBody_TokenRevokeKyc{
			TokenRevokeKyc: tx.buildProtoBody(),
		},
	}
}

func (tx *TokenRevokeKycTransaction) buildScheduled() (*services.SchedulableTransactionBody, error) {
	return &services.SchedulableTransactionBody{
		TransactionFee: tx.transactionFee,
		Memo:           tx.Transaction.memo,
		Data: &services.SchedulableTransactionBody_TokenRevokeKyc{
			TokenRevokeKyc: tx.buildProtoBody(),
		},
	}, nil
}

func (tx *TokenRevokeKycTransaction) buildProtoBody() *services.TokenRevokeKycTransactionBody {
	body := &services.TokenRevokeKycTransactionBody{}
	if tx.tokenID != nil {
		body.Token = tx.tokenID._ToProtobuf()
	}

	if tx.accountID != nil {
		body.Account = tx.accountID._ToProtobuf()
	}

	return body
}

func (tx *TokenRevokeKycTransaction) getMethod(channel *_Channel) _Method {
	return _Method{
		transaction: channel._GetToken().RevokeKycFromTokenAccount,
	}
}

func (tx *TokenRevokeKycTransaction) _ConstructScheduleProtobuf() (*services.SchedulableTransactionBody, error) {
	return tx.buildScheduled()
}
