package hedera

/*-
 *
 * Hedera Go SDK
 *
 * Copyright (C) 2020 - 2022 Hedera Hashgraph, LLC
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

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

// Client is the Hedera protocol wrapper for the SDK used by all
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

func ClientForNetwork(network map[string]AccountID) *Client {
	net := _NewNetwork()
	client := _NewClient(net, []string{}, "mainnet")
	_ = client.SetNetwork(network)
	net._SetLedgerID(*NewLedgerIDMainnet())
	return client
}

// ClientForMainnet returns a preconfigured client for use with the standard
// Hedera mainnet.
// Most users will want to set an _Operator account with .SetOperator so
// transactions can be automatically given TransactionIDs and signed.
func ClientForMainnet() *Client {
	return _NewClient(*_NetworkForMainnet(mainnetNodes._ToMap()), mainnetMirror, NetworkNameMainnet)
}

// ClientForTestnet returns a preconfigured client for use with the standard
// Hedera testnet.
// Most users will want to set an _Operator account with .SetOperator so
// transactions can be automatically given TransactionIDs and signed.
func ClientForTestnet() *Client {
	return _NewClient(*_NetworkForTestnet(testnetNodes._ToMap()), testnetMirror, NetworkNameTestnet)
}

// ClientForPreviewnet returns a preconfigured client for use with the standard
// Hedera previewnet.
// Most users will want to set an _Operator account with .SetOperator so
// transactions can be automatically given TransactionIDs and signed.
func ClientForPreviewnet() *Client {
	return _NewClient(*_NetworkForPreviewnet(previewnetNodes._ToMap()), previewnetMirror, NetworkNamePreviewnet)
}

// newClient takes in a map of _Node addresses to their respective IDS (_Network)
// and returns a Client instance which can be used to
func _NewClient(network _Network, mirrorNetwork []string, name NetworkName) *Client {
	ctx, cancel := context.WithCancel(context.Background())
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
	}

	client.SetMirrorNetwork(mirrorNetwork)
	client.network._SetNetworkName(name)

	// We can't ask for AddressBook from non existent Mirror node
	if len(mirrorNetwork) > 0 {
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

func (client *Client) CancelScheduledNetworkUpdate() {
	client.cancelNetworkUpdate()
}

func (client *Client) SetNetworkUpdatePeriod(period time.Duration) *Client {
	client.defaultNetworkUpdatePeriod = period
	client.CancelScheduledNetworkUpdate()
	client.networkUpdateContext, client.cancelNetworkUpdate = context.WithCancel(context.Background())
	go client._ScheduleNetworkUpdate(client.networkUpdateContext, period)
	return client
}

func (client *Client) GetNetworkUpdatePeriod() time.Duration {
	return client.defaultNetworkUpdatePeriod
}

func ClientForName(name string) (*Client, error) {
	switch name {
	case string(NetworkNameTestnet):
		return ClientForTestnet(), nil
	case string(NetworkNamePreviewnet):
		return ClientForPreviewnet(), nil
	case string(NetworkNameMainnet):
		return ClientForMainnet(), nil
	default:
		return &Client{}, fmt.Errorf("%q is not recognized as a valid Hedera _Network", name)
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
		client = _NewClient(network, arr, NetworkNameMainnet)
	case string:
		if len(mirror) > 0 {
			switch mirror {
			case string(NetworkNameMainnet):
				client = _NewClient(network, mainnetMirror, NetworkNameMainnet)
			case string(NetworkNamePreviewnet):
				client = _NewClient(network, previewnetMirror, NetworkNamePreviewnet)
			case string(NetworkNameTestnet):
				client = _NewClient(network, testnetMirror, NetworkNameTestnet)
			}
		}
	default:
		return client, errors.New("mirrorNetwork is expected to be either string or an array of strings")
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

func (client *Client) SetNetwork(network map[string]AccountID) error {
	return client.network.SetNetwork(network)
}

func (client *Client) GetNetwork() map[string]AccountID {
	return client.network._GetNetwork()
}

func (client *Client) SetMaxNodeReadmitTime(readmitTime time.Duration) {
	client.network._SetMaxNodeReadmitPeriod(readmitTime)
}

func (client *Client) GetMaxNodeReadmitPeriod() time.Duration {
	return client.network._GetMaxNodeReadmitPeriod()
}

func (client *Client) SetMinNodeReadmitTime(readmitTime time.Duration) {
	client.network._SetMinNodeReadmitPeriod(readmitTime)
}

func (client *Client) GetMinNodeReadmitPeriod() time.Duration {
	return client.network._GetMinNodeReadmitPeriod()
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
	return client.maxBackoff
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
	client.network._SetMaxNodeAttempts(max)
}

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

func (client *Client) SetNodeMinBackoff(nodeWait time.Duration) {
	client.network._SetNodeMinBackoff(nodeWait)
}

func (client *Client) GetNodeMinBackoff() time.Duration {
	return client.network._GetNodeMinBackoff()
}

func (client *Client) SetNodeMaxBackoff(nodeWait time.Duration) {
	client.network._SetNodeMaxBackoff(nodeWait)
}

func (client *Client) GetNodeMaxBackoff() time.Duration {
	return client.network._GetNodeMaxBackoff()
}

func (client *Client) SetMaxNodesPerTransaction(max int) {
	client.network._SetMaxNodesPerTransaction(max)
}

// SetNetwork replaces all _Nodes in the Client with a new set of _Nodes.
// (e.g. for an Address Book update).
func (client *Client) SetMirrorNetwork(mirrorNetwork []string) {
	_ = client.mirrorNetwork._SetNetwork(mirrorNetwork)
}

func (client *Client) GetMirrorNetwork() []string {
	return client.mirrorNetwork._GetNetwork()
}

func (client *Client) SetTransportSecurity(tls bool) *Client {
	client.network._SetTransportSecurity(tls)

	return client
}

func (client *Client) SetCertificateVerification(verify bool) *Client {
	client.network._SetVerifyCertificate(verify)

	return client
}

func (client *Client) GetCertificateVerification() bool {
	return client.network._GetVerifyCertificate()
}

// Deprecated: Use SetLedgerID instead
func (client *Client) SetNetworkName(name NetworkName) {
	client.network._SetNetworkName(name)
}

// Deprecated: Use GetLedgerID instead
func (client *Client) GetNetworkName() *NetworkName {
	return client.network._GetNetworkName()
}

func (client *Client) SetLedgerID(id LedgerID) {
	client.network._SetLedgerID(id)
}

func (client *Client) GetLedgerID() *LedgerID {
	return client.network._GetLedgerID()
}

func (client *Client) SetAutoValidateChecksums(validate bool) {
	client.autoValidateChecksums = validate
}

func (client *Client) GetAutoValidateChecksums() bool {
	return client.autoValidateChecksums
}

func (client *Client) SetDefaultRegenerateTransactionIDs(regen bool) {
	client.defaultRegenerateTransactionIDs = regen
}

func (client *Client) GetDefaultRegenerateTransactionIDs() bool {
	return client.defaultRegenerateTransactionIDs
}

func (client *Client) SetNodeMinReadmitPeriod(period time.Duration) {
	client.network._SetNodeMinReadmitPeriod(period)
}

func (client *Client) SetNodeMaxReadmitPeriod(period time.Duration) {
	client.network._SetNodeMaxReadmitPeriod(period)
}

func (client *Client) GetNodeMinReadmitPeriod() time.Duration {
	return client.network._GetNodeMinReadmitPeriod()
}

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

func (client *Client) SetRequestTimeout(timeout *time.Duration) {
	client.requestTimeout = timeout
}

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

func (client *Client) SetNetworkFromAddressBook(addressBook NodeAddressBook) *Client {
	client.network._SetNetworkFromAddressBook(addressBook)
	return client
}

func (client *Client) SetDefaultMaxQueryPayment(defaultMaxQueryPayment Hbar) error {
	if defaultMaxQueryPayment.AsTinybar() < 0 {
		return errors.New("DefaultMaxQueryPayment must be non-negative")
	}

	client.defaultMaxQueryPayment = defaultMaxQueryPayment
	return nil
}

func (client *Client) GetDefaultMaxQueryPayment() Hbar {
	return client.defaultMaxQueryPayment
}

func (client *Client) SetDefaultMaxTransactionFee(defaultMaxTransactionFee Hbar) error {
	if defaultMaxTransactionFee.AsTinybar() < 0 {
		return errors.New("DefaultMaxTransactionFee must be non-negative")
	}

	client.defaultMaxTransactionFee = defaultMaxTransactionFee
	return nil
}

func (client *Client) GetDefaultMaxTransactionFee() Hbar {
	return client.defaultMaxTransactionFee
}
