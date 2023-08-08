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

	"github.com/hashgraph/hedera-protobufs-go/services"
)

// TokenAssociateTransaction Associates the provided account with the provided tokens. Must be signed by the provided Account's key.
// If the provided account is not found, the transaction will resolve to
// INVALID_ACCOUNT_ID.
// If the provided account has been deleted, the transaction will resolve to
// ACCOUNT_DELETED.
// If any of the provided tokens is not found, the transaction will resolve to
// INVALID_TOKEN_REF.
// If any of the provided tokens has been deleted, the transaction will resolve to
// TOKEN_WAS_DELETED.
// If an association between the provided account and any of the tokens already exists, the
// transaction will resolve to
// TOKEN_ALREADY_ASSOCIATED_TO_ACCOUNT.
// If the provided account's associations count exceed the constraint of maximum token
// associations per account, the transaction will resolve to
// TOKENS_PER_ACCOUNT_LIMIT_EXCEEDED.
// On success, associations between the provided account and tokens are made and the account is
// ready to interact with the tokens.
type TokenAssociateTransaction struct {
	Transaction
	accountID *AccountID
	tokens    []TokenID
}

// NewTokenAssociateTransaction creates TokenAssociateTransaction which associates the provided account with the provided tokens.
// Must be signed by the provided Account's key.
// If the provided account is not found, the transaction will resolve to
// INVALID_ACCOUNT_ID.
// If the provided account has been deleted, the transaction will resolve to
// ACCOUNT_DELETED.
// If any of the provided tokens is not found, the transaction will resolve to
// INVALID_TOKEN_REF.
// If any of the provided tokens has been deleted, the transaction will resolve to
// TOKEN_WAS_DELETED.
// If an association between the provided account and any of the tokens already exists, the
// transaction will resolve to
// TOKEN_ALREADY_ASSOCIATED_TO_ACCOUNT.
// If the provided account's associations count exceed the constraint of maximum token
// associations per account, the transaction will resolve to
// TOKENS_PER_ACCOUNT_LIMIT_EXCEEDED.
// On success, associations between the provided account and tokens are made and the account is
// ready to interact with the tokens.
func NewTokenAssociateTransaction() *TokenAssociateTransaction {
	transaction := TokenAssociateTransaction{
		Transaction: _NewTransaction(),
	}
	transaction._SetDefaultMaxTransactionFee(NewHbar(5))

	return &transaction
}

func _TokenAssociateTransactionFromProtobuf(transaction Transaction, pb *services.TransactionBody) *TokenAssociateTransaction {
	tokens := make([]TokenID, 0)
	for _, token := range pb.GetTokenAssociate().Tokens {
		if tokenID := _TokenIDFromProtobuf(token); tokenID != nil {
			tokens = append(tokens, *tokenID)
		}
	}

	return &TokenAssociateTransaction{
		Transaction: transaction,
		accountID:   _AccountIDFromProtobuf(pb.GetTokenAssociate().GetAccount()),
		tokens:      tokens,
	}
}

// When execution is attempted, a single attempt will timeout when this deadline is reached. (The SDK may subsequently retry the execution.)
func (transaction *TokenAssociateTransaction) SetGrpcDeadline(deadline *time.Duration) *TokenAssociateTransaction {
	transaction.Transaction.SetGrpcDeadline(deadline)
	return transaction
}

// SetAccountID Sets the account to be associated with the provided tokens
func (transaction *TokenAssociateTransaction) SetAccountID(accountID AccountID) *TokenAssociateTransaction {
	transaction._RequireNotFrozen()
	transaction.accountID = &accountID
	return transaction
}

// GetAccountID returns the account to be associated with the provided tokens
func (transaction *TokenAssociateTransaction) GetAccountID() AccountID {
	if transaction.accountID == nil {
		return AccountID{}
	}

	return *transaction.accountID
}

// SetTokenIDs Sets the tokens to be associated with the provided account
func (transaction *TokenAssociateTransaction) SetTokenIDs(ids ...TokenID) *TokenAssociateTransaction {
	transaction._RequireNotFrozen()
	transaction.tokens = make([]TokenID, len(ids))
	copy(transaction.tokens, ids)

	return transaction
}

// AddTokenID Adds the token to a token list to be associated with the provided account
func (transaction *TokenAssociateTransaction) AddTokenID(id TokenID) *TokenAssociateTransaction {
	transaction._RequireNotFrozen()
	if transaction.tokens == nil {
		transaction.tokens = make([]TokenID, 0)
	}

	transaction.tokens = append(transaction.tokens, id)

	return transaction
}

// GetTokenIDs returns the tokens to be associated with the provided account
func (transaction *TokenAssociateTransaction) GetTokenIDs() []TokenID {
	tokenIDs := make([]TokenID, len(transaction.tokens))
	copy(tokenIDs, transaction.tokens)

	return tokenIDs
}

