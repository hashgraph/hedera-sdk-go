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
	"errors"

	"github.com/hashgraph/hedera-protobufs-go/services"

	"time"
)

// LiveHashDeleteTransaction At consensus, deletes a livehash associated to the given account. The transaction must be signed
// by either the key of the owning account, or at least one of the keys associated to the livehash.
type LiveHashDeleteTransaction struct {
	Transaction
	accountID *AccountID
	hash      []byte
}

// NewLiveHashDeleteTransaction creates LiveHashDeleteTransaction which at consensus, deletes a livehash associated to the given account.
// The transaction must be signed by either the key of the owning account, or at least one of the keys associated to the livehash.
func NewLiveHashDeleteTransaction() *LiveHashDeleteTransaction {
	tx := LiveHashDeleteTransaction{
		Transaction: _NewTransaction(),
	}
	tx._SetDefaultMaxTransactionFee(NewHbar(2))

	return &tx
}

func _LiveHashDeleteTransactionFromProtobuf(tx Transaction, pb *services.TransactionBody) *LiveHashDeleteTransaction {
	return &LiveHashDeleteTransaction{
		Transaction: tx,
		accountID:   _AccountIDFromProtobuf(pb.GetCryptoDeleteLiveHash().GetAccountOfLiveHash()),
		hash:        pb.GetCryptoDeleteLiveHash().LiveHashToDelete,
	}
}

// SetHash Set the SHA-384 livehash to delete from the account
func (tx *LiveHashDeleteTransaction) SetHash(hash []byte) *LiveHashDeleteTransaction {
	tx._RequireNotFrozen()
	tx.hash = hash
	return tx
}

// GetHash returns the SHA-384 livehash to delete from the account
func (tx *LiveHashDeleteTransaction) GetHash() []byte {
	return tx.hash
}

// SetAccountID Sets the account owning the livehash
func (tx *LiveHashDeleteTransaction) SetAccountID(accountID AccountID) *LiveHashDeleteTransaction {
	tx._RequireNotFrozen()
	tx.accountID = &accountID
	return tx
}

// GetAccountID returns the account owning the livehash
func (tx *LiveHashDeleteTransaction) GetAccountID() AccountID {
	if tx.accountID == nil {
		return AccountID{}
	}

	return *tx.accountID
}

// ---- Required Interfaces ---- //

// Sign uses the provided privateKey to sign the transaction.
func (tx *LiveHashDeleteTransaction) Sign(privateKey PrivateKey) *LiveHashDeleteTransaction {
	tx.Transaction.Sign(privateKey)
	return tx
}

// SignWithOperator signs the transaction with client's operator privateKey.
func (tx *LiveHashDeleteTransaction) SignWithOperator(client *Client) (*LiveHashDeleteTransaction, error) {
	_, err := tx.Transaction.signWithOperator(client, tx)
	if err != nil {
		return nil, err
	}
	return tx, nil
}

// SignWith executes the TransactionSigner and adds the resulting signature data to the Transaction's signature map
// with the publicKey as the map key.
func (tx *LiveHashDeleteTransaction) SignWith(
	publicKey PublicKey,
	signer TransactionSigner,
) *LiveHashDeleteTransaction {
	tx.Transaction.SignWith(publicKey, signer)
	return tx
}

// AddSignature adds a signature to the transaction.
func (tx *LiveHashDeleteTransaction) AddSignature(publicKey PublicKey, signature []byte) *LiveHashDeleteTransaction {
	tx.Transaction.AddSignature(publicKey, signature)
	return tx
}

// When execution is attempted, a single attempt will timeout when this deadline is reached. (The SDK may subsequently retry the execution.)
func (tx *LiveHashDeleteTransaction) SetGrpcDeadline(deadline *time.Duration) *LiveHashDeleteTransaction {
	tx.Transaction.SetGrpcDeadline(deadline)
	return tx
}

func (tx *LiveHashDeleteTransaction) Freeze() (*LiveHashDeleteTransaction, error) {
	return tx.FreezeWith(nil)
}

func (tx *LiveHashDeleteTransaction) FreezeWith(client *Client) (*LiveHashDeleteTransaction, error) {
	_, err := tx.Transaction.freezeWith(client, tx)
	return tx, err
}

// SetMaxTransactionFee sets the maximum transaction fee for this LiveHashDeleteTransaction.
func (tx *LiveHashDeleteTransaction) SetMaxTransactionFee(fee Hbar) *LiveHashDeleteTransaction {
	tx.Transaction.SetMaxTransactionFee(fee)
	return tx
}

// SetRegenerateTransactionID sets if transaction IDs should be regenerated when `TRANSACTION_EXPIRED` is received
func (tx *LiveHashDeleteTransaction) SetRegenerateTransactionID(regenerateTransactionID bool) *LiveHashDeleteTransaction {
	tx.Transaction.SetRegenerateTransactionID(regenerateTransactionID)
	return tx
}

