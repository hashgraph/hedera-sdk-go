package hedera

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"time"
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

	network               network
	mirrorNetwork         *mirrorNetwork
	networkName           *NetworkName
	autoValidateChecksums bool
	maxAttempts           *int

	maxBackoff time.Duration
	minBackoff time.Duration
}

// TransactionSigner is a closure or function that defines how transactions will be signed
type TransactionSigner func(message []byte) []byte

type operator struct {
	accountID  AccountID
	privateKey *PrivateKey
	publicKey  PublicKey
	signer     TransactionSigner
}

var mainnetNodes = map[string]AccountID{
	"35.237.200.180:50211": {Account: 3},
	"35.186.191.247:50211": {Account: 4},
	"35.192.2.25:50211":    {Account: 5},
	"35.199.161.108:50211": {Account: 6},
	"35.203.82.240:50211":  {Account: 7},
	"35.236.5.219:50211":   {Account: 8},
	"35.197.192.225:50211": {Account: 9},
	"35.242.233.154:50211": {Account: 10},
	"35.240.118.96:50211":  {Account: 11},
	"35.204.86.32:50211":   {Account: 12},
	"35.234.132.107:50211": {Account: 13},
	"35.236.2.27:50211":    {Account: 14},
	"35.228.11.53:50211":   {Account: 15},
	"34.91.181.183:50211":  {Account: 16},
	"34.86.212.247:50211":  {Account: 17},
	"172.105.247.67:50211": {Account: 18},
	"34.89.87.138:50211":   {Account: 19},
	"34.82.78.255:50211":   {Account: 20},
}

var testnetNodes = map[string]AccountID{
	"0.testnet.hedera.com:50211": {Account: 3},
	"1.testnet.hedera.com:50211": {Account: 4},
	"2.testnet.hedera.com:50211": {Account: 5},
	"3.testnet.hedera.com:50211": {Account: 6},
	"4.testnet.hedera.com:50211": {Account: 7},
}

var previewnetNodes = map[string]AccountID{
	"0.previewnet.hedera.com:50211": {Account: 3},
	"1.previewnet.hedera.com:50211": {Account: 4},
	"2.previewnet.hedera.com:50211": {Account: 5},
	"3.previewnet.hedera.com:50211": {Account: 6},
	"4.previewnet.hedera.com:50211": {Account: 7},
}

var mainnetMirror = []string{"hcs.mainnet.mirrornode.hedera.com:5600"}
var testnetMirror = []string{"hcs.testnet.mirrornode.hedera.com:5600"}
var previewnetMirror = []string{"hcs.previewnet.mirrornode.hedera.com:5600"}

func ClientForNetwork(network map[string]AccountID) *Client {
	return newClient(network, []string{}, "mainnet")
}

// ClientForMainnet returns a preconfigured client for use with the standard
// Hedera mainnet.
// Most users will want to set an operator account with .SetOperator so
// transactions can be automatically given TransactionIDs and signed.
func ClientForMainnet() *Client {
	return newClient(mainnetNodes, mainnetMirror, NetworkNameMainnet)
}

// ClientForTestnet returns a preconfigured client for use with the standard
// Hedera testnet.
// Most users will want to set an operator account with .SetOperator so
// transactions can be automatically given TransactionIDs and signed.
func ClientForTestnet() *Client {
	return newClient(testnetNodes, testnetMirror, NetworkNameTestnet)
}

// ClientForPreviewnet returns a preconfigured client for use with the standard
// Hedera previewnet.
// Most users will want to set an operator account with .SetOperator so
// transactions can be automatically given TransactionIDs and signed.
func ClientForPreviewnet() *Client {
	return newClient(previewnetNodes, previewnetMirror, NetworkNamePreviewnet)
}

