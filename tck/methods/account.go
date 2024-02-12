package methods

import (
	"context"
	"fmt"
	"github.com/creachadair/jrpc2"
	"github.com/hashgraph/hedera-sdk-go/v2"
	"time"
)

type AccountService struct {
	sdkService *SDKService
	hedera.AccountID
}

type CreateAccountParams struct {
	PublicKey                     string        `json:"publicKey"`
	InitialBalance                int64         `json:"initialBalance"`
	ReceiverSignatureRequired     bool          `json:"receiverSignatureRequired"`
	MaxAutomaticTokenAssociations uint32        `json:"maxAutomaticTokenAssociations"`
	StakedAccountId               string        `json:"stakedAccountId"`
	StakedNodeId                  int64         `json:"stakedNodeId"`
	DeclineStakingReward          bool          `json:"declineStakingReward"`
	AccountMemo                   string        `json:"accountMemo"`
	AutoRenewPeriod               time.Duration `json:"autoRenewPeriod"`
	PrivateKey                    string        `json:"privateKey"`
}
type AccountFromAliasParams struct {
	OperatorId     string `json:"operator_id"`
	AliasAccountId string `json:"aliasAccountId"`
	InitialBalance int64  `json:"initialBalance"`
}
type UpdateAccountKeyParams struct {
	AccountId     string `json:"accountId"`
	NewPublicKey  string `json:"newPublicKey"`
	OldPrivateKey string `json:"oldPrivateKey"`
	NewPrivateKey string `json:"newPrivateKey"`
}
type UpdateAccountMemoParams struct {
	AccountId string `json:"accountId"`
	Key       string `json:"key"`
	Memo      string `json:"memo"`
}
type DeleteAccountParams struct {
	AccountId   string `json:"accountId"`
	AccountKey  string `json:"accountKey"`
	RecipientId string `json:"recipientId"`
}

// NOTE: Structs are created with uppercase, otherwise they are private and JSON serializer cannot encode/decode them.
// That's why we set `json:"{expected field name}"`, which is lowercase.

// TODO: Move this struct in dedicated file for hepler structs, or something like this. We want to use it for all methods
type ErrorStatus struct {
	Status string `json:"status"`
}

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

var threeSec = time.Second * 3

// SetSdkService We set object, which is holding our client params. Pass it by referance, because TCK is dynamically updating it
func (a *AccountService) SetSdkService(service *SDKService) {
	a.sdkService = service
}

// TODO Find better way to handle error. We want some wrapper, which would intercept the error and contruct the correct
// TODO jrpc error based on that  "jrpc2.Errorf(-32603, "Internal error"/"Hedera error", ErrorStatus{Status: err.Error()})"
// CreateAccount jRPC method for createAccount
func (a *AccountService) CreateAccount(_ context.Context, accountParams CreateAccountParams) (*AccountResponse, error) {
	transaction := hedera.NewAccountCreateTransaction().SetGrpcDeadline(&threeSec)

	if accountParams.PublicKey != "" {
		key, err := hedera.PublicKeyFromString(accountParams.PublicKey)
		if err != nil {
			return nil, jrpc2.Errorf(-32603, "Internal error", ErrorStatus{Status: err.Error()})
		}
		transaction.SetKey(key)
	}
	if accountParams.InitialBalance != 0 {
		transaction.SetInitialBalance(hedera.HbarFromTinybar(accountParams.InitialBalance))
	}
	if accountParams.ReceiverSignatureRequired {
		transaction.SetReceiverSignatureRequired(accountParams.ReceiverSignatureRequired)
	}
	if accountParams.MaxAutomaticTokenAssociations != 0 {
		transaction.SetMaxAutomaticTokenAssociations(accountParams.MaxAutomaticTokenAssociations)
	}
	if accountParams.StakedAccountId != "" {
		accountId, err := hedera.AccountIDFromString(accountParams.StakedAccountId)
		if err != nil {
			return nil, jrpc2.Errorf(-32603, "Internal error", ErrorStatus{Status: err.Error()})
		}
		transaction.SetStakedAccountID(accountId)
	}
	if accountParams.StakedNodeId != 0 {
		transaction.SetStakedNodeID(accountParams.StakedNodeId)
	}
	if accountParams.DeclineStakingReward {
		transaction.SetDeclineStakingReward(accountParams.DeclineStakingReward)
	}
	if accountParams.AccountMemo != "" {
		transaction.SetAccountMemo(accountParams.AccountMemo)
	}
	if accountParams.AutoRenewPeriod != 0 {
		transaction.SetAutoRenewPeriod(accountParams.AutoRenewPeriod)
	}
	if accountParams.PrivateKey != "" {
		key, err := hedera.PrivateKeyFromString(accountParams.PrivateKey)
		if err != nil {
			return nil, jrpc2.Errorf(-32603, "Internal error", ErrorStatus{Status: err.Error()})
		}

		_, err = transaction.FreezeWith(a.sdkService.Client)
		if err != nil {
			return nil, jrpc2.Errorf(-32603, err.Error())
		}
		transaction.Sign(key)
	}

	txResponse, err := transaction.Execute(a.sdkService.Client)
	if err != nil {
		fmt.Println(err)
		return nil, jrpc2.Errorf(-32603, "Internal error").WithData(ErrorStatus{Status: err.Error()})
	}
	receipt, err := txResponse.GetReceipt(a.sdkService.Client)
	var accId string
	if receipt.Status == hedera.StatusSuccess {
		accId = receipt.AccountID.String()
	}

	return &AccountResponse{AccountId: accId, Status: receipt.Status.String()}, nil
}

