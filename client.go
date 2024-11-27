package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"context"
	_ "embed"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"time"
)

//go:embed addressbook/mainnet.pb
var mainnetAddress []byte
var mainnetNodes, _ = NodeAddressBookFromBytes(mainnetAddress)

//go:embed addressbook/previewnet.pb
var previewnetAddress []byte
var previewnetNodes, _ = NodeAddressBookFromBytes(previewnetAddress)

//go:embed addressbook/testnet.pb
var testnetAddress []byte
var testnetNodes, _ = NodeAddressBookFromBytes(testnetAddress)

// Client is the Hiero protocol wrapper for the SDK used by all
// transaction and query types.
type Client struct {
	defaultMaxTransactionFee Hbar
	defaultMaxQueryPayment   Hbar

	operator *_Operator

	network                         _Network
	mirrorNetwork                   *_MirrorNetwork
	autoValidateChecksums           bool
	defaultRegenerateTransactionIDs bool
	maxAttempts                     *int

	maxBackoff time.Duration
	minBackoff time.Duration

	requestTimeout             *time.Duration
	defaultNetworkUpdatePeriod time.Duration
	networkUpdateContext       context.Context
	cancelNetworkUpdate        context.CancelFunc
	logger                     Logger
}

// TransactionSigner is a closure or function that defines how transactions will be signed
type TransactionSigner func(message []byte) []byte

type _Operator struct {
	accountID  AccountID
	privateKey *PrivateKey
	publicKey  PublicKey
	signer     TransactionSigner
}

var mainnetMirror = []string{"mainnet-public.mirrornode.hedera.com:443"}
var testnetMirror = []string{"testnet.mirrornode.hedera.com:443"}
var previewnetMirror = []string{"previewnet.mirrornode.hedera.com:443"}

// ClientForMirrorNetwork constructs a client given a set of mirror network nodes.
func ClientForMirrorNetwork(mirrorNetwork []string) (*Client, error) {
	net := _NewNetwork()
	client := _NewClient(net, mirrorNetwork, nil, true)
	addressbook, err := NewAddressBookQuery().
		SetFileID(FileIDForAddressBook()).
		Execute(client)
	if err != nil {
		return nil, fmt.Errorf("failed to query address book: %v", err)
	}
	client.SetNetworkFromAddressBook(addressbook)
	return client, nil
}

// ClientForNetwork constructs a client given a set of nodes.
func ClientForNetwork(network map[string]AccountID) *Client {
	net := _NewNetwork()
	client := _NewClient(net, []string{}, nil, true)
	_ = client.SetNetwork(network)
	net._SetLedgerID(*NewLedgerIDMainnet())
	return client
}

// ClientForMainnet returns a preconfigured client for use with the standard
// Hiero mainnet.
// Most users will want to set an _Operator account with .SetOperator so
// transactions can be automatically given TransactionIDs and signed.
func ClientForMainnet() *Client {
	return _NewClient(*_NetworkForMainnet(mainnetNodes._ToMap()), mainnetMirror, NewLedgerIDMainnet(), true)
}

// ClientForTestnet returns a preconfigured client for use with the standard
// Hiero testnet.
// Most users will want to set an _Operator account with .SetOperator so
// transactions can be automatically given TransactionIDs and signed.
func ClientForTestnet() *Client {
	return _NewClient(*_NetworkForTestnet(testnetNodes._ToMap()), testnetMirror, NewLedgerIDTestnet(), true)
}

// ClientForPreviewnet returns a preconfigured client for use with the standard
// Hiero previewnet.
// Most users will want to set an _Operator account with .SetOperator so
// transactions can be automatically given TransactionIDs and signed.
func ClientForPreviewnet() *Client {
	return _NewClient(*_NetworkForPreviewnet(previewnetNodes._ToMap()), previewnetMirror, NewLedgerIDPreviewnet(), true)
}

