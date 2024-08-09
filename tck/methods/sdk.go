package methods

import (
	"context"

	"github.com/hashgraph/hedera-sdk-go/tck/param"
	"github.com/hashgraph/hedera-sdk-go/tck/response"
	"github.com/hashgraph/hedera-sdk-go/v2"
)

type SDKService struct {
	Client *hedera.Client
}

// Setup function for the SDK
func (s *SDKService) Setup(_ context.Context, params param.SetupParams) (response.SetupResponse, error) {
	var clientType string

	if params.NodeIp != nil && params.NodeAccountId != nil && params.MirrorNetworkIp != nil {
		// Custom client setup
		nodeId, err := hedera.AccountIDFromString(*params.NodeAccountId)
		if err != nil {
			return response.SetupResponse{}, err
		}
		node := map[string]hedera.AccountID{
			*params.NodeIp: nodeId,
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

	return response.SetupResponse{
		Message: "Successfully setup " + clientType + " client.",
		Status:  "SUCCESS",
	}, nil
}

// Reset function for the SDK
func (s *SDKService) Reset(_ context.Context) response.SetupResponse {
	s.Client = nil
	return response.SetupResponse{
		Status: "SUCCESS",
	}
}