// SetTransactionMemo sets the memo for this LiveHashDeleteTransaction.
func (tx *LiveHashDeleteTransaction) SetTransactionMemo(memo string) *LiveHashDeleteTransaction {
	tx.Transaction.SetTransactionMemo(memo)
	return tx
}

// SetTransactionValidDuration sets the valid duration for this LiveHashDeleteTransaction.
func (tx *LiveHashDeleteTransaction) SetTransactionValidDuration(duration time.Duration) *LiveHashDeleteTransaction {
	tx.Transaction.SetTransactionValidDuration(duration)
	return tx
}

// SetTransactionID sets the TransactionID for this LiveHashDeleteTransaction.
func (tx *LiveHashDeleteTransaction) SetTransactionID(transactionID TransactionID) *LiveHashDeleteTransaction {
	tx.Transaction.SetTransactionID(transactionID)
	return tx
}

// SetNodeAccountIDs sets the _Node AccountID for this LiveHashDeleteTransaction.
func (tx *LiveHashDeleteTransaction) SetNodeAccountIDs(nodeID []AccountID) *LiveHashDeleteTransaction {
	tx.Transaction.SetNodeAccountIDs(nodeID)
	return tx
}

// SetMaxRetry sets the max number of errors before execution will fail.
func (tx *LiveHashDeleteTransaction) SetMaxRetry(count int) *LiveHashDeleteTransaction {
	tx.Transaction.SetMaxRetry(count)
	return tx
}

// SetMaxBackoff The maximum amount of time to wait between retries.
// Every retry attempt will increase the wait time exponentially until it reaches this time.
func (tx *LiveHashDeleteTransaction) SetMaxBackoff(max time.Duration) *LiveHashDeleteTransaction {
	tx.Transaction.SetMaxBackoff(max)
	return tx
}

// SetMinBackoff sets the minimum amount of time to wait between retries.
func (tx *LiveHashDeleteTransaction) SetMinBackoff(min time.Duration) *LiveHashDeleteTransaction {
	tx.Transaction.SetMinBackoff(min)
	return tx
}

func (tx *LiveHashDeleteTransaction) SetLogLevel(level LogLevel) *LiveHashDeleteTransaction {
	tx.Transaction.SetLogLevel(level)
	return tx
}

func (tx *LiveHashDeleteTransaction) Execute(client *Client) (TransactionResponse, error) {
	return tx.Transaction.execute(client, tx)
}

func (tx *LiveHashDeleteTransaction) Schedule() (*ScheduleCreateTransaction, error) {
	return tx.Transaction.schedule(tx)
}

// ----------- Overridden functions ----------------

func (tx *LiveHashDeleteTransaction) getName() string {
	return "LiveHashDeleteTransaction"
}

func (tx *LiveHashDeleteTransaction) validateNetworkOnIDs(client *Client) error {
	if client == nil || !client.autoValidateChecksums {
		return nil
	}

	if tx.accountID != nil {
		if err := tx.accountID.ValidateChecksum(client); err != nil {
			return err
		}
	}

	return nil
}

func (tx *LiveHashDeleteTransaction) build() *services.TransactionBody {
	return &services.TransactionBody{
		TransactionFee:           tx.transactionFee,
		Memo:                     tx.Transaction.memo,
		TransactionValidDuration: _DurationToProtobuf(tx.GetTransactionValidDuration()),
		TransactionID:            tx.transactionID._ToProtobuf(),
		Data: &services.TransactionBody_CryptoDeleteLiveHash{
			CryptoDeleteLiveHash: tx.buildProtoBody(),
		},
	}
}

func (tx *LiveHashDeleteTransaction) buildScheduled() (*services.SchedulableTransactionBody, error) {
	return nil, errors.New("cannot schedule `LiveHashDeleteTransaction`")
}

func (tx *LiveHashDeleteTransaction) buildProtoBody() *services.CryptoDeleteLiveHashTransactionBody {
	body := &services.CryptoDeleteLiveHashTransactionBody{}

	if tx.accountID != nil {
		body.AccountOfLiveHash = tx.accountID._ToProtobuf()
	}

	if tx.hash != nil {
		body.LiveHashToDelete = tx.hash
	}

	return body
}

func (tx *LiveHashDeleteTransaction) getMethod(channel *_Channel) _Method {
	return _Method{
		transaction: channel._GetCrypto().DeleteLiveHash,
	}
}
func (tx *LiveHashDeleteTransaction) _ConstructScheduleProtobuf() (*services.SchedulableTransactionBody, error) {
	return tx.buildScheduled()
}
