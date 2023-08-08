package hedera

/*-
 *
 * Hedera Go SDK
 *
 * Copyright (C) 2020 - 2023 Hedera Hashgraph, LLC
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
)

// ContractLogInfo is the log info for events returned by a function
type ContractLogInfo struct {
	ContractID ContractID
	Bloom      []byte
	Topics     [][]byte
	Data       []byte
}

func _ContractLogInfoFromProtobuf(pb *services.ContractLoginfo) ContractLogInfo {
	if pb == nil {
		return ContractLogInfo{}
	}

	contractID := ContractID{}
	if pb.ContractID != nil {
		contractID = *_ContractIDFromProtobuf(pb.ContractID)
	}

	return ContractLogInfo{
		ContractID: contractID,
		Bloom:      pb.Bloom,
		Topics:     pb.Topic,
		Data:       pb.Data,
	}
}

func (logInfo ContractLogInfo) _ToProtobuf() *services.ContractLoginfo {
	return &services.ContractLoginfo{
		ContractID: logInfo.ContractID._ToProtobuf(),
		Bloom:      logInfo.Bloom,
		Topic:      logInfo.Topics,
		Data:       logInfo.Data,
	}
}
