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
	"fmt"
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
	transaction := TokenFeeScheduleUpdateTransaction{
		Transaction: _NewTransaction(),
	}
	transaction._SetDefaultMaxTransactionFee(NewHbar(5))

	return &transaction
}

func _TokenFeeScheduleUpdateTransactionFromProtobuf(transaction Transaction, pb *services.TransactionBody) *TokenFeeScheduleUpdateTransaction {
	customFees := make([]Fee, 0)

	for _, fee := range pb.GetTokenFeeScheduleUpdate().GetCustomFees() {
		customFees = append(customFees, _CustomFeeFromProtobuf(fee))
	}

	return &TokenFeeScheduleUpdateTransaction{
		Transaction: transaction,
		tokenID:     _TokenIDFromProtobuf(pb.GetTokenFeeScheduleUpdate().TokenId),
		customFees:  customFees,
	}
}

// When execution is attempted, a single attempt will timeout when this deadline is reached. (The SDK may subsequently retry the execution.)
func (transaction *TokenFeeScheduleUpdateTransaction) SetGrpcDeadline(deadline *time.Duration) *TokenFeeScheduleUpdateTransaction {
	transaction.Transaction.SetGrpcDeadline(deadline)
	return transaction
}

// SetTokenID Sets the token whose fee schedule is to be updated
func (transaction *TokenFeeScheduleUpdateTransaction) SetTokenID(tokenID TokenID) *TokenFeeScheduleUpdateTransaction {
	transaction._RequireNotFrozen()
	transaction.tokenID = &tokenID
	return transaction
}

// GetTokenID returns the token whose fee schedule is to be updated
func (transaction *TokenFeeScheduleUpdateTransaction) GetTokenID() TokenID {
	if transaction.tokenID == nil {
		return TokenID{}
	}

	return *transaction.tokenID
}

// SetCustomFees Sets the new custom fees to be assessed during a CryptoTransfer that transfers units of this token
func (transaction *TokenFeeScheduleUpdateTransaction) SetCustomFees(fees []Fee) *TokenFeeScheduleUpdateTransaction {
	transaction._RequireNotFrozen()
	transaction.customFees = fees
	return transaction
}

// GetCustomFees returns the new custom fees to be assessed during a CryptoTransfer that transfers units of this token
func (transaction *TokenFeeScheduleUpdateTransaction) GetCustomFees() []Fee {
	return transaction.customFees
}

func (transaction *TokenFeeScheduleUpdateTransaction) _ValidateNetworkOnIDs(client *Client) error {
	if client == nil || !client.autoValidateChecksums {
		return nil
	}

	if transaction.tokenID != nil {
		if err := transaction.tokenID.ValidateChecksum(client); err != nil {
			return err
		}
	}

	for _, customFee := range transaction.customFees {
		if err := customFee._ValidateNetworkOnIDs(client); err != nil {
			return err
		}
	}

	return nil
}

func (transaction *TokenFeeScheduleUpdateTransaction) _Build() *services.TransactionBody {
	body := &services.TokenFeeScheduleUpdateTransactionBody{}
	if transaction.tokenID != nil {
		body.TokenId = transaction.tokenID._ToProtobuf()
	}

	if len(transaction.customFees) > 0 {
		for _, customFee := range transaction.customFees {
			if body.CustomFees == nil {
				body.CustomFees = make([]*services.CustomFee, 0)
			}
			body.CustomFees = append(body.CustomFees, customFee._ToProtobuf())
		}
	}

	return &services.TransactionBody{
		TransactionFee:           transaction.transactionFee,
		Memo:                     transaction.Transaction.memo,
		TransactionValidDuration: _DurationToProtobuf(transaction.GetTransactionValidDuration()),
		TransactionID:            transaction.transactionID._ToProtobuf(),
		Data: &services.TransactionBody_TokenFeeScheduleUpdate{
			TokenFeeScheduleUpdate: body,
		},
	}
}

