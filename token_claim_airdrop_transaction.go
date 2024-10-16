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

type TokenClaimAirdropTransaction struct {
	Transaction
	pendingAirdropIds []*PendingAirdropId
}

func NewTokenClaimAirdropTransaction() *TokenClaimAirdropTransaction {
	tx := TokenClaimAirdropTransaction{
		Transaction:       _NewTransaction(),
		pendingAirdropIds: make([]*PendingAirdropId, 0),
	}

	tx._SetDefaultMaxTransactionFee(NewHbar(1))

	return &tx
}

func _TokenClaimAirdropTransactionFromProtobuf(tx Transaction, pb *services.TransactionBody) *TokenClaimAirdropTransaction {
	tokenClaim := &TokenClaimAirdropTransaction{
		Transaction: tx,
	}

	for _, pendingAirdrops := range pb.GetTokenClaimAirdrop().PendingAirdrops {
		tokenClaim.pendingAirdropIds = append(tokenClaim.pendingAirdropIds, _PendingAirdropIdFromProtobuf(pendingAirdrops))
	}

	return tokenClaim
}

// SetPendingAirdropIds sets the pending airdrop IDs for this TokenClaimAirdropTransaction.
func (tx *TokenClaimAirdropTransaction) SetPendingAirdropIds(ids []*PendingAirdropId) *TokenClaimAirdropTransaction {
	tx._RequireNotFrozen()
	tx.pendingAirdropIds = ids
	return tx
}

// AddPendingAirdropId adds a pending airdrop ID to this TokenClaimAirdropTransaction.
func (tx *TokenClaimAirdropTransaction) AddPendingAirdropId(id PendingAirdropId) *TokenClaimAirdropTransaction {
	tx._RequireNotFrozen()
	tx.pendingAirdropIds = append(tx.pendingAirdropIds, &id)
	return tx
}

// GetPendingAirdropIds returns the pending airdrop IDs for this TokenClaimAirdropTransaction.
func (tx *TokenClaimAirdropTransaction) GetPendingAirdropIds() []*PendingAirdropId {
	return tx.pendingAirdropIds
}

// ---- Required Interfaces ---- //

// Sign uses the provided privateKey to sign the transaction.
func (tx *TokenClaimAirdropTransaction) Sign(privateKey PrivateKey) *TokenClaimAirdropTransaction {
	tx.Transaction.Sign(privateKey)
	return tx
}

// SignWithOperator signs the transaction with client's operator privateKey.
func (tx *TokenClaimAirdropTransaction) SignWithOperator(client *Client) (*TokenClaimAirdropTransaction, error) {
	_, err := tx.Transaction.signWithOperator(client, tx)
	if err != nil {
		return nil, err
	}
	return tx, nil
}

// SignWith executes the TransactionSigner and adds the resulting signature data to the Transaction's signature map
// with the publicKey as the map key.
func (tx *TokenClaimAirdropTransaction) SignWith(
	publicKey PublicKey,
	signer TransactionSigner,
) *TokenClaimAirdropTransaction {
	tx.Transaction.SignWith(publicKey, signer)
	return tx
}

// AddSignature adds a signature to the transaction.
func (tx *TokenClaimAirdropTransaction) AddSignature(publicKey PublicKey, signature []byte) *TokenClaimAirdropTransaction {
	tx.Transaction.AddSignature(publicKey, signature)
	return tx
}

// When execution is attempted, a single attempt will timeout when this deadline is reached. (The SDK may subsequently retry the execution.)
func (tx *TokenClaimAirdropTransaction) SetGrpcDeadline(deadline *time.Duration) *TokenClaimAirdropTransaction {
	tx.Transaction.SetGrpcDeadline(deadline)
	return tx
}

func (tx *TokenClaimAirdropTransaction) Freeze() (*TokenClaimAirdropTransaction, error) {
	return tx.FreezeWith(nil)
}

func (tx *TokenClaimAirdropTransaction) FreezeWith(client *Client) (*TokenClaimAirdropTransaction, error) {
	_, err := tx.Transaction.freezeWith(client, tx)
	return tx, err
}

// SetMaxTransactionFee sets the max transaction fee for this TokenClaimAirdropTransaction.
func (tx *TokenClaimAirdropTransaction) SetMaxTransactionFee(fee Hbar) *TokenClaimAirdropTransaction {
	tx.Transaction.SetMaxTransactionFee(fee)
	return tx
}

// SetRegenerateTransactionID sets if transaction IDs should be regenerated when `TRANSACTION_EXPIRED` is received
func (tx *TokenClaimAirdropTransaction) SetRegenerateTransactionID(regenerateTransactionID bool) *TokenClaimAirdropTransaction {
	tx.Transaction.SetRegenerateTransactionID(regenerateTransactionID)
	return tx
}

