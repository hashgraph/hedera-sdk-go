package hedera

import (
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
