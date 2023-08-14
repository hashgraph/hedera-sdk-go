//go:build all || e2e
// +build all e2e

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
	"os"
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/require"
)

var callCount int64

func incrementCallCount() {
	atomic.AddInt64(&callCount, 1)
}

func getCallCount() int64 {
	return atomic.LoadInt64(&callCount)
}

func TestIntegrationOneSignature(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)

	client := ClientForNetwork(env.Client.GetNetwork()).SetOperatorWith(env.OriginalOperatorID, env.OriginalOperatorKey, signingServiceTwo)
	response, err := NewTransferTransaction().
		AddHbarTransfer(env.OriginalOperatorID, NewHbar(-1)).
		AddHbarTransfer(AccountID{Account: 3}, NewHbar(1)).
		SetNodeAccountIDs(env.NodeAccountIDs).
		Execute(client)
	require.NoError(t, err)

	_, err = response.GetReceipt(client)
	require.NoError(t, err)

	require.Equal(t, int64(1), getCallCount())
	client.Close()
	err = CloseIntegrationTestEnv(env, nil)
	require.NoError(t, err)
}

func signingServiceTwo(txBytes []byte) []byte {
	localOperatorPrivateKey, _ := PrivateKeyFromString(os.Getenv("OPERATOR_KEY"))
	incrementCallCount()

	signature := localOperatorPrivateKey.Sign(txBytes)
	return signature
}
