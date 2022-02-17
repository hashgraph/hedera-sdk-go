//go:build all || unit
// +build all unit

package hedera

import (
	"bytes"
	"testing"

	"github.com/hashgraph/hedera-protobufs-go/services"

	"github.com/stretchr/testify/assert"

	"github.com/stretchr/testify/require"
)

func TestUnitFileContentsQueryValidate(t *testing.T) {
	client := ClientForTestnet()
	client.SetAutoValidateChecksums(true)
	fileID, err := FileIDFromString("0.0.123-rmkyk")
	require.NoError(t, err)

	fileContents := NewFileContentsQuery().
		SetFileID(fileID)

	err = fileContents._ValidateNetworkOnIDs(client)
	require.NoError(t, err)
}

func TestUnitFileContentsQueryValidateWrong(t *testing.T) {
	client := ClientForTestnet()
	client.SetAutoValidateChecksums(true)
	fileID, err := FileIDFromString("0.0.123-rmkykd")
	require.NoError(t, err)

	fileContents := NewFileContentsQuery().
		SetFileID(fileID)

	err = fileContents._ValidateNetworkOnIDs(client)
	assert.Error(t, err)
	if err != nil {
		assert.Equal(t, "network mismatch or wrong checksum given, given checksum: rmkykd, correct checksum rmkyk, network: testnet", err.Error())
	}
}

func TestUnitMockContractContentsQuery(t *testing.T) {
	responses := [][]interface{}{{
		&services.Response{
			Response: &services.Response_FileGetContents{
				FileGetContents: &services.FileGetContentsResponse{
					Header: &services.ResponseHeader{NodeTransactionPrecheckCode: services.ResponseCodeEnum_BUSY, ResponseType: services.ResponseType_ANSWER_ONLY},
				},
			},
		},
		&services.Response{
			Response: &services.Response_FileGetContents{
				FileGetContents: &services.FileGetContentsResponse{
					Header: &services.ResponseHeader{NodeTransactionPrecheckCode: services.ResponseCodeEnum_OK, ResponseType: services.ResponseType_ANSWER_ONLY, Cost: 2},
					FileContents: &services.FileGetContentsResponse_FileContents{
						FileID:   &services.FileID{FileNum: 3},
						Contents: []byte{123},
					},
				},
			},
		},
	}}

	client, server := NewMockClientAndServer(responses)

	result, err := NewFileContentsQuery().
		SetFileID(FileID{File: 3}).
		SetMaxQueryPayment(NewHbar(1)).
		SetQueryPayment(HbarFromTinybar(25)).
		SetNodeAccountIDs([]AccountID{{Account: 3}}).
		Execute(client)
	require.NoError(t, err)

	require.Equal(t, bytes.Compare(result, []byte{123}), 0)

	server.Close()
}
