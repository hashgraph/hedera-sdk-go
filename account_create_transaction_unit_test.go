//go:build all || unit
// +build all unit

package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"encoding/hex"
	"testing"
	"time"

	"github.com/hiero-ledger/hiero-sdk-go/v2/proto/services"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/stretchr/testify/assert"

	"github.com/stretchr/testify/require"
)

func TestUnitAccountCreateTransactionValidate(t *testing.T) {
	t.Parallel()

	client, err := _NewMockClient()
	client.SetLedgerID(*NewLedgerIDTestnet())
	require.NoError(t, err)
	client.SetAutoValidateChecksums(true)
	accountID, err := AccountIDFromString("0.0.123-esxsf")
	require.NoError(t, err)

	createAccount := NewAccountCreateTransaction().
		SetProxyAccountID(accountID)

	err = createAccount.validateNetworkOnIDs(client)
	require.NoError(t, err)
}

func TestUnitAccountCreateTransactionValidateWrong(t *testing.T) {
	t.Parallel()

	client, err := _NewMockClient()
	client.SetLedgerID(*NewLedgerIDTestnet())
	require.NoError(t, err)
	client.SetAutoValidateChecksums(true)
	accountID, err := AccountIDFromString("0.0.123-rmkykd")
	require.NoError(t, err)

	createAccount := NewAccountCreateTransaction().
		SetProxyAccountID(accountID)

	err = createAccount.validateNetworkOnIDs(client)
	assert.Error(t, err)
	if err != nil {
		assert.Equal(t, "network mismatch or wrong checksum given, given checksum: rmkykd, correct checksum esxsf, network: testnet", err.Error())
	}
}

