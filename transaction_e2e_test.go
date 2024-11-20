//go:build all || e2e
// +build all e2e

package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"testing"

	"github.com/hiero-ledger/hiero-sdk-go/v2/proto/sdk"
	"github.com/hiero-ledger/hiero-sdk-go/v2/proto/services"
	"github.com/stretchr/testify/assert"

	"github.com/stretchr/testify/require"
	protobuf "google.golang.org/protobuf/proto"
)

func TestIntegrationTransactionAddSignature(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)

	newKey, err := PrivateKeyGenerateEd25519()
	require.NoError(t, err)

	resp, err := NewAccountCreateTransaction().
		SetKey(newKey.PublicKey()).
		SetNodeAccountIDs(env.NodeAccountIDs).
		Execute(env.Client)
	require.NoError(t, err)

	receipt, err := resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	tx, err := NewAccountDeleteTransaction().
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		SetAccountID(*receipt.AccountID).
		SetTransferAccountID(env.Client.GetOperatorAccountID()).
		FreezeWith(env.Client)
	require.NoError(t, err)

	updateBytes, err := tx.ToBytes()
	require.NoError(t, err)

	sig1, err := newKey.SignTransaction(tx)
	require.NoError(t, err)

	tx2, err := TransactionFromBytes(updateBytes)
	require.NoError(t, err)

	if newTx, ok := tx2.(AccountDeleteTransaction); ok {
		resp, err = newTx.AddSignature(newKey.PublicKey(), sig1).Execute(env.Client)
		require.NoError(t, err)
	}

	_, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

}

func TestIntegrationTransactionSignTransaction(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)

	newKey, err := PrivateKeyGenerateEd25519()
	require.NoError(t, err)

	resp, err := NewAccountCreateTransaction().
		SetKey(newKey.PublicKey()).
		SetNodeAccountIDs(env.NodeAccountIDs).
		Execute(env.Client)
	require.NoError(t, err)

	receipt, err := resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

	tx, err := NewAccountDeleteTransaction().
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		SetAccountID(*receipt.AccountID).
		SetTransferAccountID(env.Client.GetOperatorAccountID()).
		FreezeWith(env.Client)
	require.NoError(t, err)

	_, err = newKey.SignTransaction(tx)
	require.NoError(t, err)

	resp, err = tx.Execute(env.Client)
	require.NoError(t, err)

	_, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
	require.NoError(t, err)

}

func TestIntegrationTransactionGetHash(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)

	newKey, err := PrivateKeyGenerateEd25519()
	require.NoError(t, err)

	tx, err := NewAccountCreateTransaction().
		SetKey(newKey.PublicKey()).
		SetNodeAccountIDs(env.NodeAccountIDs).
		FreezeWith(env.Client)
	require.NoError(t, err)

	tx, err = tx.SignWithOperator(env.Client)
	require.NoError(t, err)

	hash, err := tx.GetTransactionHash()
	require.NoError(t, err)

	resp, err := tx.Execute(env.Client)
	require.NoError(t, err)

	record, err := resp.GetRecord(env.Client)
	require.NoError(t, err)

	assert.Equal(t, hash, record.TransactionHash)

}

