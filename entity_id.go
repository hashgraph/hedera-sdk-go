package hedera

import (
	"fmt"
	"github.com/pkg/errors"
	"regexp"
	"strconv"
)

// EntityID is an interface for various IDs of entities (Account, Contract, File, etc)
type EntityID interface {
	isEntityID()
}

type parseAddressResult struct {
	status             int
	num1               int64
	num2               int64
	num3               int64
	correctChecksum    string
	givenChecksum      string
	noChecksumFormat   string
	withChecksumFormat string
}

func checksumVerify(num int) error {
	switch num {
	case 0:
		return errors.New("Invalid ID: format should look like 0.0.123 or 0.0.123-laujm")
	case 1:
		return errors.New("Invalid ID: checksum does not match")
	case 2:
		return nil
	case 3:
		return nil
	default:
		return errors.New("Unrecognized status")
	}
}

func checksumParseAddress(ledgerID string, address string) (parseAddressResult, error) {
	match, err := regexp.Compile(`(0|(?:[1-9]\d*))\.(0|(?:[1-9]\d*))\.(0|(?:[1-9]\d*))(?:-([a-z]{5}))?$`)
	if err != nil {
		return parseAddressResult{status: 0}, err
	}

	matchArray := match.FindStringSubmatch(address)

	a := make([]int64, len(matchArray))
	for i := 1; i < len(matchArray)-1; i++ {
		a[i], err = strconv.ParseInt(matchArray[i], 10, 64)
		if err != nil {
			return parseAddressResult{status: 0}, err
		}
	}

	ad := fmt.Sprintf("%s.%s.%s", matchArray[1], matchArray[2], matchArray[3])

	checksum := checkChecksum(ledgerID, ad)

	var status int
	if matchArray[4] == "" {
		status = 2
	} else if checksum == matchArray[4] {
		status = 3
	} else {
		status = 1
	}
	return parseAddressResult{
		status:             status,
		num1:               a[1],
		num2:               a[2],
		num3:               a[3],
		correctChecksum:    checksum,
		givenChecksum:      matchArray[4],
		noChecksumFormat:   ad,
		withChecksumFormat: ad + "(" + checksum + ")",
	}, nil
}

func checkChecksum(ledgerID string, address string) string {
	answer := ""
	digits := make([]int, 0)
	s0 := 0
	s1 := 0
	s := 0
	sh := 0
	checksum := 0
	n := len(address)
	p3 := 26 * 26 * 26
	p5 := 26 * 26 * 26 * 26 * 26
	m := 1000003
	asciiA := []rune("a")[0]
	w := 31

	id := ledgerID + "000000000000"
	h := make([]int64, 0)

	for i := 0; i < len(id); i += 2 {
		processed, _ := strconv.ParseInt(id[i:i+2], 16, 64)
		h = append(h, processed)
		if i+3 == len(id) {
			processed, _ = strconv.ParseInt(id[i:len(id)-1], 16, 64)
			h = append(h, processed)
			break
		}
	}

	for _, j := range address {
		if string(j) == "." {
			digits = append(digits, 10)
		} else {
			processed, _ := strconv.Atoi(string(j))
			digits = append(digits, processed)
		}
	}

	for i := 0; i < len(digits); i++ {
		s = (w*s + digits[i]) % p3
		if i%2 == 0 {
			s0 = (s0 + digits[i]) % 11
		} else {
			s1 = (s1 + digits[i]) % 11
		}
	}

	for i := 0; i < len(h); i++ {
		sh = (w*sh + int(h[i])) % p5
	}

	checksum = ((((n%5)*11+s0)*11+s1)*p3 + s + sh) % p5
	checksum = (checksum * m) % p5

	for i := 0; i < 5; i++ {
		answer = string(asciiA+rune(checksum%26)) + answer
		checksum /= 26
	}

	return answer
}

func (id AccountID) isEntityID() {}

// func (id FileID) isEntityID()     {}
// func (id ContractID) isEntityID() {}
