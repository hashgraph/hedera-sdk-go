//go:build all || unit
// +build all unit

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
	"testing"

	"github.com/hashgraph/hedera-protobufs-go/services"
	"github.com/stretchr/testify/assert"
)

func TestTokenAssociationFromProtobuf(t *testing.T) {
	t.Parallel()

	var pbAssociation *services.TokenAssociation
	var association TokenAssociation

	association = tokenAssociationFromProtobuf(pbAssociation)
	assert.Equal(t, TokenAssociation{}, association)

	pbAssociation = &services.TokenAssociation{}
	association = tokenAssociationFromProtobuf(pbAssociation)
	assert.Equal(t, TokenAssociation{}, association)

	pbAssociation = &services.TokenAssociation{
		TokenId:   &services.TokenID{ShardNum: 0, RealmNum: 0, TokenNum: 3},
		AccountId: &services.AccountID{ShardNum: 0, RealmNum: 0, Account: &services.AccountID_AccountNum{AccountNum: 6}},
	}

	association = tokenAssociationFromProtobuf(pbAssociation)
	assert.Equal(t, TokenAssociation{
		TokenID:   &TokenID{Shard: 0, Realm: 0, Token: 3},
		AccountID: &AccountID{Shard: 0, Realm: 0, Account: 6},
	}, association)
}

func TestTokenAssociationToProtobuf(t *testing.T) {
	t.Parallel()

	var association TokenAssociation
	var pbAssociation *services.TokenAssociation

	association = TokenAssociation{}
	pbAssociation = association.toProtobuf()
	assert.Equal(t, &services.TokenAssociation{}, pbAssociation)

	association = TokenAssociation{
		TokenID:   &TokenID{Shard: 0, Realm: 0, Token: 3},
		AccountID: &AccountID{Shard: 0, Realm: 0, Account: 6},
	}
	pbAssociation = association.toProtobuf()
	assert.Equal(t, &services.TokenAssociation{
		TokenId:   &services.TokenID{ShardNum: 0, RealmNum: 0, TokenNum: 3},
		AccountId: &services.AccountID{ShardNum: 0, RealmNum: 0, Account: &services.AccountID_AccountNum{AccountNum: 6}},
	}, pbAssociation)

}

func TestTokenAssociationToAndFromBytes(t *testing.T) {
	t.Parallel()

	association := TokenAssociation{
		TokenID:   &TokenID{Shard: 0, Realm: 0, Token: 3},
		AccountID: &AccountID{Shard: 0, Realm: 0, Account: 6},
	}

	bytes := association.ToBytes()
	fromBytes, err := TokenAssociationFromBytes(bytes)
	assert.NoError(t, err)
	assert.Equal(t, association, fromBytes)

	association = TokenAssociation{}
	bytes = association.ToBytes()

	// Test empty bytes
	data, err := TokenAssociationFromBytes(bytes)
	assert.Nil(t, err)
	assert.Equal(t, TokenAssociation{}, data)

	// Test invalid bytes
	_, err = TokenAssociationFromBytes([]byte{0x00})
	assert.Error(t, err)

	// Test nil bytes
	_, err = TokenAssociationFromBytes(nil)
	assert.Error(t, err)
}
