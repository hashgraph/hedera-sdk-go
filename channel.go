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
	"google.golang.org/grpc"
)

type _Channel struct {
	crypto   services.CryptoServiceClient
	file     services.FileServiceClient
	contract services.SmartContractServiceClient
	topic    services.ConsensusServiceClient
	freeze   services.FreezeServiceClient
	network  services.NetworkServiceClient
	token    services.TokenServiceClient
	schedule services.ScheduleServiceClient
	util     services.UtilServiceClient
	client   *grpc.ClientConn
}

func _NewChannel(client *grpc.ClientConn) _Channel {
	return _Channel{
		client: client,
	}
}

func (channel _Channel) _GetCrypto() services.CryptoServiceClient {
	if channel.crypto == nil {
		channel.crypto = services.NewCryptoServiceClient(channel.client)
	}

	return channel.crypto
}

func (channel _Channel) _GetFile() services.FileServiceClient {
	if channel.file == nil {
		channel.file = services.NewFileServiceClient(channel.client)
	}

	return channel.file
}

func (channel _Channel) _GetContract() services.SmartContractServiceClient {
	if channel.contract == nil {
		channel.contract = services.NewSmartContractServiceClient(channel.client)
	}

	return channel.contract
}

func (channel _Channel) _GetTopic() services.ConsensusServiceClient {
	if channel.topic == nil {
		channel.topic = services.NewConsensusServiceClient(channel.client)
	}

	return channel.topic
}

func (channel _Channel) _GetFreeze() services.FreezeServiceClient {
	if channel.freeze == nil {
		channel.freeze = services.NewFreezeServiceClient(channel.client)
	}

	return channel.freeze
}

func (channel _Channel) _GetNetwork() services.NetworkServiceClient {
	if channel.network == nil {
		channel.network = services.NewNetworkServiceClient(channel.client)
	}

	return channel.network
}

func (channel _Channel) _GetToken() services.TokenServiceClient {
	if channel.token == nil {
		channel.token = services.NewTokenServiceClient(channel.client)
	}

	return channel.token
}

func (channel _Channel) _GetSchedule() services.ScheduleServiceClient {
	if channel.schedule == nil {
		channel.schedule = services.NewScheduleServiceClient(channel.client)
	}

	return channel.schedule
}

func (channel _Channel) _GetUtil() services.UtilServiceClient {
	if channel.util == nil {
		channel.util = services.NewUtilServiceClient(channel.client)
	}

	return channel.util
}
