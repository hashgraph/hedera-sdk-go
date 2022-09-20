//go:build all || unit
// +build all unit

package hedera

/*-
 *
 * Hedera Go SDK
 *
 * Copyright (C) 2020 - 2022 Hedera Hashgraph, LLC
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

import (
	"encoding/base64"
	"testing"
	"time"

	"github.com/hashgraph/hedera-protobufs-go/services"

	"github.com/stretchr/testify/assert"

	"github.com/stretchr/testify/require"
)

func TestUnitTokenInfoQueryValidate(t *testing.T) {
	client := ClientForTestnet()
	client.SetAutoValidateChecksums(true)
	tokenID, err := TokenIDFromString("0.0.123-esxsf")
	require.NoError(t, err)

	tokenInfo := NewTokenInfoQuery().
		SetTokenID(tokenID)

	err = tokenInfo._ValidateNetworkOnIDs(client)
	require.NoError(t, err)
}

func TestUnitTokenInfoQueryValidateWrong(t *testing.T) {
	client := ClientForTestnet()
	client.SetAutoValidateChecksums(true)
	tokenID, err := TokenIDFromString("0.0.123-rmkykd")
	require.NoError(t, err)

	tokenInfo := NewTokenInfoQuery().
		SetTokenID(tokenID)

	err = tokenInfo._ValidateNetworkOnIDs(client)
	assert.Error(t, err)
	if err != nil {
		assert.Equal(t, "network mismatch or wrong checksum given, given checksum: rmkykd, correct checksum esxsf, network: testnet", err.Error())
	}
}

func TestUnitTokenInfoFromBytesBadBytes(t *testing.T) {
	bytes, err := base64.StdEncoding.DecodeString("tfhyY++/Q4BycortAgD4cmMKACB/")
	require.NoError(t, err)

	_, err = TokenInfoFromBytes(bytes)
	require.NoError(t, err)
}

func TestUnitTokenInfoFromBytesEmptyBytes(t *testing.T) {
	_, err := TokenInfoFromBytes([]byte{})
	require.NoError(t, err)
}

func TestUnitTokenInfoQueryGet(t *testing.T) {
	tokenID := TokenID{Token: 7}

	balance := NewTokenInfoQuery().
		SetTokenID(tokenID).
		SetQueryPayment(NewHbar(2)).
		SetMaxQueryPayment(NewHbar(1)).
		SetQueryPayment(HbarFromTinybar(25)).
		SetNodeAccountIDs([]AccountID{{Account: 10}, {Account: 11}, {Account: 12}})

	balance.GetTokenID()
	balance.GetNodeAccountIDs()
	balance.GetMinBackoff()
	balance.GetMaxBackoff()
	balance.GetMaxRetryCount()
	balance.GetPaymentTransactionID()
	balance.GetQueryPayment()
	balance.GetMaxQueryPayment()
}

func TestUnitTokenInfoQueryNothingSet(t *testing.T) {
	balance := NewTokenInfoQuery()

	balance.GetTokenID()
	balance.GetNodeAccountIDs()
	balance.GetMinBackoff()
	balance.GetMaxBackoff()
	balance.GetMaxRetryCount()
	balance.GetPaymentTransactionID()
	balance.GetQueryPayment()
	balance.GetMaxQueryPayment()
}

func TestUnitTokenInfoQueryCoverage(t *testing.T) {
	checksum := "dmqui"
	grpc := time.Second * 3
	token := TokenID{Token: 3, checksum: &checksum}
	nodeAccountID := []AccountID{{Account: 10}}
	transactionID := TransactionIDGenerate(AccountID{Account: 324})

	client := ClientForTestnet()
	client.SetAutoValidateChecksums(true)

	query := NewTokenInfoQuery().
		SetTokenID(token).
		SetMaxRetry(3).
		SetMaxBackoff(time.Second * 30).
		SetMinBackoff(time.Second * 10).
		SetNodeAccountIDs(nodeAccountID).
		SetPaymentTransactionID(transactionID).
		SetMaxQueryPayment(NewHbar(23)).
		SetQueryPayment(NewHbar(3)).
		SetGrpcDeadline(&grpc)

	err := query._ValidateNetworkOnIDs(client)
	require.NoError(t, err)
	query.GetNodeAccountIDs()
	query.GetMaxBackoff()
	query.GetMinBackoff()
	query._GetLogID()
	query.GetTokenID()
	query.GetQueryPayment()
	query.GetMaxQueryPayment()
}

func TestUnitTokenInfoQueryMock(t *testing.T) {
	responses := [][]interface{}{{
		&services.Response{
			Response: &services.Response_TokenGetInfo{
				TokenGetInfo: &services.TokenGetInfoResponse{
					Header: &services.ResponseHeader{NodeTransactionPrecheckCode: services.ResponseCodeEnum_OK, ResponseType: services.ResponseType_COST_ANSWER, Cost: 2},
				},
			},
		},
		&services.Response{
			Response: &services.Response_TokenGetInfo{
				TokenGetInfo: &services.TokenGetInfoResponse{
					Header: &services.ResponseHeader{NodeTransactionPrecheckCode: services.ResponseCodeEnum_OK, ResponseType: services.ResponseType_COST_ANSWER, Cost: 2},
				},
			},
		},
		&services.Response{
			Response: &services.Response_TokenGetInfo{
				TokenGetInfo: &services.TokenGetInfoResponse{
					Header: &services.ResponseHeader{NodeTransactionPrecheckCode: services.ResponseCodeEnum_OK, ResponseType: services.ResponseType_ANSWER_ONLY, Cost: 2},
					TokenInfo: &services.TokenInfo{
						TokenId:             nil,
						Name:                "",
						Symbol:              "",
						Decimals:            0,
						TotalSupply:         0,
						Treasury:            nil,
						AdminKey:            nil,
						KycKey:              nil,
						FreezeKey:           nil,
						WipeKey:             nil,
						SupplyKey:           nil,
						DefaultFreezeStatus: 0,
						DefaultKycStatus:    0,
						Deleted:             false,
						AutoRenewAccount:    nil,
						AutoRenewPeriod:     nil,
						Expiry:              nil,
						Memo:                "",
						TokenType:           0,
						SupplyType:          0,
						MaxSupply:           0,
						FeeScheduleKey:      nil,
						CustomFees:          nil,
						PauseKey:            nil,
						PauseStatus:         0,
						LedgerId:            nil,
					},
				},
			},
		},
	}}

	client, server := NewMockClientAndServer(responses)
	defer server.Close()

	query := NewTokenInfoQuery().
		SetTokenID(TokenID{Token: 3}).
		SetNodeAccountIDs([]AccountID{{Account: 3}}).
		SetMaxQueryPayment(NewHbar(1))

	_, err := query.GetCost(client)
	require.NoError(t, err)

	_, err = query.Execute(client)
	require.NoError(t, err)
}
