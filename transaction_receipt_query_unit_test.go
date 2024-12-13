//go:build all || unit
// +build all unit

package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"encoding/hex"
	"encoding/json"
	"testing"

	"github.com/hiero-ledger/hiero-sdk-go/v2/proto/services"
	"github.com/stretchr/testify/assert"

	"github.com/stretchr/testify/require"
)

func TestUnitTransactionReceiptQueryValidate(t *testing.T) {
	t.Parallel()

	client, err := _NewMockClient()
	client.SetLedgerID(*NewLedgerIDTestnet())
	require.NoError(t, err)
	client.SetAutoValidateChecksums(true)
	accountID, err := AccountIDFromString("0.0.123-esxsf")
	transactionID := TransactionIDGenerate(accountID)
	require.NoError(t, err)

	receiptQuery := NewTransactionReceiptQuery().
		SetTransactionID(transactionID)

	err = receiptQuery.validateNetworkOnIDs(client)
	require.NoError(t, err)
}

func TestUnitTransactionReceiptQueryValidateWrong(t *testing.T) {
	t.Parallel()

	client, err := _NewMockClient()
	client.SetLedgerID(*NewLedgerIDTestnet())
	require.NoError(t, err)
	client.SetAutoValidateChecksums(true)
	accountID, err := AccountIDFromString("0.0.123-rmkykd")
	transactionID := TransactionIDGenerate(accountID)
	require.NoError(t, err)

	receiptQuery := NewTransactionReceiptQuery().
		SetTransactionID(transactionID)

	err = receiptQuery.validateNetworkOnIDs(client)
	assert.Error(t, err)
	if err != nil {
		assert.Equal(t, "network mismatch or wrong checksum given, given checksum: rmkykd, correct checksum esxsf, network: testnet", err.Error())
	}
}

func TestUnitTransactionReceiptQueryGet(t *testing.T) {
	t.Parallel()

	txID := TransactionIDGenerate(AccountID{Account: 7})

	balance := NewTransactionReceiptQuery().
		SetTransactionID(txID).
		SetIncludeDuplicates(true).
		SetIncludeChildren(true).
		SetQueryPayment(NewHbar(2)).
		SetMaxQueryPayment(NewHbar(1)).
		SetNodeAccountIDs([]AccountID{{Account: 10}, {Account: 11}, {Account: 12}})

	balance.GetTransactionID()
	balance.GetIncludeChildren()
	balance.GetIncludeDuplicates()
	balance.GetNodeAccountIDs()
	balance.GetMinBackoff()
	balance.GetMaxBackoff()
	balance.GetMaxRetryCount()
	balance.GetPaymentTransactionID()
	balance.GetQueryPayment()
	balance.GetMaxQueryPayment()
}

func TestUnitTransactionReceiptQueryNothingSet(t *testing.T) {
	t.Parallel()

	balance := NewTransactionReceiptQuery()

	balance.GetTransactionID()
	balance.GetIncludeChildren()
	balance.GetIncludeDuplicates()
	balance.GetNodeAccountIDs()
	balance.GetMinBackoff()
	balance.GetMaxBackoff()
	balance.GetMaxRetryCount()
	balance.GetPaymentTransactionID()
	balance.GetQueryPayment()
	balance.GetMaxQueryPayment()
}

func TestUnitTransactionPlatformNotActiveGracefulHandling(t *testing.T) {
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
						Status: services.ResponseCodeEnum_PLATFORM_NOT_ACTIVE,
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
	_, err = tx.SetValidateStatus(true).GetReceipt(client)
	require.NoError(t, err)
}

func TestUnitTransactionReceiptNotFound(t *testing.T) {
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
	receipt, err := tx.SetValidateStatus(true).GetReceipt(client)
	require.Error(t, err)
	require.Equal(t, "exceptional precheck status RECEIPT_NOT_FOUND", err.Error())
	require.Equal(t, StatusReceiptNotFound, receipt.Status)
}

func TestUnitTransactionReceiptUknown(t *testing.T) {
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
						Status: services.ResponseCodeEnum_UNKNOWN,
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
	receipt, err := tx.SetValidateStatus(true).GetReceipt(client)
	require.NoError(t, err)
	require.Equal(t, StatusSuccess, receipt.Status)
}

