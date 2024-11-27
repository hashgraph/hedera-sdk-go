package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"time"

	"github.com/hiero-ledger/hiero-sdk-go/v2/proto/services"
)

// Delete a file or smart contract - can only be done with a Hiero admin.
// When it is deleted, it immediately disappears from the system as seen by the user,
// but is still stored internally until the expiration time, at which time it
// is truly and permanently deleted.
// Until that time, it can be undeleted by the Hiero admin.
// When a smart contract is deleted, the cryptocurrency account within it continues
// to exist, and is not affected by the expiration time here.
type SystemDeleteTransaction struct {
	*Transaction[*SystemDeleteTransaction]
	contractID     *ContractID
	fileID         *FileID
	expirationTime *time.Time
}

// NewSystemDeleteTransaction creates a SystemDeleteTransaction transaction which can be
// used to construct and execute a System Delete Transaction.
func NewSystemDeleteTransaction() *SystemDeleteTransaction {
	tx := &SystemDeleteTransaction{}
	tx.Transaction = _NewTransaction(tx)
	tx._SetDefaultMaxTransactionFee(NewHbar(2))

	return tx
}

func _SystemDeleteTransactionFromProtobuf(tx Transaction[*SystemDeleteTransaction], pb *services.TransactionBody) SystemDeleteTransaction {
	expiration := time.Date(
		time.Now().Year(), time.Now().Month(), time.Now().Day(),
		time.Now().Hour(), time.Now().Minute(),
		int(pb.GetSystemDelete().ExpirationTime.Seconds), time.Now().Nanosecond(), time.Now().Location(),
	)
	systemDeleteTransaction := SystemDeleteTransaction{
		contractID:     _ContractIDFromProtobuf(pb.GetSystemDelete().GetContractID()),
		fileID:         _FileIDFromProtobuf(pb.GetSystemDelete().GetFileID()),
		expirationTime: &expiration,
	}

	tx.childTransaction = &systemDeleteTransaction
	systemDeleteTransaction.Transaction = &tx
	return systemDeleteTransaction
}

// SetExpirationTime sets the time at which this transaction will expire.
func (tx *SystemDeleteTransaction) SetExpirationTime(expiration time.Time) *SystemDeleteTransaction {
	tx._RequireNotFrozen()
	tx.expirationTime = &expiration
	return tx
}

// GetExpirationTime returns the time at which this transaction will expire.
func (tx *SystemDeleteTransaction) GetExpirationTime() int64 {
	if tx.expirationTime != nil {
		return tx.expirationTime.Unix()
	}

	return 0
}

// SetContractID sets the ContractID of the contract which will be deleted.
func (tx *SystemDeleteTransaction) SetContractID(contractID ContractID) *SystemDeleteTransaction {
	tx._RequireNotFrozen()
	tx.contractID = &contractID
	return tx
}

// GetContractID returns the ContractID of the contract which will be deleted.
func (tx *SystemDeleteTransaction) GetContractID() ContractID {
	if tx.contractID == nil {
		return ContractID{}
	}

	return *tx.contractID
}

// SetFileID sets the FileID of the file which will be deleted.
func (tx *SystemDeleteTransaction) SetFileID(fileID FileID) *SystemDeleteTransaction {
	tx._RequireNotFrozen()
	tx.fileID = &fileID
	return tx
}

// GetFileID returns the FileID of the file which will be deleted.
func (tx *SystemDeleteTransaction) GetFileID() FileID {
	if tx.fileID == nil {
		return FileID{}
	}

	return *tx.fileID
}

// ----------- Overridden functions ----------------

func (tx SystemDeleteTransaction) getName() string {
	return "SystemDeleteTransaction"
}

func (tx SystemDeleteTransaction) validateNetworkOnIDs(client *Client) error {
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

func (tx SystemDeleteTransaction) build() *services.TransactionBody {
	return &services.TransactionBody{
		TransactionFee:           tx.transactionFee,
		Memo:                     tx.Transaction.memo,
		TransactionValidDuration: _DurationToProtobuf(tx.GetTransactionValidDuration()),
		TransactionID:            tx.transactionID._ToProtobuf(),
		Data: &services.TransactionBody_SystemDelete{
			SystemDelete: tx.buildProtoBody(),
		},
	}
}

func (tx SystemDeleteTransaction) buildScheduled() (*services.SchedulableTransactionBody, error) {
	return &services.SchedulableTransactionBody{
		TransactionFee: tx.transactionFee,
		Memo:           tx.Transaction.memo,
		Data: &services.SchedulableTransactionBody_SystemDelete{
			SystemDelete: tx.buildProtoBody(),
		},
	}, nil
}

func (tx SystemDeleteTransaction) buildProtoBody() *services.SystemDeleteTransactionBody {
	body := &services.SystemDeleteTransactionBody{}

	if tx.expirationTime != nil {
		body.ExpirationTime = &services.TimestampSeconds{
			Seconds: tx.expirationTime.Unix(),
		}
	}

	if tx.contractID != nil {
		body.Id = &services.SystemDeleteTransactionBody_ContractID{
			ContractID: tx.contractID._ToProtobuf(),
		}
	}

	if tx.fileID != nil {
		body.Id = &services.SystemDeleteTransactionBody_FileID{
			FileID: tx.fileID._ToProtobuf(),
		}
	}

	return body
}

func (tx SystemDeleteTransaction) getMethod(channel *_Channel) _Method {
	if channel._GetContract() == nil {
		return _Method{
			transaction: channel._GetFile().SystemDelete,
		}
	}

	return _Method{
		transaction: channel._GetContract().SystemDelete,
	}
}

func (tx SystemDeleteTransaction) constructScheduleProtobuf() (*services.SchedulableTransactionBody, error) {
	return tx.buildScheduled()
}

func (tx SystemDeleteTransaction) getBaseTransaction() *Transaction[TransactionInterface] {
	return castFromConcreteToBaseTransaction(tx.Transaction, &tx)
}
