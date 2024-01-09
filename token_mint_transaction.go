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

	"github.com/hashgraph/hedera-protobufs-go/services"
)

// TokenMintTransaction
// Mints tokens from the Token's treasury Account. If no Supply Key is defined, the transaction
// will resolve to TOKEN_HAS_NO_SUPPLY_KEY.
// The operation decreases the Total Supply of the Token. Total supply cannot go below zero.
// The amount provided must be in the lowest denomination possible. Example:
// Token A has 2 decimals. In order to mint 100 tokens, one must provide amount of 10000. In order
// to mint 100.55 tokens, one must provide amount of 10055.
type TokenMintTransaction struct {
	Transaction
	tokenID *TokenID
	amount  uint64
	meta    [][]byte
}

// NewTokenMintTransaction creates TokenMintTransaction which
// mints tokens from the Token's treasury Account. If no Supply Key is defined, the transaction
// will resolve to TOKEN_HAS_NO_SUPPLY_KEY.
// The operation decreases the Total Supply of the Token. Total supply cannot go below zero.
// The amount provided must be in the lowest denomination possible. Example:
// Token A has 2 decimals. In order to mint 100 tokens, one must provide amount of 10000. In order
// to mint 100.55 tokens, one must provide amount of 10055.
func NewTokenMintTransaction() *TokenMintTransaction {
	tx := TokenMintTransaction{
		Transaction: _NewTransaction(),
	}

	tx._SetDefaultMaxTransactionFee(NewHbar(30))

	return &tx
}

func _TokenMintTransactionFromProtobuf(tx Transaction, pb *services.TransactionBody) *TokenMintTransaction {
	return &TokenMintTransaction{
		Transaction: tx,
		tokenID:     _TokenIDFromProtobuf(pb.GetTokenMint().GetToken()),
		amount:      pb.GetTokenMint().GetAmount(),
		meta:        pb.GetTokenMint().GetMetadata(),
	}
}

// SetTokenID Sets the token for which to mint tokens. If token does not exist, transaction results in
// INVALID_TOKEN_ID
func (tx *TokenMintTransaction) SetTokenID(tokenID TokenID) *TokenMintTransaction {
	tx._RequireNotFrozen()
	tx.tokenID = &tokenID
	return tx
}

// GetTokenID returns the TokenID for this TokenMintTransaction
func (tx *TokenMintTransaction) GetTokenID() TokenID {
	if tx.tokenID == nil {
		return TokenID{}
	}

	return *tx.tokenID
}

// SetAmount Sets the amount to mint from the Treasury Account. Amount must be a positive non-zero number, not
// bigger than the token balance of the treasury account (0; balance], represented in the lowest
// denomination.
func (tx *TokenMintTransaction) SetAmount(amount uint64) *TokenMintTransaction {
	tx._RequireNotFrozen()
	tx.amount = amount
	return tx
}

// GetAmount returns the amount to mint from the Treasury Account
func (tx *TokenMintTransaction) GetAmount() uint64 {
	return tx.amount
}

// SetMetadatas
// Applicable to tokens of type NON_FUNGIBLE_UNIQUE. A list of metadata that are being created.
// Maximum allowed size of each metadata is 100 bytes
func (tx *TokenMintTransaction) SetMetadatas(meta [][]byte) *TokenMintTransaction {
	tx._RequireNotFrozen()
	tx.meta = meta
	return tx
}

// SetMetadata
// Applicable to tokens of type NON_FUNGIBLE_UNIQUE. A list of metadata that are being created.
// Maximum allowed size of each metadata is 100 bytes
func (tx *TokenMintTransaction) SetMetadata(meta []byte) *TokenMintTransaction {
	tx._RequireNotFrozen()
	if tx.meta == nil {
		tx.meta = make([][]byte, 0)
	}
	tx.meta = append(tx.meta, [][]byte{meta}...)
	return tx
}

// GetMetadatas returns the metadata that are being created.
func (tx *TokenMintTransaction) GetMetadatas() [][]byte {
	return tx.meta
}

// ---- Required Interfaces ---- //

// Sign uses the provided privateKey to sign the transaction.
func (tx *TokenMintTransaction) Sign(privateKey PrivateKey) *TokenMintTransaction {
	tx.Transaction.Sign(privateKey)
	return tx
}

// SignWithOperator signs the transaction with client's operator privateKey.
func (tx *TokenMintTransaction) SignWithOperator(client *Client) (*TokenMintTransaction, error) {
	_, err := tx.Transaction.signWithOperator(client, tx)
	if err != nil {
		return nil, err
	}
	return tx, nil
}

// SignWith executes the TransactionSigner and adds the resulting signature data to the Transaction's signature map
// with the publicKey as the map key.
func (tx *TokenMintTransaction) SignWith(
	publicKey PublicKey,
	signer TransactionSigner,
) *TokenMintTransaction {
	tx.Transaction.SignWith(publicKey, signer)
	return tx
}

// AddSignature adds a signature to the transaction.
func (tx *TokenMintTransaction) AddSignature(publicKey PublicKey, signature []byte) *TokenMintTransaction {
	tx.Transaction.AddSignature(publicKey, signature)
	return tx
}

