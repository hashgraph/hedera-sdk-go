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

	"github.com/pkg/errors"

	"github.com/hashgraph/hedera-protobufs-go/services"
)

// TokenFeeScheduleUpdateTransaction
// At consensus, updates a token type's fee schedule to the given list of custom fees.
//
// If the target token type has no fee_schedule_key, resolves to TOKEN_HAS_NO_FEE_SCHEDULE_KEY.
// Otherwise this transaction must be signed to the fee_schedule_key, or the transaction will
// resolve to INVALID_SIGNATURE.
//
// If the custom_fees list is empty, clears the fee schedule or resolves to
// CUSTOM_SCHEDULE_ALREADY_HAS_NO_FEES if the fee schedule was already empty.
type TokenFeeScheduleUpdateTransaction struct {
	Transaction
	tokenID    *TokenID
	customFees []Fee
}

// NewTokenFeeScheduleUpdateTransaction creates TokenFeeScheduleUpdateTransaction which
// at consensus, updates a token type's fee schedule to the given list of custom fees.
//
// If the target token type has no fee_schedule_key, resolves to TOKEN_HAS_NO_FEE_SCHEDULE_KEY.
// Otherwise this transaction must be signed to the fee_schedule_key, or the transaction will
// resolve to INVALID_SIGNATURE.
//
// If the custom_fees list is empty, clears the fee schedule or resolves to
// CUSTOM_SCHEDULE_ALREADY_HAS_NO_FEES if the fee schedule was already empty.
func NewTokenFeeScheduleUpdateTransaction() *TokenFeeScheduleUpdateTransaction {
	tx := TokenFeeScheduleUpdateTransaction{
		Transaction: _NewTransaction(),
	}

	tx._SetDefaultMaxTransactionFee(NewHbar(5))

	return &tx
}

func _TokenFeeScheduleUpdateTransactionFromProtobuf(transaction Transaction, pb *services.TransactionBody) *TokenFeeScheduleUpdateTransaction {
	customFees := make([]Fee, 0)

	for _, fee := range pb.GetTokenFeeScheduleUpdate().GetCustomFees() {
		customFees = append(customFees, _CustomFeeFromProtobuf(fee))
	}

	resultTx := &TokenFeeScheduleUpdateTransaction{
		Transaction: transaction,
		tokenID:     _TokenIDFromProtobuf(pb.GetTokenFeeScheduleUpdate().TokenId),
		customFees:  customFees,
	}
	return resultTx
}

// SetTokenID Sets the token whose fee schedule is to be updated
func (tx *TokenFeeScheduleUpdateTransaction) SetTokenID(tokenID TokenID) *TokenFeeScheduleUpdateTransaction {
	tx._RequireNotFrozen()
	tx.tokenID = &tokenID
	return tx
}

// GetTokenID returns the token whose fee schedule is to be updated
func (tx *TokenFeeScheduleUpdateTransaction) GetTokenID() TokenID {
	if tx.tokenID == nil {
		return TokenID{}
	}

	return *tx.tokenID
}

// SetCustomFees Sets the new custom fees to be assessed during a CryptoTransfer that transfers units of this token
func (tx *TokenFeeScheduleUpdateTransaction) SetCustomFees(fees []Fee) *TokenFeeScheduleUpdateTransaction {
	tx._RequireNotFrozen()
	tx.customFees = fees
	return tx
}

// GetCustomFees returns the new custom fees to be assessed during a CryptoTransfer that transfers units of this token
func (tx *TokenFeeScheduleUpdateTransaction) GetCustomFees() []Fee {
	return tx.customFees
}

// ---- Required Interfaces ---- //

// Sign uses the provided privateKey to sign the transaction.
func (tx *TokenFeeScheduleUpdateTransaction) Sign(privateKey PrivateKey) *TokenFeeScheduleUpdateTransaction {
	tx.Transaction.Sign(privateKey)
	return tx
}

