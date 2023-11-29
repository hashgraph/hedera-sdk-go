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

// AccountAllowanceApproveTransaction
// Creates one or more hbar/token approved allowances <b>relative to the owner account specified in the allowances of
// this transaction</b>. Each allowance grants a spender the right to transfer a pre-determined amount of the owner's
// hbar/token to any other account of the spender's choice. If the owner is not specified in any allowance, the payer
// of transaction is considered to be the owner for that particular allowance.
// Setting the amount to zero in CryptoAllowance or TokenAllowance will remove the respective allowance for the spender.
//
// (So if account <tt>0.0.X</tt> pays for this transaction and owner is not specified in the allowance,
// then at consensus each spender account will have new allowances to spend hbar or tokens from <tt>0.0.X</tt>).
type AccountAllowanceApproveTransaction struct {
	transaction
	hbarAllowances  []*HbarAllowance
	tokenAllowances []*TokenAllowance
	nftAllowances   []*TokenNftAllowance
}

// NewAccountAllowanceApproveTransaction
// Creates an AccountAloowanceApproveTransaction which creates
// one or more hbar/token approved allowances relative to the owner account specified in the allowances of
// this transaction. Each allowance grants a spender the right to transfer a pre-determined amount of the owner's
// hbar/token to any other account of the spender's choice. If the owner is not specified in any allowance, the payer
// of transaction is considered to be the owner for that particular allowance.
// Setting the amount to zero in CryptoAllowance or TokenAllowance will remove the respective allowance for the spender.
//
// (So if account 0.0.X pays for this transaction and owner is not specified in the allowance,
// then at consensus each spender account will have new allowances to spend hbar or tokens from 0.0.X).
func NewAccountAllowanceApproveTransaction() *AccountAllowanceApproveTransaction {
	this := AccountAllowanceApproveTransaction{
		transaction: _NewTransaction(),
	}
	this.e = &this
	this._SetDefaultMaxTransactionFee(NewHbar(2))

	return &this
}

func _AccountAllowanceApproveTransactionFromProtobuf(this transaction, pb *services.TransactionBody) *AccountAllowanceApproveTransaction {
	accountApproval := make([]*HbarAllowance, 0)
	tokenApproval := make([]*TokenAllowance, 0)
	nftApproval := make([]*TokenNftAllowance, 0)

	for _, ap := range pb.GetCryptoApproveAllowance().GetCryptoAllowances() {
		temp := _HbarAllowanceFromProtobuf(ap)
		accountApproval = append(accountApproval, &temp)
	}

	for _, ap := range pb.GetCryptoApproveAllowance().GetTokenAllowances() {
		temp := _TokenAllowanceFromProtobuf(ap)
		tokenApproval = append(tokenApproval, &temp)
	}

	for _, ap := range pb.GetCryptoApproveAllowance().GetNftAllowances() {
		temp := _TokenNftAllowanceFromProtobuf(ap)
		nftApproval = append(nftApproval, &temp)
	}

	return &AccountAllowanceApproveTransaction{
		transaction:     this,
		hbarAllowances:  accountApproval,
		tokenAllowances: tokenApproval,
		nftAllowances:   nftApproval,
	}
}

func (this *AccountAllowanceApproveTransaction) _ApproveHbarApproval(ownerAccountID *AccountID, id AccountID, amount Hbar) *AccountAllowanceApproveTransaction {
	this._RequireNotFrozen()
	this.hbarAllowances = append(this.hbarAllowances, &HbarAllowance{
		SpenderAccountID: &id,
		Amount:           amount.AsTinybar(),
		OwnerAccountID:   ownerAccountID,
	})

	return this
}

// AddHbarApproval
// Deprecated - Use ApproveHbarAllowance instead
func (this *AccountAllowanceApproveTransaction) AddHbarApproval(id AccountID, amount Hbar) *AccountAllowanceApproveTransaction {
	return this._ApproveHbarApproval(nil, id, amount)
}

