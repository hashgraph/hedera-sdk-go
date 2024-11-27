package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"fmt"
	"time"

	"github.com/hiero-ledger/hiero-sdk-go/v2/proto/services"
	protobuf "google.golang.org/protobuf/proto"
)

type FeeSchedule struct {
	TransactionFeeSchedules []TransactionFeeSchedule
	ExpirationTime          *time.Time
}

func _FeeScheduleFromProtobuf(feeSchedule *services.FeeSchedule) (FeeSchedule, error) {
	if feeSchedule == nil {
		return FeeSchedule{}, errParameterNull
	}

	txFeeSchedules := make([]TransactionFeeSchedule, 0)
	for _, txFeeSchedule := range feeSchedule.GetTransactionFeeSchedule() {
		txFeeScheduleFromProto, err := _TransactionFeeScheduleFromProtobuf(txFeeSchedule)
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

func (feeSchedule FeeSchedule) _ToProtobuf() *services.FeeSchedule {
	txFeeSchedules := make([]*services.TransactionFeeSchedule, 0)
	for _, txFeeSchedule := range feeSchedule.TransactionFeeSchedules {
		txFeeSchedules = append(txFeeSchedules, txFeeSchedule._ToProtobuf())
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

// ToBytes returns the byte representation of the FeeSchedule
func (feeSchedule FeeSchedule) ToBytes() []byte {
	data, err := protobuf.Marshal(feeSchedule._ToProtobuf())
	if err != nil {
		return make([]byte, 0)
	}

	return data
}

// FeeScheduleFromBytes returns a FeeSchedule from a raw protobuf byte array
func FeeScheduleFromBytes(data []byte) (FeeSchedule, error) {
	if data == nil {
		return FeeSchedule{}, errByteArrayNull
	}
	pb := services.FeeSchedule{}
	err := protobuf.Unmarshal(data, &pb)
	if err != nil {
		return FeeSchedule{}, err
	}

	info, err := _FeeScheduleFromProtobuf(&pb)
	if err != nil {
		return FeeSchedule{}, err
	}

	return info, nil
}

// String returns a string representation of the FeeSchedule
func (feeSchedule FeeSchedule) String() string {
	array := "\n"
	for _, i := range feeSchedule.TransactionFeeSchedules {
		array = array + i.String() + "\n"
	}
	return fmt.Sprintf("TransactionFeeSchedules: %s, Expiration: %s", array, feeSchedule.ExpirationTime)
}
