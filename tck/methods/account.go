package methods

// SPDX-License-Identifier: Apache-2.0

import (
	"bytes"
	"context"
	"encoding/hex"
	"encoding/json"
	"errors"
	"strconv"
	"time"

	"github.com/hiero-ledger/hiero-sdk-go/tck/param"
	"github.com/hiero-ledger/hiero-sdk-go/tck/response"
	"github.com/hiero-ledger/hiero-sdk-go/tck/utils"
	hiero "github.com/hiero-ledger/hiero-sdk-go/v2/sdk"
)

// ---- Struct to hold hiero.Client implementation and to implement the methods of the specification ----
type AccountService struct {
	sdkService *SDKService
}

// Variable to be set to `SetGrpcDeadline` for all transactions
var threeSecondsDuration = time.Second * 3

// SetSdkService We set object, which is holding our client param. Pass it by reference, because TCK is dynamically updating it
func (a *AccountService) SetSdkService(service *SDKService) {
	a.sdkService = service
}

// CreateAccount jRPC method for createAccount
func (a *AccountService) CreateAccount(_ context.Context, params param.CreateAccountParams) (*response.AccountResponse, error) {
	transaction := hiero.NewAccountCreateTransaction().SetGrpcDeadline(&threeSecondsDuration)

	// Set key
	if err := utils.SetKeyIfPresent(params.Key, transaction.SetKeyWithoutAlias); err != nil {
		return nil, err
	}
	if params.InitialBalance != nil {
		initialBalance, err := strconv.ParseInt(*params.InitialBalance, 10, 64)
		if err != nil {
			return nil, err
		}
		transaction.SetInitialBalance(hiero.HbarFromTinybar(initialBalance))
	}
	if params.ReceiverSignatureRequired != nil {
		transaction.SetReceiverSignatureRequired(*params.ReceiverSignatureRequired)
	}
	if params.MaxAutomaticTokenAssociations != nil {
		transaction.SetMaxAutomaticTokenAssociations(*params.MaxAutomaticTokenAssociations)
	}
	// Set staked account ID
	if err := utils.SetAccountIDIfPresent(params.StakedAccountId, transaction.SetStakedAccountID); err != nil {
		return nil, err
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
		autoRenewPeriodSeconds, err := strconv.ParseInt(*params.AutoRenewPeriod, 10, 64)
		if err != nil {
			return nil, err
		}

		transaction.SetAutoRenewPeriod(time.Duration(autoRenewPeriodSeconds) * time.Second)
	}
	if params.Alias != nil {
		transaction.SetAlias(*params.Alias)
	}
	if params.CommonTransactionParams != nil {
		err := params.CommonTransactionParams.FillOutTransaction(transaction, a.sdkService.GetClient(params.SessionId))
		if err != nil {
			return nil, err
		}
	}
	txResponse, err := transaction.Execute(a.sdkService.GetClient(params.SessionId))
	if err != nil {
		return nil, err
	}
	receipt, err := txResponse.GetReceipt(a.sdkService.GetClient(params.SessionId))
	if err != nil {
		return nil, err
	}
	var accId string
	if receipt.Status == hiero.StatusSuccess {
		accId = receipt.AccountID.String()
	}
	return &response.AccountResponse{AccountId: accId, Status: receipt.Status.String(), TransactionId: txResponse.TransactionID.String()}, nil
}

