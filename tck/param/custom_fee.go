package param

// SPDX-License-Identifier: Apache-2.0

type CustomFee struct {
	FeeCollectorAccountId string         `json:"feeCollectorAccountId"`
	FeeCollectorsExempt   *bool          `json:"feeCollectorsExempt"`
	FixedFee              *FixedFee      `json:"fixedFee,omitempty"`
	FractionalFee         *FractionalFee `json:"fractionalFee,omitempty"`
	RoyaltyFee            *RoyaltyFee    `json:"royaltyFee,omitempty"`
}

// FixedFee represents the fixed fee structure.
type FixedFee struct {
	Amount              string  `json:"amount"`
	DenominatingTokenId *string `json:"denominatingTokenId,omitempty"`
}

// FractionalFee represents the fractional fee structure.
type FractionalFee struct {
	Numerator        string `json:"numerator"`
	Denominator      string `json:"denominator"`
	MinimumAmount    string `json:"minimumAmount"`
	MaximumAmount    string `json:"maximumAmount"`
	AssessmentMethod string `json:"assessmentMethod"`
}

// RoyaltyFee represents the royalty fee structure.
type RoyaltyFee struct {
	Numerator   string    `json:"numerator"`
	Denominator string    `json:"denominator"`
	FallbackFee *FixedFee `json:"fallbackFee,omitempty"`
}
