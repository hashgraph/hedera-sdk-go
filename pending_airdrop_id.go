package hedera

import "github.com/hashgraph/hedera-protobufs-go/services"

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

type PendingAirdropId struct {
	sender   *AccountID
	receiver *AccountID
	tokenID  *TokenID
	nftID    *NftID
}

func (pendingAirdropId *PendingAirdropId) NewPendingAirdropId() *PendingAirdropId {
	return &PendingAirdropId{}
}

func (pendingAirdropId *PendingAirdropId) GetSender() *AccountID {
	return pendingAirdropId.sender
}

func (pendingAirdropId *PendingAirdropId) SetSender(sender AccountID) *PendingAirdropId {
	pendingAirdropId.sender = &sender
	return pendingAirdropId
}

func (pendingAirdropId *PendingAirdropId) GetReceiver() *AccountID {
	return pendingAirdropId.receiver
}

func (pendingAirdropId *PendingAirdropId) SetReceiver(receiver AccountID) *PendingAirdropId {
	pendingAirdropId.receiver = &receiver
	return pendingAirdropId
}

func (pendingAirdropId *PendingAirdropId) GetTokenID() *TokenID {
	return pendingAirdropId.tokenID
}

func (pendingAirdropId *PendingAirdropId) SetTokenID(tokenID TokenID) *PendingAirdropId {
	pendingAirdropId.tokenID = &tokenID
	return pendingAirdropId
}

func (pendingAirdropId *PendingAirdropId) GetNftID() *NftID {
	return pendingAirdropId.nftID
}

func (pendingAirdropId *PendingAirdropId) SetNftID(nftID NftID) *PendingAirdropId {
	pendingAirdropId.nftID = &nftID
	return pendingAirdropId
}

func _PendingAirdropIdFromProtobuf(pb *services.PendingAirdropId) *PendingAirdropId {
	if pb.GetFungibleTokenType() != nil {
		return &PendingAirdropId{
			sender:   _AccountIDFromProtobuf(pb.GetSenderId()),
			receiver: _AccountIDFromProtobuf(pb.GetReceiverId()),
			tokenID:  _TokenIDFromProtobuf(pb.GetFungibleTokenType()),
		}
	} else {
		nftID := _NftIDFromProtobuf(pb.GetNonFungibleToken())
		return &PendingAirdropId{
			sender:   _AccountIDFromProtobuf(pb.GetSenderId()),
			receiver: _AccountIDFromProtobuf(pb.GetReceiverId()),
			nftID:    &nftID,
		}
	}
}

func (pendingAirdropId *PendingAirdropId) _ToProtobuf() *services.PendingAirdropId {
	pb := &services.PendingAirdropId{
		SenderId:   pendingAirdropId.sender._ToProtobuf(),
		ReceiverId: pendingAirdropId.receiver._ToProtobuf(),
	}

	if pendingAirdropId.tokenID != nil {
		pb.TokenReference = &services.PendingAirdropId_FungibleTokenType{
			FungibleTokenType: pendingAirdropId.tokenID._ToProtobuf(),
		}
	} else {
		pb.TokenReference = &services.PendingAirdropId_NonFungibleToken{
			NonFungibleToken: pendingAirdropId.nftID._ToProtobuf(),
		}
	}
	return pb
}