func (transaction *TokenFeeScheduleUpdateTransaction) _ConstructScheduleProtobuf() (*services.SchedulableTransactionBody, error) {
	return nil, errors.New("cannot schedule `TokenFeeScheduleUpdateTransaction")
}

func _TokenFeeScheduleUpdateTransactionGetMethod(request interface{}, channel *_Channel) _Method {
	return _Method{
		transaction: channel._GetToken().UpdateTokenFeeSchedule,
	}
}

func (transaction *TokenFeeScheduleUpdateTransaction) IsFrozen() bool {
	return transaction._IsFrozen()
}

// Sign uses the provided privateKey to sign the transaction.
func (transaction *TokenFeeScheduleUpdateTransaction) Sign(
	privateKey PrivateKey,
) *TokenFeeScheduleUpdateTransaction {
	return transaction.SignWith(privateKey.PublicKey(), privateKey.Sign)
}

// SignWithOperator signs the transaction with client's operator privateKey.
func (transaction *TokenFeeScheduleUpdateTransaction) SignWithOperator(
	client *Client,
) (*TokenFeeScheduleUpdateTransaction, error) {
	// If the transaction is not signed by the _Operator, we need
	// to sign the transaction with the _Operator

	if client == nil {
		return nil, errNoClientProvided
	} else if client.operator == nil {
		return nil, errClientOperatorSigning
	}

	if !transaction.IsFrozen() {
		_, err := transaction.FreezeWith(client)
		if err != nil {
			return transaction, err
		}
	}
	return transaction.SignWith(client.operator.publicKey, client.operator.signer), nil
}

// SignWith executes the TransactionSigner and adds the resulting signature data to the Transaction's signature map
// with the publicKey as the map key.
func (transaction *TokenFeeScheduleUpdateTransaction) SignWith(
	publicKey PublicKey,
	signer TransactionSigner,
) *TokenFeeScheduleUpdateTransaction {
	if !transaction._KeyAlreadySigned(publicKey) {
		transaction._SignWith(publicKey, signer)
	}

	return transaction
}

// Execute executes the Transaction with the provided client
func (transaction *TokenFeeScheduleUpdateTransaction) Execute(
	client *Client,
) (TransactionResponse, error) {
	if client == nil {
		return TransactionResponse{}, errNoClientProvided
	}

	if transaction.freezeError != nil {
		return TransactionResponse{}, transaction.freezeError
	}

	if !transaction.IsFrozen() {
		_, err := transaction.FreezeWith(client)
		if err != nil {
			return TransactionResponse{}, err
		}
	}

	transactionID := transaction.transactionIDs._GetCurrent().(TransactionID)

	if !client.GetOperatorAccountID()._IsZero() && client.GetOperatorAccountID()._Equals(*transactionID.AccountID) {
		transaction.SignWith(
			client.GetOperatorPublicKey(),
			client.operator.signer,
		)
	}

	resp, err := _Execute(
		client,
		&transaction.Transaction,
		_TransactionShouldRetry,
		_TransactionMakeRequest,
		_TransactionAdvanceRequest,
		_TransactionGetNodeAccountID,
		_TokenFeeScheduleUpdateTransactionGetMethod,
		_TransactionMapStatusError,
		_TransactionMapResponse,
		transaction._GetLogID(),
		transaction.grpcDeadline,
		transaction.maxBackoff,
		transaction.minBackoff,
		transaction.maxRetry,
	)

	if err != nil {
		return TransactionResponse{
			TransactionID:  transaction.GetTransactionID(),
			NodeID:         resp.(TransactionResponse).NodeID,
			ValidateStatus: true,
		}, err
	}

	return TransactionResponse{
		TransactionID:  transaction.GetTransactionID(),
		NodeID:         resp.(TransactionResponse).NodeID,
		Hash:           resp.(TransactionResponse).Hash,
		ValidateStatus: true,
	}, nil
}

func (transaction *TokenFeeScheduleUpdateTransaction) Freeze() (*TokenFeeScheduleUpdateTransaction, error) {
	return transaction.FreezeWith(nil)
}

