//go:build all || unit
// +build all unit

package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/hiero-ledger/hiero-sdk-go/v2/proto/sdk"
	"github.com/hiero-ledger/hiero-sdk-go/v2/proto/services"
	protobuf "google.golang.org/protobuf/proto"

	"github.com/stretchr/testify/assert"

	"github.com/stretchr/testify/require"
)

func TestUnitTransactionSerializationDeserialization(t *testing.T) {
	t.Parallel()

	transaction, err := _NewMockTransaction()
	require.NoError(t, err)

	_, err = transaction.Freeze()
	require.NoError(t, err)

	_, err = transaction.GetSignatures()
	require.NoError(t, err)

	_, err = transaction.GetTransactionHash()
	require.NoError(t, err)

	transaction.
		SetTransactionMemo("memo").
		SetMaxTransactionFee(NewHbar(5))

	txBytes, err := transaction.ToBytes()
	require.NoError(t, err)

	deserializedTX, err := TransactionFromBytes(txBytes)
	require.NoError(t, err)

	var deserializedTXTyped TransferTransaction
	switch tx := deserializedTX.(type) {
	case TransferTransaction:
		deserializedTXTyped = tx
	default:
		panic("Transaction was not TransferTransaction")
	}

	require.Equal(t, "memo", deserializedTXTyped.memo)
	require.Equal(t, NewHbar(5), deserializedTXTyped.GetMaxTransactionFee())
	assert.Equal(t, transaction.String(), deserializedTXTyped.String())
}

func TestUnitTransactionValidateBodiesEqual(t *testing.T) {
	t.Parallel()

	key, err := PrivateKeyFromString(mockPrivateKey)
	require.NoError(t, err)
	transaction := services.TransactionBody{
		TransactionID:            testTransactionID._ToProtobuf(),
		NodeAccountID:            AccountID{Account: 3}._ToProtobuf(),
		TransactionFee:           0,
		TransactionValidDuration: nil,
		GenerateRecord:           false,
		Memo:                     "",
		Data: &services.TransactionBody_CryptoCreateAccount{
			CryptoCreateAccount: &services.CryptoCreateTransactionBody{
				Key:                           key._ToProtoKey(),
				InitialBalance:                0,
				ProxyAccountID:                AccountID{Account: 123}._ToProtobuf(),
				SendRecordThreshold:           0,
				ReceiveRecordThreshold:        0,
				ReceiverSigRequired:           false,
				AutoRenewPeriod:               nil,
				ShardID:                       nil,
				RealmID:                       nil,
				NewRealmAdminKey:              nil,
				Memo:                          "",
				MaxAutomaticTokenAssociations: 0,
			},
		},
	}

	transactionBody, err := protobuf.Marshal(&transaction)
	require.NoError(t, err)

	signed, err := protobuf.Marshal(&services.SignedTransaction{
		BodyBytes: transactionBody,
	})
	require.NoError(t, err)
	list, err := protobuf.Marshal(&sdk.TransactionList{
		TransactionList: []*services.Transaction{
			{
				SignedTransactionBytes: signed,
			},
			{
				SignedTransactionBytes: signed,
			},
			{
				SignedTransactionBytes: signed,
			},
		},
	})

	deserializedTX, err := TransactionFromBytes(list)
	require.NoError(t, err)

	var deserializedTXTyped AccountCreateTransaction
	switch tx := deserializedTX.(type) {
	case AccountCreateTransaction:
		deserializedTXTyped = tx
	default:
		panic("Transaction was not AccountCreateTransaction")
	}

	assert.Equal(t, uint64(transaction.TransactionID.AccountID.GetAccountNum()), deserializedTXTyped.GetTransactionID().AccountID.Account)
}

