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
	"github.com/hashgraph/hedera-protobufs-go/services"

	"time"
)

type FreezeTransaction struct {
	Transaction
	startTime  time.Time
	endTime    time.Time
	fileID     *FileID
	fileHash   []byte
	freezeType FreezeType
}

func NewFreezeTransaction() *FreezeTransaction {
	tx := FreezeTransaction{
		Transaction: _NewTransaction(),
	}

	tx._SetDefaultMaxTransactionFee(NewHbar(2))

	return &tx
}

func _FreezeTransactionFromProtobuf(tx Transaction, pb *services.TransactionBody) *FreezeTransaction {
	startTime := time.Date(
		time.Now().Year(), time.Now().Month(), time.Now().Day(),
		int(pb.GetFreeze().GetStartHour()), int(pb.GetFreeze().GetStartMin()), // nolint
		0, time.Now().Nanosecond(), time.Now().Location(),
	)

	endTime := time.Date(
		time.Now().Year(), time.Now().Month(), time.Now().Day(),
		int(pb.GetFreeze().GetEndHour()), int(pb.GetFreeze().GetEndMin()), // nolint
		0, time.Now().Nanosecond(), time.Now().Location(),
	)

	resultTx := &FreezeTransaction{
		Transaction: tx,
		startTime:   startTime,
		endTime:     endTime,
		fileID:      _FileIDFromProtobuf(pb.GetFreeze().GetUpdateFile()),
		fileHash:    pb.GetFreeze().FileHash,
	}
	return resultTx
}

func (tx *FreezeTransaction) SetStartTime(startTime time.Time) *FreezeTransaction {
	tx._RequireNotFrozen()
	tx.startTime = startTime
	return tx
}

func (tx *FreezeTransaction) GetStartTime() time.Time {
	return tx.startTime
}

// Deprecated
func (tx *FreezeTransaction) SetEndTime(endTime time.Time) *FreezeTransaction {
	tx._RequireNotFrozen()
	tx.endTime = endTime
	return tx
}

// Deprecated
func (tx *FreezeTransaction) GetEndTime() time.Time {
	return tx.endTime
}

func (tx *FreezeTransaction) SetFileID(id FileID) *FreezeTransaction {
	tx._RequireNotFrozen()
	tx.fileID = &id
	return tx
}

func (tx *FreezeTransaction) GetFileID() *FileID {
	return tx.fileID
}

func (tx *FreezeTransaction) SetFreezeType(freezeType FreezeType) *FreezeTransaction {
	tx._RequireNotFrozen()
	tx.freezeType = freezeType
	return tx
}

func (tx *FreezeTransaction) GetFreezeType() FreezeType {
	return tx.freezeType
}

func (tx *FreezeTransaction) SetFileHash(hash []byte) *FreezeTransaction {
	tx._RequireNotFrozen()
	tx.fileHash = hash
	return tx
}

func (tx *FreezeTransaction) GetFileHash() []byte {
	return tx.fileHash
}

// ---- Required Interfaces ---- //

// Sign uses the provided privateKey to sign the transaction.
func (tx *FreezeTransaction) Sign(
	privateKey PrivateKey,
) *FreezeTransaction {
	tx.Transaction.Sign(privateKey)
	return tx
}

// SignWithOperator signs the transaction with client's operator privateKey.
func (tx *FreezeTransaction) SignWithOperator(
	client *Client,
) (*FreezeTransaction, error) {
	// If the transaction is not signed by the _Operator, we need
	// to sign the transaction with the _Operator
	_, err := tx.Transaction.signWithOperator(client, tx)
	if err != nil {
		return nil, err
	}
	return tx, nil
}

// SignWith executes the TransactionSigner and adds the resulting signature data to the transaction's signature map
// with the publicKey as the map key.
func (tx *FreezeTransaction) SignWith(
	publicKey PublicKey,
	signer TransactionSigner,
) *FreezeTransaction {
	tx.Transaction.SignWith(publicKey, signer)
	return tx
}

// AddSignature adds a signature to the transaction.
func (tx *FreezeTransaction) AddSignature(publicKey PublicKey, signature []byte) *FreezeTransaction {
	tx.Transaction.AddSignature(publicKey, signature)
	return tx
}

func (tx *FreezeTransaction) Freeze() (*FreezeTransaction, error) {
	return tx.FreezeWith(nil)
}

func (tx *FreezeTransaction) FreezeWith(client *Client) (*FreezeTransaction, error) {
	_, err := tx.Transaction.freezeWith(client, tx)
	return tx, err
}

// GetMaxTransactionFee returns the maximum transaction fee the operator (paying account) is willing to pay.
func (tx *FreezeTransaction) GetMaxTransactionFee() Hbar {
	return tx.Transaction.GetMaxTransactionFee()
}

