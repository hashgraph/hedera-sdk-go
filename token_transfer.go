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

// TokenTransfer is a token transfer record.
type TokenTransfer struct {
	AccountID  AccountID
	Amount     int64
	IsApproved bool
}

type _TokenTransfers struct {
	transfers []TokenTransfer
}

// NewTokenTransfer creates a TokenTransfer with the given accountID and amount
func NewTokenTransfer(accountID AccountID, amount int64) TokenTransfer {
	return TokenTransfer{
		AccountID: accountID,
		Amount:    amount,
	}
}

func _TokenTransferFromProtobuf(pb *services.AccountAmount) TokenTransfer {
	if pb == nil {
		return TokenTransfer{}
	}

	accountID := AccountID{}
	if pb.AccountID != nil {
		accountID = *_AccountIDFromProtobuf(pb.AccountID)
	}

	return TokenTransfer{
		AccountID:  accountID,
		Amount:     pb.Amount,
		IsApproved: pb.IsApproval,
	}
}

func (transfer *TokenTransfer) _ToProtobuf() *services.AccountAmount {
	return &services.AccountAmount{
		AccountID:  transfer.AccountID._ToProtobuf(),
		Amount:     transfer.Amount,
		IsApproval: transfer.IsApproved,
	}
}

// ToBytes returns a protobuf encoded version of the TokenTransfer
func (transfer TokenTransfer) ToBytes() []byte {
	data, err := protobuf.Marshal(transfer._ToProtobuf())
	if err != nil {
		return make([]byte, 0)
	}

	return data
}

// TokenTransferFromBytes returns a TokenTransfer struct from a protobuf encoded byte array
func TokenTransferFromBytes(data []byte) (TokenTransfer, error) {
	if data == nil {
		return TokenTransfer{}, errByteArrayNull
	}
	pb := services.AccountAmount{}
	err := protobuf.Unmarshal(data, &pb)
	if err != nil {
		return TokenTransfer{}, err
	}

	return _TokenTransferFromProtobuf(&pb), nil
}

func (transfer TokenTransfer) String() string {
	return fmt.Sprintf("accountID: %s, amount: %d", transfer.AccountID.String(), transfer.Amount)
}

func (transfers _TokenTransfers) Len() int {
	return len(transfers.transfers)
}
func (transfers _TokenTransfers) Swap(i, j int) {
	transfers.transfers[i], transfers.transfers[j] = transfers.transfers[j], transfers.transfers[i]
}

func (transfers _TokenTransfers) Less(i, j int) bool {
	if transfers.transfers[i].AccountID.Compare(transfers.transfers[j].AccountID) < 0 { //nolint
		return true
	}

	return false
}
