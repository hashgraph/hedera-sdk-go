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

// TokenDeleteTransaction
// Marks a token as deleted, though it will remain in the ledger.
// The operation must be signed by the specified Admin Key of the Token. If
// admin key is not set, transaction will result in TOKEN_IS_IMMUTABlE.
// Once deleted update, mint, burn, wipe, freeze, unfreeze, grant kyc, revoke
// kyc and token transfer transactions will resolve to TOKEN_WAS_DELETED.
type TokenDeleteTransaction struct {
	transaction
	tokenID *TokenID
}

// NewTokenDeleteTransaction creates TokenDeleteTransaction which marks a token as deleted,
// though it will remain in the ledger.
// The operation must be signed by the specified Admin Key of the Token. If
// admin key is not set, transaction will result in TOKEN_IS_IMMUTABlE.
// Once deleted update, mint, burn, wipe, freeze, unfreeze, grant kyc, revoke
// kyc and token transfer transactions will resolve to TOKEN_WAS_DELETED.
func NewTokenDeleteTransaction() *TokenDeleteTransaction {
	tx := TokenDeleteTransaction{
		transaction: _NewTransaction(),
	}

	tx.e = &tx
	tx._SetDefaultMaxTransactionFee(NewHbar(30))

	return &tx
}

func _TokenDeleteTransactionFromProtobuf(tx transaction, pb *services.TransactionBody) *TokenDeleteTransaction {
	resultTx := &TokenDeleteTransaction{
		transaction: tx,
		tokenID:     _TokenIDFromProtobuf(pb.GetTokenDeletion().GetToken()),
	}
	resultTx.e = resultTx
	return resultTx
}

// SetTokenID Sets the Token to be deleted
func (tx *TokenDeleteTransaction) SetTokenID(tokenID TokenID) *TokenDeleteTransaction {
	tx._RequireNotFrozen()
	tx.tokenID = &tokenID
	return tx
}

// GetTokenID returns the TokenID of the token to be deleted
func (tx *TokenDeleteTransaction) GetTokenID() TokenID {
	if tx.tokenID == nil {
		return TokenID{}
	}

	return *tx.tokenID
}

// ---- Required Interfaces ---- //

// Sign uses the provided privateKey to sign the transaction.
func (tx *TokenDeleteTransaction) Sign(privateKey PrivateKey) *TokenDeleteTransaction {
	tx.transaction.Sign(privateKey)
	return tx
}

// SignWithOperator signs the transaction with client's operator privateKey.
func (tx *TokenDeleteTransaction) SignWithOperator(client *Client) (*TokenDeleteTransaction, error) {
	_, err := tx.transaction.SignWithOperator(client)
	return tx, err
}

// SignWith executes the TransactionSigner and adds the resulting signature data to the transaction's signature map
// with the publicKey as the map key.
func (tx *TokenDeleteTransaction) SignWith(
	publicKey PublicKey,
	signer TransactionSigner,
) *TokenDeleteTransaction {
	tx.transaction.SignWith(publicKey, signer)
	return tx
}

// AddSignature adds a signature to the transaction.
func (tx *TokenDeleteTransaction) AddSignature(publicKey PublicKey, signature []byte) *TokenDeleteTransaction {
	tx.transaction.AddSignature(publicKey, signature)
	return tx
}

// When execution is attempted, a single attempt will timeout when this deadline is reached. (The SDK may subsequently retry the execution.)
func (tx *TokenDeleteTransaction) SetGrpcDeadline(deadline *time.Duration) *TokenDeleteTransaction {
	tx.transaction.SetGrpcDeadline(deadline)
	return tx
}

func (tx *TokenDeleteTransaction) Freeze() (*TokenDeleteTransaction, error) {
	return tx.FreezeWith(nil)
}

func (tx *TokenDeleteTransaction) FreezeWith(client *Client) (*TokenDeleteTransaction, error) {
	_, err := tx.transaction.FreezeWith(client)
	return tx, err
}

// SetMaxTransactionFee sets the max transaction fee for this TokenDeleteTransaction.
func (tx *TokenDeleteTransaction) SetMaxTransactionFee(fee Hbar) *TokenDeleteTransaction {
	tx.transaction.SetMaxTransactionFee(fee)
	return tx
}

