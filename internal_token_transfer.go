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
	"github.com/hashgraph/hedera-protobufs-go/services"
)

type _TokenTransfer struct {
	Transfers        []*_HbarTransfer
	ExpectedDecimals *uint32
}

func _TokenTransferPrivateFromProtobuf(pb *services.TokenTransferList) *_TokenTransfer {
	if pb == nil {
		return &_TokenTransfer{}
	}

	var decimals *uint32
	if pb.ExpectedDecimals != nil {
		temp := pb.ExpectedDecimals.GetValue()
		decimals = &temp
	}

	return &_TokenTransfer{
		Transfers:        _HbarTransferFromProtobuf(pb.Transfers),
		ExpectedDecimals: decimals,
	}
}

func (transfer *_TokenTransfer) _ToProtobuf() []*services.AccountAmount {
	transfers := make([]*services.AccountAmount, 0)
	for _, t := range transfer.Transfers {
		transfers = append(transfers, &services.AccountAmount{
			AccountID:  t.accountID._ToProtobuf(),
			Amount:     t.Amount.AsTinybar(),
			IsApproval: t.IsApproved,
		})
	}
	return transfers
}
