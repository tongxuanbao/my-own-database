package main

import (
	"encoding/binary"
	"testing"
)

/*
func createTestBNodeInternal() BNode {
	// Create a BNode with sufficient data slice
	node := BNode{data: make([]byte, 30)} // Adjust size as needed

	// Set up the data

	// Headers
	binary.LittleEndian.PutUint16(node.data[0:], 1) // type = 1 (internal node)
	binary.LittleEndian.PutUint16(node.data[2:4], 2) // nkeys = 2
	// Pointers
	binary.LittleEndian.PutUint64(node.data[HEADER:], 0x1111) // pointer1
	binary.LittleEndian.PutUint64(node.data[HEADER+8:], 0x2222) // pointer2

	// Offsets
	binary.LittleEndian.PutUint16(node.data[HEADER+16:], 160) // offset 1
	binary.LittleEndian.PutUint16(node.data[HEADER+18:], 160) // offset 2

	// key-value 1
	binary.LittleEndian.PutUint16(node.data[HEADER+20:], 4) // klen
	binary.LittleEndian.PutUint16(node.data[HEADER+20:], 4) // vlen
	binary.LittleEndian.PutUint64(node.data[HEADER+20:], 0x11111111) // key
	binary.LittleEndian.PutUint16(node.data[HEADER+20:], 0x22222222) // val

	// key-value 2


	return node
}*/

// | type | nkeys | pointers    | offsets     | key-values
// | 2B   | 2B    | nkeys(3)*8B | nkeys * 2B  | ...
// | 2    |  3    | none        | 0 | 14 | 30 | ...
// This is the format of the KV pair. Lengths followed by data.
// | klen | vlen | key  | val    |
// | 2B   | 2B   | ...  | ...    |
// | 4    | 6    | key1 | value1 |
func createMockedLeafNode() BNode {
	node := BNode{
		data: make([]byte, BTREE_PAGE_SIZE),
	}
	node.setHeader(BNODE_LEAF, 3)

	// Pointers are not used in leaf nodes, so we skip the first 8*nkeys bytes.

	// Offsets
	offsets := []uint16{0, 14, 30}
	for i, offset := range offsets {
		node.setOffset(uint16(i+1), offset)
	}

	// Key-Value pairs
	// klen, vlen, key, value
	kvPairs := [][]byte{
		//            |--------|-this is klen 2 bytes
		//            |        |  |--------|-this is vlen 2 bytes
		append([]byte{0x04, 0x00, 0x06, 0x00}, []byte("key1value1")...),
		append([]byte{0x04, 0x00, 0x06, 0x00}, []byte("key2value2")...),
		append([]byte{0x04, 0x00, 0x06, 0x00}, []byte("key3value3")...),
	}

	pos := HEADER + 8*node.nkeys() + 2*node.nkeys()
	for _, kv := range kvPairs {
		copy(node.data[pos:], kv)
		pos += uint16(len(kv))
	}

	return node
}

// | type | nkeys | pointers    | offsets     | key-values
// | 2B   | 2B    | nkeys(3)*8B | nkeys * 2B  | ...
// | 2    |  3    | none        | 0 | 14 | 30 | ...
// This is the format of the KV pair. Lengths followed by data.
// | klen | vlen | key  | val    |
// | 2B   | 2B   | ...  | ...    |
// | 4    | 6    | key1 | value1 |
func createMockedInternalNode() BNode {
	node := BNode{
		data: make([]byte, BTREE_PAGE_SIZE),
	}
	node.setHeader(BNODE_NODE, 3)

	// Pointers
	pointers := []uint64{1, 2, 3}
	for i, ptr := range pointers {
		node.setPtr(uint16(i), ptr)
	}

	// Offsets
	offsets := []uint16{0, 10}
	for i, offset := range offsets {
		node.setOffset(uint16(i+1), offset)
	}

	// Key-Value pairs
	// klen, vlen, key, value
	kvPairs := [][]byte{
		append([]byte{0x04, 0x00, 0x00, 0x00}, []byte("key1")...),
		append([]byte{0x04, 0x00, 0x00, 0x00}, []byte("key2")...),
	}

	pos := HEADER + 8*node.nkeys() + 2*node.nkeys()
	for _, kv := range kvPairs {
		copy(node.data[pos:], kv)
		pos += uint16(len(kv))
	}

	return node
}