// buildCreateAccount builds an account create transaction without executing it (for scheduling)
func (a *AccountService) buildCreateAccount(params param.CreateAccountParams) (*hiero.AccountCreateTransaction, error) {
	transaction := hiero.NewAccountCreateTransaction().SetGrpcDeadline(&threeSecondsDuration)

	// Set key
	if err := utils.SetKeyIfPresent(params.Key, transaction.SetKeyWithoutAlias); err != nil {
		return nil, err
	}
	if params.InitialBalance != nil {
		initialBalance, err := strconv.ParseInt(*params.InitialBalance, 10, 64)
		if err != nil {
			return nil, err
		}
		transaction.SetInitialBalance(hiero.HbarFromTinybar(initialBalance))
	}
	if params.ReceiverSignatureRequired != nil {
		transaction.SetReceiverSignatureRequired(*params.ReceiverSignatureRequired)
	}
	if params.MaxAutomaticTokenAssociations != nil {
		transaction.SetMaxAutomaticTokenAssociations(*params.MaxAutomaticTokenAssociations)
	}
	// Set staked account ID
	if err := utils.SetAccountIDIfPresent(params.StakedAccountId, transaction.SetStakedAccountID); err != nil {
		return nil, err
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
		autoRenewPeriodSeconds, err := strconv.ParseInt(*params.AutoRenewPeriod, 10, 64)
		if err != nil {
			return nil, err
		}

		transaction.SetAutoRenewPeriod(time.Duration(autoRenewPeriodSeconds) * time.Second)
	}
	if params.Alias != nil {
		transaction.SetAlias(*params.Alias)
	}
	if params.CommonTransactionParams != nil {
		err := params.CommonTransactionParams.FillOutTransaction(transaction, a.sdkService.GetClient(params.SessionId))
		if err != nil {
			return nil, err
		}
	}

	return transaction, nil
}

// UpdateAccount jRPC method for updateAccount
func (a *AccountService) UpdateAccount(_ context.Context, params param.UpdateAccountParams) (*response.AccountResponse, error) {
	transaction := hiero.NewAccountUpdateTransaction().SetGrpcDeadline(&threeSecondsDuration)
	// Set account ID
	if err := utils.SetAccountIDIfPresent(params.AccountId, transaction.SetAccountID); err != nil {
		return nil, err
	}

	// Set key
	if err := utils.SetKeyIfPresent(params.Key, transaction.SetKey); err != nil {
		return nil, err
	}

	if params.ExpirationTime != nil {
		expirationTime, err := strconv.ParseInt(*params.ExpirationTime, 10, 64)
		if err != nil {
			return nil, err
		}
		transaction.SetExpirationTime(time.Unix(expirationTime, 0))
	}

	if params.ReceiverSignatureRequired != nil {
		transaction.SetReceiverSignatureRequired(*params.ReceiverSignatureRequired)
	}

	if params.MaxAutomaticTokenAssociations != nil {
		transaction.SetMaxAutomaticTokenAssociations(*params.MaxAutomaticTokenAssociations)
	}

	// Set staked account ID
	if err := utils.SetAccountIDIfPresent(params.StakedAccountId, transaction.SetStakedAccountID); err != nil {
		return nil, err
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
		autoRenewPeriodSeconds, err := strconv.ParseInt(*params.AutoRenewPeriod, 10, 64)
		if err != nil {
			return nil, err
		}
		transaction.SetAutoRenewPeriod(time.Duration(autoRenewPeriodSeconds) * time.Second)
	}

	if params.CommonTransactionParams != nil {
		err := params.CommonTransactionParams.FillOutTransaction(transaction, a.sdkService.GetClient(params.SessionId))
		if err != nil {
			return nil, err
		}
	}

	txResponse, err := transaction.Execute(a.sdkService.GetClient(params.SessionId))
	if err != nil {
		return nil, err
	}
	receipt, err := txResponse.GetReceipt(a.sdkService.GetClient(params.SessionId))
	if err != nil {
		return nil, err
	}
	return &response.AccountResponse{Status: receipt.Status.String()}, nil
}

// DeleteAccount jRPC method for deleteAccount
func (a *AccountService) DeleteAccount(_ context.Context, params param.DeleteAccountParams) (*response.AccountResponse, error) {
	transaction := hiero.NewAccountDeleteTransaction().SetGrpcDeadline(&threeSecondsDuration)
	// Set account ID
	if err := utils.SetAccountIDIfPresent(params.DeleteAccountId, transaction.SetAccountID); err != nil {
		return nil, err
	}

	// Set transfer account ID
	if err := utils.SetAccountIDIfPresent(params.TransferAccountId, transaction.SetTransferAccountID); err != nil {
		return nil, err
	}

	if params.CommonTransactionParams != nil {
		err := params.CommonTransactionParams.FillOutTransaction(transaction, a.sdkService.GetClient(params.SessionId))
		if err != nil {
			return nil, err
		}
	}

	txResponse, err := transaction.Execute(a.sdkService.GetClient(params.SessionId))
	if err != nil {
		return nil, err
	}
	receipt, err := txResponse.GetReceipt(a.sdkService.GetClient(params.SessionId))
	if err != nil {
		return nil, err
	}
	return &response.AccountResponse{Status: receipt.Status.String()}, nil
}

