package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"fmt"

	"github.com/hiero-ledger/hiero-sdk-go/v2/proto/services"
)

// KeyList is a list of keys
type KeyList struct {
	keys      []Key
	threshold int
}

// NewKeyListWithThreshold creates a new KeyList with the given threshold
func KeyListWithThreshold(threshold uint) *KeyList {
	return &KeyList{
		keys:      make([]Key, 0),
		threshold: int(threshold),
	}
}

// NewKeyList creates a new KeyList with no threshold
func NewKeyList() *KeyList {
	return &KeyList{
		keys:      make([]Key, 0),
		threshold: -1,
	}
}

// SetThreshold sets the threshold of the KeyList
func (kl *KeyList) SetThreshold(threshold int) *KeyList {
	kl.threshold = threshold
	return kl
}

// Add adds a key to the KeyList
func (kl *KeyList) Add(key Key) *KeyList {
	kl.keys = append(kl.keys, key)
	return kl
}

// AddAll adds all the keys to the KeyList
func (kl *KeyList) AddAll(keys []Key) *KeyList {
	for _, key := range keys {
		kl.Add(key)
	}

	return kl
}

// AddAllPublicKeys adds all the public keys to the KeyList
func (kl *KeyList) AddAllPublicKeys(keys []PublicKey) *KeyList {
	for _, key := range keys {
		kl.Add(key)
	}

	return kl
}

// GetKeys returns the internal list of Keys.
func (kl KeyList) GetKeys() []Key {
	return kl.keys
}

// GetThreshold returns the threshold value set on the KeyList.
// A value of -1 means that there is no threshold set.
func (kl KeyList) GetThreshold() int {
	return kl.threshold
}

// String returns a string representation of the KeyList
func (kl KeyList) String() string {
	var s string
	if kl.threshold > 0 {
		s = "{threshold:" + fmt.Sprint(kl.threshold) + ",["
	} else {
		s = "{["
	}

	for i, key := range kl.keys {
		s += key.String()
		if i != len(kl.keys)-1 {
			s += ","
		}
	}

	s += "]}"

	return s
}

func (kl KeyList) _ToProtoKey() *services.Key {
	keys := make([]*services.Key, len(kl.keys))
	for i, key := range kl.keys {
		keys[i] = key._ToProtoKey()
	}

	if kl.threshold >= 0 {
		return &services.Key{
			Key: &services.Key_ThresholdKey{
				ThresholdKey: &services.ThresholdKey{
					Threshold: uint32(kl.threshold),
					Keys: &services.KeyList{
						Keys: keys,
					},
				},
			},
		}
	}

	return &services.Key{
		Key: &services.Key_KeyList{
			KeyList: &services.KeyList{
				Keys: keys,
			},
		},
	}
}

func (kl *KeyList) _ToProtoKeyList() *services.KeyList {
	keys := make([]*services.Key, len(kl.keys))
	for i, key := range kl.keys {
		keys[i] = key._ToProtoKey()
	}

	return &services.KeyList{
		Keys: keys,
	}
}

func _KeyListFromProtobuf(pb *services.KeyList) (KeyList, error) {
	if pb == nil {
		return KeyList{}, errParameterNull
	}
	var keys = make([]Key, len(pb.Keys))

	for i, pbKey := range pb.Keys {
		key, err := _KeyFromProtobuf(pbKey)

		if err != nil {
			return KeyList{}, err
		}

		keys[i] = key
	}

	return KeyList{
		keys:      keys,
		threshold: -1,
	}, nil
}
