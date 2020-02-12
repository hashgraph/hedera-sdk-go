package hedera

import (
	"github.com/bradleyjkemp/cupaloy"
	"testing"
	"time"
)

func TestMirrorConsensusTopicQuery(t *testing.T) {
	query := NewMirrorConsensusTopicQuery().
		SetTopicID(ConsensusTopicID{Topic: 99}).
		SetStartTime(time.Unix(0, 0)).
		SetEndTime(time.Unix(9, 9)).
		SetLimit(100)

	cupaloy.SnapshotT(t, query.pb.String())
}