// When execution is attempted, a single attempt will timeout when this deadline is reached. (The SDK may subsequently retry the execution.)
func (tx *TokenMintTransaction) SetGrpcDeadline(deadline *time.Duration) *TokenMintTransaction {
	tx.Transaction.SetGrpcDeadline(deadline)
	return tx
}

func (tx *TokenMintTransaction) Freeze() (*TokenMintTransaction, error) {
	return tx.FreezeWith(nil)
}

func (tx *TokenMintTransaction) FreezeWith(client *Client) (*TokenMintTransaction, error) {
	_, err := tx.Transaction.freezeWith(client, tx)
	return tx, err
}

// SetMaxTransactionFee sets the max transaction fee for this TokenMintTransaction.
func (tx *TokenMintTransaction) SetMaxTransactionFee(fee Hbar) *TokenMintTransaction {
	tx.Transaction.SetMaxTransactionFee(fee)
	return tx
}

// SetRegenerateTransactionID sets if transaction IDs should be regenerated when `TRANSACTION_EXPIRED` is received
func (tx *TokenMintTransaction) SetRegenerateTransactionID(regenerateTransactionID bool) *TokenMintTransaction {
	tx.Transaction.SetRegenerateTransactionID(regenerateTransactionID)
	return tx
}

// SetTransactionMemo sets the memo for this TokenMintTransaction.
func (tx *TokenMintTransaction) SetTransactionMemo(memo string) *TokenMintTransaction {
	tx.Transaction.SetTransactionMemo(memo)
	return tx
}

// SetTransactionValidDuration sets the valid duration for this TokenMintTransaction.
func (tx *TokenMintTransaction) SetTransactionValidDuration(duration time.Duration) *TokenMintTransaction {
	tx.Transaction.SetTransactionValidDuration(duration)
	return tx
}

// SetTransactionID sets the TransactionID for this TokenMintTransaction.
func (tx *TokenMintTransaction) SetTransactionID(transactionID TransactionID) *TokenMintTransaction {
	tx.Transaction.SetTransactionID(transactionID)
	return tx
}

// SetNodeAccountIDs sets the _Node AccountID for this TokenMintTransaction.
func (tx *TokenMintTransaction) SetNodeAccountIDs(nodeID []AccountID) *TokenMintTransaction {
	tx.Transaction.SetNodeAccountIDs(nodeID)
	return tx
}

// SetMaxRetry sets the max number of errors before execution will fail.
func (tx *TokenMintTransaction) SetMaxRetry(count int) *TokenMintTransaction {
	tx.Transaction.SetMaxRetry(count)
	return tx
}

// SetMaxBackoff The maximum amount of time to wait between retries.
// Every retry attempt will increase the wait time exponentially until it reaches this time.
func (tx *TokenMintTransaction) SetMaxBackoff(max time.Duration) *TokenMintTransaction {
	tx.Transaction.SetMaxBackoff(max)
	return tx
}

// SetMinBackoff sets the minimum amount of time to wait between retries.
func (tx *TokenMintTransaction) SetMinBackoff(min time.Duration) *TokenMintTransaction {
	tx.Transaction.SetMinBackoff(min)
	return tx
}

func (tx *TokenMintTransaction) SetLogLevel(level LogLevel) *TokenMintTransaction {
	tx.Transaction.SetLogLevel(level)
	return tx
}

func (tx *TokenMintTransaction) Execute(client *Client) (TransactionResponse, error) {
	return tx.Transaction.execute(client, tx)
}

func (tx *TokenMintTransaction) Schedule() (*ScheduleCreateTransaction, error) {
	return tx.Transaction.schedule(tx)
}

// ----------- Overridden functions ----------------

func (tx *TokenMintTransaction) getName() string {
	return "TokenMintTransaction"
}

func (tx *TokenMintTransaction) validateNetworkOnIDs(client *Client) error {
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

func (tx *TokenMintTransaction) build() *services.TransactionBody {
	return &services.TransactionBody{
		TransactionFee:           tx.transactionFee,
		Memo:                     tx.Transaction.memo,
		TransactionValidDuration: _DurationToProtobuf(tx.GetTransactionValidDuration()),
		TransactionID:            tx.transactionID._ToProtobuf(),
		Data: &services.TransactionBody_TokenMint{
			TokenMint: tx.buildProtoBody(),
		},
	}
}

func (tx *TokenMintTransaction) buildScheduled() (*services.SchedulableTransactionBody, error) {
	return &services.SchedulableTransactionBody{
		TransactionFee: tx.transactionFee,
		Memo:           tx.Transaction.memo,
		Data: &services.SchedulableTransactionBody_TokenMint{
			TokenMint: tx.buildProtoBody(),
		},
	}, nil
}

func (tx *TokenMintTransaction) buildProtoBody() *services.TokenMintTransactionBody {
	body := &services.TokenMintTransactionBody{
		Amount: tx.amount,
	}

	if tx.meta != nil {
		body.Metadata = tx.meta
	}

	if tx.tokenID != nil {
		body.Token = tx.tokenID._ToProtobuf()
	}

	return body
}

func (tx *TokenMintTransaction) getMethod(channel *_Channel) _Method {
	return _Method{
		transaction: channel._GetToken().MintToken,
	}
}
func (tx *TokenMintTransaction) _ConstructScheduleProtobuf() (*services.SchedulableTransactionBody, error) {
	return tx.buildScheduled()
}
