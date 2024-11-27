package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"github.com/hiero-ledger/hiero-sdk-go/v2/proto/services"
)

// ContractDeleteTransaction marks a contract as deleted and transfers its remaining hBars, if any, to a
// designated receiver. After a contract is deleted, it can no longer be called.
type ContractDeleteTransaction struct {
	*Transaction[*ContractDeleteTransaction]
	contractID        *ContractID
	transferContactID *ContractID
	transferAccountID *AccountID
	permanentRemoval  bool
}

// NewContractDeleteTransaction creates ContractDeleteTransaction which marks a contract as deleted and transfers its remaining hBars, if any, to a
// designated receiver. After a contract is deleted, it can no longer be called.
func NewContractDeleteTransaction() *ContractDeleteTransaction {
	tx := &ContractDeleteTransaction{}
	tx.Transaction = _NewTransaction(tx)
	tx._SetDefaultMaxTransactionFee(NewHbar(2))

	return tx
}

func _ContractDeleteTransactionFromProtobuf(tx Transaction[*ContractDeleteTransaction], pb *services.TransactionBody) ContractDeleteTransaction {
	contractDeleteTransaction := ContractDeleteTransaction{
		contractID:        _ContractIDFromProtobuf(pb.GetContractDeleteInstance().GetContractID()),
		transferContactID: _ContractIDFromProtobuf(pb.GetContractDeleteInstance().GetTransferContractID()),
		transferAccountID: _AccountIDFromProtobuf(pb.GetContractDeleteInstance().GetTransferAccountID()),
		permanentRemoval:  pb.GetContractDeleteInstance().GetPermanentRemoval(),
	}
	tx.childTransaction = &contractDeleteTransaction
	contractDeleteTransaction.Transaction = &tx
	return contractDeleteTransaction
}

// Sets the contract ID which will be deleted.
func (tx *ContractDeleteTransaction) SetContractID(contractID ContractID) *ContractDeleteTransaction {
	tx._RequireNotFrozen()
	tx.contractID = &contractID
	return tx
}

// Returns the contract ID which will be deleted.
func (tx *ContractDeleteTransaction) GetContractID() ContractID {
	if tx.contractID == nil {
		return ContractID{}
	}

	return *tx.contractID
}

// Sets the contract ID which will receive all remaining hbars.
func (tx *ContractDeleteTransaction) SetTransferContractID(transferContactID ContractID) *ContractDeleteTransaction {
	tx._RequireNotFrozen()
	tx.transferContactID = &transferContactID
	return tx
}

// Returns the contract ID which will receive all remaining hbars.
func (tx *ContractDeleteTransaction) GetTransferContractID() ContractID {
	if tx.transferContactID == nil {
		return ContractID{}
	}

	return *tx.transferContactID
}

// Sets the account ID which will receive all remaining hbars.
func (tx *ContractDeleteTransaction) SetTransferAccountID(accountID AccountID) *ContractDeleteTransaction {
	tx._RequireNotFrozen()
	tx.transferAccountID = &accountID

	return tx
}

// Returns the account ID which will receive all remaining hbars.
func (tx *ContractDeleteTransaction) GetTransferAccountID() AccountID {
	if tx.transferAccountID == nil {
		return AccountID{}
	}

	return *tx.transferAccountID
}

// SetPermanentRemoval
// If set to true, means this is a "synthetic" system transaction being used to
// alert mirror nodes that the contract is being permanently removed from the ledger.
// IMPORTANT: User transactions cannot set this field to true, as permanent
// removal is always managed by the ledger itself. Any ContractDeleteTransaction
// submitted to HAPI with permanent_removal=true will be rejected with precheck status
// PERMANENT_REMOVAL_REQUIRES_SYSTEM_INITIATION.
func (tx *ContractDeleteTransaction) SetPermanentRemoval(remove bool) *ContractDeleteTransaction {
	tx._RequireNotFrozen()
	tx.permanentRemoval = remove

	return tx
}

// GetPermanentRemoval returns true if this is a "synthetic" system transaction.
func (tx *ContractDeleteTransaction) GetPermanentRemoval() bool {
	return tx.permanentRemoval
}

// ----------- Overridden functions ----------------

func (tx ContractDeleteTransaction) getName() string {
	return "ContractDeleteTransaction"
}
func (tx ContractDeleteTransaction) validateNetworkOnIDs(client *Client) error {
	if client == nil || !client.autoValidateChecksums {
		return nil
	}

	if tx.contractID != nil {
		if err := tx.contractID.ValidateChecksum(client); err != nil {
			return err
		}
	}

	if tx.transferContactID != nil {
		if err := tx.transferContactID.ValidateChecksum(client); err != nil {
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

func (tx ContractDeleteTransaction) build() *services.TransactionBody {
	pb := services.TransactionBody{
		TransactionFee:           tx.transactionFee,
		Memo:                     tx.Transaction.memo,
		TransactionValidDuration: _DurationToProtobuf(tx.GetTransactionValidDuration()),
		TransactionID:            tx.transactionID._ToProtobuf(),
		Data: &services.TransactionBody_ContractDeleteInstance{
			ContractDeleteInstance: tx.buildProtoBody(),
		},
	}

	return &pb
}

func (tx ContractDeleteTransaction) buildProtoBody() *services.ContractDeleteTransactionBody {
	body := &services.ContractDeleteTransactionBody{
		PermanentRemoval: tx.permanentRemoval,
	}

	if tx.contractID != nil {
		body.ContractID = tx.contractID._ToProtobuf()
	}

	if tx.transferContactID != nil {
		body.Obtainers = &services.ContractDeleteTransactionBody_TransferContractID{
			TransferContractID: tx.transferContactID._ToProtobuf(),
		}
	}

	if tx.transferAccountID != nil {
		body.Obtainers = &services.ContractDeleteTransactionBody_TransferAccountID{
			TransferAccountID: tx.transferAccountID._ToProtobuf(),
		}
	}

	return body
}

func (tx ContractDeleteTransaction) buildScheduled() (*services.SchedulableTransactionBody, error) {
	return &services.SchedulableTransactionBody{
		TransactionFee: tx.transactionFee,
		Memo:           tx.Transaction.memo,
		Data: &services.SchedulableTransactionBody_ContractDeleteInstance{
			ContractDeleteInstance: tx.buildProtoBody(),
		},
	}, nil
}

func (tx ContractDeleteTransaction) getMethod(channel *_Channel) _Method {
	return _Method{
		transaction: channel._GetContract().DeleteContract,
	}
}

func (tx ContractDeleteTransaction) constructScheduleProtobuf() (*services.SchedulableTransactionBody, error) {
	return tx.buildScheduled()
}

func (tx ContractDeleteTransaction) getBaseTransaction() *Transaction[TransactionInterface] {
	return castFromConcreteToBaseTransaction(tx.Transaction, &tx)
}
