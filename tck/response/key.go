package response

// SPDX-License-Identifier: Apache-2.0

type GenerateKeyResponse struct {
	Key         string   `json:"key"`
	PrivateKeys []string `json:"privateKeys"`
}
