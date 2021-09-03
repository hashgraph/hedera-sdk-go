package hedera

import (
	"github.com/hashgraph/hedera-sdk-go/v2/proto"
	"google.golang.org/grpc"
)

type _Channel struct {
	crypto   proto.CryptoServiceClient
	file     proto.FileServiceClient
	contract proto.SmartContractServiceClient
	topic    proto.ConsensusServiceClient
	freeze   proto.FreezeServiceClient
	network  proto.NetworkServiceClient
	token    proto.TokenServiceClient
	schedule proto.ScheduleServiceClient
	client   *grpc.ClientConn
}

func _NewChannel(client *grpc.ClientConn) _Channel {
	return _Channel{
		client: client,
	}
}

func (channel _Channel) _GetCrypto() proto.CryptoServiceClient {
	if channel.crypto == nil {
		channel.crypto = proto.NewCryptoServiceClient(channel.client)
	}

	return channel.crypto
}

func (channel _Channel) _GetFile() proto.FileServiceClient {
	if channel.file == nil {
		channel.file = proto.NewFileServiceClient(channel.client)
	}

	return channel.file
}

func (channel _Channel) _GetContract() proto.SmartContractServiceClient {
	if channel.contract == nil {
		channel.contract = proto.NewSmartContractServiceClient(channel.client)
	}

	return channel.contract
}

func (channel _Channel) _GetTopic() proto.ConsensusServiceClient {
	if channel.topic == nil {
		channel.topic = proto.NewConsensusServiceClient(channel.client)
	}

	return channel.topic
}

func (channel _Channel) _GetFreeze() proto.FreezeServiceClient {
	if channel.freeze == nil {
		channel.freeze = proto.NewFreezeServiceClient(channel.client)
	}

	return channel.freeze
}

func (channel _Channel) _GetNetwork() proto.NetworkServiceClient {
	if channel.network == nil {
		channel.network = proto.NewNetworkServiceClient(channel.client)
	}

	return channel.network
}

func (channel _Channel) _GetToken() proto.TokenServiceClient {
	if channel.token == nil {
		channel.token = proto.NewTokenServiceClient(channel.client)
	}

	return channel.token
}

func (channel _Channel) _GetSchedule() proto.ScheduleServiceClient {
	if channel.schedule == nil {
		channel.schedule = proto.NewScheduleServiceClient(channel.client)
	}

	return channel.schedule
}
