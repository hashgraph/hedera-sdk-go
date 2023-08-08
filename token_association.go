package hedera

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

import (
	"github.com/hashgraph/hedera-protobufs-go/services"
	protobuf "google.golang.org/protobuf/proto"
)

// A token - account association
type TokenAssociation struct {
	TokenID   *TokenID
	AccountID *AccountID
}

func tokenAssociationFromProtobuf(pb *services.TokenAssociation) TokenAssociation {
	if pb == nil {
		return TokenAssociation{}
	}

	return TokenAssociation{
		TokenID:   _TokenIDFromProtobuf(pb.TokenId),
		AccountID: _AccountIDFromProtobuf(pb.AccountId),
	}
}

func (association *TokenAssociation) toProtobuf() *services.TokenAssociation {
	var tokenID *services.TokenID
	if association.TokenID != nil {
		tokenID = association.TokenID._ToProtobuf()
	}

	var accountID *services.AccountID
	if association.AccountID != nil {
		accountID = association.AccountID._ToProtobuf()
	}

	return &services.TokenAssociation{
		TokenId:   tokenID,
		AccountId: accountID,
	}
}

// ToBytes returns the byte representation of the TokenAssociation
func (association *TokenAssociation) ToBytes() []byte {
	data, err := protobuf.Marshal(association.toProtobuf())
	if err != nil {
		return make([]byte, 0)
	}

	return data
}

// TokenAssociationFromBytes returns a TokenAssociation from a raw protobuf byte array
func TokenAssociationFromBytes(data []byte) (TokenAssociation, error) {
	if data == nil {
		return TokenAssociation{}, errByteArrayNull
	}
	pb := services.TokenAssociation{}
	err := protobuf.Unmarshal(data, &pb)
	if err != nil {
		return TokenAssociation{}, err
	}

	association := tokenAssociationFromProtobuf(&pb)

	return association, nil
}
