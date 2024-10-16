package hedera

/*-
 *
 * Hedera Go SDK
 *
 * Copyright (C) 2020 - 2024 Hedera Hashgraph, LLC
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

import "github.com/hashgraph/hedera-sdk-go/v2/proto/services"

type _HbarTransfer struct {
	accountID  *AccountID
	Amount     Hbar
	IsApproved bool
}

func _HbarTransferFromProtobuf(pb []*services.AccountAmount) []*_HbarTransfer {
	result := make([]*_HbarTransfer, 0)
	for _, acc := range pb {
		result = append(result, &_HbarTransfer{
			accountID:  _AccountIDFromProtobuf(acc.AccountID),
			Amount:     HbarFromTinybar(acc.Amount),
			IsApproved: acc.GetIsApproval(),
		})
	}

	return result
}

func (transfer *_HbarTransfer) _ToProtobuf() *services.AccountAmount { //nolint
	var account *services.AccountID
	if transfer.accountID != nil {
		account = transfer.accountID._ToProtobuf()
	}

	return &services.AccountAmount{
		AccountID:  account,
		Amount:     transfer.Amount.AsTinybar(),
		IsApproval: transfer.IsApproved,
	}
}
