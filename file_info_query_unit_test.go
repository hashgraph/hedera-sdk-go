//go:build all || unit
// +build all unit

package hedera

import (
	"testing"

	"github.com/hashgraph/hedera-protobufs-go/services"

	"github.com/stretchr/testify/assert"

	"github.com/stretchr/testify/require"
)

func TestUnitFileInfoQueryValidate(t *testing.T) {
	client := ClientForTestnet()
	client.SetAutoValidateChecksums(true)
	fileID, err := FileIDFromString("0.0.123-rmkyk")
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
		assert.Equal(t, "network mismatch or wrong checksum given, given checksum: rmkykd, correct checksum rmkyk, network: testnet", err.Error())
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

	server.Close()
}
