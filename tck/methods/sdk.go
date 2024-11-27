package methods

// SPDX-License-Identifier: Apache-2.0

import (
	"context"

	"github.com/hiero-ledger/hiero-sdk-go/tck/param"
	"github.com/hiero-ledger/hiero-sdk-go/tck/response"
	"github.com/hiero-ledger/hiero-sdk-go/v2"
)

type SDKService struct {
	Client *hiero.Client
}

// Setup function for the SDK
func (s *SDKService) Setup(_ context.Context, params param.SetupParams) (response.SetupResponse, error) {
	var clientType string

	if params.NodeIp != nil && params.NodeAccountId != nil && params.MirrorNetworkIp != nil {
		// Custom client setup
		nodeId, err := hiero.AccountIDFromString(*params.NodeAccountId)
		if err != nil {
			return response.SetupResponse{}, err
		}
		node := map[string]hiero.AccountID{
			*params.NodeIp: nodeId,
		}
		s.Client = hiero.ClientForNetwork(node)
		clientType = "custom"
		s.Client.SetMirrorNetwork([]string{*params.MirrorNetworkIp})
	} else {
		// Default to testnet
		s.Client = hiero.ClientForTestnet()
		clientType = "testnet"
	}

	// Set operator (adjustments may be needed based on the Hiero SDK)
	operatorId, _ := hiero.AccountIDFromString(params.OperatorAccountId)
	operatorKey, _ := hiero.PrivateKeyFromString(params.OperatorPrivateKey)
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