// ApproveHbarApproval
// Deprecated - Use ApproveHbarAllowance instead
func (this *AccountAllowanceApproveTransaction) ApproveHbarApproval(ownerAccountID AccountID, id AccountID, amount Hbar) *AccountAllowanceApproveTransaction {
	return this._ApproveHbarApproval(&ownerAccountID, id, amount)
}

// ApproveHbarAllowance
// Approves allowance of hbar transfers for a spender.
func (this *AccountAllowanceApproveTransaction) ApproveHbarAllowance(ownerAccountID AccountID, id AccountID, amount Hbar) *AccountAllowanceApproveTransaction {
	return this._ApproveHbarApproval(&ownerAccountID, id, amount)
}

// List of hbar allowance records
func (this *AccountAllowanceApproveTransaction) GetHbarAllowances() []*HbarAllowance {
	return this.hbarAllowances
}

func (this *AccountAllowanceApproveTransaction) _ApproveTokenApproval(tokenID TokenID, ownerAccountID *AccountID, accountID AccountID, amount int64) *AccountAllowanceApproveTransaction {
	this._RequireNotFrozen()
	tokenApproval := TokenAllowance{
		TokenID:          &tokenID,
		SpenderAccountID: &accountID,
		Amount:           amount,
		OwnerAccountID:   ownerAccountID,
	}

	this.tokenAllowances = append(this.tokenAllowances, &tokenApproval)
	return this
}

// Deprecated - Use ApproveTokenAllowance instead
func (this *AccountAllowanceApproveTransaction) AddTokenApproval(tokenID TokenID, accountID AccountID, amount int64) *AccountAllowanceApproveTransaction {
	return this._ApproveTokenApproval(tokenID, nil, accountID, amount)
}

// ApproveTokenApproval
// Deprecated - Use ApproveTokenAllowance instead
func (this *AccountAllowanceApproveTransaction) ApproveTokenApproval(tokenID TokenID, ownerAccountID AccountID, accountID AccountID, amount int64) *AccountAllowanceApproveTransaction {
	return this._ApproveTokenApproval(tokenID, &ownerAccountID, accountID, amount)
}

// ApproveTokenAllowance
// Approve allowance of fungible token transfers for a spender.
func (this *AccountAllowanceApproveTransaction) ApproveTokenAllowance(tokenID TokenID, ownerAccountID AccountID, accountID AccountID, amount int64) *AccountAllowanceApproveTransaction {
	return this._ApproveTokenApproval(tokenID, &ownerAccountID, accountID, amount)
}

// List of token allowance records
func (this *AccountAllowanceApproveTransaction) GetTokenAllowances() []*TokenAllowance {
	return this.tokenAllowances
}

func (this *AccountAllowanceApproveTransaction) _ApproveTokenNftApproval(nftID NftID, ownerAccountID *AccountID, spenderAccountID *AccountID, delegatingSpenderAccountId *AccountID) *AccountAllowanceApproveTransaction {
	this._RequireNotFrozen()

	for _, t := range this.nftAllowances {
		if t.TokenID.String() == nftID.TokenID.String() {
			if t.SpenderAccountID.String() == spenderAccountID.String() {
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
		TokenID:           &nftID.TokenID,
		SpenderAccountID:  spenderAccountID,
		SerialNumbers:     []int64{nftID.SerialNumber},
		AllSerials:        false,
		OwnerAccountID:    ownerAccountID,
		DelegatingSpender: delegatingSpenderAccountId,
	})
	return this
}

// AddTokenNftApproval
// Deprecated - Use ApproveTokenNftAllowance instead
func (this *AccountAllowanceApproveTransaction) AddTokenNftApproval(nftID NftID, accountID AccountID) *AccountAllowanceApproveTransaction {
	return this._ApproveTokenNftApproval(nftID, nil, &accountID, nil)
}

// ApproveTokenNftApproval
// Deprecated - Use ApproveTokenNftAllowance instead
func (this *AccountAllowanceApproveTransaction) ApproveTokenNftApproval(nftID NftID, ownerAccountID AccountID, accountID AccountID) *AccountAllowanceApproveTransaction {
	return this._ApproveTokenNftApproval(nftID, &ownerAccountID, &accountID, nil)
}

func (this *AccountAllowanceApproveTransaction) ApproveTokenNftAllowanceWithDelegatingSpender(nftID NftID, ownerAccountID AccountID, spenderAccountId AccountID, delegatingSpenderAccountID AccountID) *AccountAllowanceApproveTransaction {
	this._RequireNotFrozen()
	return this._ApproveTokenNftApproval(nftID, &ownerAccountID, &spenderAccountId, &delegatingSpenderAccountID)
}

// ApproveTokenNftAllowance
// Approve allowance of non-fungible token transfers for a spender.
func (this *AccountAllowanceApproveTransaction) ApproveTokenNftAllowance(nftID NftID, ownerAccountID AccountID, accountID AccountID) *AccountAllowanceApproveTransaction {
	return this._ApproveTokenNftApproval(nftID, &ownerAccountID, &accountID, nil)
}

func (this *AccountAllowanceApproveTransaction) _ApproveTokenNftAllowanceAllSerials(tokenID TokenID, ownerAccountID *AccountID, spenderAccount AccountID) *AccountAllowanceApproveTransaction {
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
		SerialNumbers:    []int64{},
		AllSerials:       true,
		OwnerAccountID:   ownerAccountID,
	})
	return this
}

