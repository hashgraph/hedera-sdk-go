package hedera

import (
	"github.com/hashgraph/hedera-sdk-go/proto"
)

type KeyList struct {
	keys []*proto.Key
}

func NewKeyList() *KeyList {
	return &KeyList{
		keys: []*proto.Key{},
	}
}

func (kl *KeyList) Add(key PublicKey) *KeyList {
	kl.keys = append(kl.keys, key.toProto())
	return kl
}

func (kl *KeyList) AddAll(keys []PublicKey) *KeyList {
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

func (kl *KeyList) toProto() *proto.Key {
	return &proto.Key{Key: &proto.Key_KeyList{KeyList: &proto.KeyList{Keys: kl.keys}}}
}

func (kl *KeyList) toProtoKeyList() *proto.KeyList {
	return &proto.KeyList{Keys: kl.keys}
}
