package hedera

/*-
 *
 * Hedera Go SDK
 *
 * Copyright (C) 2020 - 2023 Hedera Hashgraph, LLC
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

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

// ToBytes returns the byte representation of the FeeSchedules
func (feeSchedules FeeSchedules) ToBytes() []byte {
	data, err := protobuf.Marshal(feeSchedules._ToProtobuf())
	if err != nil {
		return make([]byte, 0)
	}

	return data
}

// FeeSchedulesFromBytes returns a FeeSchedules object from a raw byte array
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

// String returns a string representation of the FeeSchedules
func (feeSchedules FeeSchedules) String() string {
	return fmt.Sprintf("Current: %s, Next: %s", feeSchedules.current.String(), feeSchedules.next.String())
}
