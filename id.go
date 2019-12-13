package hedera

import (
	"fmt"
	"strconv"
	"strings"
)

func idFromString(s string) (shard int, realm int, num int, err error) {
	values := strings.SplitN(s, ".", 3)
	if len(values) != 3 {
		// Was not three values separated by periods
		return 0, 0, 0, fmt.Errorf("expected {shard}.{realm}.{num}")
	}

	shard, err = strconv.Atoi(values[0])
	if err != nil {
		return 0, 0, 0, err
	}

	realm, err = strconv.Atoi(values[1])
	if err != nil {
		return 0, 0, 0, err
	}

	num, err = strconv.Atoi(values[2])
	if err != nil {
		return 0, 0, 0, err
	}

	return
}
