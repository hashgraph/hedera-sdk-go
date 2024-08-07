package methods

import (
	"context"
	"strconv"
	"time"

	"github.com/hashgraph/hedera-sdk-go/tck/param"
	"github.com/hashgraph/hedera-sdk-go/tck/response"
	"github.com/hashgraph/hedera-sdk-go/tck/utils"
	"github.com/hashgraph/hedera-sdk-go/v2"
)

// ---- Struct to hold hedera.Client implementation and to implement the methods of the specification ----
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
	transaction := hedera.NewAccountCreateTransaction().SetGrpcDeadline(&threeSecondsDuration)

	if accountCreateParams.Key != nil {
		key, err := utils.GetKeyFromString(*accountCreateParams.Key)
		if err != nil {
			return nil, err
		}
		transaction.SetKey(key)
	}
	if accountCreateParams.InitialBalance != nil {
		transaction.SetInitialBalance(hedera.HbarFromTinybar(*accountCreateParams.InitialBalance))
	}
	if accountCreateParams.ReceiverSignatureRequired != nil {
		transaction.SetReceiverSignatureRequired(*accountCreateParams.ReceiverSignatureRequired)
	}
	if accountCreateParams.MaxAutomaticTokenAssociations != nil {
		transaction.SetMaxAutomaticTokenAssociations(*accountCreateParams.MaxAutomaticTokenAssociations)
	}
	if accountCreateParams.StakedAccountId != nil {
		accountId, err := hedera.AccountIDFromString(*accountCreateParams.StakedAccountId)
		if err != nil {
			return nil, err
		}
		transaction.SetStakedAccountID(accountId)
	}
	if accountCreateParams.StakedNodeId != nil {
		stakedNodeID, err := strconv.ParseInt(accountCreateParams.StakedNodeId.String(), 10, 64)
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
		accountCreateParams.CommonTransactionParams.FillOutTransaction(transaction, &transaction.Transaction, a.sdkService.Client)
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
	if receipt.Status == hedera.StatusSuccess {
		accId = receipt.AccountID.String()
	}

	return &response.AccountResponse{AccountId: accId, Status: receipt.Status.String()}, nil
}
