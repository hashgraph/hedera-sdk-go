package hedera

import (
	"github.com/hashgraph/hedera-sdk-go/v2/proto"
)

type ScheduleInfo struct {
	ScheduleID       ScheduleID
	CreatorAccountID AccountID
	PayerAccountID   AccountID
	TransactionBody  []byte
	Signers          *KeyList
	AdminKey         Key
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

	return ScheduleInfo{
		ScheduleID:       scheduleIDFromProtobuf(pb.ScheduleID),
		CreatorAccountID: accountIDFromProtobuf(pb.CreatorAccountID),
		PayerAccountID:   accountIDFromProtobuf(pb.PayerAccountID),
		TransactionBody:  pb.TransactionBody,
		Signers:          &signers,
		AdminKey:         adminKey,
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

	return &proto.ScheduleInfo{
		ScheduleID:       scheduleInfo.ScheduleID.toProtobuf(),
		CreatorAccountID: scheduleInfo.CreatorAccountID.toProtobuf(),
		PayerAccountID:   scheduleInfo.PayerAccountID.toProtobuf(),
		TransactionBody:  scheduleInfo.TransactionBody,
		Signers:          signers,
		AdminKey:         adminKey,
	}
}
