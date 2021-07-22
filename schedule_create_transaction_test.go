package hedera

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

//func TestScheduleCreateTransaction_Execute(t *testing.T) {
//	env := NewIntegrationTestEnv(t)
//
//	newKey, err := GeneratePrivateKey()
//	assert.NoError(t, err)
//
//	newBalance := NewHbar(1)
//
//	assert.Equal(t, HbarUnits.Hbar.numberOfTinybar(), newBalance.tinybar)
//
//	transactionID := TransactionIDGenerate(env.Client.GetOperatorAccountID())
//
//	tx := NewAccountCreateTransaction().
//		SetTransactionID(transactionID).
//		SetKey(newKey.PublicKey()).
//		SetNodeAccountIDs(env.NodeAccountIDs).
//		SetMaxTransactionFee(NewHbar(2)).
//		SetInitialBalance(newBalance)
//
//	assert.NoError(t, err)
//
//	scheduleTx, err := tx.Schedule()
//	assert.NoError(t, err)
//
//	resp, err := scheduleTx.
//		SetPayerAccountID(env.Client.GetOperatorAccountID()).
//		SetAdminKey(env.Client.GetOperatorPublicKey()).
//		Execute(env.Client)
//	assert.NoError(t, err)
//
//	receipt, err := resp.GetReceipt(env.Client)
//	assert.NoError(t, err)
//
//	info, err := NewScheduleInfoQuery().
//		SetScheduleID(*receipt.ScheduleID).
//		SetQueryPayment(NewHbar(2)).
//		Execute(env.Client)
//	assert.NoError(t, err)
//
//	infoTx, err := info.GetScheduledTransaction()
//	assert.NoError(t, err)
//	assert.NotNil(t, infoTx)
//
//	tx2, err := NewScheduleDeleteTransaction().
//		SetScheduleID(*receipt.ScheduleID).
//		FreezeWith(env.Client)
//	assert.NoError(t, err)
//
//	resp, err = tx2.
//		Sign(newKey).
//		Execute(env.Client)
//	assert.NoError(t, err)
//
//	_, err = resp.GetReceipt(env.Client)
//	assert.Error(t, err)
//	if err != nil {
//		assert.Equal(t, fmt.Sprintf("exceptional receipt status: SCHEDULE_ALREADY_EXECUTED"), err.Error())
//	}
//}
//
//func TestScheduleCreateTransaction_SetTransaction_Execute(t *testing.T) {
//	env := NewIntegrationTestEnv(t)
//
//	newKey, err := GeneratePrivateKey()
//	assert.NoError(t, err)
//
//	newBalance := NewHbar(1)
//
//	assert.Equal(t, HbarUnits.Hbar.numberOfTinybar(), newBalance.tinybar)
//
//	transactionID := TransactionIDGenerate(env.Client.GetOperatorAccountID())
//
//	tx := NewAccountCreateTransaction().
//		SetKey(newKey.PublicKey()).
//		SetTransactionID(transactionID).
//		SetNodeAccountIDs(env.NodeAccountIDs).
//		SetMaxTransactionFee(NewHbar(2)).
//		SetInitialBalance(newBalance)
//
//	scheduleTx, err := NewScheduleCreateTransaction().
//		SetPayerAccountID(env.Client.GetOperatorAccountID()).
//		SetAdminKey(env.Client.GetOperatorPublicKey()).
//		SetScheduledTransaction(tx)
//	assert.NoError(t, err)
//
//	resp, err := scheduleTx.Execute(env.Client)
//	assert.NoError(t, err)
//
//	receipt, err := resp.GetReceipt(env.Client)
//	assert.NoError(t, err)
//
//	_, err = NewScheduleInfoQuery().
//		SetScheduleID(*receipt.ScheduleID).
//		Execute(env.Client)
//	assert.NoError(t, err)
//
//	tx2, err := NewScheduleDeleteTransaction().
//		SetScheduleID(*receipt.ScheduleID).
//		FreezeWith(env.Client)
//	assert.NoError(t, err)
//
//	resp, err = tx2.
//		Sign(newKey).
//		Execute(env.Client)
//	assert.NoError(t, err)
//
//	_, err = resp.GetReceipt(env.Client)
//	assert.Error(t, err)
//	if err != nil {
//		assert.Equal(t, fmt.Sprintf("exceptional receipt status: SCHEDULE_ALREADY_EXECUTED"), err.Error())
//	}
//}
//
//func TestScheduleCreateTransaction_MultiSig_Execute(t *testing.T) {
//	env := NewIntegrationTestEnv(t)
//
//	keys := make([]PrivateKey, 3)
//	pubKeys := make([]PublicKey, 3)
//
//	for i := range keys {
//		newKey, err := GeneratePrivateKey()
//		assert.NoError(t, err)
//
//		keys[i] = newKey
//		pubKeys[i] = newKey.PublicKey()
//	}
//
//	keyList := NewKeyList().
//		AddAllPublicKeys(pubKeys)
//
//	createResponse, err := NewAccountCreateTransaction().
//		SetKey(keyList).
//		SetInitialBalance(NewHbar(10)).
//		SetNodeAccountIDs(env.NodeAccountIDs).
//		Execute(env.Client)
//	assert.NoError(t, err)
//
//	transactionReceipt, err := createResponse.GetReceipt(env.Client)
//	assert.NoError(t, err)
//
//	transactionID := TransactionIDGenerate(env.Client.GetOperatorAccountID())
//
//	newAccountID := *transactionReceipt.AccountID
//
//	transferTx := NewTransferTransaction().
//		SetTransactionID(transactionID).
//		AddHbarTransfer(newAccountID, HbarFrom(-1, HbarUnits.Hbar)).
//		AddHbarTransfer(env.Client.GetOperatorAccountID(), HbarFrom(1, HbarUnits.Hbar))
//
//	scheduled, err := transferTx.Schedule()
//	assert.NoError(t, err)
//
//	scheduleResponse, err := scheduled.Execute(env.Client)
//	assert.NoError(t, err)
//
//	scheduleReceipt, err := scheduleResponse.GetReceipt(env.Client)
//	assert.NoError(t, err)
//
//	scheduleID := *scheduleReceipt.ScheduleID
//
//	info, err := NewScheduleInfoQuery().
//		SetNodeAccountIDs([]AccountID{createResponse.NodeID}).
//		SetScheduleID(scheduleID).
//		Execute(env.Client)
//	assert.NoError(t, err)
//
//	transfer, err := info.GetScheduledTransaction()
//	assert.NoError(t, err)
//	assert.NotNil(t, transfer)
//
//	signTransaction, err := NewScheduleSignTransaction().
//		SetNodeAccountIDs([]AccountID{createResponse.NodeID}).
//		SetScheduleID(scheduleID).
//		FreezeWith(env.Client)
//	assert.NoError(t, err)
//
//	signTransaction.Sign(keys[0])
//	signTransaction.Sign(keys[1])
//	signTransaction.Sign(keys[2])
//
//	resp, err := signTransaction.Execute(env.Client)
//	assert.NoError(t, err)
//
//	_, err = resp.GetReceipt(env.Client)
//	assert.NoError(t, err)
//
//	info2, err := NewScheduleInfoQuery().
//		SetScheduleID(scheduleID).
//		SetNodeAccountIDs([]AccountID{createResponse.NodeID}).
//		Execute(env.Client)
//	assert.NoError(t, err)
//	assert.False(t, info2.ExecutedAt.IsZero())
//}
//
//func TestScheduleCreateTransaction_Delete_Execute(t *testing.T) {
//	env := NewIntegrationTestEnv(t)
//
//	key, err := GeneratePrivateKey()
//	key2, err := GeneratePrivateKey()
//	assert.NoError(t, err)
//
//	createResponse, err := NewAccountCreateTransaction().
//		SetKey(key).
//		SetInitialBalance(NewHbar(10)).
//		SetNodeAccountIDs(env.NodeAccountIDs).
//		Execute(env.Client)
//	assert.NoError(t, err)
//
//	transactionReceipt, err := createResponse.GetReceipt(env.Client)
//	assert.NoError(t, err)
//
//	transactionID := TransactionIDGenerate(env.Client.GetOperatorAccountID())
//
//	newAccountID := *transactionReceipt.AccountID
//
//	transferTx := NewTransferTransaction().
//		SetTransactionID(transactionID).
//		AddHbarTransfer(newAccountID, HbarFrom(-1, HbarUnits.Hbar)).
//		AddHbarTransfer(env.Client.GetOperatorAccountID(), HbarFrom(1, HbarUnits.Hbar))
//
//	scheduled, err := transferTx.Schedule()
//	assert.NoError(t, err)
//
//	fr, err := scheduled.SetAdminKey(key2).FreezeWith(env.Client)
//	assert.NoError(t, err)
//
//	scheduleResponse, err := fr.Sign(key2).Execute(env.Client)
//	assert.NoError(t, err)
//
//	scheduleReceipt, err := scheduleResponse.GetReceipt(env.Client)
//	assert.NoError(t, err)
//
//	scheduleID := *scheduleReceipt.ScheduleID
//
//	info, err := NewScheduleInfoQuery().
//		SetNodeAccountIDs([]AccountID{createResponse.NodeID}).
//		SetScheduleID(scheduleID).
//		Execute(env.Client)
//	assert.NoError(t, err)
//
//	transfer, err := info.GetScheduledTransaction()
//	assert.NoError(t, err)
//	assert.NotNil(t, transfer)
//	assert.Nil(t, info.ExecutedAt)
//	assert.Nil(t, info.DeletedAt)
//
//	tx2, err := NewScheduleDeleteTransaction().
//		SetScheduleID(scheduleID).
//		FreezeWith(env.Client)
//	assert.NoError(t, err)
//
//	resp, err := tx2.
//		Sign(key2).
//		Execute(env.Client)
//	assert.NoError(t, err)
//
//	_, err = resp.GetReceipt(env.Client)
//	assert.NoError(t, err)
//
//	info2, err := NewScheduleInfoQuery().
//		SetScheduleID(scheduleID).
//		SetNodeAccountIDs([]AccountID{createResponse.NodeID}).
//		Execute(env.Client)
//	assert.NoError(t, err)
//	assert.False(t, info2.DeletedAt.IsZero())
//}
//
//func TestScheduleCreateTransaction_CheckValidGetTransaction_Execute(t *testing.T) {
//	env := NewIntegrationTestEnv(t)
//
//	newKey, err := GeneratePrivateKey()
//	assert.NoError(t, err)
//
//	newBalance := NewHbar(1)
//
//	assert.Equal(t, HbarUnits.Hbar.numberOfTinybar(), newBalance.tinybar)
//
//	transactionID := TransactionIDGenerate(env.Client.GetOperatorAccountID())
//
//	tx := NewAccountCreateTransaction().
//		SetTransactionID(transactionID).
//		SetKey(newKey.PublicKey()).
//		SetNodeAccountIDs(env.NodeAccountIDs).
//		SetMaxTransactionFee(NewHbar(2)).
//		SetInitialBalance(newBalance)
//
//	assert.NoError(t, err)
//
//	scheduleTx, err := tx.Schedule()
//	assert.NoError(t, err)
//
//	resp, err := scheduleTx.
//		SetPayerAccountID(env.Client.GetOperatorAccountID()).
//		SetAdminKey(env.Client.GetOperatorPublicKey()).
//		Execute(env.Client)
//	assert.NoError(t, err)
//
//	receipt, err := resp.GetReceipt(env.Client)
//	assert.NoError(t, err)
//
//	info, err := NewScheduleInfoQuery().
//		SetScheduleID(*receipt.ScheduleID).
//		SetQueryPayment(NewHbar(2)).
//		Execute(env.Client)
//	assert.NoError(t, err)
//
//	infoTx, err := info.GetScheduledTransaction()
//	assert.NoError(t, err)
//
//	assert.NotNil(t, infoTx)
//
//	switch createTx := infoTx.(type) {
//	case *AccountCreateTransaction:
//		assert.Equal(t, createTx.pbBody.GetCryptoCreateAccount().InitialBalance, uint64(NewHbar(1).tinybar))
//	}
//
//	tx2, err := NewScheduleDeleteTransaction().
//		SetScheduleID(*receipt.ScheduleID).
//		FreezeWith(env.Client)
//	assert.NoError(t, err)
//
//	resp, err = tx2.
//		Sign(newKey).
//		Execute(env.Client)
//	assert.NoError(t, err)
//
//	_, err = resp.GetReceipt(env.Client)
//	assert.Error(t, err)
//	if err != nil {
//		assert.Equal(t, fmt.Sprintf("exceptional receipt status: SCHEDULE_ALREADY_EXECUTED"), err.Error())
//	}
//}
//
//func TestScheduleCreateTransaction_Duplicate_Execute(t *testing.T) {
//	env := NewIntegrationTestEnv(t)
//
//	newKey, err := GeneratePrivateKey()
//	assert.NoError(t, err)
//
//	newBalance := NewHbar(1)
//
//	assert.Equal(t, HbarUnits.Hbar.numberOfTinybar(), newBalance.tinybar)
//
//	transactionID := TransactionIDGenerate(env.Client.GetOperatorAccountID())
//
//	tx := NewAccountCreateTransaction().
//		SetTransactionID(transactionID).
//		SetNodeAccountIDs(env.NodeAccountIDs).
//		SetKey(newKey.PublicKey()).
//		SetMaxTransactionFee(NewHbar(2)).
//		SetInitialBalance(newBalance)
//
//	assert.NoError(t, err)
//
//	scheduleTx, err := tx.Schedule()
//	assert.NoError(t, err)
//
//	scheduleTx = scheduleTx.
//		SetPayerAccountID(env.Client.GetOperatorAccountID()).
//		SetAdminKey(env.Client.GetOperatorPublicKey()).
//		SetTransactionID(TransactionIDGenerate(env.Client.GetOperatorAccountID()))
//
//	resp, err := scheduleTx.Execute(env.Client)
//	assert.NoError(t, err)
//
//	_, err = resp.GetReceipt(env.Client)
//	assert.NoError(t, err)
//
//	_, err = scheduleTx.Execute(env.Client)
//	assert.Error(t, err)
//
//	scheduleTx, err = tx.Schedule()
//	assert.NoError(t, err)
//
//	scheduleTx = scheduleTx.
//		SetPayerAccountID(env.Client.GetOperatorAccountID()).
//		SetAdminKey(env.Client.GetOperatorPublicKey()).
//		SetTransactionID(TransactionIDGenerate(env.Client.GetOperatorAccountID()))
//
//	resp, err = scheduleTx.Execute(env.Client)
//	assert.NoError(t, err)
//
//	_, err = resp.GetReceipt(env.Client)
//	assert.NoError(t, err)
//}