func DisabledTestUnitTransactionValidateBodiesNotEqual(t *testing.T) {
	t.Parallel()

	key, err := PrivateKeyFromString(mockPrivateKey)
	require.NoError(t, err)
	transaction := services.TransactionBody{
		TransactionID:            testTransactionID._ToProtobuf(),
		NodeAccountID:            AccountID{Account: 3}._ToProtobuf(),
		TransactionFee:           0,
		TransactionValidDuration: nil,
		GenerateRecord:           false,
		Memo:                     "",
		Data: &services.TransactionBody_CryptoCreateAccount{
			CryptoCreateAccount: &services.CryptoCreateTransactionBody{
				Key:                           key._ToProtoKey(),
				InitialBalance:                0,
				ProxyAccountID:                AccountID{Account: 123}._ToProtobuf(),
				SendRecordThreshold:           0,
				ReceiveRecordThreshold:        0,
				ReceiverSigRequired:           false,
				AutoRenewPeriod:               nil,
				ShardID:                       nil,
				RealmID:                       nil,
				NewRealmAdminKey:              nil,
				Memo:                          "",
				MaxAutomaticTokenAssociations: 0,
			},
		},
	}

	transaction2 := services.TransactionBody{
		TransactionID:            testTransactionID._ToProtobuf(),
		NodeAccountID:            AccountID{Account: 3}._ToProtobuf(),
		TransactionFee:           0,
		TransactionValidDuration: nil,
		GenerateRecord:           false,
		Memo:                     "",
		Data: &services.TransactionBody_CryptoCreateAccount{
			CryptoCreateAccount: &services.CryptoCreateTransactionBody{
				Key:                           key._ToProtoKey(),
				InitialBalance:                0,
				ProxyAccountID:                AccountID{Account: 1}._ToProtobuf(),
				SendRecordThreshold:           0,
				ReceiveRecordThreshold:        0,
				ReceiverSigRequired:           false,
				AutoRenewPeriod:               nil,
				ShardID:                       nil,
				RealmID:                       nil,
				NewRealmAdminKey:              nil,
				Memo:                          "",
				MaxAutomaticTokenAssociations: 0,
			},
		},
	}

	transactionBody, err := protobuf.Marshal(&transaction)
	require.NoError(t, err)

	signed, err := protobuf.Marshal(&services.SignedTransaction{
		BodyBytes: transactionBody,
	})

	transactionBody2, err := protobuf.Marshal(&transaction2)
	require.NoError(t, err)

	signed2, err := protobuf.Marshal(&services.SignedTransaction{
		BodyBytes: transactionBody2,
	})

	require.NoError(t, err)
	list, err := protobuf.Marshal(&sdk.TransactionList{
		TransactionList: []*services.Transaction{
			{
				SignedTransactionBytes: signed,
			},
			{
				SignedTransactionBytes: signed2,
			},
			{
				SignedTransactionBytes: signed2,
			},
		},
	})

	_, err = TransactionFromBytes(list)
	require.Error(t, err)
	if err != nil {
		assert.Equal(t, fmt.Sprintf("failed to validate transaction bodies"), err.Error())
	}
}

func TestUnitTransactionToFromBytes(t *testing.T) {
	t.Parallel()

	duration := time.Second * 10
	operatorID := AccountID{Account: 5}
	recepientID := AccountID{Account: 4}
	node := []AccountID{{Account: 3}}
	transaction, err := NewTransferTransaction().
		SetTransactionID(testTransactionID).
		SetNodeAccountIDs(node).
		AddHbarTransfer(operatorID, NewHbar(-1)).
		AddHbarTransfer(recepientID, NewHbar(1)).
		SetTransactionMemo("go sdk example multi_app_transfer/main.go").
		SetTransactionValidDuration(duration).
		Freeze()
	require.NoError(t, err)

	_ = transaction.GetTransactionID()
	nodeID := transaction.GetNodeAccountIDs()
	require.NotEmpty(t, nodeID)
	require.False(t, nodeID[0]._IsZero())

	var tx services.TransactionBody
	_ = protobuf.Unmarshal(transaction.signedTransactions._Get(0).(*services.SignedTransaction).BodyBytes, &tx)
	require.Equal(t, tx.TransactionID.String(), testTransactionID._ToProtobuf().String())
	require.Equal(t, tx.NodeAccountID.String(), node[0]._ToProtobuf().String())
	require.Equal(t, tx.Memo, "go sdk example multi_app_transfer/main.go")
	require.Equal(t, duration, _DurationFromProtobuf(tx.TransactionValidDuration))
	require.Equal(t, tx.GetCryptoTransfer().Transfers.AccountAmounts,
		[]*services.AccountAmount{
			{
				AccountID: &services.AccountID{Account: &services.AccountID_AccountNum{5}},
				Amount:    -100000000,
			},
			{
				AccountID: &services.AccountID{Account: &services.AccountID_AccountNum{4}},
				Amount:    100000000,
			},
		})

	txBytes, err := transaction.ToBytes()
	require.NoError(t, err)

	newTransaction, err := TransactionFromBytes(txBytes)

	_ = protobuf.Unmarshal(newTransaction.(TransferTransaction).signedTransactions._Get(0).(*services.SignedTransaction).BodyBytes, &tx)
	require.Equal(t, tx.TransactionID.String(), testTransactionID._ToProtobuf().String())
	require.Equal(t, tx.NodeAccountID.String(), node[0]._ToProtobuf().String())
	require.Equal(t, tx.Memo, "go sdk example multi_app_transfer/main.go")
	require.Equal(t, duration, _DurationFromProtobuf(tx.TransactionValidDuration))
	require.Equal(t, tx.GetCryptoTransfer().Transfers.AccountAmounts,
		[]*services.AccountAmount{
			{
				AccountID: &services.AccountID{Account: &services.AccountID_AccountNum{5}},
				Amount:    -100000000,
			},
			{
				AccountID: &services.AccountID{Account: &services.AccountID_AccountNum{4}},
				Amount:    100000000,
			},
		})
}