// buildApproveAllowance builds an AccountAllowanceApproveTransaction from parameters
func (a *AccountService) buildApproveAllowance(params param.AccountAllowanceApproveParams) (*hiero.AccountAllowanceApproveTransaction, error) {
	transaction := hiero.NewAccountAllowanceApproveTransaction().SetGrpcDeadline(&threeSecondsDuration)

	allowances := *params.Allowances

	for _, allowance := range allowances {
		owner, err := hiero.AccountIDFromString(*allowance.OwnerAccountId)
		if err != nil {
			return nil, err
		}

		spender, err := hiero.AccountIDFromString(*allowance.SpenderAccountId)
		if err != nil {
			return nil, err
		}

		hbar := allowance.Hbar
		token := allowance.Token
		nft := allowance.Nft

		switch {
		case hbar != nil:
			// Process Hbar allowance
			hbarAmount, err := strconv.ParseInt(*hbar.Amount, 10, 64)
			if err != nil {
				return nil, err
			}
			transaction.ApproveHbarAllowance(owner, spender, hiero.HbarFromTinybar(hbarAmount))

		case token != nil:
			// Process Token allowance
			tokenID, err := hiero.TokenIDFromString(*token.TokenId)
			if err != nil {
				return nil, err
			}
			tokenAmount, err := strconv.ParseInt(*token.Amount, 10, 64)
			if err != nil {
				return nil, err
			}
			transaction.ApproveTokenAllowance(tokenID, owner, spender, tokenAmount)

		case nft != nil:
			// Process Nft allowance
			tokenID, err := hiero.TokenIDFromString(*nft.TokenId)
			if err != nil {
				return nil, err
			}

			switch {
			case nft.SerialNumbers != nil:
				for _, serialNumber := range *nft.SerialNumbers {
					serialNumberParsed, err := strconv.ParseInt(serialNumber, 10, 64)
					if err != nil {
						return nil, err
					}

					nftID := hiero.NftID{
						TokenID:      tokenID,
						SerialNumber: serialNumberParsed,
					}

					if nft.DelegateSpenderAccountId != nil {
						delegateSpenderAccountId, err := hiero.AccountIDFromString(*nft.DelegateSpenderAccountId)
						if err != nil {
							return nil, err
						}

						transaction.ApproveTokenNftAllowanceWithDelegatingSpender(
							nftID,
							owner,
							spender,
							delegateSpenderAccountId,
						)
					} else {
						transaction.ApproveTokenNftAllowance(
							nftID,
							owner,
							spender,
						)
					}
				}
			case nft.ApprovedForAll != nil && *nft.ApprovedForAll:
				transaction.ApproveTokenNftAllowanceAllSerials(
					tokenID,
					owner,
					spender,
				)
			default:
				transaction.DeleteTokenNftAllowanceAllSerials(tokenID, owner, spender)
			}

		default:
			return nil, errors.New("no valid allowance type provided")
		}
	}

	if params.CommonTransactionParams != nil {
		err := params.CommonTransactionParams.FillOutTransaction(transaction, a.sdkService.GetClient(params.SessionId))
		if err != nil {
			return nil, err
		}
	}
	return transaction, nil
}

// ApproveAllowance jRPC method for approveAllowance
func (a *AccountService) ApproveAllowance(_ context.Context, params param.AccountAllowanceApproveParams) (*response.AccountResponse, error) {
	transaction, err := a.buildApproveAllowance(params)
	if err != nil {
		return nil, err
	}

	txResponse, err := transaction.Execute(a.sdkService.GetClient(params.SessionId))
	if err != nil {
		return nil, err
	}
	receipt, err := txResponse.GetReceipt(a.sdkService.GetClient(params.SessionId))
	if err != nil {
		return nil, err
	}
	return &response.AccountResponse{Status: receipt.Status.String()}, nil
}

