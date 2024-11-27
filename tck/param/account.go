package param

// SPDX-License-Identifier: Apache-2.0

import (
	"encoding/json"
)

type CreateAccountParams struct {
	Key                           *string                  `json:"key"`
	InitialBalance                *int64                   `json:"initialBalance"`
	ReceiverSignatureRequired     *bool                    `json:"receiverSignatureRequired"`
	AutoRenewPeriod               *int64                   `json:"autoRenewPeriod"`
	Memo                          *string                  `json:"memo"`
	MaxAutomaticTokenAssociations *int32                   `json:"maxAutoTokenAssociations"`
	StakedAccountId               *string                  `json:"stakedAccountId"`
	StakedNodeId                  *json.Number             `json:"stakedNodeId"`
	DeclineStakingReward          *bool                    `json:"declineStakingReward"`
	Alias                         *string                  `json:"alias"`
	CommonTransactionParams       *CommonTransactionParams `json:"commonTransactionParams"`
}

type UpdateAccountParams struct {
	AccountId                     *string                  `json:"accountId"`
	Key                           *string                  `json:"key"`
	ReceiverSignatureRequired     *bool                    `json:"receiverSignatureRequired"`
	AutoRenewPeriod               *int64                   `json:"autoRenewPeriod"`
	ExpirationTime                *int64                   `json:"expirationTime"`
	Memo                          *string                  `json:"memo"`
	MaxAutomaticTokenAssociations *int32                   `json:"maxAutoTokenAssociations"`
	StakedAccountId               *string                  `json:"stakedAccountId"`
	StakedNodeId                  *json.Number             `json:"stakedNodeId"`
	DeclineStakingReward          *bool                    `json:"declineStakingReward"`
	CommonTransactionParams       *CommonTransactionParams `json:"commonTransactionParams"`
}
type DeleteAccountParams struct {
	DeleteAccountId         *string                  `json:"deleteAccountId"`
	TransferAccountId       *string                  `json:"transferAccountId"`
	CommonTransactionParams *CommonTransactionParams `json:"commonTransactionParams"`
}
