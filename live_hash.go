package hedera

import (
	protobuf "github.com/golang/protobuf/proto"
	"github.com/hashgraph/hedera-sdk-go/v2/proto"

	"time"
)

type LiveHash struct {
	AccountID AccountID
	Hash      []byte
	Keys      KeyList
	Duration  time.Time
}

func newLiveHash(accountId AccountID, hash []byte, keys KeyList, duration time.Time) LiveHash {
	return LiveHash{
		AccountID: accountId,
		Hash:      hash,
		Keys:      keys,
		Duration:  duration,
	}
}

func (liveHash *LiveHash) toProtobuf() *proto.LiveHash {
	return &proto.LiveHash{
		AccountId: liveHash.AccountID.toProtobuf(),
		Hash:      liveHash.Hash,
		Keys:      liveHash.Keys.toProtoKeyList(),
		Duration: &proto.Duration{
			Seconds: int64(liveHash.Duration.Second()),
		},
	}
}

func liveHashFromProtobuf(hash *proto.LiveHash) (LiveHash, error) {
	if hash == nil {
		return LiveHash{}, errParameterNull
	}
	keyList, err := keyListFromProtobuf(hash.Keys)
	if err != nil {
		return LiveHash{}, err
	}

	return LiveHash{
		AccountID: accountIDFromProtobuf(hash.GetAccountId()),
		Hash:      hash.Hash,
		Keys:      keyList,
		Duration: time.Date(time.Now().Year(), time.Now().Month(),
			time.Now().Day(), time.Now().Hour(), time.Now().Minute(),
			int(hash.Duration.Seconds), time.Now().Nanosecond(), time.Now().Location()),
	}, nil
}

func (liveHash LiveHash) ToBytes() []byte {
	data, err := protobuf.Marshal(liveHash.toProtobuf())
	if err != nil {
		return make([]byte, 0)
	}

	return data
}

func LiveHashFromBytes(data []byte) (LiveHash, error) {
	if data == nil {
		return LiveHash{}, errByteArrayNull
	}
	pb := proto.LiveHash{}
	err := protobuf.Unmarshal(data, &pb)
	if err != nil {
		return LiveHash{}, err
	}

	liveHash, err := liveHashFromProtobuf(&pb)
	if err != nil {
		return LiveHash{}, err
	}

	return liveHash, nil
}
