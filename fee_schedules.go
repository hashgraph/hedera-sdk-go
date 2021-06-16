package hedera

import (
	"fmt"
	protobuf "github.com/golang/protobuf/proto"
	"github.com/hashgraph/hedera-sdk-go/v2/proto"
)

type FeeSchedules struct {
	current *FeeSchedule
	next    *FeeSchedule
}

func newFeeSchedules() FeeSchedules {
	return FeeSchedules{
		current: nil,
		next:    nil,
	}
}

func feeSchedulesFromProtobuf(feeSchedules *proto.CurrentAndNextFeeSchedule) (FeeSchedules, error) {
	if feeSchedules == nil {
		return FeeSchedules{}, errParameterNull
	}

	var current FeeSchedule
	var err error
	if feeSchedules.CurrentFeeSchedule != nil {
		current, err = feeScheduleFromProtobuf(feeSchedules.GetCurrentFeeSchedule())
		if err != nil {
			return FeeSchedules{}, err
		}
	}

	var next FeeSchedule
	if feeSchedules.NextFeeSchedule != nil {
		next, err = feeScheduleFromProtobuf(feeSchedules.GetNextFeeSchedule())
		if err != nil {
			return FeeSchedules{}, err
		}
	}

	return FeeSchedules{
		current: &current,
		next:    &next,
	}, nil
}

func (feeSchedules FeeSchedules) toProtobuf() *proto.CurrentAndNextFeeSchedule {
	var current *proto.FeeSchedule
	if feeSchedules.current != nil {
		current = feeSchedules.current.toProtobuf()
	}

	var next *proto.FeeSchedule
	if feeSchedules.next != nil {
		next = feeSchedules.next.toProtobuf()
	}

	return &proto.CurrentAndNextFeeSchedule{
		CurrentFeeSchedule: current,
		NextFeeSchedule:    next,
	}
}

func (feeSchedules FeeSchedules) ToBytes() []byte {
	data, err := protobuf.Marshal(feeSchedules.toProtobuf())
	if err != nil {
		return make([]byte, 0)
	}

	return data
}

func FeeSchedulesFromBytes(data []byte) (FeeSchedules, error) {
	if data == nil {
		return FeeSchedules{}, errByteArrayNull
	}
	pb := proto.CurrentAndNextFeeSchedule{}
	err := protobuf.Unmarshal(data, &pb)
	if err != nil {
		return FeeSchedules{}, err
	}

	info, err := feeSchedulesFromProtobuf(&pb)
	if err != nil {
		return FeeSchedules{}, err
	}

	return info, nil
}

func (feeSchedules FeeSchedules) String() string {
	return fmt.Sprintf("Current: %s, Next: %s", feeSchedules.current.String(), feeSchedules.next.String())
}
