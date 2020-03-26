package hedera

import (
	"testing"
	"time"
)

func TestMirrorConsensusTopicQuery(t *testing.T) {
	NewMirrorConsensusTopicQuery().
		SetTopicID(ConsensusTopicID{Topic: 99}).
		SetStartTime(time.Unix(0, 0)).
		SetEndTime(time.Unix(9, 9)).
		SetLimit(100)
}
