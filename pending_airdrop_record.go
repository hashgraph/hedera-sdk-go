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

import (
	"fmt"

	"github.com/hashgraph/hedera-sdk-go/v2/proto/services"
)

type PendingAirdropRecord struct {
	pendingAirdropId     PendingAirdropId
	pendingAirdropAmount uint64
}

func (pendingAirdropRecord *PendingAirdropRecord) GetPendingAirdropId() PendingAirdropId {
	return pendingAirdropRecord.pendingAirdropId
}

func (pendingAirdropRecord *PendingAirdropRecord) GetPendingAirdropAmount() uint64 {
	return pendingAirdropRecord.pendingAirdropAmount
}

func _PendingAirdropRecordFromProtobuf(pb *services.PendingAirdropRecord) PendingAirdropRecord {
	return PendingAirdropRecord{
		pendingAirdropId:     *(_PendingAirdropIdFromProtobuf(pb.GetPendingAirdropId())),
		pendingAirdropAmount: pb.PendingAirdropValue.GetAmount(),
	}
}

func (pendingAirdropRecord *PendingAirdropRecord) _ToProtobuf() *services.PendingAirdropRecord {
	return &services.PendingAirdropRecord{
		PendingAirdropId: pendingAirdropRecord.pendingAirdropId._ToProtobuf(),
		PendingAirdropValue: &services.PendingAirdropValue{
			Amount: pendingAirdropRecord.pendingAirdropAmount,
		},
	}
}

func (pendingAirdropRecord *PendingAirdropRecord) String() string {
	return fmt.Sprintf("PendingAirdropRecord{PendingAirdropId: %s, PendingAirdropAmount: %d}", pendingAirdropRecord.pendingAirdropId.String(), pendingAirdropRecord.pendingAirdropAmount)
}
