package hedera

import "google.golang.org/grpc"

// fixme: use the actual structs when they are available
type PublicKey []byte
type PrivateKey []byte

type Client struct {
	// todo: support multiple nodes
	nodeId AccountId
	operator *operator
	conn *grpc.ClientConn
}

type operator struct {
	accountId AccountId
	privateKey PrivateKey
}

func NewClient(nodeId AccountId, address string) (*Client, error) {
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	client := Client {
		nodeId,
		nil,
		conn,
	}

	return &client, nil
}

func (client Client) SetOperator(accountId AccountId, privateKey PrivateKey) Client {
	operator := operator {
		accountId,
		privateKey,
	}

	client.operator = &operator

	return client
}

func (client Client) OperatorId() *AccountId {
	if client.operator == nil {
		return nil
	}

	return &client.operator.accountId
}

func (client Client) OperatorPrivateKey() *PrivateKey {
	if client.operator == nil {
		return nil
	}

	return &client.operator.privateKey
}