// SignWithOperator signs the transaction with client's operator privateKey.
func (tx *TokenFeeScheduleUpdateTransaction) SignWithOperator(client *Client) (*TokenFeeScheduleUpdateTransaction, error) {
	_, err := tx.Transaction.signWithOperator(client, tx)
	if err != nil {
		return nil, err
	}
	return tx, nil
}

// SignWith executes the TransactionSigner and adds the resulting signature data to the transaction's signature map
// with the publicKey as the map key.
func (tx *TokenFeeScheduleUpdateTransaction) SignWith(
	publicKey PublicKey,
	signer TransactionSigner,
) *TokenFeeScheduleUpdateTransaction {
	tx.Transaction.SignWith(publicKey, signer)
	return tx
}

// AddSignature adds a signature to the transaction.
func (tx *TokenFeeScheduleUpdateTransaction) AddSignature(publicKey PublicKey, signature []byte) *TokenFeeScheduleUpdateTransaction {
	tx.Transaction.AddSignature(publicKey, signature)
	return tx
}

// When execution is attempted, a single attempt will timeout when this deadline is reached. (The SDK may subsequently retry the execution.)
func (tx *TokenFeeScheduleUpdateTransaction) SetGrpcDeadline(deadline *time.Duration) *TokenFeeScheduleUpdateTransaction {
	tx.Transaction.SetGrpcDeadline(deadline)
	return tx
}

func (tx *TokenFeeScheduleUpdateTransaction) Freeze() (*TokenFeeScheduleUpdateTransaction, error) {
	return tx.FreezeWith(nil)
}

func (tx *TokenFeeScheduleUpdateTransaction) FreezeWith(client *Client) (*TokenFeeScheduleUpdateTransaction, error) {
	_, err := tx.Transaction.freezeWith(client, tx)
	return tx, err
}

// SetMaxTransactionFee sets the max transaction fee for this TokenFeeScheduleUpdateTransaction.
func (tx *TokenFeeScheduleUpdateTransaction) SetMaxTransactionFee(fee Hbar) *TokenFeeScheduleUpdateTransaction {
	tx.Transaction.SetMaxTransactionFee(fee)
	return tx
}

// SetRegenerateTransactionID sets if transaction IDs should be regenerated when `TRANSACTION_EXPIRED` is received
func (tx *TokenFeeScheduleUpdateTransaction) SetRegenerateTransactionID(regenerateTransactionID bool) *TokenFeeScheduleUpdateTransaction {
	tx.Transaction.SetRegenerateTransactionID(regenerateTransactionID)
	return tx
}

// SetTransactionMemo sets the memo for this TokenFeeScheduleUpdateTransaction.
func (tx *TokenFeeScheduleUpdateTransaction) SetTransactionMemo(memo string) *TokenFeeScheduleUpdateTransaction {
	tx.Transaction.SetTransactionMemo(memo)
	return tx
}

// SetTransactionValidDuration sets the valid duration for this TokenFeeScheduleUpdateTransaction.
func (tx *TokenFeeScheduleUpdateTransaction) SetTransactionValidDuration(duration time.Duration) *TokenFeeScheduleUpdateTransaction {
	tx.Transaction.SetTransactionValidDuration(duration)
	return tx
}

// SetTransactionID sets the TransactionID for this TokenFeeScheduleUpdateTransaction.
func (tx *TokenFeeScheduleUpdateTransaction) SetTransactionID(transactionID TransactionID) *TokenFeeScheduleUpdateTransaction {
	tx.Transaction.SetTransactionID(transactionID)
	return tx
}

// SetNodeAccountIDs sets the _Node AccountID for this TokenFeeScheduleUpdateTransaction.
func (tx *TokenFeeScheduleUpdateTransaction) SetNodeAccountIDs(nodeID []AccountID) *TokenFeeScheduleUpdateTransaction {
	tx.Transaction.SetNodeAccountIDs(nodeID)
	return tx
}

