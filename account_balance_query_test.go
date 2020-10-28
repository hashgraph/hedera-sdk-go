package hedera

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewAccountBalanceQuery(t *testing.T) {
	mockTransaction, err := newMockTransaction()
	assert.NoError(t, err)

	query := NewAccountBalanceQuery().
		SetAccountID(AccountID{Account: 3}).
		SetQueryPaymentTransaction(mockTransaction)

	assert.Equal(t, `cryptogetAccountBalance:{header:{payment:{bodyBytes:"\n\x0e\n\x08\x08\xdc\xc9\x07\x10۟\t\x12\x02\x18\x03\x12\x02\x18\x03\x18\x80\xc2\xd7/\"\x02\x08xr\x14\n\x12\n\x07\n\x02\x18\x02\x10\xc7\x01\n\x07\n\x02\x18\x03\x10\xc8\x01"sigMap:{sigPair:{pubKeyPrefix:"\xe4\xf1\xc0\xebL}\xcd\xc3\xe7\xeb\x11p\xb3\x08\x8a=\x12\xa2\x97\xf4\xa3\xeb\xe2\xf2\x85\x03\xfdg5F\xed\x8e"ed25519:"\x12&5\x96\xfb\xb4\x1c]P\xbb%\xecP\x9bk͙\x0b߼\xac)\xa6+\xd2<\x97+\xbb\x8c\x8af\xcb\xdai\x17T4{\xf7\xf3UYn\n\x8f\xabep\x04\xf6\x83\x0f\xbaFUP\xa3\xd1/\x1d\x9d\x1a\x0b"}}}}accountID:{accountNum:3}}`, strings.ReplaceAll(query.QueryBuilder.pb.String(), " ", ""))
}

func TestNewAccountBalanceQuery_ForContract(t *testing.T) {
	mockTransaction, err := newMockTransaction()
	assert.NoError(t, err)

	query := NewAccountBalanceQuery().
		SetContractID(ContractID{Contract: 3}).
		SetQueryPaymentTransaction(mockTransaction)

	assert.Equal(t, `cryptogetAccountBalance:{header:{payment:{bodyBytes:"\n\x0e\n\x08\x08\xdc\xc9\x07\x10۟\t\x12\x02\x18\x03\x12\x02\x18\x03\x18\x80\xc2\xd7/\"\x02\x08xr\x14\n\x12\n\x07\n\x02\x18\x02\x10\xc7\x01\n\x07\n\x02\x18\x03\x10\xc8\x01"sigMap:{sigPair:{pubKeyPrefix:"\xe4\xf1\xc0\xebL}\xcd\xc3\xe7\xeb\x11p\xb3\x08\x8a=\x12\xa2\x97\xf4\xa3\xeb\xe2\xf2\x85\x03\xfdg5F\xed\x8e"ed25519:"\x12&5\x96\xfb\xb4\x1c]P\xbb%\xecP\x9bk͙\x0b߼\xac)\xa6+\xd2<\x97+\xbb\x8c\x8af\xcb\xdai\x17T4{\xf7\xf3UYn\n\x8f\xabep\x04\xf6\x83\x0f\xbaFUP\xa3\xd1/\x1d\x9d\x1a\x0b"}}}}contractID:{contractNum:3}}`, strings.ReplaceAll(query.QueryBuilder.pb.String(), " ", ""))
}

func TestAccountBalanceQuery_Execute(t *testing.T) {
	client := newTestClient(t)

	newKey, err := GenerateEd25519PrivateKey()
	assert.NoError(t, err)

	newBalance := NewHbar(1)

	assert.Equal(t, HbarUnits.Hbar.numberOfTinybar(), newBalance.tinybar)

	txID, err := NewAccountCreateTransaction().
		SetKey(newKey.PublicKey()).
		SetMaxTransactionFee(NewHbar(2)).
		SetInitialBalance(newBalance).
		Execute(client)
	assert.NoError(t, err)

	receipt, err := txID.GetReceipt(client)
	assert.NoError(t, err)

	accountID := receipt.GetAccountID()
	assert.NoError(t, err)

	balance, err := NewAccountBalanceQuery().
		SetAccountID(accountID).
		Execute(client)
	assert.NoError(t, err)
	assert.Equal(t, newBalance, balance)

	tx, err := NewAccountDeleteTransaction().
		SetDeleteAccountID(accountID).
		SetTransferAccountID(client.GetOperatorID()).
		SetMaxTransactionFee(NewHbar(1)).
		SetTransactionID(NewTransactionID(accountID)).
		Build(client)
	assert.NoError(t, err)

	txID, err = tx.
		Sign(newKey).
		Execute(client)
	assert.NoError(t, err)

	_, err = txID.GetReceipt(client)
	assert.NoError(t, err)
}
