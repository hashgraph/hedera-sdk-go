//go:build all || unit
// +build all unit

package hiero

import (
	"testing"

	"github.com/hiero-ledger/hiero-sdk-go/v2/proto/services"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"
)

func TestCustomFractionalFee_SettersAndGetters(t *testing.T) {
	t.Parallel()

	fee := NewCustomFractionalFee()
	const FeeAssessmentMethod_TRUE FeeAssessmentMethod = true

	fee.SetFeeCollectorAccountID(AccountID{Account: 1234}).
		SetNumerator(1).
		SetDenominator(2).
		SetMin(10).
		SetMax(100).
		SetAllCollectorsAreExempt(true).
		SetAssessmentMethod(FeeAssessmentMethod_TRUE)

	require.Equal(t, AccountID{Account: 1234}, fee.GetFeeCollectorAccountID())
	require.Equal(t, int64(1), fee.GetNumerator())
	require.Equal(t, int64(2), fee.GetDenominator())
	require.Equal(t, int64(10), fee.GetMin())
	require.Equal(t, int64(100), fee.GetMax())
	require.Equal(t, true, fee.GetAllCollectorsAreExempt())
	require.Equal(t, FeeAssessmentMethod_TRUE, fee.GetAssessmentMethod())
}

func TestCustomFractionalFee_ToBytes(t *testing.T) {
	t.Parallel()

	fee := NewCustomFractionalFee()
	const FeeAssessmentMethod_TRUE FeeAssessmentMethod = true

	fee.SetFeeCollectorAccountID(AccountID{Account: 1234}).
		SetNumerator(1).
		SetDenominator(2).
		SetMin(10).
		SetMax(100).
		SetAllCollectorsAreExempt(true).
		SetAssessmentMethod(FeeAssessmentMethod_TRUE)

	bytes := fee.ToBytes()
	require.NotEmpty(t, bytes)

	protoFee := &services.CustomFee{}
	err := proto.Unmarshal(bytes, protoFee)
	require.NoError(t, err)

	require.Equal(t, int64(1), protoFee.GetFractionalFee().GetFractionalAmount().GetNumerator())
	require.Equal(t, int64(2), protoFee.GetFractionalFee().GetFractionalAmount().GetDenominator())
	require.Equal(t, int64(10), protoFee.GetFractionalFee().GetMinimumAmount())
	require.Equal(t, int64(100), protoFee.GetFractionalFee().GetMaximumAmount())
	require.Equal(t, true, protoFee.GetFractionalFee().GetNetOfTransfers())
	require.Equal(t, int64(1234), protoFee.GetFeeCollectorAccountId().GetAccountNum())
	require.Equal(t, true, protoFee.GetAllCollectorsAreExempt())
}

func TestCustomFractionalFee_ToProtobuf(t *testing.T) {
	t.Parallel()

	fee := NewCustomFractionalFee()
	const FeeAssessmentMethod_TRUE FeeAssessmentMethod = true

	fee.SetFeeCollectorAccountID(AccountID{Account: 1234}).
		SetNumerator(1).
		SetDenominator(2).
		SetMin(10).
		SetMax(100).
		SetAllCollectorsAreExempt(true).
		SetAssessmentMethod(FeeAssessmentMethod_TRUE)

	protoFee := fee._ToProtobuf()
	require.Equal(t, int64(1), protoFee.GetFractionalFee().GetFractionalAmount().GetNumerator())
	require.Equal(t, int64(2), protoFee.GetFractionalFee().GetFractionalAmount().GetDenominator())
	require.Equal(t, int64(10), protoFee.GetFractionalFee().GetMinimumAmount())
	require.Equal(t, int64(100), protoFee.GetFractionalFee().GetMaximumAmount())
	require.Equal(t, true, protoFee.GetFractionalFee().GetNetOfTransfers())
	require.Equal(t, int64(1234), protoFee.GetFeeCollectorAccountId().GetAccountNum())
	require.Equal(t, true, protoFee.GetAllCollectorsAreExempt())
	// CustomFeeFromProtobuf calls _CustomFractionalFeeFromProtobuf
	deserializedFee := _CustomFeeFromProtobuf(protoFee)

	require.Equal(t, *fee, deserializedFee)
}
