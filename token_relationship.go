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
import "github.com/pkg/errors"

// TokenRelationship is the information about a token relationship
type TokenRelationship struct {
	TokenID              TokenID
	Balance              float64
	KycStatus            *bool
	FreezeStatus         *bool
	Decimals             float64
	AutomaticAssociation bool
}

func TokenRelationshipFromJson(tokenObject interface{}) (*TokenRelationship, error) {
	tokenJSON, ok := tokenObject.(map[string]interface{})
	if !ok {
		return nil, errors.New("Invalid token JSON object")
	}

	tokenId, err := TokenIDFromString(tokenJSON["token_id"].(string))
	if err != nil {
		return nil, err
	}

	var freezeStatus *bool
	if tokenFreezeStatus, ok := tokenJSON["freeze_status"].(string); ok {
		var freezeStatusValue bool
		if tokenFreezeStatus == "FROZEN" || tokenFreezeStatus == "UNFROZEN" {
			freezeStatusValue = tokenFreezeStatus == "FROZEN"
			freezeStatus = &freezeStatusValue
		} else {
			freezeStatus = nil
		}
	}

	var kycStatus *bool
	if tokenKycStatus, ok := tokenJSON["kyc_status"].(string); ok {
		var kycStatusValue bool
		if tokenKycStatus == "GRANTED" || tokenKycStatus == "REVOKED" {
			kycStatusValue = tokenKycStatus == "GRANTED"
			kycStatus = &kycStatusValue
		} else {
			kycStatus = nil
		}
	}

	tokenRelationship := &TokenRelationship{
		TokenID:              tokenId,
		Balance:              float64(tokenJSON["balance"].(float64)),
		KycStatus:            kycStatus,
		FreezeStatus:         freezeStatus,
		Decimals:             tokenJSON["decimals"].(float64),
		AutomaticAssociation: tokenJSON["automatic_association"].(bool),
	}
	return tokenRelationship, nil
}