// DeleteAllowance jRPC method for deleteAllowance
func (a *AccountService) DeleteAllowance(_ context.Context, params param.AccountAllowanceDeleteParams) (*response.AccountResponse, error) {
	transaction := hiero.NewAccountAllowanceDeleteTransaction().SetGrpcDeadline(&threeSecondsDuration)

	allowances := *params.Allowances

	// Loop through each allowance and process
	for _, allowance := range allowances {
		owner, err := hiero.AccountIDFromString(*allowance.OwnerAccountId)
		if err != nil {
			return nil, err
		}

		tokenID, err := hiero.TokenIDFromString(*allowance.TokenId)
		if err != nil {
			return nil, err
		}

		// Process NFT serial numbers if provided
		if allowance.SerialNumbers != nil {
			for _, serialNumber := range *allowance.SerialNumbers {
				serialNumberParsed, err := strconv.ParseInt(serialNumber, 10, 64)
				if err != nil {
					return nil, err
				}

				nftID := hiero.NftID{
					TokenID:      tokenID,
					SerialNumber: serialNumberParsed,
				}

				transaction.DeleteAllTokenNftAllowances(nftID, &owner)
			}
		} else {
			transaction.DeleteAllTokenNftAllowances(hiero.NftID{TokenID: tokenID}, &owner)
		}
	}

	if params.CommonTransactionParams != nil {
		err := params.CommonTransactionParams.FillOutTransaction(transaction, a.sdkService.GetClient(params.SessionId))
		if err != nil {
			return nil, err
		}
	}

	txResponse, err := transaction.Execute(a.sdkService.GetClient(params.SessionId))
	if err != nil {
		return nil, err
	}
	receipt, err := txResponse.GetReceipt(a.sdkService.GetClient(params.SessionId))
	if err != nil {
		return nil, err
	}
	return &response.AccountResponse{Status: receipt.Status.String()}, nil
}

// buildTransferCrypto builds a TransferTransaction from parameters
func (a *AccountService) buildTransferCrypto(params param.TransferCryptoParams) (*hiero.TransferTransaction, error) {
	transaction := hiero.NewTransferTransaction().SetGrpcDeadline(&threeSecondsDuration)

	if params.Transfers == nil {
		return nil, response.NewInternalError("transferParams is required")
	}

	transferParams := *params.Transfers
	if len(transferParams) == 0 {
		return nil, response.NewInternalError("transferParams is required")
	}

	for _, transferParam := range transferParams {
		if err := utils.HandleTransferParam(transaction, transferParam); err != nil {
			return nil, err
		}
	}

	if params.CommonTransactionParams != nil {
		err := params.CommonTransactionParams.FillOutTransaction(transaction, a.sdkService.GetClient(params.SessionId))
		if err != nil {
			return nil, err
		}
	}
	return transaction, nil
}

// TransferCrypto jRPC method for transferCrypto
func (a *AccountService) TransferCrypto(_ context.Context, params param.TransferCryptoParams) (*response.AccountResponse, error) {
	transaction, err := a.buildTransferCrypto(params)
	if err != nil {
		return nil, err
	}

	txResponse, err := transaction.Execute(a.sdkService.GetClient(params.SessionId))
	if err != nil {
		return nil, err
	}

	receipt, err := txResponse.GetReceipt(a.sdkService.GetClient(params.SessionId))
	if err != nil {
		return nil, err
	}

	return &response.AccountResponse{Status: receipt.Status.String()}, nil
}

