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

// AccountDeleteTransaction
// Mark an account as deleted, moving all its current hbars to another account. It will remain in
// the ledger, marked as deleted, until it expires. Transfers into it a deleted account fail. But a
// deleted account can still have its expiration extended in the normal way.
type AccountDeleteTransaction struct {
	Transaction
	transferAccountID *AccountID
	deleteAccountID   *AccountID
}

func _AccountDeleteTransactionFromProtobuf(transaction Transaction, pb *services.TransactionBody) *AccountDeleteTransaction {
	return &AccountDeleteTransaction{
		Transaction:       transaction,
		transferAccountID: _AccountIDFromProtobuf(pb.GetCryptoDelete().GetTransferAccountID()),
		deleteAccountID:   _AccountIDFromProtobuf(pb.GetCryptoDelete().GetDeleteAccountID()),
	}
}

// NewAccountDeleteTransaction creates AccountDeleteTransaction which marks an account as deleted, moving all its current hbars to another account. It will remain in
// the ledger, marked as deleted, until it expires. Transfers into it a deleted account fail. But a
// deleted account can still have its expiration extended in the normal way.
func NewAccountDeleteTransaction() *AccountDeleteTransaction {
	tx := AccountDeleteTransaction{
		Transaction: _NewTransaction(),
	}

	tx._SetDefaultMaxTransactionFee(NewHbar(2))

	return &tx
}

// SetNodeAccountID sets the _Node AccountID for this AccountDeleteTransaction.
func (tx *AccountDeleteTransaction) SetAccountID(accountID AccountID) *AccountDeleteTransaction {
	tx._RequireNotFrozen()
	tx.deleteAccountID = &accountID
	return tx
}

// GetAccountID returns the AccountID which will be deleted.
func (tx *AccountDeleteTransaction) GetAccountID() AccountID {
	if tx.deleteAccountID == nil {
		return AccountID{}
	}

	return *tx.deleteAccountID
}

// SetTransferAccountID sets the AccountID which will receive all remaining hbars.
func (tx *AccountDeleteTransaction) SetTransferAccountID(transferAccountID AccountID) *AccountDeleteTransaction {
	tx._RequireNotFrozen()
	tx.transferAccountID = &transferAccountID
	return tx
}

// GetTransferAccountID returns the AccountID which will receive all remaining hbars.
func (tx *AccountDeleteTransaction) GetTransferAccountID() AccountID {
	if tx.transferAccountID == nil {
		return AccountID{}
	}

	return *tx.transferAccountID
}

// ---- Required Interfaces ---- //

// Sign uses the provided privateKey to sign the transaction.
func (tx *AccountDeleteTransaction) Sign(
	privateKey PrivateKey,
) *AccountDeleteTransaction {
	tx.Transaction.Sign(privateKey)
	return tx
}

// SignWithOperator signs the transaction with client's operator privateKey.
func (tx *AccountDeleteTransaction) SignWithOperator(
	client *Client,
) (*AccountDeleteTransaction, error) {
	_, err := tx.Transaction.signWithOperator(client, tx)
	if err != nil {
		return nil, err
	}
	return tx, nil
}

// SignWith executes the TransactionSigner and adds the resulting signature data to the Transaction's signature map
// with the publicKey as the map key.
func (tx *AccountDeleteTransaction) SignWith(
	publicKey PublicKey,
	signer TransactionSigner,
) *AccountDeleteTransaction {
	tx.Transaction.SignWith(publicKey, signer)
	return tx
}

// AddSignature adds a signature to the transaction.
func (tx *AccountDeleteTransaction) AddSignature(publicKey PublicKey, signature []byte) *AccountDeleteTransaction {
	tx.Transaction.AddSignature(publicKey, signature)
	return tx
}

// When execution is attempted, a single attempt will timeout when this deadline is reached. (The SDK may subsequently retry the execution.)
func (tx *AccountDeleteTransaction) SetGrpcDeadline(deadline *time.Duration) *AccountDeleteTransaction {
	tx.Transaction.SetGrpcDeadline(deadline)
	return tx
}

func (tx *AccountDeleteTransaction) Freeze() (*AccountDeleteTransaction, error) {
	return tx.FreezeWith(nil)
}

func (tx *AccountDeleteTransaction) FreezeWith(client *Client) (*AccountDeleteTransaction, error) {
	_, err := tx.Transaction.freezeWith(client, tx)
	return tx, err
}

// SetMaxTransactionFee sets the maximum transaction fee the operator (paying account) is willing to pay.
func (tx *AccountDeleteTransaction) SetMaxTransactionFee(fee Hbar) *AccountDeleteTransaction {
	tx._RequireNotFrozen()
	tx.Transaction.SetMaxTransactionFee(fee)
	return tx
}

// SetRegenerateTransactionID sets if transaction IDs should be regenerated when `TRANSACTION_EXPIRED` is received
func (tx *AccountDeleteTransaction) SetRegenerateTransactionID(regenerateTransactionID bool) *AccountDeleteTransaction {
	tx._RequireNotFrozen()
	tx.Transaction.SetRegenerateTransactionID(regenerateTransactionID)
	return tx
}

