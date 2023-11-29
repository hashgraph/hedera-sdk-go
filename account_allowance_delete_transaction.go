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
	transaction
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
	this := AccountAllowanceDeleteTransaction{
		transaction: _NewTransaction(),
	}
	this.e = &this
	this._SetDefaultMaxTransactionFee(NewHbar(2))

	return &this
}

func _AccountAllowanceDeleteTransactionFromProtobuf(transaction transaction, pb *services.TransactionBody) *AccountAllowanceDeleteTransaction {
	nftWipe := make([]*TokenNftAllowance, 0)

	for _, ap := range pb.GetCryptoDeleteAllowance().GetNftAllowances() {
		temp := _TokenNftWipeAllowanceProtobuf(ap)
		nftWipe = append(nftWipe, &temp)
	}

	return &AccountAllowanceDeleteTransaction{
		transaction: transaction,
		nftWipe:     nftWipe,
	}
}

// Deprecated
func (this *AccountAllowanceDeleteTransaction) DeleteAllHbarAllowances(ownerAccountID *AccountID) *AccountAllowanceDeleteTransaction {
	this._RequireNotFrozen()
	this.hbarWipe = append(this.hbarWipe, &HbarAllowance{
		OwnerAccountID: ownerAccountID,
	})

	return this
}

// Deprecated
func (this *AccountAllowanceDeleteTransaction) GetAllHbarDeleteAllowances() []*HbarAllowance {
	return this.hbarWipe
}

// Deprecated
func (this *AccountAllowanceDeleteTransaction) DeleteAllTokenAllowances(tokenID TokenID, ownerAccountID *AccountID) *AccountAllowanceDeleteTransaction {
	this._RequireNotFrozen()
	tokenApproval := TokenAllowance{
		TokenID:        &tokenID,
		OwnerAccountID: ownerAccountID,
	}

	this.tokenWipe = append(this.tokenWipe, &tokenApproval)
	return this
}

// Deprecated
func (this *AccountAllowanceDeleteTransaction) GetAllTokenDeleteAllowances() []*TokenAllowance {
	return this.tokenWipe
}

// DeleteAllTokenNftAllowances
// The non-fungible token allowance/allowances to remove.
func (this *AccountAllowanceDeleteTransaction) DeleteAllTokenNftAllowances(nftID NftID, ownerAccountID *AccountID) *AccountAllowanceDeleteTransaction {
	this._RequireNotFrozen()

	for _, t := range this.nftWipe {
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
				return this
			}
		}
	}

	this.nftWipe = append(this.nftWipe, &TokenNftAllowance{
		TokenID:        &nftID.TokenID,
		OwnerAccountID: ownerAccountID,
		SerialNumbers:  []int64{nftID.SerialNumber},
		AllSerials:     false,
	})
	return this
}

// GetAllTokenNftDeleteAllowances
// Get the non-fungible token allowance/allowances that will be removed.
func (this *AccountAllowanceDeleteTransaction) GetAllTokenNftDeleteAllowances() []*TokenNftAllowance {
	return this.nftWipe
}
func (this *AccountAllowanceDeleteTransaction) Schedule() (*ScheduleCreateTransaction, error) {
	this._RequireNotFrozen()

	scheduled, err := this.buildScheduled()
	if err != nil {
		return nil, err
	}

	return NewScheduleCreateTransaction()._SetSchedulableTransactionBody(scheduled), nil
}

func (this *AccountAllowanceDeleteTransaction) IsFrozen() bool {
	return this._IsFrozen()
}

// Sign uses the provided privateKey to sign the transaction.
func (this *AccountAllowanceDeleteTransaction) Sign(
	privateKey PrivateKey,
) *AccountAllowanceDeleteTransaction {
	this.transaction.Sign(privateKey)
	return this
}

