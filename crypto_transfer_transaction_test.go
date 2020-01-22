package hedera

import (
	"github.com/bradleyjkemp/cupaloy"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSerializeCryptoTransferTransaction(t *testing.T) {
	tx, err := newMockTransaction()
	assert.NoError(t, err)

	cupaloy.SnapshotT(t, tx.String())
}
