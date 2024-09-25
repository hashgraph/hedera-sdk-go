//go:build all || unit
// +build all unit

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
	"testing"

	"github.com/hashgraph/hedera-protobufs-go/services"
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