// GetAccountInfo jRPC method for getAccountInfo
func (a *AccountService) GetAccountInfo(_ context.Context, params param.GetAccountInfoParams) (*response.AccountInfoResponse, error) {
	query := hiero.NewAccountInfoQuery().SetGrpcDeadline(&threeSecondsDuration)

	if err := utils.SetAccountIDIfPresent(params.AccountId, query.SetAccountID); err != nil {
		return nil, err
	}

	info, err := query.Execute(a.sdkService.GetClient(params.SessionId))
	if err != nil {
		return nil, err
	}

	result := &response.AccountInfoResponse{
		AccountId:                     info.AccountID.String(),
		ContractAccountId:             info.ContractAccountID,
		IsDeleted:                     info.IsDeleted,
		ProxyReceived:                 strconv.FormatInt(info.ProxyReceived.AsTinybar(), 10),
		Balance:                       strconv.FormatInt(info.Balance.AsTinybar(), 10),
		SendRecordThreshold:           strconv.FormatInt(info.GenerateSendRecordThreshold.AsTinybar(), 10),
		ReceiveRecordThreshold:        strconv.FormatInt(info.GenerateReceiveRecordThreshold.AsTinybar(), 10),
		IsReceiverSignatureRequired:   info.ReceiverSigRequired,
		ExpirationTime:                strconv.FormatInt(info.ExpirationTime.Unix(), 10),
		AutoRenewPeriod:               strconv.FormatInt(int64(info.AutoRenewPeriod.Seconds()), 10),
		AccountMemo:                   info.AccountMemo,
		OwnedNfts:                     strconv.FormatInt(info.OwnedNfts, 10),
		MaxAutomaticTokenAssociations: strconv.FormatUint(uint64(info.MaxAutomaticTokenAssociations), 10),
		LedgerId:                      info.LedgerID.String(),
		EthereumNonce:                 strconv.FormatInt(info.EthereumNonce, 10),
		LiveHashes:                    liveHashesToResponse(info.LiveHashes),
		HbarAllowances:                []response.HbarAllowanceResponse{},
		TokenAllowances:               []response.TokenAllowanceResponse{},
		NftAllowances:                 []response.TokenNftAllowanceResponse{},
		TokenRelationships:            tokenRelationshipsToResponse(info.TokenRelationships),
		StakingInfo:                   stakingInfoToResponse(info.StakingInfo),
	}

	if info.Key != nil {
		result.Key = info.Key.String()
	}

	if info.AliasKey != nil {
		aliasKey := info.AliasKey.String()
		result.AliasKey = &aliasKey
	}

	return result, nil
}

// Map sdk LiveHash to jRPC LiveHashResponse
func liveHashesToResponse(liveHashes []*hiero.LiveHash) []response.LiveHashResponse {
	if len(liveHashes) == 0 {
		return []response.LiveHashResponse{}
	}

	result := make([]response.LiveHashResponse, 0, len(liveHashes))

	for _, liveHash := range liveHashes {
		keys := []string{}
		for _, key := range liveHash.Keys.GetKeys() {
			keys = append(keys, key.String())
		}

		result = append(result, response.LiveHashResponse{
			AccountId: liveHash.AccountID.String(),
			Hash:      hex.EncodeToString(liveHash.Hash),
			Keys:      keys,
			Duration:  strconv.FormatInt(int64(liveHash.LiveHashDuration.Seconds()), 10),
		})
	}

	return result
}

// Map sdk TokenRelationship to jRPC TokenRelationshipResponse
func tokenRelationshipsToResponse(relationships []*hiero.TokenRelationship) map[string]response.TokenRelationshipInfo {
	if len(relationships) == 0 {
		return map[string]response.TokenRelationshipInfo{}
	}

	result := make(map[string]response.TokenRelationshipInfo, len(relationships))

	for _, relationship := range relationships {
		if relationship == nil {
			continue
		}

		tokenID := relationship.TokenID.String()

		result[tokenID] = response.TokenRelationshipInfo{
			TokenId:              tokenID,
			Symbol:               relationship.Symbol,
			Balance:              strconv.FormatUint(relationship.Balance, 10),
			IsKycGranted:         relationship.KycStatus,
			IsFrozen:             relationship.FreezeStatus,
			AutomaticAssociation: relationship.AutomaticAssociation,
		}
	}

	return result
}

