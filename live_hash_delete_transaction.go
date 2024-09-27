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
	"errors"

	"github.com/hashgraph/hedera-protobufs-go/services"
)

// LiveHashDeleteTransaction At consensus, deletes a livehash associated to the given account. The transaction must be signed
// by either the key of the owning account, or at least one of the keys associated to the livehash.
type LiveHashDeleteTransaction struct {
	*Transaction[*LiveHashDeleteTransaction]
	accountID *AccountID
	hash      []byte
}

// NewLiveHashDeleteTransaction creates LiveHashDeleteTransaction which at consensus, deletes a livehash associated to the given account.
// The transaction must be signed by either the key of the owning account, or at least one of the keys associated to the livehash.
func NewLiveHashDeleteTransaction() *LiveHashDeleteTransaction {
	tx := &LiveHashDeleteTransaction{}
	tx.Transaction = _NewTransaction(tx)
	tx._SetDefaultMaxTransactionFee(NewHbar(2))

	return tx
}

func _LiveHashDeleteTransactionFromProtobuf(tx Transaction[*LiveHashDeleteTransaction], pb *services.TransactionBody) *LiveHashDeleteTransaction {
	return &LiveHashDeleteTransaction{
		Transaction: &tx,
		accountID:   _AccountIDFromProtobuf(pb.GetCryptoDeleteLiveHash().GetAccountOfLiveHash()),
		hash:        pb.GetCryptoDeleteLiveHash().LiveHashToDelete,
	}
}

// SetHash Set the SHA-384 livehash to delete from the account
func (tx *LiveHashDeleteTransaction) SetHash(hash []byte) *LiveHashDeleteTransaction {
	tx._RequireNotFrozen()
	tx.hash = hash
	return tx
}

// GetHash returns the SHA-384 livehash to delete from the account
func (tx *LiveHashDeleteTransaction) GetHash() []byte {
	return tx.hash
}

// SetAccountID Sets the account owning the livehash
func (tx *LiveHashDeleteTransaction) SetAccountID(accountID AccountID) *LiveHashDeleteTransaction {
	tx._RequireNotFrozen()
	tx.accountID = &accountID
	return tx
}

// GetAccountID returns the account owning the livehash
func (tx *LiveHashDeleteTransaction) GetAccountID() AccountID {
	if tx.accountID == nil {
		return AccountID{}
	}

	return *tx.accountID
}

// ----------- Overridden functions ----------------

func (tx *LiveHashDeleteTransaction) getName() string {
	return "LiveHashDeleteTransaction"
}

func (tx *LiveHashDeleteTransaction) validateNetworkOnIDs(client *Client) error {
	if client == nil || !client.autoValidateChecksums {
		return nil
	}

	if tx.accountID != nil {
		if err := tx.accountID.ValidateChecksum(client); err != nil {
			return err
		}
	}

	return nil
}

func (tx *LiveHashDeleteTransaction) build() *services.TransactionBody {
	return &services.TransactionBody{
		TransactionFee:           tx.transactionFee,
		Memo:                     tx.Transaction.memo,
		TransactionValidDuration: _DurationToProtobuf(tx.GetTransactionValidDuration()),
		TransactionID:            tx.transactionID._ToProtobuf(),
		Data: &services.TransactionBody_CryptoDeleteLiveHash{
			CryptoDeleteLiveHash: tx.buildProtoBody(),
		},
	}
}

func (tx *LiveHashDeleteTransaction) buildScheduled() (*services.SchedulableTransactionBody, error) {
	return nil, errors.New("cannot schedule `LiveHashDeleteTransaction`")
}

func (tx *LiveHashDeleteTransaction) buildProtoBody() *services.CryptoDeleteLiveHashTransactionBody {
	body := &services.CryptoDeleteLiveHashTransactionBody{}

	if tx.accountID != nil {
		body.AccountOfLiveHash = tx.accountID._ToProtobuf()
	}

	if tx.hash != nil {
		body.LiveHashToDelete = tx.hash
	}

	return body
}

func (tx *LiveHashDeleteTransaction) getMethod(channel *_Channel) _Method {
	return _Method{
		transaction: channel._GetCrypto().DeleteLiveHash,
	}
}

func (tx *LiveHashDeleteTransaction) _ConstructScheduleProtobuf() (*services.SchedulableTransactionBody, error) {
	return tx.buildScheduled()
}

func (tx *LiveHashDeleteTransaction) getBaseTransaction() *Transaction[TransactionInterface] {
	return castFromConcreteToBaseTransaction[*LiveHashDeleteTransaction](tx.Transaction)
}
