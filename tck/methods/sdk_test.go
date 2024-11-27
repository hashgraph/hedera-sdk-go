package methods

import (
	"context"
	"testing"

	"github.com/hiero-ledger/hiero-sdk-go/tck/param"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSetup(t *testing.T) {
	// // Given
	sdkService := &SDKService{}
	params := param.SetupParams{
		OperatorAccountId:  "0.0.2",
		OperatorPrivateKey: "302e020100300506032b65700422042091132178e72057a1d7528025956fe39b0b847f200ab59b2fdd367017f3087137",
		NodeIp:             stringPointer("127.0.0.1:50211"),
		NodeAccountId:      stringPointer("0.0.3"),
		MirrorNetworkIp:    stringPointer("http://127.0.0.1:5551"),
	}

	// When
	response, _ := sdkService.Setup(context.Background(), params)

	// Then
	assert.Equal(t, "Successfully setup custom client.", response.Message)
	assert.Equal(t, "SUCCESS", response.Status)
}

func TestSetupFail(t *testing.T) {
	// Given
	sdkService := &SDKService{}
	params := param.SetupParams{
		OperatorAccountId:  "operatorAccountId",
		OperatorPrivateKey: "operatorPrivateKey",
		NodeIp:             stringPointer("nodeIp"),
		NodeAccountId:      stringPointer("3asdf"),
		MirrorNetworkIp:    stringPointer("127.0.0.1:50211"),
	}

	// Then
	_, err := sdkService.Setup(context.Background(), params)
	require.Error(t, err)
}

func TestReset(t *testing.T) {
	// Given
	sdkService := &SDKService{}
	params := param.SetupParams{
		OperatorAccountId:  "0.0.2",
		OperatorPrivateKey: "302e020100300506032b65700422042091132178e72057a1d7528025956fe39b0b847f200ab59b2fdd367017f3087137",
		NodeIp:             stringPointer("127.0.0.1:50211"),
		NodeAccountId:      stringPointer("0.0.3"),
		MirrorNetworkIp:    stringPointer("http://127.0.0.1:5551"),
	}

	// Setup first to initialize the client
	_, err := sdkService.Setup(context.Background(), params)
	require.NoError(t, err)

	// When
	response := sdkService.Reset(context.Background())

	// Then
	assert.Equal(t, "SUCCESS", response.Status)
	assert.Nil(t, sdkService.Client)
}

func stringPointer(s string) *string {
	return &s
}

func intPointer(i int) *int {
	return &i
}
