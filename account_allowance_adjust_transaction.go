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
	transaction
	hbarAllowances  []*HbarAllowance
	tokenAllowances []*TokenAllowance
	nftAllowances   []*TokenNftAllowance
}

func NewAccountAllowanceAdjustTransaction() *AccountAllowanceAdjustTransaction {
	this := AccountAllowanceAdjustTransaction{
		transaction: _NewTransaction(),
	}
	this.e = &this
	this._SetDefaultMaxTransactionFee(NewHbar(2))

	return &this
}

func (this *AccountAllowanceAdjustTransaction) _AdjustHbarAllowance(ownerAccountID *AccountID, id AccountID, amount Hbar) *AccountAllowanceAdjustTransaction {
	this._RequireNotFrozen()
	this.hbarAllowances = append(this.hbarAllowances, &HbarAllowance{
		SpenderAccountID: &id,
		OwnerAccountID:   ownerAccountID,
		Amount:           amount.AsTinybar(),
	})

	return this
}

// Deprecated
func (this *AccountAllowanceAdjustTransaction) AddHbarAllowance(id AccountID, amount Hbar) *AccountAllowanceAdjustTransaction {
	return this._AdjustHbarAllowance(nil, id, amount)
}

// Deprecated
func (this *AccountAllowanceAdjustTransaction) GrantHbarAllowance(ownerAccountID AccountID, id AccountID, amount Hbar) *AccountAllowanceAdjustTransaction {
	return this._AdjustHbarAllowance(&ownerAccountID, id, amount)
}

// Deprecated
func (this *AccountAllowanceAdjustTransaction) RevokeHbarAllowance(ownerAccountID AccountID, id AccountID, amount Hbar) *AccountAllowanceAdjustTransaction {
	return this._AdjustHbarAllowance(&ownerAccountID, id, amount.Negated())
}

// Deprecated
func (this *AccountAllowanceAdjustTransaction) GetHbarAllowances() []*HbarAllowance {
	return this.hbarAllowances
}

func (this *AccountAllowanceAdjustTransaction) _AdjustTokenAllowance(tokenID TokenID, ownerAccountID *AccountID, accountID AccountID, amount int64) *AccountAllowanceAdjustTransaction {
	this._RequireNotFrozen()
	tokenApproval := TokenAllowance{
		TokenID:          &tokenID,
		SpenderAccountID: &accountID,
		OwnerAccountID:   ownerAccountID,
		Amount:           amount,
	}

	this.tokenAllowances = append(this.tokenAllowances, &tokenApproval)
	return this
}

// Deprecated
func (this *AccountAllowanceAdjustTransaction) AddTokenAllowance(tokenID TokenID, accountID AccountID, amount int64) *AccountAllowanceAdjustTransaction {
	return this._AdjustTokenAllowance(tokenID, nil, accountID, amount)
}

// Deprecated
func (this *AccountAllowanceAdjustTransaction) GrantTokenAllowance(tokenID TokenID, ownerAccountID AccountID, accountID AccountID, amount int64) *AccountAllowanceAdjustTransaction {
	return this._AdjustTokenAllowance(tokenID, &ownerAccountID, accountID, amount)
}

// Deprecated
func (this *AccountAllowanceAdjustTransaction) RevokeTokenAllowance(tokenID TokenID, ownerAccountID AccountID, accountID AccountID, amount uint64) *AccountAllowanceAdjustTransaction {
	return this._AdjustTokenAllowance(tokenID, &ownerAccountID, accountID, -int64(amount))
}

// Deprecated
func (this *AccountAllowanceAdjustTransaction) GetTokenAllowances() []*TokenAllowance {
	return this.tokenAllowances
}

func (this *AccountAllowanceAdjustTransaction) _AdjustTokenNftAllowance(nftID NftID, ownerAccountID *AccountID, accountID AccountID) *AccountAllowanceAdjustTransaction {
	this._RequireNotFrozen()

	for _, t := range this.nftAllowances {
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
				return this
			}
		}
	}

	this.nftAllowances = append(this.nftAllowances, &TokenNftAllowance{
		TokenID:          &nftID.TokenID,
		SpenderAccountID: &accountID,
		OwnerAccountID:   ownerAccountID,
		SerialNumbers:    []int64{nftID.SerialNumber},
		AllSerials:       false,
	})
	return this
}

// Deprecated
func (this *AccountAllowanceAdjustTransaction) AddTokenNftAllowance(nftID NftID, accountID AccountID) *AccountAllowanceAdjustTransaction {
	return this._AdjustTokenNftAllowance(nftID, nil, accountID)
}