func (transaction *TokenAssociateTransaction) _ValidateNetworkOnIDs(client *Client) error {
	if client == nil || !client.autoValidateChecksums {
		return nil
	}

	if transaction.accountID != nil {
		if err := transaction.accountID.ValidateChecksum(client); err != nil {
			return err
		}
	}

	for _, tokenID := range transaction.tokens {
		if err := tokenID.ValidateChecksum(client); err != nil {
			return err
		}
	}

	return nil
}

func (transaction *TokenAssociateTransaction) _Build() *services.TransactionBody {
	body := &services.TokenAssociateTransactionBody{}
	if transaction.accountID != nil {
		body.Account = transaction.accountID._ToProtobuf()
	}

	if len(transaction.tokens) > 0 {
		for _, tokenID := range transaction.tokens {
			if body.Tokens == nil {
				body.Tokens = make([]*services.TokenID, 0)
			}
			body.Tokens = append(body.Tokens, tokenID._ToProtobuf())
		}
	}

	return &services.TransactionBody{
		TransactionFee:           transaction.transactionFee,
		Memo:                     transaction.Transaction.memo,
		TransactionValidDuration: _DurationToProtobuf(transaction.GetTransactionValidDuration()),
		TransactionID:            transaction.transactionID._ToProtobuf(),
		Data: &services.TransactionBody_TokenAssociate{
			TokenAssociate: body,
		},
	}
}

func (transaction *TokenAssociateTransaction) Schedule() (*ScheduleCreateTransaction, error) {
	transaction._RequireNotFrozen()

	scheduled, err := transaction._ConstructScheduleProtobuf()
	if err != nil {
		return nil, err
	}

	return NewScheduleCreateTransaction()._SetSchedulableTransactionBody(scheduled), nil
}

func (transaction *TokenAssociateTransaction) _ConstructScheduleProtobuf() (*services.SchedulableTransactionBody, error) {
	body := &services.TokenAssociateTransactionBody{}
	if transaction.accountID != nil {
		body.Account = transaction.accountID._ToProtobuf()
	}

	if len(transaction.tokens) > 0 {
		for _, tokenID := range transaction.tokens {
			if body.Tokens == nil {
				body.Tokens = make([]*services.TokenID, 0)
			}
			body.Tokens = append(body.Tokens, tokenID._ToProtobuf())
		}
	}
	return &services.SchedulableTransactionBody{
		TransactionFee: transaction.transactionFee,
		Memo:           transaction.Transaction.memo,
		Data: &services.SchedulableTransactionBody_TokenAssociate{
			TokenAssociate: body,
		},
	}, nil
}

func _TokenAssociateTransactionGetMethod(request interface{}, channel *_Channel) _Method {
	return _Method{
		transaction: channel._GetToken().AssociateTokens,
	}
}

func (transaction *TokenAssociateTransaction) IsFrozen() bool {
	return transaction._IsFrozen()
}

// Sign uses the provided privateKey to sign the transaction.
func (transaction *TokenAssociateTransaction) Sign(
	privateKey PrivateKey,
) *TokenAssociateTransaction {
	return transaction.SignWith(privateKey.PublicKey(), privateKey.Sign)
}

