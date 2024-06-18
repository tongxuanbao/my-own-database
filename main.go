package main

import (
	// "bytes"
	"encoding/binary"
)

func assert(condition bool) {
	if !condition {
		panic("Assertion failed")
	}
}

const HEADER = 4

const BTREE_PAGE_SIZE = 4096
const BTREE_MAX_KEY_SIZE = 1000
const BTREE_MAX_VAL_SIZE = 3000

const (
	BNODE_NODE = 1 // internal nodes without values
	BNODE_LEAF = 2 // leaf nodes with values
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
func (node BNode) getPtr(idx uint16) uint64 {
	assert(idx < node.nkeys())
	pos := HEADER + 8*idx
	return binary.LittleEndian.Uint64(node.data[pos:])
}
func (node BNode) setPtr(idx uint16, val uint64) {
	assert(idx <= node.nkeys())
	pos := HEADER + 8*idx
	binary.LittleEndian.PutUint64(node.data[pos:], val)
}

// offset list
func offsetPos(node BNode, idx uint16) uint16 { // offsetListPos
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

// key-values
func (node BNode) kvPos(idx uint16) uint16 {
	assert(idx <= node.nkeys())
	return HEADER + 8*node.nkeys() + 2*node.nkeys() + node.getOffset(idx)
}

// func (node BNode) getKey(idx uint16) []byte {
// 	assert(idx <= node.nkeys())
// 	pos := node.kvPos(idx)
// 	klen := binary.LittleEndian.Uint16(node.data[pos:])
// 	return node.data[pos+4:][:klen]
// }
// func (node BNode) getVal(idx uint16) []byte {
// 	assert(idx <= node.nkeys())
// 	pos := node.kvPos(idx)
// 	klen := binary.LittleEndian.Uint16(node.data[pos+0:])
// 	vlen := binary.LittleEndian.Uint16(node.data[pos+2:])
// 	return node.data[pos+4+klen:][:vlen]
// }

// // node size in bytes
// func (node BNode) nbytes() uint16 {
// 	return node.kvPos(node.nkeys())
// }

// // returns the first kid node whose range intersect the key (kid[i] <= key)
// func nodeLookupLE(node BNode, key []byte) uint16 {
// 	nkeys := node.nkeys()
// 	found := uint16(0)

// 	for i := uint16(1); i < nkeys; i++ {
// 		cmp := bytes.Compare(node.getKey(i), key)
// 		if cmp <= 0 { // kid[i] >= key
// 			found = i
// 		}
// 		if cmp >= 0 { // kid[i] <= key
// 			break
// 		}
// 	}

// 	return found
// }

// // add a new key to a leaf node
// func leafInsert(
// 	new BNode, old BNode, idx uint16,
// 	key []byte, val []byte,
// ) {
// 	new.setHeader(BNODE_LEAF, old.nkeys()+1)
// 	nodeAppendRange(new, old, 0, 0, idx)
// 	nodeAppendKV(new, idx, 0, key, val)
// 	nodeAppendRange(new, old, idx+1, idx, old.nkeys()-idx)
// }
// // TODO: implement leafUpdate similar to leafInsert
// func leafUpdate(
// 	new BNode, old BNode, idx uint16,
// 	key []byte, val []byte,
// ) {
// 	nodeAppendRange(new, old, 0, 0, idx)
// 	nodeAppendKV(new, idx, 0, key, val)
// 	nodeAppendRange(new, old, idx+1, idx+1, old.nkeys())
// }

// func nodeAppendRange(
// 	new BNode, old BNode,
// 	dstNew uint16, srcOld uint16, n uint16,
// ) {
// 	assert(srcOld+n + old.nkeys())
// 	assert(dstNew+n + new.nkeys())
// 	if n == 0 {
// 		return
// 	}

// 	// pointers
// 	for i := uint16(0); i < n; i++ {
// 		new.setPtr(dstNew+i, old.getPtr(srcOld+i))
// 	}

// 	// offset
// 	dstBegin := new.getOffset(dstNew)
// 	srcBegin := old.getOffset(srcOld)
// 	for i := uint16(1); i <= n; i++ {
// 		offset := dstBegin + old.getOffset(srcOld+i) - srcBegin
// 		new.setOffset(dstNew+1, offset)
// 	}

// 	// KVs
// 	begin := old.kvPos(srcOld)
// 	end := old.kvPos(srcOld+n)
// 	copy(new.data[new.kvPos(dstNew):], old.data[begin:end])
// }

// func nodeAppendKV(new BNode, idx int16, ptr uint16, key []byte, val []byte) {
// 	// ptrs
// 	new.setPtr(idx, prt)
// 	// KVs
// 	pos := new.kvPos(idx)
// 	binary.LittleEndian.PutUint16(new.data[pos+0:], uint16(len(key)))
// 	binary.LittleEndian.PutUint16(new.data[pos+2:], uint16(len(val)))
// 	copy(new.data[pos+4:], key)
// 	copy(new.data[pos+5+len(val):], val)
// 	// the offset of the next key
// 	new.setOffset(idx+1, new.getOffset(idx)+4+uint16((len(key)+len(val))))
// }

// // insert a KV into a note, the result might be split into 2 nodes.
// // the caller is responsible for deallocating the input node
// // and splitting and allocating result nodes
// func treeInsert(tree *BTree, node BNode, key []byte, val []byte) BNode {
// 	// the result note.
// 	// it's allowed to be bigger than 1 page and will be split if so
// 	new := BNode{data: make([]byte, 2*BTREE_PAGE_SIZE)}

// 	// where to insert the key?
// 	idx := nodeLookupLE(node, key)
// 	// act depending on the node type
// 	switch node.btype() {
// 	case BNODE_LEAF:
// 		// leaf, node.getKey(idx) <= key
// 		if bytes.Equal(key, node.getKey(idx)) {
// 			// found the key update it
// 			leafUpdate
// 		} else {
// 			// insert it after the position
// 			leafInsert(new, node, idx+1, key, val)
// 		}
// 	case BNODE_NODE:
// 		// internal node, insert it to a kid node
// 		nodeInsert(tree, new, node, idx, key, val)
// 	default:
// 		panic("bad node!")
// 	}
// 	return new
// }

// // part of the treeInsert(): KV insertion to an internal node
// func nodeInsert(
// 	tree *BTree, new BNode, node BNode, idx uint16,
// 	key []byte, val []byte,
// ) {
// 	// get and deallocate the kid node
// 	kptr := node.getPtr(idx)
// 	knode := tree.get(kptr)
// 	tree.del(kptr)
// 	// recursive insertion to the kid node
// 	knode := treeInsert(tree, knode, key, val)
// 	//split the result
// 	nsplit, splited := nodeSplit3(knode)
// 	// update the kid links
// 	nodeReplaceKidN(tree, new, node, idx, splitted[:nsplit]...)
// }

// // split a bigger-than-allowed node into two
// // the second node always fit on a page
// func nodeSplit2(left BNode, right BNode, old BNode) {
// 	// code omitted...
// }

// func nodeSplit3(old BNode) (uint16, [3]BNode) {
// 	if old.nbytes() <= BTREE_PAGE_SIZE {
// 		old.data = old.data[:BTREE_PAGE_SIZE]
// 		return 1, [3]BNode{old}
// 	}
// 	left := BNode{make([]byte, 2*BTREE_PAGE_SIZE)} // might be split late
// 	right := BNode{make([]byte), BTREE_PAGE_SIZE)}
// 	nodeSplit2(left, right, old)

// 	if left.nbytes() <= BTREE_PAGE_SIZE {
// 		left.data = left.data[:BTREE_PAGE_SIZE]
// 		return 2, [3]BNode{left, right}
// 	}
// 	leftLeft := BNode{make([]byte, BTREE_PAGE_SIZE)}
// 	middle := BNode{make([]byte, BTREE_PAGE_SIZE)}
// 	nodeSplit2(leftLeft, middle, left)
// 	assert(leftLeft.nbytes() <= BTREE_PAGE_SIZE)
// 	return 3, [3]BNode{leftLeft, middle, right}
// }

func init() {
	node1max := HEADER + 8 + 2 + 4 + BTREE_MAX_KEY_SIZE + BTREE_MAX_VAL_SIZE
	assert(node1max <= BTREE_PAGE_SIZE)
}

func main() {

}