// newClient takes in a map of _Node addresses to their respective IDS (_Network)
// and returns a Client instance which can be used to
func _NewClient(network _Network, mirrorNetwork []string, ledgerId *LedgerID, shouldScheduleNetworkUpdate bool) *Client {
	ctx, cancel := context.WithCancel(context.Background())
	logger := NewLogger("hiero-sdk-go", LogLevel(os.Getenv("HEDERA_SDK_GO_LOG_LEVEL")))
	var defaultLogger Logger = logger

	client := Client{
		defaultMaxQueryPayment:          NewHbar(1),
		network:                         network,
		mirrorNetwork:                   _NewMirrorNetwork(),
		autoValidateChecksums:           false,
		maxAttempts:                     nil,
		minBackoff:                      250 * time.Millisecond,
		maxBackoff:                      8 * time.Second,
		defaultRegenerateTransactionIDs: true,
		defaultNetworkUpdatePeriod:      24 * time.Hour,
		networkUpdateContext:            ctx,
		cancelNetworkUpdate:             cancel,
		logger:                          defaultLogger,
	}

	client.SetMirrorNetwork(mirrorNetwork)
	if ledgerId != nil {
		client.SetLedgerID(*ledgerId)
	}

	// We can't ask for AddressBook from non existent Mirror node
	if len(mirrorNetwork) > 0 && shouldScheduleNetworkUpdate {
		// Update the Addressbook, before the default timeout starts
		client._UpdateAddressBook()
		go client._ScheduleNetworkUpdate(ctx, client.defaultNetworkUpdatePeriod)
	}

	return &client
}

func (client *Client) _UpdateAddressBook() {
	addressbook, err := NewAddressBookQuery().
		SetFileID(FileIDForAddressBook()).
		Execute(client)
	if err == nil && len(addressbook.NodeAddresses) > 0 {
		client.SetNetworkFromAddressBook(addressbook)
	}
}

func (client *Client) _ScheduleNetworkUpdate(ctx context.Context, duration time.Duration) {
	for {
		select {
		case <-ctx.Done():
			return
		case <-time.After(duration):
			client._UpdateAddressBook()
		}
	}
}

// CancelScheduledNetworkUpdate cancels the scheduled network update the network address book
func (client *Client) CancelScheduledNetworkUpdate() {
	client.cancelNetworkUpdate()
}

// SetNetworkUpdatePeriod sets how often the client will update the network address book
func (client *Client) SetNetworkUpdatePeriod(period time.Duration) *Client {
	client.defaultNetworkUpdatePeriod = period
	client.CancelScheduledNetworkUpdate()
	client.networkUpdateContext, client.cancelNetworkUpdate = context.WithCancel(context.Background())
	go client._ScheduleNetworkUpdate(client.networkUpdateContext, period)
	return client
}

// GetNetworkUpdatePeriod returns the current network update period
func (client *Client) GetNetworkUpdatePeriod() time.Duration {
	return client.defaultNetworkUpdatePeriod
}

// ClientForName set up the client for the selected network.
func ClientForName(name string) (*Client, error) {
	switch name {
	case string(NetworkNameTestnet):
		return ClientForTestnet(), nil
	case string(NetworkNamePreviewnet):
		return ClientForPreviewnet(), nil
	case string(NetworkNameMainnet):
		return ClientForMainnet(), nil
	case "local", "localhost":
		network := make(map[string]AccountID)
		network["127.0.0.1:50213"] = AccountID{Account: 3}
		mirror := []string{"127.0.0.1:5600"}
		client := ClientForNetwork(network)
		client.SetMirrorNetwork(mirror)
		return client, nil
	default:
		return &Client{}, fmt.Errorf("%q is not recognized as a valid Hiero _Network", name)
	}
}

type _ConfigOperator struct {
	AccountID  string `json:"accountId"`
	PrivateKey string `json:"privateKey"`
}

// TODO: Implement complete spec: https://gitlab.com/launchbadge/hedera/sdk/python/-/issues/45
type _ClientConfig struct {
	Network       interface{}      `json:"network"`
	MirrorNetwork interface{}      `json:"mirrorNetwork"`
	Operator      *_ConfigOperator `json:"operator"`
}

