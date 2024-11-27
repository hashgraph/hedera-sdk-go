package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"github.com/hiero-ledger/hiero-sdk-go/v2/proto/services"

	"time"
)

type FreezeTransaction struct {
	*Transaction[*FreezeTransaction]
	startTime  time.Time
	endTime    time.Time
	fileID     *FileID
	fileHash   []byte
	freezeType FreezeType
}

func NewFreezeTransaction() *FreezeTransaction {
	tx := &FreezeTransaction{}
	tx.Transaction = _NewTransaction(tx)

	tx._SetDefaultMaxTransactionFee(NewHbar(2))

	return tx
}

func _FreezeTransactionFromProtobuf(tx Transaction[*FreezeTransaction], pb *services.TransactionBody) FreezeTransaction {
	startTime := time.Date(
		time.Now().Year(), time.Now().Month(), time.Now().Day(),
		int(pb.GetFreeze().GetStartHour()), int(pb.GetFreeze().GetStartMin()), // nolint
		0, time.Now().Nanosecond(), time.Now().Location(),
	)

	endTime := time.Date(
		time.Now().Year(), time.Now().Month(), time.Now().Day(),
		int(pb.GetFreeze().GetEndHour()), int(pb.GetFreeze().GetEndMin()), // nolint
		0, time.Now().Nanosecond(), time.Now().Location(),
	)

	freezeTransaction := FreezeTransaction{
		startTime: startTime,
		endTime:   endTime,
		fileID:    _FileIDFromProtobuf(pb.GetFreeze().GetUpdateFile()),
		fileHash:  pb.GetFreeze().FileHash,
	}
	tx.childTransaction = &freezeTransaction
	freezeTransaction.Transaction = &tx
	return freezeTransaction
}

func (tx *FreezeTransaction) SetStartTime(startTime time.Time) *FreezeTransaction {
	tx._RequireNotFrozen()
	tx.startTime = startTime
	return tx
}

func (tx *FreezeTransaction) GetStartTime() time.Time {
	return tx.startTime
}

// Deprecated
func (tx *FreezeTransaction) SetEndTime(endTime time.Time) *FreezeTransaction {
	tx._RequireNotFrozen()
	tx.endTime = endTime
	return tx
}

// Deprecated
func (tx *FreezeTransaction) GetEndTime() time.Time {
	return tx.endTime
}

func (tx *FreezeTransaction) SetFileID(id FileID) *FreezeTransaction {
	tx._RequireNotFrozen()
	tx.fileID = &id
	return tx
}

func (tx *FreezeTransaction) GetFileID() *FileID {
	return tx.fileID
}

func (tx *FreezeTransaction) SetFreezeType(freezeType FreezeType) *FreezeTransaction {
	tx._RequireNotFrozen()
	tx.freezeType = freezeType
	return tx
}

func (tx *FreezeTransaction) GetFreezeType() FreezeType {
	return tx.freezeType
}

func (tx *FreezeTransaction) SetFileHash(hash []byte) *FreezeTransaction {
	tx._RequireNotFrozen()
	tx.fileHash = hash
	return tx
}

func (tx *FreezeTransaction) GetFileHash() []byte {
	return tx.fileHash
}

// ----------- Overridden functions ----------------

func (tx FreezeTransaction) getName() string {
	return "FreezeTransaction"
}
func (tx FreezeTransaction) build() *services.TransactionBody {
	return &services.TransactionBody{
		TransactionFee:           tx.transactionFee,
		Memo:                     tx.Transaction.memo,
		TransactionValidDuration: _DurationToProtobuf(tx.GetTransactionValidDuration()),
		TransactionID:            tx.transactionID._ToProtobuf(),
		Data: &services.TransactionBody_Freeze{
			Freeze: tx.buildProtoBody(),
		},
	}
}
func (tx FreezeTransaction) buildScheduled() (*services.SchedulableTransactionBody, error) {
	return &services.SchedulableTransactionBody{
		TransactionFee: tx.transactionFee,
		Memo:           tx.Transaction.memo,
		Data: &services.SchedulableTransactionBody_Freeze{
			Freeze: tx.buildProtoBody(),
		},
	}, nil
}
func (tx FreezeTransaction) buildProtoBody() *services.FreezeTransactionBody {
	body := &services.FreezeTransactionBody{
		FileHash:   tx.fileHash,
		StartTime:  _TimeToProtobuf(tx.startTime),
		FreezeType: services.FreezeType(tx.freezeType),
	}

	if tx.fileID != nil {
		body.UpdateFile = tx.fileID._ToProtobuf()
	}

	return body
}
func (tx FreezeTransaction) getMethod(channel *_Channel) _Method {
	return _Method{
		transaction: channel._GetFreeze().Freeze,
	}
}
func (tx FreezeTransaction) constructScheduleProtobuf() (*services.SchedulableTransactionBody, error) {
	return tx.buildScheduled()
}

func (tx FreezeTransaction) validateNetworkOnIDs(client *Client) error {
	return nil
}

func (tx FreezeTransaction) getBaseTransaction() *Transaction[TransactionInterface] {
	return castFromConcreteToBaseTransaction(tx.Transaction, &tx)
}