// SetRegenerateTransactionID sets if transaction IDs should be regenerated when `TRANSACTION_EXPIRED` is received
func (tx *TokenDeleteTransaction) SetRegenerateTransactionID(regenerateTransactionID bool) *TokenDeleteTransaction {
	tx.transaction.SetRegenerateTransactionID(regenerateTransactionID)
	return tx
}

// SetTransactionMemo sets the memo for this TokenDeleteTransaction.
func (tx *TokenDeleteTransaction) SetTransactionMemo(memo string) *TokenDeleteTransaction {
	tx.transaction.SetTransactionMemo(memo)
	return tx
}

// SetTransactionValidDuration sets the valid duration for this TokenDeleteTransaction.
func (tx *TokenDeleteTransaction) SetTransactionValidDuration(duration time.Duration) *TokenDeleteTransaction {
	tx.transaction.SetTransactionValidDuration(duration)
	return tx
}

// SetTransactionID sets the TransactionID for this TokenDeleteTransaction.
func (tx *TokenDeleteTransaction) SetTransactionID(transactionID TransactionID) *TokenDeleteTransaction {
	tx.transaction.SetTransactionID(transactionID)
	return tx
}

// SetNodeAccountIDs sets the _Node AccountID for this TokenDeleteTransaction.
func (tx *TokenDeleteTransaction) SetNodeAccountIDs(nodeID []AccountID) *TokenDeleteTransaction {
	tx.transaction.SetNodeAccountIDs(nodeID)
	return tx
}

// SetMaxRetry sets the max number of errors before execution will fail.
func (tx *TokenDeleteTransaction) SetMaxRetry(count int) *TokenDeleteTransaction {
	tx.transaction.SetMaxRetry(count)
	return tx
}

// SetMaxBackoff The maximum amount of time to wait between retries.
// Every retry attempt will increase the wait time exponentially until it reaches this time.
func (tx *TokenDeleteTransaction) SetMaxBackoff(max time.Duration) *TokenDeleteTransaction {
	tx.transaction.SetMaxBackoff(max)
	return tx
}

// SetMinBackoff sets the minimum amount of time to wait between retries.
func (tx *TokenDeleteTransaction) SetMinBackoff(min time.Duration) *TokenDeleteTransaction {
	tx.transaction.SetMinBackoff(min)
	return tx
}

func (tx *TokenDeleteTransaction) SetLogLevel(level LogLevel) *TokenDeleteTransaction {
	tx.transaction.SetLogLevel(level)
	return tx
}

// ----------- overriden functions ----------------

func (tx *TokenDeleteTransaction) getName() string {
	return "TokenDeleteTransaction"
}

func (tx *TokenDeleteTransaction) validateNetworkOnIDs(client *Client) error {
	if client == nil || !client.autoValidateChecksums {
		return nil
	}

	if tx.tokenID != nil {
		if err := tx.tokenID.ValidateChecksum(client); err != nil {
			return err
		}
	}

	return nil
}

func (tx *TokenDeleteTransaction) build() *services.TransactionBody {
	return &services.TransactionBody{
		TransactionFee:           tx.transactionFee,
		Memo:                     tx.transaction.memo,
		TransactionValidDuration: _DurationToProtobuf(tx.GetTransactionValidDuration()),
		TransactionID:            tx.transactionID._ToProtobuf(),
		Data: &services.TransactionBody_TokenDeletion{
			TokenDeletion: tx.buildProtoBody(),
		},
	}
}

func (tx *TokenDeleteTransaction) buildScheduled() (*services.SchedulableTransactionBody, error) {
	return &services.SchedulableTransactionBody{
		TransactionFee: tx.transactionFee,
		Memo:           tx.transaction.memo,
		Data: &services.SchedulableTransactionBody_TokenDeletion{
			TokenDeletion: tx.buildProtoBody(),
		},
	}, nil
}

func (tx *TokenDeleteTransaction) buildProtoBody() *services.TokenDeleteTransactionBody {
	body := &services.TokenDeleteTransactionBody{}
	if tx.tokenID != nil {
		body.Token = tx.tokenID._ToProtobuf()
	}

	return body
}

func (tx *TokenDeleteTransaction) getMethod(channel *_Channel) _Method {
	return _Method{
		transaction: channel._GetToken().DeleteToken,
	}
}
