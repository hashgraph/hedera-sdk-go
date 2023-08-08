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

// AccountAllowanceDeleteTransaction
// Deletes one or more non-fungible approved allowances from an owner's account. This operation
// will remove the allowances granted to one or more specific non-fungible token serial numbers. Each owner account
// listed as wiping an allowance must sign the transaction. Hbar and fungible token allowances
// can be removed by setting the amount to zero in CryptoApproveAllowance.
type AccountAllowanceDeleteTransaction struct {
	Transaction
	hbarWipe  []*HbarAllowance
	tokenWipe []*TokenAllowance
	nftWipe   []*TokenNftAllowance
}

// NewAccountAllowanceDeleteTransaction
// Creates AccountAllowanceDeleteTransaction whoch deletes one or more non-fungible approved allowances from an owner's account. This operation
// will remove the allowances granted to one or more specific non-fungible token serial numbers. Each owner account
// listed as wiping an allowance must sign the transaction. Hbar and fungible token allowances
// can be removed by setting the amount to zero in CryptoApproveAllowance.
func NewAccountAllowanceDeleteTransaction() *AccountAllowanceDeleteTransaction {
	transaction := AccountAllowanceDeleteTransaction{
		Transaction: _NewTransaction(),
	}

	transaction._SetDefaultMaxTransactionFee(NewHbar(2))

	return &transaction
}

func _AccountAllowanceDeleteTransactionFromProtobuf(transaction Transaction, pb *services.TransactionBody) *AccountAllowanceDeleteTransaction {
	nftWipe := make([]*TokenNftAllowance, 0)

	for _, ap := range pb.GetCryptoDeleteAllowance().GetNftAllowances() {
		temp := _TokenNftWipeAllowanceProtobuf(ap)
		nftWipe = append(nftWipe, &temp)
	}

	return &AccountAllowanceDeleteTransaction{
		Transaction: transaction,
		nftWipe:     nftWipe,
	}
}

// Deprecated
func (transaction *AccountAllowanceDeleteTransaction) DeleteAllHbarAllowances(ownerAccountID *AccountID) *AccountAllowanceDeleteTransaction {
	transaction._RequireNotFrozen()
	transaction.hbarWipe = append(transaction.hbarWipe, &HbarAllowance{
		OwnerAccountID: ownerAccountID,
	})

	return transaction
}

// Deprecated
func (transaction *AccountAllowanceDeleteTransaction) GetAllHbarDeleteAllowances() []*HbarAllowance {
	return transaction.hbarWipe
}

// Deprecated
func (transaction *AccountAllowanceDeleteTransaction) DeleteAllTokenAllowances(tokenID TokenID, ownerAccountID *AccountID) *AccountAllowanceDeleteTransaction {
	transaction._RequireNotFrozen()
	tokenApproval := TokenAllowance{
		TokenID:        &tokenID,
		OwnerAccountID: ownerAccountID,
	}

	transaction.tokenWipe = append(transaction.tokenWipe, &tokenApproval)
	return transaction
}

// Deprecated
func (transaction *AccountAllowanceDeleteTransaction) GetAllTokenDeleteAllowances() []*TokenAllowance {
	return transaction.tokenWipe
}

// DeleteAllTokenNftAllowances
// The non-fungible token allowance/allowances to remove.
func (transaction *AccountAllowanceDeleteTransaction) DeleteAllTokenNftAllowances(nftID NftID, ownerAccountID *AccountID) *AccountAllowanceDeleteTransaction {
	transaction._RequireNotFrozen()

	for _, t := range transaction.nftWipe {
		if t.TokenID.String() == nftID.TokenID.String() {
			if t.OwnerAccountID.String() == ownerAccountID.String() {
				b := false
				for _, s := range t.SerialNumbers {
					if s == nftID.SerialNumber {
						b = true
					}
				}
				if !b {
					t.SerialNumbers = append(t.SerialNumbers, nftID.SerialNumber)
				}
				return transaction
			}
		}
	}

	transaction.nftWipe = append(transaction.nftWipe, &TokenNftAllowance{
		TokenID:        &nftID.TokenID,
		OwnerAccountID: ownerAccountID,
		SerialNumbers:  []int64{nftID.SerialNumber},
		AllSerials:     false,
	})
	return transaction
}

// GetAllTokenNftDeleteAllowances
// Get the non-fungible token allowance/allowances that will be removed.
func (transaction *AccountAllowanceDeleteTransaction) GetAllTokenNftDeleteAllowances() []*TokenNftAllowance {
	return transaction.nftWipe
}

func (transaction *AccountAllowanceDeleteTransaction) _ValidateNetworkOnIDs(client *Client) error {
	if client == nil || !client.autoValidateChecksums {
		return nil
	}

	for _, ap := range transaction.nftWipe {
		if ap.TokenID != nil {
			if err := ap.TokenID.ValidateChecksum(client); err != nil {
				return err
			}
		}

		if ap.OwnerAccountID != nil {
			if err := ap.OwnerAccountID.ValidateChecksum(client); err != nil {
				return err
			}
		}
	}

	return nil
}