// AddAllTokenNftApproval
// Approve allowance of non-fungible token transfers for a spender.
// Spender has access to all of the owner's NFT units of type tokenId (currently
// owned and any in the future).
func (this *AccountAllowanceApproveTransaction) AddAllTokenNftApproval(tokenID TokenID, spenderAccount AccountID) *AccountAllowanceApproveTransaction {
	return this._ApproveTokenNftAllowanceAllSerials(tokenID, nil, spenderAccount)
}

// ApproveTokenNftAllowanceAllSerials
// Approve allowance of non-fungible token transfers for a spender.
// Spender has access to all of the owner's NFT units of type tokenId (currently
// owned and any in the future).
func (this *AccountAllowanceApproveTransaction) ApproveTokenNftAllowanceAllSerials(tokenID TokenID, ownerAccountID AccountID, spenderAccount AccountID) *AccountAllowanceApproveTransaction {
	return this._ApproveTokenNftAllowanceAllSerials(tokenID, &ownerAccountID, spenderAccount)
}

// List of NFT allowance records
func (this *AccountAllowanceApproveTransaction) GetTokenNftAllowances() []*TokenNftAllowance {
	return this.nftAllowances
}

func (this *AccountAllowanceApproveTransaction) Schedule() (*ScheduleCreateTransaction, error) {
	this._RequireNotFrozen()

	scheduled, err := this.buildProtoBody()
	if err != nil {
		return nil, err
	}

	return NewScheduleCreateTransaction()._SetSchedulableTransactionBody(scheduled), nil
}

func (this *AccountAllowanceApproveTransaction) IsFrozen() bool {
	return this._IsFrozen()
}

// Sign uses the provided privateKey to sign the transaction.
func (this *AccountAllowanceApproveTransaction) Sign(
	privateKey PrivateKey,
) *AccountAllowanceApproveTransaction {
	this.transaction.SignWith(privateKey.PublicKey(), privateKey.Sign)
	return this
}

// SignWithOperator signs the transaction with client's operator privateKey.
func (this *AccountAllowanceApproveTransaction) SignWithOperator(
	client *Client,
) (*AccountAllowanceApproveTransaction, error) {
	_,err := this.transaction.SignWithOperator(client)
	return this, err
}

// SignWith executes the TransactionSigner and adds the resulting signature data to the transaction's signature map
// with the publicKey as the map key.
func (this *AccountAllowanceApproveTransaction) SignWith(
	publicKey PublicKey,
	signer TransactionSigner,
) *AccountAllowanceApproveTransaction {
	this.transaction.SignWith(publicKey, signer)
	return this
}

func (this *AccountAllowanceApproveTransaction) Freeze() (*AccountAllowanceApproveTransaction, error) {
	_,err := this.transaction.Freeze()
	return this, err
}