// SetMaxTransactionFee sets the maximum transaction fee the operator (paying account) is willing to pay.
func (tx *FreezeTransaction) SetMaxTransactionFee(fee Hbar) *FreezeTransaction {
	tx._RequireNotFrozen()
	tx.Transaction.SetMaxTransactionFee(fee)
	return tx
}

// SetRegenerateTransactionID sets if transaction IDs should be regenerated when `TRANSACTION_EXPIRED` is received
func (tx *FreezeTransaction) SetRegenerateTransactionID(regenerateTransactionID bool) *FreezeTransaction {
	tx._RequireNotFrozen()
	tx.Transaction.SetRegenerateTransactionID(regenerateTransactionID)
	return tx
}

// SetTransactionMemo sets the memo for tx FreezeTransaction.
func (tx *FreezeTransaction) SetTransactionMemo(memo string) *FreezeTransaction {
	tx._RequireNotFrozen()
	tx.Transaction.SetTransactionMemo(memo)
	return tx
}

// SetTransactionValidDuration sets the valid duration for tx FreezeTransaction.
func (tx *FreezeTransaction) SetTransactionValidDuration(duration time.Duration) *FreezeTransaction {
	tx._RequireNotFrozen()
	tx.Transaction.SetTransactionValidDuration(duration)
	return tx
}

// SetTransactionID sets the TransactionID for tx FreezeTransaction.
func (tx *FreezeTransaction) SetTransactionID(transactionID TransactionID) *FreezeTransaction {
	tx._RequireNotFrozen()

	tx.Transaction.SetTransactionID(transactionID)
	return tx
}

// SetNodeAccountID sets the _Node AccountID for tx FreezeTransaction.
func (tx *FreezeTransaction) SetNodeAccountIDs(nodeID []AccountID) *FreezeTransaction {
	tx._RequireNotFrozen()
	tx.Transaction.SetNodeAccountIDs(nodeID)
	return tx
}

// SetMaxRetry sets the max number of errors before execution will fail.
func (tx *FreezeTransaction) SetMaxRetry(count int) *FreezeTransaction {
	tx.Transaction.SetMaxRetry(count)
	return tx
}

// SetMaxBackoff The maximum amount of time to wait between retries.
// Every retry attempt will increase the wait time exponentially until it reaches tx time.
func (tx *FreezeTransaction) SetMaxBackoff(max time.Duration) *FreezeTransaction {
	tx.Transaction.SetMaxBackoff(max)
	return tx
}

// SetMinBackoff sets the minimum amount of time to wait between retries.
func (tx *FreezeTransaction) SetMinBackoff(min time.Duration) *FreezeTransaction {
	tx.Transaction.SetMinBackoff(min)
	return tx
}

func (tx *FreezeTransaction) SetLogLevel(level LogLevel) *FreezeTransaction {
	tx.Transaction.SetLogLevel(level)
	return tx
}

func (tx *FreezeTransaction) Execute(client *Client) (TransactionResponse, error) {
	return tx.Transaction.execute(client, tx)
}

func (tx *FreezeTransaction) Schedule() (*ScheduleCreateTransaction, error) {
	return tx.Transaction.schedule(tx)
}

// ----------- Overridden functions ----------------

func (tx *FreezeTransaction) getName() string {
	return "FreezeTransaction"
}
func (tx *FreezeTransaction) build() *services.TransactionBody {
	return &services.TransactionBody{
		TransactionFee:           tx.transactionFee,
		Memo:                     tx.Transaction.memo,
		TransactionValidDuration: _DurationToProtobuf(tx.GetTransactionValidDuration()),
		TransactionID:            tx.transactionID._ToProtobuf(),
		Data: &services.TransactionBody_Freeze{
			Freeze: tx.buildProtoBody(),
		},
	}
}
func (tx *FreezeTransaction) buildScheduled() (*services.SchedulableTransactionBody, error) {
	return &services.SchedulableTransactionBody{
		TransactionFee: tx.transactionFee,
		Memo:           tx.Transaction.memo,
		Data: &services.SchedulableTransactionBody_Freeze{
			Freeze: tx.buildProtoBody(),
		},
	}, nil
}
func (tx *FreezeTransaction) buildProtoBody() *services.FreezeTransactionBody {
	body := &services.FreezeTransactionBody{
		FileHash:   tx.fileHash,
		StartTime:  _TimeToProtobuf(tx.startTime),
		FreezeType: services.FreezeType(tx.freezeType),
	}

	if tx.fileID != nil {
		body.UpdateFile = tx.fileID._ToProtobuf()
	}

	return body
}
func (tx *FreezeTransaction) getMethod(channel *_Channel) _Method {
	return _Method{
		transaction: channel._GetFreeze().Freeze,
	}
}
func (tx *FreezeTransaction) _ConstructScheduleProtobuf() (*services.SchedulableTransactionBody, error) {
	return tx.buildScheduled()
}