func (transaction *AccountAllowanceDeleteTransaction) _Build() *services.TransactionBody {
	nftWipe := make([]*services.NftRemoveAllowance, 0)

	for _, ap := range transaction.nftWipe {
		nftWipe = append(nftWipe, ap._ToWipeProtobuf())
	}

	return &services.TransactionBody{
		TransactionID:            transaction.transactionID._ToProtobuf(),
		TransactionFee:           transaction.transactionFee,
		TransactionValidDuration: _DurationToProtobuf(transaction.GetTransactionValidDuration()),
		Memo:                     transaction.Transaction.memo,
		Data: &services.TransactionBody_CryptoDeleteAllowance{
			CryptoDeleteAllowance: &services.CryptoDeleteAllowanceTransactionBody{
				NftAllowances: nftWipe,
			},
		},
	}
}

func (transaction *AccountAllowanceDeleteTransaction) Schedule() (*ScheduleCreateTransaction, error) {
	transaction._RequireNotFrozen()

	scheduled, err := transaction._ConstructScheduleProtobuf()
	if err != nil {
		return nil, err
	}

	return NewScheduleCreateTransaction()._SetSchedulableTransactionBody(scheduled), nil
}

func (transaction *AccountAllowanceDeleteTransaction) _ConstructScheduleProtobuf() (*services.SchedulableTransactionBody, error) {
	nftWipe := make([]*services.NftRemoveAllowance, 0)

	for _, ap := range transaction.nftWipe {
		nftWipe = append(nftWipe, ap._ToWipeProtobuf())
	}

	return &services.SchedulableTransactionBody{
		TransactionFee: transaction.transactionFee,
		Memo:           transaction.Transaction.memo,
		Data: &services.SchedulableTransactionBody_CryptoDeleteAllowance{
			CryptoDeleteAllowance: &services.CryptoDeleteAllowanceTransactionBody{
				NftAllowances: nftWipe,
			},
		},
	}, nil
}

func _AccountDeleteAllowanceTransactionGetMethod(request interface{}, channel *_Channel) _Method {
	return _Method{
		transaction: channel._GetCrypto().DeleteAllowances,
	}
}

func (transaction *AccountAllowanceDeleteTransaction) IsFrozen() bool {
	return transaction._IsFrozen()
}

// Sign uses the provided privateKey to sign the transaction.
func (transaction *AccountAllowanceDeleteTransaction) Sign(
	privateKey PrivateKey,
) *AccountAllowanceDeleteTransaction {
	return transaction.SignWith(privateKey.PublicKey(), privateKey.Sign)
}

