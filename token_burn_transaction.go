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

// TokenBurnTransaction Burns tokens from the Token's treasury Account.
// If no Supply Key is defined, the transaction
// will resolve to TOKEN_HAS_NO_SUPPLY_KEY.
// The operation decreases the Total Supply of the Token. Total supply cannot go below
// zero.
// The amount provided must be in the lowest denomination possible. Example:
// Token A has 2 decimals. In order to burn 100 tokens, one must provide amount of 10000. In order
// to burn 100.55 tokens, one must provide amount of 10055.
type TokenBurnTransaction struct {
	Transaction
	tokenID *TokenID
	amount  uint64
	serial  []int64
}

// NewTokenBurnTransaction creates TokenBurnTransaction which burns tokens from the Token's treasury Account.
// If no Supply Key is defined, the transaction
// will resolve to TOKEN_HAS_NO_SUPPLY_KEY.
// The operation decreases the Total Supply of the Token. Total supply cannot go below
// zero.
// The amount provided must be in the lowest denomination possible. Example:
// Token A has 2 decimals. In order to burn 100 tokens, one must provide amount of 10000. In order
// to burn 100.55 tokens, one must provide amount of 10055.
func NewTokenBurnTransaction() *TokenBurnTransaction {
	tx := TokenBurnTransaction{
		Transaction: _NewTransaction(),
	}

	tx._SetDefaultMaxTransactionFee(NewHbar(2))

	return &tx
}

func _TokenBurnTransactionFromProtobuf(tx Transaction, pb *services.TransactionBody) *TokenBurnTransaction {
	return &TokenBurnTransaction{
		Transaction: tx,
		tokenID:     _TokenIDFromProtobuf(pb.GetTokenBurn().Token),
		amount:      pb.GetTokenBurn().GetAmount(),
		serial:      pb.GetTokenBurn().GetSerialNumbers(),
	}
}

// SetTokenID Sets the token for which to burn tokens. If token does not exist, transaction results in
// INVALID_TOKEN_ID
func (tx *TokenBurnTransaction) SetTokenID(tokenID TokenID) *TokenBurnTransaction {
	tx._RequireNotFrozen()
	tx.tokenID = &tokenID
	return tx
}

// GetTokenID returns the TokenID for the token which will be burned.
func (tx *TokenBurnTransaction) GetTokenID() TokenID {
	if tx.tokenID == nil {
		return TokenID{}
	}

	return *tx.tokenID
}

// SetAmount Sets the amount to burn from the Treasury Account. Amount must be a positive non-zero number, not
// bigger than the token balance of the treasury account (0; balance], represented in the lowest
// denomination.
func (tx *TokenBurnTransaction) SetAmount(amount uint64) *TokenBurnTransaction {
	tx._RequireNotFrozen()
	tx.amount = amount
	return tx
}

// Deprecated: Use TokenBurnTransaction.GetAmount() instead.
func (tx *TokenBurnTransaction) GetAmmount() uint64 {
	return tx.amount
}

func (tx *TokenBurnTransaction) GetAmount() uint64 {
	return tx.amount
}

// SetSerialNumber
// Applicable to tokens of type NON_FUNGIBLE_UNIQUE.
// The list of serial numbers to be burned.
func (tx *TokenBurnTransaction) SetSerialNumber(serial int64) *TokenBurnTransaction {
	tx._RequireNotFrozen()
	if tx.serial == nil {
		tx.serial = make([]int64, 0)
	}
	tx.serial = append(tx.serial, serial)
	return tx
}

// SetSerialNumbers sets the list of serial numbers to be burned.
func (tx *TokenBurnTransaction) SetSerialNumbers(serial []int64) *TokenBurnTransaction {
	tx._RequireNotFrozen()
	tx.serial = serial
	return tx
}

// GetSerialNumbers returns the list of serial numbers to be burned.
func (tx *TokenBurnTransaction) GetSerialNumbers() []int64 {
	return tx.serial
}

// ---- Required Interfaces ---- //

// Sign uses the provided privateKey to sign the transaction.
func (tx *TokenBurnTransaction) Sign(privateKey PrivateKey) *TokenBurnTransaction {
	tx.Transaction.Sign(privateKey)
	return tx
}

// SignWithOperator signs the transaction with client's operator privateKey.
func (tx *TokenBurnTransaction) SignWithOperator(client *Client) (*TokenBurnTransaction, error) {
	_, err := tx.Transaction.signWithOperator(client, tx)
	if err != nil {
		return nil, err
	}
	return tx, nil
}

// SignWith executes the TransactionSigner and adds the resulting signature data to the Transaction's signature map
// with the publicKey as the map key.
func (tx *TokenBurnTransaction) SignWith(
	publicKey PublicKey,
	signer TransactionSigner,
) *TokenBurnTransaction {
	tx.Transaction.SignWith(publicKey, signer)
	return tx
}

// AddSignature adds a signature to the transaction.
func (tx *TokenBurnTransaction) AddSignature(publicKey PublicKey, signature []byte) *TokenBurnTransaction {
	tx.Transaction.AddSignature(publicKey, signature)
	return tx
}

// When execution is attempted, a single attempt will timeout when this deadline is reached. (The SDK may subsequently retry the execution.)
func (tx *TokenBurnTransaction) SetGrpcDeadline(deadline *time.Duration) *TokenBurnTransaction {
	tx.Transaction.SetGrpcDeadline(deadline)
	return tx
}