func TestUnitAccountCreateTransactionMock(t *testing.T) {
	t.Parallel()

	responses := [][]interface{}{{
		status.New(codes.Unavailable, "node is UNAVAILABLE").Err(),
		status.New(codes.Internal, "Received RST_STREAM with code 0").Err(),
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
	defer server.Close()

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
}

func TestUnitAccountCreateTransactionGet(t *testing.T) {
	t.Parallel()

	nodeAccountID := []AccountID{{Account: 10}, {Account: 11}, {Account: 12}}

	key, err := PrivateKeyGenerateEd25519()

	transactionID := TransactionIDGenerate(AccountID{Account: 324})

	transaction, err := NewAccountCreateTransaction().
		SetTransactionID(transactionID).
		SetNodeAccountIDs(nodeAccountID).
		SetKey(key).
		SetAccountMemo("").
		SetReceiverSignatureRequired(true).
		SetMaxAutomaticTokenAssociations(2).
		SetAutoRenewPeriod(60 * time.Second).
		SetTransactionMemo("").
		SetTransactionValidDuration(60 * time.Second).
		Freeze()
	require.NoError(t, err)

	transaction.GetTransactionID()
	transaction.GetNodeAccountIDs()

	_, err = transaction.GetTransactionHash()
	require.NoError(t, err)

	transaction.GetMaxTransactionFee()
	transaction.GetTransactionMemo()
	transaction.GetRegenerateTransactionID()
	_, err = transaction.GetSignatures()
	require.NoError(t, err)
	transaction.GetRegenerateTransactionID()
	transaction.GetAccountMemo()
	transaction.GetMaxTransactionFee()
	transaction.GetMaxAutomaticTokenAssociations()
	transaction.GetRegenerateTransactionID()
	transaction.GetKey()
	transaction.GetInitialBalance()
	transaction.GetAutoRenewPeriod()
	transaction.GetReceiverSignatureRequired()
}

func TestUnitAccountCreateTransactionSetNothing(t *testing.T) {
	t.Parallel()

	nodeAccountID := []AccountID{{Account: 10}, {Account: 11}, {Account: 12}}
	transactionID := TransactionIDGenerate(AccountID{Account: 324})

	transaction, err := NewAccountCreateTransaction().
		SetTransactionID(transactionID).
		SetNodeAccountIDs(nodeAccountID).
		Freeze()
	require.NoError(t, err)

	transaction.GetTransactionID()
	transaction.GetNodeAccountIDs()

	_, err = transaction.GetTransactionHash()
	require.NoError(t, err)

	transaction.GetMaxTransactionFee()
	transaction.GetTransactionMemo()
	transaction.GetRegenerateTransactionID()
	_, err = transaction.GetSignatures()
	require.NoError(t, err)
	transaction.GetRegenerateTransactionID()
	transaction.GetAccountMemo()
	transaction.GetMaxTransactionFee()
	transaction.GetMaxAutomaticTokenAssociations()
	transaction.GetProxyAccountID()
	transaction.GetRegenerateTransactionID()
	transaction.GetKey()
	transaction.GetInitialBalance()
	transaction.GetAutoRenewPeriod()
	transaction.GetReceiverSignatureRequired()
}

func TestUnitAccountCreateTransactionProtoCheck(t *testing.T) {
	t.Parallel()

	nodeAccountID := []AccountID{{Account: 10}, {Account: 11}, {Account: 12}}
	stackedAccountID := AccountID{Account: 5}

	key, err := PrivateKeyGenerateEd25519()
	require.NoError(t, err)

	transactionID := TransactionIDGenerate(AccountID{Account: 324})

	alias := "5c562e90feaf0eebd33ea75d21024f249d451417"

	transaction, err := NewAccountCreateTransaction().
		SetTransactionID(transactionID).
		SetNodeAccountIDs(nodeAccountID).
		SetKey(key).
		SetInitialBalance(NewHbar(3)).
		SetAccountMemo("ty").
		SetReceiverSignatureRequired(true).
		SetMaxAutomaticTokenAssociations(2).
		SetStakedAccountID(stackedAccountID).
		SetDeclineStakingReward(true).
		SetAutoRenewPeriod(60 * time.Second).
		SetTransactionMemo("").
		SetTransactionValidDuration(60 * time.Second).
		SetAlias(alias).
		Freeze()
	require.NoError(t, err)

	transaction.GetTransactionID()
	transaction.GetNodeAccountIDs()

	proto := transaction.build().GetCryptoCreateAccount()
	require.Equal(t, proto.Key.String(), key._ToProtoKey().String())
	require.Equal(t, proto.InitialBalance, uint64(NewHbar(3).AsTinybar()))
	require.Equal(t, proto.Memo, "ty")
	require.Equal(t, proto.ReceiverSigRequired, true)
	require.Equal(t, proto.MaxAutomaticTokenAssociations, int32(2))
	require.Equal(t, proto.StakedId.(*services.CryptoCreateTransactionBody_StakedAccountId).StakedAccountId.String(),
		stackedAccountID._ToProtobuf().String())
	require.Equal(t, proto.DeclineReward, true)
	require.Equal(t, proto.AutoRenewPeriod.String(), _DurationToProtobuf(60*time.Second).String())
	require.Equal(t, hex.EncodeToString(proto.Alias), alias)
}

func TestUnitAccountCreateTransactionCoverage(t *testing.T) {
	t.Parallel()

	checksum := "dmqui"
	account := AccountID{Account: 3, checksum: &checksum}
	nodeAccountID := []AccountID{{Account: 10}}
	transactionID := TransactionIDGenerate(AccountID{Account: 324})

	key, err := PrivateKeyGenerateEd25519()
	require.NoError(t, err)

	alias := "5c562e90feaf0eebd33ea75d21024f249d451417"

	client, err := _NewMockClient()
	client.SetLedgerID(*NewLedgerIDTestnet())
	require.NoError(t, err)
	client.SetAutoValidateChecksums(true)

	trx, err := NewAccountCreateTransaction().
		SetTransactionID(transactionID).
		SetNodeAccountIDs(nodeAccountID).
		SetKey(key).
		SetInitialBalance(NewHbar(3)).
		SetAccountMemo("ty").
		SetReceiverSignatureRequired(true).
		SetMaxAutomaticTokenAssociations(2).
		SetStakedAccountID(account).
		SetStakedNodeID(4).
		SetDeclineStakingReward(true).
		SetAutoRenewPeriod(60 * time.Second).
		SetTransactionMemo("").
		SetTransactionValidDuration(60 * time.Second).
		SetMaxTransactionFee(NewHbar(3)).
		SetMaxRetry(3).
		SetMaxBackoff(time.Second * 30).
		SetMinBackoff(time.Second * 10).
		SetTransactionMemo("no").
		SetTransactionValidDuration(time.Second * 30).
		SetRegenerateTransactionID(false).
		SetAlias(alias).
		Freeze()
	require.NoError(t, err)

	trx.validateNetworkOnIDs(client)
	_, err = trx.Schedule()
	require.NoError(t, err)
	trx.GetTransactionID()
	trx.GetNodeAccountIDs()
	trx.GetMaxRetry()
	trx.GetMaxTransactionFee()
	trx.GetMaxBackoff()
	trx.GetMinBackoff()
	trx.GetRegenerateTransactionID()
	byt, err := trx.ToBytes()
	require.NoError(t, err)
	txFromBytes, err := TransactionFromBytes(byt)
	require.NoError(t, err)
	sig, err := key.SignTransaction(trx)
	require.NoError(t, err)

	_, err = trx.GetTransactionHash()
	require.NoError(t, err)
	trx.GetMaxTransactionFee()
	trx.GetTransactionMemo()
	trx.GetRegenerateTransactionID()
	trx.GetStakedAccountID()
	trx.GetStakedNodeID()
	trx.GetDeclineStakingReward()
	trx.GetAlias()
	_, err = trx.GetSignatures()
	require.NoError(t, err)
	trx.getName()
	switch b := txFromBytes.(type) {
	case AccountCreateTransaction:
		b.AddSignature(key.PublicKey(), sig)
	}
}

func TestUnitAccountCreateSetStakedNodeID(t *testing.T) {
	t.Parallel()

	checksum := "dmqui"
	account := AccountID{Account: 3, checksum: &checksum}
	tx := NewAccountCreateTransaction().
		SetStakedAccountID(account).
		SetStakedNodeID(4)

	require.Equal(t, AccountID{}, tx.GetStakedAccountID())
	require.Equal(t, int64(4), tx.GetStakedNodeID())
}
func TestUnitAccountCreateSetStakedAccountID(t *testing.T) {
	t.Parallel()

	checksum := "dmqui"
	account := AccountID{Account: 3, checksum: &checksum}
	tx := NewAccountCreateTransaction().
		SetStakedNodeID(4).
		SetStakedAccountID(account)

	require.Equal(t, int64(0), tx.GetStakedNodeID())
	require.Equal(t, account, tx.GetStakedAccountID())
}
