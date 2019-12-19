package hedera

import (
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSolidtySerialization(t *testing.T) {
	params := NewCallParams(nil).SetFunction("f")

	params, err := params.AddInt(uint32(16909060))
	if err != nil {
		println("AddInt uint32 paniced")
		panic(err)
	}

	params, err = params.AddBytesFixed([]byte{0, 1, 0, 0, 4, 0, 0, 0, 0, 8}, 10)
	if err != nil {
		println("AddBytesFixed paniced")
		panic(err)
	}

	params, err = params.AddInt(uint64(0xffffffff00000000))
	if err != nil {
		println("AddInt uint64 paniced")
		panic(err)
	}

	params.AddString("this is a grin: 😁")

	params, err = params.AddInt(uint16(1515))
	if err != nil {
		println("AddInt uint16 paniced")
		panic(err)
	}

	result := params.Finish()

	function, err := hex.DecodeString("f174a8d9")
	if err != nil {
		println("Paniced decoding hex string")
		panic(err)
	}

	first, err := hex.DecodeString("0000000000000000000000000000000000000000000000000000000001020304")
	if err != nil {
		println("Paniced decoding hex string")
		panic(err)
	}

	second, err := hex.DecodeString("0001000004000000000800000000000000000000000000000000000000000000")
	if err != nil {
		println("Paniced decoding hex string")
		panic(err)
	}

	third, err := hex.DecodeString("000000000000000000000000000000000000000000000000ffffffff00000000")
	if err != nil {
		println("Paniced decoding hex string")
		panic(err)
	}

	forthOffset, err := hex.DecodeString("00000000000000000000000000000000000000000000000000000000000000a0")
	if err != nil {
		println("Paniced decoding hex string")
		panic(err)
	}

	forthValueLength, err := hex.DecodeString("0000000000000000000000000000000000000000000000000000000000000014")
	if err != nil {
		println("Paniced decoding hex string")
		panic(err)
	}

	forthValue, err := hex.DecodeString("746869732069732061206772696e3a20f09f9881000000000000000000000000")
	if err != nil {
		println("Paniced decoding hex string")
		panic(err)
	}

	fifth, err := hex.DecodeString("00000000000000000000000000000000000000000000000000000000000005eb")
	if err != nil {
		println("Paniced decoding hex string")
		panic(err)
	}

	assert.Equal(t, 228, len(result), "Length of params does not match")
	assert.Equal(t, function, result[0:4], "Function signature doesn't match")
	assert.Equal(t, first, result[4:32+4], "First argument doesn't match")
	assert.Equal(t, second, result[32+4:(32*2)+4], "Second argument doesn't match")
	assert.Equal(t, third, result[(32*2)+4:(32*3)+4], "Third argument doesn't match")
	assert.Equal(t, forthOffset, result[(32*3)+4:(32*4)+4], "ForthOffset argument doesn't match")
	assert.Equal(t, fifth, result[(32*4)+4:(32*5)+4], "Fifth argument doesn't match")
	assert.Equal(t, forthValueLength, result[(32*5)+4:(32*6)+4], "ForthValueLength argument doesn't match")
	assert.Equal(t, forthValue, result[(32*6)+4:(32*7)+4], "ForthValue argument doesn't match")
}

func TestSolidtyStringArraySerialization(t *testing.T) {
	result := NewCallParams(nil).
		SetFunction("f").
		AddStringArray([]string{"one", "four"}).
		Finish()

	function, err := hex.DecodeString("e9cc8780")
	if err != nil {
		println("Paniced decoding hex string")
		panic(err)
	}

	firstOffset, err := hex.DecodeString("0000000000000000000000000000000000000000000000000000000000000020")
	if err != nil {
		println("Paniced decoding hex string")
		panic(err)
	}

	firstLengthOfArray, err := hex.DecodeString("0000000000000000000000000000000000000000000000000000000000000002")
	if err != nil {
		println("Paniced decoding hex string")
		panic(err)
	}

	firstElementOffset, err := hex.DecodeString("0000000000000000000000000000000000000000000000000000000000000040")
	if err != nil {
		println("Paniced decoding hex string")
		panic(err)
	}

	secondElementOffset, err := hex.DecodeString("0000000000000000000000000000000000000000000000000000000000000080")
	if err != nil {
		println("Paniced decoding hex string")
		panic(err)
	}

	firstElementLength, err := hex.DecodeString("0000000000000000000000000000000000000000000000000000000000000003")
	if err != nil {
		println("Paniced decoding hex string")
		panic(err)
	}

	secondElementLength, err := hex.DecodeString("0000000000000000000000000000000000000000000000000000000000000004")
	if err != nil {
		println("Paniced decoding hex string")
		panic(err)
	}

	firstElementValue, err := hex.DecodeString("6f6e650000000000000000000000000000000000000000000000000000000000")
	if err != nil {
		println("Paniced decoding hex string")
		panic(err)
	}

	secondElementValue, err := hex.DecodeString("666f757200000000000000000000000000000000000000000000000000000000")
	if err != nil {
		println("Paniced decoding hex string")
		panic(err)
	}

	assert.Equal(t, 260, len(result), "Length of params does not match")
	assert.Equal(t, function, result[0:4], "Function signature doesn't match")
	assert.Equal(t, firstOffset, result[4:32+4], "FirstOffset doesn't match")
	assert.Equal(t, firstLengthOfArray, result[32+4:(32*2)+4], "FirstLengthOfArray doesn't match")
	assert.Equal(t, firstElementOffset, result[(32*2)+4:(32*3)+4], "FirstElementOffset doesn't match")
	assert.Equal(t, secondElementOffset, result[(32*3)+4:(32*4)+4], "SecondElementOffset doesn't match")
	assert.Equal(t, firstElementLength, result[(32*4)+4:(32*5)+4], " firstElementLength doesn't match")
	assert.Equal(t, firstElementValue, result[(32*5)+4:(32*6)+4], "firstElementLength doesn't match")
	assert.Equal(t, secondElementLength, result[(32*6)+4:(32*7)+4], " secondElementLength doesn't match")
	assert.Equal(t, secondElementValue, result[(32*7)+4:(32*8)+4], "secondElementLength doesn't match")

}