// SetMaxRetry sets the max number of errors before execution will fail.
func (tx *TokenFeeScheduleUpdateTransaction) SetMaxRetry(count int) *TokenFeeScheduleUpdateTransaction {
	tx.Transaction.SetMaxRetry(count)
	return tx
}

// SetMaxBackoff The maximum amount of time to wait between retries.
// Every retry attempt will increase the wait time exponentially until it reaches this time.
func (tx *TokenFeeScheduleUpdateTransaction) SetMaxBackoff(max time.Duration) *TokenFeeScheduleUpdateTransaction {
	tx.Transaction.SetMaxBackoff(max)
	return tx
}

// SetMinBackoff sets the minimum amount of time to wait between retries.
func (tx *TokenFeeScheduleUpdateTransaction) SetMinBackoff(min time.Duration) *TokenFeeScheduleUpdateTransaction {
	tx.Transaction.SetMinBackoff(min)
	return tx
}

func (tx *TokenFeeScheduleUpdateTransaction) SetLogLevel(level LogLevel) *TokenFeeScheduleUpdateTransaction {
	tx.Transaction.SetLogLevel(level)
	return tx
}

func (tx *TokenFeeScheduleUpdateTransaction) Execute(client *Client) (TransactionResponse, error) {
	return tx.Transaction.execute(client, tx)
}

func (tx *TokenFeeScheduleUpdateTransaction) Schedule() (*ScheduleCreateTransaction, error) {
	return tx.Transaction.schedule(tx)
}

// ----------- Overridden functions ----------------

func (tx *TokenFeeScheduleUpdateTransaction) getName() string {
	return "TokenFeeScheduleUpdateTransaction"
}

func (tx *TokenFeeScheduleUpdateTransaction) validateNetworkOnIDs(client *Client) error {
	if client == nil || !client.autoValidateChecksums {
		return nil
	}

	if tx.tokenID != nil {
		if err := tx.tokenID.ValidateChecksum(client); err != nil {
			return err
		}
	}

	for _, customFee := range tx.customFees {
		if err := customFee.validateNetworkOnIDs(client); err != nil {
			return err
		}
	}

	return nil
}

func (tx *TokenFeeScheduleUpdateTransaction) build() *services.TransactionBody {
	return &services.TransactionBody{
		TransactionFee:           tx.transactionFee,
		Memo:                     tx.Transaction.memo,
		TransactionValidDuration: _DurationToProtobuf(tx.GetTransactionValidDuration()),
		TransactionID:            tx.transactionID._ToProtobuf(),
		Data: &services.TransactionBody_TokenFeeScheduleUpdate{
			TokenFeeScheduleUpdate: tx.buildProtoBody(),
		},
	}
}

func (tx *TokenFeeScheduleUpdateTransaction) buildScheduled() (*services.SchedulableTransactionBody, error) {
	return nil, errors.New("cannot schedule `TokenFeeScheduleUpdateTransaction")
}

func (tx *TokenFeeScheduleUpdateTransaction) buildProtoBody() *services.TokenFeeScheduleUpdateTransactionBody {
	body := &services.TokenFeeScheduleUpdateTransactionBody{}
	if tx.tokenID != nil {
		body.TokenId = tx.tokenID._ToProtobuf()
	}

	if len(tx.customFees) > 0 {
		for _, customFee := range tx.customFees {
			if body.CustomFees == nil {
				body.CustomFees = make([]*services.CustomFee, 0)
			}
			body.CustomFees = append(body.CustomFees, customFee._ToProtobuf())
		}
	}

	return body
}

func (tx *TokenFeeScheduleUpdateTransaction) getMethod(channel *_Channel) _Method {
	return _Method{
		transaction: channel._GetToken().UpdateTokenFeeSchedule,
	}
}
func (tx *TokenFeeScheduleUpdateTransaction) _ConstructScheduleProtobuf() (*services.SchedulableTransactionBody, error) {
	return tx.buildScheduled()
}
