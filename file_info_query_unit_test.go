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
	"testing"

	"github.com/hashgraph/hedera-protobufs-go/services"

	"github.com/stretchr/testify/assert"

	"github.com/stretchr/testify/require"
)

func TestUnitFileInfoQueryValidate(t *testing.T) {
	client := ClientForTestnet()
	client.SetAutoValidateChecksums(true)
	fileID, err := FileIDFromString("0.0.123-esxsf")
	require.NoError(t, err)

	fileInfo := NewFileInfoQuery().
		SetFileID(fileID)

	err = fileInfo._ValidateNetworkOnIDs(client)
	require.NoError(t, err)
}

func TestUnitFileInfoQueryValidateWrong(t *testing.T) {
	client := ClientForTestnet()
	client.SetAutoValidateChecksums(true)
	fileID, err := FileIDFromString("0.0.123-rmkykd")
	require.NoError(t, err)

	fileInfo := NewFileInfoQuery().
		SetFileID(fileID)

	err = fileInfo._ValidateNetworkOnIDs(client)
	assert.Error(t, err)
	if err != nil {
		assert.Equal(t, "network mismatch or wrong checksum given, given checksum: rmkykd, correct checksum esxsf, network: testnet", err.Error())
	}
}

func TestUnitMockFileInfoQuery(t *testing.T) {
	newKey, err := PrivateKeyFromStringEd25519("302e020100300506032b657004220420a869f4c6191b9c8c99933e7f6b6611711737e4b1a1a5a4cb5370e719a1f6df98")
	require.NoError(t, err)
	key := newKey.PublicKey().BytesRaw()

	responses := [][]interface{}{{
		&services.Response{
			Response: &services.Response_FileGetInfo{
				FileGetInfo: &services.FileGetInfoResponse{
					Header: &services.ResponseHeader{NodeTransactionPrecheckCode: services.ResponseCodeEnum_BUSY, ResponseType: services.ResponseType_ANSWER_ONLY},
				},
			},
		},
		&services.Response{
			Response: &services.Response_FileGetInfo{
				FileGetInfo: &services.FileGetInfoResponse{
					Header: &services.ResponseHeader{NodeTransactionPrecheckCode: services.ResponseCodeEnum_OK, ResponseType: services.ResponseType_ANSWER_ONLY, Cost: 2},
					FileInfo: &services.FileGetInfoResponse_FileInfo{
						FileID:         &services.FileID{FileNum: 3},
						Size:           10,
						ExpirationTime: nil,
						Deleted:        false,
						Keys: &services.KeyList{
							Keys: []*services.Key{
								{
									Key: &services.Key_Ed25519{
										Ed25519: key,
									},
								},
							},
						},
						Memo:     "no memo",
						LedgerId: []byte{0},
					},
				},
			},
		},
	}}

	client, server := NewMockClientAndServer(responses)
	defer server.Close()

	result, err := NewFileInfoQuery().
		SetFileID(FileID{File: 3}).
		SetNodeAccountIDs([]AccountID{{Account: 3}}).
		SetMaxQueryPayment(NewHbar(1)).
		SetQueryPayment(HbarFromTinybar(25)).
		Execute(client)
	require.NoError(t, err)

	require.Equal(t, result.Keys.keys[0].String(), newKey.PublicKey().String())
	require.Equal(t, result.FileMemo, "no memo")
	require.Equal(t, result.IsDeleted, false)
	require.True(t, result.LedgerID.IsMainnet())
}
