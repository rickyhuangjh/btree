package btree

import (
	"fmt"
)

type BTree interface {
	Set(kvp *KVP) error
	// Delete(key int) error
	Get(key int) (string, bool)
	Traverse() []*KVP
	Size() int
	// Height() int
}

type btree struct {
	order int // max children
	root  *btreeNode
    size int
}

func NewBTree(order int) BTree {
	return &btree{
		order: order,
		root:  nil,
	}
}

func (btree *btree) Size() int {
	return btree.size
}

func (btree *btree) Get(key int) (string, bool) {
	return btree.root.get(key)
}

func (btree *btree) Traverse() []*KVP {
	if btree.root == nil {
		return nil
	}
	btree.root.printKeys()
	fmt.Println(btree.root.children)
    res := btree.root.traverse(make([]*KVP, 0, btree.size))
	return res
}

func (btree *btree) Set(kvp *KVP) error {
	if btree.root == nil {
		btree.root = newNode(btree.order, true)
	}

    wasAdded, err := btree.root.set(kvp)
    if btree.root.parent != nil {
        btree.root = btree.root.parent
    }
    if wasAdded {
        btree.size++
    }
	 //fmt.Println(a)
    return err
}
