//+build all unit

package hedera

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUnitTokenPause(t *testing.T) {
	accountID, err := AccountIDFromString("0.0.5005")
	assert.NoError(t, err)
	tokenID, err := TokenIDFromString("0.0.5005")
	assert.NoError(t, err)

	tx, err := NewTokenPauseTransaction().
		SetNodeAccountIDs([]AccountID{accountID}).
		SetTransactionID(TransactionIDGenerate(accountID)).
		SetTokenID(tokenID).
		Freeze()
	assert.NoError(t, err)

	pb := tx._Build()
	assert.Equal(t, pb.GetTokenPause().GetToken().String(), tokenID._ToProtobuf().String())
}

func TestUnitTokenUnpause(t *testing.T) {
	accountID, err := AccountIDFromString("0.0.5005")
	assert.NoError(t, err)
	tokenID, err := TokenIDFromString("0.0.5005")
	assert.NoError(t, err)

	tx, err := NewTokenUnpauseTransaction().
		SetNodeAccountIDs([]AccountID{accountID}).
		SetTransactionID(TransactionIDGenerate(accountID)).
		SetTokenID(tokenID).
		Freeze()
	assert.NoError(t, err)

	pb := tx._Build()
	assert.Equal(t, pb.GetTokenUnpause().GetToken().String(), tokenID._ToProtobuf().String())
}
