//go:build all || unit
// +build all unit

package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"testing"

	"github.com/hiero-ledger/hiero-sdk-go/v2/proto/services"
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
