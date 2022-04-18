package hedera

/*-
 *
 * Hedera Go SDK
 *
 * Copyright (C) 2020 - 2022 Hedera Hashgraph, LLC
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
	"github.com/hashgraph/hedera-protobufs-go/services"
	protobuf "google.golang.org/protobuf/proto"

	"time"
)

type LiveHash struct {
	AccountID AccountID
	Hash      []byte
	Keys      KeyList
	Duration  time.Time
}

func (liveHash *LiveHash) _ToProtobuf() *services.LiveHash {
	return &services.LiveHash{
		AccountId: liveHash.AccountID._ToProtobuf(),
		Hash:      liveHash.Hash,
		Keys:      liveHash.Keys._ToProtoKeyList(),
		Duration: &services.Duration{
			Seconds: int64(liveHash.Duration.Second()),
		},
	}
}

func _LiveHashFromProtobuf(hash *services.LiveHash) (LiveHash, error) {
	if hash == nil {
		return LiveHash{}, errParameterNull
	}
	keyList, err := _KeyListFromProtobuf(hash.Keys)
	if err != nil {
		return LiveHash{}, err
	}

	accountID := AccountID{}
	if hash.AccountId != nil {
		accountID = *_AccountIDFromProtobuf(hash.AccountId)
	}

	return LiveHash{
		AccountID: accountID,
		Hash:      hash.Hash,
		Keys:      keyList,
		Duration: time.Date(time.Now().Year(), time.Now().Month(),
			time.Now().Day(), time.Now().Hour(), time.Now().Minute(),
			int(hash.Duration.Seconds), time.Now().Nanosecond(), time.Now().Location()),
	}, nil
}

func (liveHash LiveHash) ToBytes() []byte {
	data, err := protobuf.Marshal(liveHash._ToProtobuf())
	if err != nil {
		return make([]byte, 0)
	}

	return data
}

func LiveHashFromBytes(data []byte) (LiveHash, error) {
	if data == nil {
		return LiveHash{}, errByteArrayNull
	}
	pb := services.LiveHash{}
	err := protobuf.Unmarshal(data, &pb)
	if err != nil {
		return LiveHash{}, err
	}

	liveHash, err := _LiveHashFromProtobuf(&pb)
	if err != nil {
		return LiveHash{}, err
	}

	return liveHash, nil
}
