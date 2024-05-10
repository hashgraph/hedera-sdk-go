package hedera

/*-
 *
 * Hedera Go SDK
 *
 * Copyright (C) 2020 - 2024 Hedera Hashgraph, LLC
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */
import (
	"encoding/json"
	"fmt"
	"net/http"
)

// Function to obtain the token relationships of the specified account
func TokenReleationshipMirrorNodeQuery(network string, id string) (map[string]interface{}, error) {
	tokenRelationshipUrl := BuildUrl(network, "accounts", id, "tokens")
	return MakeGetRequest(tokenRelationshipUrl)
}

// Function to obtain account info for given account ID. Return the pure JSON response as mapping
func AccountInfoMirrorNodeQuery(network string, accountId string) (map[string]interface{}, error) {
	accountInfoUrl := BuildUrl(network, "accounts", accountId)
	return MakeGetRequest(accountInfoUrl)
}

// Function to obtain balance of tokens for given account ID. Return the pure JSON response as mapping
func AccountBalanceMirrorNodeQuery(network string, accountId string) (map[string]interface{}, error) {
	// accountInfoMirrorNodeQuery provides the needed data this function exists only for the convenience of naming
	info, err := AccountInfoMirrorNodeQuery(network, accountId)
	return info["balance"].(map[string]interface{}), err
}

// Function to obtain balance of tokens for given contract ID. Return the pure JSON response as mapping
func ContractInfoMirrorNodeQuery(network string, contractId string) (map[string]interface{}, error) { // nolint
	contractInfoUrl := BuildUrl(network, "contracts", contractId)
	return MakeGetRequest(contractInfoUrl)
}

// Make a GET HTTP request to provided URL and map it's json response to a generic `interface` map and return it
func MakeGetRequest(url string) (response map[string]interface{}, e error) {
	// Make an HTTP request
	resp, err := http.Get(url) //nolint

	if err != nil {
		return nil, err
	}
	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("HTTP request failed with status code: %d", resp.StatusCode)
	}
	defer resp.Body.Close()

	// Decode the JSON response into a map
	var resultMap map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&resultMap)
	if err != nil {
		return nil, err
	}

	return resultMap, nil
}

// Uses the client to deduce the current network as the network is ambiguous during Mirror Node calls
func ObtainUrlForMirrorNode(client *Client) string {
	const localNetwork = "127.0.0.1"
	if client.GetMirrorNetwork()[0] == localNetwork+":5600" || client.GetMirrorNetwork()[0] == localNetwork+":443" {
		return localNetwork + "5551"
	} else {
		return client.GetMirrorNetwork()[0]
	}
}

// This function takes the current network(localhost,testnet,previewnet,mainnet) adds the current api version hardcore style
// and concatenates further parameters for the call to MirrorNode
func BuildUrl(network string, params ...string) string {
	url := "https://" + network + "/api/v1"
	for _, arg := range params {
		url += "/" + arg
	}
	return url
}