// ClientFromConfig takes in the byte slice representation of a JSON string or
// document and returns Client based on the configuration.
func ClientFromConfig(jsonBytes []byte) (*Client, error) {
	return clientFromConfig(jsonBytes, true)
}

// ClientFromConfigWithoutScheduleNetworkUpdate does not schedule network update
// the user has to call SetNetworkUpdatePeriod manually
func ClientFromConfigWithoutScheduleNetworkUpdate(jsonBytes []byte) (*Client, error) {
	return clientFromConfig(jsonBytes, false)
}

func clientFromConfig(jsonBytes []byte, shouldScheduleNetworkUpdate bool) (*Client, error) {
	var clientConfig _ClientConfig
	var client *Client

	err := json.Unmarshal(jsonBytes, &clientConfig)
	if err != nil {
		return nil, err
	}

	network := _NewNetwork()
	networkAddresses := make(map[string]AccountID)

	switch net := clientConfig.Network.(type) {
	case map[string]interface{}:
		for url, inter := range net {
			switch id := inter.(type) {
			case string:
				accountID, err := AccountIDFromString(id)
				if err != nil {
					return client, err
				}
				networkAddresses[url] = accountID
			default:
				return client, errors.New("network is expected to be map of string to string, or string")
			}
		}
		err = network.SetNetwork(networkAddresses)
		if err != nil {
			return &Client{}, err
		}
	case string:
		if len(net) > 0 {
			switch net {
			case string(NetworkNameMainnet):
				network = *_NetworkForMainnet(mainnetNodes._ToMap())
			case string(NetworkNamePreviewnet):
				network = *_NetworkForPreviewnet(previewnetNodes._ToMap())
			case string(NetworkNameTestnet):
				network = *_NetworkForTestnet(testnetNodes._ToMap())
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
		client = _NewClient(network, arr, nil, shouldScheduleNetworkUpdate)
	case string:
		if len(mirror) > 0 {
			switch mirror {
			case string(NetworkNameMainnet):
				client = _NewClient(network, mainnetMirror, NewLedgerIDMainnet(), shouldScheduleNetworkUpdate)
			case string(NetworkNameTestnet):
				client = _NewClient(network, testnetMirror, NewLedgerIDTestnet(), shouldScheduleNetworkUpdate)
			case string(NetworkNamePreviewnet):
				client = _NewClient(network, previewnetMirror, NewLedgerIDPreviewnet(), shouldScheduleNetworkUpdate)
			}
		}
	case nil:
		client = _NewClient(network, []string{}, nil, true)
	default:
		return client, errors.New("mirrorNetwork is expected to be a string, an array of strings or nil")
	}

	// if the _Operator is not provided, finish here
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

	operator := _Operator{
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

	configBytes, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	return ClientFromConfig(configBytes)
}

// Close is used to disconnect the Client from the _Network
func (client *Client) Close() error {
	client.CancelScheduledNetworkUpdate()
	err := client.network._Close()
	if err != nil {
		return err
	}
	err = client.mirrorNetwork._Close()
	if err != nil {
		return err
	}

	return nil
}

// SetNetwork replaces all nodes in this Client with a new set of nodes.
func (client *Client) SetNetwork(network map[string]AccountID) error {
	return client.network.SetNetwork(network)
}

// GetNetwork returns the current set of nodes in this Client.
func (client *Client) GetNetwork() map[string]AccountID {
	return client.network._GetNetwork()
}

// SetMaxNodeReadmitTime The maximum amount of time to wait before attempting to
// reconnect to a node that has been removed from the network.
func (client *Client) SetMaxNodeReadmitTime(readmitTime time.Duration) {
	client.network._SetMaxNodeReadmitPeriod(readmitTime)
}

// GetMaxNodeReadmitTime returns the maximum amount of time to wait before attempting to
// reconnect to a node that has been removed from the network.
func (client *Client) GetMaxNodeReadmitPeriod() time.Duration {
	return client.network._GetMaxNodeReadmitPeriod()
}

// SetMinNodeReadmitTime The minimum amount of time to wait before attempting to
// reconnect to a node that has been removed from the network.
func (client *Client) SetMinNodeReadmitTime(readmitTime time.Duration) {
	client.network._SetMinNodeReadmitPeriod(readmitTime)
}

// GetMinNodeReadmitTime returns the minimum amount of time to wait before attempting to
// reconnect to a node that has been removed from the network.
func (client *Client) GetMinNodeReadmitPeriod() time.Duration {
	return client.network._GetMinNodeReadmitPeriod()
}

// SetMaxBackoff The maximum amount of time to wait between retries.
// Every retry attempt will increase the wait time exponentially until it reaches this time.
func (client *Client) SetMaxBackoff(max time.Duration) {
	if max.Nanoseconds() < 0 {
		panic("maxBackoff must be a positive duration")
	} else if max.Nanoseconds() < client.minBackoff.Nanoseconds() {
		panic("maxBackoff must be greater than or equal to minBackoff")
	}
	client.maxBackoff = max
}

// GetMaxBackoff returns the maximum amount of time to wait between retries.
func (client *Client) GetMaxBackoff() time.Duration {
	return client.maxBackoff
}

// SetMinBackoff sets the minimum amount of time to wait between retries.
func (client *Client) SetMinBackoff(min time.Duration) {
	if min.Nanoseconds() < 0 {
		panic("minBackoff must be a positive duration")
	} else if client.maxBackoff.Nanoseconds() < min.Nanoseconds() {
		panic("minBackoff must be less than or equal to maxBackoff")
	}
	client.minBackoff = min
}

// GetMinBackoff returns the minimum amount of time to wait between retries.
func (client *Client) GetMinBackoff() time.Duration {
	return client.minBackoff
}

// SetMaxAttempts sets the maximum number of times to attempt a transaction or query.
func (client *Client) SetMaxAttempts(max int) {
	client.maxAttempts = &max
}

// GetMaxAttempts returns the maximum number of times to attempt a transaction or query.
func (client *Client) GetMaxAttempts() int {
	if client.maxAttempts == nil {
		return -1
	}

	return *client.maxAttempts
}

// SetMaxNodeAttempts sets the maximum number of times to attempt a transaction or query on a single node.
func (client *Client) SetMaxNodeAttempts(max int) {
	client.network._SetMaxNodeAttempts(max)
}

// GetMaxNodeAttempts returns the maximum number of times to attempt a transaction or query on a single node.
func (client *Client) GetMaxNodeAttempts() int {
	return client.network._GetMaxNodeAttempts()
}

// Deprecated: use SetNodeMinBackoff
func (client *Client) SetNodeWaitTime(nodeWait time.Duration) {
	client.network._SetNodeMinBackoff(nodeWait)
}

// Deprecated: use GetNodeMinBackoff
func (client *Client) GetNodeWaitTime() time.Duration {
	return client.network._GetNodeMinBackoff()
}

// SetNodeMinBackoff sets the minimum amount of time to wait between retries on a single node.
func (client *Client) SetNodeMinBackoff(nodeWait time.Duration) {
	client.network._SetNodeMinBackoff(nodeWait)
}

// GetNodeMinBackoff returns the minimum amount of time to wait between retries on a single node.
func (client *Client) GetNodeMinBackoff() time.Duration {
	return client.network._GetNodeMinBackoff()
}

// SetNodeMaxBackoff sets the maximum amount of time to wait between retries on a single node.
func (client *Client) SetNodeMaxBackoff(nodeWait time.Duration) {
	client.network._SetNodeMaxBackoff(nodeWait)
}

// GetNodeMaxBackoff returns the maximum amount of time to wait between retries on a single node.
func (client *Client) GetNodeMaxBackoff() time.Duration {
	return client.network._GetNodeMaxBackoff()
}

// SetMaxNodesPerTransaction sets the maximum number of nodes to try for a single transaction.
func (client *Client) SetMaxNodesPerTransaction(max int) {
	client.network._SetMaxNodesPerTransaction(max)
}

// SetNetwork replaces all _Nodes in the Client with a new set of _Nodes.
// (e.g. for an Address Book update).
func (client *Client) SetMirrorNetwork(mirrorNetwork []string) {
	_ = client.mirrorNetwork._SetNetwork(mirrorNetwork)
}

// GetNetwork returns the mirror network node list.
func (client *Client) GetMirrorNetwork() []string {
	return client.mirrorNetwork._GetNetwork()
}

// SetTransportSecurity sets if transport security should be used to connect to consensus nodes.
// If transport security is enabled all connections to consensus nodes will use TLS, and
// the server's certificate hash will be compared to the hash stored in the NodeAddressBook
// for the given network.
// *Note*: If transport security is enabled, but {@link Client#isVerifyCertificates()} is disabled
// then server certificates will not be verified.
func (client *Client) SetTransportSecurity(tls bool) *Client {
	client.network._SetTransportSecurity(tls)

	return client
}

// SetCertificateVerification sets if server certificates should be verified against an existing address book.
func (client *Client) SetCertificateVerification(verify bool) *Client {
	client.network._SetVerifyCertificate(verify)

	return client
}

// GetCertificateVerification returns if server certificates should be verified against an existing address book.
func (client *Client) GetCertificateVerification() bool {
	return client.network._GetVerifyCertificate()
}

// Deprecated: Use SetLedgerID instead
func (client *Client) SetNetworkName(name NetworkName) {
	ledgerID, _ := LedgerIDFromNetworkName(name)
	client.SetLedgerID(*ledgerID)
}

// Deprecated: Use GetLedgerID instead
func (client *Client) GetNetworkName() *NetworkName {
	name, _ := client.GetLedgerID().ToNetworkName()
	return &name
}

// SetLedgerID sets the ledger ID for the Client.
func (client *Client) SetLedgerID(id LedgerID) {
	client.network._SetLedgerID(id)
}

// GetLedgerID returns the ledger ID for the Client.
func (client *Client) GetLedgerID() *LedgerID {
	return client.network._GetLedgerID()
}

// SetAutoValidateChecksums sets if an automatic entity ID checksum validation should be performed.
func (client *Client) SetAutoValidateChecksums(validate bool) {
	client.autoValidateChecksums = validate
}

// GetAutoValidateChecksums returns if an automatic entity ID checksum validation should be performed.
func (client *Client) GetAutoValidateChecksums() bool {
	return client.autoValidateChecksums
}

// SetDefaultRegenerateTransactionIDs sets if an automatic transaction ID regeneration should be performed.
func (client *Client) SetDefaultRegenerateTransactionIDs(regen bool) {
	client.defaultRegenerateTransactionIDs = regen
}

// GetDefaultRegenerateTransactionIDs returns if an automatic transaction ID regeneration should be performed.
func (client *Client) GetDefaultRegenerateTransactionIDs() bool {
	return client.defaultRegenerateTransactionIDs
}

// SetNodeMinReadmitPeriod sets the minimum amount of time to wait before attempting to
// reconnect to a node that has been removed from the network.
func (client *Client) SetNodeMinReadmitPeriod(period time.Duration) {
	client.network._SetNodeMinReadmitPeriod(period)
}

// SetNodeMaxReadmitPeriod sets the maximum amount of time to wait before attempting to
// reconnect to a node that has been removed from the network.
func (client *Client) SetNodeMaxReadmitPeriod(period time.Duration) {
	client.network._SetNodeMaxReadmitPeriod(period)
}

// GetNodeMinReadmitPeriod returns the minimum amount of time to wait before attempting to
// reconnect to a node that has been removed from the network.
func (client *Client) GetNodeMinReadmitPeriod() time.Duration {
	return client.network._GetNodeMinReadmitPeriod()
}

// GetNodeMaxReadmitPeriod returns the maximum amount of time to wait before attempting to
// reconnect to a node that has been removed from the network.
func (client *Client) GetNodeMaxReadmitPeriod() time.Duration {
	return client.network._GetNodeMaxReadmitPeriod()
}

// SetOperator sets that account that will, by default, be paying for
// transactions and queries built with the client and the associated key
// with which to automatically sign transactions.
func (client *Client) SetOperator(accountID AccountID, privateKey PrivateKey) *Client {
	client.operator = &_Operator{
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
	client.operator = &_Operator{
		accountID:  accountID,
		privateKey: nil,
		publicKey:  publicKey,
		signer:     signer,
	}

	return client
}

// SetRequestTimeout sets the timeout for all requests made by the client.
func (client *Client) SetRequestTimeout(timeout *time.Duration) {
	client.requestTimeout = timeout
}

// GetRequestTimeout returns the timeout for all requests made by the client.
func (client *Client) GetRequestTimeout() *time.Duration {
	return client.requestTimeout
}

// GetOperatorAccountID returns the ID for the _Operator
func (client *Client) GetOperatorAccountID() AccountID {
	if client.operator != nil {
		return client.operator.accountID
	}

	return AccountID{}
}

// GetOperatorPublicKey returns the Key for the _Operator
func (client *Client) GetOperatorPublicKey() PublicKey {
	if client.operator != nil {
		return client.operator.publicKey
	}

	return PublicKey{}
}

// Ping sends an AccountBalanceQuery to the specified _Node returning nil if no
// problems occur. Otherwise, an error representing the status of the _Node will
// be returned.
func (client *Client) Ping(nodeID AccountID) error {
	_, err := NewAccountBalanceQuery().
		SetNodeAccountIDs([]AccountID{nodeID}).
		SetAccountID(client.GetOperatorAccountID()).
		Execute(client)

	return err
}

func (client *Client) PingAll() {
	for _, s := range client.GetNetwork() {
		_ = client.Ping(s)
	}
}

// SetNetworkFromAddressBook replaces all nodes in this Client with the nodes in the Address Book.
func (client *Client) SetNetworkFromAddressBook(addressBook NodeAddressBook) *Client {
	client.network._SetNetworkFromAddressBook(addressBook)
	return client
}

// SetDefaultMaxQueryPayment sets the default maximum payment allowed for queries.
func (client *Client) SetDefaultMaxQueryPayment(defaultMaxQueryPayment Hbar) error {
	if defaultMaxQueryPayment.AsTinybar() < 0 {
		return errors.New("DefaultMaxQueryPayment must be non-negative")
	}

	client.defaultMaxQueryPayment = defaultMaxQueryPayment
	return nil
}

// GetDefaultMaxQueryPayment returns the default maximum payment allowed for queries.
func (client *Client) GetDefaultMaxQueryPayment() Hbar {
	return client.defaultMaxQueryPayment
}

// SetDefaultMaxTransactionFee sets the default maximum fee allowed for transactions.
func (client *Client) SetDefaultMaxTransactionFee(defaultMaxTransactionFee Hbar) error {
	if defaultMaxTransactionFee.AsTinybar() < 0 {
		return errors.New("DefaultMaxTransactionFee must be non-negative")
	}

	client.defaultMaxTransactionFee = defaultMaxTransactionFee
	return nil
}

// GetDefaultMaxTransactionFee returns the default maximum fee allowed for transactions.
func (client *Client) GetDefaultMaxTransactionFee() Hbar {
	return client.defaultMaxTransactionFee
}

func (client *Client) SetLogger(logger Logger) *Client {
	client.logger = logger
	return client
}

func (client *Client) GetLogger() Logger {
	return client.logger
}

func (client *Client) SetLogLevel(level LogLevel) *Client {
	client.logger.SetLevel(level)
	return client
}
