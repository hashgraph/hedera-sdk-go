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

// TokenUnpauseTransaction
// Unpauses the Token. Must be signed with the Token's pause key.
// If the provided token is not found, the transaction will resolve to INVALID_TOKEN_ID.
// If the provided token has been deleted, the transaction will resolve to TOKEN_WAS_DELETED.
// If no Pause Key is defined, the transaction will resolve to TOKEN_HAS_NO_PAUSE_KEY.
// Once executed the Token is marked as Unpaused and can be used in Transactions.
// The operation is idempotent - becomes a no-op if the Token is already unpaused.
type TokenUnpauseTransaction struct {
	Transaction
	tokenID *TokenID
}

// NewTokenUnpauseTransaction creates TokenUnpauseTransaction which unpauses the Token.
// Must be signed with the Token's pause key.
// If the provided token is not found, the transaction will resolve to INVALID_TOKEN_ID.
// If the provided token has been deleted, the transaction will resolve to TOKEN_WAS_DELETED.
// If no Pause Key is defined, the transaction will resolve to TOKEN_HAS_NO_PAUSE_KEY.
// Once executed the Token is marked as Unpaused and can be used in Transactions.
// The operation is idempotent - becomes a no-op if the Token is already unpaused.
func NewTokenUnpauseTransaction() *TokenUnpauseTransaction {
	tx := TokenUnpauseTransaction{
		Transaction: _NewTransaction(),
	}

	tx._SetDefaultMaxTransactionFee(NewHbar(30))

	return &tx
}

func _TokenUnpauseTransactionFromProtobuf(tx Transaction, pb *services.TransactionBody) *TokenUnpauseTransaction {
	resultTx := &TokenUnpauseTransaction{
		Transaction: tx,
		tokenID:     _TokenIDFromProtobuf(pb.GetTokenDeletion().GetToken()),
	}
	return resultTx
}

// SetTokenID Sets the token to be unpaused.
func (tx *TokenUnpauseTransaction) SetTokenID(tokenID TokenID) *TokenUnpauseTransaction {
	tx._RequireNotFrozen()
	tx.tokenID = &tokenID
	return tx
}

// GetTokenID returns the TokenID of the token to be unpaused.
func (tx *TokenUnpauseTransaction) GetTokenID() TokenID {
	if tx.tokenID == nil {
		return TokenID{}
	}

	return *tx.tokenID
}

// ---- Required Interfaces ---- //

// Sign uses the provided privateKey to sign the transaction.
func (tx *TokenUnpauseTransaction) Sign(privateKey PrivateKey) *TokenUnpauseTransaction {
	tx.Transaction.Sign(privateKey)
	return tx
}

// SignWithOperator signs the transaction with client's operator privateKey.
func (tx *TokenUnpauseTransaction) SignWithOperator(client *Client) (*TokenUnpauseTransaction, error) {
	_, err := tx.Transaction.signWithOperator(client, tx)
	if err != nil {
		return nil, err
	}
	return tx, nil
}

// SignWith executes the TransactionSigner and adds the resulting signature data to the transaction's signature map
// with the publicKey as the map key.
func (tx *TokenUnpauseTransaction) SignWith(
	publicKey PublicKey,
	signer TransactionSigner,
) *TokenUnpauseTransaction {
	tx.Transaction.SignWith(publicKey, signer)
	return tx
}

// AddSignature adds a signature to the transaction.
func (tx *TokenUnpauseTransaction) AddSignature(publicKey PublicKey, signature []byte) *TokenUnpauseTransaction {
	tx.Transaction.AddSignature(publicKey, signature)
	return tx
}

// When execution is attempted, a single attempt will timeout when this deadline is reached. (The SDK may subsequently retry the execution.)
func (tx *TokenUnpauseTransaction) SetGrpcDeadline(deadline *time.Duration) *TokenUnpauseTransaction {
	tx.Transaction.SetGrpcDeadline(deadline)
	return tx
}

func (tx *TokenUnpauseTransaction) Freeze() (*TokenUnpauseTransaction, error) {
	return tx.FreezeWith(nil)
}

func (tx *TokenUnpauseTransaction) FreezeWith(client *Client) (*TokenUnpauseTransaction, error) {
	_, err := tx.Transaction.freezeWith(client, tx)
	return tx, err
}

// SetMaxTransactionFee sets the max transaction fee for this TokenUnpauseTransaction.
func (tx *TokenUnpauseTransaction) SetMaxTransactionFee(fee Hbar) *TokenUnpauseTransaction {
	tx.Transaction.SetMaxTransactionFee(fee)
	return tx
}

