package hedera

import (
	"github.com/hashgraph/hedera-sdk-go/proto"
)

type KeyList struct {
	keys      []*proto.Key
	threshold int
}

func KeyListWithThreshold(threshold uint) *KeyList {
	return &KeyList{
		keys:      []*proto.Key{},
		threshold: int(threshold),
	}
}

func NewKeyList() *KeyList {
	return &KeyList{
		keys:      []*proto.Key{},
		threshold: -1,
	}
}

func (kl *KeyList) Add(key Key) *KeyList {
	kl.keys = append(kl.keys, key.toProtobuf())
	return kl
}

func (kl *KeyList) AddAll(keys []Key) *KeyList {
	for _, key := range keys {
		kl.Add(key)
	}

	return kl
}

func (kl *KeyList) toProtobuf() *proto.Key {
	if kl.threshold >= 0 {
		return &proto.Key{
			Key: &proto.Key_ThresholdKey{
				ThresholdKey: &proto.ThresholdKey{
					Threshold: uint32(kl.threshold),
					Keys: &proto.KeyList{
						Keys: kl.keys,
					},
				},
			},
		}
	} else {
		return &proto.Key{
			Key: &proto.Key_KeyList{
				KeyList: &proto.KeyList{
					Keys: kl.keys,
				},
			},
		}
	}
}