// SignWithOperator signs the transaction with client's operator privateKey.
func (transaction *AccountAllowanceDeleteTransaction) SignWithOperator(
	client *Client,
) (*AccountAllowanceDeleteTransaction, error) {
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
func (transaction *AccountAllowanceDeleteTransaction) SignWith(
	publicKey PublicKey,
	signer TransactionSigner,
) *AccountAllowanceDeleteTransaction {
	if !transaction._KeyAlreadySigned(publicKey) {
		transaction._SignWith(publicKey, signer)
	}

	return transaction
}

// Execute executes the Transaction with the provided client
func (transaction *AccountAllowanceDeleteTransaction) Execute(
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
		_AccountDeleteAllowanceTransactionGetMethod,
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

func (transaction *AccountAllowanceDeleteTransaction) Freeze() (*AccountAllowanceDeleteTransaction, error) {
	return transaction.FreezeWith(nil)
}

func (transaction *AccountAllowanceDeleteTransaction) FreezeWith(client *Client) (*AccountAllowanceDeleteTransaction, error) {
	if transaction.IsFrozen() {
		return transaction, nil
	}
	transaction._InitFee(client)
	if err := transaction._InitTransactionID(client); err != nil {
		return transaction, err
	}
	err := transaction._ValidateNetworkOnIDs(client)
	body := transaction._Build()
	if err != nil {
		return &AccountAllowanceDeleteTransaction{}, err
	}

	return transaction, _TransactionFreezeWith(&transaction.Transaction, client, body)
}

// GetMaxTransactionFee returns the maximum transaction fee the operator (paying account) is willing to pay.
func (transaction *AccountAllowanceDeleteTransaction) GetMaxTransactionFee() Hbar {
	return transaction.Transaction.GetMaxTransactionFee()
}

// SetMaxTransactionFee sets the max transaction fee for this AccountAllowanceDeleteTransaction.
func (transaction *AccountAllowanceDeleteTransaction) SetMaxTransactionFee(fee Hbar) *AccountAllowanceDeleteTransaction {
	transaction._RequireNotFrozen()
	transaction.Transaction.SetMaxTransactionFee(fee)
	return transaction
}

// SetRegenerateTransactionID sets if transaction IDs should be regenerated when `TRANSACTION_EXPIRED` is received
func (transaction *AccountAllowanceDeleteTransaction) SetRegenerateTransactionID(regenerateTransactionID bool) *AccountAllowanceDeleteTransaction {
	transaction._RequireNotFrozen()
	transaction.Transaction.SetRegenerateTransactionID(regenerateTransactionID)
	return transaction
}

// GetRegenerateTransactionID returns true if transaction ID regeneration is enabled.
func (transaction *AccountAllowanceDeleteTransaction) GetRegenerateTransactionID() bool {
	return transaction.Transaction.GetRegenerateTransactionID()
}

// GetTransactionMemo returns the memo for this AccountAllowanceDeleteTransaction.
func (transaction *AccountAllowanceDeleteTransaction) GetTransactionMemo() string {
	return transaction.Transaction.GetTransactionMemo()
}

// SetTransactionMemo sets the memo for this AccountAllowanceDeleteTransaction.
func (transaction *AccountAllowanceDeleteTransaction) SetTransactionMemo(memo string) *AccountAllowanceDeleteTransaction {
	transaction._RequireNotFrozen()
	transaction.Transaction.SetTransactionMemo(memo)
	return transaction
}

// GetTransactionValidDuration returns the duration that this transaction is valid for.
func (transaction *AccountAllowanceDeleteTransaction) GetTransactionValidDuration() time.Duration {
	return transaction.Transaction.GetTransactionValidDuration()
}

// SetTransactionValidDuration sets the valid duration for this AccountAllowanceDeleteTransaction.
func (transaction *AccountAllowanceDeleteTransaction) SetTransactionValidDuration(duration time.Duration) *AccountAllowanceDeleteTransaction {
	transaction._RequireNotFrozen()
	transaction.Transaction.SetTransactionValidDuration(duration)
	return transaction
}

// GetTransactionID returns the TransactionID for this AccountAllowanceDeleteTransaction.
func (transaction *AccountAllowanceDeleteTransaction) GetTransactionID() TransactionID {
	return transaction.Transaction.GetTransactionID()
}

// SetTransactionID sets the TransactionID for this AccountAllowanceDeleteTransaction.
func (transaction *AccountAllowanceDeleteTransaction) SetTransactionID(transactionID TransactionID) *AccountAllowanceDeleteTransaction {
	transaction._RequireNotFrozen()

	transaction.Transaction.SetTransactionID(transactionID)
	return transaction
}

// SetNodeAccountIDs sets the _Node AccountID for this AccountAllowanceDeleteTransaction.
func (transaction *AccountAllowanceDeleteTransaction) SetNodeAccountIDs(nodeID []AccountID) *AccountAllowanceDeleteTransaction {
	transaction._RequireNotFrozen()
	transaction.Transaction.SetNodeAccountIDs(nodeID)
	return transaction
}

// SetMaxRetry sets the max number of errors before execution will fail.
func (transaction *AccountAllowanceDeleteTransaction) SetMaxRetry(count int) *AccountAllowanceDeleteTransaction {
	transaction.Transaction.SetMaxRetry(count)
	return transaction
}

// AddSignature adds a signature to the transaction.
func (transaction *AccountAllowanceDeleteTransaction) AddSignature(publicKey PublicKey, signature []byte) *AccountAllowanceDeleteTransaction {
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
func (transaction *AccountAllowanceDeleteTransaction) SetMaxBackoff(max time.Duration) *AccountAllowanceDeleteTransaction {
	if max.Nanoseconds() < 0 {
		panic("maxBackoff must be a positive duration")
	} else if max.Nanoseconds() < transaction.minBackoff.Nanoseconds() {
		panic("maxBackoff must be greater than or equal to minBackoff")
	}
	transaction.maxBackoff = &max
	return transaction
}

// GetMaxBackoff returns the max back off for this AccountAllowanceDeleteTransaction.
func (transaction *AccountAllowanceDeleteTransaction) GetMaxBackoff() time.Duration {
	if transaction.maxBackoff != nil {
		return *transaction.maxBackoff
	}

	return 8 * time.Second
}

// SetMinBackoff sets the min back off for this AccountAllowanceDeleteTransaction.
func (transaction *AccountAllowanceDeleteTransaction) SetMinBackoff(min time.Duration) *AccountAllowanceDeleteTransaction {
	if min.Nanoseconds() < 0 {
		panic("minBackoff must be a positive duration")
	} else if transaction.maxBackoff.Nanoseconds() < min.Nanoseconds() {
		panic("minBackoff must be less than or equal to maxBackoff")
	}
	transaction.minBackoff = &min
	return transaction
}

// GetMinBackoff returns the minimum amount of time to wait between retries.
func (transaction *AccountAllowanceDeleteTransaction) GetMinBackoff() time.Duration {
	if transaction.minBackoff != nil {
		return *transaction.minBackoff
	}

	return 250 * time.Millisecond
}

func (transaction *AccountAllowanceDeleteTransaction) _GetLogID() string {
	timestamp := transaction.transactionIDs._GetCurrent().(TransactionID).ValidStart
	return fmt.Sprintf("AccountAllowanceDeleteTransaction:%d", timestamp.UnixNano())
}
