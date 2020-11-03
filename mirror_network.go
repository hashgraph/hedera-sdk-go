package hedera

import (
	"math/rand"

	"github.com/hashgraph/hedera-sdk-go/proto/mirror"
	"google.golang.org/grpc"
)

type mirrorNetwork struct {
	channels map[string]mirror.ConsensusServiceClient
	network  []string
	index    uint
}

func newMirrorNetwork(network []string) mirrorNetwork {
	if len(network) > 0 {
		rand.Shuffle(len(network), func(i, j int) {
			network[i], network[j] = network[j], network[i]
		})
	}

	return mirrorNetwork{
		channels: make(map[string]mirror.ConsensusServiceClient),
		network:  network,
		index:    0,
	}
}

func contains(arr []string, str string) bool {
	for _, a := range arr {
		if a == str {
			return true
		}
	}
	return false
}

func (network mirrorNetwork) setNetwork(newNetwork []string) {
	for _, n := range network.network {
		if !contains(newNetwork, n) {
			delete(network.channels, n)
		}
	}

	for _, n := range newNetwork {
		if !contains(network.network, n) {
			network.network = append(network.network, n)
		}
	}

	network.index = 0

	if len(network.network) > 0 {
		rand.Shuffle(len(network.network), func(i, j int) {
			network.network[i], network.network[j] = network.network[j], network.network[i]
		})
	}
}

func (network mirrorNetwork) getNextChannel() (mirror.ConsensusServiceClient, error) {
	if channel, ok := network.channels[network.network[network.index]]; ok {
		network.index = (network.index + 1) % uint(len(network.network))
		return channel, nil
	}

	conn, err := grpc.Dial(network.network[network.index], grpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	channel := mirror.NewConsensusServiceClient(conn)

	network.channels[network.network[network.index]] = channel
	network.index = (network.index + 1) % uint(len(network.network))
	return channel, nil
}
