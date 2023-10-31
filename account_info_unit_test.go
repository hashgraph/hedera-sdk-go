//go:build all || unit
// +build all unit

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
	"fmt"
	"reflect"
	"testing"
	"time"
)

func TestAccountInfoToBytes(t *testing.T) {
	t.Parallel()

	accInfoOriginal := _MockAccountInfo();
	fmt.Println("Original AccountInfo = ", accInfoOriginal.);
	accInfoBytes := accInfoOriginal.ToBytes();

	accInfoFromBytes,err := AccountInfoFromBytes(accInfoBytes);

	fmt.Println("AccountInfo obtained from bytes = ", accInfoFromBytes);

	if(err != nil){
		t.Fatalf("Error trying to parse from bytes to accountInfo. Error = %v", err );
	}
	if !reflect.DeepEqual(accInfoOriginal, accInfoFromBytes){
		t.Fatalf("AccountInfoToBytes() = %v, want %v", accInfoFromBytes, accInfoOriginal);
	}
}		

func _MockAccountInfo() *AccountInfo{

	privateKey,_ := PrivateKeyFromString(mockPrivateKey);
	accountID,_ := AccountIDFromString("0.0.123-esxsf")

		return &AccountInfo{
			AccountID:                      accountID,
			ContractAccountID:              "",
			IsDeleted:                      false,
			ProxyAccountID:                 AccountID{},
			ProxyReceived:                  Hbar{},
			Key:                            privateKey,
			Balance:                        Hbar{},
			GenerateSendRecordThreshold:    Hbar{},
			GenerateReceiveRecordThreshold: Hbar{},
			ReceiverSigRequired:            false,
			ExpirationTime:                 time.Now(),
			AutoRenewPeriod:                time.Hour,
			LiveHashes:                     []*LiveHash{},
			TokenRelationships:             []*TokenRelationship{},
			AccountMemo:                    "",
			OwnedNfts:                      0,
			MaxAutomaticTokenAssociations:  0,
			AliasKey:                       &PublicKey{},
			LedgerID:                       LedgerID{},
			HbarAllowances:                 []HbarAllowance{},
			NftAllowances:                  []TokenNftAllowance{},
			TokenAllowances:                []TokenAllowance{},
			EthereumNonce:                  0,
			StakingInfo:                    &StakingInfo{},
		}

}