// SetRegenerateTransactionID sets if transaction IDs should be regenerated when `TRANSACTION_EXPIRED` is received
func (tx *TokenUnpauseTransaction) SetRegenerateTransactionID(regenerateTransactionID bool) *TokenUnpauseTransaction {
	tx.Transaction.SetRegenerateTransactionID(regenerateTransactionID)
	return tx
}

// SetTransactionMemo sets the memo for this TokenUnpauseTransaction.
func (tx *TokenUnpauseTransaction) SetTransactionMemo(memo string) *TokenUnpauseTransaction {
	tx.Transaction.SetTransactionMemo(memo)
	return tx
}

// SetTransactionValidDuration sets the valid duration for this TokenUnpauseTransaction.
func (tx *TokenUnpauseTransaction) SetTransactionValidDuration(duration time.Duration) *TokenUnpauseTransaction {
	tx.Transaction.SetTransactionValidDuration(duration)
	return tx
}

// SetTransactionID sets the TransactionID for this TokenUnpauseTransaction.
func (tx *TokenUnpauseTransaction) SetTransactionID(transactionID TransactionID) *TokenUnpauseTransaction {
	tx.Transaction.SetTransactionID(transactionID)
	return tx
}

// SetNodeAccountIDs sets the _Node AccountID for this TokenUnpauseTransaction.
func (tx *TokenUnpauseTransaction) SetNodeAccountIDs(nodeID []AccountID) *TokenUnpauseTransaction {
	tx.Transaction.SetNodeAccountIDs(nodeID)
	return tx
}

// SetMaxRetry sets the max number of errors before execution will fail.
func (tx *TokenUnpauseTransaction) SetMaxRetry(count int) *TokenUnpauseTransaction {
	tx.Transaction.SetMaxRetry(count)
	return tx
}

// SetMaxBackoff The maximum amount of time to wait between retries.
// Every retry attempt will increase the wait time exponentially until it reaches this time.
func (tx *TokenUnpauseTransaction) SetMaxBackoff(max time.Duration) *TokenUnpauseTransaction {
	tx.Transaction.SetMaxBackoff(max)
	return tx
}

// SetMinBackoff sets the minimum amount of time to wait between retries.
func (tx *TokenUnpauseTransaction) SetMinBackoff(min time.Duration) *TokenUnpauseTransaction {
	tx.Transaction.SetMinBackoff(min)
	return tx
}

func (tx *TokenUnpauseTransaction) SetLogLevel(level LogLevel) *TokenUnpauseTransaction {
	tx.Transaction.SetLogLevel(level)
	return tx
}

func (tx *TokenUnpauseTransaction) Execute(client *Client) (TransactionResponse, error) {
	return tx.Transaction.execute(client, tx)
}

func (tx *TokenUnpauseTransaction) Schedule() (*ScheduleCreateTransaction, error) {
	return tx.Transaction.schedule(tx)
}

// ----------- Overridden functions ----------------

func (tx *TokenUnpauseTransaction) getName() string {
	return "TokenUnpauseTransaction"
}

func (tx *TokenUnpauseTransaction) validateNetworkOnIDs(client *Client) error {
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

func (tx *TokenUnpauseTransaction) build() *services.TransactionBody {
	return &services.TransactionBody{
		TransactionFee:           tx.transactionFee,
		Memo:                     tx.Transaction.memo,
		TransactionValidDuration: _DurationToProtobuf(tx.GetTransactionValidDuration()),
		TransactionID:            tx.transactionID._ToProtobuf(),
		Data: &services.TransactionBody_TokenUnpause{
			TokenUnpause: tx.buildProtoBody(),
		},
	}
}

func (tx *TokenUnpauseTransaction) buildScheduled() (*services.SchedulableTransactionBody, error) { //nolint
	return &services.SchedulableTransactionBody{
		TransactionFee: tx.transactionFee,
		Memo:           tx.Transaction.memo,
		Data: &services.SchedulableTransactionBody_TokenUnpause{
			TokenUnpause: tx.buildProtoBody(),
		},
	}, nil
}

func (tx *TokenUnpauseTransaction) buildProtoBody() *services.TokenUnpauseTransactionBody { //nolint
	body := &services.TokenUnpauseTransactionBody{}
	if tx.tokenID != nil {
		body.Token = tx.tokenID._ToProtobuf()
	}

	return body
}

func (tx *TokenUnpauseTransaction) getMethod(channel *_Channel) _Method {
	return _Method{
		transaction: channel._GetToken().DeleteToken,
	}
}
func (tx *TokenUnpauseTransaction) _ConstructScheduleProtobuf() (*services.SchedulableTransactionBody, error) {
	return tx.buildScheduled()
}
