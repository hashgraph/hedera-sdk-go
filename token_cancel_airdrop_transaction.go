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

type TokenCancelAirdropTransaction struct {
	Transaction
	pendingAirdropIds []*PendingAirdropId
}

func NewTokenCancelAirdropTransaction() *TokenCancelAirdropTransaction {
	tx := TokenCancelAirdropTransaction{
		Transaction:       _NewTransaction(),
		pendingAirdropIds: make([]*PendingAirdropId, 0),
	}

	tx._SetDefaultMaxTransactionFee(NewHbar(1))

	return &tx
}

func _TokenCancelAirdropTransactionFromProtobuf(tx Transaction, pb *services.TransactionBody) *TokenCancelAirdropTransaction {
	TokenCancel := &TokenCancelAirdropTransaction{
		Transaction: tx,
	}

	for _, pendingAirdrops := range pb.GetTokenCancelAirdrop().PendingAirdrops {
		TokenCancel.pendingAirdropIds = append(TokenCancel.pendingAirdropIds, _PendingAirdropIdFromProtobuf(pendingAirdrops))
	}

	return TokenCancel
}

// SetPendingAirdropIds sets the pending airdrop IDs for this TokenCancelAirdropTransaction.
func (tx *TokenCancelAirdropTransaction) SetPendingAirdropIds(ids []*PendingAirdropId) *TokenCancelAirdropTransaction {
	tx._RequireNotFrozen()
	tx.pendingAirdropIds = ids
	return tx
}

// AddPendingAirdropId adds a pending airdrop ID to this TokenCancelAirdropTransaction.
func (tx *TokenCancelAirdropTransaction) AddPendingAirdropId(id PendingAirdropId) *TokenCancelAirdropTransaction {
	tx._RequireNotFrozen()
	tx.pendingAirdropIds = append(tx.pendingAirdropIds, &id)
	return tx
}

// GetPendingAirdropIds returns the pending airdrop IDs for this TokenCancelAirdropTransaction.
func (tx *TokenCancelAirdropTransaction) GetPendingAirdropIds() []*PendingAirdropId {
	return tx.pendingAirdropIds
}

// ---- Required Interfaces ---- //

// Sign uses the provided privateKey to sign the transaction.
func (tx *TokenCancelAirdropTransaction) Sign(privateKey PrivateKey) *TokenCancelAirdropTransaction {
	tx.Transaction.Sign(privateKey)
	return tx
}

// SignWithOperator signs the transaction with client's operator privateKey.
func (tx *TokenCancelAirdropTransaction) SignWithOperator(client *Client) (*TokenCancelAirdropTransaction, error) {
	_, err := tx.Transaction.signWithOperator(client, tx)
	if err != nil {
		return nil, err
	}
	return tx, nil
}

// SignWith executes the TransactionSigner and adds the resulting signature data to the Transaction's signature map
// with the publicKey as the map key.
func (tx *TokenCancelAirdropTransaction) SignWith(
	publicKey PublicKey,
	signer TransactionSigner,
) *TokenCancelAirdropTransaction {
	tx.Transaction.SignWith(publicKey, signer)
	return tx
}

// AddSignature adds a signature to the transaction.
func (tx *TokenCancelAirdropTransaction) AddSignature(publicKey PublicKey, signature []byte) *TokenCancelAirdropTransaction {
	tx.Transaction.AddSignature(publicKey, signature)
	return tx
}

// When execution is attempted, a single attempt will timeout when this deadline is reached. (The SDK may subsequently retry the execution.)
func (tx *TokenCancelAirdropTransaction) SetGrpcDeadline(deadline *time.Duration) *TokenCancelAirdropTransaction {
	tx.Transaction.SetGrpcDeadline(deadline)
	return tx
}

func (tx *TokenCancelAirdropTransaction) Freeze() (*TokenCancelAirdropTransaction, error) {
	return tx.FreezeWith(nil)
}

func (tx *TokenCancelAirdropTransaction) FreezeWith(client *Client) (*TokenCancelAirdropTransaction, error) {
	_, err := tx.Transaction.freezeWith(client, tx)
	return tx, err
}

// SetMaxTransactionFee sets the max transaction fee for this TokenCancelAirdropTransaction.
func (tx *TokenCancelAirdropTransaction) SetMaxTransactionFee(fee Hbar) *TokenCancelAirdropTransaction {
	tx.Transaction.SetMaxTransactionFee(fee)
	return tx
}

// SetRegenerateTransactionID sets if transaction IDs should be regenerated when `TRANSACTION_EXPIRED` is received
func (tx *TokenCancelAirdropTransaction) SetRegenerateTransactionID(regenerateTransactionID bool) *TokenCancelAirdropTransaction {
	tx.Transaction.SetRegenerateTransactionID(regenerateTransactionID)
	return tx
}

// SetTransactionMemo sets the memo for this TokenCancelAirdropTransaction.
func (tx *TokenCancelAirdropTransaction) SetTransactionMemo(memo string) *TokenCancelAirdropTransaction {
	tx.Transaction.SetTransactionMemo(memo)
	return tx
}

// SetTransactionValidDuration sets the valid duration for this TokenCancelAirdropTransaction.
func (tx *TokenCancelAirdropTransaction) SetTransactionValidDuration(duration time.Duration) *TokenCancelAirdropTransaction {
	tx.Transaction.SetTransactionValidDuration(duration)
	return tx
}

