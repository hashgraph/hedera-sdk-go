package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"github.com/hiero-ledger/hiero-sdk-go/v2/proto/services"
)

// Undelete a file or smart contract that was deleted by AdminDelete.
// Can only be done with a Hiero admin.
type SystemUndeleteTransaction struct {
	*Transaction[*SystemUndeleteTransaction]
	contractID *ContractID
	fileID     *FileID
}

// NewSystemUndeleteTransaction creates a SystemUndeleteTransaction transaction which can be
// used to construct and execute a System Undelete Transaction.
func NewSystemUndeleteTransaction() *SystemUndeleteTransaction {
	tx := &SystemUndeleteTransaction{}
	tx.Transaction = _NewTransaction(tx)
	tx._SetDefaultMaxTransactionFee(NewHbar(2))

	return tx
}

func _SystemUndeleteTransactionFromProtobuf(tx Transaction[*SystemUndeleteTransaction], pb *services.TransactionBody) SystemUndeleteTransaction {
	systemUndeleteTransaction := SystemUndeleteTransaction{
		contractID: _ContractIDFromProtobuf(pb.GetSystemUndelete().GetContractID()),
		fileID:     _FileIDFromProtobuf(pb.GetSystemUndelete().GetFileID()),
	}

	tx.childTransaction = &systemUndeleteTransaction
	systemUndeleteTransaction.Transaction = &tx
	return systemUndeleteTransaction
}

// SetContractID sets the ContractID of the contract whose deletion is being undone.
func (tx *SystemUndeleteTransaction) SetContractID(contractID ContractID) *SystemUndeleteTransaction {
	tx._RequireNotFrozen()
	tx.contractID = &contractID
	return tx
}

// GetContractID returns the ContractID of the contract whose deletion is being undone.
func (tx *SystemUndeleteTransaction) GetContractID() ContractID {
	if tx.contractID == nil {
		return ContractID{}
	}

	return *tx.contractID
}

// SetFileID sets the FileID of the file whose deletion is being undone.
func (tx *SystemUndeleteTransaction) SetFileID(fileID FileID) *SystemUndeleteTransaction {
	tx._RequireNotFrozen()
	tx.fileID = &fileID
	return tx
}

// GetFileID returns the FileID of the file whose deletion is being undone.
func (tx *SystemUndeleteTransaction) GetFileID() FileID {
	if tx.fileID == nil {
		return FileID{}
	}

	return *tx.fileID
}

// ----------- Overridden functions ----------------

func (tx SystemUndeleteTransaction) getName() string {
	return "SystemUndeleteTransaction"
}

func (tx SystemUndeleteTransaction) validateNetworkOnIDs(client *Client) error {
	if client == nil || !client.autoValidateChecksums {
		return nil
	}

	if tx.contractID != nil {
		if err := tx.contractID.ValidateChecksum(client); err != nil {
			return err
		}
	}

	if tx.fileID != nil {
		if err := tx.fileID.ValidateChecksum(client); err != nil {
			return err
		}
	}

	return nil
}

func (tx SystemUndeleteTransaction) build() *services.TransactionBody {
	return &services.TransactionBody{
		TransactionFee:           tx.transactionFee,
		Memo:                     tx.Transaction.memo,
		TransactionValidDuration: _DurationToProtobuf(tx.GetTransactionValidDuration()),
		TransactionID:            tx.transactionID._ToProtobuf(),
		Data: &services.TransactionBody_SystemUndelete{
			SystemUndelete: tx.buildProtoBody(),
		},
	}
}

func (tx SystemUndeleteTransaction) buildScheduled() (*services.SchedulableTransactionBody, error) {
	return &services.SchedulableTransactionBody{
		TransactionFee: tx.transactionFee,
		Memo:           tx.Transaction.memo,
		Data: &services.SchedulableTransactionBody_SystemUndelete{
			SystemUndelete: tx.buildProtoBody(),
		},
	}, nil
}

func (tx SystemUndeleteTransaction) buildProtoBody() *services.SystemUndeleteTransactionBody {
	body := &services.SystemUndeleteTransactionBody{}
	if tx.contractID != nil {
		body.Id = &services.SystemUndeleteTransactionBody_ContractID{
			ContractID: tx.contractID._ToProtobuf(),
		}
	}

	if tx.fileID != nil {
		body.Id = &services.SystemUndeleteTransactionBody_FileID{
			FileID: tx.fileID._ToProtobuf(),
		}
	}

	return body
}

func (tx SystemUndeleteTransaction) getMethod(channel *_Channel) _Method {
	if channel._GetContract() == nil {
		return _Method{
			transaction: channel._GetFile().SystemUndelete,
		}
	}

	return _Method{
		transaction: channel._GetContract().SystemUndelete,
	}
}

func (tx SystemUndeleteTransaction) constructScheduleProtobuf() (*services.SchedulableTransactionBody, error) {
	return tx.buildScheduled()
}

func (tx SystemUndeleteTransaction) getBaseTransaction() *Transaction[TransactionInterface] {
	return castFromConcreteToBaseTransaction(tx.Transaction, &tx)
}
