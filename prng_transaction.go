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

// PrngTransaction is used to generate a random number in a given range
type PrngTransaction struct {
	transaction
	rang uint32
}

// NewPrngTransaction creates a PrngTransaction transaction which can be used to construct and execute
// a Prng transaction.
func NewPrngTransaction() *PrngTransaction {
	tx := PrngTransaction{
		transaction: _NewTransaction(),
	}

	tx._SetDefaultMaxTransactionFee(NewHbar(5))
	tx.e = &tx

	return &tx
}

func _PrngTransactionFromProtobuf(tx transaction, pb *services.TransactionBody) *PrngTransaction {
	resultTx := &PrngTransaction{
		transaction: tx,
		rang:        uint32(pb.GetUtilPrng().GetRange()),
	}
	resultTx.e = resultTx
	return resultTx
}

// SetPayerAccountID Sets an optional id of the account to be charged the service fee for the scheduled transaction at
// the consensus time that it executes (if ever); defaults to the ScheduleCreate payer if not
// given
func (tx *PrngTransaction) SetRange(r uint32) *PrngTransaction {
	tx._RequireNotFrozen()
	tx.rang = r

	return tx
}

// GetRange returns the range of the prng
func (tx *PrngTransaction) GetRange() uint32 {
	return tx.rang
}

// ---- Required Interfaces ---- //

// Sign uses the provided privateKey to sign the transaction.
func (tx *PrngTransaction) Sign(privateKey PrivateKey) *PrngTransaction {
	tx.transaction.Sign(privateKey)
	return tx
}

// SignWithOperator signs the transaction with client's operator privateKey.
func (tx *PrngTransaction) SignWithOperator(client *Client) (*PrngTransaction, error) {
	_, err := tx.transaction.SignWithOperator(client)
	return tx, err
}

// SignWith executes the TransactionSigner and adds the resulting signature data to the transaction's signature map
// with the publicKey as the map key.
func (tx *PrngTransaction) SignWith(
	publicKey PublicKey,
	signer TransactionSigner,
) *PrngTransaction {
	tx.transaction.SignWith(publicKey, signer)
	return tx
}

// AddSignature adds a signature to the transaction.
func (tx *PrngTransaction) AddSignature(publicKey PublicKey, signature []byte) *PrngTransaction {
	tx.transaction.AddSignature(publicKey, signature)
	return tx
}

// When execution is attempted, a single attempt will timeout when this deadline is reached. (The SDK may subsequently retry the execution.)
func (tx *PrngTransaction) SetGrpcDeadline(deadline *time.Duration) *PrngTransaction {
	tx.transaction.SetGrpcDeadline(deadline)
	return tx
}

func (tx *PrngTransaction) Freeze() (*PrngTransaction, error) {
	return tx.FreezeWith(nil)
}

func (tx *PrngTransaction) FreezeWith(client *Client) (*PrngTransaction, error) {
	_, err := tx.transaction.FreezeWith(client)
	return tx, err
}

// SetMaxTransactionFee sets the max transaction fee for this PrngTransaction.
func (tx *PrngTransaction) SetMaxTransactionFee(fee Hbar) *PrngTransaction {
	tx.transaction.SetMaxTransactionFee(fee)
	return tx
}

// SetRegenerateTransactionID sets if transaction IDs should be regenerated when `TRANSACTION_EXPIRED` is received
func (tx *PrngTransaction) SetRegenerateTransactionID(regenerateTransactionID bool) *PrngTransaction {
	tx.transaction.SetRegenerateTransactionID(regenerateTransactionID)
	return tx
}

// SetTransactionMemo sets the memo for this PrngTransaction.
func (tx *PrngTransaction) SetTransactionMemo(memo string) *PrngTransaction {
	tx.transaction.SetTransactionMemo(memo)
	return tx
}

// SetTransactionValidDuration sets the valid duration for this PrngTransaction.
func (tx *PrngTransaction) SetTransactionValidDuration(duration time.Duration) *PrngTransaction {
	tx.transaction.SetTransactionValidDuration(duration)
	return tx
}

// SetTransactionID sets the TransactionID for this PrngTransaction.
func (tx *PrngTransaction) SetTransactionID(transactionID TransactionID) *PrngTransaction {
	tx.transaction.SetTransactionID(transactionID)
	return tx
}

// SetNodeAccountIDs sets the _Node AccountID for this PrngTransaction.
func (tx *PrngTransaction) SetNodeAccountIDs(nodeID []AccountID) *PrngTransaction {
	tx.transaction.SetNodeAccountIDs(nodeID)
	return tx
}

// SetMaxRetry sets the max number of errors before execution will fail.
func (tx *PrngTransaction) SetMaxRetry(count int) *PrngTransaction {
	tx.transaction.SetMaxRetry(count)
	return tx
}

// SetMaxBackoff The maximum amount of time to wait between retries.
// Every retry attempt will increase the wait time exponentially until it reaches this time.
func (tx *PrngTransaction) SetMaxBackoff(max time.Duration) *PrngTransaction {
	tx.transaction.SetMaxBackoff(max)
	return tx
}

// SetMinBackoff sets the minimum amount of time to wait between retries.
func (tx *PrngTransaction) SetMinBackoff(min time.Duration) *PrngTransaction {
	tx.transaction.SetMinBackoff(min)
	return tx
}

func (tx *PrngTransaction) SetLogLevel(level LogLevel) *PrngTransaction {
	tx.transaction.SetLogLevel(level)
	return tx
}

// ----------- overriden functions ----------------

func (tx *PrngTransaction) getName() string {
	return "PrngTransaction"
}

func (tx *PrngTransaction) build() *services.TransactionBody {
	return &services.TransactionBody{
		TransactionFee:           tx.transactionFee,
		Memo:                     tx.transaction.memo,
		TransactionValidDuration: _DurationToProtobuf(tx.GetTransactionValidDuration()),
		TransactionID:            tx.transactionID._ToProtobuf(),
		Data: &services.TransactionBody_UtilPrng{
			UtilPrng: tx.buildProtoBody(),
		},
	}
}

func (tx *PrngTransaction) buildScheduled() (*services.SchedulableTransactionBody, error) {
	return &services.SchedulableTransactionBody{
		TransactionFee: tx.transactionFee,
		Memo:           tx.transaction.memo,
		Data: &services.SchedulableTransactionBody_UtilPrng{
			UtilPrng: tx.buildProtoBody(),
		},
	}, nil
}

func (tx *PrngTransaction) buildProtoBody() *services.UtilPrngTransactionBody {
	body := &services.UtilPrngTransactionBody{
		Range: int32(tx.rang),
	}

	return body
}

func (tx *PrngTransaction) getMethod(channel *_Channel) _Method {
	return _Method{
		transaction: channel._GetUtil().Prng,
	}
}
func (tx *PrngTransaction) _ConstructScheduleProtobuf() (*services.SchedulableTransactionBody, error) {
	return tx.buildScheduled()
}
