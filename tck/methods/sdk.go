package methods

import (
	"context"
	"strconv"

	"github.com/hashgraph/hedera-sdk-go/tck/param"
	"github.com/hashgraph/hedera-sdk-go/tck/response"
	"github.com/hashgraph/hedera-sdk-go/v2"
)

// Define the SDK service structure
type SDKService struct {
	Client *hedera.Client
}

// Setup function for the SDK
func (s *SDKService) Setup(_ context.Context, params param.SetupParams) (*response.SetupResponse, error) {
	var clientType string
	if params.NodeIp != nil && params.NodeAccountId != nil && params.MirrorNetworkIp != nil {
		// Custom client setup
		num, err := strconv.ParseUint(*params.NodeAccountId, 10, 64)
		if err != nil {
			return nil, response.InvalidParams.WithData(err.Error())
		}
		node := map[string]hedera.AccountID{
			*params.NodeIp: {Account: num},
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
	return response.NewSetupReponse("Successfully setup " + clientType + " client."), nil
}

// Reset function for the SDK
func (s *SDKService) Reset(_ context.Context) (*response.SetupResponse, error) {
	s.Client = nil
	return response.NewSetupReponse(""), nil
}
