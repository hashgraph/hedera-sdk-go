//go:build all || unit
// +build all unit

package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"encoding/hex"
	"testing"
	"time"

	"github.com/hiero-ledger/hiero-sdk-go/v2/proto/services"
	"github.com/stretchr/testify/require"
)

func TestUnitTransactionRecordQueryValidate(t *testing.T) {
	t.Parallel()

	client, err := _NewMockClient()
	client.SetLedgerID(*NewLedgerIDTestnet())
	require.NoError(t, err)
	client.SetAutoValidateChecksums(true)
	accountID, err := AccountIDFromString("0.0.123-esxsf")
	transactionID := TransactionIDGenerate(accountID)
	require.NoError(t, err)

	recordQuery := NewTransactionRecordQuery().
		SetTransactionID(transactionID)

	err = recordQuery.validateNetworkOnIDs(client)
	require.NoError(t, err)
}

func TestUnitTransactionRecordQueryValidateWrong(t *testing.T) {
	t.Parallel()

	client, err := _NewMockClient()
	client.SetLedgerID(*NewLedgerIDTestnet())
	require.NoError(t, err)
	client.SetAutoValidateChecksums(true)
	accountID, err := AccountIDFromString("0.0.123-rmkykd")
	transactionID := TransactionIDGenerate(accountID)
	require.NoError(t, err)

	recordQuery := NewTransactionRecordQuery().
		SetTransactionID(transactionID)

	err = recordQuery.validateNetworkOnIDs(client)
	require.Error(t, err)
	if err != nil {
		require.Equal(t, "network mismatch or wrong checksum given, given checksum: rmkykd, correct checksum esxsf, network: testnet", err.Error())
	}
}

func TestUnitTransactionRecordQueryGet(t *testing.T) {
	t.Parallel()

	txID := TransactionIDGenerate(AccountID{Account: 7})
	deadline := time.Duration(time.Minute)
	accountId := AccountID{Account: 123}
	transactionID := TransactionIDGenerate(accountId)
	query := NewTransactionRecordQuery().
		SetTransactionID(txID).
		SetIncludeDuplicates(true).
		SetIncludeChildren(true).
		SetQueryPayment(NewHbar(2)).
		SetMaxQueryPayment(NewHbar(1)).
		SetQueryPayment(HbarFromTinybar(25)).
		SetNodeAccountIDs([]AccountID{{Account: 10}, {Account: 11}, {Account: 12}}).
		SetMaxRetry(3).
		SetMinBackoff(300 * time.Millisecond).
		SetMaxBackoff(10 * time.Second).
		SetPaymentTransactionID(transactionID).
		SetMaxQueryPayment(NewHbar(500)).
		SetGrpcDeadline(&deadline)

	require.Equal(t, txID, query.GetTransactionID())
	require.True(t, query.GetIncludeChildren())
	require.True(t, query.GetIncludeDuplicates())
	require.Equal(t, []AccountID{{Account: 10}, {Account: 11}, {Account: 12}}, query.GetNodeAccountIDs())
	require.Equal(t, 300*time.Millisecond, query.GetMinBackoff())
	require.Equal(t, 10*time.Second, query.GetMaxBackoff())
	require.Equal(t, 3, query.GetMaxRetryCount())
	require.Equal(t, transactionID, query.GetPaymentTransactionID())
	require.Equal(t, HbarFromTinybar(25), query.GetQueryPayment())
	require.Equal(t, NewHbar(500), query.GetMaxQueryPayment())
	require.Equal(t, &deadline, query.GetGrpcDeadline())
}

