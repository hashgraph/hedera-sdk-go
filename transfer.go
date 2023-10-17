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

// Transfer is a transfer of hbars or tokens from one account to another
type Transfer struct {
	AccountID  AccountID
	Amount     Hbar
	IsApproved bool
}

func _TransferFromProtobuf(pb *services.AccountAmount) Transfer {
	if pb == nil {
		return Transfer{}
	}

	accountID := AccountID{}
	if pb.AccountID != nil {
		accountID = *_AccountIDFromProtobuf(pb.AccountID)
	}

	return Transfer{
		AccountID:  accountID,
		Amount:     HbarFromTinybar(pb.Amount),
		IsApproved: pb.GetIsApproval(),
	}
}

func (transfer Transfer) _ToProtobuf() *services.TransferList { // nolint
	var ammounts = make([]*services.AccountAmount, 0)
	ammounts = append(ammounts, &services.AccountAmount{
		AccountID: transfer.AccountID._ToProtobuf(),
		Amount:    transfer.Amount.AsTinybar(),
	})

	return &services.TransferList{
		AccountAmounts: ammounts,
	}
}
