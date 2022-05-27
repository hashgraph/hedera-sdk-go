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
	"time"

	"github.com/hashgraph/hedera-protobufs-go/services"
	protobuf "google.golang.org/protobuf/proto"
)

type StakingInfo struct {
	DeclineStakingReward bool
	StakePeriodStart     *time.Time
	PendingReward        int64
	StakedToMe           Hbar
	StakedAccountID      *AccountID
	StakedNodeID         *int64
}

func _StakingInfoFromProtobuf(pb *services.StakingInfo) StakingInfo {
	var start time.Time
	if pb.StakePeriodStart != nil {
		start = _TimeFromProtobuf(pb.StakePeriodStart)
	}

	body := StakingInfo{
		DeclineStakingReward: pb.DeclineReward,
		StakePeriodStart:     &start,
		PendingReward:        pb.PendingReward,
		StakedToMe:           HbarFromTinybar(pb.StakedToMe),
	}

	switch temp := pb.StakedId.(type) {
	case *services.StakingInfo_StakedAccountId:
		body.StakedAccountID = _AccountIDFromProtobuf(temp.StakedAccountId)
	case *services.StakingInfo_StakedNodeId:
		body.StakedNodeID = &temp.StakedNodeId
	}

	return body
}

func (stakingInfo *StakingInfo) _ToProtobuf() *services.StakingInfo { // nolint
	body := services.StakingInfo{
		DeclineReward: stakingInfo.DeclineStakingReward,
		PendingReward: stakingInfo.PendingReward,
		StakedToMe:    stakingInfo.StakedToMe.AsTinybar(),
	}

	if stakingInfo.StakePeriodStart != nil {
		body.StakePeriodStart = _TimeToProtobuf(*stakingInfo.StakePeriodStart)
	}

	if stakingInfo.StakedAccountID != nil {
		body.StakedId = &services.StakingInfo_StakedAccountId{StakedAccountId: stakingInfo.StakedAccountID._ToProtobuf()}
	} else if stakingInfo.StakedNodeID != nil {
		body.StakedId = &services.StakingInfo_StakedNodeId{StakedNodeId: *stakingInfo.StakedNodeID}
	}

	return &body
}

func (stakingInfo *StakingInfo) ToBytes() []byte {
	data, err := protobuf.Marshal(stakingInfo._ToProtobuf())
	if err != nil {
		return make([]byte, 0)
	}

	return data
}

func StakingInfoFromBytes(data []byte) (StakingInfo, error) {
	if data == nil {
		return StakingInfo{}, errByteArrayNull
	}
	pb := services.StakingInfo{}
	err := protobuf.Unmarshal(data, &pb)
	if err != nil {
		return StakingInfo{}, err
	}

	info := _StakingInfoFromProtobuf(&pb)

	return info, nil
}
