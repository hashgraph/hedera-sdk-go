package methods

// SPDX-License-Identifier: Apache-2.0

import (
	"context"
	"time"

	"github.com/hiero-ledger/hiero-sdk-go/tck/param"
	"github.com/hiero-ledger/hiero-sdk-go/tck/response"
	"github.com/hiero-ledger/hiero-sdk-go/tck/utils"
	"github.com/hiero-ledger/hiero-sdk-go/v2"
)

// ---- Struct to hold hiero.Client implementation and to implement the methods of the specification ----
type AccountService struct {
	sdkService *SDKService
}

// Variable to be set to `SetGrpcDeadline` for all transactions
var threeSecondsDuration = time.Second * 3

// SetSdkService We set object, which is holding our client param. Pass it by referance, because TCK is dynamically updating it
func (a *AccountService) SetSdkService(service *SDKService) {
	a.sdkService = service
}

// CreateAccount jRPC method for createAccount
func (a *AccountService) CreateAccount(_ context.Context, params param.CreateAccountParams) (*response.AccountResponse, error) {
	transaction := hiero.NewAccountCreateTransaction().SetGrpcDeadline(&threeSecondsDuration)

	if params.Key != nil {
		key, err := utils.GetKeyFromString(*params.Key)
		if err != nil {
			return nil, err
		}
		transaction.SetKey(key)
	}
	if params.InitialBalance != nil {
		transaction.SetInitialBalance(hiero.HbarFromTinybar(*params.InitialBalance))
	}
	if params.ReceiverSignatureRequired != nil {
		transaction.SetReceiverSignatureRequired(*params.ReceiverSignatureRequired)
	}
	if params.MaxAutomaticTokenAssociations != nil {
		transaction.SetMaxAutomaticTokenAssociations(*params.MaxAutomaticTokenAssociations)
	}
	if params.StakedAccountId != nil {
		accountId, err := hiero.AccountIDFromString(*params.StakedAccountId)
		if err != nil {
			return nil, err
		}
		transaction.SetStakedAccountID(accountId)
	}
	if params.StakedNodeId != nil {
		stakedNodeID, err := params.StakedNodeId.Int64()
		if err != nil {
			return nil, response.InvalidParams.WithData(err.Error())
		}
		transaction.SetStakedNodeID(stakedNodeID)
	}
	if params.DeclineStakingReward != nil {
		transaction.SetDeclineStakingReward(*params.DeclineStakingReward)
	}
	if params.Memo != nil {
		transaction.SetAccountMemo(*params.Memo)
	}
	if params.AutoRenewPeriod != nil {
		transaction.SetAutoRenewPeriod(time.Duration(*params.AutoRenewPeriod) * time.Second)
	}
	if params.Alias != nil {
		transaction.SetAlias(*params.Alias)
	}
	if params.CommonTransactionParams != nil {
		params.CommonTransactionParams.FillOutTransaction(transaction, a.sdkService.Client)
	}
	txResponse, err := transaction.Execute(a.sdkService.Client)
	if err != nil {
		return nil, err
	}
	receipt, err := txResponse.GetReceipt(a.sdkService.Client)
	if err != nil {
		return nil, err
	}
	var accId string
	if receipt.Status == hiero.StatusSuccess {
		accId = receipt.AccountID.String()
	}
	return &response.AccountResponse{AccountId: accId, Status: receipt.Status.String()}, nil
}

// UpdateAccount jRPC method for updateAccount
func (a *AccountService) UpdateAccount(_ context.Context, params param.UpdateAccountParams) (*response.AccountResponse, error) {
	transaction := hiero.NewAccountUpdateTransaction().SetGrpcDeadline(&threeSecondsDuration)
	if params.AccountId != nil {
		accountId, _ := hiero.AccountIDFromString(*params.AccountId)
		transaction.SetAccountID(accountId)
	}

	if params.Key != nil {
		key, err := utils.GetKeyFromString(*params.Key)
		if err != nil {
			return nil, err
		}
		transaction.SetKey(key)
	}

	if params.ExpirationTime != nil {
		transaction.SetExpirationTime(time.Unix(*params.ExpirationTime, 0))
	}

	if params.ReceiverSignatureRequired != nil {
		transaction.SetReceiverSignatureRequired(*params.ReceiverSignatureRequired)
	}

	if params.MaxAutomaticTokenAssociations != nil {
		transaction.SetMaxAutomaticTokenAssociations(*params.MaxAutomaticTokenAssociations)
	}

	if params.StakedAccountId != nil {
		accountId, err := hiero.AccountIDFromString(*params.StakedAccountId)
		if err != nil {
			return nil, err
		}
		transaction.SetStakedAccountID(accountId)
	}

	if params.StakedNodeId != nil {
		stakedNodeID, err := params.StakedNodeId.Int64()
		if err != nil {
			return nil, response.InvalidParams.WithData(err.Error())
		}
		transaction.SetStakedNodeID(stakedNodeID)
	}

	if params.DeclineStakingReward != nil {
		transaction.SetDeclineStakingReward(*params.DeclineStakingReward)
	}

	if params.Memo != nil {
		transaction.SetAccountMemo(*params.Memo)
	}

	if params.AutoRenewPeriod != nil {
		transaction.SetAutoRenewPeriod(time.Duration(*params.AutoRenewPeriod) * time.Second)
	}

	if params.CommonTransactionParams != nil {
		params.CommonTransactionParams.FillOutTransaction(transaction, a.sdkService.Client)
	}

	txResponse, err := transaction.Execute(a.sdkService.Client)
	if err != nil {
		return nil, err
	}
	receipt, err := txResponse.GetReceipt(a.sdkService.Client)
	if err != nil {
		return nil, err
	}
	return &response.AccountResponse{Status: receipt.Status.String()}, nil
}

// DeleteAccount jRPC method for deleteAccount
func (a *AccountService) DeleteAccount(_ context.Context, params param.DeleteAccountParams) (*response.AccountResponse, error) {
	transaction := hiero.NewAccountDeleteTransaction().SetGrpcDeadline(&threeSecondsDuration)
	if params.DeleteAccountId != nil {
		accountId, _ := hiero.AccountIDFromString(*params.DeleteAccountId)
		transaction.SetAccountID(accountId)
	}

	if params.TransferAccountId != nil {
		accountId, _ := hiero.AccountIDFromString(*params.TransferAccountId)
		transaction.SetTransferAccountID(accountId)
	}

	if params.CommonTransactionParams != nil {
		params.CommonTransactionParams.FillOutTransaction(transaction, a.sdkService.Client)
	}

	txResponse, err := transaction.Execute(a.sdkService.Client)
	if err != nil {
		return nil, err
	}
	receipt, err := txResponse.GetReceipt(a.sdkService.Client)
	if err != nil {
		return nil, err
	}
	return &response.AccountResponse{Status: receipt.Status.String()}, nil
}
