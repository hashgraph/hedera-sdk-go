package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"fmt"

	"github.com/hiero-ledger/hiero-sdk-go/v2/proto/services"
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
