package main

import (
	"encoding/binary"
	"testing"
)

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
}

func TestBNodeNBytes(t *testing.T) {
	// Create a sample BNode instance
	nodeData := []byte{1, 0} // little endian
	node := BNode{data: nodeData}

	expectedBType := uint16(1)
	actualBType := node.btype()

	if actualBType != expectedBType {
		t.Errorf("Expected btype to be %d, but got %d", expectedBType, actualBType)
	}
}

func TestBNodeNKeys(t *testing.T) {
	nodeData := []byte{1, 0, 2, 0}
	node := BNode{data: nodeData}

	expectedNKeys := uint16(2)
	actualNKeys := node.nkeys()

	if actualNKeys != expectedNKeys {
		t.Errorf("Expected btype to be %d, but got %d", expectedNKeys, actualNKeys)
	}
}

func TestBNodeSetHeader(t *testing.T) {
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

func TestBNodeGetPtr(t *testing.T) {
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

func TestBNodeSetPtr(t *testing.T) {
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

func TestBNodeOffsetPos(t *testing.T) {
	
}
