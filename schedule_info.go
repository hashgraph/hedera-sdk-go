package hedera

import (
	"github.com/hashgraph/hedera-sdk-go/proto"
)

type ScheduleInfo struct {
	ScheduleID       ScheduleID
	CreatorAccountID AccountID
	PayerAccountID   AccountID
	TransactionBody  []byte
	Signers          []PublicKey
	AdminKey         PublicKey
}

func scheduleInfoFromProtobuf(pb *proto.ScheduleInfo) ScheduleInfo {
	var adminKey PublicKey
	if pb.AdminKey != nil {
		adminKey, _ = publicKeyFromProto(pb.AdminKey)
	}

	var signers []PublicKey
	if pb.Signatories != nil {
		signers, _ = publicKeyListFromProto(pb.Signatories)
	}

	return ScheduleInfo{
		ScheduleID:       scheduleIDFromProto(pb.ScheduleID),
		CreatorAccountID: accountIDFromProto(pb.CreatorAccountID),
		PayerAccountID:   accountIDFromProto(pb.PayerAccountID),
		TransactionBody:  pb.TransactionBody,
		Signers:          signers,
		AdminKey:         adminKey,
	}
}

func (scheduleInfo *ScheduleInfo) toProtobuf() *proto.ScheduleInfo {
	var adminKey *proto.Key
	if scheduleInfo.AdminKey != nil {
		adminKey = scheduleInfo.AdminKey.toProto()
	}

	var temp KeyList
	if scheduleInfo.Signers != nil {
		temp.AddAll(scheduleInfo.Signers)
	}

	var signers *proto.KeyList
	if temp.keys != nil {
		signers = &proto.KeyList{Keys: temp.keys}
	}

	return &proto.ScheduleInfo{
		ScheduleID:       scheduleInfo.ScheduleID.toProto(),
		CreatorAccountID: scheduleInfo.CreatorAccountID.toProto(),
		PayerAccountID:   scheduleInfo.PayerAccountID.toProto(),
		TransactionBody:  scheduleInfo.TransactionBody,
		Signatories:      signers,
		AdminKey:         adminKey,
	}
}
