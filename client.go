package hedera

import (
	"context"
	"google.golang.org/grpc"
	"math/rand"
)

// Default max fees and payments to 1 h-bar
const defaultMaxTransactionFee uint64 = 100_000_000
const defaultMaxQueryPayment uint64 = 100_000_000

type Client struct {
	maxTransactionFee uint64
	maxQueryPayment   uint64

	operator *operator

	networkNodes   map[AccountID]*node
	networkNodeIds []AccountID
}

type node struct {
	conn    *grpc.ClientConn
	id      AccountID
	address string
}

type operator struct {
	accountID  AccountID
	privateKey Ed25519PrivateKey
}

func NewClient(network map[string]AccountID) *Client {
	networkNodes := map[AccountID]*node{}
	var networkNodeIds []AccountID

	for address, id := range network {
		networkNodeIds = append(networkNodeIds, id)
		networkNodes[id] = &node{
			id:      id,
			address: address,
		}
	}

	return &Client{
		maxQueryPayment:   defaultMaxQueryPayment,
		maxTransactionFee: defaultMaxTransactionFee,
		networkNodes:      networkNodes,
		networkNodeIds:    networkNodeIds,
	}
}

func (client *Client) Close() error {
	for _, node := range client.networkNodes {
		if node.conn != nil {
			err := node.conn.Close()
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (client *Client) SetOperator(accountID AccountID, privateKey Ed25519PrivateKey) *Client {
	client.operator = &operator{
		accountID,
		privateKey,
	}

	return client
}

func (client *Client) SetMaxTransactionFee(tinyBars uint64) *Client {
	client.maxTransactionFee = tinyBars
	return client
}

func (client *Client) SetMaxQueryPayment(tinyBars uint64) *Client {
	client.maxTransactionFee = tinyBars
	return client
}

func (client *Client) randomNode() *node {
	nodeIndex := rand.Intn(len(client.networkNodeIds))
	nodeId := client.networkNodeIds[nodeIndex]

	return client.networkNodes[nodeId]
}

func (client *Client) node(id AccountID) *node {
	return client.networkNodes[id]
}

func (node *node) invoke(method string, in interface{}, out interface{}) error {
	if node.conn == nil {
		conn, err := grpc.Dial(node.address, grpc.WithInsecure())
		if err != nil {
			return err
		}

		node.conn = conn
	}

	return node.conn.Invoke(context.TODO(), method, in, out)
}
