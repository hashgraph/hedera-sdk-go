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
	"github.com/hashgraph/hedera-protobufs-go/services"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestUnitScheduleInfoQueryValidate(t *testing.T) {
	client := ClientForTestnet()
	client.SetAutoValidateChecksums(true)
	scheduleID, err := ScheduleIDFromString("0.0.123-esxsf")
	require.NoError(t, err)

	scheduleInfo := NewScheduleInfoQuery().
		SetScheduleID(scheduleID)

	err = scheduleInfo._ValidateNetworkOnIDs(client)
	require.NoError(t, err)
}

func TestUnitScheduleInfoQueryValidateWrong(t *testing.T) {
	client := ClientForTestnet()
	client.SetAutoValidateChecksums(true)
	scheduleID, err := ScheduleIDFromString("0.0.123-rmkykd")
	require.NoError(t, err)

	scheduleInfo := NewScheduleInfoQuery().
		SetScheduleID(scheduleID)

	err = scheduleInfo._ValidateNetworkOnIDs(client)
	assert.Error(t, err)
	if err != nil {
		assert.Equal(t, "network mismatch or wrong checksum given, given checksum: rmkykd, correct checksum esxsf, network: testnet", err.Error())
	}
}

func TestUnitScheduleInfoQueryGet(t *testing.T) {
	scheduleID := ScheduleID{Schedule: 7}

	balance := NewScheduleInfoQuery().
		SetScheduleID(scheduleID).
		SetQueryPayment(NewHbar(2)).
		SetMaxQueryPayment(NewHbar(1)).
		SetQueryPayment(HbarFromTinybar(25)).
		SetNodeAccountIDs([]AccountID{{Account: 10}, {Account: 11}, {Account: 12}})

	balance.GetScheduleID()
	balance.GetNodeAccountIDs()
	balance.GetMinBackoff()
	balance.GetMaxBackoff()
	balance.GetMaxRetryCount()
	balance.GetPaymentTransactionID()
	balance.GetQueryPayment()
	balance.GetMaxQueryPayment()
}

func TestUnitScheduleInfoQuerySetNothing(t *testing.T) {
	balance := NewScheduleInfoQuery()

	balance.GetScheduleID()
	balance.GetNodeAccountIDs()
	balance.GetMinBackoff()
	balance.GetMaxBackoff()
	balance.GetMaxRetryCount()
	balance.GetPaymentTransactionID()
	balance.GetQueryPayment()
	balance.GetMaxQueryPayment()
}

func TestUnitScheduleInfoQueryCoverage(t *testing.T) {
	checksum := "dmqui"
	grpc := time.Second * 3
	schedule := ScheduleID{Schedule: 3, checksum: &checksum}
	nodeAccountID := []AccountID{{Account: 10}}
	transactionID := TransactionIDGenerate(AccountID{Account: 324})

	client := ClientForTestnet()
	client.SetAutoValidateChecksums(true)

	query := NewScheduleInfoQuery().
		SetScheduleID(schedule).
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
	query.GetScheduleID()
	query.GetQueryPayment()
	query.GetMaxQueryPayment()
}

func TestUnitScheduleInfoQueryMock(t *testing.T) {
	responses := [][]interface{}{{
		&services.Response{
			Response: &services.Response_ScheduleGetInfo{
				ScheduleGetInfo: &services.ScheduleGetInfoResponse{
					Header: &services.ResponseHeader{NodeTransactionPrecheckCode: services.ResponseCodeEnum_OK, ResponseType: services.ResponseType_COST_ANSWER, Cost: 2},
				},
			},
		},
		&services.Response{
			Response: &services.Response_ScheduleGetInfo{
				ScheduleGetInfo: &services.ScheduleGetInfoResponse{
					Header: &services.ResponseHeader{NodeTransactionPrecheckCode: services.ResponseCodeEnum_OK, ResponseType: services.ResponseType_COST_ANSWER, Cost: 2},
				},
			},
		},
		&services.Response{
			Response: &services.Response_ScheduleGetInfo{
				ScheduleGetInfo: &services.ScheduleGetInfoResponse{
					Header: &services.ResponseHeader{NodeTransactionPrecheckCode: services.ResponseCodeEnum_OK, ResponseType: services.ResponseType_ANSWER_ONLY, Cost: 2},
					ScheduleInfo: &services.ScheduleInfo{
						ScheduleID:               nil,
						Data:                     nil,
						ExpirationTime:           nil,
						ScheduledTransactionBody: nil,
						Memo:                     "",
						AdminKey:                 nil,
						Signers:                  nil,
						CreatorAccountID:         nil,
						PayerAccountID:           nil,
						ScheduledTransactionID:   nil,
						LedgerId:                 nil,
						WaitForExpiry:            false,
					},
				},
			},
		},
	}}

	client, server := NewMockClientAndServer(responses)
	defer server.Close()

	query := NewScheduleInfoQuery().
		SetScheduleID(ScheduleID{Schedule: 3}).
		SetNodeAccountIDs([]AccountID{{Account: 3}}).
		SetMaxQueryPayment(NewHbar(1))

	_, err := query.GetCost(client)
	require.NoError(t, err)

	_, err = query.Execute(client)
	require.NoError(t, err)
}
