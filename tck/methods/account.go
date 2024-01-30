package methods

import (
	"context"
	"github.com/hashgraph/hedera-sdk-go/v2"
	"time"
)

type AccountService struct {
	sdkService *SDKService
	hedera.AccountID
}

type AccountParams struct {
	publicKey                     string
	initialBalance                int64
	receiverSignatureRequired     bool
	maxAutomaticTokenAssociations uint32
	stakedAccountId               string
	stakedNodeId                  int64
	declineStakingReward          bool
	accountMemo                   string
	autoRenewPeriod               time.Duration
	privateKey                    string
}

type AccountResponse struct {
	AccountId string
	Status    string
}

func (a *AccountService) SetSdkService(service *SDKService) {
	a.sdkService = service
}

func (a *AccountService) CreateAccount(_ context.Context, accountParams AccountParams) AccountResponse {
	durationValue := time.Second * 30000
	transaction := hedera.NewAccountCreateTransaction().SetGrpcDeadline(&durationValue)

	if accountParams.publicKey != "" {
		key, err := hedera.PublicKeyFromString(accountParams.publicKey)
		if err != nil {
			panic(err)
		}
		transaction.SetKey(key)
	}
	if accountParams.initialBalance != 0 {
		transaction.SetInitialBalance(hedera.HbarFromTinybar(accountParams.initialBalance))
	}
	if accountParams.receiverSignatureRequired {
		transaction.SetReceiverSignatureRequired(accountParams.receiverSignatureRequired)
	}
	if accountParams.maxAutomaticTokenAssociations != 0 {
		transaction.SetMaxAutomaticTokenAssociations(accountParams.maxAutomaticTokenAssociations)
	}
	if accountParams.stakedAccountId != "" {
		accountId, err := hedera.AccountIDFromString(accountParams.stakedAccountId)
		if err != nil {
			return AccountResponse{}
		}
		transaction.SetStakedAccountID(accountId)
	}
	if accountParams.stakedNodeId != 0 {
		transaction.SetStakedNodeID(accountParams.stakedNodeId)
	}
	if accountParams.declineStakingReward {
		transaction.SetDeclineStakingReward(accountParams.declineStakingReward)
	}
	if accountParams.accountMemo != "" {
		transaction.SetAccountMemo(accountParams.accountMemo)
	}
	if accountParams.autoRenewPeriod != 0 {
		transaction.SetAutoRenewPeriod(accountParams.autoRenewPeriod)
	}
	if accountParams.privateKey != "" {
		key, err := hedera.PrivateKeyFromString(accountParams.privateKey)
		if err != nil {
			panic(err)
		}

		_, err = transaction.FreezeWith(a.sdkService.Client)
		if err != nil {
			return AccountResponse{}
		}
		transaction.Sign(key)
	}

	txResponse, err := transaction.Execute(a.sdkService.Client)
	if err != nil {
		return AccountResponse{}
	}
	receipt, err := txResponse.GetReceipt(a.sdkService.Client)

	return AccountResponse{AccountId: receipt.AccountID.String(), Status: receipt.Status.String()}
}
