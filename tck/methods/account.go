package methods

import (
	"context"
	"encoding/json"
	"strconv"
	"time"

	"github.com/hashgraph/hedera-sdk-go/tck/param"
	"github.com/hashgraph/hedera-sdk-go/tck/response"
	"github.com/hashgraph/hedera-sdk-go/v2"
)

// ---- Struct to hold hedera.Client implementation and to implement the methods of the specification ----
type AccountService struct {
	sdkService *SDKService
	hedera.AccountID
}

// Variable to be set to `SetGrpcDeadline` for all transactions
var threeSecondsDuration = time.Second * 3

// SetSdkService We set object, which is holding our client param. Pass it by referance, because TCK is dynamically updating it
func (a *AccountService) SetSdkService(service *SDKService) {
	a.sdkService = service
}

// CreateAccount jRPC method for createAccount
func (a *AccountService) CreateAccount(_ context.Context, accountParams param.CreateAccountParams) (*response.AccountResponse, error) {
	transaction := hedera.NewAccountCreateTransaction().SetGrpcDeadline(&threeSecondsDuration)

	if accountParams.PublicKey != "" {
		key, err := hedera.PublicKeyFromString(accountParams.PublicKey)
		if err != nil {
			return nil, err
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
	if accountParams.StakedAccountId != nil {
		accountId, err := hedera.AccountIDFromString(*accountParams.StakedAccountId)
		if err != nil {
			return nil, err
		}
		transaction.SetStakedAccountID(accountId)
	}
	if accountParams.StakedNodeId.String() != "" {
		stakedNodeID, err := strconv.ParseInt(accountParams.StakedNodeId.String(), 10, 64)
		if err != nil {
			return nil, response.InvalidParams.WithData(err.Error())
		}
		transaction.SetStakedNodeID(stakedNodeID)
	}
	if accountParams.DeclineStakingReward {
		transaction.SetDeclineStakingReward(accountParams.DeclineStakingReward)
	}
	if accountParams.AccountMemo != "" {
		transaction.SetAccountMemo(accountParams.AccountMemo)
	}
	if accountParams.AutoRenewPeriod != "" {
		autoRenewPeriod, err := strconv.ParseInt(accountParams.AutoRenewPeriod, 10, 64)
		if err != nil {
			return nil, response.InvalidParams.WithData(err.Error())
		}
		transaction.SetAutoRenewPeriod(time.Duration(autoRenewPeriod) * time.Second)
	}
	if accountParams.PrivateKey != "" {
		key, err := hedera.PrivateKeyFromString(accountParams.PrivateKey)
		if err != nil {
			return nil, err
		}
		_, err = transaction.FreezeWith(a.sdkService.Client)
		if err != nil {
			return nil, err
		}
		transaction.Sign(key)
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

// CreateAccountFromAlias Create an account from aliasId by transferring some HBAR amount from given account (opperatorId)
func (a *AccountService) CreateAccountFromAlias(_ context.Context, fromAliasParams param.AccountFromAliasParams) (*hedera.TransactionReceipt, error) {
	operator, err := hedera.AccountIDFromString(fromAliasParams.OperatorId)
	if err != nil {
		return nil, err
	}
	var aliasId param.Alias
	if err := json.Unmarshal([]byte(fromAliasParams.AliasAccountId), &aliasId); err != nil {
		return nil, response.InvalidParams.WithData(err.Error())
	}

	alias, err := hedera.AccountIDFromString(string(aliasId))
	if err != nil {
		return nil, err
	}

	resp, err := hedera.NewTransferTransaction().SetGrpcDeadline(&threeSecondsDuration).
		AddHbarTransfer(operator, hedera.NewHbar(float64(fromAliasParams.InitialBalance)).Negated()).
		AddHbarTransfer(alias, hedera.NewHbar(float64(fromAliasParams.InitialBalance))).
		Execute(a.sdkService.Client)
	if err != nil {
		return nil, err
	}

	receipt, err := resp.GetReceipt(a.sdkService.Client)
	if err != nil {
		return nil, err
	}
	return &receipt, nil
}

// GetAccountInfo Get info for a given accountId and return a custom object containing aggregated information about the account
func (a *AccountService) GetAccountInfo(_ context.Context, accountId string) (*response.AccountInfoResponse, error) {
	account, err := hedera.AccountIDFromString(accountId)
	if err != nil {
		return nil, err
	}
	resp, err := hedera.NewAccountInfoQuery().SetAccountID(account).SetGrpcDeadline(&threeSecondsDuration).Execute(a.sdkService.Client)
	if err != nil {
		return nil, err
	}

	return &response.AccountInfoResponse{AccountID: resp.AccountID.String(),
		Balance:                       resp.Balance.String(),
		Key:                           resp.Key.String(),
		AccountMemo:                   resp.AccountMemo,
		MaxAutomaticTokenAssociations: resp.MaxAutomaticTokenAssociations,
		AutoRenewPeriod:               resp.AutoRenewPeriod}, nil
}

// DeleteAccount deletes a provided account by signing the transaction with the key of that account
func (a *AccountService) DeleteAccount(_ context.Context, param param.DeleteAccountParams) (*response.AccountResponse, error) {
	accId, _ := hedera.AccountIDFromString(param.AccountId)
	recipientId, _ := hedera.AccountIDFromString(param.RecipientId)
	tx, err := hedera.NewAccountDeleteTransaction().SetGrpcDeadline(&threeSecondsDuration).SetAccountID(accId).SetTransferAccountID(recipientId).FreezeWith(a.sdkService.Client)
	if err != nil {
		return nil, err
	}
	signature, _ := hedera.PrivateKeyFromString(param.AccountKey)

	resp, err := tx.Sign(signature).Execute(a.sdkService.Client)
	if err != nil {
		return nil, err
	}

	receipt, _ := resp.GetReceipt(a.sdkService.Client)
	return &response.AccountResponse{Status: receipt.Status.String()}, nil
}

// UpdateAccountKey updates an existing acoount id keys with provided params
func (a *AccountService) UpdateAccountKey(_ context.Context, params param.UpdateAccountParams) (*response.AccountResponse, error) {
	accId, err := hedera.AccountIDFromString(params.AccountId)
	if err != nil {
		return nil, err
	}
	key, err := hedera.PublicKeyFromString(params.NewPublicKey)
	if err != nil {
		return nil, err
	}

	tx, err := hedera.NewAccountUpdateTransaction().SetGrpcDeadline(&threeSecondsDuration).SetAccountID(accId).SetKey(key).FreezeWith(a.sdkService.Client)
	if err != nil {
		return nil, err
	}

	oldKey, err := hedera.PrivateKeyFromString(params.OldPrivateKey)
	if err != nil {
		return nil, err
	}
	newKey, err := hedera.PrivateKeyFromString(params.NewPrivateKey)
	if err != nil {
		return nil, err
	}
	tx.Sign(oldKey)
	tx.Sign(newKey)

	resp, err := tx.Execute(a.sdkService.Client)
	if err != nil {
		return nil, err
	}

	receipt, err := resp.GetReceipt(a.sdkService.Client)
	if err != nil {
		return nil, err
	}

	return &response.AccountResponse{Status: receipt.Status.String()}, nil
}

// UpdateAccountMemo updates account memo of an existing account ID
func (a *AccountService) UpdateAccountMemo(_ context.Context, params param.UpdateAccountParams) (*response.AccountResponse, error) {
	accId, err := hedera.AccountIDFromString(params.AccountId)
	if err != nil {
		return nil, err
	}
	tx, err := hedera.NewAccountUpdateTransaction().SetGrpcDeadline(&threeSecondsDuration).SetAccountID(accId).SetAccountMemo(params.Memo).FreezeWith(a.sdkService.Client)
	if err != nil {
		return nil, err
	}
	signature, _ := hedera.PrivateKeyFromString(params.Key)

	resp, err := tx.Sign(signature).Execute(a.sdkService.Client)
	if err != nil {
		return nil, err
	}

	receipt, err := resp.GetReceipt(a.sdkService.Client)
	if err != nil {
		return nil, err
	}
	return &response.AccountResponse{Status: receipt.Status.String()}, nil
}
