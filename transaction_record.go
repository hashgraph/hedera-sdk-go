package hedera

import (
	"github.com/hashgraph/hedera-sdk-go/proto"
	"time"
)

type TransactionRecord struct {
	Receipt            TransactionReceipt
	Hash               []byte
	ConsensusTimestamp time.Time
	TransactionID      TransactionID
	Memo               string
	TransactionFee     uint64
	TransferList       []AccountAmount
}

func transactionRecordFromProto(pb *proto.TransactionRecord) TransactionRecord {
	var transferList = []AccountAmount{}

	for _, element := range pb.TransferList.AccountAmounts {
		transferList = append(transferList, accountAmountFromProto(element))
	}

	return TransactionRecord{
		Receipt:            transactionReceiptFromProto(pb.Receipt),
		Hash:               pb.TransactionHash,
		ConsensusTimestamp: timeFromProto(pb.ConsensusTimestamp),
		TransactionID:      transactionIDFromProto(pb.TransactionID),
		Memo:               pb.Memo,
		TransactionFee:     pb.TransactionFee,
		TransferList:       transferList,
	}
}
