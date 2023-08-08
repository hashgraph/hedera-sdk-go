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

type TransactionFeeSchedule struct {
	RequestType RequestType
	// Deprecated: use Fees
	FeeData *FeeData
	Fees    []*FeeData
}

func _TransactionFeeScheduleFromProtobuf(txFeeSchedule *services.TransactionFeeSchedule) (TransactionFeeSchedule, error) {
	if txFeeSchedule == nil {
		return TransactionFeeSchedule{}, errParameterNull
	}

	feeData := make([]*FeeData, 0)

	for _, d := range txFeeSchedule.GetFees() {
		temp, err := _FeeDataFromProtobuf(d)
		if err != nil {
			return TransactionFeeSchedule{}, err
		}
		feeData = append(feeData, &temp)
	}

	return TransactionFeeSchedule{
		RequestType: RequestType(txFeeSchedule.GetHederaFunctionality()),
		Fees:        feeData,
	}, nil
}

func (txFeeSchedule TransactionFeeSchedule) _ToProtobuf() *services.TransactionFeeSchedule {
	feeData := make([]*services.FeeData, 0)
	if txFeeSchedule.Fees != nil {
		for _, data := range txFeeSchedule.Fees {
			feeData = append(feeData, data._ToProtobuf())
		}
	}

	var singleFee *services.FeeData
	if txFeeSchedule.FeeData != nil {
		singleFee = txFeeSchedule.FeeData._ToProtobuf()
	}

	return &services.TransactionFeeSchedule{
		HederaFunctionality: services.HederaFunctionality(txFeeSchedule.RequestType),
		Fees:                feeData,
		FeeData:             singleFee,
	}
}

func (txFeeSchedule TransactionFeeSchedule) ToBytes() []byte {
	data, err := protobuf.Marshal(txFeeSchedule._ToProtobuf())
	if err != nil {
		return make([]byte, 0)
	}

	return data
}

func (txFeeSchedule TransactionFeeSchedule) String() string {
	str := ""
	for _, dat := range txFeeSchedule.Fees {
		str = str + dat.String() + ", "
	}

	if txFeeSchedule.FeeData != nil {
		return fmt.Sprintf("RequestType: %s, Feedata: %s", txFeeSchedule.RequestType.String(), txFeeSchedule.FeeData.String())
	}

	return fmt.Sprintf("RequestType: %s, Feedata: %s", txFeeSchedule.RequestType.String(), str)
}