func TestUnitTransactionToFromBytesWithClient(t *testing.T) {
	t.Parallel()

	duration := time.Second * 10
	operatorID := AccountID{Account: 5}
	recepientID := AccountID{Account: 4}
	client, err := _NewMockClient()
	client.SetLedgerID(*NewLedgerIDTestnet())
	require.NoError(t, err)
	privateKey, err := PrivateKeyFromString(mockPrivateKey)
	client.SetOperator(AccountID{Account: 2}, privateKey)

	transaction, err := NewTransferTransaction().
		AddHbarTransfer(operatorID, NewHbar(-1)).
		AddHbarTransfer(recepientID, NewHbar(1)).
		SetTransactionMemo("go sdk example multi_app_transfer/main.go").
		SetTransactionValidDuration(duration).
		FreezeWith(client)
	require.NoError(t, err)

	var tx services.TransactionBody
	_ = protobuf.Unmarshal(transaction.signedTransactions._Get(0).(*services.SignedTransaction).BodyBytes, &tx)
	require.NotNil(t, tx.TransactionID, tx.NodeAccountID)
	require.Equal(t, tx.Memo, "go sdk example multi_app_transfer/main.go")
	require.Equal(t, duration, _DurationFromProtobuf(tx.TransactionValidDuration))
	require.Equal(t, tx.GetCryptoTransfer().Transfers.AccountAmounts,
		[]*services.AccountAmount{
			{
				AccountID: &services.AccountID{Account: &services.AccountID_AccountNum{5}},
				Amount:    -100000000,
			},
			{
				AccountID: &services.AccountID{Account: &services.AccountID_AccountNum{4}},
				Amount:    100000000,
			},
		})

	initialTxID := tx.TransactionID
	initialNode := tx.NodeAccountID

	txBytes, err := transaction.ToBytes()
	require.NoError(t, err)

	newTransaction, err := TransactionFromBytes(txBytes)

	_ = protobuf.Unmarshal(newTransaction.(TransferTransaction).signedTransactions._Get(0).(*services.SignedTransaction).BodyBytes, &tx)
	require.NotNil(t, tx.TransactionID, tx.NodeAccountID)
	require.Equal(t, tx.TransactionID.String(), initialTxID.String())
	require.Equal(t, tx.NodeAccountID.String(), initialNode.String())
	require.Equal(t, tx.Memo, "go sdk example multi_app_transfer/main.go")
	require.Equal(t, duration, _DurationFromProtobuf(tx.TransactionValidDuration))
	require.Equal(t, tx.GetCryptoTransfer().Transfers.AccountAmounts,
		[]*services.AccountAmount{
			{
				AccountID: &services.AccountID{Account: &services.AccountID_AccountNum{5}},
				Amount:    -100000000,
			},
			{
				AccountID: &services.AccountID{Account: &services.AccountID_AccountNum{4}},
				Amount:    100000000,
			},
		})
}