// CreateAccountFromAlias Create an account from aliasId by transferring  some HBAR amount from given account (opperatorId)
func (a *AccountService) CreateAccountFromAlias(_ context.Context, fromAliasParams AccountFromAliasParams) (*AccountResponse, error) {
	operator, err := hedera.AccountIDFromString(fromAliasParams.OperatorId)
	alias, err := hedera.AccountIDFromString(fromAliasParams.AliasAccountId)

	if err != nil {
		return nil, jrpc2.Errorf(-32603, "Internal error", ErrorStatus{Status: err.Error()})
	}

	resp, err := hedera.NewTransferTransaction().SetGrpcDeadline(&threeSec).
		AddHbarTransfer(operator, hedera.NewHbar(float64(fromAliasParams.InitialBalance)).Negated()).
		AddHbarTransfer(alias, hedera.NewHbar(float64(fromAliasParams.InitialBalance))).
		Execute(a.sdkService.Client)

	if err != nil {
		fmt.Println(err)
		return nil, jrpc2.Errorf(-32603, "Internal error").WithData(ErrorStatus{Status: err.Error()})
	}
	receipt, err := resp.GetReceipt(a.sdkService.Client)
	var accId string
	if receipt.Status == hedera.StatusSuccess {
		accId = receipt.AccountID.String()
	}

	return &AccountResponse{AccountId: accId, Status: receipt.Status.String()}, nil
}

// GetAccountInfo Get info for a given accountId and return a custom object containing aggregated information about the account
func (a *AccountService) GetAccountInfo(_ context.Context, accountId string) (*AccountInfoResponse, error) {
	account, _ := hedera.AccountIDFromString(accountId)

	resp, err := hedera.NewAccountInfoQuery().SetAccountID(account).SetGrpcDeadline(&threeSec).Execute(a.sdkService.Client)
	if err != nil {
		return nil, jrpc2.Errorf(-32603, "Internal error", ErrorStatus{Status: err.Error()})
	}

	return &AccountInfoResponse{AccountID: resp.AccountID.String(),
		Balance:                       resp.Balance.String(),
		Key:                           resp.Key.String(),
		AccountMemo:                   resp.AccountMemo,
		MaxAutomaticTokenAssociations: resp.MaxAutomaticTokenAssociations,
		AutoRenewPeriod:               resp.AutoRenewPeriod}, nil
}

// UpdateAccountKey updates an existing acoount id keys with provided params
func (a *AccountService) UpdateAccountKey(_ context.Context, params UpdateAccountKeyParams) (*AccountResponse, error) {
	accId, _ := hedera.AccountIDFromString(params.AccountId)
	key, _ := hedera.PublicKeyFromString(params.NewPublicKey)

	tx, err := hedera.NewAccountUpdateTransaction().SetGrpcDeadline(&threeSec).SetAccountID(accId).SetKey(key).FreezeWith(a.sdkService.Client)

	if err != nil {
		return nil, jrpc2.Errorf(-32603, "Internal error", ErrorStatus{Status: err.Error()})
	}
	oldKey, _ := hedera.PrivateKeyFromString(params.OldPrivateKey)
	newKey, _ := hedera.PrivateKeyFromString(params.NewPrivateKey)
	tx.Sign(oldKey)
	tx.Sign(newKey)

	resp, err := tx.Execute(a.sdkService.Client)
	if err != nil {
		return nil, jrpc2.Errorf(-32603, "Internal error", ErrorStatus{Status: err.Error()})
	}

	receipt, err := resp.GetReceipt(a.sdkService.Client)
	if err != nil {
		return nil, jrpc2.Errorf(-32603, "Internal error", ErrorStatus{Status: err.Error()})
	}

	return &AccountResponse{Status: receipt.Status.String()}, nil
}

// UpdateAccountMemo updates account memo of an existing account ID
func (a *AccountService) UpdateAccountMemo(_ context.Context, params UpdateAccountMemoParams) (*AccountResponse, error) {
	accId, _ := hedera.AccountIDFromString(params.AccountId)
	tx, err := hedera.NewAccountUpdateTransaction().SetGrpcDeadline(&threeSec).SetAccountID(accId).SetAccountMemo(params.Memo).FreezeWith(a.sdkService.Client)
	if err != nil {
		return nil, jrpc2.Errorf(-32603, "Internal error", ErrorStatus{Status: err.Error()})
	}
	signature, _ := hedera.PrivateKeyFromString(params.Key)

	resp, err := tx.Sign(signature).Execute(a.sdkService.Client)
	if err != nil {
		return nil, jrpc2.Errorf(-32603, "Internal error", ErrorStatus{Status: err.Error()})
	}

	receipt, _ := resp.GetReceipt(a.sdkService.Client)
	return &AccountResponse{Status: receipt.Status.String()}, nil
}

// DeleteAccount deletes a provided account by signing the transaction with the key of that account
func (a *AccountService) DeleteAccount(_ context.Context, params DeleteAccountParams) (*AccountResponse, error) {
	accId, _ := hedera.AccountIDFromString(params.AccountId)
	recipientId, _ := hedera.AccountIDFromString(params.RecipientId)
	tx, err := hedera.NewAccountDeleteTransaction().SetGrpcDeadline(&threeSec).SetAccountID(accId).SetTransferAccountID(recipientId).FreezeWith(a.sdkService.Client)
	if err != nil {
		return nil, jrpc2.Errorf(-32603, "Internal error", ErrorStatus{Status: err.Error()})
	}
	signature, _ := hedera.PrivateKeyFromString(params.AccountKey)

	resp, err := tx.Sign(signature).Execute(a.sdkService.Client)
	if err != nil {
		return nil, jrpc2.Errorf(-32603, "Internal error", ErrorStatus{Status: err.Error()})
	}

	receipt, _ := resp.GetReceipt(a.sdkService.Client)
	return &AccountResponse{Status: receipt.Status.String()}, nil
}