// SignWithOperator signs the transaction with client's operator privateKey.
func (transaction *TokenAssociateTransaction) SignWithOperator(
	client *Client,
) (*TokenAssociateTransaction, error) {
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
func (transaction *TokenAssociateTransaction) SignWith(
	publicKey PublicKey,
	signer TransactionSigner,
) *TokenAssociateTransaction {
	if !transaction._KeyAlreadySigned(publicKey) {
		transaction._SignWith(publicKey, signer)
	}

	return transaction
}

// Execute executes the Transaction with the provided client
func (transaction *TokenAssociateTransaction) Execute(
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
		_TokenAssociateTransactionGetMethod,
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

func (transaction *TokenAssociateTransaction) Freeze() (*TokenAssociateTransaction, error) {
	return transaction.FreezeWith(nil)
}

func (transaction *TokenAssociateTransaction) FreezeWith(client *Client) (*TokenAssociateTransaction, error) {
	if transaction.IsFrozen() {
		return transaction, nil
	}
	transaction._InitFee(client)
	err := transaction._ValidateNetworkOnIDs(client)
	if err != nil {
		return &TokenAssociateTransaction{}, err
	}
	if err := transaction._InitTransactionID(client); err != nil {
		return transaction, err
	}
	body := transaction._Build()

	return transaction, _TransactionFreezeWith(&transaction.Transaction, client, body)
}

// GetMaxTransactionFee returns the maximum transaction fee the operator (paying account) is willing to pay.
func (transaction *TokenAssociateTransaction) GetMaxTransactionFee() Hbar {
	return transaction.Transaction.GetMaxTransactionFee()
}

// SetMaxTransactionFee sets the max transaction fee for this TokenAssociateTransaction.
func (transaction *TokenAssociateTransaction) SetMaxTransactionFee(fee Hbar) *TokenAssociateTransaction {
	transaction._RequireNotFrozen()
	transaction.Transaction.SetMaxTransactionFee(fee)
	return transaction
}

// SetRegenerateTransactionID sets if transaction IDs should be regenerated when `TRANSACTION_EXPIRED` is received
func (transaction *TokenAssociateTransaction) SetRegenerateTransactionID(regenerateTransactionID bool) *TokenAssociateTransaction {
	transaction._RequireNotFrozen()
	transaction.Transaction.SetRegenerateTransactionID(regenerateTransactionID)
	return transaction
}

// GetRegenerateTransactionID returns true if transaction ID regeneration is enabled.
func (transaction *TokenAssociateTransaction) GetRegenerateTransactionID() bool {
	return transaction.Transaction.GetRegenerateTransactionID()
}

// GetTransactionMemo returns the memo for this	TokenAssociateTransaction.
func (transaction *TokenAssociateTransaction) GetTransactionMemo() string {
	return transaction.Transaction.GetTransactionMemo()
}

// SetTransactionMemo sets the memo for this TokenAssociateTransaction.
func (transaction *TokenAssociateTransaction) SetTransactionMemo(memo string) *TokenAssociateTransaction {
	transaction._RequireNotFrozen()
	transaction.Transaction.SetTransactionMemo(memo)
	return transaction
}

// GetTransactionValidDuration returns the duration that this transaction is valid for.
func (transaction *TokenAssociateTransaction) GetTransactionValidDuration() time.Duration {
	return transaction.Transaction.GetTransactionValidDuration()
}

// SetTransactionValidDuration sets the valid duration for this TokenAssociateTransaction.
func (transaction *TokenAssociateTransaction) SetTransactionValidDuration(duration time.Duration) *TokenAssociateTransaction {
	transaction._RequireNotFrozen()
	transaction.Transaction.SetTransactionValidDuration(duration)
	return transaction
}

func (transaction *TokenAssociateTransaction) GetTransactionID() TransactionID {
	return transaction.Transaction.GetTransactionID()
}

// SetTransactionID sets the TransactionID for this TokenAssociateTransaction.
func (transaction *TokenAssociateTransaction) SetTransactionID(transactionID TransactionID) *TokenAssociateTransaction {
	transaction._RequireNotFrozen()

	transaction.Transaction.SetTransactionID(transactionID)
	return transaction
}

// SetNodeTokenID sets the _Node TokenID for this TokenAssociateTransaction.
func (transaction *TokenAssociateTransaction) SetNodeAccountIDs(nodeID []AccountID) *TokenAssociateTransaction {
	transaction._RequireNotFrozen()
	transaction.Transaction.SetNodeAccountIDs(nodeID)
	return transaction
}

// SetMaxRetry sets the max number of errors before execution will fail.
func (transaction *TokenAssociateTransaction) SetMaxRetry(count int) *TokenAssociateTransaction {
	transaction.Transaction.SetMaxRetry(count)
	return transaction
}

// AddSignature adds a signature to the Transaction.
func (transaction *TokenAssociateTransaction) AddSignature(publicKey PublicKey, signature []byte) *TokenAssociateTransaction {
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
func (transaction *TokenAssociateTransaction) SetMaxBackoff(max time.Duration) *TokenAssociateTransaction {
	if max.Nanoseconds() < 0 {
		panic("maxBackoff must be a positive duration")
	} else if max.Nanoseconds() < transaction.minBackoff.Nanoseconds() {
		panic("maxBackoff must be greater than or equal to minBackoff")
	}
	transaction.maxBackoff = &max
	return transaction
}

// GetMaxBackoff returns the maximum amount of time to wait between retries.
func (transaction *TokenAssociateTransaction) GetMaxBackoff() time.Duration {
	if transaction.maxBackoff != nil {
		return *transaction.maxBackoff
	}

	return 8 * time.Second
}

// SetMinBackoff sets the minimum amount of time to wait between retries.
func (transaction *TokenAssociateTransaction) SetMinBackoff(min time.Duration) *TokenAssociateTransaction {
	if min.Nanoseconds() < 0 {
		panic("minBackoff must be a positive duration")
	} else if transaction.maxBackoff.Nanoseconds() < min.Nanoseconds() {
		panic("minBackoff must be less than or equal to maxBackoff")
	}
	transaction.minBackoff = &min
	return transaction
}

// GetMinBackoff returns the minimum amount of time to wait between retries.
func (transaction *TokenAssociateTransaction) GetMinBackoff() time.Duration {
	if transaction.minBackoff != nil {
		return *transaction.minBackoff
	}

	return 250 * time.Millisecond
}

func (transaction *TokenAssociateTransaction) _GetLogID() string {
	timestamp := transaction.transactionIDs._GetCurrent().(TransactionID).ValidStart
	return fmt.Sprintf("TokenAssociateTransaction:%d", timestamp.UnixNano())
}

func (transaction *TokenAssociateTransaction) SetLogLevel(level LogLevel) *TokenAssociateTransaction {
	transaction.Transaction.SetLogLevel(level)
	return transaction
}
