package hedera

import (
	"encoding/json"
	"fmt"
	"net/http"
)

/*-
 *
 * Hedera Go SDK
 *
 * Copyright (C) 2020 - 2023 Hedera Hashgraph, LLC
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

const apiVersion = "v1/"

var restUrls = map[string]string{
	"mainnet":    "https://mainnet-public.mirrornode.hedera.com/api/" + apiVersion,
	"testnet":    "https://testnet.mirrornode.hedera.com/api/" + apiVersion,
	"previewnet": "https://previewnet.mirrornode.hedera.com/api/" + apiVersion,
}
var queryTypes = map[string]string{
	"account":  "accounts/",
	"contract": "contracts/"}

func accountBalanceQuery(network string, accountId string) (map[string]interface{}, error) {
	info, err := accountInfoQuery(network, accountId)
	// Cast balance body to map
	return info["balance"].(map[string]interface{}), err
}

func accountInfoQuery(network string, accountId string) (map[string]interface{}, error) {
	accountInfoUrl := restUrls[network] + queryTypes["account"] + accountId
	return makeGetRequest(accountInfoUrl)
}

func contractInfoQuery(network string, accountId string) (map[string]interface{}, error) {
	contractInfoUrl := restUrls[network] + queryTypes["contract"] + accountId
	return makeGetRequest(contractInfoUrl)
}

// Make a GET HTTP request to provided URL and map it's json response to a generic `interface` map and return it
func makeGetRequest(url string) (response map[string]interface{}, e error) {
	// Make an HTTP request
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Error making HTTP request:", err)
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
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
