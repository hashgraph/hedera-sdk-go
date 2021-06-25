package hedera

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"
)

func idFromString(s string) (shard int, realm int, num int, checksum *string, err error) {
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

	return
}

func idFromSolidityAddress(s string) (uint64, uint64, uint64, error) {
	bytes, err := hex.DecodeString(s)
	if err != nil {
		return 0, 0, 0, err
	}

	if len(bytes) != 20 {
		return 0, 0, 0, fmt.Errorf("Solidity address must be 20 bytes")
	}

	return uint64(binary.BigEndian.Uint32(bytes[0:4])), binary.BigEndian.Uint64(bytes[4:12]), binary.BigEndian.Uint64(bytes[12:20]), nil
}

func idToSolidityAddress(shard uint64, realm uint64, num uint64) string {
	bytes := make([]byte, 20)
	binary.BigEndian.PutUint32(bytes[0:4], uint32(shard))
	binary.BigEndian.PutUint64(bytes[4:12], realm)
	binary.BigEndian.PutUint64(bytes[12:20], num)
	return hex.EncodeToString(bytes)
}