func TestScheduleCreateTransaction_Transfer_Execute(t *testing.T) {
	env := NewIntegrationTestEnv(t)

	key, err := GeneratePrivateKey()
	assert.NoError(t, err)

	response, err := NewAccountCreateTransaction().
		SetKey(key).
		SetInitialBalance(NewHbar(2)).
		Execute(env.Client)
	assert.NoError(t, err)

	receipt, err := response.GetReceipt(env.Client)
	assert.NoError(t, err)

	accountID := *receipt.AccountID

	tx := NewTransferTransaction().
		AddHbarTransfer(accountID, NewHbar(1).Negated()).
		AddHbarTransfer(env.Client.GetOperatorAccountID(), NewHbar(1))

	scheduleTx, err := tx.Schedule()
	assert.NoError(t, err)

	scheduleTx = scheduleTx.
		SetNodeAccountIDs(env.NodeAccountIDs).
		SetPayerAccountID(env.Client.GetOperatorAccountID()).
		SetAdminKey(env.Client.GetOperatorPublicKey()).
		SetTransactionID(TransactionIDGenerate(env.Client.GetOperatorAccountID()))

	response, err = scheduleTx.Execute(env.Client)
	assert.NoError(t, err)

	receipt, err = response.GetReceipt(env.Client)
	assert.NoError(t, err)

	scheduleID := *receipt.ScheduleID

	scheduleSignTx, err := NewScheduleSignTransaction().
		SetNodeAccountIDs(env.NodeAccountIDs).
		SetScheduleID(scheduleID).
		FreezeWith(env.Client)
	assert.NoError(t, err)

	scheduleSignTx.Sign(key)

	response, err = scheduleSignTx.Execute(env.Client)
	assert.NoError(t, err)

	_, err = response.GetReceipt(env.Client)
	assert.NoError(t, err)

	info, err := NewScheduleInfoQuery().
		SetNodeAccountIDs(env.NodeAccountIDs).
		SetScheduleID(scheduleID).
		Execute(env.Client)
	assert.NoError(t, err)
	assert.NotNil(t, info.ExecutedAt)

	err = CloseIntegrationTestEnv(env, nil)
	assert.NoError(t, err)
}

