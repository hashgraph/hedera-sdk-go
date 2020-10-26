package hedera

import (
	"encoding/json"
	"io/ioutil"
	"os"

	// "github.com/hashgraph/hedera-sdk-go/proto"
	"google.golang.org/grpc"
)

// Default max fees and payments to 1 h-bar
var defaultMaxTransactionFee Hbar = NewHbar(1)
var defaultMaxQueryPayment Hbar = NewHbar(1)

// Client is the Hedera protocol wrapper for the SDK used by all
// transaction and query types.
type Client struct {
	maxTransactionFee Hbar
	maxQueryPayment   Hbar

	operator *operator

	networkChannels map[AccountID]*channel
	networkNodeIds  []NodeID
	network         map[string]NodeID

	mirrorChannels map[string]*grpc.ClientConn
	mirrorNetwork  []string

	nextNodeIndex uint
}

// TransactionSigner is a closure or function that defines how transactions will be signed
type TransactionSigner func(message []byte) []byte

type operator struct {
	accountID  AccountID
	privateKey *PrivateKey
	publicKey  PublicKey
	signer     TransactionSigner
}

var mainnetNodes = map[string]NodeID{
	"35.237.200.180:50211" : NewNodeID(AccountID{Shard:   0, Realm:   0, Account: 3}, "35.237.200.180:50211"),
	"35.186.191.247:50211" : NewNodeID(AccountID{Shard:   0, Realm:   0, Account: 4}, "35.186.191.247:50211"),
	"35.192.2.25:50211" : NewNodeID(AccountID{Shard:   0, Realm:   0, Account: 5}, "35.192.2.25:50211"),
	"35.199.161.108:50211" : NewNodeID(AccountID{Shard:   0, Realm:   0, Account: 6}, "35.199.161.108:50211"),
	"35.203.82.240:50211" : NewNodeID(AccountID{Shard:   0, Realm:   0, Account: 7}, "35.203.82.240:50211"),
	"35.236.5.219:50211" : NewNodeID(AccountID{Shard:   0, Realm:   0, Account: 8}, "35.236.5.219:50211"),
	"35.197.192.225:50211" : NewNodeID(AccountID{Shard:   0, Realm:   0, Account: 9}, "35.197.192.225:50211"),
	"35.242.233.154:50211" : NewNodeID(AccountID{Shard:   0, Realm:   0, Account: 10}, "35.242.233.154:50211"),
	"35.240.118.96:50211" : NewNodeID(AccountID{Shard:   0, Realm:   0, Account: 11}, "35.240.118.96:50211"),
	"35.204.86.32:50211" : NewNodeID(AccountID{Shard:   0, Realm:   0, Account: 12}, "35.204.86.32:50211"),
}

var testnetNodes = map[string]NodeID{
	"0.testnet.hedera.com:50211": NewNodeID(AccountID{Shard:   0, Realm:   0, Account: 3}, "0.previewnet.hedera.com:50211"),
	"1.testnet.hedera.com:50211": NewNodeID(AccountID{Shard:   0, Realm:   0, Account: 3}, "1.previewnet.hedera.com:50211"),
	"2.testnet.hedera.com:50211": NewNodeID(AccountID{Shard:   0, Realm:   0, Account: 3}, "2.previewnet.hedera.com:50211"),
	"3.testnet.hedera.com:50211": NewNodeID(AccountID{Shard:   0, Realm:   0, Account: 3}, "3.previewnet.hedera.com:50211"),
}

var previewnetNodes = map[string]NodeID {
	"0.previewnet.hedera.com:50211": NewNodeID(AccountID{Shard:   0, Realm:   0, Account: 3}, "0.previewnet.hedera.com:50211"),
	"1.previewnet.hedera.com:50211": NewNodeID(AccountID{Shard:   0, Realm:   0, Account: 4}, "1.previewnet.hedera.com:50211"),
	"2.previewnet.hedera.com:50211": NewNodeID(AccountID{Shard:   0, Realm:   0, Account: 5}, "2.previewnet.hedera.com:50211"),
	"3.previewnet.hedera.com:50211": NewNodeID(AccountID{Shard:   0, Realm:   0, Account: 6}, "3.previewnet.hedera.com:50211"),
}

