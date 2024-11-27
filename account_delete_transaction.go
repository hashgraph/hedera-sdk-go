package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"github.com/hiero-ledger/hiero-sdk-go/v2/proto/services"
)

// AccountDeleteTransaction
// Mark an account as deleted, moving all its current hbars to another account. It will remain in
// the ledger, marked as deleted, until it expires. Transfers into it a deleted account fail. But a
// deleted account can still have its expiration extended in the normal way.
type AccountDeleteTransaction struct {
	*Transaction[*AccountDeleteTransaction]
	transferAccountID *AccountID
	deleteAccountID   *AccountID
}

func _AccountDeleteTransactionFromProtobuf(tx Transaction[*AccountDeleteTransaction], pb *services.TransactionBody) AccountDeleteTransaction {
	accountDeleteTransaction := AccountDeleteTransaction{
		transferAccountID: _AccountIDFromProtobuf(pb.GetCryptoDelete().GetTransferAccountID()),
		deleteAccountID:   _AccountIDFromProtobuf(pb.GetCryptoDelete().GetDeleteAccountID()),
	}
	tx.childTransaction = &accountDeleteTransaction
	accountDeleteTransaction.Transaction = &tx
	return accountDeleteTransaction
}

// NewAccountDeleteTransaction creates AccountDeleteTransaction which marks an account as deleted, moving all its current hbars to another account. It will remain in
// the ledger, marked as deleted, until it expires. Transfers into it a deleted account fail. But a
// deleted account can still have its expiration extended in the normal way.
func NewAccountDeleteTransaction() *AccountDeleteTransaction {
	tx := &AccountDeleteTransaction{}
	tx.Transaction = _NewTransaction(tx)
	tx._SetDefaultMaxTransactionFee(NewHbar(2))

	return tx
}

// SetNodeAccountID sets the _Node AccountID for this AccountDeleteTransaction.
func (tx *AccountDeleteTransaction) SetAccountID(accountID AccountID) *AccountDeleteTransaction {
	tx._RequireNotFrozen()
	tx.deleteAccountID = &accountID
	return tx
}

// GetAccountID returns the AccountID which will be deleted.
func (tx *AccountDeleteTransaction) GetAccountID() AccountID {
	if tx.deleteAccountID == nil {
		return AccountID{}
	}

	return *tx.deleteAccountID
}

// SetTransferAccountID sets the AccountID which will receive all remaining hbars.
func (tx *AccountDeleteTransaction) SetTransferAccountID(transferAccountID AccountID) *AccountDeleteTransaction {
	tx._RequireNotFrozen()
	tx.transferAccountID = &transferAccountID
	return tx
}

// GetTransferAccountID returns the AccountID which will receive all remaining hbars.
func (tx *AccountDeleteTransaction) GetTransferAccountID() AccountID {
	if tx.transferAccountID == nil {
		return AccountID{}
	}

	return *tx.transferAccountID
}

// ----------- Overridden functions ----------------

func (tx AccountDeleteTransaction) getName() string {
	return "AccountDeleteTransaction"
}
func (tx AccountDeleteTransaction) validateNetworkOnIDs(client *Client) error {
	if client == nil || !client.autoValidateChecksums {
		return nil
	}

	if tx.deleteAccountID != nil {
		if err := tx.deleteAccountID.ValidateChecksum(client); err != nil {
			return err
		}
	}

	if tx.transferAccountID != nil {
		if err := tx.transferAccountID.ValidateChecksum(client); err != nil {
			return err
		}
	}

	return nil
}

func (tx AccountDeleteTransaction) build() *services.TransactionBody {
	return &services.TransactionBody{
		TransactionFee:           tx.transactionFee,
		Memo:                     tx.Transaction.memo,
		TransactionValidDuration: _DurationToProtobuf(tx.GetTransactionValidDuration()),
		TransactionID:            tx.transactionID._ToProtobuf(),
		Data: &services.TransactionBody_CryptoDelete{
			CryptoDelete: tx.buildProtoBody(),
		},
	}
}

func (tx AccountDeleteTransaction) buildScheduled() (*services.SchedulableTransactionBody, error) {
	return &services.SchedulableTransactionBody{
		TransactionFee: tx.transactionFee,
		Memo:           tx.Transaction.memo,
		Data: &services.SchedulableTransactionBody_CryptoDelete{
			CryptoDelete: tx.buildProtoBody(),
		},
	}, nil
}

func (tx AccountDeleteTransaction) buildProtoBody() *services.CryptoDeleteTransactionBody {
	body := &services.CryptoDeleteTransactionBody{}

	if tx.transferAccountID != nil {
		body.TransferAccountID = tx.transferAccountID._ToProtobuf()
	}

	if tx.deleteAccountID != nil {
		body.DeleteAccountID = tx.deleteAccountID._ToProtobuf()
	}

	return body
}

func (tx AccountDeleteTransaction) getMethod(channel *_Channel) _Method {
	return _Method{
		transaction: channel._GetCrypto().ApproveAllowances,
	}
}

func (tx AccountDeleteTransaction) constructScheduleProtobuf() (*services.SchedulableTransactionBody, error) {
	return tx.buildScheduled()
}

func (tx AccountDeleteTransaction) getBaseTransaction() *Transaction[TransactionInterface] {
	return castFromConcreteToBaseTransaction(tx.Transaction, tx)
}
