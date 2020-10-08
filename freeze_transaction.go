package hedera

import (
	"github.com/hashgraph/hedera-sdk-go/proto"
	"time"
)

type FreezeTransaction struct {
	Transaction
	pb *proto.FreezeTransactionBody
}

func NewFreezeTransaction() *FreezeTransaction {
	pb := &proto.FreezeTransactionBody{}

	transaction := FreezeTransaction{
		pb:          pb,
		Transaction: newTransaction(),
	}

	return &transaction
}

func (transaction *FreezeTransaction) SetStartTime(startTime time.Time) *FreezeTransaction {
	transaction.pb.StartHour = int32(startTime.Hour())
	transaction.pb.StartMin = int32(startTime.Minute())
	return transaction
}

func (transaction *FreezeTransaction) GetStartTime() time.Time {
	t1 := time.Date(
		time.Now().Year(), time.Now().Month(), time.Now().Day(),
		int(transaction.pb.StartHour), int(transaction.pb.StartMin),
		0, time.Now().Nanosecond(), time.Now().Location(),
		)
	return t1
}

func (transaction *FreezeTransaction) SetEndTime(endTime time.Time) *FreezeTransaction {
	transaction.pb.StartHour = int32(endTime.Hour())
	transaction.pb.StartMin = int32(endTime.Minute())
	return transaction
}

func (transaction *FreezeTransaction) GetEndTime() time.Time {
	t1 := time.Date(
		time.Now().Year(), time.Now().Month(), time.Now().Day(),
		int(transaction.pb.EndHour), int(transaction.pb.EndMin),
		0, time.Now().Nanosecond(), time.Now().Location(),
		)
	return t1
}