func TestUnitQueryRegression(t *testing.T) {
	t.Parallel()

	accountID := AccountID{Account: 5}
	node := []AccountID{{Account: 3}}
	client, err := _NewMockClient()
	client.SetLedgerID(*NewLedgerIDTestnet())
	require.NoError(t, err)
	privateKey, err := PrivateKeyFromString(mockPrivateKey)
	client.SetOperator(AccountID{Account: 2}, privateKey)

	query := NewAccountInfoQuery().
		SetAccountID(accountID).
		SetNodeAccountIDs(node).
		SetPaymentTransactionID(testTransactionID).
		SetMaxQueryPayment(NewHbar(1)).
		SetQueryPayment(HbarFromTinybar(25))

	body := query.buildQuery()
	_, err = query.generatePayments(client, HbarFromTinybar(20))
	require.NoError(t, err)

	var paymentTx services.TransactionBody
	_ = protobuf.Unmarshal(query.Query.paymentTransactions[0].BodyBytes, &paymentTx)

	require.Equal(t, body.GetCryptoGetInfo().AccountID.String(), accountID._ToProtobuf().String())
	require.Equal(t, paymentTx.NodeAccountID.String(), node[0]._ToProtobuf().String())
	require.Equal(t, paymentTx.TransactionFee, uint64(NewHbar(1).tinybar))
	require.Equal(t, paymentTx.TransactionValidDuration, &services.Duration{Seconds: 120})
	require.Equal(t, paymentTx.Data, &services.TransactionBody_CryptoTransfer{
		CryptoTransfer: &services.CryptoTransferTransactionBody{
			Transfers: &services.TransferList{
				AccountAmounts: []*services.AccountAmount{
					{
						AccountID: node[0]._ToProtobuf(),
						Amount:    HbarFromTinybar(20).AsTinybar(),
					},
					{
						AccountID: client.GetOperatorAccountID()._ToProtobuf(),
						Amount:    -HbarFromTinybar(20).AsTinybar(),
					},
				},
			},
		},
	})
}
func TestUnitTransactionInitFeeMaxTransactionWithouthSettingFee(t *testing.T) {
	t.Parallel()

	//Default Max Fee for TransferTransaction
	fee := NewHbar(1)
	client, err := _NewMockClient()
	client.SetLedgerID(*NewLedgerIDTestnet())
	require.NoError(t, err)
	transaction, err := NewTransferTransaction().
		AddHbarTransfer(AccountID{Account: 2}, HbarFromTinybar(-100)).
		AddHbarTransfer(AccountID{Account: 3}, HbarFromTinybar(100)).
		FreezeWith(client)
	require.NoError(t, err)
	require.Equal(t, uint64(fee.AsTinybar()), transaction.transactionFee)
}

func TestUnitTransactionInitFeeMaxTransactionFeeSetExplicitly(t *testing.T) {
	t.Parallel()

	clientMaxFee := NewHbar(14)
	explicitMaxFee := NewHbar(15)
	client, err := _NewMockClient()
	client.SetLedgerID(*NewLedgerIDTestnet())
	client.SetDefaultMaxTransactionFee(clientMaxFee)
	require.NoError(t, err)
	transaction, err := NewTransferTransaction().
		AddHbarTransfer(AccountID{Account: 2}, HbarFromTinybar(-100)).
		AddHbarTransfer(AccountID{Account: 3}, HbarFromTinybar(100)).
		SetMaxTransactionFee(explicitMaxFee).
		FreezeWith(client)
	require.NoError(t, err)
	require.Equal(t, uint64(explicitMaxFee.AsTinybar()), transaction.transactionFee)
}

func TestUnitTransactionInitFeeMaxTransactionFromClientDefault(t *testing.T) {
	t.Parallel()

	fee := NewHbar(14)
	client, err := _NewMockClient()
	client.SetLedgerID(*NewLedgerIDTestnet())
	client.SetDefaultMaxTransactionFee(fee)
	require.NoError(t, err)
	transaction, err := NewTransferTransaction().
		AddHbarTransfer(AccountID{Account: 2}, HbarFromTinybar(-100)).
		AddHbarTransfer(AccountID{Account: 3}, HbarFromTinybar(100)).
		FreezeWith(client)
	require.NoError(t, err)
	require.Equal(t, uint64(fee.AsTinybar()), transaction.transactionFee)
}

func TestUnitTransactionSignSwitchCases(t *testing.T) {
	t.Parallel()

	newKey, client, nodeAccountId := signSwitchCaseaSetup(t)

	txs := []TransactionInterface{
		NewAccountCreateTransaction(),
		NewAccountDeleteTransaction(),
		NewAccountUpdateTransaction(),
		NewAccountAllowanceApproveTransaction(),
		NewAccountAllowanceDeleteTransaction(),
		NewFileCreateTransaction(),
		NewFileDeleteTransaction(),
		NewFileUpdateTransaction(),
		NewLiveHashAddTransaction(),
		NewLiveHashDeleteTransaction(),
		NewTokenAssociateTransaction(),
		NewTokenBurnTransaction(),
		NewTokenCreateTransaction(),
		NewTokenDeleteTransaction(),
		NewTokenDissociateTransaction(),
		NewTokenFeeScheduleUpdateTransaction(),
		NewTokenFreezeTransaction(),
		NewTokenGrantKycTransaction(),
		NewTokenMintTransaction(),
		NewTokenRevokeKycTransaction(),
		NewTokenUnfreezeTransaction(),
		NewTokenUpdateTransaction(),
		NewTokenWipeTransaction(),
		NewTopicCreateTransaction(),
		NewTopicDeleteTransaction(),
		NewTopicUpdateTransaction(),
		NewTransferTransaction(),
	}

	for _, tx := range txs {

		txVal, signature, transferTxBytes := signSwitchCaseaHelper(t, tx, newKey, client)

		signTests := signTestsForTransaction(txVal, newKey, signature, client)

		for _, tt := range signTests {
			t.Run(tt.name, func(t *testing.T) {
				transactionInterface, err := TransactionFromBytes(transferTxBytes)
				require.NoError(t, err)

				tx, err := tt.sign(transactionInterface, newKey)
				assert.NoError(t, err)
				assert.NotEmpty(t, tx)

				signs, err := TransactionGetSignatures(transactionInterface)
				assert.NoError(t, err)

				// verify with range because var signs = map[AccountID]map[*PublicKey][]byte, where *PublicKey is unknown memory address
				for key := range signs[nodeAccountId] {
					assert.Equal(t, signs[nodeAccountId][key], signature)
				}
			})
		}
	}
}

