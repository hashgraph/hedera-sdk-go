package hedera

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	"time"
)

type MirrorClient struct {
	conn *grpc.ClientConn
}

func newMirrorClient(endpoint string) (MirrorClient, error) {
	client, err := grpc.Dial(endpoint, grpc.WithKeepaliveParams(keepalive.ClientParameters{Time: 2 * time.Minute}))

	if err != nil {
		return MirrorClient{}, err
	}

	return MirrorClient{client}, err
}

func (mc MirrorClient) Close() error {
	if mc.conn == nil {
		return nil
	}

	return mc.conn.Close()
}