func TestScheduledTokenNftTransferTransaction_Execute(t *testing.T) {
	env := NewIntegrationTestEnv(t)

	keys := make([]PrivateKey, 3)
	pubKeys := make([]PublicKey, 3)

	for i := range keys {
		newKey, err := GeneratePrivateKey()
		assert.NoError(t, err)

		keys[i] = newKey
		pubKeys[i] = newKey.PublicKey()
	}

	keyList := NewKeyList().
		AddAllPublicKeys(pubKeys)

	newBalance := NewHbar(2)

	assert.Equal(t, 2*HbarUnits.Hbar.numberOfTinybar(), newBalance.tinybar)

	resp, err := NewAccountCreateTransaction().
		SetNodeAccountIDs(env.NodeAccountIDs).
		SetKey(keyList).
		SetInitialBalance(newBalance).
		Execute(env.Client)
	assert.NoError(t, err)

	receipt, err := resp.GetReceipt(env.Client)
	assert.NoError(t, err)

	accountID := *receipt.AccountID

	resp, err = NewTokenCreateTransaction().
		SetNodeAccountIDs(env.NodeAccountIDs).
		SetTokenName("ffff").
		SetTokenSymbol("F").
		SetTokenType(TokenTypeNonFungibleUnique).
		SetSupplyType(TokenSupplyTypeFinite).
		SetMaxSupply(5).
		SetTreasuryAccountID(env.Client.GetOperatorAccountID()).
		SetAdminKey(env.Client.GetOperatorPublicKey()).
		SetFreezeKey(env.Client.GetOperatorPublicKey()).
		SetWipeKey(env.Client.GetOperatorPublicKey()).
		SetKycKey(env.Client.GetOperatorPublicKey()).
		SetSupplyKey(env.Client.GetOperatorPublicKey()).
		SetFreezeDefault(false).
		Execute(env.Client)
	assert.NoError(t, err)

	receipt, err = resp.GetReceipt(env.Client)
	assert.NoError(t, err)

	tokenID := *receipt.TokenID
	metaData := [][]byte{{50}, {50}}

	mint, err := NewTokenMintTransaction().
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		SetTokenID(tokenID).
		SetMetadatas(metaData).
		Execute(env.Client)
	assert.NoError(t, err)

	mintReceipt, err := mint.GetReceipt(env.Client)
	assert.NoError(t, err)

	transaction, err := NewTokenAssociateTransaction().
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		SetAccountID(accountID).
		SetTokenIDs(tokenID).
		FreezeWith(env.Client)
	assert.NoError(t, err)

	resp, err = transaction.
		Sign(keys[0]).
		Sign(keys[1]).
		Sign(keys[2]).
		Execute(env.Client)
	assert.NoError(t, err)

	_, err = resp.GetReceipt(env.Client)
	assert.NoError(t, err)

	resp, err = NewTokenGrantKycTransaction().
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		SetAccountID(accountID).
		SetTokenID(tokenID).
		Execute(env.Client)
	assert.NoError(t, err)

	_, err = resp.GetReceipt(env.Client)
	assert.NoError(t, err)

	tx := NewTransferTransaction().
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		AddNftTransfer(tokenID.Nft(mintReceipt.SerialNumbers[0]), env.OperatorID, accountID).
		AddNftTransfer(tokenID.Nft(mintReceipt.SerialNumbers[1]), env.OperatorID, accountID)

	scheduleTx, err := tx.Schedule()
	assert.NoError(t, err)

	scheduleTx = scheduleTx.
		SetNodeAccountIDs(env.NodeAccountIDs).
		SetPayerAccountID(accountID).
		SetAdminKey(env.OperatorKey).
		SetTransactionID(TransactionIDGenerate(env.Client.GetOperatorAccountID()))

	resp, err = scheduleTx.Execute(env.Client)
	assert.NoError(t, err)

	receipt, err = resp.GetReceipt(env.Client)
	assert.NoError(t, err)

	scheduleID := *receipt.ScheduleID

	info, err := NewScheduleInfoQuery().
		SetNodeAccountIDs(env.NodeAccountIDs).
		SetScheduleID(scheduleID).
		Execute(env.Client)
	assert.NoError(t, err)
	assert.Equal(t, info.CreatorAccountID.String(), env.OperatorID.String())

	signTransaction, err := NewScheduleSignTransaction().
		SetNodeAccountIDs(env.NodeAccountIDs).
		SetScheduleID(scheduleID).
		FreezeWith(env.Client)
	assert.NoError(t, err)

	signTransaction.Sign(keys[0])
	signTransaction.Sign(keys[1])
	signTransaction.Sign(keys[2])

	resp, err = signTransaction.Execute(env.Client)
	assert.NoError(t, err)

	_, err = resp.GetReceipt(env.Client)
	assert.NoError(t, err)

	info2, err := NewScheduleInfoQuery().
		SetNodeAccountIDs(env.NodeAccountIDs).
		SetScheduleID(scheduleID).
		Execute(env.Client)
	assert.NoError(t, err)
	assert.NotNil(t, info2.ExecutedAt)

	nftInfo, err := NewTokenNftInfoQuery().
		ByNftID(tokenID.Nft(mintReceipt.SerialNumbers[0])).
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		Execute(env.Client)
	assert.NoError(t, err)

	assert.Equal(t, accountID.String(), nftInfo[0].AccountID.String())

	err = CloseIntegrationTestEnv(env, &tokenID)
	assert.NoError(t, err)
}