func TestUnitTransactionSignSwitchCasesPointers(t *testing.T) {
	t.Parallel()

	newKey, client, nodeAccountId := signSwitchCaseaSetup(t)

	txs := []TransactionInterface{
		NewAccountCreateTransaction(),
		NewAccountDeleteTransaction(),
		NewAccountUpdateTransaction(),
		NewAccountAllowanceApproveTransaction(),
		NewAccountAllowanceDeleteTransaction(),
		NewFileCreateTransaction(),
		NewFileDeleteTransaction(),
		NewFileUpdateTransaction(),
		NewLiveHashAddTransaction(),
		NewLiveHashDeleteTransaction(),
		NewTokenAssociateTransaction(),
		NewTokenBurnTransaction(),
		NewTokenCreateTransaction(),
		NewTokenDeleteTransaction(),
		NewTokenDissociateTransaction(),
		NewTokenFeeScheduleUpdateTransaction(),
		NewTokenFreezeTransaction(),
		NewTokenGrantKycTransaction(),
		NewTokenMintTransaction(),
		NewTokenRevokeKycTransaction(),
		NewTokenUnfreezeTransaction(),
		NewTokenUpdateTransaction(),
		NewTokenWipeTransaction(),
		NewTopicCreateTransaction(),
		NewTopicDeleteTransaction(),
		NewTopicUpdateTransaction(),
		NewTransferTransaction(),
	}

	for _, tx := range txs {

		txVal, signature, transferTxBytes := signSwitchCaseaHelper(t, tx, newKey, client)
		signTests := signTestsForTransaction(txVal, newKey, signature, client)

		for _, tt := range signTests {
			t.Run(tt.name, func(t *testing.T) {
				transactionInterface, err := TransactionFromBytes(transferTxBytes)
				require.NoError(t, err)

				signs, err := TransactionGetSignatures(transactionInterface)
				assert.NoError(t, err)

				// verify with range because var signs = map[AccountID]map[*PublicKey][]byte, where *PublicKey is unknown memory address
				for key := range signs[nodeAccountId] {
					assert.Equal(t, signs[nodeAccountId][key], signature)
				}
			})
		}
	}
}

func TestUnitTransactionAttributes(t *testing.T) {
	t.Parallel()

	client, err := _NewMockClient()
	require.NoError(t, err)
	client.SetLedgerID(*NewLedgerIDTestnet())
	nodeAccountIds := client.network._GetNodeAccountIDsForExecute()

	txs := []TransactionInterface{
		NewAccountCreateTransaction(),
		NewAccountDeleteTransaction(),
		NewAccountUpdateTransaction(),
		NewAccountAllowanceApproveTransaction(),
		NewAccountAllowanceDeleteTransaction(),
		NewContractCreateTransaction(),
		NewContractDeleteTransaction(),
		NewContractExecuteTransaction(),
		NewContractUpdateTransaction(),
		NewFileAppendTransaction(),
		NewFileCreateTransaction(),
		NewFileDeleteTransaction(),
		NewFileUpdateTransaction(),
		NewLiveHashAddTransaction(),
		NewLiveHashDeleteTransaction(),
		NewScheduleCreateTransaction(),
		NewScheduleDeleteTransaction(),
		NewScheduleSignTransaction(),
		NewSystemDeleteTransaction(),
		NewSystemUndeleteTransaction(),
		NewTokenAssociateTransaction(),
		NewTokenBurnTransaction(),
		NewTokenCreateTransaction(),
		NewTokenDeleteTransaction(),
		NewTokenDissociateTransaction(),
		NewTokenFeeScheduleUpdateTransaction(),
		NewTokenFreezeTransaction(),
		NewTokenGrantKycTransaction(),
		NewTokenMintTransaction(),
		NewTokenRevokeKycTransaction(),
		NewTokenUnfreezeTransaction(),
		NewTokenUpdateTransaction(),
		NewTokenWipeTransaction(),
		NewTopicCreateTransaction(),
		NewTopicDeleteTransaction(),
		NewTopicUpdateTransaction(),
		NewTransferTransaction(),
	}

	for _, tx := range txs {
		txName := reflect.TypeOf(tx).Elem().Name()

		tests := createTransactionTests(txName, nodeAccountIds)

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				txSet, err := tt.set(tx)
				require.NoError(t, err)

				txGet, err := tt.get(txSet)
				require.NoError(t, err)

				tt.assert(t, txGet)
			})
		}
	}
}