func TestBNode_NBytes(t *testing.T) {
	// Create a sample BNode instance
	node := createMockedLeafNode()

	expectedBType := uint16(2)
	actualBType := node.btype()

	if actualBType != expectedBType {
		t.Errorf("Expected btype to be %d, but got %d", expectedBType, actualBType)
	}
}

func TestBNode_NKeys(t *testing.T) {
	node := createMockedLeafNode()

	expectedNKeys := uint16(3)
	actualNKeys := node.nkeys()

	if actualNKeys != expectedNKeys {
		t.Errorf("Expected btype to be %d, but got %d", expectedNKeys, actualNKeys)
	}
}

func TestBNode_SetHeader(t *testing.T) {
	nodeData := []byte{1, 0, 2, 0}
	node := BNode{data: nodeData}

	btype := uint16(2)
	nkeys := uint16(3)

	node.setHeader(btype, nkeys)

	actualBType := node.btype()

	if actualBType != btype {
		t.Errorf("Expected btype to be %d, but got %d", btype, actualBType)
	}

	actualNKeys := node.nkeys()

	if actualNKeys != nkeys {
		t.Errorf("Expected btype to be %d, but got %d", nkeys, actualNKeys)
	}
}

func TestBNode_GetPtr(t *testing.T) {
	// Create a BNode with sufficient data slice
	node := BNode{data: make([]byte, 30)} // Adjust size as needed

	// Set up the data
	binary.LittleEndian.PutUint16(node.data[2:4], 3) // nkeys = 3
	binary.LittleEndian.PutUint64(node.data[HEADER:], 0x1111)
	binary.LittleEndian.PutUint64(node.data[HEADER+8:], 0x2222)
	binary.LittleEndian.PutUint64(node.data[HEADER+16:], 0x3333)

	// Define test cases
	tests := []struct {
		idx  uint16
		want uint64
	}{
		{0, 0x1111},
		{1, 0x2222},
		{2, 0x3333},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			got := node.getPtr(tt.idx)
			if got != tt.want {
				t.Errorf("getPtr(%d) = 0x%x; want 0x%x", tt.idx, got, tt.want)
			}
		})
	}
}

func TestBNode_SetPtr(t *testing.T) {
	// Create a BNode with sufficient data slice
	node := BNode{data: make([]byte, 28)} // Adjust size as needed

	// Set up the data
	binary.LittleEndian.PutUint16(node.data[2:4], 3) // nkeys = 3

	// Define test cases
	tests := []struct {
		idx  uint16
		val  uint64
		want uint64
	}{
		{0, 0x1111, 0x1111},
		{1, 0x2222_2222_2222_2222, 0x2222_2222_2222_2222},
		{2, 0x3333_3333_3333_3333, 0x3333_3333_3333_3333},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			node.setPtr(tt.idx, tt.val)
			pos := HEADER + 8*tt.idx
			got := binary.LittleEndian.Uint64(node.data[pos:])
			if got != tt.want {
				t.Errorf("setPtr(%d, 0x%x) = 0x%x; want 0x%x", tt.idx, tt.val, got, tt.want)
			}
		})
	}
}

func TestBNode_OffsetPos(t *testing.T) {
	node := createMockedLeafNode()

	expectedOffsetPos := uint16(28)
	actualOffsetPos := offsetPos(node, 1)

	if expectedOffsetPos != actualOffsetPos {
		t.Errorf("Expected btype to be %d, but got %d", expectedOffsetPos, actualOffsetPos)
	}
}

func TestBNode_getOffset(t *testing.T) {
	node := createMockedLeafNode()

	tests := []struct {
		idx  uint16
		want uint16
	}{
		{1, uint16(0)},
		{2, uint16(14)},
		{3, uint16(30)},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			got := node.getOffset(tt.idx)
			if got != tt.want {
				t.Errorf("getOffset(%d) = %d, want %d", tt.idx, got, tt.want)
			}
		})
	}
}

func TestBNode_setOffset(t *testing.T) {
	node := createMockedLeafNode()

	node.setOffset(2, uint16(16))

	want := uint16(16)
	got := node.getOffset(2)

	if got != want {
		t.Errorf("setOffset(2, 16) = %d, want %d", got, want)
	}
}

// func TestBNode_kvPos(t *testing.T) {
// 	node := createMockedLeafNode()

// }
