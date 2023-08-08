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

// Deprecated
type AccountAllowanceAdjustTransaction struct {
	Transaction
	hbarAllowances  []*HbarAllowance
	tokenAllowances []*TokenAllowance
	nftAllowances   []*TokenNftAllowance
}

func NewAccountAllowanceAdjustTransaction() *AccountAllowanceAdjustTransaction {
	transaction := AccountAllowanceAdjustTransaction{
		Transaction: _NewTransaction(),
	}

	transaction._SetDefaultMaxTransactionFee(NewHbar(2))

	return &transaction
}

func (transaction *AccountAllowanceAdjustTransaction) _AdjustHbarAllowance(ownerAccountID *AccountID, id AccountID, amount Hbar) *AccountAllowanceAdjustTransaction {
	transaction._RequireNotFrozen()
	transaction.hbarAllowances = append(transaction.hbarAllowances, &HbarAllowance{
		SpenderAccountID: &id,
		OwnerAccountID:   ownerAccountID,
		Amount:           amount.AsTinybar(),
	})

	return transaction
}

// Deprecated
func (transaction *AccountAllowanceAdjustTransaction) AddHbarAllowance(id AccountID, amount Hbar) *AccountAllowanceAdjustTransaction {
	return transaction._AdjustHbarAllowance(nil, id, amount)
}

// Deprecated
func (transaction *AccountAllowanceAdjustTransaction) GrantHbarAllowance(ownerAccountID AccountID, id AccountID, amount Hbar) *AccountAllowanceAdjustTransaction {
	return transaction._AdjustHbarAllowance(&ownerAccountID, id, amount)
}

// Deprecated
func (transaction *AccountAllowanceAdjustTransaction) RevokeHbarAllowance(ownerAccountID AccountID, id AccountID, amount Hbar) *AccountAllowanceAdjustTransaction {
	return transaction._AdjustHbarAllowance(&ownerAccountID, id, amount.Negated())
}

// Deprecated
func (transaction *AccountAllowanceAdjustTransaction) GetHbarAllowances() []*HbarAllowance {
	return transaction.hbarAllowances
}

func (transaction *AccountAllowanceAdjustTransaction) _AdjustTokenAllowance(tokenID TokenID, ownerAccountID *AccountID, accountID AccountID, amount int64) *AccountAllowanceAdjustTransaction {
	transaction._RequireNotFrozen()
	tokenApproval := TokenAllowance{
		TokenID:          &tokenID,
		SpenderAccountID: &accountID,
		OwnerAccountID:   ownerAccountID,
		Amount:           amount,
	}

	transaction.tokenAllowances = append(transaction.tokenAllowances, &tokenApproval)
	return transaction
}

// Deprecated
func (transaction *AccountAllowanceAdjustTransaction) AddTokenAllowance(tokenID TokenID, accountID AccountID, amount int64) *AccountAllowanceAdjustTransaction {
	return transaction._AdjustTokenAllowance(tokenID, nil, accountID, amount)
}

// Deprecated
func (transaction *AccountAllowanceAdjustTransaction) GrantTokenAllowance(tokenID TokenID, ownerAccountID AccountID, accountID AccountID, amount int64) *AccountAllowanceAdjustTransaction {
	return transaction._AdjustTokenAllowance(tokenID, &ownerAccountID, accountID, amount)
}

// Deprecated
func (transaction *AccountAllowanceAdjustTransaction) RevokeTokenAllowance(tokenID TokenID, ownerAccountID AccountID, accountID AccountID, amount uint64) *AccountAllowanceAdjustTransaction {
	return transaction._AdjustTokenAllowance(tokenID, &ownerAccountID, accountID, -int64(amount))
}

// Deprecated
func (transaction *AccountAllowanceAdjustTransaction) GetTokenAllowances() []*TokenAllowance {
	return transaction.tokenAllowances
}

