package hedera

import "github.com/hashgraph/hedera-sdk-go/proto"

type ThresholdKey struct {
	threshold uint32
	keyList   KeyList
}

func NewThresholdKey(threshold uint32) ThresholdKey {
	return ThresholdKey{threshold: threshold}
}

func (tk ThresholdKey) Add(key PublicKey) ThresholdKey {
	tk.keyList.Add(key)
	return tk
}

func (tk ThresholdKey) AddAll(keys []PublicKey) ThresholdKey {
	tk.keyList.AddAll(keys)
	return tk
}

func (tk ThresholdKey) toProto() *proto.Key {
	keyList := tk.keyList.toProto().GetKeyList()

	tkProto := proto.ThresholdKey{Keys: keyList, Threshold: tk.threshold}
	keyProto := proto.Key{Key: &proto.Key_ThresholdKey{&tkProto}}

	return &keyProto
}
