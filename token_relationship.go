package hedera

import "errors"

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

// TokenRelationship is the information about a token relationship
type TokenRelationship struct {
	TokenID              TokenID
	Symbol               string
	Balance              uint64
	KycStatus            *bool
	FreezeStatus         *bool
	Decimals             uint32
	AutomaticAssociation bool
}

//	func _TokenRelationshipFromProtobuf(pb *services.TokenRelationship) TokenRelationship {
//		if pb == nil {
//			return TokenRelationship{}
//		}
//
//		tokenID := TokenID{}
//		if pb.TokenId != nil {
//			tokenID = *_TokenIDFromProtobuf(pb.TokenId)
//		}
//
//		return TokenRelationship{
//			TokenID:              tokenID,
//			Symbol:               pb.Symbol,
//			Balance:              pb.Balance,
//			KycStatus:            _KycStatusFromProtobuf(pb.KycStatus),
//			FreezeStatus:         _FreezeStatusFromProtobuf(pb.FreezeStatus),
//			Decimals:             pb.Decimals,
//			AutomaticAssociation: pb.AutomaticAssociation,
//		}
//	}
//
//	func (relationship *TokenRelationship) _ToProtobuf() *services.TokenRelationship {
//		var freezeStatus services.TokenFreezeStatus
//		switch *relationship.FreezeStatus {
//		case true:
//			freezeStatus = 1
//		case false:
//			freezeStatus = 2
//		default:
//			freezeStatus = 0
//		}
//
//		var kycStatus services.TokenKycStatus
//		switch *relationship.KycStatus {
//		case true:
//			kycStatus = 1
//		case false:
//			kycStatus = 2
//		default:
//			kycStatus = 0
//		}
//
//		return &services.TokenRelationship{
//			TokenId:              relationship.TokenID._ToProtobuf(),
//			Symbol:               relationship.Symbol,
//			Balance:              relationship.Balance,
//			KycStatus:            kycStatus,
//			FreezeStatus:         freezeStatus,
//			Decimals:             relationship.Decimals,
//			AutomaticAssociation: relationship.AutomaticAssociation,
//		}
//	}
//
//	func (relationship TokenRelationship) ToBytes() []byte {
//		data, err := protobuf.Marshal(relationship._ToProtobuf())
//		if err != nil {
//			return make([]byte, 0)
//		}
//
//		return data
//	}
//
//	func TokenRelationshipFromBytes(data []byte) (TokenRelationship, error) {
//		if data == nil {
//			return TokenRelationship{}, errors.New("byte array can't be null")
//		}
//		pb := services.TokenRelationship{}
//		err := protobuf.Unmarshal(data, &pb)
//		if err != nil {
//			return TokenRelationship{}, err
//		}
//
//		return _TokenRelationshipFromProtobuf(&pb), nil
//	}
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
