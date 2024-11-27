package param

// SPDX-License-Identifier: Apache-2.0

type CreateTokenParams struct {
	Name                    *string                  `json:"name,omitempty"`
	Symbol                  *string                  `json:"symbol,omitempty"`
	Decimals                *int                     `json:"decimals,omitempty"`
	InitialSupply           *string                  `json:"initialSupply,omitempty"`
	TreasuryAccountId       *string                  `json:"treasuryAccountId,omitempty"`
	AdminKey                *string                  `json:"adminKey,omitempty"`
	KycKey                  *string                  `json:"kycKey,omitempty"`
	FreezeKey               *string                  `json:"freezeKey,omitempty"`
	WipeKey                 *string                  `json:"wipeKey,omitempty"`
	SupplyKey               *string                  `json:"supplyKey,omitempty"`
	FeeScheduleKey          *string                  `json:"feeScheduleKey,omitempty"`
	PauseKey                *string                  `json:"pauseKey,omitempty"`
	MetadataKey             *string                  `json:"metadataKey,omitempty"`
	FreezeDefault           *bool                    `json:"freezeDefault,omitempty"`
	ExpirationTime          *string                  `json:"expirationTime,omitempty"`
	AutoRenewAccountId      *string                  `json:"autoRenewAccountId,omitempty"`
	AutoRenewPeriod         *string                  `json:"autoRenewPeriod,omitempty"`
	Memo                    *string                  `json:"memo,omitempty"`
	TokenType               *string                  `json:"tokenType,omitempty"`
	SupplyType              *string                  `json:"supplyType,omitempty"`
	MaxSupply               *string                  `json:"maxSupply,omitempty"`
	CustomFees              *[]CustomFee             `json:"customFees,omitempty"`
	Metadata                *string                  `json:"metadata,omitempty"`
	CommonTransactionParams *CommonTransactionParams `json:"commonTransactionParams,omitempty"`
}