func TestUnitTransactionReceiptToJson(t *testing.T) {
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
						AccountID: &services.AccountID{Account: &services.AccountID_AccountNum{
							AccountNum: 123,
						}},
						ContractID: &services.ContractID{Contract: &services.ContractID_ContractNum{
							ContractNum: 456,
						}},
						FileID:        &services.FileID{FileNum: 789},
						TokenID:       &services.TokenID{TokenNum: 987},
						SerialNumbers: []int64{1, 2, 3},
						TopicID:       &services.TopicID{TopicNum: 654},
						ScheduleID:    &services.ScheduleID{ScheduleNum: 321},
						NodeId:        1,
						ScheduledTransactionID: &services.TransactionID{
							AccountID: &services.AccountID{Account: &services.AccountID_AccountNum{
								AccountNum: 123,
							}},
							TransactionValidStart: &services.Timestamp{
								Seconds: 1694689200,
							},
						},
						TopicSequenceNumber: 10,
						TopicRunningHash:    []byte{10},
						ExchangeRate: &services.ExchangeRateSet{
							CurrentRate: &services.ExchangeRate{
								HbarEquiv: 30000,
								CentEquiv: 154271,
								ExpirationTime: &services.TimestampSeconds{
									Seconds: 1694689200,
								},
							},
							NextRate: &services.ExchangeRate{
								HbarEquiv: 50000,
								CentEquiv: 154271,
								ExpirationTime: &services.TimestampSeconds{
									Seconds: 1694689200,
								},
							},
						},
					},
					ChildTransactionReceipts: []*services.TransactionReceipt{
						{
							Status: services.ResponseCodeEnum_SUCCESS,
							AccountID: &services.AccountID{Account: &services.AccountID_AccountNum{
								AccountNum: 123,
							}},
							ContractID: &services.ContractID{Contract: &services.ContractID_ContractNum{
								ContractNum: 456,
							}},
							FileID:        &services.FileID{FileNum: 789},
							TokenID:       &services.TokenID{TokenNum: 987},
							SerialNumbers: []int64{1, 2, 3},
							TopicID:       &services.TopicID{TopicNum: 654},
							ScheduleID:    &services.ScheduleID{ScheduleNum: 321},
							ScheduledTransactionID: &services.TransactionID{
								AccountID: &services.AccountID{Account: &services.AccountID_AccountNum{
									AccountNum: 123,
								}},
								TransactionValidStart: &services.Timestamp{
									Seconds: 1694689200,
								},
							},
							TopicSequenceNumber: 10,
							TopicRunningHash:    []byte{10},
							ExchangeRate: &services.ExchangeRateSet{
								CurrentRate: &services.ExchangeRate{
									HbarEquiv: 30000,
									CentEquiv: 154271,
									ExpirationTime: &services.TimestampSeconds{
										Seconds: 1694689200,
									},
								},
							},
						},
					},
					DuplicateTransactionReceipts: []*services.TransactionReceipt{
						{
							Status: services.ResponseCodeEnum_SUCCESS,
							AccountID: &services.AccountID{Account: &services.AccountID_AccountNum{
								AccountNum: 123,
							}},
							ContractID: &services.ContractID{Contract: &services.ContractID_ContractNum{
								ContractNum: 456,
							}},
							FileID:        &services.FileID{FileNum: 789},
							TokenID:       &services.TokenID{TokenNum: 987},
							SerialNumbers: []int64{1, 2, 3},
							TopicID:       &services.TopicID{TopicNum: 654},
							ScheduleID:    &services.ScheduleID{ScheduleNum: 321},
							ScheduledTransactionID: &services.TransactionID{
								AccountID: &services.AccountID{Account: &services.AccountID_AccountNum{
									AccountNum: 123,
								}},
								TransactionValidStart: &services.Timestamp{
									Seconds: 1694689200,
								},
							},
							TopicSequenceNumber: 10,
							TopicRunningHash:    []byte{10},
							ExchangeRate: &services.ExchangeRateSet{
								NextRate: &services.ExchangeRate{
									HbarEquiv: 50000,
									CentEquiv: 154271,
									ExpirationTime: &services.TimestampSeconds{
										Seconds: 1694689200,
									},
								},
							},
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
	require.NoError(t, err)
	receipt, err := tx.SetValidateStatus(true).GetReceipt(client)
	require.NoError(t, err)
	receiptBytes := receipt.ToBytes()
	receipt, err = TransactionReceiptFromBytes(receiptBytes)
	require.NoError(t, err)
	jsonBytes, err := receipt.MarshalJSON()
	require.NoError(t, err)
	expected := `
{
   "accountId":"0.0.123",
   "children":[
      {
         "accountId":"0.0.123",
         "children":[
            
         ],
         "contractId":"0.0.456",
         "duplicates":[
            
         ],
         "exchangeRate":{
            "hbars":30000,
            "cents":154271,
            "expirationTime":"2023-09-14T11:00:00.000Z"
         },
         "fileId":"0.0.789",
         "nodeId":0,
         "scheduleId":"0.0.321",
         "scheduledTransactionId":"0.0.123@1694689200.000000000",
         "serialNumbers":[
            1,
            2,
            3
         ],
         "status":"SUCCESS",
         "tokenId":"0.0.987",
         "topicId":"0.0.654",
         "topicRunningHash":"0a",
         "topicRunningHashVersion":0,
         "topicSequenceNumber":10,
         "totalSupply":0
      }
   ],
   "contractId":"0.0.456",
   "duplicates":[
      {
         "accountId":"0.0.123",
         "children":[
            
         ],
         "contractId":"0.0.456",
         "duplicates":[
            
         ],
         "fileId":"0.0.789",
         "nextExchangeRate":{
            "cents":154271,
            "expirationTime":"2023-09-14T11:00:00.000Z",
            "hbars":50000
         },
         "nodeId":0,
         "scheduleId":"0.0.321",
         "scheduledTransactionId":"0.0.123@1694689200.000000000",
         "serialNumbers":[
            1,
            2,
            3
         ],
         "status":"SUCCESS",
         "tokenId":"0.0.987",
         "topicId":"0.0.654",
         "topicRunningHash":"0a",
         "topicRunningHashVersion":0,
         "topicSequenceNumber":10,
         "totalSupply":0
      }
   ],
   "exchangeRate":{
      "hbars":30000,
      "cents":154271,
      "expirationTime":"2023-09-14T11:00:00.000Z"
   },
   "fileId":"0.0.789",
   "nextExchangeRate":{
      "cents":154271,
      "expirationTime":"2023-09-14T11:00:00.000Z",
      "hbars":50000
   },
   "nodeId":1,
   "scheduleId":"0.0.321",
   "scheduledTransactionId":"0.0.123@1694689200.000000000",
   "serialNumbers":[
      1,
      2,
      3
   ],
   "status":"SUCCESS",
   "tokenId":"0.0.987",
   "topicId":"0.0.654",
   "topicRunningHash":"0a",
   "topicRunningHashVersion":0,
   "topicSequenceNumber":10,
   "totalSupply":0
}
	`

	assert.JSONEqf(t, expected, string(jsonBytes), "json should be equal")

}

func TestUnitTransactionResponseToJson(t *testing.T) {
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
	}}
	client, server := NewMockClientAndServer(responses)
	defer server.Close()
	tx, err := NewTransferTransaction().
		SetNodeAccountIDs([]AccountID{{Account: 3}}).
		AddHbarTransfer(AccountID{Account: 2}, HbarFromTinybar(-1)).
		AddHbarTransfer(AccountID{Account: 3}, HbarFromTinybar(1)).
		Execute(client)
	require.NoError(t, err)
	jsonBytes, err := tx.MarshalJSON()
	require.NoError(t, err)
	obj := make(map[string]interface{})
	obj["nodeID"] = tx.NodeID.String()
	obj["hash"] = hex.EncodeToString(tx.Hash)
	obj["transactionID"] = tx.TransactionID.String()
	expectedJSON, err := json.Marshal(obj)
	require.NoError(t, err)
	assert.Equal(t, expectedJSON, jsonBytes)
}
