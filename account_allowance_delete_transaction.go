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
	"github.com/hashgraph/hedera-protobufs-go/services"
)

// AccountAllowanceDeleteTransaction
// Deletes one or more non-fungible approved allowances from an owner's account. This operation
// will remove the allowances granted to one or more specific non-fungible token serial numbers. Each owner account
// listed as wiping an allowance must sign the transaction. Hbar and fungible token allowances
// can be removed by setting the amount to zero in CryptoApproveAllowance.
type AccountAllowanceDeleteTransaction struct {
	*Transaction[*AccountAllowanceDeleteTransaction]
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
	tx := &AccountAllowanceDeleteTransaction{}
	tx.Transaction = _NewTransaction(tx)
	tx._SetDefaultMaxTransactionFee(NewHbar(2))

	return tx
}

func _AccountAllowanceDeleteTransactionFromProtobuf(pb *services.TransactionBody) *AccountAllowanceDeleteTransaction {
	nftWipe := make([]*TokenNftAllowance, 0)

	for _, ap := range pb.GetCryptoDeleteAllowance().GetNftAllowances() {
		temp := _TokenNftWipeAllowanceProtobuf(ap)
		nftWipe = append(nftWipe, &temp)
	}

	return &AccountAllowanceDeleteTransaction{
		nftWipe: nftWipe,
	}
}

// Deprecated
func (tx *AccountAllowanceDeleteTransaction) DeleteAllHbarAllowances(ownerAccountID *AccountID) *AccountAllowanceDeleteTransaction {
	tx._RequireNotFrozen()
	tx.hbarWipe = append(tx.hbarWipe, &HbarAllowance{
		OwnerAccountID: ownerAccountID,
	})

	return tx
}

// Deprecated
func (tx *AccountAllowanceDeleteTransaction) GetAllHbarDeleteAllowances() []*HbarAllowance {
	return tx.hbarWipe
}

// Deprecated
func (tx *AccountAllowanceDeleteTransaction) DeleteAllTokenAllowances(tokenID TokenID, ownerAccountID *AccountID) *AccountAllowanceDeleteTransaction {
	tx._RequireNotFrozen()
	tokenApproval := TokenAllowance{
		TokenID:        &tokenID,
		OwnerAccountID: ownerAccountID,
	}

	tx.tokenWipe = append(tx.tokenWipe, &tokenApproval)
	return tx
}

// Deprecated
func (tx *AccountAllowanceDeleteTransaction) GetAllTokenDeleteAllowances() []*TokenAllowance {
	return tx.tokenWipe
}

// DeleteAllTokenNftAllowances
// The non-fungible token allowance/allowances to remove.
func (tx *AccountAllowanceDeleteTransaction) DeleteAllTokenNftAllowances(nftID NftID, ownerAccountID *AccountID) *AccountAllowanceDeleteTransaction {
	tx._RequireNotFrozen()

	for _, t := range tx.nftWipe {
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
				return tx
			}
		}
	}

	tx.nftWipe = append(tx.nftWipe, &TokenNftAllowance{
		TokenID:        &nftID.TokenID,
		OwnerAccountID: ownerAccountID,
		SerialNumbers:  []int64{nftID.SerialNumber},
		AllSerials:     false,
	})
	return tx
}

// GetAllTokenNftDeleteAllowances
// Get the non-fungible token allowance/allowances that will be removed.
func (tx *AccountAllowanceDeleteTransaction) GetAllTokenNftDeleteAllowances() []*TokenNftAllowance {
	return tx.nftWipe
}

// ----------- Overridden functions ----------------

func (tx *AccountAllowanceDeleteTransaction) getName() string {
	return "AccountAllowanceDeleteTransaction"
}

func (tx *AccountAllowanceDeleteTransaction) validateNetworkOnIDs(client *Client) error {
	if client == nil || !client.autoValidateChecksums {
		return nil
	}

	for _, ap := range tx.nftWipe {
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

func (tx *AccountAllowanceDeleteTransaction) build() *services.TransactionBody {
	return &services.TransactionBody{
		TransactionID:            tx.transactionID._ToProtobuf(),
		TransactionFee:           tx.transactionFee,
		TransactionValidDuration: _DurationToProtobuf(tx.GetTransactionValidDuration()),
		Memo:                     tx.Transaction.memo,
		Data: &services.TransactionBody_CryptoDeleteAllowance{
			CryptoDeleteAllowance: tx.buildProtoBody(),
		},
	}
}

func (tx *AccountAllowanceDeleteTransaction) buildScheduled() (*services.SchedulableTransactionBody, error) {
	return &services.SchedulableTransactionBody{
		TransactionFee: tx.transactionFee,
		Memo:           tx.Transaction.memo,
		Data: &services.SchedulableTransactionBody_CryptoDeleteAllowance{
			CryptoDeleteAllowance: tx.buildProtoBody(),
		},
	}, nil
}

func (tx *AccountAllowanceDeleteTransaction) buildProtoBody() *services.CryptoDeleteAllowanceTransactionBody {
	body := &services.CryptoDeleteAllowanceTransactionBody{}
	nftWipe := make([]*services.NftRemoveAllowance, 0)

	for _, ap := range tx.nftWipe {
		nftWipe = append(nftWipe, ap._ToWipeProtobuf())
	}

	body.NftAllowances = nftWipe
	return body
}

func (tx *AccountAllowanceDeleteTransaction) getMethod(channel *_Channel) _Method {
	return _Method{
		transaction: channel._GetCrypto().DeleteAllowances,
	}
}

func (this *AccountAllowanceDeleteTransaction) _ConstructScheduleProtobuf() (*services.SchedulableTransactionBody, error) {
	return this.buildScheduled()
}

func (tx *AccountAllowanceDeleteTransaction) getBaseTransaction() *Transaction[TransactionInterface] {
	return castFromConcreteToBaseTransaction(tx.Transaction)
}

func (tx *AccountAllowanceDeleteTransaction) setBaseTransaction(baseTx Transaction[TransactionInterface]) {
	tx.Transaction = castFromBaseToConcreteTransaction[*AccountAllowanceDeleteTransaction](baseTx)
}
