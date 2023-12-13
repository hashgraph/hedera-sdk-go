package hedera

/*-
 *
 * Hedera Go SDK
 *
 * Copyright (C) 2020 - 2023 Hedera Hashgraph, LLC
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use tx file except in compliance with the License.
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

// Deprecated
type AccountAllowanceAdjustTransaction struct {
	Transaction
	hbarAllowances  []*HbarAllowance
	tokenAllowances []*TokenAllowance
	nftAllowances   []*TokenNftAllowance
}

func NewAccountAllowanceAdjustTransaction() *AccountAllowanceAdjustTransaction {
	tx := AccountAllowanceAdjustTransaction{
		Transaction: _NewTransaction(),
	}
	tx.e = &tx
	tx._SetDefaultMaxTransactionFee(NewHbar(2))

	return &tx
}

func (tx *AccountAllowanceAdjustTransaction) _AdjustHbarAllowance(ownerAccountID *AccountID, id AccountID, amount Hbar) *AccountAllowanceAdjustTransaction {
	tx._RequireNotFrozen()
	tx.hbarAllowances = append(tx.hbarAllowances, &HbarAllowance{
		SpenderAccountID: &id,
		OwnerAccountID:   ownerAccountID,
		Amount:           amount.AsTinybar(),
	})

	return tx
}

// Deprecated
func (tx *AccountAllowanceAdjustTransaction) AddHbarAllowance(id AccountID, amount Hbar) *AccountAllowanceAdjustTransaction {
	return tx._AdjustHbarAllowance(nil, id, amount)
}

// Deprecated
func (tx *AccountAllowanceAdjustTransaction) GrantHbarAllowance(ownerAccountID AccountID, id AccountID, amount Hbar) *AccountAllowanceAdjustTransaction {
	return tx._AdjustHbarAllowance(&ownerAccountID, id, amount)
}

// Deprecated
func (tx *AccountAllowanceAdjustTransaction) RevokeHbarAllowance(ownerAccountID AccountID, id AccountID, amount Hbar) *AccountAllowanceAdjustTransaction {
	return tx._AdjustHbarAllowance(&ownerAccountID, id, amount.Negated())
}

// Deprecated
func (tx *AccountAllowanceAdjustTransaction) GetHbarAllowances() []*HbarAllowance {
	return tx.hbarAllowances
}

func (tx *AccountAllowanceAdjustTransaction) _AdjustTokenAllowance(tokenID TokenID, ownerAccountID *AccountID, accountID AccountID, amount int64) *AccountAllowanceAdjustTransaction {
	tx._RequireNotFrozen()
	tokenApproval := TokenAllowance{
		TokenID:          &tokenID,
		SpenderAccountID: &accountID,
		OwnerAccountID:   ownerAccountID,
		Amount:           amount,
	}

	tx.tokenAllowances = append(tx.tokenAllowances, &tokenApproval)
	return tx
}

// Deprecated
func (tx *AccountAllowanceAdjustTransaction) AddTokenAllowance(tokenID TokenID, accountID AccountID, amount int64) *AccountAllowanceAdjustTransaction {
	return tx._AdjustTokenAllowance(tokenID, nil, accountID, amount)
}

// Deprecated
func (tx *AccountAllowanceAdjustTransaction) GrantTokenAllowance(tokenID TokenID, ownerAccountID AccountID, accountID AccountID, amount int64) *AccountAllowanceAdjustTransaction {
	return tx._AdjustTokenAllowance(tokenID, &ownerAccountID, accountID, amount)
}

// Deprecated
func (tx *AccountAllowanceAdjustTransaction) RevokeTokenAllowance(tokenID TokenID, ownerAccountID AccountID, accountID AccountID, amount uint64) *AccountAllowanceAdjustTransaction {
	return tx._AdjustTokenAllowance(tokenID, &ownerAccountID, accountID, -int64(amount))
}

// Deprecated
func (tx *AccountAllowanceAdjustTransaction) GetTokenAllowances() []*TokenAllowance {
	return tx.tokenAllowances
}

func (tx *AccountAllowanceAdjustTransaction) _AdjustTokenNftAllowance(nftID NftID, ownerAccountID *AccountID, accountID AccountID) *AccountAllowanceAdjustTransaction {
	tx._RequireNotFrozen()

	for _, t := range tx.nftAllowances {
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
				return tx
			}
		}
	}

	tx.nftAllowances = append(tx.nftAllowances, &TokenNftAllowance{
		TokenID:          &nftID.TokenID,
		SpenderAccountID: &accountID,
		OwnerAccountID:   ownerAccountID,
		SerialNumbers:    []int64{nftID.SerialNumber},
		AllSerials:       false,
	})
	return tx
}

// Deprecated
func (tx *AccountAllowanceAdjustTransaction) AddTokenNftAllowance(nftID NftID, accountID AccountID) *AccountAllowanceAdjustTransaction {
	return tx._AdjustTokenNftAllowance(nftID, nil, accountID)
}

// Deprecated
func (tx *AccountAllowanceAdjustTransaction) GrantTokenNftAllowance(nftID NftID, ownerAccountID AccountID, accountID AccountID) *AccountAllowanceAdjustTransaction {
	return tx._AdjustTokenNftAllowance(nftID, &ownerAccountID, accountID)
}

// Deprecated
func (tx *AccountAllowanceAdjustTransaction) RevokeTokenNftAllowance(nftID NftID, ownerAccountID AccountID, accountID AccountID) *AccountAllowanceAdjustTransaction {
	return tx._AdjustTokenNftAllowance(nftID, &ownerAccountID, accountID)
}

func (tx *AccountAllowanceAdjustTransaction) _AdjustTokenNftAllowanceAllSerials(tokenID TokenID, ownerAccountID *AccountID, spenderAccount AccountID, allSerials bool) *AccountAllowanceAdjustTransaction {
	for _, t := range tx.nftAllowances {
		if t.TokenID.String() == tokenID.String() {
			if t.SpenderAccountID.String() == spenderAccount.String() {
				t.SerialNumbers = []int64{}
				t.AllSerials = true
				return tx
			}
		}
	}

	tx.nftAllowances = append(tx.nftAllowances, &TokenNftAllowance{
		TokenID:          &tokenID,
		SpenderAccountID: &spenderAccount,
		OwnerAccountID:   ownerAccountID,
		SerialNumbers:    []int64{},
		AllSerials:       allSerials,
	})
	return tx
}

// Deprecated
func (tx *AccountAllowanceAdjustTransaction) AddAllTokenNftAllowance(tokenID TokenID, spenderAccount AccountID) *AccountAllowanceAdjustTransaction {
	return tx._AdjustTokenNftAllowanceAllSerials(tokenID, nil, spenderAccount, true)
}

// Deprecated
func (tx *AccountAllowanceAdjustTransaction) GrantTokenNftAllowanceAllSerials(ownerAccountID AccountID, tokenID TokenID, spenderAccount AccountID) *AccountAllowanceAdjustTransaction {
	return tx._AdjustTokenNftAllowanceAllSerials(tokenID, &ownerAccountID, spenderAccount, true)
}

// Deprecated
func (tx *AccountAllowanceAdjustTransaction) RevokeTokenNftAllowanceAllSerials(ownerAccountID AccountID, tokenID TokenID, spenderAccount AccountID) *AccountAllowanceAdjustTransaction {
	return tx._AdjustTokenNftAllowanceAllSerials(tokenID, &ownerAccountID, spenderAccount, false)
}

// Deprecated
func (tx *AccountAllowanceAdjustTransaction) GetTokenNftAllowances() []*TokenNftAllowance {
	return tx.nftAllowances
}

// Deprecated
func (tx *AccountAllowanceAdjustTransaction) Sign(
	privateKey PrivateKey,
) *AccountAllowanceAdjustTransaction {
	tx.Transaction.Sign(privateKey)
	return tx
}

// Deprecated
func (tx *AccountAllowanceAdjustTransaction) SignWithOperator(
	client *Client,
) (*AccountAllowanceAdjustTransaction, error) {
	// If the transaction is not signed by the _Operator, we need
	// to sign the transaction with the _Operator
	_, err := tx.Transaction.SignWithOperator(client)
	if err != nil {
		return nil, err
	}
	return tx, nil
}

// Deprecated
func (tx *AccountAllowanceAdjustTransaction) SignWith(
	publicKey PublicKey,
	signer TransactionSigner,
) *AccountAllowanceAdjustTransaction {
	tx.Transaction.SignWith(publicKey, signer)
	return tx
}

// Deprecated
func (this *AccountAllowanceAdjustTransaction) Freeze() (*AccountAllowanceAdjustTransaction, error) {
	_, err := this.Transaction.Freeze()
	return this, err
}

// Deprecated
func (this *AccountAllowanceAdjustTransaction) FreezeWith(client *Client) (*AccountAllowanceAdjustTransaction, error) {
	_, err := this.Transaction.FreezeWith(client)
	return this, err
}

// SetMaxTransactionFee sets the max transaction fee for tx AccountAllowanceAdjustTransaction.
func (tx *AccountAllowanceAdjustTransaction) SetMaxTransactionFee(fee Hbar) *AccountAllowanceAdjustTransaction {
	tx._RequireNotFrozen()
	tx.Transaction.SetMaxTransactionFee(fee)
	return tx
}

// SetRegenerateTransactionID sets if transaction IDs should be regenerated when `TRANSACTION_EXPIRED` is received
func (tx *AccountAllowanceAdjustTransaction) SetRegenerateTransactionID(regenerateTransactionID bool) *AccountAllowanceAdjustTransaction {
	tx._RequireNotFrozen()
	tx.Transaction.SetRegenerateTransactionID(regenerateTransactionID)
	return tx
}

// SetTransactionMemo sets the memo for tx AccountAllowanceAdjustTransaction.
func (tx *AccountAllowanceAdjustTransaction) SetTransactionMemo(memo string) *AccountAllowanceAdjustTransaction {
	tx._RequireNotFrozen()
	tx.Transaction.SetTransactionMemo(memo)
	return tx
}

func (tx *AccountAllowanceAdjustTransaction) GetTransactionValidDuration() time.Duration {
	return tx.Transaction.GetTransactionValidDuration()
}

// SetTransactionValidDuration sets the valid duration for tx AccountAllowanceAdjustTransaction.
func (tx *AccountAllowanceAdjustTransaction) SetTransactionValidDuration(duration time.Duration) *AccountAllowanceAdjustTransaction {
	tx._RequireNotFrozen()
	tx.Transaction.SetTransactionValidDuration(duration)
	return tx
}

func (tx *AccountAllowanceAdjustTransaction) GetTransactionID() TransactionID {
	return tx.Transaction.GetTransactionID()
}

// SetTransactionID sets the TransactionID for tx AccountAllowanceAdjustTransaction.
func (tx *AccountAllowanceAdjustTransaction) SetTransactionID(transactionID TransactionID) *AccountAllowanceAdjustTransaction {
	tx._RequireNotFrozen()

	tx.Transaction.SetTransactionID(transactionID)
	return tx
}

// SetNodeAccountIDs sets the _Node AccountID for tx AccountAllowanceAdjustTransaction.
func (tx *AccountAllowanceAdjustTransaction) SetNodeAccountIDs(nodeID []AccountID) *AccountAllowanceAdjustTransaction {
	tx.Transaction.SetNodeAccountIDs(nodeID)
	return tx
}

func (tx *AccountAllowanceAdjustTransaction) SetMaxRetry(count int) *AccountAllowanceAdjustTransaction {
	tx.Transaction.SetMaxRetry(count)
	return tx
}

func (tx *AccountAllowanceAdjustTransaction) AddSignature(publicKey PublicKey, signature []byte) *AccountAllowanceAdjustTransaction {
	tx.Transaction.AddSignature(publicKey, signature)
	return tx
}

func (tx *AccountAllowanceAdjustTransaction) SetMaxBackoff(max time.Duration) *AccountAllowanceAdjustTransaction {
	tx.Transaction.SetMaxBackoff(max)
	return tx
}

func (tx *AccountAllowanceAdjustTransaction) SetMinBackoff(min time.Duration) *AccountAllowanceAdjustTransaction {
	tx.Transaction.SetMinBackoff(min)
	return tx
}

// ----------- overridden functions ----------------

func (transaction *AccountAllowanceAdjustTransaction) getName() string {
	return "AccountAllowanceAdjustTransaction"
}

func (tx *AccountAllowanceAdjustTransaction) validateNetworkOnIDs(client *Client) error {
	if client == nil || !client.autoValidateChecksums {
		return nil
	}

	for _, ap := range tx.hbarAllowances {
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

	for _, ap := range tx.tokenAllowances {
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

	for _, ap := range tx.nftAllowances {
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

func (tx *AccountAllowanceAdjustTransaction) build() *services.TransactionBody {
	return &services.TransactionBody{}
}

func (tx *AccountAllowanceAdjustTransaction) buildScheduled() (*services.SchedulableTransactionBody, error) {
	return &services.SchedulableTransactionBody{}, nil
}

func (this *AccountAllowanceAdjustTransaction) _ConstructScheduleProtobuf() (*services.SchedulableTransactionBody, error) {
	return this.buildScheduled()
}
