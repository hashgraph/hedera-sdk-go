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

	operator *_Operator

	network                         _Network
	mirrorNetwork                   *_MirrorNetwork
	autoValidateChecksums           bool
	defaultRegenerateTransactionIDs bool
	maxAttempts                     *int

	maxBackoff time.Duration
	minBackoff time.Duration

	requestTimeout *time.Duration
}

// TransactionSigner is a closure or function that defines how transactions will be signed
type TransactionSigner func(message []byte) []byte

type _Operator struct {
	accountID  AccountID
	privateKey *PrivateKey
	publicKey  PublicKey
	signer     TransactionSigner
}

var mainnetNodes = map[string]AccountID{
	"35.237.200.180:50211":  {Account: 3},
	"34.239.82.6:50211":     {Account: 3},
	"13.82.40.153:50211":    {Account: 3},
	"13.124.142.126:50211":  {Account: 3},
	"15.164.44.66:50211":    {Account: 3},
	"15.165.118.251:50211":  {Account: 3},
	"35.186.191.247:50211":  {Account: 4},
	"3.130.52.236:50211":    {Account: 4},
	"137.116.36.18:50211":   {Account: 4},
	"35.192.2.25:50211":     {Account: 5},
	"3.18.18.254:50211":     {Account: 5},
	"104.43.194.202:50211":  {Account: 5},
	"23.111.186.250:50211":  {Account: 5},
	"74.50.117.35:50211":    {Account: 5},
	"107.155.64.98:50211":   {Account: 5},
	"35.199.161.108:50211":  {Account: 6},
	"13.52.108.243:50211":   {Account: 6},
	"13.64.151.232:50211":   {Account: 6},
	"13.235.15.32:50211":    {Account: 6},
	"104.211.205.124:50211": {Account: 6},
	"13.71.90.154:50211":    {Account: 6},
	"35.203.82.240:50211":   {Account: 7},
	"3.114.54.4:50211":      {Account: 7},
	"23.102.74.34:50211":    {Account: 7},
	"35.236.5.219:50211":    {Account: 8},
	"35.183.66.150:50211":   {Account: 8},
	"23.96.185.18:50211":    {Account: 8},
	"35.197.192.225:50211":  {Account: 9},
	"35.181.158.250:50211":  {Account: 9},
	"23.97.237.125:50211":   {Account: 9},
	"31.214.8.131:50211":    {Account: 9},
	"35.242.233.154:50211":  {Account: 10},
	"3.248.27.48:50211":     {Account: 10},
	"65.52.68.254:50211":    {Account: 10},
	"179.190.33.184:50211":  {Account: 10},
	"35.240.118.96:50211":   {Account: 11},
	"13.53.119.185:50211":   {Account: 11},
	"23.97.247.27:50211":    {Account: 11},
	"69.87.222.61:50211":    {Account: 11},
	"96.126.72.172:50211":   {Account: 11},
	"69.87.221.231:50211":   {Account: 11},
	"35.204.86.32:50211":    {Account: 12},
	"35.177.162.180:50211":  {Account: 12},
	"51.140.102.228:50211":  {Account: 12},
	"35.234.132.107:50211":  {Account: 13},
	"34.215.192.104:50211":  {Account: 13},
	"13.77.158.252:50211":   {Account: 13},
	"35.236.2.27:50211":     {Account: 14},
	"52.8.21.141:50211":     {Account: 14},
	"40.114.107.85:50211":   {Account: 14},
	"35.228.11.53:50211":    {Account: 15},
	"3.121.238.26:50211":    {Account: 15},
	"40.89.139.247:50211":   {Account: 15},
	"34.91.181.183:50211":   {Account: 16},
	"18.157.223.230:50211":  {Account: 16},
	"13.69.120.73:50211":    {Account: 16},
	"50.7.176.235:50211":    {Account: 16},
	"198.16.99.40:50211":    {Account: 16},
	"50.7.124.46:50211":     {Account: 16},
	"34.86.212.247:50211":   {Account: 17},
	"18.232.251.19:50211":   {Account: 17},
	"40.114.92.39:50211":    {Account: 17},
	"172.105.247.67:50211":  {Account: 18},
	"172.104.150.132:50211": {Account: 18},
	"139.162.156.222:50211": {Account: 18},
	"34.89.87.138:50211":    {Account: 19},
	"18.168.4.59:50211":     {Account: 19},
	"51.140.43.81:50211":    {Account: 19},
	"34.82.78.255:50211":    {Account: 20},
	"13.77.151.212:50211":   {Account: 20},
	"34.76.140.109:50211":   {Account: 21},
	"13.36.123.209:50211":   {Account: 21},
	"34.64.141.166:50211":   {Account: 22},
	"52.78.202.34:50211":    {Account: 22},
	"35.232.244.145:50211":  {Account: 23},
	"3.18.91.176:50211":     {Account: 23},
	"34.89.103.38:50211":    {Account: 24},
	"18.135.7.211:50211":    {Account: 24},
	"34.93.112.7:50211":     {Account: 25},
	"13.232.240.207:50211":  {Account: 25},
	"34.87.150.174:50211":   {Account: 25},
	"13.228.103.14:50211":   {Account: 25},
}