func TestUnitTransactionRecordQueryNothingSet(t *testing.T) {
	t.Parallel()

	query := NewTransactionRecordQuery()

	require.Equal(t, TransactionID{}, query.GetTransactionID())
	require.False(t, query.GetIncludeChildren())
	require.False(t, query.GetIncludeDuplicates())
	require.Empty(t, query.GetNodeAccountIDs())
	require.Equal(t, 250*time.Millisecond, query.GetMinBackoff())
	require.Equal(t, 8*time.Second, query.GetMaxBackoff())
	require.Equal(t, 10, query.GetMaxRetryCount())
	require.Equal(t, TransactionID{}, query.GetPaymentTransactionID())
	require.Equal(t, Hbar{}, query.GetQueryPayment())
	require.Equal(t, Hbar{}, query.GetMaxQueryPayment())
}

func TestUnitTransactionRecordPlatformNotActiveGracefulHandling(t *testing.T) {
	t.Parallel()

	responses := [][]interface{}{{
		&services.TransactionResponse{
			NodeTransactionPrecheckCode: services.ResponseCodeEnum_OK,
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
					},
				},
			},
		},
		&services.Response{
			Response: &services.Response_TransactionGetRecord{
				TransactionGetRecord: &services.TransactionGetRecordResponse{
					Header: &services.ResponseHeader{
						Cost:         0,
						ResponseType: services.ResponseType_ANSWER_ONLY,
					},
					TransactionRecord: &services.TransactionRecord{
						Receipt: &services.TransactionReceipt{
							Status: services.ResponseCodeEnum_PLATFORM_NOT_ACTIVE,
						},
					},
				},
			},
		},
		&services.Response{
			Response: &services.Response_TransactionGetRecord{
				TransactionGetRecord: &services.TransactionGetRecordResponse{
					Header: &services.ResponseHeader{
						Cost:         0,
						ResponseType: services.ResponseType_ANSWER_ONLY,
					},
					TransactionRecord: &services.TransactionRecord{
						Receipt: &services.TransactionReceipt{
							Status: services.ResponseCodeEnum_SUCCESS,
						},
					},
				},
			},
		},
		&services.Response{
			Response: &services.Response_TransactionGetRecord{
				TransactionGetRecord: &services.TransactionGetRecordResponse{
					Header: &services.ResponseHeader{
						Cost:         0,
						ResponseType: services.ResponseType_ANSWER_ONLY,
					},
					TransactionRecord: &services.TransactionRecord{
						Receipt: &services.TransactionReceipt{
							Status: services.ResponseCodeEnum_PLATFORM_NOT_ACTIVE,
						},
					},
				},
			},
		},
		&services.Response{
			Response: &services.Response_TransactionGetRecord{
				TransactionGetRecord: &services.TransactionGetRecordResponse{
					Header: &services.ResponseHeader{
						Cost:         0,
						ResponseType: services.ResponseType_ANSWER_ONLY,
					},
					TransactionRecord: &services.TransactionRecord{
						Receipt: &services.TransactionReceipt{
							Status: services.ResponseCodeEnum_SUCCESS,
						},
					},
				},
			},
		},
	}}

	client, server := NewMockClientAndServer(responses)
	defer server.Close()
	tx, err := NewTransferTransaction().
		SetNodeAccountIDs([]AccountID{{Account: 3}}).
		AddHbarTransfer(AccountID{Account: 2}, HbarFromTinybar(-1)).
		AddHbarTransfer(AccountID{Account: 3}, HbarFromTinybar(1)).
		Execute(client)
	client.SetMaxAttempts(2)
	require.NoError(t, err)
	_, err = tx.SetValidateStatus(true).GetRecord(client)
	require.NoError(t, err)
}

func TestUnitTransactionRecordReceiptNotFound(t *testing.T) {
	t.Parallel()

	responses := [][]interface{}{{
		&services.TransactionResponse{
			NodeTransactionPrecheckCode: services.ResponseCodeEnum_OK,
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
						Status: services.ResponseCodeEnum_RECEIPT_NOT_FOUND,
					},
				},
			},
		},
	}}
	client, server := NewMockClientAndServer(responses)
	defer server.Close()
	tx, err := NewTransferTransaction().
		SetNodeAccountIDs([]AccountID{{Account: 3}}).
		AddHbarTransfer(AccountID{Account: 2}, HbarFromTinybar(-1)).
		AddHbarTransfer(AccountID{Account: 3}, HbarFromTinybar(1)).
		Execute(client)
	client.SetMaxAttempts(2)
	require.NoError(t, err)
	record, err := tx.SetValidateStatus(true).GetRecord(client)
	require.Error(t, err)
	require.Equal(t, "exceptional precheck status RECEIPT_NOT_FOUND", err.Error())
	require.Equal(t, StatusReceiptNotFound, record.Receipt.Status)
}

