package hedera

import (
	"github.com/hashgraph/hedera-sdk-go/proto/mirror"
	"google.golang.org/grpc"
)

type MirrorClient struct {
	connection *grpc.ClientConn
	client     mirror.ConsensusServiceClient
}

func NewMirrorClient(endpoint string) (MirrorClient, error) {
	conn, err := grpc.Dial(endpoint, grpc.WithInsecure())
	if err != nil {
		return MirrorClient{}, err
	}

	return MirrorClient{connection: conn, client: mirror.NewConsensusServiceClient(conn)}, nil
}

func (mc MirrorClient) Close() error {
	if mc.connection == nil {
		return nil
	}

	if err := mc.connection.Close(); err != nil {
		return err
	}

	mc.connection = nil

	return nil
}
