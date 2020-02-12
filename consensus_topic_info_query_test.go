package hedera

import (
	"github.com/bradleyjkemp/cupaloy"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestConsensusTopicInfoQuery(t *testing.T) {
	mockTransaction, err := newMockTransaction()
	assert.NoError(t, err)

	query := NewConsensusTopicInfoQuery().
		SetTopicID(ConsensusTopicID{Topic: 99}).
		SetQueryPaymentTransaction(mockTransaction)

	cupaloy.SnapshotT(t, query.pb.String())
}
