package hedera

import (
	"time"

	"github.com/hashgraph/hedera-protobufs-go/services"
)

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

type TokenRejectTransaction struct {
	Transaction
}

func NewTokenRejectTransaction() *TokenRejectTransaction {
	tx := TokenRejectTransaction{
		Transaction: _NewTransaction(),
	}
	return &tx
}

// ---- Required Interfaces ---- //

// Sign uses the provided privateKey to sign the transaction.
func (tx *TokenRejectTransaction) Sign(privateKey PrivateKey) *TokenRejectTransaction {
	tx.Transaction.Sign(privateKey)
	return tx
}

// SignWithOperator signs the transaction with client's operator privateKey.
func (tx *TokenRejectTransaction) SignWithOperator(client *Client) (*TokenRejectTransaction, error) {
	_, err := tx.Transaction.signWithOperator(client, tx)
	if err != nil {
		return nil, err
	}
	return tx, nil
}

// SignWith executes the TransactionSigner and adds the resulting signature data to the Transaction's signature map
// with the publicKey as the map key.
func (tx *TokenRejectTransaction) SignWith(
	publicKey PublicKey,
	signer TransactionSigner,
) *TokenRejectTransaction {
	tx.Transaction.SignWith(publicKey, signer)
	return tx
}

// AddSignature adds a signature to the transaction.
func (tx *TokenRejectTransaction) AddSignature(publicKey PublicKey, signature []byte) *TokenRejectTransaction {
	tx.Transaction.AddSignature(publicKey, signature)
	return tx
}

// When execution is attempted, a single attempt will timeout when this deadline is reached. (The SDK may subsequently retry the execution.)
func (tx *TokenRejectTransaction) SetGrpcDeadline(deadline *time.Duration) *TokenRejectTransaction {
	tx.Transaction.SetGrpcDeadline(deadline)
	return tx
}

func (tx *TokenRejectTransaction) Freeze() (*TokenRejectTransaction, error) {
	return tx.FreezeWith(nil)
}

func (tx *TokenRejectTransaction) FreezeWith(client *Client) (*TokenRejectTransaction, error) {
	_, err := tx.Transaction.freezeWith(client, tx)
	return tx, err
}

// SetMaxTransactionFee sets the max transaction fee for this TokenRejectTransaction.
func (tx *TokenRejectTransaction) SetMaxTransactionFee(fee Hbar) *TokenRejectTransaction {
	tx.Transaction.SetMaxTransactionFee(fee)
	return tx
}

// SetRegenerateTransactionID sets if transaction IDs should be regenerated when `TRANSACTION_EXPIRED` is received
func (tx *TokenRejectTransaction) SetRegenerateTransactionID(regenerateTransactionID bool) *TokenRejectTransaction {
	tx.Transaction.SetRegenerateTransactionID(regenerateTransactionID)
	return tx
}

// SetTransactionMemo sets the memo for this TokenRejectTransaction.
func (tx *TokenRejectTransaction) SetTransactionMemo(memo string) *TokenRejectTransaction {
	tx.Transaction.SetTransactionMemo(memo)
	return tx
}

// SetTransactionValidDuration sets the valid duration for this TokenRejectTransaction.
func (tx *TokenRejectTransaction) SetTransactionValidDuration(duration time.Duration) *TokenRejectTransaction {
	tx.Transaction.SetTransactionValidDuration(duration)
	return tx
}

// ToBytes serialise the tx to bytes, no matter if it is signed (locked), or not
func (tx *TokenRejectTransaction) ToBytes() ([]byte, error) {
	bytes, err := tx.Transaction.toBytes(tx)
	if err != nil {
		return nil, err
	}
	return bytes, nil
}

// SetTransactionID sets the TransactionID for this TokenRejectTransaction.
func (tx *TokenRejectTransaction) SetTransactionID(transactionID TransactionID) *TokenRejectTransaction {
	tx.Transaction.SetTransactionID(transactionID)
	return tx
}

// SetNodeAccountIDs sets the _Node AccountID for this TokenRejectTransaction.
func (tx *TokenRejectTransaction) SetNodeAccountIDs(nodeID []AccountID) *TokenRejectTransaction {
	tx.Transaction.SetNodeAccountIDs(nodeID)
	return tx
}

// SetMaxRetry sets the max number of errors before execution will fail.
func (tx *TokenRejectTransaction) SetMaxRetry(count int) *TokenRejectTransaction {
	tx.Transaction.SetMaxRetry(count)
	return tx
}

// SetMaxBackoff The maximum amount of time to wait between retries.
// Every retry attempt will increase the wait time exponentially until it reaches this time.
func (tx *TokenRejectTransaction) SetMaxBackoff(max time.Duration) *TokenRejectTransaction {
	tx.Transaction.SetMaxBackoff(max)
	return tx
}

// SetMinBackoff sets the minimum amount of time to wait between retries.
func (tx *TokenRejectTransaction) SetMinBackoff(min time.Duration) *TokenRejectTransaction {
	tx.Transaction.SetMinBackoff(min)
	return tx
}

func (tx *TokenRejectTransaction) SetLogLevel(level LogLevel) *TokenRejectTransaction {
	tx.Transaction.SetLogLevel(level)
	return tx
}

func (tx *TokenRejectTransaction) Execute(client *Client) (TransactionResponse, error) {
	return tx.Transaction.execute(client, tx)
}

func (tx *TokenRejectTransaction) Schedule() (*ScheduleCreateTransaction, error) {
	return tx.Transaction.schedule(tx)
}

// ----------- Overridden functions ----------------

func (tx *TokenRejectTransaction) getName() string {
	return "TokenRejectTransaction"
}

func (tx *TokenRejectTransaction) validateNetworkOnIDs(client *Client) error {
	return nil
}

func (tx *TokenRejectTransaction) build() *services.TransactionBody {
	return &services.TransactionBody{
		TransactionFee:           tx.transactionFee,
		Memo:                     tx.Transaction.memo,
		TransactionValidDuration: _DurationToProtobuf(tx.GetTransactionValidDuration()),
		TransactionID:            tx.transactionID._ToProtobuf(),
		// Data:                     &services.TransactionBody_TokenCreation{
		// TokenCreation: tx.buildProtoBody(),
		// },
	}
}

func (tx *TokenRejectTransaction) buildScheduled() (*services.SchedulableTransactionBody, error) {
	return &services.SchedulableTransactionBody{
		TransactionFee: tx.transactionFee,
		Memo:           tx.Transaction.memo,
		// Data:           &services.SchedulableTransactionBody_TokenCreation{
		// TokenCreation: tx.buildProtoBody(),
		// },
	}, nil
}

// func (tx *TokenRejectTransaction) buildProtoBody() *services.TokenRejectTransaction {
// 	return nil
// }

func (tx *TokenRejectTransaction) getMethod(channel *_Channel) _Method {
	return _Method{
		transaction: channel._GetToken().CreateToken,
	}
}

func (tx *TokenRejectTransaction) preFreezeWith(client *Client) {
}

func (tx *TokenRejectTransaction) _ConstructScheduleProtobuf() (*services.SchedulableTransactionBody, error) {
	return tx.buildScheduled()
}
