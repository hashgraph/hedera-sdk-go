package hedera

import (
	"fmt"
	"time"

	"github.com/hashgraph/hedera-sdk-go/proto"
)

type TransactionID struct {
	AccountID  AccountID
	ValidStart time.Time
}

func NewTransactionID(accountID AccountID) TransactionID {
	now := time.Now().Add(-10 * time.Second)

	return TransactionID{accountID, now}
}

func TransactionIDWithValidStart(accountID AccountID, validStart time.Time) TransactionID {
	return TransactionID{accountID, validStart}
}

func (id TransactionID) GetReceipt(client *Client) (TransactionReceipt, error) {
	// TODO: Return an error on an exceptional receipt status and ensure that
	// 		 TransactionReceiptQuery does not return error on exceptional receipt status

	return NewTransactionReceiptQuery().
		SetTransactionID(id).
		Execute(client)
}

func (id TransactionID) GetRecord(client *Client) (TransactionRecord, error) {
	_, err := id.GetReceipt(client)
	if err != nil {
		return TransactionRecord{}, err
	}

	return NewTransactionRecordQuery().SetTransactionID(id).
		Execute(client)
}

func (id TransactionID) String() string {
	pb := timeToProto(id.ValidStart)
	return fmt.Sprintf("%v@%v.%v", id.AccountID, pb.Seconds, pb.Nanos)
}

func (id TransactionID) toProto() *proto.TransactionID {
	return &proto.TransactionID{
		TransactionValidStart: timeToProto(id.ValidStart),
		AccountID:             id.AccountID.toProto(),
	}
}

func transactionIDFromProto(pb *proto.TransactionID) TransactionID {
	validStart := timeFromProto(pb.TransactionValidStart)
	accountID := accountIDFromProto(pb.AccountID)

	return TransactionID{accountID, validStart}
}
