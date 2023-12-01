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
	tx := AccountAllowanceDeleteTransaction{
		transaction: _NewTransaction(),
	}
	tx.e = &tx
	tx._SetDefaultMaxTransactionFee(NewHbar(2))

	return &tx
}

func _AccountAllowanceDeleteTransactionFromProtobuf(transaction transaction, pb *services.TransactionBody) *AccountAllowanceDeleteTransaction {
	nftWipe := make([]*TokenNftAllowance, 0)

	for _, ap := range pb.GetCryptoDeleteAllowance().GetNftAllowances() {
		temp := _TokenNftWipeAllowanceProtobuf(ap)
		nftWipe = append(nftWipe, &temp)
	}

	resultTx := &AccountAllowanceDeleteTransaction{
		transaction: transaction,
		nftWipe:     nftWipe,
	}
	resultTx.e = resultTx
	return resultTx
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

// ---- Required Interfaces ---- //

// Sign uses the provided privateKey to sign the transaction.
func (tx *AccountAllowanceDeleteTransaction) Sign(
	privateKey PrivateKey,
) *AccountAllowanceDeleteTransaction {
	tx.transaction.Sign(privateKey)
	return tx
}

// SignWithOperator signs the transaction with client's operator privateKey.
func (tx *AccountAllowanceDeleteTransaction) SignWithOperator(
	client *Client,
) (*AccountAllowanceDeleteTransaction, error) {
	// If the transaction is not signed by the _Operator, we need
	// to sign the transaction with the _Operator
	_, err := tx.transaction.SignWithOperator(client)
	return tx, err
}

// SignWith executes the TransactionSigner and adds the resulting signature data to the transaction's signature map
// with the publicKey as the map key.
func (tx *AccountAllowanceDeleteTransaction) SignWith(
	publicKey PublicKey,
	signer TransactionSigner,
) *AccountAllowanceDeleteTransaction {
	tx.transaction.SignWith(publicKey, signer)
	return tx
}

// AddSignature adds a signature to the transaction.
func (tx *AccountAllowanceDeleteTransaction) AddSignature(publicKey PublicKey, signature []byte) *AccountAllowanceDeleteTransaction {
	tx.transaction.AddSignature(publicKey, signature)
	return tx
}

// When execution is attempted, a single attempt will timeout when this deadline is reached. (The SDK may subsequently retry the execution.)
func (tx *AccountAllowanceDeleteTransaction) SetGrpcDeadline(deadline *time.Duration) *AccountAllowanceDeleteTransaction {
	tx.transaction.SetGrpcDeadline(deadline)
	return tx
}

func (tx *AccountAllowanceDeleteTransaction) Freeze() (*AccountAllowanceDeleteTransaction, error) {
	_, err := tx.transaction.Freeze()
	return tx, err
}

func (tx *AccountAllowanceDeleteTransaction) FreezeWith(client *Client) (*AccountAllowanceDeleteTransaction, error) {
	_, err := tx.transaction.FreezeWith(client)
	return tx, err
}

// SetMaxTransactionFee sets the max transaction fee for tx AccountAllowanceDeleteTransaction.
func (tx *AccountAllowanceDeleteTransaction) SetMaxTransactionFee(fee Hbar) *AccountAllowanceDeleteTransaction {
	tx._RequireNotFrozen()
	tx.transaction.SetMaxTransactionFee(fee)
	return tx
}

// SetRegenerateTransactionID sets if transaction IDs should be regenerated when `TRANSACTION_EXPIRED` is received
func (tx *AccountAllowanceDeleteTransaction) SetRegenerateTransactionID(regenerateTransactionID bool) *AccountAllowanceDeleteTransaction {
	tx._RequireNotFrozen()
	tx.transaction.SetRegenerateTransactionID(regenerateTransactionID)
	return tx
}

// SetTransactionMemo sets the memo for tx AccountAllowanceDeleteTransaction.
func (tx *AccountAllowanceDeleteTransaction) SetTransactionMemo(memo string) *AccountAllowanceDeleteTransaction {
	tx.transaction.SetTransactionMemo(memo)
	return tx
}

// SetTransactionValidDuration sets the valid duration for tx AccountAllowanceDeleteTransaction.
func (tx *AccountAllowanceDeleteTransaction) SetTransactionValidDuration(duration time.Duration) *AccountAllowanceDeleteTransaction {
	tx.transaction.SetTransactionValidDuration(duration)
	return tx
}

// SetTransactionID sets the TransactionID for tx AccountAllowanceDeleteTransaction.
func (tx *AccountAllowanceDeleteTransaction) SetTransactionID(transactionID TransactionID) *AccountAllowanceDeleteTransaction {
	tx._RequireNotFrozen()

	tx.transaction.SetTransactionID(transactionID)
	return tx
}

// SetNodeAccountIDs sets the _Node AccountID for tx AccountAllowanceDeleteTransaction.
func (tx *AccountAllowanceDeleteTransaction) SetNodeAccountIDs(nodeID []AccountID) *AccountAllowanceDeleteTransaction {
	tx._RequireNotFrozen()
	tx.transaction.SetNodeAccountIDs(nodeID)
	return tx
}

// SetMaxRetry sets the max number of errors before execution will fail.
func (tx *AccountAllowanceDeleteTransaction) SetMaxRetry(count int) *AccountAllowanceDeleteTransaction {
	tx.transaction.SetMaxRetry(count)
	return tx
}

// SetMaxBackoff The maximum amount of time to wait between retries.
// Every retry attempt will increase the wait time exponentially until it reaches tx time.
func (tx *AccountAllowanceDeleteTransaction) SetMaxBackoff(max time.Duration) *AccountAllowanceDeleteTransaction {
	tx.transaction.SetMaxBackoff(max)
	return tx
}

// SetMinBackoff sets the min back off for tx AccountAllowanceDeleteTransaction.
func (tx *AccountAllowanceDeleteTransaction) SetMinBackoff(min time.Duration) *AccountAllowanceDeleteTransaction {
	tx.transaction.SetMinBackoff(min)
	return tx
}

// ----------- overriden functions ----------------
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
		Memo:                     tx.transaction.memo,
		Data: &services.TransactionBody_CryptoDeleteAllowance{
			CryptoDeleteAllowance: tx.buildProtoBody(),
		},
	}
}

func (tx *AccountAllowanceDeleteTransaction) buildScheduled() (*services.SchedulableTransactionBody, error) {
	return &services.SchedulableTransactionBody{
		TransactionFee: tx.transactionFee,
		Memo:           tx.transaction.memo,
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
