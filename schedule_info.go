package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"time"

	"github.com/hiero-ledger/hiero-sdk-go/v2/proto/services"
)

type ScheduleInfo struct {
	ScheduleID       ScheduleID
	CreatorAccountID AccountID
	PayerAccountID   AccountID
	ExecutedAt       *time.Time
	DeletedAt        *time.Time
	ExpirationTime   time.Time
	Signatories      *KeyList
	// Deprecated: Use ScheduleInfo.Signatories instead
	Signers                  *KeyList
	AdminKey                 Key
	Memo                     string
	ScheduledTransactionID   *TransactionID
	scheduledTransactionBody *services.SchedulableTransactionBody
	LedgerID                 LedgerID
	WaitForExpiry            bool
}

func _ScheduleInfoFromProtobuf(pb *services.ScheduleInfo) ScheduleInfo {
	if pb == nil {
		return ScheduleInfo{}
	}
	var adminKey Key
	if pb.AdminKey != nil {
		adminKey, _ = _KeyFromProtobuf(pb.AdminKey)
	}

	var signatories KeyList
	if pb.Signers != nil {
		signatories, _ = _KeyListFromProtobuf(pb.Signers)
	}

	var scheduledTransactionID TransactionID
	if pb.ScheduledTransactionID != nil {
		scheduledTransactionID = _TransactionIDFromProtobuf(pb.ScheduledTransactionID)
	}

	var executed *time.Time
	var deleted *time.Time
	switch t := pb.Data.(type) {
	case *services.ScheduleInfo_ExecutionTime:
		temp := _TimeFromProtobuf(t.ExecutionTime)
		executed = &temp
	case *services.ScheduleInfo_DeletionTime:
		temp := _TimeFromProtobuf(t.DeletionTime)
		deleted = &temp
	}

	creatorAccountID := AccountID{}
	if pb.CreatorAccountID != nil {
		creatorAccountID = *_AccountIDFromProtobuf(pb.CreatorAccountID)
	}

	payerAccountID := AccountID{}
	if pb.PayerAccountID != nil {
		payerAccountID = *_AccountIDFromProtobuf(pb.PayerAccountID)
	}

	scheduleID := ScheduleID{}
	if pb.ScheduleID != nil {
		scheduleID = *_ScheduleIDFromProtobuf(pb.ScheduleID)
	}

	return ScheduleInfo{
		ScheduleID:               scheduleID,
		CreatorAccountID:         creatorAccountID,
		PayerAccountID:           payerAccountID,
		ExecutedAt:               executed,
		DeletedAt:                deleted,
		ExpirationTime:           _TimeFromProtobuf(pb.ExpirationTime),
		Signatories:              &signatories,
		Signers:                  &signatories,
		AdminKey:                 adminKey,
		Memo:                     pb.Memo,
		ScheduledTransactionID:   &scheduledTransactionID,
		scheduledTransactionBody: pb.ScheduledTransactionBody,
		LedgerID:                 LedgerID{pb.LedgerId},
		WaitForExpiry:            pb.WaitForExpiry,
	}
}

func (scheduleInfo *ScheduleInfo) _ToProtobuf() *services.ScheduleInfo { // nolint
	var adminKey *services.Key
	if scheduleInfo.AdminKey != nil {
		adminKey = scheduleInfo.AdminKey._ToProtoKey()
	}

	var signatories *services.KeyList
	if scheduleInfo.Signatories != nil {
		signatories = scheduleInfo.Signatories._ToProtoKeyList()
	} else if scheduleInfo.Signers != nil {
		signatories = scheduleInfo.Signers._ToProtoKeyList()
	}

	info := &services.ScheduleInfo{
		ScheduleID:               scheduleInfo.ScheduleID._ToProtobuf(),
		ExpirationTime:           _TimeToProtobuf(scheduleInfo.ExpirationTime),
		ScheduledTransactionBody: scheduleInfo.scheduledTransactionBody,
		Memo:                     scheduleInfo.Memo,
		AdminKey:                 adminKey,
		Signers:                  signatories,
		CreatorAccountID:         scheduleInfo.CreatorAccountID._ToProtobuf(),
		PayerAccountID:           scheduleInfo.PayerAccountID._ToProtobuf(),
		ScheduledTransactionID:   scheduleInfo.ScheduledTransactionID._ToProtobuf(),
		LedgerId:                 scheduleInfo.LedgerID.ToBytes(),
		WaitForExpiry:            scheduleInfo.WaitForExpiry,
	}

	if scheduleInfo.ExecutedAt != nil {
		info.Data = &services.ScheduleInfo_DeletionTime{
			DeletionTime: _TimeToProtobuf(*scheduleInfo.DeletedAt),
		}
	} else if scheduleInfo.DeletedAt != nil {
		info.Data = &services.ScheduleInfo_ExecutionTime{
			ExecutionTime: _TimeToProtobuf(*scheduleInfo.ExecutedAt),
		}
	}

	return info
}

// GetScheduledTransaction returns the scheduled transaction associated with this schedule
func (scheduleInfo *ScheduleInfo) GetScheduledTransaction() (TransactionInterface, error) { // nolint
	pb := scheduleInfo.scheduledTransactionBody
	return transactionFromScheduledTransaction(pb)
}
