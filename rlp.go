package main

import (
	"bytes"
	"encoding/hex"
	"errors"
	"fmt"
)

// RLP encoding constants representing the prefixes for short strings and lists.
const (
	OffsetShortString = 0x80 // Offset for short strings
	OffsetLongString  = 0xb7 // Offset for long strings
	OffsetShortList   = 0xc0 // Offset for short lists
	OffsetLongList    = 0xf7 // Offset for long lists
)

// Decode takes an RLP-encoded byte slice and decodes it into an interface{}.
// The result can be a string, byte slice, or a list of decoded values.
func Decode(input []byte) (interface{}, error) {
	if len(input) == 0 {
		return nil, errors.New("input is empty") // Guard clause for empty input
	}
	return decodeRLP(input) // Start the recursive decoding process
}

// decodeRLP decodes an RLP byte slice based on its prefix type.
func decodeRLP(input []byte) (interface{}, error) {
	if len(input) == 0 {
		return nil, errors.New("input is empty") // Guard clause for empty input
	}

	prefix := input[0] // Read the first byte to determine the type

	// Handle small integers (0-127)
	if prefix < OffsetShortString {
		return uint64(input[0]), nil
	}

	// Handle short strings
	if prefix < OffsetLongString {
		length := int(prefix - OffsetShortString)
		if len(input) < 1+length {
			return nil, errors.New("input too short for short string")
		}
		return input[1 : 1+length], nil // Return the string slice
	}

	// Handle long strings
	if prefix < OffsetShortList {
		lengthLength := int(prefix - OffsetLongString)
		if len(input) < 1+lengthLength {
			return nil, errors.New("input too short for long string length")
		}
		length, err := parseLength(input[1 : 1+lengthLength])
		if err != nil {
			return nil, err
		}
		if len(input) < 1+lengthLength+length {
			return nil, errors.New("input too short for long string data")
		}
		return input[1+lengthLength : 1+lengthLength+length], nil
	}

	// Handle short lists
	if prefix < OffsetLongList {
		length := int(prefix - OffsetShortList)
		if len(input) < 1+length {
			return nil, errors.New("input too short for short list")
		}
		return decodeList(input[1 : 1+length]) // Decode list elements
	}

	// Handle long lists
	lengthLength := int(prefix - OffsetLongList)
	if len(input) < 1+lengthLength {
		return nil, errors.New("input too short for long list length")
	}
	length, err := parseLength(input[1 : 1+lengthLength])
	if err != nil {
		return nil, err
	}
	if len(input) < 1+lengthLength+length {
		return nil, errors.New("input too short for long list data")
	}
	return decodeList(input[1+lengthLength : 1+lengthLength+length]) // Decode list elements
}

// parseLength decodes the length of a byte slice in big-endian format.
func parseLength(data []byte) (int, error) {
	length := 0
	for _, b := range data {
		length = (length << 8) + int(b) // Convert byte to length
	}
	return length, nil
}

// decodeList recursively decodes a list of RLP-encoded elements.
func decodeList(input []byte) ([]interface{}, error) {
	var result []interface{}
	for len(input) > 0 {
		element, err := decodeRLP(input)
		if err != nil {
			return nil, err // Handle decoding errors
		}

		elemBytes, err := encodeRLP(element) // Encode element to calculate size
		if err != nil {
			return nil, err
		}
		result = append(result, element) // Append the decoded element to the result
		input = input[len(elemBytes):]   // Move to the next element
	}
	return result, nil
}

// encodeRLP encodes an element into its RLP-encoded format.
func encodeRLP(value interface{}) ([]byte, error) {
	switch v := value.(type) {
	case byte:
		if v < OffsetShortString {
			return []byte{v}, nil // Encode small byte directly
		}
	case uint64:
		return encodeRLP([]byte{byte(v)}) // Recursively encode uint64
	case []byte:
		length := len(v)
		// Encode short byte array
		if length == 1 && v[0] < OffsetShortString {
			return v, nil
		} else if length <= 55 {
			return append([]byte{byte(OffsetShortString + length)}, v...), nil
		}
		// Encode long byte array
		lengthBytes := encodeLength(length)
		return append(append([]byte{byte(OffsetLongString + len(lengthBytes))}, lengthBytes...), v...), nil
	case []interface{}:
		var buffer bytes.Buffer
		for _, elem := range v {
			encodedElem, err := encodeRLP(elem)
			if err != nil {
				return nil, err
			}
			buffer.Write(encodedElem) // Write each encoded element to the buffer
		}
		listBytes := buffer.Bytes()
		// Encode short list
		if len(listBytes) <= 55 {
			return append([]byte{byte(OffsetShortList + len(listBytes))}, listBytes...), nil
		}
		// Encode long list
		lengthBytes := encodeLength(len(listBytes))
		return append(append([]byte{byte(OffsetLongList + len(lengthBytes))}, lengthBytes...), listBytes...), nil
	}
	return nil, errors.New("unsupported RLP encoding type") // Handle unsupported types
}

// encodeLength converts an integer length to a big-endian byte slice format.
func encodeLength(length int) []byte {
	if length == 0 {
		return []byte{0} // Return zero-length byte
	}
	var result []byte
	for length > 0 {
		result = append([]byte{byte(length & 0xff)}, result...) // Build the byte slice in reverse order
		length >>= 8                                            // Shift right to get the next byte
	}
	return result
}

// toUint64 converts a byte slice or interface to uint64.
func toUint64(data interface{}) uint64 {
	bytes, ok := data.([]byte) // Check if data is a byte slice
	if !ok {
		return 0 // Return 0 if conversion fails
	}
	var result uint64
	for _, b := range bytes {
		result = (result << 8) + uint64(b) // Convert byte slice to uint64
	}
	return result
}

func main() {
	// Sample RLP-encoded Ethereum transaction (in hex string format)
	encodedTxHex := "f87282012a808085d1385c7bf082609094927e41ff8307835a1c081e0d7fd250625f2d4d0e8502540be4008406fdde03c080a0d0af512df4e0b4abdf5e3ff7a8b86aaf6bec465b9b89a2c3a4cd157e55c1afe1a0b1b2549650ff426312476c3f9e52f1a51242b9e726251430b0be13da81d94d4f"

	// Decode the hex string into a byte slice
	encodedTx, err := hex.DecodeString(encodedTxHex)
	if err != nil {
		fmt.Println("Failed to decode transaction hex:", err)
		return
	}

	// Decode the RLP data
	decoded, err := Decode(encodedTx)
	if err != nil {
		fmt.Println("Decoding error:", err)
		return
	}

	// Display the decoded data
	fmt.Println("Decoded Ethereum Transaction Fields:")
	decodedData, ok := decoded.([]interface{}) // Type assertion to a list of interface{}
	if ok {
		for i, field := range decodedData {
			fmt.Printf("Field %d: %v\n", i, field) // Print each decoded field
		}
	} else {
		fmt.Println("Unexpected decoded format.")
	}
}