func (this *AccountAllowanceApproveTransaction) FreezeWith(client *Client) (*AccountAllowanceApproveTransaction, error) {
	_, err := this.transaction.FreezeWith(client)
	return this, err
}

// GetMaxTransactionFee returns the maximum transaction fee the operator (paying account) is willing to pay.
func (this *AccountAllowanceApproveTransaction) GetMaxTransactionFee() Hbar {
	return this.transaction.GetMaxTransactionFee()
}

// SetMaxTransactionFee sets the maximum transaction fee the operator (paying account) is willing to pay.
func (this *AccountAllowanceApproveTransaction) SetMaxTransactionFee(fee Hbar) *AccountAllowanceApproveTransaction {
	this._RequireNotFrozen()
	this.transaction.SetMaxTransactionFee(fee)
	return this
}

// SetRegenerateTransactionID sets if transaction IDs should be regenerated when `TRANSACTION_EXPIRED` is received
func (this *AccountAllowanceApproveTransaction) SetRegenerateTransactionID(regenerateTransactionID bool) *AccountAllowanceApproveTransaction {
	this._RequireNotFrozen()
	this.transaction.SetRegenerateTransactionID(regenerateTransactionID)
	return this
}

// GetRegenerateTransactionID returns true if transaction ID regeneration is enabled.
func (this *AccountAllowanceApproveTransaction) GetRegenerateTransactionID() bool {
	return this.transaction.GetRegenerateTransactionID()
}

// GetTransactionMemo returns the memo for this AccountAllowanceApproveTransaction.
func (this *AccountAllowanceApproveTransaction) GetTransactionMemo() string {
	return this.transaction.GetTransactionMemo()
}

// SetTransactionMemo sets the memo for this AccountAllowanceApproveTransaction.
func (this *AccountAllowanceApproveTransaction) SetTransactionMemo(memo string) *AccountAllowanceApproveTransaction {
	this._RequireNotFrozen()
	this.transaction.SetTransactionMemo(memo)
	return this
}

// GetTransactionValidDuration returns the duration that this transaction is valid for.
func (this *AccountAllowanceApproveTransaction) GetTransactionValidDuration() time.Duration {
	return this.transaction.GetTransactionValidDuration()
}

// SetTransactionValidDuration sets the valid duration for this AccountAllowanceApproveTransaction.
func (this *AccountAllowanceApproveTransaction) SetTransactionValidDuration(duration time.Duration) *AccountAllowanceApproveTransaction {
	this._RequireNotFrozen()
	this.transaction.SetTransactionValidDuration(duration)
	return this
}

// GetTransactionID gets the TransactionID for this AccountAllowanceApproveTransaction.
func (this *AccountAllowanceApproveTransaction) GetTransactionID() TransactionID {
	return this.transaction.GetTransactionID()
}

// SetTransactionID sets the TransactionID for this AccountAllowanceApproveTransaction.
func (this *AccountAllowanceApproveTransaction) SetTransactionID(transactionID TransactionID) *AccountAllowanceApproveTransaction {
	this._RequireNotFrozen()

	this.transaction.SetTransactionID(transactionID)
	return this
}

// SetNodeAccountIDs sets the _Node AccountID for this AccountAllowanceApproveTransaction.
func (this *AccountAllowanceApproveTransaction) SetNodeAccountIDs(nodeID []AccountID) *AccountAllowanceApproveTransaction {
	this._RequireNotFrozen()
	this.transaction.SetNodeAccountIDs(nodeID)
	return this
}

// SetMaxRetry sets the max number of errors before execution will fail.
func (this *AccountAllowanceApproveTransaction) SetMaxRetry(count int) *AccountAllowanceApproveTransaction {
	this.transaction.SetMaxRetry(count)
	return this
}

// AddSignature adds a signature to the transaction.
func (this *AccountAllowanceApproveTransaction) AddSignature(publicKey PublicKey, signature []byte) *AccountAllowanceApproveTransaction {
	this.transaction.AddSignature(publicKey, signature)
	return this
}