// Map sdk StakingInfo to jRPC StakingInfoResponse
func stakingInfoToResponse(info *hiero.StakingInfo) *response.StakingInfoResponse {
	if info == nil {
		return nil
	}

	stakeInfo := &response.StakingInfoResponse{
		DeclineStakingReward: info.DeclineStakingReward,
		PendingReward:        strconv.FormatInt(info.PendingHbarReward.AsTinybar(), 10),
		StakedToMe:           strconv.FormatInt(info.StakedToMe.AsTinybar(), 10),
	}

	if info.StakePeriodStart != nil && !info.StakePeriodStart.IsZero() {
		stakePeriod := strconv.FormatInt(info.StakePeriodStart.Unix(), 10)
		stakeInfo.StakePeriodStart = &stakePeriod
	}

	if info.StakedAccountID != nil {
		stakedAccountId := info.StakedAccountID.String()
		stakeInfo.StakedAccountId = &stakedAccountId
	}

	if info.StakedNodeID != nil {
		stakedNodeId := strconv.FormatInt(*info.StakedNodeID, 10)
		stakeInfo.StakedNodeId = &stakedNodeId
	}

	return stakeInfo
}

// GetTransactionReceipt jRPC method for getTransactionReceipt
func (a *AccountService) GetTransactionReceipt(_ context.Context, params param.GetTransactionReceiptParams) (*response.TransactionReceiptResponse, error) {
	if params.TransactionID == nil {
		return nil, hiero.ErrHederaPreCheckStatus{Status: hiero.StatusInvalidTransactionID}
	}

	transactionId, err := hiero.TransactionIdFromString(*params.TransactionID)
	if err != nil {
		return nil, err
	}

	query := hiero.NewTransactionReceiptQuery().SetGrpcDeadline(&threeSecondsDuration).SetTransactionID(transactionId)

	if params.IncludeChildren != nil {
		query.SetIncludeChildren(*params.IncludeChildren)
	}

	if params.IncludeDuplicates != nil {
		query.SetIncludeDuplicates(*params.IncludeDuplicates)
	}

	receipt, err := query.Execute(a.sdkService.GetClient(params.SessionId))
	if err != nil {
		return nil, err
	}

	validateStatus := true
	if params.ValidateStatus != nil {
		validateStatus = *params.ValidateStatus
	}

	if err := receipt.ValidateStatus(validateStatus); err != nil {
		return nil, err
	}

	result := transactionReceiptToResponse(receipt)

	for _, duplicate := range receipt.Duplicates {
		result.Duplicates = append(result.Duplicates, transactionReceiptToResponse(duplicate))
	}

	for _, child := range receipt.Children {
		result.Children = append(result.Children, transactionReceiptToResponse(child))
	}

	return &result, nil

}

func transactionReceiptToResponse(receipt hiero.TransactionReceipt) response.TransactionReceiptResponse {
	result := response.TransactionReceiptResponse{
		Status:     receipt.Status.String(),
		Serials:    utils.MapSerialNumbersToString(receipt.SerialNumbers),
		Duplicates: []response.TransactionReceiptResponse{},
		Children:   []response.TransactionReceiptResponse{},
	}

	if receipt.AccountID != nil {
		accountId := receipt.AccountID.String()
		result.AccountId = &accountId
	}

	if receipt.ContractID != nil {
		contractId := receipt.ContractID.String()
		result.ContractId = &contractId
	}

	if receipt.FileID != nil {
		fileId := receipt.FileID.String()
		result.FileId = &fileId
	}

	if receipt.TopicID != nil {
		topicId := receipt.TopicID.String()
		result.TopicId = &topicId
	}

	if receipt.TokenID != nil {
		tokenId := receipt.TokenID.String()
		result.TokenId = &tokenId
	}

	if receipt.ScheduleID != nil {
		scheduleId := receipt.ScheduleID.String()
		result.ScheduleId = &scheduleId
	}

	if receipt.ScheduledTransactionID != nil {
		scheduledTransactionId := receipt.ScheduledTransactionID.String()
		result.ScheduledTransactionId = &scheduledTransactionId
	}

	if receipt.ExchangeRate != nil {
		result.ExchangeRate = &response.ExchangeRate{
			Hbars: receipt.ExchangeRate.Hbars,
			Cents: receipt.ExchangeRate.GetCents(),
		}

		if expirationTime := receipt.ExchangeRate.GetExpirationTime(); expirationTime != nil {
			expiration := expirationTime.UTC().Format(time.RFC3339)
			result.ExchangeRate.ExpirationTime = &expiration
		}
	}

	if receipt.TopicSequenceNumber != 0 {
		topicSequenceNumber := strconv.FormatUint(receipt.TopicSequenceNumber, 10)
		result.TopicSequenceNumber = &topicSequenceNumber
	}

	if len(receipt.TopicRunningHash) > 0 {
		topicRunningHash := hex.EncodeToString(receipt.TopicRunningHash)
		result.TopicRunningHash = &topicRunningHash
	}

	if receipt.TotalSupply != 0 {
		totalSupply := strconv.FormatUint(receipt.TotalSupply, 10)
		result.TotalSupply = &totalSupply
	}

	if receipt.NodeID != 0 {
		nodeId := strconv.FormatUint(receipt.NodeID, 10)
		result.NodeId = &nodeId
	}

	return result

}