// Deprecated
func (this *AccountAllowanceAdjustTransaction) GrantTokenNftAllowance(nftID NftID, ownerAccountID AccountID, accountID AccountID) *AccountAllowanceAdjustTransaction {
	return this._AdjustTokenNftAllowance(nftID, &ownerAccountID, accountID)
}

// Deprecated
func (this *AccountAllowanceAdjustTransaction) RevokeTokenNftAllowance(nftID NftID, ownerAccountID AccountID, accountID AccountID) *AccountAllowanceAdjustTransaction {
	return this._AdjustTokenNftAllowance(nftID, &ownerAccountID, accountID)
}

func (this *AccountAllowanceAdjustTransaction) _AdjustTokenNftAllowanceAllSerials(tokenID TokenID, ownerAccountID *AccountID, spenderAccount AccountID, allSerials bool) *AccountAllowanceAdjustTransaction {
	for _, t := range this.nftAllowances {
		if t.TokenID.String() == tokenID.String() {
			if t.SpenderAccountID.String() == spenderAccount.String() {
				t.SerialNumbers = []int64{}
				t.AllSerials = true
				return this
			}
		}
	}

	this.nftAllowances = append(this.nftAllowances, &TokenNftAllowance{
		TokenID:          &tokenID,
		SpenderAccountID: &spenderAccount,
		OwnerAccountID:   ownerAccountID,
		SerialNumbers:    []int64{},
		AllSerials:       allSerials,
	})
	return this
}

// Deprecated
func (this *AccountAllowanceAdjustTransaction) AddAllTokenNftAllowance(tokenID TokenID, spenderAccount AccountID) *AccountAllowanceAdjustTransaction {
	return this._AdjustTokenNftAllowanceAllSerials(tokenID, nil, spenderAccount, true)
}

// Deprecated
func (this *AccountAllowanceAdjustTransaction) GrantTokenNftAllowanceAllSerials(ownerAccountID AccountID, tokenID TokenID, spenderAccount AccountID) *AccountAllowanceAdjustTransaction {
	return this._AdjustTokenNftAllowanceAllSerials(tokenID, &ownerAccountID, spenderAccount, true)
}

// Deprecated
func (this *AccountAllowanceAdjustTransaction) RevokeTokenNftAllowanceAllSerials(ownerAccountID AccountID, tokenID TokenID, spenderAccount AccountID) *AccountAllowanceAdjustTransaction {
	return this._AdjustTokenNftAllowanceAllSerials(tokenID, &ownerAccountID, spenderAccount, false)
}

// Deprecated
func (this *AccountAllowanceAdjustTransaction) GetTokenNftAllowances() []*TokenNftAllowance {
	return this.nftAllowances
}

// Deprecated
func (this *AccountAllowanceAdjustTransaction) Schedule() (*ScheduleCreateTransaction, error) {
	this._RequireNotFrozen()

	scheduled, err := this.buildScheduled()
	if err != nil {
		return nil, err
	}

	return NewScheduleCreateTransaction()._SetSchedulableTransactionBody(scheduled), nil
}

// Deprecated
func (this *AccountAllowanceAdjustTransaction) IsFrozen() bool {
	return this._IsFrozen()
}

// Deprecated
func (this *AccountAllowanceAdjustTransaction) Sign(
	privateKey PrivateKey,
) *AccountAllowanceAdjustTransaction {
	this.transaction.Sign(privateKey)
	return this
}

// Deprecated
func (this *AccountAllowanceAdjustTransaction) SignWithOperator(
	client *Client,
) (*AccountAllowanceAdjustTransaction, error) {
	// If the transaction is not signed by the _Operator, we need
	// to sign the transaction with the _Operator
	_,err := this.transaction.SignWithOperator(client)
	return this,err
}

// Deprecated
func (this *AccountAllowanceAdjustTransaction) SignWith(
	publicKey PublicKey,
	signer TransactionSigner,
) *AccountAllowanceAdjustTransaction {
	this.transaction.SignWith(publicKey,signer);
	return this
}

// Deprecated
func (this *AccountAllowanceAdjustTransaction) Freeze() (*AccountAllowanceAdjustTransaction, error) {
	_,err := this.transaction.Freeze()
	return this,err
}

// Deprecated
func (this *AccountAllowanceAdjustTransaction) FreezeWith(client *Client) (*AccountAllowanceAdjustTransaction, error) {
	_,err := this.transaction.FreezeWith(client)
	return this, err
}

func (this *AccountAllowanceAdjustTransaction) GetMaxTransactionFee() Hbar {
	return this.transaction.GetMaxTransactionFee()
}

// SetMaxTransactionFee sets the max transaction fee for this AccountAllowanceAdjustTransaction.
func (this *AccountAllowanceAdjustTransaction) SetMaxTransactionFee(fee Hbar) *AccountAllowanceAdjustTransaction {
	this._RequireNotFrozen()
	this.transaction.SetMaxTransactionFee(fee)
	return this
}