func TestScheduledTokenNftTransferTransactionSigned_Execute(t *testing.T) {
	env := NewIntegrationTestEnv(t)

	keys := make([]PrivateKey, 3)
	pubKeys := make([]PublicKey, 3)

	for i := range keys {
		newKey, err := GeneratePrivateKey()
		assert.NoError(t, err)

		keys[i] = newKey
		pubKeys[i] = newKey.PublicKey()
	}

	keyList := NewKeyList().
		AddAllPublicKeys(pubKeys)

	newBalance := NewHbar(2)

	assert.Equal(t, 2*HbarUnits.Hbar.numberOfTinybar(), newBalance.tinybar)

	resp, err := NewAccountCreateTransaction().
		SetNodeAccountIDs(env.NodeAccountIDs).
		SetKey(keyList).
		SetInitialBalance(newBalance).
		Execute(env.Client)
	assert.NoError(t, err)

	receipt, err := resp.GetReceipt(env.Client)
	assert.NoError(t, err)

	accountID := *receipt.AccountID

	resp, err = NewTokenCreateTransaction().
		SetNodeAccountIDs(env.NodeAccountIDs).
		SetTokenName("ffff").
		SetTokenSymbol("F").
		SetTokenType(TokenTypeNonFungibleUnique).
		SetSupplyType(TokenSupplyTypeFinite).
		SetMaxSupply(5).
		SetTreasuryAccountID(env.Client.GetOperatorAccountID()).
		SetAdminKey(env.Client.GetOperatorPublicKey()).
		SetFreezeKey(env.Client.GetOperatorPublicKey()).
		SetWipeKey(env.Client.GetOperatorPublicKey()).
		SetKycKey(env.Client.GetOperatorPublicKey()).
		SetSupplyKey(env.Client.GetOperatorPublicKey()).
		SetFreezeDefault(false).
		Execute(env.Client)
	assert.NoError(t, err)

	receipt, err = resp.GetReceipt(env.Client)
	assert.NoError(t, err)

	tokenID := *receipt.TokenID
	metaData := [][]byte{{50}, {50}}

	mint, err := NewTokenMintTransaction().
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		SetTokenID(tokenID).
		SetMetadatas(metaData).
		Execute(env.Client)
	assert.NoError(t, err)

	mintReceipt, err := mint.GetReceipt(env.Client)
	assert.NoError(t, err)

	transaction, err := NewTokenAssociateTransaction().
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		SetAccountID(accountID).
		SetTokenIDs(tokenID).
		FreezeWith(env.Client)
	assert.NoError(t, err)

	resp, err = transaction.
		Sign(keys[0]).
		Sign(keys[1]).
		Sign(keys[2]).
		Execute(env.Client)
	assert.NoError(t, err)

	_, err = resp.GetReceipt(env.Client)
	assert.NoError(t, err)

	resp, err = NewTokenGrantKycTransaction().
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		SetAccountID(accountID).
		SetTokenID(tokenID).
		Execute(env.Client)
	assert.NoError(t, err)

	_, err = resp.GetReceipt(env.Client)
	assert.NoError(t, err)

	tx := NewTransferTransaction().
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		AddNftTransfer(tokenID.Nft(mintReceipt.SerialNumbers[0]), env.OperatorID, accountID).
		AddNftTransfer(tokenID.Nft(mintReceipt.SerialNumbers[1]), env.OperatorID, accountID)

	scheduleTx, err := tx.Schedule()
	assert.NoError(t, err)

	scheduleTx, err = scheduleTx.
		SetNodeAccountIDs(env.NodeAccountIDs).
		SetPayerAccountID(accountID).
		SetAdminKey(env.OperatorKey).
		SetTransactionID(TransactionIDGenerate(env.Client.GetOperatorAccountID())).
		FreezeWith(env.Client)
	assert.NoError(t, err)

	scheduleTx.Sign(keys[0])
	scheduleTx.Sign(keys[1])
	scheduleTx.Sign(keys[2])

	resp, err = scheduleTx.Execute(env.Client)
	assert.NoError(t, err)

	receipt, err = resp.GetReceipt(env.Client)
	assert.NoError(t, err)

	scheduleID := *receipt.ScheduleID

	info2, err := NewScheduleInfoQuery().
		SetNodeAccountIDs(env.NodeAccountIDs).
		SetScheduleID(scheduleID).
		Execute(env.Client)
	assert.NoError(t, err)
	assert.NotNil(t, info2.ExecutedAt)

	nftInfo, err := NewTokenNftInfoQuery().
		ByNftID(tokenID.Nft(mintReceipt.SerialNumbers[0])).
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		Execute(env.Client)
	assert.NoError(t, err)

	assert.Equal(t, accountID.String(), nftInfo[0].AccountID.String())

	err = CloseIntegrationTestEnv(env, &tokenID)
	assert.NoError(t, err)
}