// SignWithOperator signs the transaction with client's operator privateKey.
func (this *AccountAllowanceDeleteTransaction) SignWithOperator(
	client *Client,
) (*AccountAllowanceDeleteTransaction, error) {
	// If the transaction is not signed by the _Operator, we need
	// to sign the transaction with the _Operator
	_,err := this.transaction.SignWithOperator(client)
	return this,err
}

// SignWith executes the TransactionSigner and adds the resulting signature data to the transaction's signature map
// with the publicKey as the map key.
func (this *AccountAllowanceDeleteTransaction) SignWith(
	publicKey PublicKey,
	signer TransactionSigner,
) *AccountAllowanceDeleteTransaction {
	this.transaction.SignWith(publicKey,signer);
	return this
}

func (this *AccountAllowanceDeleteTransaction) Freeze() (*AccountAllowanceDeleteTransaction, error) {
	_,err := this.transaction.Freeze()
	return this,err
}

func (this *AccountAllowanceDeleteTransaction) FreezeWith(client *Client) (*AccountAllowanceDeleteTransaction, error) {
	_,err := this.transaction.FreezeWith(client)
	return this, err
}

// GetMaxTransactionFee returns the maximum transaction fee the operator (paying account) is willing to pay.
func (this *AccountAllowanceDeleteTransaction) GetMaxTransactionFee() Hbar {
	return this.transaction.GetMaxTransactionFee()
}

// SetMaxTransactionFee sets the max transaction fee for this AccountAllowanceDeleteTransaction.
func (this *AccountAllowanceDeleteTransaction) SetMaxTransactionFee(fee Hbar) *AccountAllowanceDeleteTransaction {
	this._RequireNotFrozen()
	this.transaction.SetMaxTransactionFee(fee)
	return this
}

// SetRegenerateTransactionID sets if transaction IDs should be regenerated when `TRANSACTION_EXPIRED` is received
func (this *AccountAllowanceDeleteTransaction) SetRegenerateTransactionID(regenerateTransactionID bool) *AccountAllowanceDeleteTransaction {
	this._RequireNotFrozen()
	this.transaction.SetRegenerateTransactionID(regenerateTransactionID)
	return this
}

// GetRegenerateTransactionID returns true if transaction ID regeneration is enabled.
func (this *AccountAllowanceDeleteTransaction) GetRegenerateTransactionID() bool {
	return this.transaction.GetRegenerateTransactionID()
}

// GetTransactionMemo returns the memo for this AccountAllowanceDeleteTransaction.
func (this *AccountAllowanceDeleteTransaction) GetTransactionMemo() string {
	return this.transaction.GetTransactionMemo()
}

// SetTransactionMemo sets the memo for this AccountAllowanceDeleteTransaction.
func (this *AccountAllowanceDeleteTransaction) SetTransactionMemo(memo string) *AccountAllowanceDeleteTransaction {
	this._RequireNotFrozen()
	this.transaction.SetTransactionMemo(memo)
	return this
}

// GetTransactionValidDuration returns the duration that this transaction is valid for.
func (transaction *AccountAllowanceDeleteTransaction) GetTransactionValidDuration() time.Duration {
	return transaction.transaction.GetTransactionValidDuration()
}

// SetTransactionValidDuration sets the valid duration for this AccountAllowanceDeleteTransaction.
func (this *AccountAllowanceDeleteTransaction) SetTransactionValidDuration(duration time.Duration) *AccountAllowanceDeleteTransaction {
	this._RequireNotFrozen()
	this.transaction.SetTransactionValidDuration(duration)
	return this
}

// GetTransactionID returns the TransactionID for this AccountAllowanceDeleteTransaction.
func (this *AccountAllowanceDeleteTransaction) GetTransactionID() TransactionID {
	return this.transaction.GetTransactionID()
}

// SetTransactionID sets the TransactionID for this AccountAllowanceDeleteTransaction.
func (this *AccountAllowanceDeleteTransaction) SetTransactionID(transactionID TransactionID) *AccountAllowanceDeleteTransaction {
	this._RequireNotFrozen()

	this.transaction.SetTransactionID(transactionID)
	return this
}