var mainnetMirror = []string{"hcs.mainnet.mirrornode.hedera.com:5600"}
var testnetMirror = []string{"hcs.testnet.mirrornode.hedera.com:5600"}
var previewnetMirror = []string{"hcs.previewnet.mirrornode.hedera.com:5600"}

func ClientForNetwork(network map[string]NodeID) *Client {
	return newClient(network, []string{})
}

// ClientForMainnet returns a preconfigured client for use with the standard
// Hedera mainnet.
// Most users will want to set an operator account with .SetOperator so
// transactions can be automatically given TransactionIDs and signed.
func ClientForMainnet() *Client {
	return newClient(mainnetNodes, mainnetMirror)
}

// ClientForTestnet returns a preconfigured client for use with the standard
// Hedera testnet.
// Most users will want to set an operator account with .SetOperator so
// transactions can be automatically given TransactionIDs and signed.
func ClientForTestnet() *Client {
	return newClient(testnetNodes, testnetMirror)
}

// ClientForPreviewnet returns a preconfigured client for use with the standard
// Hedera previewnet.
// Most users will want to set an operator account with .SetOperator so
// transactions can be automatically given TransactionIDs and signed.
func ClientForPreviewnet() *Client {
	return newClient(previewnetNodes, previewnetMirror)
}

// newClient takes in a map of node addresses to their respective IDS (network)
// and returns a Client instance which can be used to
func newClient(network map[string]NodeID, mirrorNetwork []string) *Client {
	client := Client{
		maxQueryPayment:   defaultMaxQueryPayment,
		maxTransactionFee: defaultMaxTransactionFee,
		networkChannels:   make(map[AccountID]*channel),
		networkNodeIds:    make([]NodeID, 0),
		network:           make(map[string]NodeID),
		mirrorChannels:    make(map[string]*grpc.ClientConn),
		mirrorNetwork:     make([]string, 0),
		nextNodeIndex:     0,
	}

	client.SetNetwork(network)
	client.SetMirrorNetwork(mirrorNetwork)

	return &client
}

type configOperator struct {
	AccountID  string `json:"accountId"`
	PrivateKey string `json:"privateKey"`
}

// TODO: Implement complete spec: https://gitlab.com/launchbadge/hedera/sdk/python/-/issues/45
type clientConfig struct {
	Network       map[string]string `json:"network"`
	MirrorNetwork []string          `json:"mirrorNetwork"`
	Operator      *configOperator   `json:"operator"`
}

// ClientFromJSON takes in the byte slice representation of a JSON string or
// document and returns Client based on the configuration.
func ClientFromJSON(jsonBytes []byte) (*Client, error) {
	var clientConfig clientConfig

	err := json.Unmarshal(jsonBytes, &clientConfig)
	if err != nil {
		return nil, err
	}

	var network = make(map[string]NodeID)

	for url, id := range clientConfig.Network {
		accountID, err := AccountIDFromString(id)
		if err != nil {
			return nil, err
		}

		network[url] = NewNodeID(accountID, url)
	}

	client := newClient(network, clientConfig.MirrorNetwork)

	// if the operator is not provided, finish here
	if clientConfig.Operator == nil {
		return client, nil
	}

	operatorID, err := AccountIDFromString(clientConfig.Operator.AccountID)
	if err != nil {
		return nil, err
	}

	operatorKey, err := PrivateKeyFromString(clientConfig.Operator.PrivateKey)
	if err != nil {
		return nil, err
	}

	operator := operator{
		accountID:  operatorID,
		privateKey: &operatorKey,
		publicKey:  operatorKey.PublicKey(),
		signer:     operatorKey.Sign,
	}

	client.operator = &operator

	return client, nil
}

