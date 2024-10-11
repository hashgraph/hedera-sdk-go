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

	"github.com/hashgraph/hedera-sdk-go/v2/generated/services"
)

// TokenDeleteTransaction
// Marks a token as deleted, though it will remain in the ledger.
// The operation must be signed by the specified Admin Key of the Token. If
// admin key is not set, transaction will result in TOKEN_IS_IMMUTABlE.
// Once deleted update, mint, burn, wipe, freeze, unfreeze, grant kyc, revoke
// kyc and token transfer transactions will resolve to TOKEN_WAS_DELETED.
type TokenDeleteTransaction struct {
	Transaction
	tokenID *TokenID
}

// NewTokenDeleteTransaction creates TokenDeleteTransaction which marks a token as deleted,
// though it will remain in the ledger.
// The operation must be signed by the specified Admin Key of the Token. If
// admin key is not set, Transaction will result in TOKEN_IS_IMMUTABlE.
// Once deleted update, mint, burn, wipe, freeze, unfreeze, grant kyc, revoke
// kyc and token transfer transactions will resolve to TOKEN_WAS_DELETED.
func NewTokenDeleteTransaction() *TokenDeleteTransaction {
	tx := TokenDeleteTransaction{
		Transaction: _NewTransaction(),
	}

	tx._SetDefaultMaxTransactionFee(NewHbar(30))

	return &tx
}

func _TokenDeleteTransactionFromProtobuf(tx Transaction, pb *services.TransactionBody) *TokenDeleteTransaction {
	return &TokenDeleteTransaction{
		Transaction: tx,
		tokenID:     _TokenIDFromProtobuf(pb.GetTokenDeletion().GetToken()),
	}
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
	tx.Transaction.Sign(privateKey)
	return tx
}

// SignWithOperator signs the transaction with client's operator privateKey.
func (tx *TokenDeleteTransaction) SignWithOperator(client *Client) (*TokenDeleteTransaction, error) {
	_, err := tx.Transaction.signWithOperator(client, tx)
	if err != nil {
		return nil, err
	}
	return tx, nil
}

// SignWith executes the TransactionSigner and adds the resulting signature data to the Transaction's signature map
// with the publicKey as the map key.
func (tx *TokenDeleteTransaction) SignWith(
	publicKey PublicKey,
	signer TransactionSigner,
) *TokenDeleteTransaction {
	tx.Transaction.SignWith(publicKey, signer)
	return tx
}

// AddSignature adds a signature to the transaction.
func (tx *TokenDeleteTransaction) AddSignature(publicKey PublicKey, signature []byte) *TokenDeleteTransaction {
	tx.Transaction.AddSignature(publicKey, signature)
	return tx
}

// When execution is attempted, a single attempt will timeout when this deadline is reached. (The SDK may subsequently retry the execution.)
func (tx *TokenDeleteTransaction) SetGrpcDeadline(deadline *time.Duration) *TokenDeleteTransaction {
	tx.Transaction.SetGrpcDeadline(deadline)
	return tx
}

func (tx *TokenDeleteTransaction) Freeze() (*TokenDeleteTransaction, error) {
	return tx.FreezeWith(nil)
}

func (tx *TokenDeleteTransaction) FreezeWith(client *Client) (*TokenDeleteTransaction, error) {
	_, err := tx.Transaction.freezeWith(client, tx)
	return tx, err
}

// SetMaxTransactionFee sets the max transaction fee for this TokenDeleteTransaction.
func (tx *TokenDeleteTransaction) SetMaxTransactionFee(fee Hbar) *TokenDeleteTransaction {
	tx.Transaction.SetMaxTransactionFee(fee)
	return tx
}

// SetRegenerateTransactionID sets if transaction IDs should be regenerated when `TRANSACTION_EXPIRED` is received
func (tx *TokenDeleteTransaction) SetRegenerateTransactionID(regenerateTransactionID bool) *TokenDeleteTransaction {
	tx.Transaction.SetRegenerateTransactionID(regenerateTransactionID)
	return tx
}

// SetTransactionMemo sets the memo for this TokenDeleteTransaction.
func (tx *TokenDeleteTransaction) SetTransactionMemo(memo string) *TokenDeleteTransaction {
	tx.Transaction.SetTransactionMemo(memo)
	return tx
}

// SetTransactionValidDuration sets the valid duration for this TokenDeleteTransaction.
func (tx *TokenDeleteTransaction) SetTransactionValidDuration(duration time.Duration) *TokenDeleteTransaction {
	tx.Transaction.SetTransactionValidDuration(duration)
	return tx
}

// ToBytes serialise the tx to bytes, no matter if it is signed (locked), or not
func (tx *TokenDeleteTransaction) ToBytes() ([]byte, error) {
	bytes, err := tx.Transaction.toBytes(tx)
	if err != nil {
		return nil, err
	}
	return bytes, nil
}

// SetTransactionID sets the TransactionID for this TokenDeleteTransaction.
func (tx *TokenDeleteTransaction) SetTransactionID(transactionID TransactionID) *TokenDeleteTransaction {
	tx.Transaction.SetTransactionID(transactionID)
	return tx
}

// SetNodeAccountIDs sets the _Node AccountID for this TokenDeleteTransaction.
func (tx *TokenDeleteTransaction) SetNodeAccountIDs(nodeID []AccountID) *TokenDeleteTransaction {
	tx.Transaction.SetNodeAccountIDs(nodeID)
	return tx
}

// SetMaxRetry sets the max number of errors before execution will fail.
func (tx *TokenDeleteTransaction) SetMaxRetry(count int) *TokenDeleteTransaction {
	tx.Transaction.SetMaxRetry(count)
	return tx
}

// SetMaxBackoff The maximum amount of time to wait between retries.
// Every retry attempt will increase the wait time exponentially until it reaches this time.
func (tx *TokenDeleteTransaction) SetMaxBackoff(max time.Duration) *TokenDeleteTransaction {
	tx.Transaction.SetMaxBackoff(max)
	return tx
}

// SetMinBackoff sets the minimum amount of time to wait between retries.
func (tx *TokenDeleteTransaction) SetMinBackoff(min time.Duration) *TokenDeleteTransaction {
	tx.Transaction.SetMinBackoff(min)
	return tx
}

func (tx *TokenDeleteTransaction) SetLogLevel(level LogLevel) *TokenDeleteTransaction {
	tx.Transaction.SetLogLevel(level)
	return tx
}

func (tx *TokenDeleteTransaction) Execute(client *Client) (TransactionResponse, error) {
	return tx.Transaction.execute(client, tx)
}

func (tx *TokenDeleteTransaction) Schedule() (*ScheduleCreateTransaction, error) {
	return tx.Transaction.schedule(tx)
}

// ----------- Overridden functions ----------------

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
		Memo:                     tx.Transaction.memo,
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
		Memo:           tx.Transaction.memo,
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
func (tx *TokenDeleteTransaction) _ConstructScheduleProtobuf() (*services.SchedulableTransactionBody, error) {
	return tx.buildScheduled()
}