func TestUnitTransactionAttributesDereferanced(t *testing.T) {
	t.Parallel()

	client, err := _NewMockClient()
	require.NoError(t, err)
	client.SetLedgerID(*NewLedgerIDTestnet())
	nodeAccountIds := client.network._GetNodeAccountIDsForExecute()

	txs := []TransactionInterface{
		NewAccountCreateTransaction(),
		NewAccountDeleteTransaction(),
		NewAccountUpdateTransaction(),
		NewAccountAllowanceApproveTransaction(),
		NewAccountAllowanceDeleteTransaction(),
		NewContractCreateTransaction(),
		NewContractDeleteTransaction(),
		NewContractExecuteTransaction(),
		NewContractUpdateTransaction(),
		NewFileAppendTransaction(),
		NewFileCreateTransaction(),
		NewFileDeleteTransaction(),
		NewFileUpdateTransaction(),
		NewLiveHashAddTransaction(),
		NewLiveHashDeleteTransaction(),
		NewScheduleCreateTransaction(),
		NewScheduleDeleteTransaction(),
		NewScheduleSignTransaction(),
		NewSystemDeleteTransaction(),
		NewSystemUndeleteTransaction(),
		NewTokenAssociateTransaction(),
		NewTokenBurnTransaction(),
		NewTokenCreateTransaction(),
		NewTokenDeleteTransaction(),
		NewTokenDissociateTransaction(),
		NewTokenFeeScheduleUpdateTransaction(),
		NewTokenFreezeTransaction(),
		NewTokenGrantKycTransaction(),
		NewTokenMintTransaction(),
		NewTokenRevokeKycTransaction(),
		NewTokenUnfreezeTransaction(),
		NewTokenUpdateTransaction(),
		NewTokenWipeTransaction(),
		NewTopicCreateTransaction(),
		NewTopicDeleteTransaction(),
		NewTopicUpdateTransaction(),
		NewTransferTransaction(),
	}

	for _, tx := range txs {
		txName := reflect.TypeOf(tx).Elem().Name()

		tests := createTransactionTests(txName, nodeAccountIds)

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				txSet, err := tt.set(tx)
				require.NoError(t, err)

				txGet, err := tt.get(txSet)
				require.NoError(t, err)

				tt.assert(t, txGet)
			})
		}
	}
}