// SetMaxBackoff The maximum amount of time to wait between retries.
// Every retry attempt will increase the wait time exponentially until it reaches this time.
func (this *AccountAllowanceApproveTransaction) SetMaxBackoff(max time.Duration) *AccountAllowanceApproveTransaction {
	this.SetMaxBackoff(max)
	return this
}

// GetMaxBackoff returns the max back off for this AccountAllowanceApproveTransaction.
func (this *AccountAllowanceApproveTransaction) GetMaxBackoff() time.Duration {
	return this.transaction.GetMaxBackoff()
}

// SetMinBackoff sets the min back off for this AccountAllowanceApproveTransaction.
func (this *AccountAllowanceApproveTransaction) SetMinBackoff(min time.Duration) *AccountAllowanceApproveTransaction {
	this.SetMinBackoff(min)
	return this
}

// GetMinBackoff returns the min back off for this AccountAllowanceApproveTransaction.
func (this *AccountAllowanceApproveTransaction) GetMinBackoff() time.Duration {
	return this.GetMinBackoff()
}

func (this *AccountAllowanceApproveTransaction) _GetLogID() string {
	timestamp := this.transactionIDs._GetCurrent().(TransactionID).ValidStart
	return fmt.Sprintf("AccountAllowanceApproveTransaction:%d", timestamp.UnixNano())
}


// ----------- overriden functions ----------------

func (this *AccountAllowanceApproveTransaction) getName() string {
	return "AccountAllowanceApproveTransaction"
}
func (this *AccountAllowanceApproveTransaction) validateNetworkOnIDs(client *Client) error {
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

func (this *AccountAllowanceApproveTransaction) build() *services.TransactionBody {
	accountApproval := make([]*services.CryptoAllowance, 0)
	tokenApproval := make([]*services.TokenAllowance, 0)
	nftApproval := make([]*services.NftAllowance, 0)

	for _, ap := range this.hbarAllowances {
		accountApproval = append(accountApproval, ap._ToProtobuf())
	}

	for _, ap := range this.tokenAllowances {
		tokenApproval = append(tokenApproval, ap._ToProtobuf())
	}

	for _, ap := range this.nftAllowances {
		nftApproval = append(nftApproval, ap._ToProtobuf())
	}

	return &services.TransactionBody{
		TransactionID:            this.transactionID._ToProtobuf(),
		TransactionFee:           this.transactionFee,
		TransactionValidDuration: _DurationToProtobuf(this.GetTransactionValidDuration()),
		Memo:                     this.transaction.memo,
		Data: &services.TransactionBody_CryptoApproveAllowance{
			CryptoApproveAllowance: &services.CryptoApproveAllowanceTransactionBody{
				CryptoAllowances: accountApproval,
				NftAllowances:    nftApproval,
				TokenAllowances:  tokenApproval,
			},
		},
	}
}

func (this *AccountAllowanceApproveTransaction) buildProtoBody() (*services.SchedulableTransactionBody, error) {
	accountApproval := make([]*services.CryptoAllowance, 0)
	tokenApproval := make([]*services.TokenAllowance, 0)
	nftApproval := make([]*services.NftAllowance, 0)

	for _, ap := range this.hbarAllowances {
		accountApproval = append(accountApproval, ap._ToProtobuf())
	}

	for _, ap := range this.tokenAllowances {
		tokenApproval = append(tokenApproval, ap._ToProtobuf())
	}

	for _, ap := range this.nftAllowances {
		nftApproval = append(nftApproval, ap._ToProtobuf())
	}

	return &services.SchedulableTransactionBody{
		TransactionFee: this.transactionFee,
		Memo:           this.transaction.memo,
		Data: &services.SchedulableTransactionBody_CryptoApproveAllowance{
			CryptoApproveAllowance: &services.CryptoApproveAllowanceTransactionBody{
				CryptoAllowances: accountApproval,
				NftAllowances:    nftApproval,
				TokenAllowances:  tokenApproval,
			},
		},
	}, nil
}

func (this *AccountAllowanceApproveTransaction) getMethod(channel *_Channel) _Method {
	return _Method{
		transaction: channel._GetCrypto().ApproveAllowances,
	}
 }
