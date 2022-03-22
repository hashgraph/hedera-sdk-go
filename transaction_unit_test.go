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