func DisabledTestTransactionFromBytes(t *testing.T) { // nolint
	id := TransactionIDGenerate(AccountID{0, 0, 542348, nil, nil, nil})

	TransactionBody := services.TransactionBody{
		TransactionID: &services.TransactionID{
			AccountID: &services.AccountID{
				Account: &services.AccountID_AccountNum{AccountNum: 542348},
			},
			TransactionValidStart: &services.Timestamp{
				Seconds: id.ValidStart.Unix(),
				Nanos:   int32(id.ValidStart.Nanosecond()),
			},
		},
		NodeAccountID: &services.AccountID{
			Account: &services.AccountID_AccountNum{AccountNum: 3},
		},
		TransactionFee: 200_000_000,
		TransactionValidDuration: &services.Duration{
			Seconds: 120,
		},
		GenerateRecord: false,
		Memo:           "",
		Data: &services.TransactionBody_CryptoTransfer{
			CryptoTransfer: &services.CryptoTransferTransactionBody{
				Transfers: &services.TransferList{
					AccountAmounts: []*services.AccountAmount{
						{
							AccountID: &services.AccountID{
								Account: &services.AccountID_AccountNum{AccountNum: 47439},
							},
							Amount: 10,
						},
						{
							AccountID: &services.AccountID{
								Account: &services.AccountID_AccountNum{AccountNum: 542348},
							},
							Amount: -10,
						},
					},
				},
			},
		},
	}

	BodyBytes, err := protobuf.Marshal(&TransactionBody)
	require.NoError(t, err)

	key1, _ := PrivateKeyFromString("302e020100300506032b6570042204203e7fda6dde63c3cdb3cb5ecf5264324c5faad7c9847b6db093c088838b35a110")
	key2, _ := PrivateKeyFromString("302e020100300506032b65700422042032d3d5a32e9d06776976b39c09a31fbda4a4a0208223da761c26a2ae560c1755")
	key3, _ := PrivateKeyFromString("302e020100300506032b657004220420195a919056d1d698f632c228dbf248bbbc3955adf8a80347032076832b8299f9")
	key4, _ := PrivateKeyFromString("302e020100300506032b657004220420b9962f17f94ffce73a23649718a11638cac4b47095a7a6520e88c7563865be62")
	key5, _ := PrivateKeyFromString("302e020100300506032b657004220420fef68591819080cd9d48b0cbaa10f65f919752abb50ffb3e7411ac66ab22692e")

	publicKey1 := key1.PublicKey()
	publicKey2 := key2.PublicKey()
	publicKey3 := key3.PublicKey()
	publicKey4 := key4.PublicKey()
	publicKey5 := key5.PublicKey()

	signature1 := key1.Sign(BodyBytes)
	signature2 := key2.Sign(BodyBytes)
	signature3 := key3.Sign(BodyBytes)
	signature4 := key4.Sign(BodyBytes)
	signature5 := key5.Sign(BodyBytes)

	signed := services.SignedTransaction{
		BodyBytes: BodyBytes,
		SigMap: &services.SignatureMap{
			SigPair: []*services.SignaturePair{
				{
					PubKeyPrefix: key1.PublicKey().Bytes(),
					Signature: &services.SignaturePair_Ed25519{
						Ed25519: signature1,
					},
				},
				{
					PubKeyPrefix: key2.PublicKey().Bytes(),
					Signature: &services.SignaturePair_Ed25519{
						Ed25519: signature2,
					},
				},
				{
					PubKeyPrefix: key3.PublicKey().Bytes(),
					Signature: &services.SignaturePair_Ed25519{
						Ed25519: signature3,
					},
				},
				{
					PubKeyPrefix: key4.PublicKey().Bytes(),
					Signature: &services.SignaturePair_Ed25519{
						Ed25519: signature4,
					},
				},
				{
					PubKeyPrefix: key5.PublicKey().Bytes(),
					Signature: &services.SignaturePair_Ed25519{
						Ed25519: signature5,
					},
				},
			},
		},
	}

	bytes, err := protobuf.Marshal(&signed)
	require.NoError(t, err)

	bytes, err = protobuf.Marshal(&sdk.TransactionList{
		TransactionList: []*services.Transaction{{
			SignedTransactionBytes: bytes,
		}},
	})
	require.NoError(t, err)

	transaction, err := TransactionFromBytes(bytes)
	require.NoError(t, err)

	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)

	switch tx := transaction.(type) {
	case *TransferTransaction:
		assert.Equal(t, tx.GetHbarTransfers()[AccountID{0, 0, 542348, nil, nil, nil}].AsTinybar(), int64(-10))
		assert.Equal(t, tx.GetHbarTransfers()[AccountID{0, 0, 47439, nil, nil, nil}].AsTinybar(), int64(10))

		signatures, err := tx.GetSignatures()
		require.NoError(t, err)
		assert.Contains(t, signatures[AccountID{0, 0, 3, nil, nil, nil}], &publicKey1)
		assert.Contains(t, signatures[AccountID{0, 0, 3, nil, nil, nil}], &publicKey2)
		assert.Contains(t, signatures[AccountID{0, 0, 3, nil, nil, nil}], &publicKey3)
		assert.Contains(t, signatures[AccountID{0, 0, 3, nil, nil, nil}], &publicKey4)
		assert.Contains(t, signatures[AccountID{0, 0, 3, nil, nil, nil}], &publicKey5)

		assert.Equal(t, len(tx.GetNodeAccountIDs()), 1)
		assert.True(t, tx.GetNodeAccountIDs()[0]._Equals(AccountID{0, 0, 3, nil, nil, nil}))

		resp, err := tx.Execute(env.Client)
		require.NoError(t, err)

		_, err = resp.SetValidateStatus(true).GetReceipt(env.Client)
		require.NoError(t, err)
	default:
		panic("Transaction was not a crypto transfer?")
	}
}

func TestIntegrationTransactionFailsWhenSigningWithoutFreezing(t *testing.T) {
	t.Parallel()
	env := NewIntegrationTestEnv(t)
	defer CloseIntegrationTestEnv(env, nil)

	newKey, err := PrivateKeyGenerateEd25519()
	require.NoError(t, err)

	tx := NewAccountCreateTransaction().
		SetKey(newKey.PublicKey()).
		SetNodeAccountIDs(env.NodeAccountIDs)

	_, err = tx.Sign(newKey).Execute(env.Client)
	require.ErrorContains(t, err, "transaction is not frozen")

}