func TestUnitTransactionAttributesSerialization(t *testing.T) {
	t.Parallel()

	client, err := _NewMockClient()
	require.NoError(t, err)
	client.SetLedgerID(*NewLedgerIDTestnet())

	txs := []TransactionInterface{
		NewAccountCreateTransaction(),
		NewAccountDeleteTransaction(),
		NewAccountUpdateTransaction(),
		NewAccountAllowanceApproveTransaction(),
		NewAccountAllowanceDeleteTransaction(),
		NewContractCreateTransaction(),
		NewContractDeleteTransaction(),
		NewContractExecuteTransaction(),
		NewContractUpdateTransaction(),
		NewFileCreateTransaction(),
		NewFileDeleteTransaction(),
		NewFileUpdateTransaction(),
		NewLiveHashAddTransaction(),
		NewLiveHashDeleteTransaction(),
		NewScheduleCreateTransaction(),
		NewScheduleDeleteTransaction(),
		NewScheduleSignTransaction(),
		NewSystemDeleteTransaction(),
		NewSystemUndeleteTransaction(),
		NewTokenAssociateTransaction(),
		NewTokenBurnTransaction(),
		NewTokenCreateTransaction(),
		NewTokenDeleteTransaction(),
		NewTokenDissociateTransaction(),
		NewTokenFeeScheduleUpdateTransaction(),
		NewTokenFreezeTransaction(),
		NewTokenGrantKycTransaction(),
		NewTokenMintTransaction(),
		NewTokenRevokeKycTransaction(),
		NewTokenUnfreezeTransaction(),
		NewTokenUpdateTransaction(),
		NewTokenWipeTransaction(),
		NewTopicCreateTransaction(),
		NewTopicDeleteTransaction(),
		NewTopicUpdateTransaction(),
		NewTransferTransaction(),
	}

	for _, tx := range txs {
		txName := reflect.TypeOf(tx).Elem().Name()

		// Get the reflect.Value of the pointer to the Transaction
		txPtr := reflect.ValueOf(tx)
		txPtr.MethodByName("FreezeWith").Call([]reflect.Value{reflect.ValueOf(client)})

		tests := []struct {
			name string
			act  func(transactionInterface TransactionInterface)
		}{
			{
				name: "TransactionString/" + txName,
				act: func(transactionInterface TransactionInterface) {
					txString, err := TransactionString(transactionInterface)
					require.NoError(t, err)
					require.NotEmpty(t, txString)
				},
			},
			{
				name: "TransactionToBytes/" + txName,
				act: func(transactionInterface TransactionInterface) {
					txBytes, err := TransactionToBytes(transactionInterface)
					require.NoError(t, err)
					require.NotEmpty(t, txBytes)
				},
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				tt.act(tx)
				// txValue := reflect.ValueOf(tx).Elem().Interface()
				// tt.act(txValue)
			})
		}
	}
}

func signSwitchCaseaSetup(t *testing.T) (PrivateKey, *Client, AccountID) {
	newKey, err := GeneratePrivateKey()
	require.NoError(t, err)

	client, err := _NewMockClient()
	require.NoError(t, err)
	client.SetLedgerID(*NewLedgerIDTestnet())

	nodeAccountIds := client.network._GetNodeAccountIDsForExecute()
	nodeAccountId := nodeAccountIds[0]

	return newKey, client, nodeAccountId
}

func signSwitchCaseaHelper(t *testing.T, tx TransactionInterface, newKey PrivateKey, client *Client) (txVal reflect.Value, signature []byte, transferTxBytes []byte) {
	// Get the reflect.Value of the pointer to the transaction
	txPtr := reflect.ValueOf(tx)
	txPtr.MethodByName("FreezeWith").Call([]reflect.Value{reflect.ValueOf(client)})

	// Get the reflect.Value of the transaction
	txVal = txPtr.Elem()

	// Get the transaction field by name
	// txField := txVal.FieldByName("Transaction")

	// Get the value of the Transaction field
	// txValue := txField.Interface().(Transaction[TransactionInterface])

	// refl_signature := reflect.ValueOf(newKey).MethodByName("SignTransaction").Call([]reflect.Value{reflect.ValueOf(&txValue)})
	signature, err := newKey.SignTransaction(tx)
	assert.NoError(t, err)

	transferTxBytes, err = TransactionToBytes(tx)
	assert.NoError(t, err)
	assert.NotEmpty(t, transferTxBytes)

	return txVal, signature, transferTxBytes
}

func signTestsForTransaction(txVal reflect.Value, newKey PrivateKey, signature []byte, client *Client) []struct {
	name string
	sign func(transactionInterface TransactionInterface, key Key) (TransactionInterface, error)
} {
	return []struct {
		name string
		sign func(transactionInterface TransactionInterface, key Key) (TransactionInterface, error)
	}{
		{
			name: "TransactionSign/" + txVal.Type().Name(),
			sign: func(transactionInterface TransactionInterface, key Key) (TransactionInterface, error) {
				privateKey, ok := key.(PrivateKey)
				if !ok {
					panic("key is not a PrivateKey")
				}
				return TransactionSign(transactionInterface, privateKey)
			},
		},
		{
			name: "TransactionSignWith/" + txVal.Type().Name(),
			sign: func(transactionInterface TransactionInterface, key Key) (TransactionInterface, error) {
				return TransactionSignWth(transactionInterface, newKey.PublicKey(), newKey.Sign)
			},
		},
		{
			name: "TransactionSignWithOperator/" + txVal.Type().Name(),
			sign: func(transactionInterface TransactionInterface, key Key) (TransactionInterface, error) {
				return TransactionSignWithOperator(transactionInterface, client)
			},
		},
		{
			name: "TransactionAddSignature/" + txVal.Type().Name(),
			sign: func(transactionInterface TransactionInterface, key Key) (TransactionInterface, error) {
				return TransactionAddSignature(transactionInterface, newKey.PublicKey(), signature)
			},
		},
	}
}

