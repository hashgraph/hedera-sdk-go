package hedera

import (
	protobuf "github.com/golang/protobuf/proto"
	"github.com/hashgraph/hedera-sdk-go/v2/proto"
)

type ScheduleInfo struct {
	ScheduleID       ScheduleID
	CreatorAccountID AccountID
	PayerAccountID   AccountID
	TransactionBody  []byte
	Signatories      *KeyList
	AdminKey         Key
}

func scheduleInfoFromProtobuf(pb *proto.ScheduleInfo) ScheduleInfo {
	var adminKey Key
	if pb.AdminKey != nil {
		adminKey, _ = keyFromProtobuf(pb.AdminKey)
	}

	var signers KeyList
	if pb.Signatories != nil {
		signers, _ = keyListFromProtobuf(pb.Signatories)
	}

	return ScheduleInfo{
		ScheduleID:       scheduleIDFromProtobuf(pb.ScheduleID),
		CreatorAccountID: accountIDFromProtobuf(pb.CreatorAccountID),
		PayerAccountID:   accountIDFromProtobuf(pb.PayerAccountID),
		TransactionBody:  pb.TransactionBody,
		Signatories:      &signers,
		AdminKey:         adminKey,
	}
}

func (scheduleInfo *ScheduleInfo) toProtobuf() *proto.ScheduleInfo {
	var adminKey *proto.Key
	if scheduleInfo.AdminKey != nil {
		adminKey = scheduleInfo.AdminKey.toProtoKey()
	}

	var signers *proto.KeyList
	if scheduleInfo.Signatories != nil {
		signers = scheduleInfo.Signatories.toProtoKeyList()
	}

	return &proto.ScheduleInfo{
		ScheduleID:       scheduleInfo.ScheduleID.toProtobuf(),
		CreatorAccountID: scheduleInfo.CreatorAccountID.toProtobuf(),
		PayerAccountID:   scheduleInfo.PayerAccountID.toProtobuf(),
		TransactionBody:  scheduleInfo.TransactionBody,
		Signatories:      signers,
		AdminKey:         adminKey,
	}
}

func (scheduleInfo *ScheduleInfo) GetTransaction() (interface{}, error) {
	signedBytes, err := protobuf.Marshal(&proto.SignedTransaction{
		BodyBytes: scheduleInfo.TransactionBody,
		SigMap:    &proto.SignatureMap{SigPair: make([]*proto.SignaturePair, 0)},
	})
	list := proto.TransactionList{TransactionList: []*proto.Transaction{{SignedTransactionBytes: signedBytes}}}
	listBytes, err := protobuf.Marshal(&list)
	if err != nil {
		return Transaction{}, err
	}
	tx, err := TransactionFromBytes(listBytes)
	return tx, err
}
