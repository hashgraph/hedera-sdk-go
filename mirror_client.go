package hedera

type MirrorClient struct {
	ConsensusClient
}

func newMirrorClient(endpoint string) (MirrorClient, error) {
	client, err := NewConsensusClient(endpoint)

	if err != nil {
		return MirrorClient{}, nil
	}

	return MirrorClient{*client}, nil
}

func (mc MirrorClient) Close() error {
	if mc.conn == nil {
		return nil
	}

	return mc.conn.Close()
}
