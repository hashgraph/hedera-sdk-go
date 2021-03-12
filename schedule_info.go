package hedera

import (
	"github.com/hashgraph/hedera-sdk-go/v2/proto"
	"time"
)

type ScheduleInfo struct {
	ScheduleID               ScheduleID
	CreatorAccountID         AccountID
	PayerAccountID           AccountID
	Executed                 time.Time
	Deleted                  time.Time
	ExpirationTime           time.Time
	ScheduledTransactionBody *SchedulableTransactionBody
	Signers                  *KeyList
	AdminKey                 Key
	Memo                     string
	ScheduledTransactionID   *TransactionID
}

func scheduleInfoFromProtobuf(pb *proto.ScheduleInfo) ScheduleInfo {
	var adminKey Key
	if pb.AdminKey != nil {
		adminKey, _ = keyFromProtobuf(pb.AdminKey)
	}

	var signers KeyList
	if pb.Signers != nil {
		signers, _ = keyListFromProtobuf(pb.Signers)
	}

	var scheduledTransactionID TransactionID
	if pb.ScheduledTransactionID != nil {
		scheduledTransactionID = transactionIDFromProtobuf(pb.ScheduledTransactionID)
	}

	var executed time.Time
	var deleted time.Time
	switch t := pb.Data.(type) {
	case *proto.ScheduleInfo_ExecutionTime:
		executed = timeFromProtobuf(t.ExecutionTime)
	case *proto.ScheduleInfo_DeletionTime:
		deleted = timeFromProtobuf(t.DeletionTime)
	}

	return ScheduleInfo{
		ScheduleID:               scheduleIDFromProtobuf(pb.ScheduleID),
		CreatorAccountID:         accountIDFromProtobuf(pb.CreatorAccountID),
		PayerAccountID:           accountIDFromProtobuf(pb.PayerAccountID),
		Executed:                 executed,
		Deleted:                  deleted,
		ExpirationTime:           timeFromProtobuf(pb.ExpirationTime),
		ScheduledTransactionBody: schedulableTransactionBodyFromProtobuf(pb.ScheduledTransactionBody),
		Signers:                  &signers,
		AdminKey:                 adminKey,
		Memo:                     pb.Memo,
		ScheduledTransactionID:   &scheduledTransactionID,
	}
}

func (scheduleInfo *ScheduleInfo) toProtobuf() *proto.ScheduleInfo {
	var adminKey *proto.Key
	if scheduleInfo.AdminKey != nil {
		adminKey = scheduleInfo.AdminKey.toProtoKey()
	}

	var signers *proto.KeyList
	if scheduleInfo.Signers != nil {
		signers = scheduleInfo.Signers.toProtoKeyList()
	}

	info := &proto.ScheduleInfo{
		ScheduleID:               scheduleInfo.ScheduleID.toProtobuf(),
		ExpirationTime:           timeToProtobuf(scheduleInfo.ExpirationTime),
		ScheduledTransactionBody: scheduleInfo.ScheduledTransactionBody.toProtobuf(),
		Memo:                     scheduleInfo.Memo,
		AdminKey:                 adminKey,
		Signers:                  signers,
		CreatorAccountID:         scheduleInfo.CreatorAccountID.toProtobuf(),
		PayerAccountID:           scheduleInfo.PayerAccountID.toProtobuf(),
		ScheduledTransactionID:   scheduleInfo.ScheduledTransactionID.toProtobuf(),
	}

	if scheduleInfo.Executed.IsZero() {
		info.Data = &proto.ScheduleInfo_DeletionTime{
			DeletionTime: timeToProtobuf(scheduleInfo.Deleted),
		}
	} else {
		info.Data = &proto.ScheduleInfo_ExecutionTime{
			ExecutionTime: timeToProtobuf(scheduleInfo.Executed),
		}
	}

	return info
}

func (scheduleInfo *ScheduleInfo) GetTransaction() interface{} {
	return scheduleInfo.ScheduledTransactionBody.Transaction
}