var testnetNodes = map[string]AccountID{
	"0.testnet.hedera.com:50211": {Account: 3},
	"34.94.106.61:50211":         {Account: 3},
	"50.18.132.211:50211":        {Account: 3},
	"138.91.142.219:50211":       {Account: 3},
	"1.testnet.hedera.com:50211": {Account: 4},
	"35.237.119.55:50211":        {Account: 4},
	"3.212.6.13:50211":           {Account: 4},
	"52.168.76.241:50211":        {Account: 4},
	"2.testnet.hedera.com:50211": {Account: 5},
	"35.245.27.193:50211":        {Account: 5},
	"52.20.18.86:50211":          {Account: 5},
	"40.79.83.124:50211":         {Account: 5},
	"3.testnet.hedera.com:50211": {Account: 6},
	"34.83.112.116:50211":        {Account: 6},
	"54.70.192.33:50211":         {Account: 6},
	"52.183.45.65:50211":         {Account: 6},
	"4.testnet.hedera.com:50211": {Account: 7},
	"34.94.160.4:50211":          {Account: 7},
	"54.176.199.109:50211":       {Account: 7},
	"13.64.181.136:50211":        {Account: 7},
	"5.testnet.hedera.com:50211": {Account: 8},
	"34.106.102.218:50211":       {Account: 8},
	"35.155.49.147:50211":        {Account: 8},
	"13.78.238.32:50211":         {Account: 8},
	"6.testnet.hedera.com:50211": {Account: 9},
	"34.133.197.230:50211":       {Account: 9},
	"52.14.252.207:50211":        {Account: 9},
	"52.165.17.231:50211":        {Account: 9},
}

var previewnetNodes = map[string]AccountID{
	"0.previewnet.hedera.com:50211": {Account: 3},
	"35.231.208.148:50211":          {Account: 3},
	"3.211.248.172:50211":           {Account: 3},
	"40.121.64.48:50211":            {Account: 3},
	"1.previewnet.hedera.com:50211": {Account: 4},
	"35.199.15.177:50211":           {Account: 4},
	"3.133.213.146:50211":           {Account: 4},
	"40.70.11.202:50211":            {Account: 4},
	"2.previewnet.hedera.com:50211": {Account: 5},
	"35.225.201.195:50211":          {Account: 5},
	"52.15.105.130:50211":           {Account: 5},
	"104.43.248.63:50211":           {Account: 5},
	"3.previewnet.hedera.com:50211": {Account: 6},
	"35.247.109.135:50211":          {Account: 6},
	"54.241.38.1:50211":             {Account: 6},
	"13.88.22.47:50211":             {Account: 6},
	"4.previewnet.hedera.com:50211": {Account: 7},
	"35.235.65.51:50211":            {Account: 7},
	"54.177.51.127:50211":           {Account: 7},
	"13.64.170.40:50211":            {Account: 7},
	"5.previewnet.hedera.com:50211": {Account: 8},
	"34.106.247.65:50211":           {Account: 8},
	"35.83.89.171:50211":            {Account: 8},
	"13.78.232.192:50211":           {Account: 8},
	"6.previewnet.hedera.com:50211": {Account: 9},
	"34.125.23.49:50211":            {Account: 9},
	"50.18.17.93:50211":             {Account: 9},
	"20.150.136.89:50211":           {Account: 9},
}