// SetTransactionMemo sets the memo for this TokenClaimAirdropTransaction.
func (tx *TokenClaimAirdropTransaction) SetTransactionMemo(memo string) *TokenClaimAirdropTransaction {
	tx.Transaction.SetTransactionMemo(memo)
	return tx
}

// SetTransactionValidDuration sets the valid duration for this TokenClaimAirdropTransaction.
func (tx *TokenClaimAirdropTransaction) SetTransactionValidDuration(duration time.Duration) *TokenClaimAirdropTransaction {
	tx.Transaction.SetTransactionValidDuration(duration)
	return tx
}

// ToBytes serialise the tx to bytes, no matter if it is signed (locked), or not
func (tx *TokenClaimAirdropTransaction) ToBytes() ([]byte, error) {
	bytes, err := tx.Transaction.toBytes(tx)
	if err != nil {
		return nil, err
	}
	return bytes, nil
}

// SetTransactionID sets the TransactionID for this TokenClaimAirdropTransaction.
func (tx *TokenClaimAirdropTransaction) SetTransactionID(transactionID TransactionID) *TokenClaimAirdropTransaction {
	tx.Transaction.SetTransactionID(transactionID)
	return tx
}

// SetNodeAccountIDs sets the _Node AccountID for this TokenClaimAirdropTransaction.
func (tx *TokenClaimAirdropTransaction) SetNodeAccountIDs(nodeID []AccountID) *TokenClaimAirdropTransaction {
	tx.Transaction.SetNodeAccountIDs(nodeID)
	return tx
}

// SetMaxRetry sets the max number of errors before execution will fail.
func (tx *TokenClaimAirdropTransaction) SetMaxRetry(count int) *TokenClaimAirdropTransaction {
	tx.Transaction.SetMaxRetry(count)
	return tx
}

// SetMaxBackoff The maximum amount of time to wait between retries.
// Every retry attempt will increase the wait time exponentially until it reaches this time.
func (tx *TokenClaimAirdropTransaction) SetMaxBackoff(max time.Duration) *TokenClaimAirdropTransaction {
	tx.Transaction.SetMaxBackoff(max)
	return tx
}

// SetMinBackoff sets the minimum amount of time to wait between retries.
func (tx *TokenClaimAirdropTransaction) SetMinBackoff(min time.Duration) *TokenClaimAirdropTransaction {
	tx.Transaction.SetMinBackoff(min)
	return tx
}

func (tx *TokenClaimAirdropTransaction) SetLogLevel(level LogLevel) *TokenClaimAirdropTransaction {
	tx.Transaction.SetLogLevel(level)
	return tx
}

func (tx *TokenClaimAirdropTransaction) Execute(client *Client) (TransactionResponse, error) {
	return tx.Transaction.execute(client, tx)
}

func (tx *TokenClaimAirdropTransaction) Schedule() (*ScheduleCreateTransaction, error) {
	return tx.Transaction.schedule(tx)
}

// ----------- Overridden functions ----------------

func (tx *TokenClaimAirdropTransaction) getName() string {
	return "TokenClaimAirdropTransaction"
}

func (tx *TokenClaimAirdropTransaction) validateNetworkOnIDs(client *Client) error {
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

func (tx *TokenClaimAirdropTransaction) build() *services.TransactionBody {
	return &services.TransactionBody{
		TransactionFee:           tx.transactionFee,
		Memo:                     tx.Transaction.memo,
		TransactionValidDuration: _DurationToProtobuf(tx.GetTransactionValidDuration()),
		TransactionID:            tx.transactionID._ToProtobuf(),
		Data: &services.TransactionBody_TokenClaimAirdrop{
			TokenClaimAirdrop: tx.buildProtoBody(),
		},
	}
}

func (tx *TokenClaimAirdropTransaction) buildScheduled() (*services.SchedulableTransactionBody, error) {
	return &services.SchedulableTransactionBody{
		TransactionFee: tx.transactionFee,
		Data: &services.SchedulableTransactionBody_TokenClaimAirdrop{
			TokenClaimAirdrop: tx.buildProtoBody(),
		},
	}, nil
}

func (tx *TokenClaimAirdropTransaction) buildProtoBody() *services.TokenClaimAirdropTransactionBody {
	pendingAirdrops := make([]*services.PendingAirdropId, len(tx.pendingAirdropIds))
	for i, pendingAirdropId := range tx.pendingAirdropIds {
		pendingAirdrops[i] = pendingAirdropId._ToProtobuf()
	}

	return &services.TokenClaimAirdropTransactionBody{
		PendingAirdrops: pendingAirdrops,
	}
}

func (tx *TokenClaimAirdropTransaction) getMethod(channel *_Channel) _Method {
	return _Method{
		transaction: channel._GetToken().ClaimAirdrop,
	}
}

func (tx *TokenClaimAirdropTransaction) _ConstructScheduleProtobuf() (*services.SchedulableTransactionBody, error) {
	return tx.buildScheduled()
}