// newClient takes in a map of node addresses to their respective IDS (network)
// and returns a Client instance which can be used to
func newClient(network map[string]AccountID, mirrorNetwork []string, name NetworkName) *Client {
	client := Client{
		maxQueryPayment:       defaultMaxQueryPayment,
		maxTransactionFee:     defaultMaxTransactionFee,
		network:               newNetwork(),
		mirrorNetwork:         newMirrorNetwork(),
		networkName:           &name,
		autoValidateChecksums: false,
		maxAttempts:           nil,
		minBackoff:            250 * time.Millisecond,
		maxBackoff:            8 * time.Second,
	}

	_ = client.SetNetwork(network)
	client.SetMirrorNetwork(mirrorNetwork)

	return &client
}

func ClientForName(name string) (*Client, error) {
	switch name {
	case "testnet":
		return ClientForTestnet(), nil
	case "previewnet":
		return ClientForPreviewnet(), nil
	case "mainnet":
		return ClientForMainnet(), nil
	default:
		return &Client{}, fmt.Errorf("%q is not recognized as a valid Hedera network", name)
	}
}

type configOperator struct {
	AccountID  string `json:"accountId"`
	PrivateKey string `json:"privateKey"`
}

// TODO: Implement complete spec: https://gitlab.com/launchbadge/hedera/sdk/python/-/issues/45
type clientConfig struct {
	Network       interface{}     `json:"network"`
	MirrorNetwork interface{}     `json:"mirrorNetwork"`
	Operator      *configOperator `json:"operator"`
}