// ToBytes serialise the tx to bytes, no matter if it is signed (locked), or not
func (tx *TokenCancelAirdropTransaction) ToBytes() ([]byte, error) {
	bytes, err := tx.Transaction.toBytes(tx)
	if err != nil {
		return nil, err
	}
	return bytes, nil
}

// SetTransactionID sets the TransactionID for this TokenCancelAirdropTransaction.
func (tx *TokenCancelAirdropTransaction) SetTransactionID(transactionID TransactionID) *TokenCancelAirdropTransaction {
	tx.Transaction.SetTransactionID(transactionID)
	return tx
}

// SetNodeAccountIDs sets the _Node AccountID for this TokenCancelAirdropTransaction.
func (tx *TokenCancelAirdropTransaction) SetNodeAccountIDs(nodeID []AccountID) *TokenCancelAirdropTransaction {
	tx.Transaction.SetNodeAccountIDs(nodeID)
	return tx
}

// SetMaxRetry sets the max number of errors before execution will fail.
func (tx *TokenCancelAirdropTransaction) SetMaxRetry(count int) *TokenCancelAirdropTransaction {
	tx.Transaction.SetMaxRetry(count)
	return tx
}

// SetMaxBackoff The maximum amount of time to wait between retries.
// Every retry attempt will increase the wait time exponentially until it reaches this time.
func (tx *TokenCancelAirdropTransaction) SetMaxBackoff(max time.Duration) *TokenCancelAirdropTransaction {
	tx.Transaction.SetMaxBackoff(max)
	return tx
}

// SetMinBackoff sets the minimum amount of time to wait between retries.
func (tx *TokenCancelAirdropTransaction) SetMinBackoff(min time.Duration) *TokenCancelAirdropTransaction {
	tx.Transaction.SetMinBackoff(min)
	return tx
}

func (tx *TokenCancelAirdropTransaction) SetLogLevel(level LogLevel) *TokenCancelAirdropTransaction {
	tx.Transaction.SetLogLevel(level)
	return tx
}

func (tx *TokenCancelAirdropTransaction) Execute(client *Client) (TransactionResponse, error) {
	return tx.Transaction.execute(client, tx)
}

func (tx *TokenCancelAirdropTransaction) Schedule() (*ScheduleCreateTransaction, error) {
	return tx.Transaction.schedule(tx)
}

// ----------- Overridden functions ----------------

func (tx *TokenCancelAirdropTransaction) getName() string {
	return "TokenCancelAirdropTransaction"
}

func (tx *TokenCancelAirdropTransaction) validateNetworkOnIDs(client *Client) error {
	if client == nil || !client.autoValidateChecksums {
		return nil
	}

	for _, pendingAirdropId := range tx.pendingAirdropIds {
		if pendingAirdropId.sender != nil {
			if err := pendingAirdropId.sender.ValidateChecksum(client); err != nil {
				return err
			}
		}

		if pendingAirdropId.receiver != nil {
			if err := pendingAirdropId.receiver.ValidateChecksum(client); err != nil {
				return err
			}
		}

		if pendingAirdropId.nftID != nil {
			if err := pendingAirdropId.nftID.Validate(client); err != nil {
				return err
			}
		}

		if pendingAirdropId.tokenID != nil {
			if err := pendingAirdropId.tokenID.ValidateChecksum(client); err != nil {
				return err
			}
		}
	}
	return nil
}

func (tx *TokenCancelAirdropTransaction) build() *services.TransactionBody {
	return &services.TransactionBody{
		TransactionFee:           tx.transactionFee,
		Memo:                     tx.Transaction.memo,
		TransactionValidDuration: _DurationToProtobuf(tx.GetTransactionValidDuration()),
		TransactionID:            tx.transactionID._ToProtobuf(),
		Data: &services.TransactionBody_TokenCancelAirdrop{
			TokenCancelAirdrop: tx.buildProtoBody(),
		},
	}
}

func (tx *TokenCancelAirdropTransaction) buildScheduled() (*services.SchedulableTransactionBody, error) {
	return &services.SchedulableTransactionBody{
		TransactionFee: tx.transactionFee,
		Data: &services.SchedulableTransactionBody_TokenCancelAirdrop{
			TokenCancelAirdrop: tx.buildProtoBody(),
		},
	}, nil
}

func (tx *TokenCancelAirdropTransaction) buildProtoBody() *services.TokenCancelAirdropTransactionBody {
	pendingAirdrops := make([]*services.PendingAirdropId, len(tx.pendingAirdropIds))
	for i, pendingAirdropId := range tx.pendingAirdropIds {
		pendingAirdrops[i] = pendingAirdropId._ToProtobuf()
	}

	return &services.TokenCancelAirdropTransactionBody{
		PendingAirdrops: pendingAirdrops,
	}
}

func (tx *TokenCancelAirdropTransaction) getMethod(channel *_Channel) _Method {
	return _Method{
		transaction: channel._GetToken().CancelAirdrop,
	}
}

func (tx *TokenCancelAirdropTransaction) _ConstructScheduleProtobuf() (*services.SchedulableTransactionBody, error) {
	return tx.buildScheduled()
}
