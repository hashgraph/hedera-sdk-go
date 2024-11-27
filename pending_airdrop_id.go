package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"fmt"

	"github.com/hiero-ledger/hiero-sdk-go/v2/proto/services"
)

/**
 * A unique, composite, identifier for a pending airdrop.
 *
 * Each pending airdrop SHALL be uniquely identified by a PendingAirdropId.
 * A PendingAirdropId SHALL be recorded when created and MUST be provided in any transaction
 * that would modify that pending airdrop (such as a `claimAirdrop` or `cancelAirdrop`).
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

// GetSender returns the AccountID of the sender
func (pendingAirdropId *PendingAirdropId) GetSender() *AccountID {
	return pendingAirdropId.sender
}

// SetSender sets the AccountID of the sender
func (pendingAirdropId *PendingAirdropId) SetSender(sender AccountID) *PendingAirdropId {
	pendingAirdropId.sender = &sender
	return pendingAirdropId
}

// GetReceiver returns the AccountID of the receiver
func (pendingAirdropId *PendingAirdropId) GetReceiver() *AccountID {
	return pendingAirdropId.receiver
}

// SetReceiver sets the AccountID of the receiver
func (pendingAirdropId *PendingAirdropId) SetReceiver(receiver AccountID) *PendingAirdropId {
	pendingAirdropId.receiver = &receiver
	return pendingAirdropId
}

// GetTokenID returns the TokenID of the pending airdrop
func (pendingAirdropId *PendingAirdropId) GetTokenID() *TokenID {
	return pendingAirdropId.tokenID
}

// SetTokenID sets the TokenID of the pending airdrop
func (pendingAirdropId *PendingAirdropId) SetTokenID(tokenID TokenID) *PendingAirdropId {
	pendingAirdropId.tokenID = &tokenID
	return pendingAirdropId
}

// GetNftID returns the NftID of the pending airdrop
func (pendingAirdropId *PendingAirdropId) GetNftID() *NftID {
	return pendingAirdropId.nftID
}

// SetNftID sets the NftID of the pending airdrop
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
	pb := &services.PendingAirdropId{}

	if pendingAirdropId.sender != nil {
		pb.SenderId = pendingAirdropId.sender._ToProtobuf()
	}

	if pendingAirdropId.receiver != nil {
		pb.ReceiverId = pendingAirdropId.receiver._ToProtobuf()
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

func (pendingAirdropId *PendingAirdropId) String() string {
	const nilString = "nil"
	var senderStr, receiverStr, tokenIDStr, nftIDStr string

	if pendingAirdropId.sender != nil {
		senderStr = pendingAirdropId.sender.String()
	} else {
		senderStr = nilString
	}

	if pendingAirdropId.receiver != nil {
		receiverStr = pendingAirdropId.receiver.String()
	} else {
		receiverStr = nilString
	}

	if pendingAirdropId.tokenID != nil {
		tokenIDStr = pendingAirdropId.tokenID.String()
	} else {
		tokenIDStr = nilString
	}

	if pendingAirdropId.nftID != nil {
		nftIDStr = pendingAirdropId.nftID.String()
	} else {
		nftIDStr = nilString
	}

	return fmt.Sprintf("Sender: %s, Receiver: %s, TokenID: %s, NftID: %s", senderStr, receiverStr, tokenIDStr, nftIDStr)
}
