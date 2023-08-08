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
	protobuf "google.golang.org/protobuf/proto"
)

// TokenNftTransfer is the information about a NFT transfer
type TokenNftTransfer struct {
	SenderAccountID   AccountID
	ReceiverAccountID AccountID
	SerialNumber      int64
	IsApproved        bool
}

type _TokenNftTransfers struct {
	tokenNftTransfers []*TokenNftTransfer
}

func _NftTransferFromProtobuf(pb *services.NftTransfer) TokenNftTransfer {
	if pb == nil {
		return TokenNftTransfer{}
	}

	senderAccountID := AccountID{}
	if pb.SenderAccountID != nil {
		senderAccountID = *_AccountIDFromProtobuf(pb.SenderAccountID)
	}

	receiverAccountID := AccountID{}
	if pb.ReceiverAccountID != nil {
		receiverAccountID = *_AccountIDFromProtobuf(pb.ReceiverAccountID)
	}

	return TokenNftTransfer{
		SenderAccountID:   senderAccountID,
		ReceiverAccountID: receiverAccountID,
		SerialNumber:      pb.SerialNumber,
		IsApproved:        pb.IsApproval,
	}
}

func (transfer *TokenNftTransfer) _ToProtobuf() *services.NftTransfer {
	return &services.NftTransfer{
		SenderAccountID:   transfer.SenderAccountID._ToProtobuf(),
		ReceiverAccountID: transfer.ReceiverAccountID._ToProtobuf(),
		SerialNumber:      transfer.SerialNumber,
		IsApproval:        transfer.IsApproved,
	}
}

// ToBytes returns the byte representation of the TokenNftTransfer
func (transfer TokenNftTransfer) ToBytes() []byte {
	data, err := protobuf.Marshal(transfer._ToProtobuf())
	if err != nil {
		return make([]byte, 0)
	}

	return data
}

// TokenNftTransfersFromBytes returns the TokenNftTransfer from a raw protobuf bytes representation
func NftTransferFromBytes(data []byte) (TokenNftTransfer, error) {
	if data == nil {
		return TokenNftTransfer{}, errByteArrayNull
	}
	pb := services.NftTransfer{}
	err := protobuf.Unmarshal(data, &pb)
	if err != nil {
		return TokenNftTransfer{}, err
	}

	return _NftTransferFromProtobuf(&pb), nil
}

func (tokenNftTransfers *_TokenNftTransfers) Len() int {
	return len(tokenNftTransfers.tokenNftTransfers)
}
func (tokenNftTransfers *_TokenNftTransfers) Swap(i, j int) {
	tokenNftTransfers.tokenNftTransfers[i], tokenNftTransfers.tokenNftTransfers[j] = tokenNftTransfers.tokenNftTransfers[j], tokenNftTransfers.tokenNftTransfers[i]
}

func (tokenNftTransfers *_TokenNftTransfers) Less(i, j int) bool {
	if tokenNftTransfers.tokenNftTransfers[i].SenderAccountID.Compare(tokenNftTransfers.tokenNftTransfers[j].SenderAccountID) < 0 { //nolint
		return true
	} else if tokenNftTransfers.tokenNftTransfers[i].SenderAccountID.Compare(tokenNftTransfers.tokenNftTransfers[j].SenderAccountID) > 0 {
		return false
	}

	if tokenNftTransfers.tokenNftTransfers[i].ReceiverAccountID.Compare(tokenNftTransfers.tokenNftTransfers[j].ReceiverAccountID) < 0 { //nolint
		return true
	} else if tokenNftTransfers.tokenNftTransfers[i].ReceiverAccountID.Compare(tokenNftTransfers.tokenNftTransfers[j].ReceiverAccountID) > 0 {
		return false
	}

	if tokenNftTransfers.tokenNftTransfers[i].SerialNumber < tokenNftTransfers.tokenNftTransfers[j].SerialNumber { //nolint
		return true
	} else { //nolint
		return false
	}
}
