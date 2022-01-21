package hedera

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"
)

func _AccountIDFromString(s string) (shard int, realm int, num int, checksum *string, alias *PublicKey, err error) {
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

	key, err := PublicKeyFromString(values[2])
	if err != nil {
		num, err = strconv.Atoi(values[2])
		if err != nil {
			return 0, 0, 0, nil, nil, err
		}

		return shard, realm, num, checksum, nil, nil
	}

	return shard, realm, -1, checksum, &key, nil
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
			return 0, 0, 0, nil, []byte{}, err
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