type transactionTest struct {
	name   string
	set    func(transactionInterface TransactionInterface) (TransactionInterface, error)
	get    func(transactionInterface TransactionInterface) (interface{}, error)
	assert func(t *testing.T, actual interface{})
}

func createTransactionTests(txName string, nodeAccountIds []AccountID) []transactionTest {
	return []transactionTest{
		{
			name: "TransactionTransactionID/" + txName,
			set: func(transactionInterface TransactionInterface) (TransactionInterface, error) {
				transactionID := TransactionID{AccountID: &AccountID{Account: 9999}, ValidStart: &time.Time{}, scheduled: false, Nonce: nil}
				return TransactionSetTransactionID(transactionInterface, transactionID)
			},
			get: func(transactionInterface TransactionInterface) (interface{}, error) {
				return TransactionGetTransactionID(transactionInterface)
			},
			assert: func(t *testing.T, actual interface{}) {
				transactionID := TransactionID{AccountID: &AccountID{Account: 9999}, ValidStart: &time.Time{}, scheduled: false, Nonce: nil}
				A := actual.(TransactionID)

				require.Equal(t, transactionID.AccountID, A.AccountID)
			},
		},
		{
			name: "TransactionTransactionMemo/" + txName,
			set: func(transactionInterface TransactionInterface) (TransactionInterface, error) {
				return TransactionSetTransactionMemo(transactionInterface, "test memo")
			},
			get: func(transactionInterface TransactionInterface) (interface{}, error) {
				return TransactionGetTransactionMemo(transactionInterface)
			},
			assert: func(t *testing.T, actual interface{}) {
				require.Equal(t, "test memo", actual)
			},
		},
		{
			name: "TransactionMaxTransactionFee/" + txName,
			set: func(transactionInterface TransactionInterface) (TransactionInterface, error) {
				return TransactionSetMaxTransactionFee(transactionInterface, NewHbar(1))
			},
			get: func(transactionInterface TransactionInterface) (interface{}, error) {
				return TransactionGetMaxTransactionFee(transactionInterface)
			},
			assert: func(t *testing.T, actual interface{}) {
				require.Equal(t, NewHbar(1), actual)
			},
		},
		{
			name: "TransactionTransactionValidDuration/" + txName,
			set: func(transactionInterface TransactionInterface) (TransactionInterface, error) {
				return TransactionSetTransactionValidDuration(transactionInterface, time.Second*10)
			},
			get: func(transactionInterface TransactionInterface) (interface{}, error) {
				return TransactionGetTransactionValidDuration(transactionInterface)
			},
			assert: func(t *testing.T, actual interface{}) {
				require.Equal(t, time.Second*10, actual)
			},
		},
		{
			name: "TransactionNodeAccountIDs/" + txName,
			set: func(transactionInterface TransactionInterface) (TransactionInterface, error) {
				return TransactionSetNodeAccountIDs(transactionInterface, nodeAccountIds)
			},
			get: func(transactionInterface TransactionInterface) (interface{}, error) {
				return TransactionGetNodeAccountIDs(transactionInterface)
			},
			assert: func(t *testing.T, actual interface{}) {
				require.Equal(t, nodeAccountIds, actual)
			},
		},
		{
			name: "TransactionMinBackoff/" + txName,
			set: func(transactionInterface TransactionInterface) (TransactionInterface, error) {
				tx, _ := TransactionSetMaxBackoff(transactionInterface, time.Second*200)
				return TransactionSetMinBackoff(tx, time.Second*10)
			},
			get: func(transactionInterface TransactionInterface) (interface{}, error) {
				return TransactionGetMinBackoff(transactionInterface)
			},
			assert: func(t *testing.T, actual interface{}) {
				require.Equal(t, time.Second*10, actual)
			},
		},
		{
			name: "TransactionMaxBackoff/" + txName,
			set: func(transactionInterface TransactionInterface) (TransactionInterface, error) {
				return TransactionSetMaxBackoff(transactionInterface, time.Second*200)
			},
			get: func(transactionInterface TransactionInterface) (interface{}, error) {
				return TransactionGetMaxBackoff(transactionInterface)
			},
			assert: func(t *testing.T, actual interface{}) {
				require.Equal(t, time.Second*200, actual)
			},
		},
	}
}

// TransactionGetTransactionHash //needs to be tested in e2e tests
// TransactionGetTransactionHashPerNode //needs to be tested in e2e tests
// TransactionExecute //needs to be tested in e2e tests
