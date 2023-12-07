package methods

import (
	"context"
	"github.com/hashgraph/hedera-sdk-go/v2"
)

// Define the SDK service structure
type SDKService struct {
	Client *hedera.Client
}

// Parameters for the setup function
type SetupParams struct {
	OperatorAccountId  string
	OperatorPrivateKey string
	NodeIp             *string
	NodeAccountId      *uint64
	MirrorNetworkIp    *string
}

// Response structure for the setup function
type SetupResponse struct {
	Message string
	Status  string
}

// Setup function for the SDK
func (s *SDKService) Setup(_ context.Context, params SetupParams) SetupResponse {
	var clientType string

	if params.NodeIp != nil && params.NodeAccountId != nil && params.MirrorNetworkIp != nil {
		// Custom client setup
		node := map[string]hedera.AccountID{
			*params.NodeIp: {Account: *params.NodeAccountId},
		}
		s.Client = hedera.ClientForNetwork(node)
		clientType = "custom"
		s.Client.SetMirrorNetwork([]string{*params.MirrorNetworkIp})
	} else {
		// Default to testnet
		s.Client = hedera.ClientForTestnet()
		clientType = "testnet"
	}

	// Set operator (adjustments may be needed based on the Hedera SDK)
	operatorId, _ := hedera.AccountIDFromString(params.OperatorAccountId)
	operatorKey, _ := hedera.PrivateKeyFromString(params.OperatorPrivateKey)
	s.Client.SetOperator(operatorId, operatorKey)

	return SetupResponse{
		Message: "Successfully setup " + clientType + " client.",
		Status:  "SUCCESS",
	}
}

// Reset function for the SDK
func (s *SDKService) Reset(_ context.Context) SetupResponse {
	s.Client = nil
	return SetupResponse{
		Status: "SUCCESS",
	}
}
