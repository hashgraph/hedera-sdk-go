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
	"github.com/hashgraph/hedera-sdk-go/v2/proto/services"
	protobuf "google.golang.org/protobuf/proto"
)

// _TokenNftTransfer is the information about a NFT transfer
type _TokenNftTransfer struct {
	SenderAccountID   AccountID
	ReceiverAccountID AccountID
	SerialNumber      int64
	IsApproved        bool
}

func _NftTransferFromProtobuf(pb *services.NftTransfer) _TokenNftTransfer {
	if pb == nil {
		return _TokenNftTransfer{}
	}

	senderAccountID := AccountID{}
	if pb.SenderAccountID != nil {
		senderAccountID = *_AccountIDFromProtobuf(pb.SenderAccountID)
	}

	receiverAccountID := AccountID{}
	if pb.ReceiverAccountID != nil {
		receiverAccountID = *_AccountIDFromProtobuf(pb.ReceiverAccountID)
	}

	return _TokenNftTransfer{
		SenderAccountID:   senderAccountID,
		ReceiverAccountID: receiverAccountID,
		SerialNumber:      pb.SerialNumber,
		IsApproved:        pb.IsApproval,
	}
}

func (transfer *_TokenNftTransfer) _ToProtobuf() *services.NftTransfer {
	return &services.NftTransfer{
		SenderAccountID:   transfer.SenderAccountID._ToProtobuf(),
		ReceiverAccountID: transfer.ReceiverAccountID._ToProtobuf(),
		SerialNumber:      transfer.SerialNumber,
		IsApproval:        transfer.IsApproved,
	}
}

// ToBytes returns the byte representation of the TokenNftTransfer
func (transfer _TokenNftTransfer) ToBytes() []byte {
	data, err := protobuf.Marshal(transfer._ToProtobuf())
	if err != nil {
		return make([]byte, 0)
	}

	return data
}

// TokenNftTransfersFromBytes returns the TokenNftTransfer from a raw protobuf bytes representation
func NftTransferFromBytes(data []byte) (_TokenNftTransfer, error) {
	if data == nil {
		return _TokenNftTransfer{}, errByteArrayNull
	}
	pb := services.NftTransfer{}
	err := protobuf.Unmarshal(data, &pb)
	if err != nil {
		return _TokenNftTransfer{}, err
	}

	return _NftTransferFromProtobuf(&pb), nil
}
