package hedera

import (
	"fmt"
	protobuf "github.com/golang/protobuf/proto"
	"github.com/hashgraph/hedera-protobufs-go/services"

	"time"
)

type FeeSchedule struct {
	TransactionFeeSchedules []TransactionFeeSchedule
	ExpirationTime          *time.Time
}

func newFeeSchedule() FeeSchedule {
	return FeeSchedule{
		TransactionFeeSchedules: nil,
		ExpirationTime:          nil,
	}
}

func feeScheduleFromProtobuf(feeSchedule *services.FeeSchedule) (FeeSchedule, error) {
	if feeSchedule == nil {
		return FeeSchedule{}, errParameterNull
	}

	txFeeSchedules := make([]TransactionFeeSchedule, 0)
	for _, txFeeSchedule := range feeSchedule.GetTransactionFeeSchedule() {
		txFeeScheduleFromProto, err := transactionFeeScheduleFromProtobuf(txFeeSchedule)
		if err != nil {
			return FeeSchedule{}, err
		}
		txFeeSchedules = append(txFeeSchedules, txFeeScheduleFromProto)
	}

	var expiry time.Time
	if feeSchedule.ExpiryTime != nil {
		expiry = time.Unix(feeSchedule.ExpiryTime.Seconds, 0)
	}

	return FeeSchedule{
		TransactionFeeSchedules: txFeeSchedules,
		ExpirationTime:          &expiry,
	}, nil
}

func (feeSchedule FeeSchedule) toProtobuf() *services.FeeSchedule {
	txFeeSchedules := make([]*services.TransactionFeeSchedule, 0)
	for _, txFeeSchedule := range feeSchedule.TransactionFeeSchedules {
		txFeeSchedules = append(txFeeSchedules, txFeeSchedule.toProtobuf())
	}

	var expiry services.TimestampSeconds
	if feeSchedule.ExpirationTime != nil {
		expiry = services.TimestampSeconds{Seconds: feeSchedule.ExpirationTime.Unix()}
	}

	return &services.FeeSchedule{
		TransactionFeeSchedule: txFeeSchedules,
		ExpiryTime:             &expiry,
	}
}

func (feeSchedule FeeSchedule) ToBytes() []byte {
	data, err := protobuf.Marshal(feeSchedule.toProtobuf())
	if err != nil {
		return make([]byte, 0)
	}

	return data
}

func FeeScheduleFromBytes(data []byte) (FeeSchedule, error) {
	if data == nil {
		return FeeSchedule{}, errByteArrayNull
	}
	pb := services.FeeSchedule{}
	err := protobuf.Unmarshal(data, &pb)
	if err != nil {
		return FeeSchedule{}, err
	}

	info, err := feeScheduleFromProtobuf(&pb)
	if err != nil {
		return FeeSchedule{}, err
	}

	return info, nil
}

func (feeSchedule FeeSchedule) String() string {
	array := "\n"
	for _, i := range feeSchedule.TransactionFeeSchedules {
		array = array + i.String() + "\n"
	}
	return fmt.Sprintf("TransactionFeeSchedules: %s, Expiration: %s", array, feeSchedule.ExpirationTime)
}
