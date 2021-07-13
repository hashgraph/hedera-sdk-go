package hedera

import (
	"fmt"
	"github.com/hashgraph/hedera-protobufs-go/services"
)

type KeyList struct {
	keys      []Key
	threshold int
}

func KeyListWithThreshold(threshold uint) *KeyList {
	return &KeyList{
		keys:      make([]Key, 0),
		threshold: int(threshold),
	}
}

func NewKeyList() *KeyList {
	return &KeyList{
		keys:      make([]Key, 0),
		threshold: -1,
	}
}

func (kl *KeyList) Add(key Key) *KeyList {
	kl.keys = append(kl.keys, key)
	return kl
}

func (kl *KeyList) AddAll(keys []Key) *KeyList {
	for _, key := range keys {
		kl.Add(key)
	}

	return kl
}

func (kl *KeyList) AddAllPublicKeys(keys []PublicKey) *KeyList {
	for _, key := range keys {
		kl.Add(key)
	}

	return kl
}

func (kl *KeyList) String() string {
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

func (kl *KeyList) toProtoKey() *services.Key {
	keys := make([]*services.Key, len(kl.keys))
	for i, key := range kl.keys {
		keys[i] = key.toProtoKey()
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
	} else {
		return &services.Key{
			Key: &services.Key_KeyList{
				KeyList: &services.KeyList{
					Keys: keys,
				},
			},
		}
	}
}

func (kl *KeyList) toProtoKeyList() *services.KeyList {
	keys := make([]*services.Key, len(kl.keys))
	for i, key := range kl.keys {
		keys[i] = key.toProtoKey()
	}

	return &services.KeyList{
		Keys: keys,
	}
}

func keyListFromProtobuf(pb *services.KeyList, networkName *NetworkName) (KeyList, error) {
	if pb == nil {
		return KeyList{}, errParameterNull
	}
	var keys []Key = make([]Key, len(pb.Keys))

	for i, pbKey := range pb.Keys {
		key, err := keyFromProtobuf(pbKey, networkName)

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