// SetTransactionMemo sets the memo for this AccountDeleteTransaction.
func (tx *AccountDeleteTransaction) SetTransactionMemo(memo string) *AccountDeleteTransaction {
	tx._RequireNotFrozen()
	tx.Transaction.SetTransactionMemo(memo)
	return tx
}

// SetTransactionValidDuration sets the valid duration for this AccountDeleteTransaction.
func (tx *AccountDeleteTransaction) SetTransactionValidDuration(duration time.Duration) *AccountDeleteTransaction {
	tx._RequireNotFrozen()
	tx.Transaction.SetTransactionValidDuration(duration)
	return tx
}

// SetTransactionID sets the TransactionID for this AccountDeleteTransaction.
func (tx *AccountDeleteTransaction) SetTransactionID(transactionID TransactionID) *AccountDeleteTransaction {
	tx._RequireNotFrozen()

	tx.Transaction.SetTransactionID(transactionID)
	return tx
}

// SetNodeAccountIDs sets the _Node AccountID for this AccountDeleteTransaction.
func (tx *AccountDeleteTransaction) SetNodeAccountIDs(nodeID []AccountID) *AccountDeleteTransaction {
	tx._RequireNotFrozen()
	tx.Transaction.SetNodeAccountIDs(nodeID)
	return tx
}

// SetMaxRetry sets the max number of errors before execution will fail.
func (tx *AccountDeleteTransaction) SetMaxRetry(count int) *AccountDeleteTransaction {
	tx.Transaction.SetMaxRetry(count)
	return tx
}

// SetMaxBackoff The maximum amount of time to wait between retries. Every retry attempt will increase the wait time exponentially until it reaches this time.
func (tx *AccountDeleteTransaction) SetMaxBackoff(max time.Duration) *AccountDeleteTransaction {
	tx.Transaction.SetMaxBackoff(max)
	return tx
}

// GetMaxBackoff returns the maximum amount of time to wait between retries.
func (tx *AccountDeleteTransaction) SetMinBackoff(min time.Duration) *AccountDeleteTransaction {
	tx.Transaction.SetMinBackoff(min)
	return tx
}

func (tx *AccountDeleteTransaction) SetLogLevel(level LogLevel) *AccountDeleteTransaction {
	tx.Transaction.SetLogLevel(level)
	return tx
}

func (tx *AccountDeleteTransaction) Execute(client *Client) (TransactionResponse, error) {
	return tx.Transaction.execute(client, tx)
}

func (tx *AccountDeleteTransaction) Schedule() (*ScheduleCreateTransaction, error) {
	return tx.Transaction.schedule(tx)
}

// ----------- Overridden functions ----------------

func (tx *AccountDeleteTransaction) getName() string {
	return "AccountDeleteTransaction"
}
func (tx *AccountDeleteTransaction) validateNetworkOnIDs(client *Client) error {
	if client == nil || !client.autoValidateChecksums {
		return nil
	}

	if tx.deleteAccountID != nil {
		if err := tx.deleteAccountID.ValidateChecksum(client); err != nil {
			return err
		}
	}

	if tx.transferAccountID != nil {
		if err := tx.transferAccountID.ValidateChecksum(client); err != nil {
			return err
		}
	}

	return nil
}

func (tx *AccountDeleteTransaction) build() *services.TransactionBody {
	return &services.TransactionBody{
		TransactionFee:           tx.transactionFee,
		Memo:                     tx.Transaction.memo,
		TransactionValidDuration: _DurationToProtobuf(tx.GetTransactionValidDuration()),
		TransactionID:            tx.transactionID._ToProtobuf(),
		Data: &services.TransactionBody_CryptoDelete{
			CryptoDelete: tx.buildProtoBody(),
		},
	}
}

func (tx *AccountDeleteTransaction) buildScheduled() (*services.SchedulableTransactionBody, error) {
	return &services.SchedulableTransactionBody{
		TransactionFee: tx.transactionFee,
		Memo:           tx.Transaction.memo,
		Data: &services.SchedulableTransactionBody_CryptoDelete{
			CryptoDelete: tx.buildProtoBody(),
		},
	}, nil
}

func (tx *AccountDeleteTransaction) buildProtoBody() *services.CryptoDeleteTransactionBody {
	body := &services.CryptoDeleteTransactionBody{}

	if tx.transferAccountID != nil {
		body.TransferAccountID = tx.transferAccountID._ToProtobuf()
	}

	if tx.deleteAccountID != nil {
		body.DeleteAccountID = tx.deleteAccountID._ToProtobuf()
	}

	return body
}

func (tx *AccountDeleteTransaction) getMethod(channel *_Channel) _Method {
	return _Method{
		transaction: channel._GetCrypto().ApproveAllowances,
	}
}

func (tx *AccountDeleteTransaction) _ConstructScheduleProtobuf() (*services.SchedulableTransactionBody, error) {
	return tx.buildScheduled()
}
