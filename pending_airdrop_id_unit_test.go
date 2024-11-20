//go:build all || unit
// +build all unit

package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"testing"

	"github.com/hiero-ledger/hiero-sdk-go/v2/proto/services"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPendingAirdropId_SettersAndGetters(t *testing.T) {
	t.Parallel()

	sender := AccountID{Account: 1}
	receiver := AccountID{Account: 2}
	tokenID := TokenID{Token: 3}
	nftID := NftID{TokenID: tokenID, SerialNumber: 1}

	pendingAirdropId := &PendingAirdropId{}
	pendingAirdropId.SetSender(sender)
	pendingAirdropId.SetReceiver(receiver)
	pendingAirdropId.SetTokenID(tokenID)
	pendingAirdropId.SetNftID(nftID)

	assert.Equal(t, &sender, pendingAirdropId.GetSender())
	assert.Equal(t, &receiver, pendingAirdropId.GetReceiver())
	assert.Equal(t, &tokenID, pendingAirdropId.GetTokenID())
	assert.Equal(t, &nftID, pendingAirdropId.GetNftID())
}

func TestPendingAirdropId_FromProtobuf(t *testing.T) {
	t.Parallel()

	sender := AccountID{Account: 1}
	receiver := AccountID{Account: 2}
	tokenID := TokenID{Token: 3}
	nftID := NftID{TokenID: tokenID, SerialNumber: 1}

	pb := &services.PendingAirdropId{
		SenderId:   sender._ToProtobuf(),
		ReceiverId: receiver._ToProtobuf(),
		TokenReference: &services.PendingAirdropId_FungibleTokenType{
			FungibleTokenType: tokenID._ToProtobuf(),
		},
	}

	pendingAirdropId := _PendingAirdropIdFromProtobuf(pb)
	require.NotNil(t, pendingAirdropId)
	assert.Equal(t, &sender, pendingAirdropId.GetSender())
	assert.Equal(t, &receiver, pendingAirdropId.GetReceiver())
	assert.Equal(t, &tokenID, pendingAirdropId.GetTokenID())
	assert.Nil(t, pendingAirdropId.GetNftID())

	pb = &services.PendingAirdropId{
		SenderId:   sender._ToProtobuf(),
		ReceiverId: receiver._ToProtobuf(),
		TokenReference: &services.PendingAirdropId_NonFungibleToken{
			NonFungibleToken: nftID._ToProtobuf(),
		},
	}

	pendingAirdropId = _PendingAirdropIdFromProtobuf(pb)
	require.NotNil(t, pendingAirdropId)
	assert.Equal(t, &sender, pendingAirdropId.GetSender())
	assert.Equal(t, &receiver, pendingAirdropId.GetReceiver())
	assert.Nil(t, pendingAirdropId.GetTokenID())
	assert.Equal(t, &nftID, pendingAirdropId.GetNftID())
}

func TestPendingAirdropId_ToProtobuf(t *testing.T) {
	t.Parallel()

	sender := AccountID{Account: 1}
	receiver := AccountID{Account: 2}
	tokenID := TokenID{Token: 3}
	nftID := NftID{TokenID: tokenID, SerialNumber: 1}

	pendingAirdropId := &PendingAirdropId{
		sender:   &sender,
		receiver: &receiver,
		tokenID:  &tokenID,
	}

	pb := pendingAirdropId._ToProtobuf()
	require.NotNil(t, pb)
	assert.Equal(t, sender._ToProtobuf(), pb.GetSenderId())
	assert.Equal(t, receiver._ToProtobuf(), pb.GetReceiverId())
	assert.Equal(t, tokenID._ToProtobuf(), pb.GetFungibleTokenType())
	assert.Nil(t, pb.GetNonFungibleToken())

	pendingAirdropId = &PendingAirdropId{
		sender:   &sender,
		receiver: &receiver,
		nftID:    &nftID,
	}

	pb = pendingAirdropId._ToProtobuf()
	require.NotNil(t, pb)
	assert.Equal(t, sender._ToProtobuf(), pb.GetSenderId())
	assert.Equal(t, receiver._ToProtobuf(), pb.GetReceiverId())
	assert.Nil(t, pb.GetFungibleTokenType())
	assert.Equal(t, nftID._ToProtobuf(), pb.GetNonFungibleToken())
}