// ClientFromConfig takes in the byte slice representation of a JSON string or
// document and returns Client based on the configuration.
func ClientFromConfig(jsonBytes []byte) (*Client, error) {
	var clientConfig clientConfig
	var client *Client

	err := json.Unmarshal(jsonBytes, &clientConfig)
	if err != nil {
		return nil, err
	}

	network := make(map[string]AccountID)

	switch net := clientConfig.Network.(type) {
	case map[string]interface{}:
		for url, inter := range net {
			switch id := inter.(type) {
			case string:
				accountID, err := AccountIDFromString(id)
				if err != nil {
					return client, err
				}

				network[url] = accountID
			default:
				return client, errors.New("network is expected to be map of string to string, or string")
			}
		}
	case string:
		if len(net) > 0 {
			switch net {
			case "mainnet":
				network = mainnetNodes
			case "previewnet":
				network = previewnetNodes
			case "testnet":
				network = testnetNodes
			}
		}
	default:
		return client, errors.New("network is expected to be map of string to string, or string")
	}

	switch mirror := clientConfig.MirrorNetwork.(type) {
	case []interface{}:
		arr := make([]string, len(mirror))
		for i, inter := range mirror {
			switch str := inter.(type) {
			case string:
				arr[i] = str
			default:
				return client, errors.New("mirrorNetwork is expected to be either string or an array of strings")
			}
		}
		client = newClient(network, arr, NetworkNameMainnet)
	case string:
		if len(mirror) > 0 {
			switch mirror {
			case "mainnet":
				client = newClient(network, mainnetMirror, NetworkNameMainnet)
			case "previewnet":
				client = newClient(network, previewnetMirror, NetworkNamePreviewnet)
			case "testnet":
				client = newClient(network, testnetMirror, NetworkNameTestnet)
			}
		}
	default:
		return client, errors.New("mirrorNetwork is expected to be either string or an array of strings")
	}

	// if the operator is not provided, finish here
	if clientConfig.Operator == nil {
		return client, nil
	}

	operatorID, err := AccountIDFromString(clientConfig.Operator.AccountID)
	if err != nil {
		return client, err
	}

	operatorKey, err := PrivateKeyFromString(clientConfig.Operator.PrivateKey)
	if err != nil {
		return client, err
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

// ClientFromConfigFile takes a filename string representing the path to a JSON encoded
// Client file and returns a Client based on the configuration.
func ClientFromConfigFile(filename string) (*Client, error) {
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

	return ClientFromConfig(configBytes)
}

// Close is used to disconnect the Client from the network
func (client *Client) Close() error {
	client.network.Close()

	return nil
}

// SetNetwork replaces all nodes in the Client with a new set of nodes.
// (e.g. for an Address Book update).
func (client *Client) SetNetwork(network map[string]AccountID) error {
	return client.network.SetNetwork(network)
}

func (client *Client) GetNetwork() map[string]AccountID {
	return client.network.network
}

func (client *Client) SetMaxBackoff(max time.Duration) {
	if max.Nanoseconds() < 0 {
		panic("maxBackoff must be a positive duration")
	} else if max.Nanoseconds() < client.minBackoff.Nanoseconds() {
		panic("maxBackoff must be greater than or equal to minBackoff")
	}
	client.maxBackoff = max
}

func (client *Client) GetMaxBackoff() time.Duration {
	return client.GetMaxBackoff()
}

func (client *Client) SetMinBackoff(min time.Duration) {
	if min.Nanoseconds() < 0 {
		panic("minBackoff must be a positive duration")
	} else if client.maxBackoff.Nanoseconds() < min.Nanoseconds() {
		panic("minBackoff must be less than or equal to maxBackoff")
	}
	client.minBackoff = min
}

func (client *Client) GetMinBackoff() time.Duration {
	return client.minBackoff
}

func (client *Client) SetMaxAttempts(max int) {
	client.maxAttempts = &max
}

func (client *Client) GetMaxAttempts() int {
	if client.maxAttempts == nil {
		return -1
	}

	return *client.maxAttempts
}

func (client *Client) SetMaxNodeAttempts(max int) {
	client.network.setMaxNodeAttempts(max)
}

func (client *Client) GetMaxNodeAttempts() int {
	return client.network.getMaxNodeAttempts()
}

func (client *Client) SetNodeWaitTime(nodeWait time.Duration) {
	client.network.setNodeWaitTime(nodeWait)
}

func (client *Client) GetNodeWaitTime() time.Duration {
	return client.network.getNodeWaitTime()
}

func (client *Client) SetMaxNodesPerTransaction(max int) {
	client.network.setMaxNodesPerTransaction(max)
}

// SetNetwork replaces all nodes in the Client with a new set of nodes.
// (e.g. for an Address Book update).
func (client *Client) SetMirrorNetwork(mirrorNetwork []string) {
	client.mirrorNetwork.setNetwork(mirrorNetwork)
}

func (client *Client) GetMirrorNetwork() []string {
	return client.mirrorNetwork.network
}

func (client *Client) SetNetworkName(name NetworkName) {
	client.networkName = &name
}

func (client *Client) GetNetworkName() NetworkName {
	return *client.networkName
}

func (client *Client) SetAutoValidateChecksums(validate bool) {
	client.autoValidateChecksums = validate
}

func (client *Client) GetAutoValidateChecksums() bool {
	return client.autoValidateChecksums
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

// GetOperatorAccountID returns the ID for the operator
func (client *Client) GetOperatorAccountID() AccountID {
	if client.operator != nil {
		return client.operator.accountID
	} else {
		return AccountID{}
	}
}

// GetOperatorPublicKey returns the Key for the operator
func (client *Client) GetOperatorPublicKey() PublicKey {
	if client.operator != nil {
		return client.operator.publicKey
	} else {
		return PublicKey{}
	}
}

// Ping sends an AccountBalanceQuery to the specified node returning nil if no
// problems occur. Otherwise, an error representing the status of the node will
// be returned.
func (client *Client) Ping(nodeID AccountID) error {
	_, err := NewAccountBalanceQuery().
		SetAccountID(client.GetOperatorAccountID()).
		SetNodeAccountIDs([]AccountID{nodeID}).
		Execute(client)

	return err
}

func (client *Client) PingAll() {
	for _, s := range client.network.network {
		_ = client.Ping(s)
	}

	return
}