func (transaction *TokenFeeScheduleUpdateTransaction) FreezeWith(client *Client) (*TokenFeeScheduleUpdateTransaction, error) {
	if transaction.IsFrozen() {
		return transaction, nil
	}
	transaction._InitFee(client)
	err := transaction._ValidateNetworkOnIDs(client)
	if err != nil {
		return &TokenFeeScheduleUpdateTransaction{}, err
	}
	if err := transaction._InitTransactionID(client); err != nil {
		return transaction, err
	}
	body := transaction._Build()

	return transaction, _TransactionFreezeWith(&transaction.Transaction, client, body)
}

// GetMaxTransactionFee returns the maximum transaction fee the operator (paying account) is willing to pay.
func (transaction *TokenFeeScheduleUpdateTransaction) GetMaxTransactionFee() Hbar {
	return transaction.Transaction.GetMaxTransactionFee()
}

// SetMaxTransactionFee sets the maximum transaction fee the operator (paying account) is willing to pay.
func (transaction *TokenFeeScheduleUpdateTransaction) SetMaxTransactionFee(fee Hbar) *TokenFeeScheduleUpdateTransaction {
	transaction._RequireNotFrozen()
	transaction.Transaction.SetMaxTransactionFee(fee)
	return transaction
}

// SetRegenerateTransactionID sets if transaction IDs should be regenerated when `TRANSACTION_EXPIRED` is received
func (transaction *TokenFeeScheduleUpdateTransaction) SetRegenerateTransactionID(regenerateTransactionID bool) *TokenFeeScheduleUpdateTransaction {
	transaction._RequireNotFrozen()
	transaction.Transaction.SetRegenerateTransactionID(regenerateTransactionID)
	return transaction
}

// GetRegenerateTransactionID returns true if transaction ID regeneration is enabled.
func (transaction *TokenFeeScheduleUpdateTransaction) GetRegenerateTransactionID() bool {
	return transaction.Transaction.GetRegenerateTransactionID()
}

// GetTransactionMemo returns the memo for this	TokenFeeScheduleUpdateTransaction.
func (transaction *TokenFeeScheduleUpdateTransaction) GetTransactionMemo() string {
	return transaction.Transaction.GetTransactionMemo()
}

// SetTransactionMemo sets the memo for this TokenFeeScheduleUpdateTransaction.
func (transaction *TokenFeeScheduleUpdateTransaction) SetTransactionMemo(memo string) *TokenFeeScheduleUpdateTransaction {
	transaction._RequireNotFrozen()
	transaction.Transaction.SetTransactionMemo(memo)
	return transaction
}

// GetTransactionValidDuration returns the duration that this transaction is valid for.
func (transaction *TokenFeeScheduleUpdateTransaction) GetTransactionValidDuration() time.Duration {
	return transaction.Transaction.GetTransactionValidDuration()
}

// SetTransactionValidDuration sets the valid duration for this TokenFeeScheduleUpdateTransaction.
func (transaction *TokenFeeScheduleUpdateTransaction) SetTransactionValidDuration(duration time.Duration) *TokenFeeScheduleUpdateTransaction {
	transaction._RequireNotFrozen()
	transaction.Transaction.SetTransactionValidDuration(duration)
	return transaction
}

// GetTransactionID returns the TransactionID for this TokenFeeScheduleUpdateTransaction.
func (transaction *TokenFeeScheduleUpdateTransaction) GetTransactionID() TransactionID {
	return transaction.Transaction.GetTransactionID()
}

// SetTransactionID sets the TransactionID for this TokenFeeScheduleUpdateTransaction.
func (transaction *TokenFeeScheduleUpdateTransaction) SetTransactionID(transactionID TransactionID) *TokenFeeScheduleUpdateTransaction {
	transaction._RequireNotFrozen()

	transaction.Transaction.SetTransactionID(transactionID)
	return transaction
}

