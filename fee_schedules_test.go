package hedera

// import (
// 	"github.com/stretchr/testify/assert"
// 	"testing"
// )
//
// func TestFetchFeeSchedules_Execute(t *testing.T) {
// 	env := NewIntegrationTestEnv(t)
//
// 	feeSchedulesBytes, err := NewFileContentsQuery().
// 		SetFileID(FileID{Shard: 0, Realm: 0, File: 111}).
// 		Execute(env.Client)
// 	assert.NoError(t, err)
// 	feeSchedules, err := FeeSchedulesFromBytes(feeSchedulesBytes)
// 	assert.NoError(t, err)
// 	assert.NotNil(t, feeSchedules)
// 	assert.Equal(t, feeSchedules.current.TransactionFeeSchedules[0].FeeData.NodeData.Constant, int64(4498129603))
// 	assert.Equal(t, feeSchedules.current.TransactionFeeSchedules[0].FeeData.ServiceData.Constant, int64(71970073651))
// 	assert.Equal(t, feeSchedules.current.TransactionFeeSchedules[0].RequestType, RequestTypeCryptoCreate)
//
// 	err = CloseIntegrationTestEnv(env, nil)
// 	assert.NoError(t, err)
// }