// GetMirrorNodeAccountBalance jRPC method for getMirrorNodeAccountBalance
func (a *AccountService) GetMirrorNodeAccountBalance(_ context.Context, params param.GetMirrorNodeAccountBalanceParams) (*response.MirrorNodeAccountBalanceResponse, error) {
	query := hiero.NewMirrorNodeAccountBalanceQuery()

	// A missing account ID is left to the SDK, which rejects the query before any network call
	if err := utils.SetAccountIDIfPresent(params.AccountId, query.SetAccountID); err != nil {
		return nil, err
	}

	balance, err := query.Execute(a.sdkService.GetClient(params.SessionId))
	if err != nil {
		return nil, err
	}

	return &response.MirrorNodeAccountBalanceResponse{Hbar: strconv.FormatInt(balance.Hbars.AsTinybar(), 10)}, nil
}

// ExecuteDeprecatedAccountBalanceQuery jRPC method for executeDeprecatedAccountBalanceQuery
func (a *AccountService) ExecuteDeprecatedAccountBalanceQuery(_ context.Context, params param.ExecuteDeprecatedAccountBalanceQueryParams) (*response.DeprecatedAccountBalanceQueryResponse, error) {
	if params.AccountId == nil {
		return nil, response.NewInternalError("accountId is required")
	}
	operation := "execute"
	if params.Operation != nil {
		operation = *params.Operation
	}
	if operation != "execute" && operation != "getCost" {
		return nil, response.InvalidParams.WithData("unknown operation: " + operation)
	}
	accountID, err := hiero.AccountIDFromString(*params.AccountId)
	if err != nil {
		return nil, err
	}

	query, logged, err := captureStdout(hiero.NewAccountBalanceQuery) //nolint:staticcheck // this method exists to construct the deprecated query
	if err != nil {
		return nil, err
	}
	query.SetAccountID(accountID)

	client := a.sdkService.GetClient(params.SessionId)
	if operation == "getCost" {
		_, err = query.GetCost(client) //nolint:staticcheck // this method exists to execute the deprecated query
	} else {
		_, err = query.Execute(client) //nolint:staticcheck // this method exists to execute the deprecated query
	}

	var executionError *string
	if err != nil {
		message := err.Error()
		executionError = &message
	}

	return &response.DeprecatedAccountBalanceQueryResponse{
		ConstructionWarning: loggedWarning(logged),
		ExecutionError:      executionError,
	}, nil
}

// loggedWarning returns the message of the first AccountBalanceQuery warning in the captured output, or nil.
// Only the default JSON log format is parsed, not the HEDERA_SDK_GO_LOG_PRETTY console format.
func loggedWarning(output []byte) *string {
	for line := range bytes.Lines(output) {
		var entry struct {
			Level   string `json:"level"`
			Module  string `json:"module"`
			Message string `json:"message"`
		}
		if json.Unmarshal(line, &entry) == nil && entry.Level == "warn" && entry.Module == "AccountBalanceQuery" {
			return &entry.Message
		}
	}
	return nil
}
