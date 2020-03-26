package hedera

import (
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestSerializeContractDeleteTransaction(t *testing.T) {
	mockClient, err := newMockClient()
	assert.NoError(t, err)

	privateKey, err := Ed25519PrivateKeyFromString(mockPrivateKey)
	assert.NoError(t, err)

	tx, err := NewContractDeleteTransaction().
		SetContractID(ContractID{Contract: 5}).
		SetMaxTransactionFee(HbarFromTinybar(1e6)).
		SetTransactionID(testTransactionID).
		Build(mockClient)

	assert.NoError(t, err)

	tx.Sign(privateKey)

	assert.Equal(t, `bodyBytes:"\n\016\n\010\010\334\311\007\020\333\237\t\022\002\030\003\022\002\030\003\030\300\204=\"\002\010x\262\001\004\n\002\030\005"sigMap:<sigPair:<pubKeyPrefix:"\344\361\300\353L}\315\303\347\353\021p\263\010\212=\022\242\227\364\243\353\342\362\205\003\375g5F\355\216"ed25519:"\002\037\252\273\3554\227\240V\217\231\347~S\204\227.\222\036\033reSJ\315?\240\224\341\272\271X\"\307\366\235\211k\360\264i<\224\313\220\343\022_\301w\201~e\376\203\227\2522|kg\202w\005">>transactionID:<transactionValidStart:<seconds:124124nanos:151515>accountID:<accountNum:3>>nodeAccountID:<accountNum:3>transactionFee:1000000transactionValidDuration:<seconds:120>contractDeleteInstance:<contractID:<contractNum:5>>`, strings.ReplaceAll(strings.ReplaceAll(tx.String(), " ", ""), "\n", ""))
}

func TestSerializeContractDeleteTransaction_WithAccountIDObtainer(t *testing.T) {
	mockClient, err := newMockClient()
	assert.NoError(t, err)

	privateKey, err := Ed25519PrivateKeyFromString(mockPrivateKey)
	assert.NoError(t, err)

	tx, err := NewContractDeleteTransaction().
		SetContractID(ContractID{Contract: 5}).
		SetMaxTransactionFee(HbarFromTinybar(1e6)).
		SetTransferAccountID(AccountID{Account: 3}).
		SetTransactionID(testTransactionID).
		Build(mockClient)

	assert.NoError(t, err)

	tx.Sign(privateKey)

	assert.Equal(t, `bodyBytes:"\n\016\n\010\010\334\311\007\020\333\237\t\022\002\030\003\022\002\030\003\030\300\204=\"\002\010x\262\001\010\n\002\030\005\022\002\030\003"sigMap:<sigPair:<pubKeyPrefix:"\344\361\300\353L}\315\303\347\353\021p\263\010\212=\022\242\227\364\243\353\342\362\205\003\375g5F\355\216"ed25519:"\001\221y\266\365\355\330O\004\373&\004\227;\034)\027\320\23010\240\343?\240|\004\315\326\300\317\342-\322\325\354\027\332\374\005\t\331\320\361\262K=Vr'zb\014\347Z\342\374\0356B(\336\003\017">>transactionID:<transactionValidStart:<seconds:124124nanos:151515>accountID:<accountNum:3>>nodeAccountID:<accountNum:3>transactionFee:1000000transactionValidDuration:<seconds:120>contractDeleteInstance:<contractID:<contractNum:5>transferAccountID:<accountNum:3>>`, strings.ReplaceAll(strings.ReplaceAll(tx.String(), " ", ""), "\n", ""))
}

func TestSerializeContractDeleteTransaction_WithContractIDObtainer(t *testing.T) {
	mockClient, err := newMockClient()
	assert.NoError(t, err)

	privateKey, err := Ed25519PrivateKeyFromString(mockPrivateKey)
	assert.NoError(t, err)

	tx, err := NewContractDeleteTransaction().
		SetContractID(ContractID{Contract: 5}).
		SetMaxTransactionFee(HbarFromTinybar(1e6)).
		SetTransferContractID(ContractID{Contract: 3}).
		SetTransactionID(testTransactionID).
		Build(mockClient)

	assert.NoError(t, err)

	tx.Sign(privateKey)

	assert.Equal(t, `bodyBytes:"\n\016\n\010\010\334\311\007\020\333\237\t\022\002\030\003\022\002\030\003\030\300\204=\"\002\010x\262\001\010\n\002\030\005\032\002\030\003"sigMap:<sigPair:<pubKeyPrefix:"\344\361\300\353L}\315\303\347\353\021p\263\010\212=\022\242\227\364\243\353\342\362\205\003\375g5F\355\216"ed25519:"\265\353ah\312\304mn\206ul\234\341F[pJ\"\342\352\220&wl\315\310UD\352:$GQ\326U\204\003\177\204\215\315k\277\342\376W]\377\312\037\237D\230aa\032\370>t\203\345\310c\017">>transactionID:<transactionValidStart:<seconds:124124nanos:151515>accountID:<accountNum:3>>nodeAccountID:<accountNum:3>transactionFee:1000000transactionValidDuration:<seconds:120>contractDeleteInstance:<contractID:<contractNum:5>transferContractID:<contractNum:3>>`, strings.ReplaceAll(strings.ReplaceAll(tx.String(), " ", ""), "\n", ""))
}