var mainnetMirror = []string{"mainnet-public.mirrornode.hedera.com:5600"}
var testnetMirror = []string{"hcs.testnet.mirrornode.hedera.com:5600"}
var previewnetMirror = []string{"hcs.previewnet.mirrornode.hedera.com:5600"}

func ClientForNetwork(network map[string]AccountID) *Client {
	return _NewClient(network, []string{}, "mainnet")
}

// ClientForMainnet returns a preconfigured client for use with the standard
// Hedera mainnet.
// Most users will want to set an _Operator account with .SetOperator so
// transactions can be automatically given TransactionIDs and signed.
func ClientForMainnet() *Client {
	return _NewClient(mainnetNodes, mainnetMirror, NetworkNameMainnet)
}

// ClientForTestnet returns a preconfigured client for use with the standard
// Hedera testnet.
// Most users will want to set an _Operator account with .SetOperator so
// transactions can be automatically given TransactionIDs and signed.
func ClientForTestnet() *Client {
	return _NewClient(testnetNodes, testnetMirror, NetworkNameTestnet)
}

// ClientForPreviewnet returns a preconfigured client for use with the standard
// Hedera previewnet.
// Most users will want to set an _Operator account with .SetOperator so
// transactions can be automatically given TransactionIDs and signed.
func ClientForPreviewnet() *Client {
	return _NewClient(previewnetNodes, previewnetMirror, NetworkNamePreviewnet)
}

// newClient takes in a map of _Node addresses to their respective IDS (_Network)
// and returns a Client instance which can be used to
func _NewClient(network map[string]AccountID, mirrorNetwork []string, name NetworkName) *Client {
	client := Client{
		maxQueryPayment:                 defaultMaxQueryPayment,
		maxTransactionFee:               defaultMaxTransactionFee,
		network:                         _NewNetwork(),
		mirrorNetwork:                   _NewMirrorNetwork(),
		autoValidateChecksums:           false,
		maxAttempts:                     nil,
		minBackoff:                      250 * time.Millisecond,
		maxBackoff:                      8 * time.Second,
		defaultRegenerateTransactionIDs: true,
	}

	_ = client.SetNetwork(network)
	client.SetMirrorNetwork(mirrorNetwork)
	client.network._SetNetworkName(name)

	return &client
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
			case string(NetworkNameMainnet):
				network = mainnetNodes
			case string(NetworkNamePreviewnet):
				network = previewnetNodes
			case string(NetworkNameTestnet):
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

	configBytes, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	return ClientFromConfig(configBytes)
}

// Close is used to disconnect the Client from the _Network
func (client *Client) Close() error {
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

// Deprecated use SetNodeMinBackoff
func (client *Client) SetNodeWaitTime(nodeWait time.Duration) {
	client.network._SetNodeMinBackoff(nodeWait)
}

// Deprecated use GetNodeMinBackoff
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
	client.mirrorNetwork._SetTransportSecurity(tls)

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

// Deprecated: Use SetLedgerID instead
func (client *Client) SetLedgerID(id LedgerID) {
	client.network._SetLedgerID(id)
}

// Deprecated: Use SetLedgerID instead
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