func (tx *TokenBurnTransaction) Freeze() (*TokenBurnTransaction, error) {
	return tx.FreezeWith(nil)
}

func (tx *TokenBurnTransaction) FreezeWith(client *Client) (*TokenBurnTransaction, error) {
	_, err := tx.Transaction.freezeWith(client, tx)
	return tx, err
}

// SetMaxTransactionFee sets the max transaction fee for this TokenBurnTransaction.
func (tx *TokenBurnTransaction) SetMaxTransactionFee(fee Hbar) *TokenBurnTransaction {
	tx.Transaction.SetMaxTransactionFee(fee)
	return tx
}

// SetRegenerateTransactionID sets if transaction IDs should be regenerated when `TRANSACTION_EXPIRED` is received
func (tx *TokenBurnTransaction) SetRegenerateTransactionID(regenerateTransactionID bool) *TokenBurnTransaction {
	tx.Transaction.SetRegenerateTransactionID(regenerateTransactionID)
	return tx
}

// SetTransactionMemo sets the memo for this TokenBurnTransaction.
func (tx *TokenBurnTransaction) SetTransactionMemo(memo string) *TokenBurnTransaction {
	tx.Transaction.SetTransactionMemo(memo)
	return tx
}

// SetTransactionValidDuration sets the valid duration for this TokenBurnTransaction.
func (tx *TokenBurnTransaction) SetTransactionValidDuration(duration time.Duration) *TokenBurnTransaction {
	tx.Transaction.SetTransactionValidDuration(duration)
	return tx
}

// ToBytes serialise the tx to bytes, no matter if it is signed (locked), or not
func (tx *TokenBurnTransaction) ToBytes() ([]byte, error) {
	bytes, err := tx.Transaction.toBytes(tx)
	if err != nil {
		return nil, err
	}
	return bytes, nil
}

// SetTransactionID sets the TransactionID for this TokenBurnTransaction.
func (tx *TokenBurnTransaction) SetTransactionID(transactionID TransactionID) *TokenBurnTransaction {
	tx.Transaction.SetTransactionID(transactionID)
	return tx
}

// SetNodeAccountIDs sets the _Node AccountID for this TokenBurnTransaction.
func (tx *TokenBurnTransaction) SetNodeAccountIDs(nodeID []AccountID) *TokenBurnTransaction {
	tx.Transaction.SetNodeAccountIDs(nodeID)
	return tx
}

// SetMaxRetry sets the max number of errors before execution will fail.
func (tx *TokenBurnTransaction) SetMaxRetry(count int) *TokenBurnTransaction {
	tx.Transaction.SetMaxRetry(count)
	return tx
}

// SetMaxBackoff The maximum amount of time to wait between retries.
// Every retry attempt will increase the wait time exponentially until it reaches this time.
func (tx *TokenBurnTransaction) SetMaxBackoff(max time.Duration) *TokenBurnTransaction {
	tx.Transaction.SetMaxBackoff(max)
	return tx
}

// SetMinBackoff sets the minimum amount of time to wait between retries.
func (tx *TokenBurnTransaction) SetMinBackoff(min time.Duration) *TokenBurnTransaction {
	tx.Transaction.SetMinBackoff(min)
	return tx
}

func (tx *TokenBurnTransaction) SetLogLevel(level LogLevel) *TokenBurnTransaction {
	tx.Transaction.SetLogLevel(level)
	return tx
}

func (tx *TokenBurnTransaction) Execute(client *Client) (TransactionResponse, error) {
	return tx.Transaction.execute(client, tx)
}

func (tx *TokenBurnTransaction) Schedule() (*ScheduleCreateTransaction, error) {
	return tx.Transaction.schedule(tx)
}

// ----------- Overridden functions ----------------

func (tx *TokenBurnTransaction) getName() string {
	return "TokenBurnTransaction"
}

func (tx *TokenBurnTransaction) validateNetworkOnIDs(client *Client) error {
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

func (tx *TokenBurnTransaction) build() *services.TransactionBody {
	return &services.TransactionBody{
		TransactionFee:           tx.transactionFee,
		Memo:                     tx.Transaction.memo,
		TransactionValidDuration: _DurationToProtobuf(tx.GetTransactionValidDuration()),
		TransactionID:            tx.transactionID._ToProtobuf(),
		Data: &services.TransactionBody_TokenBurn{
			TokenBurn: tx.buildProtoBody(),
		},
	}
}

func (tx *TokenBurnTransaction) buildScheduled() (*services.SchedulableTransactionBody, error) {
	return &services.SchedulableTransactionBody{
		TransactionFee: tx.transactionFee,
		Memo:           tx.Transaction.memo,
		Data: &services.SchedulableTransactionBody_TokenBurn{
			TokenBurn: tx.buildProtoBody(),
		},
	}, nil
}

func (tx *TokenBurnTransaction) buildProtoBody() *services.TokenBurnTransactionBody {
	body := &services.TokenBurnTransactionBody{
		Amount: tx.amount,
	}

	if tx.tokenID != nil {
		body.Token = tx.tokenID._ToProtobuf()
	}

	if tx.serial != nil {
		body.SerialNumbers = tx.serial
	}

	return body
}

func (tx *TokenBurnTransaction) getMethod(channel *_Channel) _Method {
	return _Method{
		transaction: channel._GetToken().BurnToken,
	}
}
func (tx *TokenBurnTransaction) _ConstructScheduleProtobuf() (*services.SchedulableTransactionBody, error) {
	return tx.buildScheduled()
}