// SetNodeTokenID sets the _Node TokenID for this TokenFeeScheduleUpdateTransaction.
func (transaction *TokenFeeScheduleUpdateTransaction) SetNodeAccountIDs(nodeID []AccountID) *TokenFeeScheduleUpdateTransaction {
	transaction._RequireNotFrozen()
	transaction.Transaction.SetNodeAccountIDs(nodeID)
	return transaction
}

// SetMaxRetry sets the max number of errors before execution will fail.
func (transaction *TokenFeeScheduleUpdateTransaction) SetMaxRetry(count int) *TokenFeeScheduleUpdateTransaction {
	transaction.Transaction.SetMaxRetry(count)
	return transaction
}

// AddSignature adds a signature to the Transaction.
func (transaction *TokenFeeScheduleUpdateTransaction) AddSignature(publicKey PublicKey, signature []byte) *TokenFeeScheduleUpdateTransaction {
	transaction._RequireOneNodeAccountID()

	if transaction._KeyAlreadySigned(publicKey) {
		return transaction
	}

	if transaction.signedTransactions._Length() == 0 {
		return transaction
	}

	transaction.transactions = _NewLockableSlice()
	transaction.publicKeys = append(transaction.publicKeys, publicKey)
	transaction.transactionSigners = append(transaction.transactionSigners, nil)
	transaction.transactionIDs.locked = true

	for index := 0; index < transaction.signedTransactions._Length(); index++ {
		var temp *services.SignedTransaction
		switch t := transaction.signedTransactions._Get(index).(type) { //nolint
		case *services.SignedTransaction:
			temp = t
		}
		temp.SigMap.SigPair = append(
			temp.SigMap.SigPair,
			publicKey._ToSignaturePairProtobuf(signature),
		)
		transaction.signedTransactions._Set(index, temp)
	}

	return transaction
}

// SetMaxBackoff The maximum amount of time to wait between retries.
// Every retry attempt will increase the wait time exponentially until it reaches this time.
func (transaction *TokenFeeScheduleUpdateTransaction) SetMaxBackoff(max time.Duration) *TokenFeeScheduleUpdateTransaction {
	if max.Nanoseconds() < 0 {
		panic("maxBackoff must be a positive duration")
	} else if max.Nanoseconds() < transaction.minBackoff.Nanoseconds() {
		panic("maxBackoff must be greater than or equal to minBackoff")
	}
	transaction.maxBackoff = &max
	return transaction
}

// GetMaxBackoff returns the maximum amount of time to wait between retries.
func (transaction *TokenFeeScheduleUpdateTransaction) GetMaxBackoff() time.Duration {
	if transaction.maxBackoff != nil {
		return *transaction.maxBackoff
	}

	return 8 * time.Second
}

// SetMinBackoff sets the minimum amount of time to wait between retries.
func (transaction *TokenFeeScheduleUpdateTransaction) SetMinBackoff(min time.Duration) *TokenFeeScheduleUpdateTransaction {
	if min.Nanoseconds() < 0 {
		panic("minBackoff must be a positive duration")
	} else if transaction.maxBackoff.Nanoseconds() < min.Nanoseconds() {
		panic("minBackoff must be less than or equal to maxBackoff")
	}
	transaction.minBackoff = &min
	return transaction
}

// GetMinBackoff returns the minimum amount of time to wait between retries.
func (transaction *TokenFeeScheduleUpdateTransaction) GetMinBackoff() time.Duration {
	if transaction.minBackoff != nil {
		return *transaction.minBackoff
	}

	return 250 * time.Millisecond
}

func (transaction *TokenFeeScheduleUpdateTransaction) _GetLogID() string {
	timestamp := transaction.transactionIDs._GetCurrent().(TransactionID).ValidStart
	return fmt.Sprintf("TokenFeeScheduleUpdateTransaction:%d", timestamp.UnixNano())
}

func (transaction *TokenFeeScheduleUpdateTransaction) SetLogLevel(level LogLevel) *TokenFeeScheduleUpdateTransaction {
	transaction.Transaction.SetLogLevel(level)
	return transaction
}
