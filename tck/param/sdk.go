package param

// SPDX-License-Identifier: Apache-2.0

type SetupParams struct {
	OperatorAccountId  string  `json:"operatorAccountId"`
	OperatorPrivateKey string  `json:"operatorPrivateKey"`
	NodeIp             *string `json:"nodeIp"`
	NodeAccountId      *string `json:"nodeAccountId"`
	MirrorNetworkIp    *string `json:"mirrorNetworkIp"`
}
