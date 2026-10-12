package response

// SPDX-License-Identifier: Apache-2.0

type AccountResponse struct {
	AccountId     string `json:"accountId"`
	Status        string `json:"status"`
	TransactionId string `json:"transactionId,omitempty"`
}

type AccountInfoResponse struct {
	AccountId                     string                           `json:"accountId"`
	ContractAccountId             string                           `json:"contractAccountId"`
	IsDeleted                     bool                             `json:"isDeleted"`
	ProxyAccountId                *string                          `json:"proxyAccountId"`
	ProxyReceived                 string                           `json:"proxyReceived"`
	Key                           string                           `json:"key"`
	Balance                       string                           `json:"balance"`
	SendRecordThreshold           string                           `json:"sendRecordThreshold"`
	ReceiveRecordThreshold        string                           `json:"receiveRecordThreshold"`
	IsReceiverSignatureRequired   bool                             `json:"isReceiverSignatureRequired"`
	ExpirationTime                string                           `json:"expirationTime"`
	AutoRenewPeriod               string                           `json:"autoRenewPeriod"`
	LiveHashes                    []LiveHashResponse               `json:"liveHashes"`
	TokenRelationships            map[string]TokenRelationshipInfo `json:"tokenRelationships"`
	AccountMemo                   string                           `json:"accountMemo"`
	OwnedNfts                     string                           `json:"ownedNfts"`
	MaxAutomaticTokenAssociations string                           `json:"maxAutomaticTokenAssociations"`
	AliasKey                      *string                          `json:"aliasKey"`
	LedgerId                      string                           `json:"ledgerId"`
	HbarAllowances                []HbarAllowanceResponse          `json:"hbarAllowances"`
	TokenAllowances               []TokenAllowanceResponse         `json:"tokenAllowances"`
	NftAllowances                 []TokenNftAllowanceResponse      `json:"nftAllowances"`
	EthereumNonce                 string                           `json:"ethereumNonce"`
	StakingInfo                   *StakingInfoResponse             `json:"stakingInfo"`
}

type LiveHashResponse struct {
	AccountId string   `json:"accountId"`
	Hash      string   `json:"hash"`
	Keys      []string `json:"keys"`
	Duration  string   `json:"duration"`
}

type TokenRelationshipInfo struct {
	TokenId              string `json:"tokenId"`
	Symbol               string `json:"symbol"`
	Balance              string `json:"balance"`
	IsKycGranted         *bool  `json:"isKycGranted"`
	IsFrozen             *bool  `json:"isFrozen"`
	AutomaticAssociation bool   `json:"automaticAssociation"`
}

type HbarAllowanceResponse struct {
	OwnerAccountId   string `json:"ownerAccountId"`
	SpenderAccountId string `json:"spenderAccountId"`
	Amount           string `json:"amount"`
}

type TokenAllowanceResponse struct {
	TokenId          string `json:"tokenId"`
	OwnerAccountId   string `json:"ownerAccountId"`
	SpenderAccountId string `json:"spenderAccountId"`
	Amount           string `json:"amount"`
}

type TokenNftAllowanceResponse struct {
	TokenId           string   `json:"tokenId"`
	OwnerAccountId    string   `json:"ownerAccountId"`
	SpenderAccountId  string   `json:"spenderAccountId"`
	SerialNumbers     []string `json:"serialNumbers"`
	AllSerials        bool     `json:"allSerials"`
	DelegatingSpender string   `json:"delegatingSpender"`
}

type StakingInfoResponse struct {
	DeclineStakingReward bool    `json:"declineStakingReward"`
	StakePeriodStart     *string `json:"stakePeriodStart"`
	PendingReward        string  `json:"pendingReward"`
	StakedToMe           string  `json:"stakedToMe"`
	StakedAccountId      *string `json:"stakedAccountId"`
	StakedNodeId         *string `json:"stakedNodeId"`
}

type TransactionReceiptResponse struct {
	Status                 string                       `json:"status"`
	AccountId              *string                      `json:"accountId"`
	FileId                 *string                      `json:"fileId"`
	ContractId             *string                      `json:"contractId"`
	TopicId                *string                      `json:"topicId"`
	TokenId                *string                      `json:"tokenId"`
	ScheduleId             *string                      `json:"scheduleId"`
	ExchangeRate           *ExchangeRate                `json:"exchangeRate"`
	TopicSequenceNumber    *string                      `json:"topicSequenceNumber"`
	TopicRunningHash       *string                      `json:"topicRunningHash"`
	TotalSupply            *string                      `json:"totalSupply"`
	ScheduledTransactionId *string                      `json:"scheduledTransactionId"`
	Serials                []string                     `json:"serials"`
	Duplicates             []TransactionReceiptResponse `json:"duplicates"`
	Children               []TransactionReceiptResponse `json:"children"`
	NodeId                 *string                      `json:"nodeId"`
}

type ExchangeRate struct {
	Hbars          int32   `json:"hbars"`
	Cents          int32   `json:"cents"`
	ExpirationTime *string `json:"expirationTime"`
}

type MirrorNodeAccountBalanceResponse struct {
	Hbar string `json:"hbars"`
}

type DeprecatedAccountBalanceQueryResponse struct {
	ConstructionWarning *string `json:"constructionWarning"`
	ExecutionError      *string `json:"executionError"`
}
