package methods

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
func (a *AccountService) CreateAccount(_ context.Context, accountCreateParams param.CreateAccountParams) (*response.AccountResponse, error) {
	transaction := hiero.NewAccountCreateTransaction().SetGrpcDeadline(&threeSecondsDuration)

	if accountCreateParams.Key != nil {
		key, err := utils.GetKeyFromString(*accountCreateParams.Key)
		if err != nil {
			return nil, err
		}
		transaction.SetKey(key)
	}
	if accountCreateParams.InitialBalance != nil {
		transaction.SetInitialBalance(hiero.HbarFromTinybar(*accountCreateParams.InitialBalance))
	}
	if accountCreateParams.ReceiverSignatureRequired != nil {
		transaction.SetReceiverSignatureRequired(*accountCreateParams.ReceiverSignatureRequired)
	}
	if accountCreateParams.MaxAutomaticTokenAssociations != nil {
		transaction.SetMaxAutomaticTokenAssociations(*accountCreateParams.MaxAutomaticTokenAssociations)
	}
	if accountCreateParams.StakedAccountId != nil {
		accountId, err := hiero.AccountIDFromString(*accountCreateParams.StakedAccountId)
		if err != nil {
			return nil, err
		}
		transaction.SetStakedAccountID(accountId)
	}
	if accountCreateParams.StakedNodeId != nil {
		stakedNodeID, err := accountCreateParams.StakedNodeId.Int64()
		if err != nil {
			return nil, response.InvalidParams.WithData(err.Error())
		}
		transaction.SetStakedNodeID(stakedNodeID)
	}
	if accountCreateParams.DeclineStakingReward != nil {
		transaction.SetDeclineStakingReward(*accountCreateParams.DeclineStakingReward)
	}
	if accountCreateParams.Memo != nil {
		transaction.SetAccountMemo(*accountCreateParams.Memo)
	}
	if accountCreateParams.AutoRenewPeriod != nil {
		transaction.SetAutoRenewPeriod(time.Duration(*accountCreateParams.AutoRenewPeriod) * time.Second)
	}
	if accountCreateParams.Alias != nil {
		transaction.SetAlias(*accountCreateParams.Alias)
	}
	if accountCreateParams.CommonTransactionParams != nil {
		accountCreateParams.CommonTransactionParams.FillOutTransaction(transaction, a.sdkService.Client)
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
func (a *AccountService) UpdateAccount(_ context.Context, accountUpdateParams param.UpdateAccountParams) (*response.AccountResponse, error) {
	transaction := hiero.NewAccountUpdateTransaction().SetGrpcDeadline(&threeSecondsDuration)
	if accountUpdateParams.AccountId != nil {
		accountId, _ := hiero.AccountIDFromString(*accountUpdateParams.AccountId)
		transaction.SetAccountID(accountId)
	}

	if accountUpdateParams.Key != nil {
		key, err := utils.GetKeyFromString(*accountUpdateParams.Key)
		if err != nil {
			return nil, err
		}
		transaction.SetKey(key)
	}

	if accountUpdateParams.ExpirationTime != nil {
		transaction.SetExpirationTime(time.Unix(*accountUpdateParams.ExpirationTime, 0))
	}

	if accountUpdateParams.ReceiverSignatureRequired != nil {
		transaction.SetReceiverSignatureRequired(*accountUpdateParams.ReceiverSignatureRequired)
	}

	if accountUpdateParams.MaxAutomaticTokenAssociations != nil {
		transaction.SetMaxAutomaticTokenAssociations(*accountUpdateParams.MaxAutomaticTokenAssociations)
	}

	if accountUpdateParams.StakedAccountId != nil {
		accountId, err := hiero.AccountIDFromString(*accountUpdateParams.StakedAccountId)
		if err != nil {
			return nil, err
		}
		transaction.SetStakedAccountID(accountId)
	}

	if accountUpdateParams.StakedNodeId != nil {
		stakedNodeID, err := accountUpdateParams.StakedNodeId.Int64()
		if err != nil {
			return nil, response.InvalidParams.WithData(err.Error())
		}
		transaction.SetStakedNodeID(stakedNodeID)
	}

	if accountUpdateParams.DeclineStakingReward != nil {
		transaction.SetDeclineStakingReward(*accountUpdateParams.DeclineStakingReward)
	}

	if accountUpdateParams.Memo != nil {
		transaction.SetAccountMemo(*accountUpdateParams.Memo)
	}

	if accountUpdateParams.AutoRenewPeriod != nil {
		transaction.SetAutoRenewPeriod(time.Duration(*accountUpdateParams.AutoRenewPeriod) * time.Second)
	}

	if accountUpdateParams.CommonTransactionParams != nil {
		accountUpdateParams.CommonTransactionParams.FillOutTransaction(transaction, a.sdkService.Client)
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
func (a *AccountService) DeleteAccount(_ context.Context, deleteAccountParams param.DeleteAccountParams) (*response.AccountResponse, error) {
	transaction := hiero.NewAccountDeleteTransaction().SetGrpcDeadline(&threeSecondsDuration)
	if deleteAccountParams.DeleteAccountId != nil {
		accountId, _ := hiero.AccountIDFromString(*deleteAccountParams.DeleteAccountId)
		transaction.SetAccountID(accountId)
	}

	if deleteAccountParams.TransferAccountId != nil {
		accountId, _ := hiero.AccountIDFromString(*deleteAccountParams.TransferAccountId)
		transaction.SetTransferAccountID(accountId)
	}

	if deleteAccountParams.CommonTransactionParams != nil {
		deleteAccountParams.CommonTransactionParams.FillOutTransaction(transaction, a.sdkService.Client)
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
