//go:build all || unit
// +build all unit

package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRLPItemEncodeInteger(t *testing.T) {
	t.Parallel()

	item := NewRLPItem(VALUE_TYPE)
	item.AssignBytes([]byte{42}) // Assign a byte slice

	encoded, err := item.Write() // Ensure this captures both return values
	require.NoError(t, err)
	expected := []byte{42}
	assert.Equal(t, expected, encoded)
}

func TestRLPItemEncodeString(t *testing.T) {
	t.Parallel()

	item := NewRLPItem(VALUE_TYPE)
	item.AssignString("hello") // Use the AssignString method

	encoded, err := item.Write() // Ensure this captures both return values
	require.NoError(t, err)
	expected := []byte{0x85, 'h', 'e', 'l', 'l', 'o'} // 0x85 is the prefix for short string of length 5
	assert.Equal(t, expected, encoded)
}

func TestRLPItemEncodeByteSlice(t *testing.T) {
	t.Parallel()

	item := NewRLPItem(VALUE_TYPE)
	item.AssignBytes([]byte{0x01, 0x02, 0x03}) // Assign a byte slice

	encoded, err := item.Write() // Ensure this captures both return values
	require.NoError(t, err)
	expected := []byte{0x83, 0x01, 0x02, 0x03} // 0x83 is the prefix for short byte slice of length 3
	assert.Equal(t, expected, encoded)
}

func TestRLPItemEncodeList(t *testing.T) {
	t.Parallel()

	listItem := NewRLPItem(LIST_TYPE)
	listItem.PushBack(NewRLPItem(VALUE_TYPE))     // Push back the first item
	listItem.childItems[0].AssignBytes([]byte{1}) // Set the value of the first item

	listItem.PushBack(NewRLPItem(VALUE_TYPE))  // Push back the second item
	listItem.childItems[1].AssignString("two") // Set the value of the second item

	listItem.PushBack(NewRLPItem(VALUE_TYPE))              // Push back the third item
	listItem.childItems[2].AssignBytes([]byte{0x03, 0x04}) // Set the value of the third item

	encoded, err := listItem.Write() // Ensure this captures both return values
	require.NoError(t, err)
	expected := []byte{0xc8, 0x01, 0x83, 't', 'w', 'o', 0x82, 0x03, 0x04}
	assert.Equal(t, expected, encoded)
}

func TestRLPItemDecodeInteger(t *testing.T) {
	t.Parallel()

	item := NewRLPItem(VALUE_TYPE)
	err := item.Read([]byte{42}) // Use the Read method for decoding
	require.NoError(t, err)
	assert.Equal(t, []byte{42}, item.itemValue)
}

func TestRLPItemDecodeString(t *testing.T) {
	t.Parallel()

	item := NewRLPItem(VALUE_TYPE)
	err := item.Read([]byte{0x85, 'h', 'e', 'l', 'l', 'o'})
	require.NoError(t, err)
	assert.Equal(t, []byte("hello"), item.itemValue)
}

func TestRLPItemDecodeByteSlice(t *testing.T) {
	t.Parallel()

	item := NewRLPItem(VALUE_TYPE)
	err := item.Read([]byte{0x83, 0x01, 0x02, 0x03})
	require.NoError(t, err)
	assert.Equal(t, []byte{0x01, 0x02, 0x03}, item.itemValue)
}

func TestRLPItemDecodeList(t *testing.T) {
	t.Parallel()

	item := NewRLPItem(LIST_TYPE)
	err := item.Read([]byte{0xc8, 0x83, 'c', 'a', 't', 0x83, 'd', 'o', 'g'})
	require.NoError(t, err)

	// The expected list contains 3 items
	expectedList := []*RLPItem{
		NewRLPItem(VALUE_TYPE), // First item (the string "cat")
		NewRLPItem(VALUE_TYPE), // Second item (the string "dog")
	}

	expectedList[0].AssignBytes([]byte{'c', 'a', 't'}) // Set the value for the first item as the byte representation of "cat"
	expectedList[1].AssignBytes([]byte{'d', 'o', 'g'}) // Set the value for the second item as the byte representation of "dog"

	// Compare the expected structure to the decoded structure
	require.Equal(t, len(expectedList), len(item.childItems)) // Check length first
	for i, expectedItem := range expectedList {
		assert.Equal(t, expectedItem.itemValue, item.childItems[i].itemValue) // Compare item values
	}
}
