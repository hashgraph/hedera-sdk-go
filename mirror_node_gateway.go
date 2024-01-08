package hedera

import (
	"encoding/json"
	"errors"
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
const httpsPrefix = "https://"
const apiPathVersion = "/api/v1"

var queryTypes = map[string]string{
	"account":  "accounts",
	"contract": "contracts",
	"token":    "tokens"}

// Function to obtain balance of tokens for given account ID. Return the pure JSON response as mapping
func accountBalanceQuery(network string, accountId string) (map[string]interface{}, error) {
	info, err := accountInfoQuery(network, accountId)
	// Cast balance body to map
	return info["balance"].(map[string]interface{}), err
}

// Function to obtain account info for given account ID. Return the pure JSON response as mapping
func accountInfoQuery(network string, accountId string) (map[string]interface{}, error) {
	accountInfoUrl := buildUrl(network, queryTypes["account"], accountId)
	return makeGetRequest(accountInfoUrl)
}

// Function to obtain balance of tokens for given contract ID. Return the pure JSON response as mapping
func contractInfoQuery(network string, contractId string) (map[string]interface{}, error) {
	contractInfoUrl := buildUrl(network, queryTypes["contract"], contractId)
	return makeGetRequest(contractInfoUrl)
}

func tokenReleationshipQuery(network string, id string) (map[string]interface{}, error) {
	tokenRelationshipUrl := buildUrl(network, queryTypes["account"], id, queryTypes["token"])
	return makeGetRequest(tokenRelationshipUrl)
}

// Make a GET HTTP request to provided URL and map it's json response to a generic `interface` map and return it
func makeGetRequest(url string) (response map[string]interface{}, e error) {
	// Make an HTTP request
	resp, err := http.Get(url)
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

func obtainUrlForMirrorNode(client *Client) string {
	const localNetwork = "127.0.0.1"
	if client.GetMirrorNetwork()[0] == localNetwork+":5600" || client.GetMirrorNetwork()[0] == localNetwork+":443" {
		return localNetwork + "5551"
	} else {
		return client.GetMirrorNetwork()[0]
	}
}

func buildUrl(network string, args ...string) string {
	url := httpsPrefix + network + apiPathVersion
	for _, arg := range args {
		url += "/" + arg
	}
	return url
}

func mapTokenRelationship(tokens []interface{}) ([]*TokenRelationship, error) {
	var tokenRelationships []*TokenRelationship

	for _, tokenObj := range tokens {
		token, ok := tokenObj.(map[string]interface{})
		if !ok {
			return nil, errors.New("invalid token object")
		}

		tokenId, err := TokenIDFromString(token["token_id"].(string))
		if err != nil {
			return nil, err
		}

		freezeStatus := false
		if tokenFreezeStatus, ok := token["freeze_status"].(string); ok {
			freezeStatus = tokenFreezeStatus == "FROZEN"
		}

		kycStatus := false
		if tokenKycStatus, ok := token["kyc_status"].(string); ok {
			kycStatus = tokenKycStatus == "GRANTED"
		}

		tokenRelationship := &TokenRelationship{
			TokenID:              tokenId,
			Balance:              uint64(token["balance"].(float64)),
			KycStatus:            &kycStatus,
			FreezeStatus:         &freezeStatus,
			AutomaticAssociation: token["automatic_association"].(bool),
		}

		tokenRelationships = append(tokenRelationships, tokenRelationship)
	}

	return tokenRelationships, nil
}
