package response

import "time"

// ---- Structs to hold different response structures expected from jRPC client ----

type AccountResponse struct {
	AccountId string `json:"accountId"`
	Status    string `json:"status"`
}
type AccountInfoResponse struct {
	AccountID                     string
	Balance                       string
	Key                           string
	AccountMemo                   string
	MaxAutomaticTokenAssociations uint32
	AutoRenewPeriod               time.Duration
}
