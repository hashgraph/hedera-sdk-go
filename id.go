package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"

	"github.com/hiero-ledger/hiero-sdk-go/v2/proto/services"
	protobuf "google.golang.org/protobuf/proto"
)

func _AccountIDFromString(s string) (shard int, realm int, num int, checksum *string, alias *PublicKey, evmAddress *[]byte, err error) {
	if _Has0xPrefix(s) {
		s = _Without0x(s)
	}
	if _IsHex(s) {
		bytes := _Hex2Bytes(s)
		if err == nil {
			if len(bytes) == 20 {
				return 0, 0, -1, nil, nil, &bytes, nil
			}
		}
	}

	if strings.Contains(s, "-") {
		values := strings.SplitN(s, "-", 2)

		if len(values) > 2 {
			return 0, 0, 0, nil, nil, nil, fmt.Errorf("expected {shard}.{realm}.{num}-{checksum}")
		}

		checksum = &values[1]
		s = values[0]
	}

	values := strings.SplitN(s, ".", 3)
	if len(values) != 3 {
		// Was not three values separated by periods
		return 0, 0, 0, nil, nil, nil, fmt.Errorf("expected {shard}.{realm}.{num}")
	}

	shard, err = strconv.Atoi(values[0])
	if err != nil {
		return 0, 0, 0, nil, nil, nil, err
	}

	realm, err = strconv.Atoi(values[1])
	if err != nil {
		return 0, 0, 0, nil, nil, nil, err
	}

	if len(values[2]) < 20 {
		num, err = strconv.Atoi(values[2])
		if err != nil {
			return 0, 0, 0, nil, nil, nil, err
		}

		return shard, realm, num, checksum, nil, nil, nil
	} else if len(values[2]) == 40 {
		temp, err2 := hex.DecodeString(values[2])
		if err2 != nil {
			return 0, 0, 0, nil, nil, nil, err2
		}
		var key services.Key
		err2 = protobuf.Unmarshal(temp, &key)
		if err2 != nil {
			return shard, realm, -1, checksum, nil, &temp, nil
		}
		aliasKey, err2 := _KeyFromProtobuf(&key)
		if err2 != nil {
			return shard, realm, -1, checksum, nil, &temp, nil
		}

		if aliasPublicKey, ok := aliasKey.(PublicKey); ok {
			return shard, realm, -1, checksum, &aliasPublicKey, nil, nil
		}

		return shard, realm, -1, checksum, nil, &temp, nil
	}

	key, err := PublicKeyFromString(values[2])
	if err != nil {
		return 0, 0, 0, nil, nil, nil, err
	}

	return shard, realm, -1, checksum, &key, nil, nil
}

func _ContractIDFromString(s string) (shard int, realm int, num int, checksum *string, evmAddress []byte, err error) {
	if strings.Contains(s, "-") {
		values := strings.SplitN(s, "-", 2)

		if len(values) > 2 {
			return 0, 0, 0, nil, nil, fmt.Errorf("expected {shard}.{realm}.{num}-{checksum}")
		}

		checksum = &values[1]
		s = values[0]
	}

	values := strings.SplitN(s, ".", 3)
	if len(values) != 3 {
		// Was not three values separated by periods
		return 0, 0, 0, nil, nil, fmt.Errorf("expected {shard}.{realm}.{num}")
	}

	shard, err = strconv.Atoi(values[0])
	if err != nil {
		return 0, 0, 0, nil, nil, err
	}

	realm, err = strconv.Atoi(values[1])
	if err != nil {
		return 0, 0, 0, nil, nil, err
	}

	num, err = strconv.Atoi(values[2])
	if err != nil {
		temp, err2 := hex.DecodeString(values[2])
		if err2 != nil {
			return 0, 0, 0, nil, nil, err
		}
		return shard, realm, -1, checksum, temp, nil
	}

	return shard, realm, num, checksum, nil, nil
}

func _IdFromString(s string) (shard int, realm int, num int, checksum *string, err error) {
	if strings.Contains(s, "-") {
		values := strings.SplitN(s, "-", 2)

		if len(values) > 2 {
			return 0, 0, 0, nil, fmt.Errorf("expected {shard}.{realm}.{num}-{checksum}")
		}

		checksum = &values[1]
		s = values[0]
	}

	values := strings.SplitN(s, ".", 3)
	if len(values) != 3 {
		// Was not three values separated by periods
		return 0, 0, 0, nil, fmt.Errorf("expected {shard}.{realm}.{num}")
	}

	shard, err = strconv.Atoi(values[0])
	if err != nil {
		return 0, 0, 0, nil, err
	}

	realm, err = strconv.Atoi(values[1])
	if err != nil {
		return 0, 0, 0, nil, err
	}

	num, err = strconv.Atoi(values[2])
	if err != nil {
		return 0, 0, 0, nil, err
	}

	return shard, realm, num, checksum, nil
}

func _IdFromSolidityAddress(s string) (uint64, uint64, uint64, error) {
	bytes, err := hex.DecodeString(s)
	if err != nil {
		return 0, 0, 0, err
	}

	if len(bytes) != 20 {
		return 0, 0, 0, fmt.Errorf("_Solidity address must be 20 bytes")
	}

	return uint64(binary.BigEndian.Uint32(bytes[0:4])), binary.BigEndian.Uint64(bytes[4:12]), binary.BigEndian.Uint64(bytes[12:20]), nil
}

func _IdToSolidityAddress(shard uint64, realm uint64, num uint64) string {
	bytes := make([]byte, 20)
	binary.BigEndian.PutUint32(bytes[0:4], uint32(shard))
	binary.BigEndian.PutUint64(bytes[4:12], realm)
	binary.BigEndian.PutUint64(bytes[12:20], num)
	return hex.EncodeToString(bytes)
}

func _Has0xPrefix(str string) bool {
	return len(str) >= 2 && str[0] == '0' && (str[1] == 'x' || str[1] == 'X')
}

func _Without0x(s string) string {
	if _Has0xPrefix(s) {
		s = s[2:]
	}
	if len(s)%2 == 1 {
		s = "0" + s
	}
	return s
}

func _Hex2Bytes(str string) []byte {
	h, _ := hex.DecodeString(str)
	return h
}

func _IsHexCharacter(c byte) bool {
	return ('0' <= c && c <= '9') || ('a' <= c && c <= 'f') || ('A' <= c && c <= 'F')
}

func _IsHex(str string) bool {
	if len(str)%2 != 0 {
		return false
	}
	for _, c := range []byte(str) {
		if !_IsHexCharacter(c) {
			return false
		}
	}
	return true
}
