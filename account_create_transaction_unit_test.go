//go:build all || unit
// +build all unit

package hedera

import (
	"testing"

	"github.com/hashgraph/hedera-protobufs-go/services"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/stretchr/testify/assert"

	"github.com/stretchr/testify/require"
)

func TestUnitAccountCreateTransactionValidate(t *testing.T) {
	client := ClientForTestnet()
	client.SetAutoValidateChecksums(true)
	accountID, err := AccountIDFromString("0.0.123-rmkyk")
	require.NoError(t, err)

	createAccount := NewAccountCreateTransaction().
		SetProxyAccountID(accountID)

	err = createAccount._ValidateNetworkOnIDs(client)
	require.NoError(t, err)
}

func TestUnitAccountCreateTransactionValidateWrong(t *testing.T) {
	client := ClientForTestnet()
	client.SetAutoValidateChecksums(true)
	accountID, err := AccountIDFromString("0.0.123-rmkykd")
	require.NoError(t, err)

	createAccount := NewAccountCreateTransaction().
		SetProxyAccountID(accountID)

	err = createAccount._ValidateNetworkOnIDs(client)
	assert.Error(t, err)
	if err != nil {
		assert.Equal(t, "network mismatch or wrong checksum given, given checksum: rmkykd, correct checksum rmkyk, network: testnet", err.Error())
	}
}

func TestUnitMockAccountCreateTransaction(t *testing.T) {
	responses := [][]interface{}{{
		status.New(codes.Unavailable, "node is UNAVAILABLE").Err(),
		status.New(codes.Internal, "Received RST_STREAM with code 0").Err(),
	}, {
		&services.TransactionResponse{
			NodeTransactionPrecheckCode: services.ResponseCodeEnum_BUSY,
		},
		&services.TransactionResponse{
			NodeTransactionPrecheckCode: services.ResponseCodeEnum_OK,
		},
		&services.Response{
			Response: &services.Response_TransactionGetReceipt{
				TransactionGetReceipt: &services.TransactionGetReceiptResponse{
					Header: &services.ResponseHeader{
						Cost:         0,
						ResponseType: services.ResponseType_COST_ANSWER,
					},
				},
			},
		},
		&services.Response{
			Response: &services.Response_TransactionGetReceipt{
				TransactionGetReceipt: &services.TransactionGetReceiptResponse{
					Header: &services.ResponseHeader{
						Cost:         0,
						ResponseType: services.ResponseType_ANSWER_ONLY,
					},
					Receipt: &services.TransactionReceipt{
						Status: services.ResponseCodeEnum_RECEIPT_NOT_FOUND,
					},
				},
			},
		},
		&services.Response{
			Response: &services.Response_TransactionGetReceipt{
				TransactionGetReceipt: &services.TransactionGetReceiptResponse{
					Header: &services.ResponseHeader{
						Cost:         0,
						ResponseType: services.ResponseType_ANSWER_ONLY,
					},
					Receipt: &services.TransactionReceipt{
						Status: services.ResponseCodeEnum_SUCCESS,
						AccountID: &services.AccountID{Account: &services.AccountID_AccountNum{
							AccountNum: 234,
						}},
					},
				},
			},
		},
	}}

	client, server := NewMockClientAndServer(responses)

	newKey, err := PrivateKeyGenerateEd25519()
	require.NoError(t, err)

	newBalance := NewHbar(2)

	tran := TransactionIDGenerate(AccountID{Account: 3})

	resp, err := NewAccountCreateTransaction().
		SetNodeAccountIDs([]AccountID{{Account: 3}, {Account: 4}}).
		SetKey(newKey).
		SetTransactionID(tran).
		SetInitialBalance(newBalance).
		SetMaxAutomaticTokenAssociations(100).
		Execute(client)
	require.NoError(t, err)

	receipt, err := resp.GetReceipt(client)
	require.NoError(t, err)
	require.Equal(t, receipt.AccountID, &AccountID{Account: 234})
	server.Close()
}
