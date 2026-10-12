package param

// SPDX-License-Identifier: Apache-2.0

import (
	"encoding/json"
)

type CreateAccountParams struct {
	BaseParams
	Key                           *string      `json:"key"`
	InitialBalance                *string      `json:"initialBalance"`
	ReceiverSignatureRequired     *bool        `json:"receiverSignatureRequired"`
	AutoRenewPeriod               *string      `json:"autoRenewPeriod"`
	Memo                          *string      `json:"memo"`
	MaxAutomaticTokenAssociations *int32       `json:"maxAutoTokenAssociations"`
	StakedAccountId               *string      `json:"stakedAccountId"`
	StakedNodeId                  *json.Number `json:"stakedNodeId"`
	DeclineStakingReward          *bool        `json:"declineStakingReward"`
	Alias                         *string      `json:"alias"`
}

type UpdateAccountParams struct {
	BaseParams
	AccountId                     *string      `json:"accountId"`
	Key                           *string      `json:"key"`
	ReceiverSignatureRequired     *bool        `json:"receiverSignatureRequired"`
	AutoRenewPeriod               *string      `json:"autoRenewPeriod"`
	ExpirationTime                *string      `json:"expirationTime"`
	Memo                          *string      `json:"memo"`
	MaxAutomaticTokenAssociations *int32       `json:"maxAutoTokenAssociations"`
	StakedAccountId               *string      `json:"stakedAccountId"`
	StakedNodeId                  *json.Number `json:"stakedNodeId"`
	DeclineStakingReward          *bool        `json:"declineStakingReward"`
}
type DeleteAccountParams struct {
	BaseParams
	DeleteAccountId   *string `json:"deleteAccountId"`
	TransferAccountId *string `json:"transferAccountId"`
}

type AccountAllowanceApproveParams struct {
	BaseParams
	Allowances *[]AllowanceParams `json:"allowances,omitempty"`
}

type AccountAllowanceDeleteParams struct {
	BaseParams
	Allowances *[]DeleteAllowanceParams `json:"allowances,omitempty"`
}

type GetAccountInfoParams struct {
	BaseParams
	AccountId *string `json:"accountId"`
}

type GetTransactionReceiptParams struct {
	BaseParams
	TransactionID     *string `json:"transactionId"`
	IncludeDuplicates *bool   `json:"includeDuplicates"`
	IncludeChildren   *bool   `json:"includeChildren"`
	ValidateStatus    *bool   `json:"validateStatus"`
}

type GetMirrorNodeAccountBalanceParams struct {
	BaseParams
	AccountId *string `json:"accountId"`
}

type ExecuteDeprecatedAccountBalanceQueryParams struct {
	BaseParams
	AccountId *string `json:"accountId"`
	Operation *string `json:"operation"`
}
