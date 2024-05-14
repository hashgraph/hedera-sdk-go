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
func TokenReleationshipMirrorNodeQuery(networkUrl string, id string) (map[string]interface{}, error) {
	tokenRelationshipUrl := BuildUrl(networkUrl, "accounts", id, "tokens")
	return MakeGetRequest(tokenRelationshipUrl)
}

// Function to obtain account info for given account ID. Return the pure JSON response as mapping
func AccountInfoMirrorNodeQuery(networkUrl string, accountId string) (map[string]interface{}, error) {
	accountInfoUrl := BuildUrl(networkUrl, "accounts", accountId)
	return MakeGetRequest(accountInfoUrl)
}

// Function to obtain balance of tokens for given account ID. Return the pure JSON response as mapping
func AccountBalanceMirrorNodeQuery(networkUrl string, accountId string) (map[string]interface{}, error) {
	// accountInfoMirrorNodeQuery provides the needed data this function exists only for the convenience of naming
	info, err := AccountInfoMirrorNodeQuery(networkUrl, accountId)
	// check in case of empty tokenBalances
	if len(info) == 0 {
		return nil, nil
	}
	return info["balance"].(map[string]interface{}), err
}

// Function to obtain balance of tokens for given contract ID. Return the pure JSON response as mapping
func ContractInfoMirrorNodeQuery(networkUrl string, contractId string) (map[string]interface{}, error) { // nolint
	contractInfoUrl := BuildUrl(networkUrl, "contracts", contractId)
	return MakeGetRequest(contractInfoUrl)
}

// Make a GET HTTP request to provided URL and map it's json response to a generic `interface` map and return it
func MakeGetRequest(networkUrl string) (response map[string]interface{}, e error) {
	// Make an HTTP request
	resp, err := http.Get(networkUrl) //nolint

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
func FetchMirrorNodeUrlFromClient(client *Client) string {
	const localNetwork = "127.0.0.1"
	const apiVersion = "/api/v1"
	if client.GetMirrorNetwork()[0] == localNetwork+":5600" || client.GetMirrorNetwork()[0] == localNetwork+":443" {
		return "http://" + localNetwork + ":5551" + apiVersion
	} else {
		return "https://" + client.GetMirrorNetwork()[0] + apiVersion
	}
}

// This function takes the current network(localhost,testnet,previewnet,mainnet) adds the current api version hardcore style
// and concatenates further parameters for the call to MirrorNode
func BuildUrl(networkUrl string, params ...string) string {
	for _, arg := range params {
		networkUrl += "/" + arg
	}
	return networkUrl
}