func (transaction *AccountAllowanceAdjustTransaction) _AdjustTokenNftAllowance(nftID NftID, ownerAccountID *AccountID, accountID AccountID) *AccountAllowanceAdjustTransaction {
	transaction._RequireNotFrozen()

	for _, t := range transaction.nftAllowances {
		if t.TokenID.String() == nftID.TokenID.String() {
			if t.SpenderAccountID.String() == accountID.String() {
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

	transaction.nftAllowances = append(transaction.nftAllowances, &TokenNftAllowance{
		TokenID:          &nftID.TokenID,
		SpenderAccountID: &accountID,
		OwnerAccountID:   ownerAccountID,
		SerialNumbers:    []int64{nftID.SerialNumber},
		AllSerials:       false,
	})
	return transaction
}

// Deprecated
func (transaction *AccountAllowanceAdjustTransaction) AddTokenNftAllowance(nftID NftID, accountID AccountID) *AccountAllowanceAdjustTransaction {
	return transaction._AdjustTokenNftAllowance(nftID, nil, accountID)
}

// Deprecated
func (transaction *AccountAllowanceAdjustTransaction) GrantTokenNftAllowance(nftID NftID, ownerAccountID AccountID, accountID AccountID) *AccountAllowanceAdjustTransaction {
	return transaction._AdjustTokenNftAllowance(nftID, &ownerAccountID, accountID)
}

// Deprecated
func (transaction *AccountAllowanceAdjustTransaction) RevokeTokenNftAllowance(nftID NftID, ownerAccountID AccountID, accountID AccountID) *AccountAllowanceAdjustTransaction {
	return transaction._AdjustTokenNftAllowance(nftID, &ownerAccountID, accountID)
}

func (transaction *AccountAllowanceAdjustTransaction) _AdjustTokenNftAllowanceAllSerials(tokenID TokenID, ownerAccountID *AccountID, spenderAccount AccountID, allSerials bool) *AccountAllowanceAdjustTransaction {
	for _, t := range transaction.nftAllowances {
		if t.TokenID.String() == tokenID.String() {
			if t.SpenderAccountID.String() == spenderAccount.String() {
				t.SerialNumbers = []int64{}
				t.AllSerials = true
				return transaction
			}
		}
	}

	transaction.nftAllowances = append(transaction.nftAllowances, &TokenNftAllowance{
		TokenID:          &tokenID,
		SpenderAccountID: &spenderAccount,
		OwnerAccountID:   ownerAccountID,
		SerialNumbers:    []int64{},
		AllSerials:       allSerials,
	})
	return transaction
}

// Deprecated
func (transaction *AccountAllowanceAdjustTransaction) AddAllTokenNftAllowance(tokenID TokenID, spenderAccount AccountID) *AccountAllowanceAdjustTransaction {
	return transaction._AdjustTokenNftAllowanceAllSerials(tokenID, nil, spenderAccount, true)
}

// Deprecated
func (transaction *AccountAllowanceAdjustTransaction) GrantTokenNftAllowanceAllSerials(ownerAccountID AccountID, tokenID TokenID, spenderAccount AccountID) *AccountAllowanceAdjustTransaction {
	return transaction._AdjustTokenNftAllowanceAllSerials(tokenID, &ownerAccountID, spenderAccount, true)
}

// Deprecated
func (transaction *AccountAllowanceAdjustTransaction) RevokeTokenNftAllowanceAllSerials(ownerAccountID AccountID, tokenID TokenID, spenderAccount AccountID) *AccountAllowanceAdjustTransaction {
	return transaction._AdjustTokenNftAllowanceAllSerials(tokenID, &ownerAccountID, spenderAccount, false)
}

// Deprecated
func (transaction *AccountAllowanceAdjustTransaction) GetTokenNftAllowances() []*TokenNftAllowance {
	return transaction.nftAllowances
}

func (transaction *AccountAllowanceAdjustTransaction) _ValidateNetworkOnIDs(client *Client) error {
	if client == nil || !client.autoValidateChecksums {
		return nil
	}

	for _, ap := range transaction.hbarAllowances {
		if ap.SpenderAccountID != nil {
			if err := ap.SpenderAccountID.ValidateChecksum(client); err != nil {
				return err
			}
		}

		if ap.OwnerAccountID != nil {
			if err := ap.OwnerAccountID.ValidateChecksum(client); err != nil {
				return err
			}
		}
	}

	for _, ap := range transaction.tokenAllowances {
		if ap.SpenderAccountID != nil {
			if err := ap.SpenderAccountID.ValidateChecksum(client); err != nil {
				return err
			}
		}

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

	for _, ap := range transaction.nftAllowances {
		if ap.SpenderAccountID != nil {
			if err := ap.SpenderAccountID.ValidateChecksum(client); err != nil {
				return err
			}
		}

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

func (transaction *AccountAllowanceAdjustTransaction) _Build() *services.TransactionBody {
	return &services.TransactionBody{}
}

// Deprecated
func (transaction *AccountAllowanceAdjustTransaction) Schedule() (*ScheduleCreateTransaction, error) {
	transaction._RequireNotFrozen()

	scheduled, err := transaction._ConstructScheduleProtobuf()
	if err != nil {
		return nil, err
	}

	return NewScheduleCreateTransaction()._SetSchedulableTransactionBody(scheduled), nil
}

func (transaction *AccountAllowanceAdjustTransaction) _ConstructScheduleProtobuf() (*services.SchedulableTransactionBody, error) {
	return &services.SchedulableTransactionBody{}, nil
}

func _AccountAdjustAllowanceTransactionGetMethod(request interface{}, channel *_Channel) _Method {
	return _Method{}
}

// Deprecated
func (transaction *AccountAllowanceAdjustTransaction) IsFrozen() bool {
	return transaction._IsFrozen()
}

// Deprecated
func (transaction *AccountAllowanceAdjustTransaction) Sign(
	privateKey PrivateKey,
) *AccountAllowanceAdjustTransaction {
	return transaction.SignWith(privateKey.PublicKey(), privateKey.Sign)
}

// Deprecated
func (transaction *AccountAllowanceAdjustTransaction) SignWithOperator(
	client *Client,
) (*AccountAllowanceAdjustTransaction, error) {
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

// Deprecated
func (transaction *AccountAllowanceAdjustTransaction) SignWith(
	publicKey PublicKey,
	signer TransactionSigner,
) *AccountAllowanceAdjustTransaction {
	if !transaction._KeyAlreadySigned(publicKey) {
		transaction._SignWith(publicKey, signer)
	}

	return transaction
}

// Deprecated
func (transaction *AccountAllowanceAdjustTransaction) Execute(
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
		_AccountAdjustAllowanceTransactionGetMethod,
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

// Deprecated
func (transaction *AccountAllowanceAdjustTransaction) Freeze() (*AccountAllowanceAdjustTransaction, error) {
	return transaction.FreezeWith(nil)
}

// Deprecated
func (transaction *AccountAllowanceAdjustTransaction) FreezeWith(client *Client) (*AccountAllowanceAdjustTransaction, error) {
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
		return &AccountAllowanceAdjustTransaction{}, err
	}

	return transaction, _TransactionFreezeWith(&transaction.Transaction, client, body)
}

func (transaction *AccountAllowanceAdjustTransaction) GetMaxTransactionFee() Hbar {
	return transaction.Transaction.GetMaxTransactionFee()
}

// SetMaxTransactionFee sets the max transaction fee for this AccountAllowanceAdjustTransaction.
func (transaction *AccountAllowanceAdjustTransaction) SetMaxTransactionFee(fee Hbar) *AccountAllowanceAdjustTransaction {
	transaction._RequireNotFrozen()
	transaction.Transaction.SetMaxTransactionFee(fee)
	return transaction
}

// SetRegenerateTransactionID sets if transaction IDs should be regenerated when `TRANSACTION_EXPIRED` is received
func (transaction *AccountAllowanceAdjustTransaction) SetRegenerateTransactionID(regenerateTransactionID bool) *AccountAllowanceAdjustTransaction {
	transaction._RequireNotFrozen()
	transaction.Transaction.SetRegenerateTransactionID(regenerateTransactionID)
	return transaction
}

// GetRegenerateTransactionID returns true if transaction ID regeneration is enabled.
func (transaction *AccountAllowanceAdjustTransaction) GetRegenerateTransactionID() bool {
	return transaction.Transaction.GetRegenerateTransactionID()
}

func (transaction *AccountAllowanceAdjustTransaction) GetTransactionMemo() string {
	return transaction.Transaction.GetTransactionMemo()
}

// SetTransactionMemo sets the memo for this AccountAllowanceAdjustTransaction.
func (transaction *AccountAllowanceAdjustTransaction) SetTransactionMemo(memo string) *AccountAllowanceAdjustTransaction {
	transaction._RequireNotFrozen()
	transaction.Transaction.SetTransactionMemo(memo)
	return transaction
}

func (transaction *AccountAllowanceAdjustTransaction) GetTransactionValidDuration() time.Duration {
	return transaction.Transaction.GetTransactionValidDuration()
}

// SetTransactionValidDuration sets the valid duration for this AccountAllowanceAdjustTransaction.
func (transaction *AccountAllowanceAdjustTransaction) SetTransactionValidDuration(duration time.Duration) *AccountAllowanceAdjustTransaction {
	transaction._RequireNotFrozen()
	transaction.Transaction.SetTransactionValidDuration(duration)
	return transaction
}

func (transaction *AccountAllowanceAdjustTransaction) GetTransactionID() TransactionID {
	return transaction.Transaction.GetTransactionID()
}

// SetTransactionID sets the TransactionID for this AccountAllowanceAdjustTransaction.
func (transaction *AccountAllowanceAdjustTransaction) SetTransactionID(transactionID TransactionID) *AccountAllowanceAdjustTransaction {
	transaction._RequireNotFrozen()

	transaction.Transaction.SetTransactionID(transactionID)
	return transaction
}

// SetNodeAccountIDs sets the _Node AccountID for this AccountAllowanceAdjustTransaction.
func (transaction *AccountAllowanceAdjustTransaction) SetNodeAccountIDs(nodeID []AccountID) *AccountAllowanceAdjustTransaction {
	transaction._RequireNotFrozen()
	transaction.Transaction.SetNodeAccountIDs(nodeID)
	return transaction
}

func (transaction *AccountAllowanceAdjustTransaction) SetMaxRetry(count int) *AccountAllowanceAdjustTransaction {
	transaction.Transaction.SetMaxRetry(count)
	return transaction
}

func (transaction *AccountAllowanceAdjustTransaction) AddSignature(publicKey PublicKey, signature []byte) *AccountAllowanceAdjustTransaction {
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

func (transaction *AccountAllowanceAdjustTransaction) SetMaxBackoff(max time.Duration) *AccountAllowanceAdjustTransaction {
	if max.Nanoseconds() < 0 {
		panic("maxBackoff must be a positive duration")
	} else if max.Nanoseconds() < transaction.minBackoff.Nanoseconds() {
		panic("maxBackoff must be greater than or equal to minBackoff")
	}
	transaction.maxBackoff = &max
	return transaction
}

func (transaction *AccountAllowanceAdjustTransaction) GetMaxBackoff() time.Duration {
	if transaction.maxBackoff != nil {
		return *transaction.maxBackoff
	}

	return 8 * time.Second
}

func (transaction *AccountAllowanceAdjustTransaction) SetMinBackoff(min time.Duration) *AccountAllowanceAdjustTransaction {
	if min.Nanoseconds() < 0 {
		panic("minBackoff must be a positive duration")
	} else if transaction.maxBackoff.Nanoseconds() < min.Nanoseconds() {
		panic("minBackoff must be less than or equal to maxBackoff")
	}
	transaction.minBackoff = &min
	return transaction
}

func (transaction *AccountAllowanceAdjustTransaction) GetMinBackoff() time.Duration {
	if transaction.minBackoff != nil {
		return *transaction.minBackoff
	}

	return 250 * time.Millisecond
}

func (transaction *AccountAllowanceAdjustTransaction) _GetLogID() string {
	timestamp := transaction.transactionIDs._GetCurrent().(TransactionID).ValidStart
	return fmt.Sprintf("AccountAllowanceAdjustTransaction:%d", timestamp.UnixNano())
}
