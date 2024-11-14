package btree

import (
	"fmt"
	"log"
)


func insert[T any](slice []T, idx int, val T) []T {
	if idx < 0 || idx > len(slice) {
		panic("Slice insert idx out of bounds")
	}
	if idx == cap(slice) {
		panic("Slice over capacity")
	}
	var dummy T
	slice = append(slice, dummy)
	copy(slice[idx+1:], slice[idx:])
	slice[idx] = val
	return slice
}


type KVP struct {
	key int
	val string
}



type btreeNode struct {
	isLeaf   bool
	order    int
	kvps     []*KVP
	children []*btreeNode
	parent   *btreeNode
}

func newNode(order int, isLeaf bool) *btreeNode {
	return &btreeNode{
		order:    order,
		kvps:     make([]*KVP, 0, order),
		children: make([]*btreeNode, 0, order+1),
		parent:   nil,
		isLeaf:   isLeaf,
	}
}

func (node *btreeNode) get(key int) (string, bool) {
	for i, curKVP := range node.kvps {
		if key == curKVP.key {
			return curKVP.val, true
		} else if key < curKVP.key {
			return node.children[i].get(key)
		}
	}
	if !node.isLeaf {
		return node.children[len(node.children)-1].get(key)
	}
	return "", false
}

func (node *btreeNode) traverse(res []*KVP) []*KVP {
	for i, curKVP := range node.kvps {
		if !node.isLeaf {
			// res = append(res, &KVP{-1, "recurse"})
			res = node.children[i].traverse(res)
		}
		res = append(res, curKVP)
	}
	if !node.isLeaf {
		res = node.children[len(node.children)-1].traverse(res)
	}
	return res
}

func (node *btreeNode) set(kvp *KVP) (bool, error) {
	if !node.isLeaf {
		for i, curKVP := range node.kvps {
			if kvp.key == curKVP.key {
				curKVP.val = kvp.val
				return false, nil
			} else if kvp.key < curKVP.key {
				return node.children[i].set(kvp)
			}
		}
		return node.children[len(node.children)-1].set(kvp)
	}

	var insertIdx int
	for i, curKVP := range node.kvps {
		insertIdx = i
		if kvp.key == curKVP.key {
			kvp.val = curKVP.val
			return false, nil
		} else if kvp.key < curKVP.key {
			break
		}
	}

	if l:=len(node.kvps); l>0 && kvp.key < node.kvps[l-1].key { // insert between keys
		node.kvps = insert(node.kvps, insertIdx, kvp)
	} else { // insert at end
		node.kvps = append(node.kvps, kvp)
	}

	// move up if len(keys) > order - 1
	err := node.split()

	return true, err
}

func (node *btreeNode) split() error {
	if len(node.kvps) < node.order {
		return nil
	}

	siblingNode := newNode(node.order, node.isLeaf)
	if node.parent == nil {
		node.parent = newNode(node.order, false)
		node.parent.children = append(node.parent.children, node)
	}
	siblingNode.parent = node.parent

	mid := len(node.kvps) / 2
	// fmt.Printf("new key %v\n", node.kvps[mid])

	for _, kvp := range node.kvps[mid+1:] {
		siblingNode.kvps = append(siblingNode.kvps, kvp)
	}
	node.kvps = node.kvps[:mid+1]
	if !node.isLeaf {
		for _, child := range node.children[mid+1:] {
			siblingNode.children = append(siblingNode.children, child)
			child.parent = siblingNode
		}
		node.children = node.children[:mid+1]
	}

	var insertIdx int
	for i, curKVP := range node.parent.kvps {
		insertIdx = i
		if node.kvps[mid].key == curKVP.key {
			log.Fatal("Duplicate keys")
		} else if node.kvps[mid].key < curKVP.key {
			// fmt.Printf("hello %v", len(node.parent.children))
			break
		}
	}

	if l:=len(node.parent.kvps);
		l>0 && node.kvps[mid].key < node.parent.kvps[l-1].key {
		// fmt.Printf("inserting %v to parent\n", node.kvps[mid].key)
		// fmt.Println("no")
		node.parent.kvps =
		insert(node.parent.kvps, insertIdx, node.kvps[mid])
		node.parent.children =
		insert(node.parent.children, insertIdx+1, siblingNode)
	} else {
		// fmt.Printf("appending %v to parent", node.kvps[mid].key)
		if len(node.parent.kvps) > 0 {
			// fmt.Printf(" with %v", node.parent.kvps[len(node.parent.kvps)-1].key)
		}
		node.parent.kvps = append(node.parent.kvps, node.kvps[mid])
		node.parent.children = append(node.parent.children, siblingNode)
		// fmt.Println(node.parent == siblingNode.parent)
		// fmt.Println(siblingNode.kvps[0].key)
		// fmt.Println(siblingNode.children)
		// node.parent.printKeys()
		// node.printKeys()
		// node.parent.children = append(node.parent.children, node)
	}
	node.kvps = node.kvps[:mid]


	return node.parent.split()
}

func (node *btreeNode) printKeys() {
	fmt.Printf("Keys: ")
	for _, kvp := range node.kvps {
		fmt.Printf("%v ", kvp.key)
	}
	fmt.Println()
}

func (node *btreeNode) verifyKeys(other *btreeNode) {
	if len(node.kvps) == 0 || len(other.kvps) == 0 {
		return
	}
	for i := range node.kvps {
		if i > 0 && node.kvps[i-1].key > node.kvps[i].key {
			panic("node keys out of order")
		}
	}
	for i := range other.kvps {
		if i > 0 && other.kvps[i-1].key > other.kvps[i].key {
			panic("other keys out of order")
		}
	}
	if node.kvps[len(node.kvps)-1].key > other.kvps[len(other.kvps)-1].key {
		panic("nodes out of order")
	}
}