func TestUnitTransactionRecordQueryMarshalJSON(t *testing.T) {
	t.Parallel()
	hexRecord, err := hex.DecodeString(`1afe010a26081612070800100018de092a130a110801100c1a0b0880ae99a4ffffffffff013800420058001230cac44f2db045ba441f3fbc295217f2eb0f956293d28b3401578f6160e66f4e47ea87952d91c4b1cb5bda6447823b979a1a0c08f3fcb495061083d9be900322190a0c08e8fcb495061098f09cf20112070800100018850918002a0030bee8f013526c0a0f0a0608001000180510d0df820118000a0f0a0608001000186210f08dff1e18000a100a070800100018a00610def1ef0318000a100a070800100018a10610def1ef0318000a110a070800100018850910fbf8b7e10718000a110a070800100018de091080a8d6b90718008a0100aa0100`)
	require.NoError(t, err)
	record, err := TransactionRecordFromBytes([]byte(hexRecord))
	require.NoError(t, err)
	accID, err := AccountIDFromString("0.0.1246")
	require.NoError(t, err)
	tokenID, err := TokenIDFromString("0.0.123")
	require.NoError(t, err)
	contractID, err := ContractIDFromString("0.0.3")
	require.NoError(t, err)
	record.Receipt.ContractID = &contractID
	record.Receipt.NodeID = 1

	tokenTransfer := TokenTransfer{
		AccountID:  accID,
		Amount:     789,
		IsApproved: true,
	}
	tokenTransferList := map[TokenID][]TokenTransfer{}
	tokenTransferList[tokenID] = []TokenTransfer{tokenTransfer}

	tokenNftTransfer := _TokenNftTransfer{
		SenderAccountID:   accID,
		ReceiverAccountID: accID,
		SerialNumber:      123,
		IsApproved:        true,
	}
	tokenNftTransferList := map[TokenID][]_TokenNftTransfer{}
	tokenNftTransferList[tokenID] = []_TokenNftTransfer{tokenNftTransfer}

	assessedCustomFee := AssessedCustomFee{
		FeeCollectorAccountId: &accID,
		Amount:                789,
		TokenID:               &tokenID,
		PayerAccountIDs:       []*AccountID{&accID},
	}

	tokenAssociation := TokenAssociation{
		AccountID: &accID,
		TokenID:   &tokenID,
	}

	plaidStaking := map[AccountID]Hbar{}
	for _, transfer := range record.Transfers {
		plaidStaking[transfer.AccountID] = transfer.Amount
	}

	pk, err := PublicKeyFromString("302a300506032b6570032100d7366c45e4d2f1a6c1d9af054f5ef8edc0b8d3875ba5d08a7f2e81ee8876e9e8")
	require.NoError(t, err)

	prngNumber := int32(123)
	evmAddressBytes, err := hex.DecodeString("deadbeef")
	require.NoError(t, err)

	record.TransactionMemo = "test"
	record.TokenTransfers = tokenTransferList
	record.NftTransfers = tokenNftTransferList
	record.ParentConsensusTimestamp = record.ConsensusTimestamp
	record.AliasKey = &pk
	record.EthereumHash = []byte{1, 2, 3, 4}
	record.PaidStakingRewards = plaidStaking
	record.PrngBytes = []byte{1, 2, 3, 4}
	record.PrngNumber = &prngNumber
	record.EvmAddress = evmAddressBytes
	record.AssessedCustomFees = []AssessedCustomFee{assessedCustomFee}
	record.AutomaticTokenAssociations = []TokenAssociation{tokenAssociation}
	record.PendingAirdropRecords = []PendingAirdropRecord{{pendingAirdropId: PendingAirdropId{&accID, &accID, &tokenID, nil}, pendingAirdropAmount: 789}}
	result, err := record.MarshalJSON()
	require.NoError(t, err)
	expected := `{
        "aliasKey":"302a300506032b6570032100d7366c45e4d2f1a6c1d9af054f5ef8edc0b8d3875ba5d08a7f2e81ee8876e9e8",
        "assessedCustomFees":[{"feeCollectorAccountId":"0.0.1246","tokenId":"0.0.123","amount":"789","payerAccountIds":["0.0.1246"]}],
        "automaticTokenAssociations":[{"tokenId":"0.0.123","accountId":"0.0.1246"}],
        "callResultIsCreate":true,
        "children":[],
        "consensusTimestamp":"2022-06-18T02:54:43.839Z",
        "duplicates":[],
        "ethereumHash":"01020304",
        "evmAddress":"deadbeef",
        "nftTransfers":{"0.0.123":[{"sender":"0.0.1246","recipient":"0.0.1246","isApproved":true,"serial":123}]},
        "paidStakingRewards":[
            {"accountId":"0.0.1157","amount":"-1041694270","isApproved":false},
            {"accountId":"0.0.1246","amount":"1000000000","isApproved":false},
            {"accountId":"0.0.5","amount":"1071080","isApproved":false},
            {"accountId":"0.0.800","amount":"4062319","isApproved":false},
            {"accountId":"0.0.801","amount":"4062319","isApproved":false},
            {"accountId":"0.0.98","amount":"32498552","isApproved":false}
        ],
        "parentConsensusTimestamp":"2022-06-18T02:54:43.839Z",
        "pendingAirdropRecords":[
            {
                "pendingAirdropAmount":"789",
                "pendingAirdropId":{
                    "nftId":"",
                    "receiver":"0.0.1246",
                    "sender":"0.0.1246",
                    "tokenId":"0.0.123"
                }
            }
        ],
        "prngBytes":"01020304",
        "prngNumber":123,
        "receipt":{
            "accountId":"0.0.1246",
            "children":[],
            "contractId":"0.0.3",
            "duplicates":[],
            "exchangeRate":{"cents":12,"expirationTime":"1963-11-25T17:31:44.000Z","hbars":1},
            "fileId":null,
            "nodeId":1,
            "scheduleId":null,
            "scheduledTransactionId":null,
            "serialNumbers":null,
            "status":"SUCCESS",
            "tokenId":null,
            "topicId":null,
            "topicRunningHash":"",
            "topicRunningHashVersion":0,
            "topicSequenceNumber":0,
            "totalSupply":0
        },
				"scheduleRef":"0.0.0",
        "tokenTransfers":{"0.0.123":{"0.0.1246":"789"}},
        "transactionFee":"41694270",
        "transactionHash":"cac44f2db045ba441f3fbc295217f2eb0f956293d28b3401578f6160e66f4e47ea87952d91c4b1cb5bda6447823b979a",
        "transactionId":"0.0.1157@1655520872.507983896",
        "transactionMemo":"test",
        "transfers":[
            {"accountId":"0.0.5","amount":"1071080","isApproved":false},
            {"accountId":"0.0.98","amount":"32498552","isApproved":false},
            {"accountId":"0.0.800","amount":"4062319","isApproved":false},
            {"accountId":"0.0.801","amount":"4062319","isApproved":false},
            {"accountId":"0.0.1157","amount":"-1041694270","isApproved":false},
            {"accountId":"0.0.1246","amount":"1000000000","isApproved":false}
        ]
    }`
	require.JSONEqf(t, expected, string(result), "json should be equal")
}
