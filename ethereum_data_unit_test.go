//go:build all || unit
// +build all unit

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
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUnitEthereumData(t *testing.T) {
	t.Parallel()

	byt, err := hex.DecodeString("02f87082012a022f2f83018000947e3a9eaf9bcc39e2ffa38eb30bf7a93feacbc181880de0b6b3a764000083123456c001a0df48f2efd10421811de2bfb125ab75b2d3c44139c4642837fb1fccce911fd479a01aaf7ae92bee896651dfc9d99ae422a296bf5d9f1ca49b2d96d82b79eb112d66")
	require.NoError(t, err)
	b, err := EthereumTransactionDataFromBytes(byt)
	require.NoError(t, err)
	k, err := b.ToBytes()
	require.Equal(t, hex.EncodeToString(k), "02f87082012a022f2f83018000947e3a9eaf9bcc39e2ffa38eb30bf7a93feacbc181880de0b6b3a764000083123456c001a0df48f2efd10421811de2bfb125ab75b2d3c44139c4642837fb1fccce911fd479a01aaf7ae92bee896651dfc9d99ae422a296bf5d9f1ca49b2d96d82b79eb112d66")
}

func TestUnitEthereumJson(t *testing.T) {
	t.Parallel()

	byt, err := hex.DecodeString("02f87082012a022f2f83018000947e3a9eaf9bcc39e2ffa38eb30bf7a93feacbc181880de0b6b3a764000083123456c001a0df48f2efd10421811de2bfb125ab75b2d3c44139c4642837fb1fccce911fd479a01aaf7ae92bee896651dfc9d99ae422a296bf5d9f1ca49b2d96d82b79eb112d66")
	require.NoError(t, err)
	b, err := EthereumTransactionDataFromBytes(byt)
	require.NoError(t, err)
	k, err := b.ToJson()
	require.Equal(t, string(k), "{\"ChainID\":298,\"Nonce\":2,\"GasTipCap\":47,\"GasFeeCap\":47,\"Gas\":98304,\"To\":\"0x7e3a9eaf9bcc39e2ffa38eb30bf7a93feacbc181\",\"Value\":1000000000000000000,\"Data\":\"EjRW\",\"AccessList\":[],\"v\":1,\"r\":100994654910787593347140909067229841565925886199208140719315441876745360233593,\"s\":12070180598849582105957204602453860723012110108194703753759584093273689501030}")
	second, err := EthereumTransactionDataFromJson(k)
	require.NoError(t, err)
	bytSecond, err := second.ToBytes()
	require.NoError(t, err)
	require.Equal(t, hex.EncodeToString(bytSecond), "02f87082012a022f2f83018000947e3a9eaf9bcc39e2ffa38eb30bf7a93feacbc181880de0b6b3a764000083123456c001a0df48f2efd10421811de2bfb125ab75b2d3c44139c4642837fb1fccce911fd479a01aaf7ae92bee896651dfc9d99ae422a296bf5d9f1ca49b2d96d82b79eb112d66")
}
