package hedera

import (
	"fmt"

	"github.com/hashgraph/hedera-protobufs-go/services"
	protobuf "google.golang.org/protobuf/proto"
)

type FeeSchedules struct {
	current *FeeSchedule
	next    *FeeSchedule
}

func _FeeSchedulesFromProtobuf(feeSchedules *services.CurrentAndNextFeeSchedule) (FeeSchedules, error) {
	if feeSchedules == nil {
		return FeeSchedules{}, errParameterNull
	}

	var current FeeSchedule
	var err error
	if feeSchedules.CurrentFeeSchedule != nil {
		current, err = _FeeScheduleFromProtobuf(feeSchedules.GetCurrentFeeSchedule())
		if err != nil {
			return FeeSchedules{}, err
		}
	}

	var next FeeSchedule
	if feeSchedules.NextFeeSchedule != nil {
		next, err = _FeeScheduleFromProtobuf(feeSchedules.GetNextFeeSchedule())
		if err != nil {
			return FeeSchedules{}, err
		}
	}

	return FeeSchedules{
		current: &current,
		next:    &next,
	}, nil
}

func (feeSchedules FeeSchedules) _ToProtobuf() *services.CurrentAndNextFeeSchedule {
	var current *services.FeeSchedule
	if feeSchedules.current != nil {
		current = feeSchedules.current._ToProtobuf()
	}

	var next *services.FeeSchedule
	if feeSchedules.next != nil {
		next = feeSchedules.next._ToProtobuf()
	}

	return &services.CurrentAndNextFeeSchedule{
		CurrentFeeSchedule: current,
		NextFeeSchedule:    next,
	}
}

func (feeSchedules FeeSchedules) ToBytes() []byte {
	data, err := protobuf.Marshal(feeSchedules._ToProtobuf())
	if err != nil {
		return make([]byte, 0)
	}

	return data
}

func FeeSchedulesFromBytes(data []byte) (FeeSchedules, error) {
	if data == nil {
		return FeeSchedules{}, errByteArrayNull
	}
	pb := services.CurrentAndNextFeeSchedule{}
	err := protobuf.Unmarshal(data, &pb)
	if err != nil {
		return FeeSchedules{}, err
	}

	info, err := _FeeSchedulesFromProtobuf(&pb)
	if err != nil {
		return FeeSchedules{}, err
	}

	return info, nil
}

func (feeSchedules FeeSchedules) String() string {
	return fmt.Sprintf("Current: %s, Next: %s", feeSchedules.current.String(), feeSchedules.next.String())
}
