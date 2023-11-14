//go:build all || testnets
// +build all testnets

package hedera

/*-
 *
 * Hedera Go SDK
 *
 * Copyright (C) 2020 - 2023 Hedera Hashgraph, LLC
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
	"bytes"
	"fmt"
	"os"
	"regexp"
	"testing"
	"time"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
)

func TestIntegrationNodeForTransaction(t *testing.T) {
	t.Parallel()

	client := ClientForTestnet()
	operatorID,err:= AccountIDFromString(os.Getenv("OPERATOR_ID"))
	require.NoError(t, err)
	operatorKey, err := PrivateKeyFromString(os.Getenv("OPERATOR_KEY"))
	require.NoError(t, err)
	client.SetOperator(operatorID, operatorKey)

	var buf bytes.Buffer
	writer := zerolog.ConsoleWriter{Out: &buf, TimeFormat: time.RFC3339}
	zerolog.SetGlobalLevel(zerolog.TraceLevel)

	l := NewLogger("test", LoggerLevelTrace)
	l.SetLevel(LoggerLevelTrace)

	logger := zerolog.New(&writer)
	l.logger = &logger
	client.SetLogger(l)
	ledger, _ := LedgerIDFromNetworkName(NetworkNameTestnet)
	client.SetTransportSecurity(true)
	client.SetLedgerID(*ledger)
	client.SetMaxAttempts(3)
	nodeAccountIDs := map[string]struct{}{}
	for i := 0; i < 5; i++ {
		_, err := NewTransferTransaction().AddHbarTransfer(operatorID, HbarFromTinybar(-1)).
		AddHbarTransfer(AccountID{Shard: 0, Realm: 0, Account: 3}, HbarFromTinybar(1)).Execute(client)
		require.NoError(t, err)
		logOutput := buf.String()
		sanitizedLogOutput := regexp.MustCompile(`\x1b\[[0-9;]*m`).ReplaceAllString(logOutput, "")
		re := regexp.MustCompile(`nodeAccountID=([\d.]+)`)
		matches := re.FindStringSubmatch(sanitizedLogOutput)
		if len(matches) > 1 {
			nodeAccountID := matches[1]
			nodeAccountIDs[nodeAccountID] = struct{}{}
		} else {
			fmt.Println("Node Account ID not found")
		}
		buf.Reset()
	}
	require.True(t, len(nodeAccountIDs) > 1, "Expected multiple different node account IDs")
}

func TestIntegrationNodeForQuery(t *testing.T) {
	t.Parallel()

	client := ClientForTestnet()
	operatorID,err:= AccountIDFromString(os.Getenv("OPERATOR_ID"))
	require.NoError(t, err)
	operatorKey, err := PrivateKeyFromString(os.Getenv("OPERATOR_KEY"))
	require.NoError(t, err)
	client.SetOperator(operatorID, operatorKey)

	var buf bytes.Buffer
	writer := zerolog.ConsoleWriter{Out: &buf, TimeFormat: time.RFC3339}
	zerolog.SetGlobalLevel(zerolog.TraceLevel)

	l := NewLogger("test", LoggerLevelTrace)
	l.SetLevel(LoggerLevelTrace)

	logger := zerolog.New(&writer)
	l.logger = &logger
	client.SetLogger(l)
	ledger, _ := LedgerIDFromNetworkName(NetworkNameTestnet)
	client.SetTransportSecurity(true)
	client.SetLedgerID(*ledger)
	client.SetMaxAttempts(3)
	nodeAccountIDs := map[string]struct{}{}
	for i := 0; i < 5; i++ {
		_, err := NewAccountBalanceQuery().
			SetAccountID(AccountID{Account: 3}).
			Execute(client)
		require.NoError(t, err)
		logOutput := buf.String()
		sanitizedLogOutput := regexp.MustCompile(`\x1b\[[0-9;]*m`).ReplaceAllString(logOutput, "")
		re := regexp.MustCompile(`nodeAccountID=([\d.]+)`)
		matches := re.FindStringSubmatch(sanitizedLogOutput)
		if len(matches) > 1 {
			nodeAccountID := matches[1]
			nodeAccountIDs[nodeAccountID] = struct{}{}
		} else {
			fmt.Println("Node Account ID not found")
		}
		buf.Reset()
	}
	require.True(t, len(nodeAccountIDs) > 1, "Expected multiple different node account IDs")
}

func TestIntegrationNodeForTransactionSourceListUnchanged(t *testing.T) {
	t.Parallel()

	client := ClientForTestnet()
	operatorID,err:= AccountIDFromString(os.Getenv("OPERATOR_ID"))
	require.NoError(t, err)
	operatorKey, err := PrivateKeyFromString(os.Getenv("OPERATOR_KEY"))
	require.NoError(t, err)
	client.SetOperator(operatorID, operatorKey)

	

	var buf bytes.Buffer
	writer := zerolog.ConsoleWriter{Out: &buf, TimeFormat: time.RFC3339}
	zerolog.SetGlobalLevel(zerolog.TraceLevel)

	l := NewLogger("test", LoggerLevelTrace)
	l.SetLevel(LoggerLevelTrace)

	logger := zerolog.New(&writer)
	l.logger = &logger
	client.SetLogger(l)
	ledger, _ := LedgerIDFromNetworkName(NetworkNameTestnet)
	client.SetTransportSecurity(true)
	client.SetLedgerID(*ledger)
	client.SetMaxAttempts(3)
	
		_, err = NewAccountBalanceQuery().
			SetAccountID(AccountID{Account: 3}).
			Execute(client)
			expectedHealthyNodes := make([]_IManagedNode, len(client.network.healthyNodes))
	copy(expectedHealthyNodes, client.network.healthyNodes)
	resultHealthyNodes := make([]_IManagedNode, len(client.network.healthyNodes))
	_, err = NewAccountBalanceQuery().
			SetAccountID(AccountID{Account: 3}).
			Execute(client)
	copy(resultHealthyNodes, client.network.healthyNodes)
	require.Equal(t, expectedHealthyNodes, resultHealthyNodes)
}