// SetNodeAccountIDs sets the _Node AccountID for this AccountAllowanceDeleteTransaction.
func (this *AccountAllowanceDeleteTransaction) SetNodeAccountIDs(nodeID []AccountID) *AccountAllowanceDeleteTransaction {
	this._RequireNotFrozen()
	this.transaction.SetNodeAccountIDs(nodeID)
	return this
}

// SetMaxRetry sets the max number of errors before execution will fail.
func (this *AccountAllowanceDeleteTransaction) SetMaxRetry(count int) *AccountAllowanceDeleteTransaction {
	this.transaction.SetMaxRetry(count)
	return this
}

// AddSignature adds a signature to the transaction.
func (this *AccountAllowanceDeleteTransaction) AddSignature(publicKey PublicKey, signature []byte) *AccountAllowanceDeleteTransaction {
	this.transaction.AddSignature(publicKey, signature)
	return this
}

// SetMaxBackoff The maximum amount of time to wait between retries.
// Every retry attempt will increase the wait time exponentially until it reaches this time.
func (this *AccountAllowanceDeleteTransaction) SetMaxBackoff(max time.Duration) *AccountAllowanceDeleteTransaction {
	this.transaction.SetMaxBackoff(max)
    return this
}

// SetMinBackoff sets the min back off for this AccountAllowanceDeleteTransaction.
func (this *AccountAllowanceDeleteTransaction) SetMinBackoff(min time.Duration) *AccountAllowanceDeleteTransaction {
	this.transaction.SetMinBackoff(min)
	return this
}

func (transaction *AccountAllowanceDeleteTransaction) _GetLogID() string {
	timestamp := transaction.transactionIDs._GetCurrent().(TransactionID).ValidStart
	return fmt.Sprintf("AccountAllowanceDeleteTransaction:%d", timestamp.UnixNano())
}

// ----------- overriden functions ----------------
func (this *AccountAllowanceDeleteTransaction) getName() string {
	return "AccountAllowanceDeleteTransaction"
}

func (this *AccountAllowanceDeleteTransaction) validateNetworkOnIDs(client *Client) error {
	if client == nil || !client.autoValidateChecksums {
		return nil
	}

	for _, ap := range this.nftWipe {
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

func (this *AccountAllowanceDeleteTransaction) build() *services.TransactionBody {
	nftWipe := make([]*services.NftRemoveAllowance, 0)

	for _, ap := range this.nftWipe {
		nftWipe = append(nftWipe, ap._ToWipeProtobuf())
	}

	return &services.TransactionBody{
		TransactionID:            this.transactionID._ToProtobuf(),
		TransactionFee:           this.transactionFee,
		TransactionValidDuration: _DurationToProtobuf(this.GetTransactionValidDuration()),
		Memo:                     this.transaction.memo,
		Data: &services.TransactionBody_CryptoDeleteAllowance{
			CryptoDeleteAllowance: &services.CryptoDeleteAllowanceTransactionBody{
				NftAllowances: nftWipe,
			},
		},
	}
}

func (this *AccountAllowanceDeleteTransaction) buildScheduled() (*services.SchedulableTransactionBody, error) {
	nftWipe := make([]*services.NftRemoveAllowance, 0)

	for _, ap := range this.nftWipe {
		nftWipe = append(nftWipe, ap._ToWipeProtobuf())
	}

	return &services.SchedulableTransactionBody{
		TransactionFee: this.transactionFee,
		Memo:           this.transaction.memo,
		Data: &services.SchedulableTransactionBody_CryptoDeleteAllowance{
			CryptoDeleteAllowance: &services.CryptoDeleteAllowanceTransactionBody{
				NftAllowances: nftWipe,
			},
		},
	}, nil
}

func (this *AccountAllowanceDeleteTransaction) getMethod(channel *_Channel) _Method {
	return _Method{
		transaction: channel._GetCrypto().DeleteAllowances,
	}
}
