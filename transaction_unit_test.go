//go:build all || unit
// +build all unit

package hedera

import (
	"fmt"
	"testing"

	"github.com/hashgraph/hedera-protobufs-go/sdk"
	"github.com/hashgraph/hedera-protobufs-go/services"
	protobuf "google.golang.org/protobuf/proto"

	"github.com/stretchr/testify/assert"

	"github.com/stretchr/testify/require"
)

func TestUnitTransactionSerializationDeserialization(t *testing.T) {
	transaction, err := _NewMockTransaction()
	require.NoError(t, err)

	_, err = transaction.Freeze()
	require.NoError(t, err)

	_, err = transaction.GetSignatures()
	require.NoError(t, err)

	_, err = transaction.GetTransactionHash()
	require.NoError(t, err)

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

	assert.Equal(t, transaction.String(), deserializedTXTyped.String())
}

func TestUnitTransactionValidateBodiesEqual(t *testing.T) {
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

	var deserializedTXTyped *AccountCreateTransaction
	switch tx := deserializedTX.(type) {
	case AccountCreateTransaction:
		deserializedTXTyped = &tx
	default:
		panic("Transaction was not AccountCreateTransaction")
	}

	assert.Equal(t, uint64(transaction.TransactionID.AccountID.GetAccountNum()), deserializedTXTyped.GetTransactionID().AccountID.Account)
}

func TestUnitTransactionValidateBodiesNotEqual(t *testing.T) {
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
	operatorID := AccountID{Account: 5}
	recepientID := AccountID{Account: 4}
	node := []AccountID{{Account: 3}}
	transaction, err := NewTransferTransaction().
		SetTransactionID(testTransactionID).
		SetNodeAccountIDs(node).
		AddHbarTransfer(operatorID, NewHbar(-1)).
		AddHbarTransfer(recepientID, NewHbar(1)).
		SetTransactionMemo("go sdk example multi_app_transfer/main.go").
		Freeze()
	require.NoError(t, err)

	var tx services.TransactionBody
	_ = protobuf.Unmarshal(transaction.signedTransactions._Get(0).(*services.SignedTransaction).BodyBytes, &tx)
	require.Equal(t, tx.TransactionID.String(), testTransactionID._ToProtobuf().String())
	require.Equal(t, tx.NodeAccountID.String(), node[0]._ToProtobuf().String())
	require.Equal(t, tx.Memo, "go sdk example multi_app_transfer/main.go")
	require.Equal(t, tx.Data, &services.TransactionBody_CryptoTransfer{
		CryptoTransfer: &services.CryptoTransferTransactionBody{
			Transfers: &services.TransferList{
				AccountAmounts: []*services.AccountAmount{
					{
						AccountID: &services.AccountID{Account: &services.AccountID_AccountNum{4}},
						Amount:    100000000,
					},
					{
						AccountID: &services.AccountID{Account: &services.AccountID_AccountNum{5}},
						Amount:    -100000000,
					},
				},
			},
		},
	})

	txBytes, err := transaction.ToBytes()
	require.NoError(t, err)

	newTransaction, err := TransactionFromBytes(txBytes)

	_ = protobuf.Unmarshal(newTransaction.(TransferTransaction).signedTransactions._Get(0).(*services.SignedTransaction).BodyBytes, &tx)
	require.Equal(t, tx.TransactionID.String(), testTransactionID._ToProtobuf().String())
	require.Equal(t, tx.NodeAccountID.String(), node[0]._ToProtobuf().String())
	require.Equal(t, tx.Memo, "go sdk example multi_app_transfer/main.go")
	require.Equal(t, tx.Data, &services.TransactionBody_CryptoTransfer{
		CryptoTransfer: &services.CryptoTransferTransactionBody{
			Transfers: &services.TransferList{
				AccountAmounts: []*services.AccountAmount{
					{
						AccountID: &services.AccountID{Account: &services.AccountID_AccountNum{4}},
						Amount:    100000000,
					},
					{
						AccountID: &services.AccountID{Account: &services.AccountID_AccountNum{5}},
						Amount:    -100000000,
					},
				},
			},
		},
	})
}

func TestUnitTransactionToFromBytesWithClient(t *testing.T) {
	operatorID := AccountID{Account: 5}
	recepientID := AccountID{Account: 4}
	client := ClientForTestnet()
	privateKey, err := PrivateKeyFromString(mockPrivateKey)
	client.SetOperator(AccountID{Account: 2}, privateKey)

	transaction, err := NewTransferTransaction().
		AddHbarTransfer(operatorID, NewHbar(-1)).
		AddHbarTransfer(recepientID, NewHbar(1)).
		SetTransactionMemo("go sdk example multi_app_transfer/main.go").
		FreezeWith(client)
	require.NoError(t, err)

	var tx services.TransactionBody
	_ = protobuf.Unmarshal(transaction.signedTransactions._Get(0).(*services.SignedTransaction).BodyBytes, &tx)
	require.NotNil(t, tx.TransactionID, tx.NodeAccountID)
	require.Equal(t, tx.Memo, "go sdk example multi_app_transfer/main.go")
	require.Equal(t, tx.Data, &services.TransactionBody_CryptoTransfer{
		CryptoTransfer: &services.CryptoTransferTransactionBody{
			Transfers: &services.TransferList{
				AccountAmounts: []*services.AccountAmount{
					{
						AccountID: &services.AccountID{Account: &services.AccountID_AccountNum{4}},
						Amount:    100000000,
					},
					{
						AccountID: &services.AccountID{Account: &services.AccountID_AccountNum{5}},
						Amount:    -100000000,
					},
				},
			},
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
	require.Equal(t, tx.Data, &services.TransactionBody_CryptoTransfer{
		CryptoTransfer: &services.CryptoTransferTransactionBody{
			Transfers: &services.TransferList{
				AccountAmounts: []*services.AccountAmount{
					{
						AccountID: &services.AccountID{Account: &services.AccountID_AccountNum{4}},
						Amount:    100000000,
					},
					{
						AccountID: &services.AccountID{Account: &services.AccountID_AccountNum{5}},
						Amount:    -100000000,
					},
				},
			},
		},
	})
}
