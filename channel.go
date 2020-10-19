package hedera

import (
	"github.com/hashgraph/hedera-sdk-go/proto"
	"google.golang.org/grpc"
)

type channel struct {
	crypto   proto.CryptoServiceClient
	file     proto.FileServiceClient
	contract proto.SmartContractServiceClient
	topic    proto.ConsensusServiceClient
	freeze   proto.FreezeServiceClient
	network  proto.NetworkServiceClient
	token    proto.TokenServiceClient
	client   *grpc.ClientConn
}

func newChannel(client *grpc.ClientConn) channel {
	return channel{
		client: client,
	}
}

func (channel channel) getCrypto() proto.CryptoServiceClient {
	if channel.crypto == nil {
		channel.crypto = proto.NewCryptoServiceClient(channel.client)
	}

	return channel.crypto
}

func (channel channel) getFile() proto.FileServiceClient {
	if channel.file == nil {
		channel.file = proto.NewFileServiceClient(channel.client)
	}

	return channel.file
}

func (channel channel) getContract() proto.SmartContractServiceClient {
	if channel.contract == nil {
		channel.contract = proto.NewSmartContractServiceClient(channel.client)
	}

	return channel.contract
}

func (channel channel) getTopic() proto.ConsensusServiceClient {
	if channel.topic == nil {
		channel.topic = proto.NewConsensusServiceClient(channel.client)
	}

	return channel.topic
}

func (channel channel) getFreeze() proto.FreezeServiceClient {
	if channel.freeze == nil {
		channel.freeze = proto.NewFreezeServiceClient(channel.client)
	}

	return channel.freeze
}

func (channel channel) getNetwork() proto.NetworkServiceClient {
	if channel.network == nil {
		channel.network = proto.NewNetworkServiceClient(channel.client)
	}

	return channel.network
}

func (channel channel) getToken() proto.TokenServiceClient {
	if channel.token == nil {
		channel.token = proto.NewTokenServiceClient(channel.client)
	}

	return channel.token
}
