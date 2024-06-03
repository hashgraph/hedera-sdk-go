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
	"strings"
)

// Function to obtain the token relationships of the specified account
func tokenRelationshipMirrorNodeQuery(networkUrl string, id string) (map[string]interface{}, error) {
	fmt.Println("accountID:", id)
	tokenRelationshipUrl := buildUrlParams(networkUrl, "accounts", id, "tokens")
	return makeGetRequest(tokenRelationshipUrl)
}

// Function to obtain account info for given account ID. Return the pure JSON response as mapping
func accountInfoMirrorNodeQuery(networkUrl string, accountId string) (map[string]interface{}, error) { // nolint
	accountInfoUrl := buildUrlParams(networkUrl, "accounts", accountId)
	return makeGetRequest(accountInfoUrl)
}

// Function to obtain balance of tokens for given contract ID. Return the pure JSON response as mapping
func contractInfoMirrorNodeQuery(networkUrl string, contractId string) (map[string]interface{}, error) { // nolint
	contractInfoUrl := buildUrlParams(networkUrl, "contracts", contractId)
	return makeGetRequest(contractInfoUrl)
}

// Function to obtain balance of tokens for given account ID. Return the pure JSON response as mapping
func accountTokenBalanceMirrorNodeQuery(networkUrl string, accountId string) (map[string]interface{}, error) {
	info, err := tokenRelationshipMirrorNodeQuery(networkUrl, accountId)

	// in case of empty info we won't be able to map to string interface
	if len(info) == 0 {
		return nil, err
	}
	return info, err
}

// Function to deduce the current network from the client as the network is ambiguous during Mirror Node calls
func fetchMirrorNodeUrlFromClient(client *Client) string {
	const localNetwork = "127.0.0.1"
	const apiVersion = "/api/v1"
	if strings.HasPrefix(client.GetMirrorNetwork()[0], localNetwork) {
		return "http://" + localNetwork + ":5551" + apiVersion
	} else {
		// prefix is mainnet, testnet or previewnet
		return "https://" + client.GetMirrorNetwork()[0] + apiVersion
	}
}

// Make a GET HTTP request to provided URL and map it's JSON response to a generic `interface` map and return it
func makeGetRequest(networkUrl string) (response map[string]interface{}, e error) {
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
	var responseMap map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&responseMap)
	if err != nil {
		return nil, err
	}

	return responseMap, nil
}

// Function to build url parameters
func buildUrlParams(networkUrl string, params ...string) string {
	for _, arg := range params {
		networkUrl += "/" + arg
	}
	return networkUrl
}
