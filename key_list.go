package hedera

import (
	"fmt"
	"github.com/hashgraph/hedera-sdk-go/v2/proto"
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

func (kl *KeyList) toProtoKey() *proto.Key {
	keys := make([]*proto.Key, len(kl.keys))
	for i, key := range kl.keys {
		keys[i] = key.toProtoKey()
	}

	if kl.threshold >= 0 {
		return &proto.Key{
			Key: &proto.Key_ThresholdKey{
				ThresholdKey: &proto.ThresholdKey{
					Threshold: uint32(kl.threshold),
					Keys: &proto.KeyList{
						Keys: keys,
					},
				},
			},
		}
	} else {
		return &proto.Key{
			Key: &proto.Key_KeyList{
				KeyList: &proto.KeyList{
					Keys: keys,
				},
			},
		}
	}
}

func (kl *KeyList) toProtoKeyList() *proto.KeyList {
	keys := make([]*proto.Key, len(kl.keys))
	for i, key := range kl.keys {
		keys[i] = key.toProtoKey()
	}

	return &proto.KeyList{
		Keys: keys,
	}
}

func keyListFromProtobuf(pb *proto.KeyList) (KeyList, error) {
	var keys []Key = make([]Key, len(pb.Keys))

	for i, pbKey := range pb.Keys {
		key, err := keyFromProtobuf(pbKey)

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
