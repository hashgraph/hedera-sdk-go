package hedera

import (
	"github.com/hashgraph/hedera-protobufs-go/services"
	"google.golang.org/grpc"
)

type channel struct {
	crypto   services.CryptoServiceClient
	file     services.FileServiceClient
	contract services.SmartContractServiceClient
	topic    services.ConsensusServiceClient
	freeze   services.FreezeServiceClient
	network  services.NetworkServiceClient
	token    services.TokenServiceClient
	schedule services.ScheduleServiceClient
	client   *grpc.ClientConn
}

func newChannel(client *grpc.ClientConn) channel {
	return channel{
		client: client,
	}
}

func (channel channel) getCrypto() services.CryptoServiceClient {
	if channel.crypto == nil {
		channel.crypto = services.NewCryptoServiceClient(channel.client)
	}

	return channel.crypto
}

func (channel channel) getFile() services.FileServiceClient {
	if channel.file == nil {
		channel.file = services.NewFileServiceClient(channel.client)
	}

	return channel.file
}

func (channel channel) getContract() services.SmartContractServiceClient {
	if channel.contract == nil {
		channel.contract = services.NewSmartContractServiceClient(channel.client)
	}

	return channel.contract
}

func (channel channel) getTopic() services.ConsensusServiceClient {
	if channel.topic == nil {
		channel.topic = services.NewConsensusServiceClient(channel.client)
	}

	return channel.topic
}

func (channel channel) getFreeze() services.FreezeServiceClient {
	if channel.freeze == nil {
		channel.freeze = services.NewFreezeServiceClient(channel.client)
	}

	return channel.freeze
}

func (channel channel) getNetwork() services.NetworkServiceClient {
	if channel.network == nil {
		channel.network = services.NewNetworkServiceClient(channel.client)
	}

	return channel.network
}

func (channel channel) getToken() services.TokenServiceClient {
	if channel.token == nil {
		channel.token = services.NewTokenServiceClient(channel.client)
	}

	return channel.token
}

func (channel channel) getSchedule() services.ScheduleServiceClient {
	if channel.schedule == nil {
		channel.schedule = services.NewScheduleServiceClient(channel.client)
	}

	return channel.schedule
}