// ClientFromJsonFile takes a filename string representing the path to a JSON encoded
// Client file and returns a Client based on the configuration.
func ClientFromJsonFile(filename string) (*Client, error) {
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

// Close is used to disconnect the Client from the network
func (client *Client) Close() error {
	for _, conn := range client.networkChannels {
		if conn != nil {
			err := conn.client.Close()
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// SetNetwork replaces all nodes in the Client with a new set of nodes.
// (e.g. for an Address Book update).
func (client *Client) SetNetwork(network map[string]NodeID) *Client {
	for address, id := range network {
		client.networkNodeIds = append(client.networkNodeIds, id)
		client.network[address] = id
	}

	return client
}

// SetNetwork replaces all nodes in the Client with a new set of nodes.
// (e.g. for an Address Book update).
func (client *Client) SetMirrorNetwork(mirrorNetwork []string) *Client {
	client.mirrorNetwork = mirrorNetwork

	return client
}

// SetOperator sets that account that will, by default, be paying for
// transactions and queries built with the client and the associated key
// with which to automatically sign transactions.
func (client *Client) SetOperator(accountID AccountID, privateKey PrivateKey) *Client {
	client.operator = &operator{
		accountID:  accountID,
		privateKey: &privateKey,
		publicKey:  privateKey.PublicKey(),
		signer:     privateKey.Sign,
	}

	return client
}

// SetOperatorWith sets that account that will, by default, be paying for
// transactions and queries built with the client, the account's PublicKey
// and a callback that will be invoked when a transaction needs to be signed.
func (client *Client) SetOperatorWith(accountID AccountID, publicKey PublicKey, signer TransactionSigner) *Client {
	client.operator = &operator{
		accountID:  accountID,
		privateKey: nil,
		publicKey:  publicKey,
		signer:     signer,
	}

	return client
}

// GetOperatorID returns the ID for the operator
func (client *Client) GetOperatorID() AccountID {
	if client.operator != nil {
		return client.operator.accountID
	} else {
		return AccountID{}
	}
}

// GetOperatorKey returns the Key for the operator
func (client *Client) GetOperatorKey() PublicKey {
	if client.operator != nil {
		return client.operator.publicKey
	} else {
		return PublicKey{}
	}
}

// SetMaxTransactionFee sets the maximum fee to be paid for the transactions
// executed by the Client.
// Because transaction fees are always maximums the actual fee assessed for
// a given transaction may be less than this value, but never greater.
func (client *Client) SetMaxTransactionFee(fee Hbar) *Client {
	client.maxTransactionFee = fee
	return client
}

// SetMaxQueryPayment sets the default maximum payment allowable for queries.
func (client *Client) SetMaxQueryPayment(payment Hbar) *Client {
	client.maxQueryPayment = payment
	return client
}

// // Ping sends an AccountBalanceQuery to the specified node returning nil if no
// // problems occur. Otherwise, an error representing the status of the node will
// // be returned.
// func (client *Client) Ping(nodeID AccountID) error {
// 	node := client.networkNodes[nodeID]
// 	if node == nil {
// 		return fmt.Errorf("node with ID %s not registered on this client", nodeID)
// 	}

// 	pingQuery := NewAccountBalanceQuery().
// 		SetAccountID(nodeID)

// 	pb := pingQuery.QueryBuilder.pb

// 	resp := new(proto.Response)

// 	err := node.invoke(methodName(pb), pb, resp)

// 	if err != nil {
// 		return newErrPingStatus(err)
// 	}

// 	respHeader := mapResponseHeader(resp)

// 	if respHeader.NodeTransactionPrecheckCode == proto.ResponseCodeEnum_BUSY {
// 		return newErrPingStatus(fmt.Errorf("%s", Status(respHeader.NodeTransactionPrecheckCode).String()))
// 	}

// 	if isResponseUnknown(resp) {
// 		return newErrPingStatus(fmt.Errorf("unknown"))
// 	}

// 	return nil
// }

func (client *Client) getNextNode() NodeID {
	nodeID := client.networkNodeIds[client.nextNodeIndex]
	client.nextNodeIndex = (client.nextNodeIndex + 1) % uint(len(client.networkNodeIds))

	return nodeID
}

func (client *Client) getNumberOfNodesForTransaction() int {
	return (len(client.networkNodeIds) + 3 - 1) / 3
}

func (client *Client) getChannel(id NodeID) (*channel, error) {
	if client.networkChannels[id.AccountID] != nil {
		return client.networkChannels[id.AccountID], nil
	}

	conn, err := grpc.Dial(id.Address, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	ch := newChannel(conn)
	client.networkChannels[id.AccountID] = &ch
	return &ch, nil
}
