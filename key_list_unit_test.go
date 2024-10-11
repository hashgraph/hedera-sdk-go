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

	"github.com/hashgraph/hedera-sdk-go/v2/generated/services"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"
)

// Mock Key and PublicKey structs and methods for testing
type MockKey struct {
	data string
}

func (k MockKey) _ToProtoKey() *services.Key {
	return &services.Key{Key: &services.Key_Ed25519{Ed25519: []byte(k.data)}}
}

func (k MockKey) String() string {
	return k.data
}

func TestNewKeyList(t *testing.T) {
	kl := NewKeyList()
	assert.NotNil(t, kl)
	assert.Equal(t, -1, kl.threshold)
	assert.Empty(t, kl.keys)
}

func TestKeyListWithThreshold(t *testing.T) {
	kl := KeyListWithThreshold(2)
	assert.NotNil(t, kl)
	assert.Equal(t, 2, kl.threshold)
	assert.Empty(t, kl.keys)
}

func TestSetThreshold(t *testing.T) {
	kl := NewKeyList()
	kl.SetThreshold(3)
	assert.Equal(t, 3, kl.threshold)
}

func TestAdd(t *testing.T) {
	kl := NewKeyList()
	key := MockKey{data: "key1"}
	kl.Add(key)
	assert.Len(t, kl.keys, 1)
	assert.Equal(t, key, kl.keys[0])
}

func TestAddAll(t *testing.T) {
	kl := NewKeyList()
	keys := []Key{MockKey{data: "key1"}, MockKey{data: "key2"}}
	kl.AddAll(keys)
	assert.Len(t, kl.keys, 2)
	assert.Equal(t, keys[0], kl.keys[0])
	assert.Equal(t, keys[1], kl.keys[1])
}

func TestAddAllPublicKeys(t *testing.T) {
	kl := NewKeyList()
	keys := []PublicKey{{ed25519PublicKey: &_Ed25519PublicKey{keyData: []byte{1, 2}}}, {ed25519PublicKey: &_Ed25519PublicKey{keyData: []byte{1}}}}
	kl.AddAllPublicKeys(keys)
	assert.Len(t, kl.keys, 2)
	assert.Equal(t, keys[0], kl.keys[0])
	assert.Equal(t, keys[1], kl.keys[1])
}

func TestStringKeyList(t *testing.T) {
	kl := KeyListWithThreshold(2)
	key := MockKey{data: "key1"}
	kl.Add(key)
	expected := "{threshold:2,[key1]}"
	assert.Equal(t, expected, kl.String())

	kl2 := NewKeyList()
	kl2.Add(key)
	expected2 := "{[key1]}"
	assert.Equal(t, expected2, kl2.String())
}

func TestToProtoKey(t *testing.T) {
	kl := KeyListWithThreshold(2)
	key := MockKey{data: "key1"}
	kl.Add(key)
	protoKey := kl._ToProtoKey()

	expected := &services.Key{
		Key: &services.Key_ThresholdKey{
			ThresholdKey: &services.ThresholdKey{
				Threshold: uint32(kl.threshold),
				Keys: &services.KeyList{
					Keys: []*services.Key{
						{Key: &services.Key_Ed25519{Ed25519: []byte(key.data)}},
					},
				},
			},
		},
	}

	assert.True(t, proto.Equal(protoKey, expected))
}

func TestToProtoKeyList(t *testing.T) {
	kl := NewKeyList()
	key := MockKey{data: "key1"}
	kl.Add(key)
	protoKeyList := kl._ToProtoKeyList()

	expected := &services.KeyList{
		Keys: []*services.Key{
			{Key: &services.Key_Ed25519{Ed25519: []byte(key.data)}},
		},
	}

	assert.True(t, proto.Equal(protoKeyList, expected))
}

func TestKeyListFromProtobuf(t *testing.T) {
	pk, _ := PrivateKeyFromStringEd25519("302e020100300506032b657004220420db484b828e64b2d8f12ce3c0a0e93a0b8cce7af1bb8f39c97732394482538e11")
	protoKeyList := &services.KeyList{
		Keys: []*services.Key{
			{Key: &services.Key_Ed25519{Ed25519: pk.PublicKey().Bytes()}},
		},
	}

	kl, err := _KeyListFromProtobuf(protoKeyList)
	require.NoError(t, err)

	assert.Len(t, kl.keys, 1)
	assert.Equal(t, -1, kl.threshold)
}

func TestKeyListFromProtobuf_Nil(t *testing.T) {
	kl, err := _KeyListFromProtobuf(nil)
	assert.Error(t, errParameterNull, err)
	assert.Equal(t, KeyList{}, kl)
}