// SetRegenerateTransactionID sets if transaction IDs should be regenerated when `TRANSACTION_EXPIRED` is received
func (this *AccountAllowanceAdjustTransaction) SetRegenerateTransactionID(regenerateTransactionID bool) *AccountAllowanceAdjustTransaction {
	this._RequireNotFrozen()
	this.transaction.SetRegenerateTransactionID(regenerateTransactionID)
	return this
}

// GetRegenerateTransactionID returns true if transaction ID regeneration is enabled.
func (this *AccountAllowanceAdjustTransaction) GetRegenerateTransactionID() bool {
	return this.transaction.GetRegenerateTransactionID()
}

func (this *AccountAllowanceAdjustTransaction) GetTransactionMemo() string {
	return this.transaction.GetTransactionMemo()
}

// SetTransactionMemo sets the memo for this AccountAllowanceAdjustTransaction.
func (this *AccountAllowanceAdjustTransaction) SetTransactionMemo(memo string) *AccountAllowanceAdjustTransaction {
	this._RequireNotFrozen()
	this.transaction.SetTransactionMemo(memo)
	return this
}

func (this *AccountAllowanceAdjustTransaction) GetTransactionValidDuration() time.Duration {
	return this.transaction.GetTransactionValidDuration()
}

// SetTransactionValidDuration sets the valid duration for this AccountAllowanceAdjustTransaction.
func (this *AccountAllowanceAdjustTransaction) SetTransactionValidDuration(duration time.Duration) *AccountAllowanceAdjustTransaction {
	this._RequireNotFrozen()
	this.transaction.SetTransactionValidDuration(duration)
	return this
}

func (this *AccountAllowanceAdjustTransaction) GetTransactionID() TransactionID {
	return this.transaction.GetTransactionID()
}

// SetTransactionID sets the TransactionID for this AccountAllowanceAdjustTransaction.
func (this *AccountAllowanceAdjustTransaction) SetTransactionID(transactionID TransactionID) *AccountAllowanceAdjustTransaction {
	this._RequireNotFrozen()

	this.transaction.SetTransactionID(transactionID)
	return this
}

// SetNodeAccountIDs sets the _Node AccountID for this AccountAllowanceAdjustTransaction.
func (this *AccountAllowanceAdjustTransaction) SetNodeAccountIDs(nodeID []AccountID) *AccountAllowanceAdjustTransaction {
	this.transaction.SetNodeAccountIDs(nodeID)
	return this
}

func (this *AccountAllowanceAdjustTransaction) SetMaxRetry(count int) *AccountAllowanceAdjustTransaction {
	this.transaction.SetMaxRetry(count)
	return this
}

func (this *AccountAllowanceAdjustTransaction) AddSignature(publicKey PublicKey, signature []byte) *AccountAllowanceAdjustTransaction {
	this.transaction.AddSignature(publicKey, signature)
	return this
}

func (this *AccountAllowanceAdjustTransaction) SetMaxBackoff(max time.Duration) *AccountAllowanceAdjustTransaction {
	this.transaction.SetMaxBackoff(max)
	return this
}

func (this *AccountAllowanceAdjustTransaction) GetMaxBackoff() time.Duration {
	return this.transaction.GetMaxBackoff()
}

func (this *AccountAllowanceAdjustTransaction) SetMinBackoff(min time.Duration) *AccountAllowanceAdjustTransaction {
	this.SetMinBackoff(min)
	return this
}

func (this *AccountAllowanceAdjustTransaction) GetMinBackoff() time.Duration {
	return this.transaction.GetMaxBackoff()
}

func (this *AccountAllowanceAdjustTransaction) _GetLogID() string {
	timestamp := this.transactionIDs._GetCurrent().(TransactionID).ValidStart
	return fmt.Sprintf("AccountAllowanceAdjustTransaction:%d", timestamp.UnixNano())
}

// ----------- overriden functions ----------------

func (transaction *AccountAllowanceAdjustTransaction) getName() string {
	return "AccountAllowanceAdjustTransaction"
}

func (this *AccountAllowanceAdjustTransaction) validateNetworkOnIDs(client *Client) error {
	if client == nil || !client.autoValidateChecksums {
		return nil
	}

	for _, ap := range this.hbarAllowances {
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

	for _, ap := range this.tokenAllowances {
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

	for _, ap := range this.nftAllowances {
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

func (this *AccountAllowanceAdjustTransaction) build() *services.TransactionBody {
	return &services.TransactionBody{}
}
func (this *AccountAllowanceAdjustTransaction) buildScheduled() (*services.SchedulableTransactionBody, error) {
	return &services.SchedulableTransactionBody{}, nil
}

