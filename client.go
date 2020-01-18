package hedera

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"math/rand"
	"os"

	"google.golang.org/grpc"
)

// Default max fees and payments to 1 h-bar
var defaultMaxTransactionFee Hbar = NewHbar(1)
var defaultMaxQueryPayment Hbar = NewHbar(1)

type Client struct {
	maxTransactionFee Hbar
	maxQueryPayment   Hbar

	operator *operator

	networkNodes   map[AccountID]*node
	networkNodeIds []AccountID
}

type node struct {
	conn    *grpc.ClientConn
	id      AccountID
	address string
}

type signer func(message []byte) []byte

type operator struct {
	accountID  AccountID
	privateKey *Ed25519PrivateKey
	publicKey  Ed25519PublicKey
	signer     signer
}

var mainnetNodes = map[string]AccountID{
	"35.237.200.180:50211": AccountID{Account: 3},
	"35.186.191.247:50211": AccountID{Account: 4},
	"35.192.2.25:50211":    AccountID{Account: 5},
	"35.199.161.108:50211": AccountID{Account: 6},
	"35.203.82.240:50211":  AccountID{Account: 7},
	"35.236.5.219:50211":   AccountID{Account: 8},
	"35.197.192.225:50211": AccountID{Account: 9},
	"35.242.233.154:50211": AccountID{Account: 10},
	"35.240.118.96:50211":  AccountID{Account: 11},
	"35.204.86.32:50211":   AccountID{Account: 12},
}

var testnetNodes = map[string]AccountID{
	"0.testnet.hedera.com:50211": AccountID{Account: 3},
	"1.testnet.hedera.com:50211": AccountID{Account: 4},
	"2.testnet.hedera.com:50211": AccountID{Account: 5},
	"3.testnet.hedera.com:50211": AccountID{Account: 6},
}

func ClientForMainnet() *Client {
	return NewClient(mainnetNodes)
}

func ClientForTestnet() *Client {
	return NewClient(testnetNodes)
}

func NewClient(network map[string]AccountID) *Client {
	client := &Client{
		maxQueryPayment:   defaultMaxQueryPayment,
		maxTransactionFee: defaultMaxTransactionFee,
		networkNodes:      map[AccountID]*node{},
		networkNodeIds:    []AccountID{},
	}

	client.ReplaceNodes(network)

	return client
}

type configOperator struct {
	AccountID  string `json:"accountId"`
	PrivateKey string `json:"privateKey"`
}

type clientConfig struct {
	Network  map[string]AccountID `json:"network"`
	Operator *configOperator      `json:"operator"`
}

func ClientFromJSON(jsonBytes []byte) (*Client, error) {
	var clientConfig clientConfig

	err := json.Unmarshal(jsonBytes, &clientConfig)
	if err != nil {
		return nil, err
	}

	client := NewClient(clientConfig.Network)

	// if the operator is not provided, finish here
	if clientConfig.Operator == nil {
		return client, nil
	}

	operatorId, err := AccountIDFromString(clientConfig.Operator.AccountID)
	if err != nil {
		return nil, err
	}

	operatorKey, err := Ed25519PrivateKeyFromString(clientConfig.Operator.PrivateKey)
	if err != nil {
		return nil, err
	}

	operator := operator{
		accountID:  operatorId,
		privateKey: &operatorKey,
		publicKey:  operatorKey.PublicKey(),
		signer:     operatorKey.Sign,
	}

	client.operator = &operator

	return client, nil
}

func ClientFromFile(filename string) (*Client, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	defer func() {
		err = file.Close()
	}()

	configBytes, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	return ClientFromJSON(configBytes)
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

func (client *Client) ReplaceNodes(network map[string]AccountID) *Client {
	for address, id := range network {
		client.networkNodeIds = append(client.networkNodeIds, id)
		client.networkNodes[id] = &node{
			id:      id,
			address: address,
		}
	}

	return client
}

func (client *Client) SetOperator(accountID AccountID, privateKey Ed25519PrivateKey) *Client {
	client.operator = &operator{
		accountID:  accountID,
		privateKey: &privateKey,
		publicKey:  privateKey.PublicKey(),
		signer:     privateKey.Sign,
	}

	return client
}

func (client *Client) SetOperatorWith(accountID AccountID, publicKey Ed25519PublicKey, signer signer) *Client {
	client.operator = &operator{
		accountID:  accountID,
		privateKey: nil,
		publicKey:  publicKey,
		signer:     signer,
	}

	return client
}

func (client *Client) SetMaxTransactionFee(fee Hbar) *Client {
	client.maxTransactionFee = fee
	return client
}

func (client *Client) SetMaxQueryPayment(payment Hbar) *Client {
	client.maxQueryPayment = payment
	return client
}

func (client *Client) randomNode() *node {
	nodeIndex := rand.Intn(len(client.networkNodeIds))
	nodeID := client.networkNodeIds[nodeIndex]

	return client.networkNodes[nodeID]
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
