package main

import "encoding/binary"

func assert(condition bool) {
	if !condition {
		panic("Assertion failed")
	}
}

const (
	BNODE_NODE = 1 // internal nodes without values
	BNODE_LEAF = 1 // leaf nodes with values
)

type BNode struct {
	data []byte
}

type BTree struct {
	// pointer (a nonzero page number)
	root uint64
	// callbacks for managing on-disk pages
	get func(uint64) BNode // dereference a pointer
	new func(BNode) uint64 // allocate a new page
	del func(uint64)       // deallocate a page
}

// header
func (node BNode) btype() uint16 {
	return binary.LittleEndian.Uint16(node.data)
}
func (node BNode) nkeys() uint16 {
	return binary.LittleEndian.Uint16(node.data[2:4])
}
func (node BNode) setHeader(btype uint16, nkeys uint16) {
	binary.LittleEndian.PutUint16(node.data[0:2], btype)
	binary.LittleEndian.PutUint16(node.data[2:4], nkeys)
}

// pointer
func (node BNode) getPtr(idx uint16) uint16 {
	assert(idx < node.nkeys())
	pos := HEADER + 8*idx
	return binary.LittleEndian.Uint16(node.data[pos:])
}
func (node BNode) setPtr(idx uint16, val uint16) {
	assert(idx < node.nkeys())
	pos := HEADER + 8*idx
	binary.LittleEndian.PutUint16(node.data[pos:], val)
}

// offset list
func offsetPos(node BNode, idx uint16) uint16 {
	assert(1 <= idx && idx <= node.nkeys())
	return HEADER + 8*node.nkeys() + 2*(idx-1)
}
func (node BNode) getOffset(idx uint16) uint16 {
	if idx == 0 {
		return 0
	}
	return binary.LittleEndian.Uint16(node.data[offsetPos(node, idx):])
}
func (node BNode) setOffset(idx uint16, offset uint16) {
	binary.LittleEndian.PutUint16(node.data[offsetPos(node, idx):], offset)
}

const HEADER = 4

const BTREE_PAGE_SIZE = 4096
const BTREE_MAX_KEY_SIZE = 1000
const BTREE_MAX_VAL_SIZE = 3000

func init() {
	node1max := HEADER + 8 + 2 + 4 + BTREE_MAX_KEY_SIZE + BTREE_MAX_VAL_SIZE
	assert(node1max <= BTREE_PAGE_SIZE)
}

func main() {

}
